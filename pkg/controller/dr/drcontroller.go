// Copyright 2018 The OpenSDS Authors.
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

package dr

import (
	"encoding/json"
	"fmt"
	"strings"

	log "github.com/golang/glog"
	"github.com/opensds/opensds/contrib/drivers/utils/config"
	c "github.com/opensds/opensds/pkg/context"
	"github.com/opensds/opensds/pkg/controller/volume"
	"github.com/opensds/opensds/pkg/db"
	. "github.com/opensds/opensds/pkg/model"
	pb "github.com/opensds/opensds/pkg/model/proto"
	"github.com/opensds/opensds/pkg/utils"
	uuid "github.com/satori/go.uuid"
)

type Controller interface {
	CreateReplication(ctx *c.Context, replica *ReplicationSpec, primaryVol, secondaryVol *VolumeSpec) (*ReplicationSpec, error)
	DeleteReplication(ctx *c.Context, replica *ReplicationSpec, primaryVol, secondaryVol *VolumeSpec) error
	EnableReplication(ctx *c.Context, replica *ReplicationSpec, primaryVol, secondaryVol *VolumeSpec) error
	DisableReplication(ctx *c.Context, replica *ReplicationSpec, primaryVol, secondaryVol *VolumeSpec) error
	FailoverReplication(ctx *c.Context, replica *ReplicationSpec, failover *FailoverReplicationSpec, primaryVol, secondaryVol *VolumeSpec) error
}

type DrController struct {
	volumeController volume.Controller
	primaryOp        ReplicationOperator
	secondaryOp      ReplicationOperator
}

// NewController method creates a controller structure and expose its pointer.
func NewController(controller volume.Controller) Controller {
	return &DrController{
		volumeController: controller,
	}
}

func (d *DrController) getVolumeDataList(ctx *c.Context) ([]*pb.VolumeData, error) {
	volList, err := db.C.ListVolumes(ctx)
	if err != nil {
		return nil, err
	}
	var volumeDataList []*pb.VolumeData
	for _, v := range volList {
		if v.ReplicationDriverData != nil && len(v.ReplicationDriverData) != 0 {
			v.ReplicationDriverData["VolumeId"] = v.Id
			volumeDataList = append(volumeDataList, &pb.VolumeData{Data: v.ReplicationDriverData})
		}
	}
	return volumeDataList, nil
}

func (d *DrController) LoadOperator(ctx *c.Context, primaryVol, secondaryVol *VolumeSpec) error {
	var err error
	d.primaryOp, err = NewPairOperator(ctx, d.volumeController, primaryVol, true)
	if err != nil {
		return err
	}
	d.secondaryOp, err = NewPairOperator(ctx, d.volumeController, secondaryVol, false)
	if err != nil {
		return err
	}
	return nil
}

func (d *DrController) CreateReplication(ctx *c.Context, replica *ReplicationSpec, primaryVol,
	secondaryVol *VolumeSpec) (*ReplicationSpec, error) {
	if primaryVol.Size != secondaryVol.Size {
		return replica, fmt.Errorf("secondary volume size(%d) is not the same as the primary size(%d)",
			secondaryVol.Size, primaryVol.Size)
	}
	pPool, _ := db.C.GetPool(ctx, primaryVol.PoolId)
	sPool, _ := db.C.GetPool(ctx, secondaryVol.PoolId)
	if pPool.ReplicationType != sPool.ReplicationType {
		return replica, fmt.Errorf("secondary replication type is not the same as the primary")
	}

	replica.PrimaryReplicationDriverData = utils.MergeStringMaps(replica.PrimaryReplicationDriverData, primaryVol.Metadata)
	replica.SecondaryReplicationDriverData = utils.MergeStringMaps(replica.SecondaryReplicationDriverData, secondaryVol.Metadata)

	// Load replication operator
	d.LoadOperator(ctx, primaryVol, secondaryVol)

	// Host-Based replication needs to do some extra operations including attaching volume and initializing volume data list
	if pPool.ReplicationType == ReplicationTypeHost {
		var err error
		replica.VolumeDataList, err = d.getVolumeDataList(ctx)
		if err != nil {
			log.Errorf("Get volume data list failed, %s", err)
			return replica, err
		}
		replica, err = d.primaryOp.Attach(ctx, replica, primaryVol)
		if err != nil {
			log.Errorf("Attach primary volume failed, %s", err)
			return replica, err
		}
		replica, err = d.secondaryOp.Attach(ctx, replica, secondaryVol)
		if err != nil {
			log.Errorf("Attach secondary volume failed, %s", err)
			return replica, err
		}
	}
	// Do replication.
	pResult, err := d.primaryOp.Create(ctx, replica, primaryVol)
	if err != nil {
		log.Errorf("Create primary replication failed, %s", err)
		return replica, err
	}

	sResult, err := d.secondaryOp.Create(ctx, replica, secondaryVol)
	if err != nil {
		log.Errorf("Create secondary replication failed, %s", err)
		return replica, err
	}

	replica.PrimaryReplicationDriverData = utils.MergeStringMaps(replica.PrimaryReplicationDriverData, pResult.PrimaryReplicationDriverData)
	replica.SecondaryReplicationDriverData = utils.MergeStringMaps(sResult.SecondaryReplicationDriverData, replica.SecondaryReplicationDriverData)
	replica.Metadata = utils.MergeStringMaps(replica.Metadata, pResult.Metadata, sResult.Metadata)

	primaryVol.ReplicationDriverData = pResult.PrimaryReplicationDriverData
	if primaryVol.ReplicationDriverData != nil {
		primaryVol.ReplicationDriverData["IsPrimary"] = "true"
		db.C.UpdateVolume(ctx, primaryVol)

	}

	secondaryVol.ReplicationDriverData = sResult.SecondaryReplicationDriverData
	if secondaryVol.ReplicationDriverData != nil {
		secondaryVol.ReplicationDriverData["IsPrimary"] = "false"
		db.C.UpdateVolume(ctx, secondaryVol)
	}

	return replica, nil
}

