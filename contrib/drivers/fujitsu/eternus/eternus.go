// Copyright 2019 The OpenSDS Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package eternus

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	log "github.com/golang/glog"
	. "github.com/opensds/opensds/contrib/drivers/utils/config"
	"github.com/opensds/opensds/pkg/model"
	pb "github.com/opensds/opensds/pkg/model/proto"

	"github.com/opensds/opensds/pkg/utils/config"
	uuid "github.com/satori/go.uuid"
)

// Driver
type Driver struct {
	conf   *EternusConfig
	client *EternusClient
}

// Setup eternus driver
func (d *Driver) Setup() (err error) {
	// Read fujitsu eternus config file
	conf := &EternusConfig{}
	d.conf = conf
	path := config.CONF.OsdsDock.Backends.FujitsuEternus.ConfigPath

	if "" == path {
		path = defaultConfPath
	}
	Parse(conf, path)
	d.client, err = NewClient(&d.conf.AuthOptions)
	if err != nil {
		log.Errorf("failed to get new client, %v", err)
		return err
	}
	err = d.client.login()
	if err != nil {
		log.Errorf("failed to login, %v", err)
		return err
	}
	return nil
}

// Unset eternus driver
func (d *Driver) Unset() error {
	d.client.Destroy()
	return nil
}

// ListPools : get pool list
func (d *Driver) ListPools() ([]*model.StoragePoolSpec, error) {
	var pols []*model.StoragePoolSpec
	sp, err := d.client.ListStoragePools()
	if err != nil {
		return nil, err
	}
	for _, p := range sp {
		c := d.conf
		if _, ok := c.Pool[p.Name]; !ok {
			continue
		}

		host, _ := os.Hostname()
		name := fmt.Sprintf("%s:%s:%s", host, d.conf.Endpoint, p.Id)
		pol := &model.StoragePoolSpec{
			BaseModel: &model.BaseModel{
				Id: uuid.NewV5(uuid.NamespaceOID, name).String(),
			},
			Name:             p.Name,
			TotalCapacity:    p.TotalCapacity,
			FreeCapacity:     p.FreeCapacity,
			StorageType:      c.Pool[p.Name].StorageType,
			Extras:           c.Pool[p.Name].Extras,
			AvailabilityZone: c.Pool[p.Name].AvailabilityZone,
		}
		if pol.AvailabilityZone == "" {
			pol.AvailabilityZone = defaultAZ
		}
		pols = append(pols, pol)
	}
	return pols, nil
}

func (d *Driver) createVolumeFromSnapshot(opt *pb.CreateVolumeOpts) (*model.VolumeSpec, error) {
	return nil, &model.NotImplementError{"method createVolumeFromSnapshot is not implement."}
}

// CreateVolume : create volume.
func (d *Driver) CreateVolume(opt *pb.CreateVolumeOpts) (*model.VolumeSpec, error) {
	log.Infof("start creating volume. opt = %v", opt)
	if opt.GetSnapshotId() != "" {
		return d.createVolumeFromSnapshot(opt)
	}

	id := opt.GetId()
	desc := opt.GetDescription()
	provPolicy := d.conf.Pool[opt.GetPoolName()].Extras.DataStorage.ProvisioningPolicy
	// execute create volume
	vol, err := d.client.CreateVolume(id, opt.GetSize(), desc, opt.GetPoolName(), provPolicy)
	if err != nil {
		log.Error("create Volume Failed:", err)
		return nil, err
	}
	log.Infof("create volume %s (%s) success.", opt.GetName(), vol.Id)
	return &model.VolumeSpec{
		BaseModel: &model.BaseModel{
			Id: opt.GetId(),
		},
		Name:             opt.GetName(),
		Size:             opt.GetSize(),
		Description:      opt.GetDescription(),
		AvailabilityZone: opt.GetAvailabilityZone(),
		Metadata: map[string]string{
			KLunId: vol.Id,
		},
	}, nil
}

// PullVolume : get volume information
func (d *Driver) PullVolume(id string) (*model.VolumeSpec, error) {
	return nil, &model.NotImplementError{"method PullVolume is not implement."}
}

