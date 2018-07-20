// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package drbd

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/LINBIT/godrbdutils"
	log "github.com/golang/glog"
	"github.com/opensds/opensds/contrib/drivers/utils/config"
	pb "github.com/opensds/opensds/pkg/dock/proto"
	"github.com/opensds/opensds/pkg/model"
)

// ReplicationDriver
type ReplicationDriver struct{}

// Setup
func (r *ReplicationDriver) Setup() error { return nil }

// Unset
func (r *ReplicationDriver) Unset() error { return nil }

// CreateReplication
func (r *ReplicationDriver) CreateReplication(opt *pb.CreateReplicationOpts) (*model.ReplicationSpec, error) {
	log.Infof("DRBD create replication ....")

	conf := drbdConf{}
	_, err := config.Parse(&conf, defaultConfPath)
	if err != nil {
		return nil, err
	}
	if len(conf.Hosts) != 2 {
		return nil, fmt.Errorf("Your configuration does not contain exactly 2 hosts")
	}

	isPrimary := opt.GetIsPrimary()
	primaryData := opt.GetPrimaryReplicationDriverData()
	secondaryData := opt.GetSecondaryReplicationDriverData()

	var myData *map[string]string
	if isPrimary {
		myData = &primaryData
	} else {
		myData = &secondaryData
	}

	var myHostName string
	var myHostIP string
	var hok bool
	if myHostName, hok = (*myData)["HostName"]; hok {
		myHostIP, hok = (*myData)["HostIp"]
	}
	if !hok {
		return nil, fmt.Errorf("Data did not contain 'HostIp' or 'HostName' key")
	}

	var myHost, peerHost godrbdutils.Host
	for _, h := range conf.Hosts {
		if h.Name == myHostName && h.IP == myHostIP {
			myHost = h
		} else {
			peerHost = h
		}
	}

	if myHost.Name == "" || myHost.IP == "" || peerHost.Name == "" || peerHost.IP == "" ||
		myHost.Name == peerHost.Name || myHost.ID == peerHost.ID || myHost.IP == peerHost.IP {
		return nil, fmt.Errorf("Could not find valid hosts")
	}

	resName := opt.GetId()

	primaryVolID := opt.GetPrimaryVolumeId()
	secondaryVolID := opt.GetSecondaryVolumeId()
	path, _ := filepath.EvalSymlinks(primaryData["Mountpoint"])
	primaryBackingDevice, _ := filepath.Abs(path)
	path, _ = filepath.EvalSymlinks(secondaryData["Mountpoint"])
	secondaryBackingDevice, _ := filepath.Abs(path)
	log.Info(primaryBackingDevice, secondaryBackingDevice)
	// as we use the same minors/ports in primary/secondary, make them a set:
	usedPort := make(map[int]bool)
	usedMinor := make(map[int]bool)
	for _, volData := range opt.GetVolumeDataList() {
		data := volData.GetData()

		// check if the current device in the DataList() already has a key for port/minor that belongs to the primary/secondary
		// That would happen for example if the data was not deleted and/or CreateReplication() was called multiple times.
		// if val, ok := data[portKey(primaryVolID)]; ok {
		// 	return nil, fmt.Errorf("Primary Volume ID (%s), already has a port number (%s)", primaryVolID, val)
		// }
		// if val, ok := data[portKey(secondaryVolID)]; ok {
		// 	return nil, fmt.Errorf("Secondary Volume ID (%s), already has a port number (%s)", secondaryVolID, val)
		// }
		// if val, ok := data[minorKey(primaryVolID)]; ok {
		// 	return nil, fmt.Errorf("Primary Volume ID (%s), already has a minor number (%s)", primaryVolID, val)
		// }
		// if val, ok := data[minorKey(secondaryVolID)]; ok {
		// 	return nil, fmt.Errorf("Secondary Volume ID (%s), already has a minor number (%s)", secondaryVolID, val)
		// }

		// get ports and minors still in use
		volID := data["VolumeId"]
		if val, ok := data[portKey(volID)]; ok {
			p, err := strconv.Atoi(val)
			if err != nil {
				return nil, err
			}
			usedPort[p] = true
		}
		if val, ok := data[minorKey(volID)]; ok {
			m, err := strconv.Atoi(val)
			if err != nil {
				return nil, err
			}
			usedMinor[m] = true
		}
	}

	portMin := cfgOrDefault(conf.PortMin, defaultPortMin)
	portMax := cfgOrDefault(conf.PortMax, defaultPortMax)
	var up []int
	for k := range usedPort {
		up = append(up, k)
	}
	port, err := godrbdutils.GetNumber(portMin, portMax, up)
	if err != nil {
		return nil, err
	}

	var um []int
	for k := range usedMinor {
		um = append(um, k)
	}
	minor, err := godrbdutils.GetNumber(minorMin, minorMax, um)
	if err != nil {
		return nil, err
	}

	res := godrbdutils.NewResource(resName, port)

	res.AddHost(myHost.ID, myHost.Name, myHost.IP)
	res.AddHost(peerHost.ID, peerHost.Name, peerHost.IP)

	drbdVolID := 0 // currently only one volume per DRBD resource
	if isPrimary {
		res.AddVolume(drbdVolID, minor, primaryBackingDevice, myHost.Name)
		res.AddVolume(drbdVolID, minor, secondaryBackingDevice, peerHost.Name)
	} else {
		res.AddVolume(drbdVolID, minor, primaryBackingDevice, peerHost.Name)
		res.AddVolume(drbdVolID, minor, secondaryBackingDevice, myHost.Name)
	}

	res.WriteConfig(resFilePath(resName))

	// Bring up the resource
	drbdadm := godrbdutils.NewDrbdAdm([]string{resName})

	drbdadm.CreateMetaData(fmt.Sprintf("--max-peers=%d", maxPeers), "--force")
	drbdadm.Up()

	if isPrimary {
		// start initial sync
		drbdadm.Primary("--force")
		drbdadm.Secondary() // switch back, rest done by auto promote
	}

	additionalPrimaryData := map[string]string{
		portKey(primaryVolID):  strconv.Itoa(port),
		minorKey(primaryVolID): strconv.Itoa(minor),
	}

	additionalSecondaryData := map[string]string{
		portKey(secondaryVolID):  strconv.Itoa(port),
		minorKey(secondaryVolID): strconv.Itoa(minor),
	}

	return &model.ReplicationSpec{
		// TODO(rck): return additional important information
		PrimaryReplicationDriverData:   additionalPrimaryData,
		SecondaryReplicationDriverData: additionalSecondaryData,
	}, nil
}