func (d *DrController) DeleteReplication(ctx *c.Context, replica *ReplicationSpec, primaryVol, secondaryVol *VolumeSpec) error {
	d.LoadOperator(ctx, primaryVol, secondaryVol)
	err := d.primaryOp.Delete(ctx, replica, primaryVol)
	if err != nil {
		return err
	}
	err = d.secondaryOp.Delete(ctx, replica, secondaryVol)
	if err != nil {
		return err
	}

	pPool, err := db.C.GetPool(ctx, primaryVol.PoolId)
	if err != nil {
		return err
	}

	if pPool.ReplicationType == ReplicationTypeHost {
		var err error
		// dettach
		err = d.primaryOp.Detach(ctx, replica, primaryVol)
		if err != nil {
			log.Errorf("Detach primary volume failed, %s", err)
			return err
		}

		err = d.secondaryOp.Detach(ctx, replica, secondaryVol)
		if err != nil {
			log.Errorf("Detach secondary volume failed, %s", err)
			return err
		}
	}

	// clean up replication driver data in volume database
	primaryVol.ReplicationDriverData = make(map[string]string)
	db.C.UpdateVolume(ctx, primaryVol)
	secondaryVol.ReplicationDriverData = make(map[string]string)
	db.C.UpdateVolume(ctx, secondaryVol)
	return nil
}

func (d *DrController) EnableReplication(ctx *c.Context, replica *ReplicationSpec, primaryVol, secondaryVol *VolumeSpec) error {
	d.LoadOperator(ctx, primaryVol, secondaryVol)
	err := d.primaryOp.Enable(ctx, replica, primaryVol)
	if err != nil {
		return err
	}
	err = d.secondaryOp.Enable(ctx, replica, secondaryVol)
	if err != nil {
		return err
	}
	return nil
}

func (d *DrController) DisableReplication(ctx *c.Context, replica *ReplicationSpec, primaryVol, secondaryVol *VolumeSpec) error {
	d.LoadOperator(ctx, primaryVol, secondaryVol)
	err := d.primaryOp.Disable(ctx, replica, primaryVol)
	if err != nil {
		return err
	}
	err = d.secondaryOp.Disable(ctx, replica, secondaryVol)
	if err != nil {
		return err
	}
	return nil
}

func (d *DrController) FailoverReplication(ctx *c.Context, replica *ReplicationSpec,
	failover *FailoverReplicationSpec, primaryVol, secondaryVol *VolumeSpec) error {
	d.LoadOperator(ctx, primaryVol, secondaryVol)
	err := d.primaryOp.Failover(ctx, replica, failover, primaryVol)
	if err != nil {
		return err
	}
	err = d.secondaryOp.Failover(ctx, replica, failover, secondaryVol)
	if err != nil {
		return err
	}
	return nil
}