// DeleteVolume : delete volume
func (d *Driver) DeleteVolume(opt *pb.DeleteVolumeOpts) error {
	log.Infof("start delete volume. opt = %v", opt)
	volID := opt.GetMetadata()[KLunId]
	err := d.client.DeleteVolume(volID)
	if err != nil {
		log.Error("remove Volume Failed:", err)
		return err
	}
	log.Infof("delete volume (%s) success.", opt.GetId())
	return nil
}

// ExtendVolume : extend volume
func (d *Driver) ExtendVolume(opt *pb.ExtendVolumeOpts) (*model.VolumeSpec, error) {
	log.Infof("start extend volume. opt = %v", opt)

	volID := opt.GetMetadata()[KLunId]
	// execute extend volume
	err := d.client.ExtendVolume(volID, opt.GetSize())
	if err != nil {
		log.Error("extend Volume Failed:", err)
		return nil, err
	}

	log.Infof("extend volume %s (%s) success.", opt.GetName(), volID)
	return &model.VolumeSpec{
		BaseModel: &model.BaseModel{
			Id: opt.GetId(),
		},
		Name:             opt.GetName(),
		Size:             opt.GetSize(),
		Description:      opt.GetDescription(),
		AvailabilityZone: opt.GetAvailabilityZone(),
	}, nil
}

// InitializeConnection :
func (d *Driver) InitializeConnection(opt *pb.CreateVolumeAttachmentOpts) (*model.ConnectionInfo, error) {
	if opt.GetAccessProtocol() == ISCSIProtocol {
		return d.initializeConnectionIscsi(opt)
	}
	if opt.GetAccessProtocol() == FCProtocol {
		return d.initializeConnectionFC(opt)
	}
	return nil, errors.New("no supported protocol for eternus driver")
}

func (d *Driver) initializeConnectionIscsi(opt *pb.CreateVolumeAttachmentOpts) (*model.ConnectionInfo, error) {

	var err error
	lunID := opt.GetMetadata()[KLunId]
	hostInfo := opt.GetHostInfo()
	hostID := ""
	lunGrpID := ""
	hostLun := ""
	hostExist := false

	// check initiator is specified
	initiator := hostInfo.GetInitiator()
	needHostAffinity := true
	if initiator == "" {
		needHostAffinity = false
	}

	// Get port info
	iscsiPortInfo, err := d.client.GetIscsiPortInfo(d.conf.CeSupport, needHostAffinity)
	if err != nil {
		log.Errorf("get iscsi port failed. error: %v", err)
		return nil, err
	}

	// Create host if not exist.
	if needHostAffinity {
		// Create resource name
		initiator := hostInfo.GetInitiator()
		ipAddr := hostInfo.GetIp()
		rscName := GetFnvHash(initiator + ipAddr)
		hostID, hostExist, err = d.client.AddIscsiHostWithCheck(rscName, initiator, ipAddr)
		if err != nil {
			log.Errorf("failed to add host, hostInfo =%v, error: %v", hostInfo, err)
			return nil, err
		}

		// Create Lun group
		lunGrpID, err = d.client.AddLunGroupWithCheck(rscName, lunID)
		if err != nil {
			log.Errorf("failed to add lun group, lun group name =%s, error: %v", rscName, err)
			return nil, err
		}
		// skip AddHostAffinity if host already exists.
		if !hostExist {
			// Create host affinity
			_, err = d.client.AddHostAffinity(lunGrpID, hostID, iscsiPortInfo.PortNumber)
			if err != nil {
				log.Errorf("failed to add host affinity, lunGrp id=%s, hostID=%s, error: %v",
					lunGrpID, hostID, err)
				return nil, err
			}
		}
		hostLun, err = d.client.GetHostLunID(lunGrpID, lunID)
		if err != nil {
			log.Error("failed to get the host lun id,", err)
			return nil, err
		}
	} else {
		hostLun, err = d.addMapping(iscsiPortInfo.PortNumber, lunID)
		if err != nil {
			return nil, err
		}
	}

	log.Infof("initialize iscsi connection (%s) success.", opt.GetId())
	connInfo := &model.ConnectionInfo{
		DriverVolumeType: ISCSIProtocol,
		ConnectionData: map[string]interface{}{
			"targetDiscovered": true,
			"targetIQN":        iscsiPortInfo.IscsiName,
			"targetPortal":     iscsiPortInfo.Ip + ":" + strconv.Itoa(iscsiPortInfo.TcpPort),
			"discard":          false,
			"targetLun":        hostLun,
		},
	}
	return connInfo, nil
}

