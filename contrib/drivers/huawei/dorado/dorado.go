// Copyright 2017 The OpenSDS Authors.
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

package dorado

import (
	"errors"
	"fmt"
	"os"
	"strings"

	log "github.com/golang/glog"
	. "github.com/opensds/opensds/contrib/drivers/utils/config"
	"github.com/opensds/opensds/pkg/model"
	pb "github.com/opensds/opensds/pkg/model/proto"
	"github.com/opensds/opensds/pkg/utils"
	"github.com/opensds/opensds/pkg/utils/config"
	uuid "github.com/satori/go.uuid"
)

type Driver struct {
	conf   *DoradoConfig
	client *DoradoClient
}

func (d *Driver) Setup() (err error) {
	// Read huawei dorado config file
	conf := &DoradoConfig{}
	d.conf = conf
	path := config.CONF.OsdsDock.Backends.HuaweiDorado.ConfigPath

	if "" == path {
		path = defaultConfPath
	}
	Parse(conf, path)
	d.client, err = NewClient(&d.conf.AuthOptions)
	if err != nil {
		log.Errorf("Get new client failed, %v", err)
		return err
	}
	return nil
}

func (d *Driver) Unset() error {
	d.client.logout()
	return nil
}

func (d *Driver) createVolumeFromSnapshot(opt *pb.CreateVolumeOpts) (*model.VolumeSpec, error) {
	metadata := opt.GetMetadata()
	if metadata["hypermetro"] == "true" && metadata["replication_enabled"] == "true" {
		msg := "Hypermetro and Replication can not be used in the same volume_type"
		log.Error(msg)
		return nil, errors.New(msg)
	}
	snapshot, e1 := d.client.GetSnapshotByName(EncodeName(opt.GetSnapshotId()))
	if e1 != nil {
		log.Infof("Get Snapshot failed : %v", e1)
		return nil, e1
	}
	volumeDesc := TruncateDescription(opt.GetDescription())
	poolId, err1 := d.client.GetPoolIdByName(opt.GetPoolName())
	if err1 != nil {
		return nil, err1
	}

	lun, err := d.client.CreateVolume(EncodeName(opt.GetId()), opt.GetSize(),
		volumeDesc, poolId)
	if err != nil {
		log.Error("Create Volume Failed:", err)
		return nil, err
	}

	log.Infof("Create Volume from snapshot, source_lun_id : %s , target_lun_id : %s", snapshot.Id, lun.Id)
	err = utils.WaitForCondition(func() (bool, error) {
		getVolumeResult, getVolumeErr := d.client.GetVolume(lun.Id)
		if nil == getVolumeErr {
			if getVolumeResult.HealthStatus == StatusHealth && getVolumeResult.RunningStatus == StatusVolumeReady {
				return true, nil
			}
			log.V(5).Infof("Current lun HealthStatus : %s , RunningStatus : %s",
				getVolumeResult.HealthStatus, getVolumeResult.RunningStatus)
			return false, nil
		}
		return false, getVolumeErr

	}, LunReadyWaitInterval, LunReadyWaitTimeout)

	if err != nil {
		log.Error(err)
		d.client.DeleteVolume(lun.Id)
		return nil, err
	}
	err = d.copyVolume(opt, snapshot.Id, lun.Id)
	if err != nil {
		d.client.DeleteVolume(lun.Id)
		return nil, err
	}
	return &model.VolumeSpec{
		BaseModel: &model.BaseModel{
			Id: opt.GetId(),
		},
		Name:             opt.GetName(),
		Size:             Sector2Gb(lun.Capacity),
		Description:      volumeDesc,
		AvailabilityZone: opt.GetAvailabilityZone(),
		Metadata: map[string]string{
			KLunId: lun.Id,
		},
	}, nil

}
func (d *Driver) copyVolume(opt *pb.CreateVolumeOpts, srcid, tgtid string) error {
	metadata := opt.GetMetadata()
	copyspeed := metadata["copyspeed"]
	luncopyid, err := d.client.CreateLunCopy(EncodeName(opt.GetId()), srcid,
		tgtid, copyspeed)

	if err != nil {
		log.Error("Create Lun Copy failed,", err)
		return err
	}

	err = d.client.StartLunCopy(luncopyid)
	if err != nil {
		log.Errorf("Start lun: %s copy failed :%v,", luncopyid, err)
		d.client.DeleteLunCopy(luncopyid)
		return err
	}

	err = utils.WaitForCondition(func() (bool, error) {
		deleteLunCopyErr := d.client.DeleteLunCopy(luncopyid)
		if nil == deleteLunCopyErr {
			return true, nil
		}

		return false, nil
	}, LunCopyWaitInterval, LunCopyWaitTimeout)

	if err != nil {
		log.Error(err)
		return err
	}

	log.Infof("Copy Volume %s success", tgtid)
	return nil
}