type ReplicationOperator interface {
	Create(ctx *c.Context, replica *ReplicationSpec, vol *VolumeSpec) (*ReplicationSpec, error)
	Delete(ctx *c.Context, replica *ReplicationSpec, vol *VolumeSpec) error
	Enable(ctx *c.Context, replica *ReplicationSpec, vol *VolumeSpec) error
	Disable(ctx *c.Context, replica *ReplicationSpec, vol *VolumeSpec) error
	Failover(ctx *c.Context, replica *ReplicationSpec, failover *FailoverReplicationSpec, vol *VolumeSpec) error
	Attach(ctx *c.Context, replica *ReplicationSpec, vol *VolumeSpec) (*ReplicationSpec, error)
	Detach(ctx *c.Context, replica *ReplicationSpec, vol *VolumeSpec) error
}

type PairOperator struct {
	volumeController volume.Controller
	isPrimary        bool
	pool             *StoragePoolSpec
	provisionDock    *DockSpec
}

func NewPairOperator(ctx *c.Context, controller volume.Controller, vol *VolumeSpec, isPrimary bool) (*PairOperator, error) {
	pool, err := db.C.GetPool(ctx, vol.PoolId)
	if err != nil {
		log.Error("get pool failed", err)
		return nil, err
	}

	provisionerDock, err := db.C.GetDockByPoolId(ctx, vol.PoolId)
	if err != nil {
		log.Error("When search dock in db by pool id", err)
		return nil, err
	}
	return &PairOperator{
		volumeController: controller,
		isPrimary:        isPrimary,
		pool:             pool,
		provisionDock:    provisionerDock,
	}, nil
}

func (p *PairOperator) getAttacherDockByProvisioner(ctx *c.Context, provisionerDock *DockSpec) (*DockSpec, error) {
	segments := strings.Split(provisionerDock.Endpoint, ":")
	endpointIp := segments[len(segments)-2]
	// Generate the attacher UUID by nodeid and endpoint ip.
	attacherDockId := uuid.NewV5(uuid.NamespaceOID, provisionerDock.NodeId+":"+endpointIp)
	return db.C.GetDock(ctx, attacherDockId.String())
}

func (p *PairOperator) doAttach(ctx *c.Context, vol *VolumeSpec, provisionerDock *DockSpec) (*VolumeAttachmentSpec, error) {
	attacherDock, err := p.getAttacherDockByProvisioner(ctx, provisionerDock)
	p.volumeController.SetDock(provisionerDock)
	attachmentId := uuid.NewV4().String()
	// Default protocol is iscsi
	protocol := config.ISCSIProtocol
	if len(p.pool.Extras.IOConnectivity.AccessProtocol) != 0 {
		protocol = p.pool.Extras.IOConnectivity.AccessProtocol
	}

	initiator := attacherDock.Metadata["Initiator"]
	if protocol == config.FCProtocol {
		initiator = attacherDock.Metadata["WWPNS"]
	}

	var createAttachOpt = &pb.CreateVolumeAttachmentOpts{
		Id:       attachmentId,
		VolumeId: vol.Id,
		HostInfo: &pb.HostInfo{
			Platform:  attacherDock.Metadata["Platform"],
			OsType:    attacherDock.Metadata["OsType"],
			Ip:        attacherDock.Metadata["HostIp"],
			Host:      attacherDock.NodeId,
			Initiator: initiator,
		},
		AccessProtocol: protocol,
		Metadata:       vol.Metadata,
		DriverName:     provisionerDock.DriverName,
		Context:        ctx.ToJson(),
	}

	atm, err := p.volumeController.CreateVolumeAttachment(createAttachOpt)
	if err != nil {
		log.Errorf("create attachment failed, %v", err)
		return nil, err
	}

	atm.Status = VolumeAvailable
	_, err = db.C.CreateVolumeAttachment(ctx, atm)
	if err != nil {
		return nil, err
	}

	rollback := false
	defer func() {
		if rollback {
			opt := &pb.DeleteVolumeAttachmentOpts{
				Id:       atm.Id,
				VolumeId: atm.VolumeId,
				HostInfo: &pb.HostInfo{
					Platform:  atm.Platform,
					OsType:    atm.OsType,
					Ip:        atm.Ip,
					Host:      atm.Host,
					Initiator: atm.Initiator,
				},
				AccessProtocol: protocol,
				Metadata:       utils.MergeStringMaps(atm.Metadata, vol.Metadata),
				DriverName:     provisionerDock.DriverName,
				Context:        ctx.ToJson(),
			}
			p.volumeController.SetDock(provisionerDock)
			p.volumeController.DeleteVolumeAttachment(opt)
			db.C.DeleteVolumeAttachment(ctx, atm.Id)
		}
	}()

	p.volumeController.SetDock(attacherDock)
	connData, _ := json.Marshal(atm.ConnectionData)
	var attachOpt = &pb.AttachVolumeOpts{
		AccessProtocol: atm.DriverVolumeType,
		ConnectionData: string(connData),
		Metadata:       map[string]string{},
		Context:        ctx.ToJson(),
	}
	mountPoint, err := p.volumeController.AttachVolume(attachOpt)
	if err != nil {
		rollback = true
		log.Errorf("attach volume failed, %v", err)
		return nil, err
	}

	atm.Mountpoint = mountPoint
	atm.AccessProtocol = atm.DriverVolumeType
	_, err = db.C.UpdateVolumeAttachment(ctx, atm.Id, atm)
	if err != nil {
		rollback = true
		return nil, err
	}

	return atm, nil
}