func (d *Driver) initializeConnectionFC(opt *pb.CreateVolumeAttachmentOpts) (*model.ConnectionInfo, error) {

	var err error
	lunID := opt.GetMetadata()[KLunId]
	hostInfo := opt.GetHostInfo()
	hostID := ""
	lunGrpID := ""
	hostLun := ""
	hostExist := false

	// check initiator is specified
	initiator := hostInfo.GetInitiator()
	needHostAffinity := true
	if initiator == "" {
		needHostAffinity = false
	}

	// Get port info
	fcPortInfo, err := d.client.GetFcPortInfo(d.conf.CeSupport, needHostAffinity)
	if err != nil {
		log.Errorf("failed to get fc port. error: %v", err)
		return nil, err
	}

	// initiator is specified
	if needHostAffinity {
		wwnName := hostInfo.GetInitiator()
		rscName := GetFnvHash(wwnName)
		// Create host if not exist.
		hostID, hostExist, err = d.client.AddFcHostWithCheck(rscName, wwnName)
		if err != nil {
			log.Errorf("failed to add host, host name =%s, error: %v", hostInfo.Host, err)
			return nil, err
		}
		// Create Lun group
		lunGrpID, err = d.client.AddLunGroupWithCheck(rscName, lunID)
		if err != nil {
			log.Errorf("failed to add lun group, lun group name =%s, error: %v", hostInfo.Host, err)
			return nil, err
		}
		// skip AddHostAffinity if host already exists.
		if !hostExist {
			// Create host affinity
			_, err = d.client.AddHostAffinity(lunGrpID, hostID, fcPortInfo.PortNumber)
			if err != nil {
				log.Errorf("failed to add host affinity, lunGrp id=%s, hostId=%s, error: %v",
					lunGrpID, hostID, err)
				return nil, err
			}
		}
		hostLun, err = d.client.GetHostLunID(lunGrpID, lunID)
		if err != nil {
			log.Error("failed to get the host lun id,", err)
			return nil, err
		}
	} else {
		hostLun, err = d.addMapping(fcPortInfo.PortNumber, lunID)
		if err != nil {
			return nil, err
		}
	}

	log.Infof("initialize fc connection (%s) success.", opt.GetId())
	fcInfo := &model.ConnectionInfo{
		DriverVolumeType: FCProtocol,
		ConnectionData: map[string]interface{}{
			"targetDiscovered": true,
			"targetWwn":        fcPortInfo.Wwpn,
			"hostname":         opt.GetHostInfo().Host,
			"targetLun":        hostLun,
		},
	}
	return fcInfo, nil
}

func (d *Driver) addMapping(PortNumber string, lunID string) (string, error) {
	hostLunID := "0"
	// get exist mapping
	mappings, err := d.client.ListMapping(PortNumber)
	if err != nil {
		log.Error("failed to get mapping,", err)
		return "", err
	}
	// get unused host lun id
	if len(mappings) >= 1024 {
		msg := "reached the upper limit to add mapping"
		log.Error("failed to get host lun id,", msg)
		return "", errors.New(msg)
	}
	for i, v := range mappings {
		if v.Lun != strconv.Itoa(i) {
			hostLunID = strconv.Itoa(i)
			break
		}
		hostLunID = strconv.Itoa(i + 1)
	}
	// add mapping
	err = d.client.AddMapping(lunID, hostLunID, PortNumber)
	if err != nil {
		log.Error("failed to add mapping,", err)
		return "", err
	}
	return hostLunID, nil
}

func (d *Driver) deleteMapping(PortNumber string, lunID string) error {
	hostLunID := ""
	// get exist mapping
	mappings, err := d.client.ListMapping(PortNumber)
	if err != nil {
		log.Error("failed to get mapping,", err)
		return err
	}
	for _, v := range mappings {
		if v.VolumeNumber == lunID {
			hostLunID = v.Lun
			break
		}
	}
	if hostLunID == "" {
		log.Infof("specified mapping already deleted, PortNumber = %s, lunID =%s", PortNumber, lunID)
		return nil
	}
	// add mapping
	err = d.client.DeleteMapping(hostLunID, PortNumber)
	if err != nil {
		log.Error("failed to delete mapping,", err)
		return err
	}
	return nil
}

