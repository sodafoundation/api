// Copyright (c) 2019 Huawei Technologies Co., Ltd. All Rights Reserved.
//
//    Licensed under the Apache License, Version 2.0 (the "License"); you may
//    not use this file except in compliance with the License. You may obtain
//    a copy of the License at
//
//         http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//    WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//    License for the specific language governing permissions and limitations
//    under the License.

package fusionstorage

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	log "github.com/golang/glog"
	. "github.com/opensds/opensds/contrib/drivers/utils/config"
	pb "github.com/opensds/opensds/pkg/dock/proto"
	. "github.com/opensds/opensds/pkg/model"
	"github.com/satori/go.uuid"
)

type Driver struct {
	cli  *Cli
	conf *Config
}

type AuthOptions struct {
	Username string   `yaml:"username"`
	Password string   `yaml:"password"`
	Url      string   `yaml:"url"`
	FmIp     string   `yaml:"fmIp,omitempty"`
	FsaIp    []string `yaml:"fsaIp,flow"`
}

type Config struct {
	AuthOptions `yaml:"authOptions"`
	Pool        map[string]PoolProperties `yaml:"pool,flow"`
}

func (d *Driver) Setup() error {
	conf := &Config{}

	d.conf = conf

	path := "./testdata/fusionstorage.yaml"
	if path == "" {
		path = ""
	}

	Parse(conf, path)

	client, err := newRestCommon(conf.Username, conf.Password, conf.Url, conf.FmIp, conf.FsaIp)
	if err != nil {
		return err
	}

	err = client.login()
	if err != nil {
		fmt.Printf("Get new client failed, %v", err)
		return err
	}

	_, err = exec.LookPath(CmdBin)
	if err != nil {
		if err == exec.ErrNotFound {
			return fmt.Errorf("%q executable not found in $PATH", CmdBin)
		}
		return err
	}

	err = client.StartServer()
	if err != nil {
		return err
	}

	d.cli = client

	return nil
}

func (d *Driver) Unset() error {
	return nil
}

func EncodeName(id string) string {
	return NamePrefix + "-" + id
}

func (d *Driver) CreateVolume(opt *pb.CreateVolumeOpts) (*VolumeSpec, error) {
	name := EncodeName(opt.GetId())
	err := d.cli.createVolume(name, opt.GetPoolName(), opt.GetSize()<<UnitGiShiftBit)
	if err != nil {
		log.Errorf("Create volume %s (%s) failed: %s", opt.GetName(), opt.GetId(), err)
		return nil, err
	}
	log.Infof("Create volume %s (%s) success.", opt.GetName(), opt.GetId())
	return &VolumeSpec{
		BaseModel: &BaseModel{
			Id: opt.GetId(),
		},
		Name:             opt.GetName(),
		Size:             opt.Size,
		Description:      opt.GetDescription(),
		AvailabilityZone: opt.GetAvailabilityZone(),
		PoolId:           opt.GetPoolId(),
		Metadata: map[string]string{
			LunId: name,
		},
	}, nil
}

func (d *Driver) DeleteVolume(opt *pb.DeleteVolumeOpts) error {
	name := EncodeName(opt.GetId())
	err := d.cli.deleteVolume(name)
	if err != nil {
		log.Errorf("Delete volume (%s) failed: %v", opt.GetId(), err)
		return err
	}
	log.Infof("Delete volume (%s) success.", opt.GetId())
	return nil
}

func (d *Driver) ListPools() ([]*StoragePoolSpec, error) {
	var pols []*StoragePoolSpec
	pools, err := d.cli.queryPoolInfo()
	if err != nil {
		return nil, err
	}

	c := d.conf
	for _, p := range pools.Pools {
		poolId := strconv.Itoa(p.PoolId)
		if _, ok := c.Pool[poolId]; !ok {
			continue
		}
		host, _ := os.Hostname()
		name := fmt.Sprintf("%s:%s:%s", host, d.conf.Url, poolId)
		pol := &StoragePoolSpec{
			BaseModel: &BaseModel{
				// Make sure uuid is unique
				Id: uuid.NewV5(uuid.NamespaceOID, name).String(),
			},
			Name:             poolId,
			TotalCapacity:    p.TotalCapacity >> UnitGiShiftBit,
			FreeCapacity:     (p.TotalCapacity - p.UsedCapacity) >> UnitGiShiftBit,
			StorageType:      c.Pool[poolId].StorageType,
			Extras:           c.Pool[poolId].Extras,
			AvailabilityZone: c.Pool[poolId].AvailabilityZone,
		}
		if pol.AvailabilityZone == "" {
			pol.AvailabilityZone = DefaultAZ
		}
		pols = append(pols, pol)
	}
	return pols, nil
}

