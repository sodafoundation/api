// Copyright (c) 2017 Huawei Technologies Co., Ltd. All Rights Reserved.
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

package dorado

import (
	"errors"
	"fmt"
	"os"
	"strings"

	log "github.com/golang/glog"
	. "github.com/opensds/opensds/contrib/drivers/utils/config"
	pb "github.com/opensds/opensds/pkg/dock/proto"
	"github.com/opensds/opensds/pkg/model"
	"github.com/opensds/opensds/pkg/utils/config"
	"github.com/satori/go.uuid"
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

	lun, err := d.client.CreateVolume(opt.GetName(), opt.GetSize(), volumeDesc, poolId)
	if err != nil {
		log.Error("Create Volume Failed:", err)
		return nil, err
	}

	log.Infof("Create Volume from snapshot, source_lun_id : %s , target_lun_id : %s", snapshot.Id, lun.Id)

	err = WaitForCondition(func() (bool, error) {
		if lun.HealthStatus == StatusHealth && lun.RunningStatus == StatusVolumeReady {
			return true, nil
		} else {
			msg := fmt.Sprintf("Volume state is not mathch, lun ID : %s , HealthStatus : %s,RunningStatus : %s",
				lun.Id, lun.HealthStatus, lun.RunningStatus)
			return false, errors.New(msg)
		}
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
	luncopyid, err := d.client.CreateLunCopy(opt.GetName(), srcid, tgtid, copyspeed)

	if err != nil {
		log.Error("Create Lun Copy failed,", err)
		return err
	}
	defer d.client.DeleteLunCopy(luncopyid)
	err = d.client.StartLunCopy(luncopyid)
	if err != nil {
		log.Errorf("Start lun: %s copy failed :%v,", luncopyid, err)
		return err
	}
	lunCopyInfo, err1 := d.client.GetLunInfo(luncopyid)
	if err1 != nil {
		log.Errorf("Get lun info failed :%v", err1)
		return err1
	}
	err = WaitForCondition(func() (bool, error) {
		if lunCopyInfo.RunningStatus == StatusLuncopyReady || lunCopyInfo.RunningStatus == StatusLunCoping {
			return true, nil
		} else if lunCopyInfo.HealthStatus != StatusHealth {
			msg := fmt.Sprintf("An error occurred during the luncopy operation. Lun name : %s  ,Lun copy health status : %s ,Lun copy running status : %s ",
				lunCopyInfo.Name, lunCopyInfo.HealthStatus, lunCopyInfo.RunningStatus)
			return false, errors.New(msg)
		}
		return true, nil
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

func (d *Driver) InitializeConnection(opt *pb.CreateAttachmentOpts) (*model.ConnectionInfo, error) {
	if opt.GetAccessProtocol() == ISCSIProtocol {
		return d.InitializeConnectionIscsi(opt)
	}
	if opt.GetAccessProtocol() == FCProtocol {
		return d.InitializeConnectionFC(opt)
	}
	return nil, errors.New("No supported protocol for dorado driver.")
}

func (d *Driver) InitializeConnectionIscsi(opt *pb.CreateAttachmentOpts) (*model.ConnectionInfo, error) {

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
		DriverVolumeType: ISCSIProtocol,
		ConnectionData: map[string]interface{}{
			"targetDiscovered": true,
			"targetIQN":        tgtIqn,
			"targetPortal":     tgtIp + ":3260",
			"discard":          false,
			"targetLun":        tgtLun,
		},
	}
	return connInfo, nil
}

func (d *Driver) TerminateConnection(opt *pb.DeleteAttachmentOpts) error {
	if opt.GetAccessProtocol() == ISCSIProtocol {
		return d.TerminateConnectionIscsi(opt)
	}
	if opt.GetAccessProtocol() == FCProtocol {
		return d.TerminateConnectionFC(opt)
	}
	return nil
}

func (d *Driver) TerminateConnectionIscsi(opt *pb.DeleteAttachmentOpts) error {
	lunId := opt.GetMetadata()[KLunId]
	hostId, err := d.client.GetHostIdByName(opt.GetHostInfo().GetHost())
	if err != nil {
		return err
	}
	lunGrpId, _ := d.client.FindLunGroup(LunGroupPrefix + hostId)
	hostGrpId, _ := d.client.FindHostGroup(HostGroupPrefix + hostId)
	viewId, _ := d.client.FindMappingView(MappingViewPrefix + hostId)
	if viewId != "" {
		d.client.RemoveLunGroupFromMappingView(viewId, lunGrpId)
		d.client.RemoveHostGroupFromMappingView(viewId, hostGrpId)
		d.client.DeleteMappingView(viewId)
	}
	if hostGrpId != "" {
		d.client.RemoveHostFromHostGroup(hostGrpId, hostId)
		d.client.DeleteHostGroup(hostGrpId)
	}
	if lunGrpId != "" {
		d.client.RemoveLunFromLunGroup(lunGrpId, lunId)
		d.client.DeleteLunGroup(lunGrpId)
	}
	d.client.RemoveIscsiFromHost(opt.GetHostInfo().GetInitiator())
	d.client.DeleteHost(hostId)
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
		}
		if pol.AvailabilityZone == "" {
			pol.AvailabilityZone = defaultAZ
		}
		pols = append(pols, pol)
	}
	return pols, nil
}

func (d *Driver) InitializeConnectionFC(opt *pb.CreateAttachmentOpts) (*model.ConnectionInfo, error) {
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
		DriverVolumeType: FCProtocol,
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

func (d *Driver) connectFCUseNoSwitch(opt *pb.CreateAttachmentOpts, wwpns string, hostId string) ([]string, map[string][]string, error) {
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

func (d *Driver) TerminateConnectionFC(opt *pb.DeleteAttachmentOpts) error {
	// Detach lun
	fcInfo, err := d.detachVolumeFC(opt)
	if err != nil {
		return err
	}
	log.Info(fmt.Sprintf("terminate connection fc, return data is: %s", fcInfo))
	return nil
}

func (d *Driver) detachVolumeFC(opt *pb.DeleteAttachmentOpts) (string, error) {
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

	lunGrpId, err := d.client.FindLunGroup(LunGroupPrefix + hostId)
	if err != nil {
		return "", "", "", "", err
	}
	hostGrpId, err := d.client.FindHostGroup(HostGroupPrefix + hostId)
	if err != nil {
		return "", "", "", "", err
	}
	viewId, err := d.client.FindMappingView(MappingViewPrefix + hostId)
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
	return nil, &model.NotImplementError{S: "Method InitializeSnapshotConnection has not been implemented yet."}
}

func (d *Driver) TerminateSnapshotConnection(opt *pb.DeleteSnapshotAttachmentOpts) error {
	return &model.NotImplementError{S: "Method TerminateSnapshotConnection has not been implemented yet."}
}

func (d *Driver) CreateVolumeGroup(opt *pb.CreateVolumeGroupOpts, vg *model.VolumeGroupSpec) (*model.VolumeGroupSpec, error) {
	return nil, &model.NotImplementError{"Method CreateVolumeGroup did not implement."}
}

func (d *Driver) UpdateVolumeGroup(opt *pb.UpdateVolumeGroupOpts, vg *model.VolumeGroupSpec, addVolumesRef []*model.VolumeSpec, removeVolumesRef []*model.VolumeSpec) (*model.VolumeGroupSpec, []*model.VolumeSpec, []*model.VolumeSpec, error) {
	return nil, nil, nil, &model.NotImplementError{"Method UpdateVolumeGroup did not implement."}
}

func (d *Driver) DeleteVolumeGroup(opt *pb.DeleteVolumeGroupOpts, vg *model.VolumeGroupSpec, volumes []*model.VolumeSpec) (*model.VolumeGroupSpec, []*model.VolumeSpec, error) {
	return nil, nil, &model.NotImplementError{"Method UpdateVolumeGroup did not implement."}
}