func (d *Driver) CreateVolume(opt *pb.CreateVolumeOpts) (*model.VolumeSpec, error) {
	if opt.GetSnapshotId() != "" {
		return d.createVolumeFromSnapshot(opt)
	}
	name := EncodeName(opt.GetId())
	desc := TruncateDescription(opt.GetDescription())
	poolId, err := d.client.GetPoolIdByName(opt.GetPoolName())
	if err != nil {
		return nil, err
	}
	lun, err := d.client.CreateVolume(name, opt.GetSize(), desc, poolId)
	if err != nil {
		log.Error("Create Volume Failed:", err)
		return nil, err
	}
	log.Infof("Create volume %s (%s) success.", opt.GetName(), lun.Id)
	return &model.VolumeSpec{
		BaseModel: &model.BaseModel{
			Id: opt.GetId(),
		},
		Name:             opt.GetName(),
		Size:             Sector2Gb(lun.Capacity),
		Description:      opt.GetDescription(),
		AvailabilityZone: opt.GetAvailabilityZone(),
		Metadata: map[string]string{
			KLunId: lun.Id,
		},
	}, nil
}

func (d *Driver) PullVolume(volID string) (*model.VolumeSpec, error) {
	name := EncodeName(volID)
	lun, err := d.client.GetVolumeByName(name)
	if err != nil {
		return nil, err
	}

	return &model.VolumeSpec{
		BaseModel: &model.BaseModel{
			Id: volID,
		},
		Size:             Sector2Gb(lun.Capacity),
		Description:      lun.Description,
		AvailabilityZone: lun.ParentName,
	}, nil
}

func (d *Driver) DeleteVolume(opt *pb.DeleteVolumeOpts) error {
	lunId := opt.GetMetadata()[KLunId]
	err := d.client.DeleteVolume(lunId)
	if err != nil {
		log.Errorf("Delete volume failed, volume id =%s , Error:%s", opt.GetId(), err)
		return err
	}
	log.Info("Remove volume success, volume id =", opt.GetId())
	return nil
}