func (p *PairOperator) doDetach(ctx *c.Context, attachmentId string, vol *VolumeSpec, provisionerDock *DockSpec) error {

	// Generate the attacher UUID by nodeid and endpoint ip.
	attacherDock, err := p.getAttacherDockByProvisioner(ctx, provisionerDock)
	if err != nil {
		log.Error("Get attacher dock failed, ", err)
		return err
	}
	atm, err := db.C.GetVolumeAttachment(ctx, attachmentId)
	if err != nil {
		log.Error("Get Volume attachment failed, ", err)
		return err
	}
	connData, _ := json.Marshal(atm.ConnectionData)
	detachOpt := &pb.DetachVolumeOpts{
		AccessProtocol: atm.DriverVolumeType,
		ConnectionData: string(connData),
		Metadata:       atm.Metadata,
		Context:        ctx.ToJson(),
	}
	p.volumeController.SetDock(attacherDock)
	if err := p.volumeController.DetachVolume(detachOpt); err != nil {
		log.Error("deatach failed,", err)
		return err
	}

	opt := &pb.DeleteVolumeAttachmentOpts{
		Id:       atm.Id,
		VolumeId: atm.VolumeId,
		HostInfo: &pb.HostInfo{
			Platform:  atm.Platform,
			OsType:    atm.OsType,
			Ip:        atm.Ip,
			Host:      atm.Host,
			Initiator: atm.Initiator,
		},
		AccessProtocol: atm.AccessProtocol,
		Metadata:       utils.MergeStringMaps(atm.Metadata, vol.Metadata),
		DriverName:     provisionerDock.DriverName,
		Context:        ctx.ToJson(),
	}

	p.volumeController.SetDock(provisionerDock)
	if err := p.volumeController.DeleteVolumeAttachment(opt); err != nil {
		log.Error("delete volume attachment failed, ", err)
		return err
	}
	db.C.DeleteVolumeAttachment(ctx, attachmentId)
	return nil
}

func (p *PairOperator) Attach(ctx *c.Context, replica *ReplicationSpec, vol *VolumeSpec) (*ReplicationSpec, error) {

	atm, err := p.doAttach(ctx, vol, p.provisionDock)
	if err != nil {
		return nil, err
	}
	data := map[string]string{
		"Mountpoint":   atm.Mountpoint,
		"AttachmentId": atm.Id,
		"HostName":     atm.Host,
		"HostIp":       atm.Ip,
	}

	if p.isPrimary {
		replica.PrimaryReplicationDriverData = utils.MergeStringMaps(replica.PrimaryReplicationDriverData, data)
	} else {
		// TODO: create replication pair in remote device.
		replica.SecondaryReplicationDriverData = utils.MergeStringMaps(replica.SecondaryReplicationDriverData, data)
	}

	return replica, nil
}

func (p *PairOperator) Detach(ctx *c.Context, replica *ReplicationSpec, vol *VolumeSpec) error {
	var attachmentId string
	if p.isPrimary {
		attachmentId = replica.PrimaryReplicationDriverData["AttachmentId"]
	} else {
		attachmentId = replica.SecondaryReplicationDriverData["AttachmentId"]
	}

	return p.doDetach(ctx, attachmentId, vol, p.provisionDock)
}

func (p *PairOperator) Create(ctx *c.Context, replica *ReplicationSpec, vol *VolumeSpec) (*ReplicationSpec, error) {
	opt := &pb.CreateReplicationOpts{
		Id:                             replica.Id,
		Name:                           replica.Name,
		Description:                    replica.Description,
		PrimaryVolumeId:                replica.PrimaryVolumeId,
		SecondaryVolumeId:              replica.SecondaryVolumeId,
		PrimaryReplicationDriverData:   replica.PrimaryReplicationDriverData,
		SecondaryReplicationDriverData: replica.SecondaryReplicationDriverData,
		ReplicationMode:                replica.ReplicationMode,
		ReplicationPeriod:              replica.ReplicationPeriod,
		ReplicationBandwidth:           replica.ReplicationBandwidth,
		PoolName:                       p.pool.Name,
		DockId:                         p.provisionDock.Id,
		DriverName:                     p.pool.ReplicationDriverName,
		Context:                        ctx.ToJson(),
		IsPrimary:                      p.isPrimary,
		VolumeDataList:                 replica.VolumeDataList,
		Metadata:                       replica.Metadata,
	}
	p.volumeController.SetDock(p.provisionDock)
	return p.volumeController.CreateReplication(opt)
}