func (d *Driver) InitializeConnection(opt *pb.CreateAttachmentOpts) (*ConnectionInfo, error) {
	lunId := opt.GetMetadata()[LunId]
	if lunId == "" {
		return nil, fmt.Errorf("Lun id is empty.")
	}

	connectorType := opt.GetMetadata()[ConnectorType]
	if connectorType == "" || connectorType != FusionstorageIscsi {
		return nil, fmt.Errorf("Connector type is empty or mismatch.")
	}

	hostInfo := opt.GetHostInfo()

	initiator := hostInfo.GetInitiator()
	hostName := hostInfo.GetHost()

	if initiator == "" || hostName == "" {
		return nil, fmt.Errorf("Host name or initiator is empty.")
	}

	// Create port if not exist.
	err := d.cli.queryPortInfo(initiator)
	if err != nil {
		if err.Error() == InitiatorNotExistErrorCode {
			err := d.cli.createPort(initiator)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	// Create host if not exist.
	isFind, err := d.cli.queryHostInfo(hostName)
	if err != nil {
		return nil, fmt.Errorf("Query host failed, host name =%s, error: %v", hostName, err)
	}

	if !isFind {
		err = d.cli.createHost(hostInfo)
		if err != nil {
			return nil, fmt.Errorf("Create host failed, host name =%s, error: %v", hostName, err)
		}
	}

	// Add port to host if port not add to the host
	hostPortMap, err := d.cli.queryHostByPort(initiator)
	if err != nil {
		return nil, err
	}

	h, ok := hostPortMap.PortHostMap[initiator]
	if ok && h[0] != hostName {
		return nil, fmt.Errorf("Initiator is already added to another host, host name =%s", h[0])
	}

	if !ok {
		err = d.cli.addPortToHost(hostName, initiator)
		if err != nil {
			return nil, fmt.Errorf("Add port to host failed, error %v", err)
		}
	}

	// Map volume to host
	err = d.cli.addLunsToHost(hostName, lunId)
	if err != nil {
		return nil, fmt.Errorf("Add luns to host failed, error %v", err)
	}

	// Get target lun id
	hostLunList, err := d.cli.queryHostLunInfo(hostName)
	if err != nil {
		return nil, err
	}

	var targetLunId int
	for _, v := range hostLunList.LunList {
		if v.Name == lunId {
			targetLunId = v.Id
		}
	}

	// Get target iscsi portal info
	targetPortalInfo, err := d.cli.queryIscsiPortal(initiator)
	if err != nil {
		return nil, err
	}

	iscsiTarget := strings.Split(targetPortalInfo, ",")

	connInfo := &ConnectionInfo{
		DriverVolumeType: ISCSIProtocol,
		ConnectionData: map[string]interface{}{
			"target_discovered": true,
			"volume_id":         opt.GetVolumeId(),
			"description":       "huawei",
			"host_name":         hostName,
			"target_lun":        targetLunId,
			"connect_type":      FusionstorageIscsi,
			"host":              hostName,
			"initiator":         initiator,
			"targetIQN":         []string{iscsiTarget[1]},
			"targetPortal":      []string{iscsiTarget[0]},
		},
	}

	return connInfo, nil
}

func (d *Driver) TerminateConnection(opt *pb.DeleteAttachmentOpts) error {
	lunId := opt.GetMetadata()[LunId]
	if lunId == "" {
		return fmt.Errorf("Lun id is empty.")
	}
	hostInfo := opt.GetHostInfo()

	initiator := hostInfo.GetInitiator()
	hostName := hostInfo.GetHost()

	if initiator == "" || hostName == "" {
		return fmt.Errorf("Host name or initiator is empty.")
	}

	// Make sure that host is exist.
	hostIsFind, err := d.cli.queryHostInfo(hostName)
	if err != nil {
		return fmt.Errorf("Query host failed, host name =%s, error: %v", hostName, err)
	}

	if !hostIsFind {
		return fmt.Errorf("Host can not be found, host name =%s", hostName)
	}

	// Check whether the volume attach to the host
	hostLunList, err := d.cli.queryHostLunInfo(hostName)
	if err != nil {
		return err
	}

	var lunIsFind = false
	for _, v := range hostLunList.LunList {
		if v.Name == lunId {
			lunIsFind = true
			break
		}
	}

	if !lunIsFind {
		return fmt.Errorf("The lun %s is not attach to the host %s", lunId, hostName)
	}

	// Remove lun from host
	err = d.cli.deleteLunFromHost(hostName, lunId)
	if err != nil {
		return err
	}

	// Remove initiator and host if there is no lun belong to the host
	fmt.Println(len(hostLunList.LunList), hostLunList.LunList)
	hostLunList, err = d.cli.queryHostLunInfo(hostName)
	if err != nil {
		return err
	}

	if len(hostLunList.LunList) == 0 {
		d.cli.deletePortFromHost(hostName, initiator)
		d.cli.deleteHost(hostName)
		d.cli.deletePort(initiator)
	}

	return nil
}

func (d *Driver) PullVolume(volIdentifier string) (*VolumeSpec, error) {
	// Not used , do nothing
	return nil, nil
}

func (d *Driver) ExtendVolume(opt *pb.ExtendVolumeOpts) (*VolumeSpec, error) {
	err := d.cli.extendVolume(EncodeName(opt.GetId()), opt.GetSize()<<UnitGiShiftBit)
	if err != nil {
		log.Errorf("Extend volume %s (%s) failed: %v", opt.GetName(), opt.GetId(), err)
		return nil, err
	}
	log.Infof("Extend volume %s (%s) success.", opt.GetName(), opt.GetId())
	return &VolumeSpec{
		BaseModel: &BaseModel{
			Id: opt.GetId(),
		},
		Name:             opt.GetName(),
		Size:             opt.GetSize(),
		Description:      opt.GetDescription(),
		AvailabilityZone: opt.GetAvailabilityZone(),
	}, nil
}

func (d *Driver) CreateSnapshot(opt *pb.CreateVolumeSnapshotOpts) (*VolumeSnapshotSpec, error) {
	snapName := EncodeName(opt.GetId())
	volName := EncodeName(opt.GetVolumeId())

	if err := d.cli.createSnapshot(snapName, volName); err != nil {
		log.Errorf("Create snapshot %s (%s) failed: %s", opt.GetName(), opt.GetId(), err)
		return nil, err
	}

	log.Errorf("Create snapshot %s (%s) success.", opt.GetName(), opt.GetId())
	return &VolumeSnapshotSpec{
		BaseModel: &BaseModel{
			Id: opt.GetId(),
		},
		Name:        opt.GetName(),
		Description: opt.GetDescription(),
		VolumeId:    opt.GetVolumeId(),
		Size:        opt.GetSize(),
	}, nil
}

func (d *Driver) PullSnapshot(snapIdentifier string) (*VolumeSnapshotSpec, error) {
	return nil, nil
}

func (d *Driver) DeleteSnapshot(opt *pb.DeleteVolumeSnapshotOpts) error {
	err := d.cli.deleteSnapshot(EncodeName(opt.GetId()))
	if err != nil {
		log.Errorf("Delete volume snapshot (%s) failed: %v", opt.GetId(), err)
		return err
	}
	log.Infof("Remove volume snapshot (%s) success", opt.GetId())
	return nil
}

func (d *Driver) InitializeSnapshotConnection(opt *pb.CreateSnapshotAttachmentOpts) (*ConnectionInfo, error) {
	return nil, &NotImplementError{S: "Method InitializeSnapshotConnection has not been implemented yet."}
}

func (d *Driver) TerminateSnapshotConnection(opt *pb.DeleteSnapshotAttachmentOpts) error {
	return &NotImplementError{S: "Method TerminateSnapshotConnection has not been implemented yet."}
}

func (d *Driver) CreateVolumeGroup(
	opt *pb.CreateVolumeGroupOpts,
	vg *VolumeGroupSpec) (*VolumeGroupSpec, error) {
	return nil, &NotImplementError{S: "Method CreateVolumeGroup has not been implemented yet."}
}

func (d *Driver) UpdateVolumeGroup(
	opt *pb.UpdateVolumeGroupOpts,
	vg *VolumeGroupSpec,
	addVolumesRef []*VolumeSpec,
	removeVolumesRef []*VolumeSpec) (*VolumeGroupSpec, []*VolumeSpec, []*VolumeSpec, error) {
	return nil, nil, nil, &NotImplementError{"Method UpdateVolumeGroup has not been implemented yet"}
}

func (d *Driver) DeleteVolumeGroup(
	opt *pb.DeleteVolumeGroupOpts,
	vg *VolumeGroupSpec,
	volumes []*VolumeSpec) (*VolumeGroupSpec, []*VolumeSpec, error) {
	return nil, nil, &NotImplementError{S: "Method DeleteVolumeGroup has not been implemented yet."}
}