// ExtendVolume ...
func (d *Driver) ExtendVolume(opt *pb.ExtendVolumeOpts) (*model.VolumeSpec, error) {
	lunId := opt.GetMetadata()[KLunId]
	err := d.client.ExtendVolume(opt.GetSize(), lunId)
	if err != nil {
		log.Error("Extend Volume Failed:", err)
		return nil, err
	}

	log.Infof("Extend volume %s (%s) success.", opt.GetName(), opt.GetId())
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

func (d *Driver) getTargetInfo() (string, string, error) {
	tgtIp := d.conf.TargetIp
	resp, err := d.client.ListTgtPort()
	if err != nil {
		return "", "", err
	}
	for _, itp := range resp.Data {
		items := strings.Split(itp.Id, ",")
		iqn := strings.Split(items[0], "+")[1]
		items = strings.Split(iqn, ":")
		ip := items[len(items)-1]
		if tgtIp == ip {
			return iqn, ip, nil
		}
	}
	msg := fmt.Sprintf("Not find configuration targetIp: %v in device", tgtIp)
	return "", "", errors.New(msg)
}

func (d *Driver) InitializeConnection(opt *pb.CreateVolumeAttachmentOpts) (*model.ConnectionInfo, error) {
	if opt.GetAccessProtocol() == ISCSIProtocol {
		return d.InitializeConnectionIscsi(opt)
	}
	if opt.GetAccessProtocol() == FCProtocol {
		return d.InitializeConnectionFC(opt)
	}
	return nil, fmt.Errorf("not supported protocol type: %s", opt.GetAccessProtocol())
}

func (d *Driver) InitializeConnectionIscsi(opt *pb.CreateVolumeAttachmentOpts) (*model.ConnectionInfo, error) {

	lunId := opt.GetMetadata()[KLunId]
	hostInfo := opt.GetHostInfo()
	// Create host if not exist.
	hostId, err := d.client.AddHostWithCheck(hostInfo)
	if err != nil {
		log.Errorf("Add host failed, host name =%s, error: %v", hostInfo.Host, err)
		return nil, err
	}

	// Add initiator to the host.
	if err = d.client.AddInitiatorToHostWithCheck(hostId, hostInfo.Initiator); err != nil {
		log.Errorf("Add initiator to host failed, host id=%s, initiator=%s, error: %v", hostId, hostInfo.Initiator, err)
		return nil, err
	}

	// Add host to hostgroup.
	hostGrpId, err := d.client.AddHostToHostGroup(hostId)
	if err != nil {
		log.Errorf("Add host to group failed, host id=%s, error: %v", hostId, err)
		return nil, err
	}

	// Mapping lungroup and hostgroup to view.
	if err = d.client.DoMapping(lunId, hostGrpId, hostId); err != nil {
		log.Errorf("Do mapping failed, lun id=%s, hostGrpId=%s, hostId=%s, error: %v",
			lunId, hostGrpId, hostId, err)
		return nil, err
	}

	tgtIqn, tgtIp, err := d.getTargetInfo()
	if err != nil {
		log.Error("Get the target info failed,", err)
		return nil, err
	}
	tgtLun, err := d.client.GetHostLunId(hostId, lunId)
	if err != nil {
		log.Error("Get the get host lun id failed,", err)
		return nil, err
	}
	connInfo := &model.ConnectionInfo{
		DriverVolumeType: opt.GetAccessProtocol(),
		ConnectionData: map[string]interface{}{
			"targetDiscovered": true,
			"targetIQN":        []string{tgtIqn},
			"targetPortal":     []string{tgtIp + ":3260"},
			"discard":          false,
			"targetLun":        tgtLun,
		},
	}
	return connInfo, nil
}

func (d *Driver) TerminateConnection(opt *pb.DeleteVolumeAttachmentOpts) error {
	if opt.GetAccessProtocol() == ISCSIProtocol {
		return d.TerminateConnectionIscsi(opt)
	}
	if opt.GetAccessProtocol() == FCProtocol {
		return d.TerminateConnectionFC(opt)
	}
	return fmt.Errorf("not supported protocal type: %s", opt.GetAccessProtocol())
}

func (d *Driver) TerminateConnectionIscsi(opt *pb.DeleteVolumeAttachmentOpts) error {
	hostId, err := d.client.GetHostIdByName(opt.GetHostInfo().GetHost())
	if err != nil {
		// host id has been delete already, ignore the host not found error
		if IsNotFoundError(err) {
			log.Warningf("host(%s) has been removed already, ignore it. "+
				"Delete volume attachment(%s)success.", hostId, opt.GetId())
			return nil
		}
		return err
	}
	// the name format of there objects blow is: xxxPrefix + hostId
	// the empty xxId means that the specified object has been removed already.
	lunGrpId, err := d.client.FindLunGroup(PrefixLunGroup + hostId)
	if err != nil && !IsNotFoundError(err) {
		return err
	}
	hostGrpId, err := d.client.FindHostGroup(PrefixHostGroup + hostId)
	if err != nil && !IsNotFoundError(err) {
		return err
	}
	viewId, err := d.client.FindMappingView(PrefixMappingView + hostId)
	if err != nil && !IsNotFoundError(err) {
		return err
	}

	lunId := opt.GetMetadata()[KLunId]
	if lunGrpId != "" {
		if d.client.IsLunGroupContainLun(lunGrpId, lunId) {
			if err := d.client.RemoveLunFromLunGroup(lunGrpId, lunId); err != nil {
				return err
			}
		}

		//  if lun group still contains other lun(s), ignore the all the operations blow,
		// and goes back with success status.
		var leftObjectCount = 0
		if leftObjectCount, err = d.client.getObjectCountFromLungroup(lunGrpId); err != nil {
			return err
		}
		if leftObjectCount > 0 {
			log.Infof("Lun group(%s) still contains %d lun(s). "+
				"Delete volume attachment(%s)success.", lunGrpId, leftObjectCount, opt.GetId())
			return nil
		}
	}

	if viewId != "" {
		if d.client.IsMappingViewContainLunGroup(viewId, lunGrpId) {
			if err := d.client.RemoveLunGroupFromMappingView(viewId, lunGrpId); err != nil {
				return err
			}
		}
		if d.client.IsMappingViewContainHostGroup(viewId, hostGrpId) {
			if err := d.client.RemoveHostGroupFromMappingView(viewId, hostGrpId); err != nil {
				return err
			}
		}
		if err := d.client.DeleteMappingView(viewId); err != nil {
			return err
		}
	}

	if lunGrpId != "" {
		if err := d.client.DeleteLunGroup(lunGrpId); err != nil {
			return err
		}
	}

	if hostGrpId != "" {
		if d.client.IsHostGroupContainHost(hostGrpId, hostId) {
			if err := d.client.RemoveHostFromHostGroup(hostGrpId, hostId); err != nil {
				return err
			}
		}
		if err := d.client.DeleteHostGroup(hostGrpId); err != nil {
			return err
		}
	}

	initiatorName := opt.GetHostInfo().GetInitiator()
	if d.client.IsHostContainInitiator(hostId, initiatorName) {
		if err := d.client.RemoveIscsiFromHost(initiatorName); err != nil {
			return err
		}
	}

	fcExist, err := d.client.checkFCInitiatorsExistInHost(hostId)
	if err != nil {
		return err
	}
	iscsiExist, err := d.client.checkIscsiInitiatorsExistInHost(hostId)
	if err != nil {
		return err
	}
	if fcExist || iscsiExist {
		log.Warningf("host (%s) still contains initiator(s), ignore delete it. "+
			"Delete volume attachment(%s)success.", hostId, opt.GetId())
		return nil
	}

	if err := d.client.DeleteHost(hostId); err != nil {
		return err
	}
	log.Infof("Delete volume attachment(%s)success.", opt.GetId())
	return nil
}

func (d *Driver) CreateSnapshot(opt *pb.CreateVolumeSnapshotOpts) (*model.VolumeSnapshotSpec, error) {
	lunId := opt.GetMetadata()[KLunId]
	name := EncodeName(opt.GetId())
	desc := TruncateDescription(opt.GetDescription())
	snap, err := d.client.CreateSnapshot(lunId, name, desc)
	if err != nil {
		return nil, err
	}
	return &model.VolumeSnapshotSpec{
		BaseModel: &model.BaseModel{
			Id: opt.GetId(),
		},
		Name:        opt.GetName(),
		Description: opt.GetDescription(),
		VolumeId:    opt.GetVolumeId(),
		Size:        0,
		Metadata: map[string]string{
			KSnapId: snap.Id,
		},
	}, nil
}

func (d *Driver) PullSnapshot(id string) (*model.VolumeSnapshotSpec, error) {
	name := EncodeName(id)
	snap, err := d.client.GetSnapshotByName(name)
	if err != nil {
		return nil, err
	}
	return &model.VolumeSnapshotSpec{
		BaseModel: &model.BaseModel{
			Id: snap.Id,
		},
		Name:        snap.Name,
		Description: snap.Description,
		Size:        0,
		VolumeId:    snap.ParentId,
	}, nil
}

func (d *Driver) DeleteSnapshot(opt *pb.DeleteVolumeSnapshotOpts) error {
	id := opt.GetMetadata()[KSnapId]
	err := d.client.DeleteSnapshot(id)
	if err != nil {
		log.Errorf("Delete volume snapshot failed, volume snapshot id = %s , error: %v", opt.GetId(), err)
		return err
	}
	log.Info("Remove volume snapshot success, volume snapshot id =", opt.GetId())
	return nil
}

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
		name := fmt.Sprintf("%s:%s:%s", host, d.conf.Endpoints, p.Id)
		pol := &model.StoragePoolSpec{
			BaseModel: &model.BaseModel{
				Id: uuid.NewV5(uuid.NamespaceOID, name).String(),
			},
			Name:             p.Name,
			TotalCapacity:    Sector2Gb(p.UserTotalCapacity),
			FreeCapacity:     Sector2Gb(p.UserFreeCapacity),
			StorageType:      c.Pool[p.Name].StorageType,
			Extras:           c.Pool[p.Name].Extras,
			AvailabilityZone: c.Pool[p.Name].AvailabilityZone,
			MultiAttach:      c.Pool[p.Name].MultiAttach,
		}
		if pol.AvailabilityZone == "" {
			pol.AvailabilityZone = defaultAZ
		}
		pols = append(pols, pol)
	}
	return pols, nil
}