// TerminateConnection :
func (d *Driver) TerminateConnection(opt *pb.DeleteVolumeAttachmentOpts) error {

	if opt.GetAccessProtocol() != ISCSIProtocol &&
		opt.GetAccessProtocol() != FCProtocol {
		return errors.New("no supported protocol for eternus driver")
	}

	lunID := opt.GetMetadata()[KLunId]
	hostInfo := opt.GetHostInfo()
	initiator := hostInfo.GetInitiator()
	needHostAffinity := true
	if initiator == "" {
		needHostAffinity = false
	}

	// Get port info
	portNumber := ""
	if opt.GetAccessProtocol() == ISCSIProtocol {
		iscsiPortInfo, err := d.client.GetIscsiPortInfo(d.conf.CeSupport, needHostAffinity)
		if err != nil {
			log.Errorf("failed to get iscsi port. error: %v", err)
			return err
		}
		portNumber = iscsiPortInfo.PortNumber
	} else if opt.GetAccessProtocol() == FCProtocol {
		fcPortInfo, err := d.client.GetFcPortInfo(d.conf.CeSupport, needHostAffinity)
		if err != nil {
			log.Errorf("failed to get fc port. error: %v", err)
			return err
		}
		portNumber = fcPortInfo.PortNumber
	}

	// if no need to delete host affinity, delete mapping
	if needHostAffinity != true {
		err := d.deleteMapping(portNumber, lunID)
		if err != nil {
			log.Errorf("failed to delete mapping. error: %v", err)
			return err
		}
		log.Infof("terminate connection (%s) success.", opt.GetId())
		return nil
	}

	// Create resource name
	rscName := ""
	if opt.GetAccessProtocol() == ISCSIProtocol {
		ipAddr := hostInfo.GetIp()
		rscName = GetFnvHash(initiator + ipAddr)
	} else {
		rscName = GetFnvHash(initiator)
	}

	// Get lun group
	lg, err := d.client.GetLunGroupByName(rscName)
	if err != nil {
		log.Errorf("failed to get lun group, error: %v", err)
		return err
	}
	// if lun group has some volumes.
	if len(lg.Volumes) > 1 {
		hostLunID := ""
		for _, v := range lg.Volumes {
			if v.Id == lunID {
				hostLunID = v.Lun
				break
			}
		}
		if hostLunID == "" {
			log.Errorf("target volume already removed from lun group, lunID: %s", lunID)
		} else {
			err = d.client.RemoveVolumeFromLunGroup(hostLunID, rscName)
			if err != nil {
				log.Errorf("failed to remove volume from lun group, error: %v", err)
				return err
			}
		}
	} else {
		// Delete host affinity
		err = d.client.DeleteHostAffinity(portNumber, rscName)
		if err != nil {
			log.Errorf("failed to delete host affinity, error: %v", err)
			return err
		}

		// Delete lun group
		err = d.client.DeleteLunGroupByName(rscName)
		if err != nil {
			log.Errorf("failed to delete lun group, error: %v", err)
			return err
		}

		// Delete host
		if opt.GetAccessProtocol() == ISCSIProtocol {
			err = d.client.DeleteIscsiHostByName(rscName)
		} else if opt.GetAccessProtocol() == FCProtocol {
			err = d.client.DeleteFcHostByName(rscName)
		}
		if err != nil {
			log.Errorf("failed to delete host, error: %v", err)
			return err
		}
	}
	log.Infof("terminate connection (%s) success.", opt.GetId())
	return nil
}