func (r *ReplicationDriver) DeleteReplication(opt *pb.DeleteReplicationOpts) error {
	log.Infof("DRBD delete replication ....")

	resName := opt.GetId()

	drbdadm := godrbdutils.NewDrbdAdm([]string{resName})
	if _, err := drbdadm.Down(); err != nil {
		return err
	}
	if err := os.Remove(resFilePath(resName)); err != nil {
		return err
	}

	// reserved minor/port are automatically deleted because they are gone from *ReplicationData

	return nil
}

func (r *ReplicationDriver) EnableReplication(opt *pb.EnableReplicationOpts) error {
	log.Infof("DRBD enable replication ....")

	drbdadm := godrbdutils.NewDrbdAdm([]string{opt.GetId()})
	_, err := drbdadm.Adjust()
	return err
}

func (r *ReplicationDriver) DisableReplication(opt *pb.DisableReplicationOpts) error {
	log.Infof("DRBD disable replication ....")

	drbdadm := godrbdutils.NewDrbdAdm([]string{opt.GetId()})
	_, err := drbdadm.Disconnect()
	return err
}

func (r *ReplicationDriver) FailoverReplication(opt *pb.FailoverReplicationOpts) error {
	log.Infof("DRBD failover replication ....")
	// nothing to do here:
	// The driver returns a block device on both nodes (/dev/drbd$minor and a symlink as /dev/drbd/by-res/$resname)
	// And the driver makes sure that it triggeres an initial sync from the primary to the secondary side
	// Then:
	// When the device is used (open(2) in RW mode, it switches that side to DRBD Primary. That is the "autopromote" feature of DRBD9
	// That happens what ever the user does (use it as raw device with 'dd', or put a file system on it and mount it,...)
	// After the user finished using the device (e.g., umount), the device switches to DRBD Secondary
	// And it can then be used on the second node by just open(2)ing the device again.
	return nil
}