func (p *PairOperator) Delete(ctx *c.Context, replica *ReplicationSpec, vol *VolumeSpec) error {

	opt := &pb.DeleteReplicationOpts{
		Id:                             replica.Id,
		Name:                           replica.Name,
		Description:                    replica.Description,
		PrimaryVolumeId:                replica.PrimaryVolumeId,
		SecondaryVolumeId:              replica.SecondaryVolumeId,
		PrimaryReplicationDriverData:   replica.PrimaryReplicationDriverData,
		SecondaryReplicationDriverData: replica.SecondaryReplicationDriverData,
		PoolName:                       p.pool.Name,
		DockId:                         p.provisionDock.Id,
		DriverName:                     p.pool.ReplicationDriverName,
		Context:                        ctx.ToJson(),
		Metadata:                       replica.Metadata,
		IsPrimary:                      p.isPrimary,
	}
	p.volumeController.SetDock(p.provisionDock)
	return p.volumeController.DeleteReplication(opt)
}

func (p *PairOperator) Enable(ctx *c.Context, replica *ReplicationSpec, vol *VolumeSpec) error {
	opt := &pb.EnableReplicationOpts{
		Id:                             replica.Id,
		Name:                           replica.Name,
		Description:                    replica.Description,
		PrimaryVolumeId:                replica.PrimaryVolumeId,
		SecondaryVolumeId:              replica.SecondaryVolumeId,
		PrimaryReplicationDriverData:   replica.PrimaryReplicationDriverData,
		SecondaryReplicationDriverData: replica.SecondaryReplicationDriverData,
		PoolName:                       p.pool.Name,
		DockId:                         p.provisionDock.Id,
		DriverName:                     p.pool.ReplicationDriverName,
		Context:                        ctx.ToJson(),
		Metadata:                       replica.Metadata,
		IsPrimary:                      p.isPrimary,
	}
	p.volumeController.SetDock(p.provisionDock)
	return p.volumeController.EnableReplication(opt)
}

func (p *PairOperator) Disable(ctx *c.Context, replica *ReplicationSpec, vol *VolumeSpec) error {
	opt := &pb.DisableReplicationOpts{
		Id:                             replica.Id,
		Name:                           replica.Name,
		Description:                    replica.Description,
		PrimaryVolumeId:                replica.PrimaryVolumeId,
		SecondaryVolumeId:              replica.SecondaryVolumeId,
		PrimaryReplicationDriverData:   replica.PrimaryReplicationDriverData,
		SecondaryReplicationDriverData: replica.SecondaryReplicationDriverData,
		PoolName:                       p.pool.Name,
		DockId:                         p.provisionDock.Id,
		DriverName:                     p.pool.ReplicationDriverName,
		Context:                        ctx.ToJson(),
		Metadata:                       replica.Metadata,
		IsPrimary:                      p.isPrimary,
	}
	p.volumeController.SetDock(p.provisionDock)
	return p.volumeController.DisableReplication(opt)
}

func (p *PairOperator) Failover(ctx *c.Context, replica *ReplicationSpec, failover *FailoverReplicationSpec, vol *VolumeSpec) error {
	opt := &pb.FailoverReplicationOpts{
		Id:                             replica.Id,
		Name:                           replica.Name,
		Description:                    replica.Description,
		PrimaryVolumeId:                replica.PrimaryVolumeId,
		SecondaryVolumeId:              replica.SecondaryVolumeId,
		PrimaryReplicationDriverData:   replica.PrimaryReplicationDriverData,
		SecondaryReplicationDriverData: replica.SecondaryReplicationDriverData,
		PoolName:                       p.pool.Name,
		DockId:                         p.provisionDock.Id,
		DriverName:                     p.pool.ReplicationDriverName,
		Context:                        ctx.ToJson(),
		Metadata:                       replica.Metadata,
		AllowAttachedVolume:            failover.AllowAttachedVolume,
		SecondaryBackendId:             failover.SecondaryBackendId,
		IsPrimary:                      p.isPrimary,
	}
	p.volumeController.SetDock(p.provisionDock)
	return p.volumeController.FailoverReplication(opt)
}