// CreateSnapshot :
func (d *Driver) CreateSnapshot(opt *pb.CreateVolumeSnapshotOpts) (*model.VolumeSnapshotSpec, error) {
	lunID := opt.GetMetadata()[KLunId]
	// get source volume information for getting pool information
	vol, err := d.client.GetVolume(lunID)
	if err != nil {
		log.Errorf("failed to get volume, error: %v", err)
		return nil, err
	}
	poolName := vol.PoolName

	// create snapshot volume
	provPolicy := d.conf.Pool[poolName].Extras.DataStorage.ProvisioningPolicy
	vol, err = d.client.CreateVolume(opt.GetId(), opt.GetSize(), opt.GetDescription(),
		poolName, provPolicy)
	if err != nil {
		log.Errorf("failed to create snapshot volume, error: %v", err)
		return nil, err
	}

	// get Client for admin role
	adminClient, err := NewClientForAdmin(&d.conf.AuthOptions)
	if err != nil {
		log.Errorf("failed to get new client, %v", err)
		return nil, err
	}
	err = adminClient.login()
	if err != nil {
		log.Errorf("failed to login, %v", err)
		return nil, err
	}
	defer adminClient.Destroy()

	// Start SnapOPC+ session (create shapshot)
	err = adminClient.CreateSnapshot(lunID, vol.Id)
	if err != nil {
		log.Errorf("failed to create snapopc+ session, error: %v", err)
		return nil, err
	}

	// get session id
	snapshotList, err := adminClient.ListSnapshot()
	snapshot := SnapShot{}
	for _, v := range snapshotList {
		if v.SrcNo == lunID && v.DestNo == vol.Id {
			snapshot = v
			break
		}
	}
	log.Info("create snapshot success, snapshot id =", opt.GetId())
	return &model.VolumeSnapshotSpec{
		BaseModel: &model.BaseModel{
			Id: opt.GetId(),
		},
		Name:        opt.GetName(),
		Description: opt.GetDescription(),
		VolumeId:    opt.GetVolumeId(),
		Size:        0,
		Metadata: map[string]string{
			KSnapId:    snapshot.Sid,
			KSnapLunId: vol.Id,
		},
	}, nil
}

// PullSnapshot :
func (d *Driver) PullSnapshot(snapIdentifier string) (*model.VolumeSnapshotSpec, error) {
	return nil, &model.NotImplementError{"method PullSnapshot is not implement."}
}

// DeleteSnapshot :
func (d *Driver) DeleteSnapshot(opt *pb.DeleteVolumeSnapshotOpts) error {
	sid := opt.GetMetadata()[KSnapId]
	volID := opt.GetMetadata()[KSnapLunId]

	// get Client for admin role
	adminClient, err := NewClientForAdmin(&d.conf.AuthOptions)
	if err != nil {
		log.Errorf("failed to create client, error: %v", err)
		return err
	}
	err = adminClient.login()
	if err != nil {
		log.Errorf("failed to login, %v", err)
		return err
	}
	defer adminClient.Destroy()

	// delete snapshot
	err = adminClient.DeleteSnapshot(sid)
	if err != nil {
		log.Errorf("failed to delete snapshot, snapshot id = %s , error: %v", opt.GetId(), err)
		return err
	}

	// delete snapshot volume
	err = d.client.DeleteVolume(volID)
	if err != nil {
		log.Errorf("failed to delete snapshot volume, volume id = %s , error: %v", volID, err)
		return err
	}
	log.Info("delete snapshot success, snapshot id =", opt.GetId())
	return nil
}

func (d *Driver) InitializeSnapshotConnection(opt *pb.CreateSnapshotAttachmentOpts) (*model.ConnectionInfo, error) {
	return nil, &model.NotImplementError{S: "method InitializeSnapshotConnection has not been implemented yet."}
}

func (d *Driver) TerminateSnapshotConnection(opt *pb.DeleteSnapshotAttachmentOpts) error {
	return &model.NotImplementError{S: "method TerminateSnapshotConnection has not been implemented yet."}
}

func (d *Driver) CreateVolumeGroup(opt *pb.CreateVolumeGroupOpts) (*model.VolumeGroupSpec, error) {
	return nil, &model.NotImplementError{"method CreateVolumeGroup has not been implemented yet"}
}

func (d *Driver) UpdateVolumeGroup(opt *pb.UpdateVolumeGroupOpts) (*model.VolumeGroupSpec, error) {
	return nil, &model.NotImplementError{"method UpdateVolumeGroup has not been implemented yet"}
}

func (d *Driver) DeleteVolumeGroup(opt *pb.DeleteVolumeGroupOpts) error {
	return &model.NotImplementError{"method DeleteVolumeGroup has not been implemented yet"}
}