func (d *Driver) InitializeConnectionFC(opt *pb.CreateVolumeAttachmentOpts) (*model.ConnectionInfo, error) {
	lunId := opt.GetMetadata()[KLunId]
	hostInfo := opt.GetHostInfo()
	// Create host if not exist.
	hostId, err := d.client.AddHostWithCheck(hostInfo)
	if err != nil {
		log.Errorf("Add host failed, host name =%s, error: %v", hostInfo.Host, err)
		return nil, err
	}

	// Add host to hostgroup.
	hostGrpId, err := d.client.AddHostToHostGroup(hostId)
	if err != nil {
		log.Errorf("Add host to group failed, host id=%s, error: %v", hostId, err)
		return nil, err
	}

	// Not use FC switch
	tgtPortWWNs, initTargMap, err := d.connectFCUseNoSwitch(opt, opt.GetHostInfo().GetInitiator(), hostId)
	if err != nil {
		return nil, err
	}

	// Mapping lungroup and hostgroup to view.
	if err = d.client.DoMapping(lunId, hostGrpId, hostId); err != nil {
		log.Errorf("Do mapping failed, lun id=%s, hostGrpId=%s, hostId=%s, error: %v",
			lunId, hostGrpId, hostId, err)
		return nil, err
	}

	tgtLun, err := d.client.GetHostLunId(hostId, lunId)
	if err != nil {
		log.Error("Get the get host lun id failed,", err)
		return nil, err
	}

	fcInfo := &model.ConnectionInfo{
		DriverVolumeType: opt.GetAccessProtocol(),
		ConnectionData: map[string]interface{}{
			"targetDiscovered":     true,
			"target_wwn":           tgtPortWWNs,
			"volume_id":            opt.GetVolumeId(),
			"initiator_target_map": initTargMap,
			"description":          "huawei",
			"host_name":            opt.GetHostInfo().Host,
			"target_lun":           tgtLun,
		},
	}
	return fcInfo, nil
}

func (d *Driver) connectFCUseNoSwitch(opt *pb.CreateVolumeAttachmentOpts, wwpns string, hostId string) ([]string, map[string][]string, error) {
	wwns := strings.Split(wwpns, ",")

	onlineWWNsInHost, err := d.client.GetHostOnlineFCInitiators(hostId)
	if err != nil {
		return nil, nil, err
	}
	onlineFreeWWNs, err := d.client.GetOnlineFreeWWNs()
	if err != nil {
		return nil, nil, err
	}
	onlineFCInitiators, err := d.client.GetOnlineFCInitiatorOnArray()
	if err != nil {
		return nil, nil, err
	}

	var wwnsNew []string
	for _, w := range wwns {
		if d.isInStringArray(w, onlineFCInitiators) {
			wwnsNew = append(wwnsNew, w)
		}
	}
	log.Infof("initialize connection, online initiators on the array:%s", wwnsNew)

	if wwnsNew == nil {
		return nil, nil, errors.New("no available host initiator")
	}

	for _, wwn := range wwnsNew {
		if !d.isInStringArray(wwn, onlineWWNsInHost) && !d.isInStringArray(wwn, onlineFreeWWNs) {
			wwnsInHost, err := d.client.GetHostFCInitiators(hostId)
			if err != nil {
				return nil, nil, err
			}
			iqnsInHost, err := d.client.GetHostIscsiInitiators(hostId)
			if err != nil {
				return nil, nil, err
			}
			flag, err := d.client.IsHostAssociatedToHostgroup(hostId)
			if err != nil {
				return nil, nil, err
			}

			if wwnsInHost == nil && iqnsInHost == nil && flag == false {
				if err = d.client.RemoveHost(hostId); err != nil {
					return nil, nil, err
				}
			}

			msg := fmt.Sprintf("host initiator occupied: Can not add FC initiator %s to host %s, please check if this initiator has been added to other host.", wwn, hostId)
			log.Errorf(msg)
			return nil, nil, errors.New(msg)
		}
	}

	for _, wwn := range wwnsNew {
		if d.isInStringArray(wwn, onlineFreeWWNs) {
			if err = d.client.AddFCPortTohost(hostId, wwn); err != nil {
				return nil, nil, err
			}
		}
	}

	tgtPortWWNs, initTargMap, err := d.client.GetIniTargMap(wwnsNew)
	if err != nil {
		return nil, nil, err
	}

	return tgtPortWWNs, initTargMap, nil

}

func (d *Driver) isInStringArray(s string, source []string) bool {
	for _, i := range source {
		if s == i {
			return true
		}
	}
	return false
}

func (d *Driver) TerminateConnectionFC(opt *pb.DeleteVolumeAttachmentOpts) error {
	// Detach lun
	fcInfo, err := d.detachVolumeFC(opt)
	if err != nil {
		return err
	}
	log.Info(fmt.Sprintf("terminate connection fc, return data is: %s", fcInfo))
	return nil
}

func (d *Driver) detachVolumeFC(opt *pb.DeleteVolumeAttachmentOpts) (string, error) {
	wwns := strings.Split(opt.GetHostInfo().GetInitiator(), ",")
	lunId := opt.GetMetadata()[KLunId]

	log.Infof("terminate connection, wwpns: %s,lun id: %s", wwns, lunId)

	hostId, lunGrpId, hostGrpId, viewId, err := d.getMappedInfo(opt.GetHostInfo().GetHost())
	if err != nil {
		return "", err
	}

	if lunId != "" && lunGrpId != "" {
		if err := d.client.RemoveLunFromLunGroup(lunGrpId, lunId); err != nil {
			return "", err
		}
	}

	var leftObjectCount = -1
	if lunGrpId != "" {
		if leftObjectCount, err = d.client.getObjectCountFromLungroup(lunGrpId); err != nil {
			return "", err
		}
	}

	var fcInfo string
	if leftObjectCount > 0 {
		fcInfo = "driver_volume_type: fibre_channel, data: {}"
	} else {
		if fcInfo, err = d.deleteZoneAndRemoveFCInitiators(wwns, hostId, hostGrpId, viewId); err != nil {
			return "", err
		}

		if err := d.clearHostRelatedResource(lunGrpId, viewId, hostId, hostGrpId); err != nil {
			return "", err
		}
	}

	log.Info(fmt.Sprintf("Return target backend FC info is: %s", fcInfo))
	return fcInfo, nil
}

func (d *Driver) deleteZoneAndRemoveFCInitiators(wwns []string, hostId, hostGrpId, viewId string) (string, error) {
	tgtPortWWNs, initTargMap, err := d.client.GetIniTargMap(wwns)
	if err != nil {
		return "", err
	}

	// Remove the initiators from host if need.
	hostGroupNum, err := d.client.getHostGroupNumFromHost(hostId)
	if err != nil {
		return "", err
	}
	if hostGrpId != "" && hostGroupNum <= 1 || (hostGrpId == "" && hostGroupNum <= 0) {
		fcInitiators, err := d.client.GetHostFCInitiators(hostId)
		if err != nil {
			return "", err
		}
		for _, wwn := range wwns {
			if d.isInStringArray(wwn, fcInitiators) {
				if err := d.client.removeFCFromHost(wwn); err != nil {
					return "", err
				}
			}
		}
	}

	return fmt.Sprintf("driver_volume_type: fibre_channel, target_wwn: %s, initiator_target_map: %s", tgtPortWWNs, initTargMap), nil
}

func (d *Driver) getMappedInfo(hostName string) (string, string, string, string, error) {
	hostId, err := d.client.GetHostIdByName(hostName)
	if err != nil {
		return "", "", "", "", err
	}

	lunGrpId, err := d.client.FindLunGroup(PrefixLunGroup + hostId)
	if err != nil {
		return "", "", "", "", err
	}
	hostGrpId, err := d.client.FindHostGroup(PrefixHostGroup + hostId)
	if err != nil {
		return "", "", "", "", err
	}
	viewId, err := d.client.FindMappingView(PrefixMappingView + hostId)
	if err != nil {
		return "", "", "", "", err
	}

	return hostId, lunGrpId, hostGrpId, viewId, nil
}

func (d *Driver) clearHostRelatedResource(lunGrpId, viewId, hostId, hostGrpId string) error {
	if lunGrpId != "" {
		if viewId != "" {
			d.client.RemoveLunGroupFromMappingView(viewId, lunGrpId)
		}
		d.client.DeleteLunGroup(lunGrpId)
	}
	if hostId != "" {
		if hostGrpId != "" {

			if viewId != "" {
				d.client.RemoveHostGroupFromMappingView(viewId, hostGrpId)
			}

			views, err := d.client.getHostgroupAssociatedViews(hostGrpId)
			if err != nil {
				return err
			}

			if len(views) <= 0 {
				if err := d.client.RemoveHostFromHostGroup(hostGrpId, hostId); err != nil {
					return err
				}
				hosts, err := d.client.getHostsInHostgroup(hostGrpId)
				if err != nil {
					return err
				}

				if len(hosts) <= 0 {
					if err := d.client.DeleteHostGroup(hostGrpId); err != nil {
						return err
					}
				}
			}
		}

		flag, err := d.client.checkFCInitiatorsExistInHost(hostId)
		if err != nil {
			return err
		}
		if !flag {
			if err := d.client.RemoveHost(hostId); err != nil {
				return err
			}
		}
	}

	if viewId != "" {
		if err := d.client.DeleteMappingView(viewId); err != nil {
			return err
		}
	}

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
