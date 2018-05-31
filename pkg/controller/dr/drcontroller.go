// Copyright (c) 2018 Huawei Technologies Co., Ltd. All Rights Reserved.
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

package dr

import (
	"encoding/json"
	"fmt"
	"strings"

	log "github.com/golang/glog"
	c "github.com/opensds/opensds/pkg/context"
	"github.com/opensds/opensds/pkg/controller/volume"
	"github.com/opensds/opensds/pkg/db"
	pb "github.com/opensds/opensds/pkg/dock/proto"
	. "github.com/opensds/opensds/pkg/model"
	"github.com/opensds/opensds/pkg/utils"
	"github.com/satori/go.uuid"
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

func (d *DrController) LoadOperator(replicationType string) {
	factory := NewReplicationFactory(replicationType)
	d.primaryOp = factory.GetPrimaryOperator(d.volumeController)
	d.secondaryOp = factory.GetSecondaryOperator(d.volumeController)
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

	// Load replication operator by replication type.
	d.LoadOperator(pPool.ReplicationType)
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
	pPool, _ := db.C.GetPool(ctx, primaryVol.PoolId)
	d.LoadOperator(pPool.ReplicationType)
	err := d.primaryOp.Delete(ctx, replica, primaryVol)
	if err != nil {
		return err
	}
	err = d.secondaryOp.Delete(ctx, replica, secondaryVol)
	if err != nil {
		return err
	}
	// clean up replication driver data in volume database
	primaryVol.ReplicationDriverData = make(map[string]string)
	db.C.UpdateVolume(ctx, primaryVol)
	secondaryVol.ReplicationDriverData = make(map[string]string)
	db.C.UpdateVolume(ctx, secondaryVol)
	return db.C.DeleteReplication(ctx, replica.Id)
}

func (d *DrController) EnableReplication(ctx *c.Context, replica *ReplicationSpec, primaryVol, secondaryVol *VolumeSpec) error {
	pPool, _ := db.C.GetPool(ctx, primaryVol.PoolId)
	d.LoadOperator(pPool.ReplicationType)
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
	pPool, _ := db.C.GetPool(ctx, primaryVol.PoolId)
	d.LoadOperator(pPool.ReplicationType)
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
	pPool, _ := db.C.GetPool(ctx, primaryVol.PoolId)
	d.LoadOperator(pPool.ReplicationType)
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
}

type ReplicationFactory interface {
	GetPrimaryOperator(controller volume.Controller) ReplicationOperator
	GetSecondaryOperator(controller volume.Controller) ReplicationOperator
}

func NewReplicationFactory(replicaType string) ReplicationFactory {
	if replicaType == ReplicationTypeArray {
		return &ArrayBasedFactory{}
	}
	return &HostBasedFactory{}
}

type HostBasedFactory struct {
	volumeController *volume.Controller
}

func (h *HostBasedFactory) GetPrimaryOperator(controller volume.Controller) ReplicationOperator {
	return NewHostPairOperator(controller, true)
}

func (h *HostBasedFactory) GetSecondaryOperator(controller volume.Controller) ReplicationOperator {
	return NewHostPairOperator(controller, false)
}

type ArrayBasedFactory struct {
	volumeController *volume.Controller
}

func (a *ArrayBasedFactory) GetPrimaryOperator(controller volume.Controller) ReplicationOperator {
	return NewArrayPairOperator(controller, true)
}

func (a *ArrayBasedFactory) GetSecondaryOperator(controller volume.Controller) ReplicationOperator {
	return NewArrayPairOperator(controller, false)
}

func NewHostPairOperator(controller volume.Controller, isPrimary bool) *HostPairOperator {
	return &HostPairOperator{
		BaseOperator{
			volumeController: controller,
			isPrimary:        isPrimary,
		},
	}
}

type HostPairOperator struct {
	BaseOperator
}

func (h *HostPairOperator) getAttacherDockByProvisioner(ctx *c.Context, provisionerDock *DockSpec) (*DockSpec, error) {
	segments := strings.Split(provisionerDock.Endpoint, ":")
	endpointIp := segments[len(segments)-2]
	// Generate the attacher UUID by nodeid and endpoint ip.
	attacherDockId := uuid.NewV5(uuid.NamespaceOID, provisionerDock.NodeId+":"+endpointIp)
	return db.C.GetDock(ctx, attacherDockId.String())
}

func (h *HostPairOperator) attach(ctx *c.Context, vol *VolumeSpec, provisionerDock *DockSpec) (*VolumeAttachmentSpec, error) {
	attacherDock, err := h.getAttacherDockByProvisioner(ctx, provisionerDock)
	h.volumeController.SetDock(provisionerDock)
	attachmentId := uuid.NewV4().String()
	var createAttachOpt = &pb.CreateAttachmentOpts{
		Id:       attachmentId,
		VolumeId: vol.Id,
		HostInfo: &pb.HostInfo{
			Platform:  attacherDock.Metadata["Platform"],
			OsType:    attacherDock.Metadata["OsType"],
			Ip:        attacherDock.Metadata["HostIp"],
			Host:      attacherDock.NodeId,
			Initiator: attacherDock.Metadata["Initiator"],
		},
		Metadata:   vol.Metadata,
		DriverName: provisionerDock.DriverName,
		Context:    ctx.ToJson(),
	}

	atm, err := h.volumeController.CreateVolumeAttachment(createAttachOpt)
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
			opt := &pb.DeleteAttachmentOpts{
				Id:       atm.Id,
				VolumeId: atm.VolumeId,
				HostInfo: &pb.HostInfo{
					Platform:  atm.Platform,
					OsType:    atm.OsType,
					Ip:        atm.Ip,
					Host:      atm.Host,
					Initiator: atm.Initiator,
				},
				Metadata:   utils.MergeStringMaps(atm.Metadata, vol.Metadata),
				DriverName: provisionerDock.DriverName,
				Context:    ctx.ToJson(),
			}
			h.volumeController.DeleteVolumeAttachment(opt)
			db.C.DeleteVolumeAttachment(ctx, atm.Id)
		}
	}()

	h.volumeController.SetDock(attacherDock)
	connData, _ := json.Marshal(atm.ConnectionData)
	var attachOpt = &pb.AttachVolumeOpts{
		AccessProtocol: atm.DriverVolumeType,
		ConnectionData: string(connData),
		Metadata:       map[string]string{},
		Context:        ctx.ToJson(),
	}
	mountPoint, err := h.volumeController.AttachVolume(attachOpt)
	if err != nil {
		rollback = true
		log.Errorf("attach volume failed, %v", err)
		return nil, err
	}

	atm.Mountpoint = mountPoint
	_, err = db.C.UpdateVolumeAttachment(ctx, atm.Id, atm)
	if err != nil {
		rollback = true
		return nil, err
	}

	return atm, nil
}

func (h *HostPairOperator) detach(ctx *c.Context, attachmentId string, vol *VolumeSpec, provisionerDock *DockSpec) error {

	// Generate the attacher UUID by nodeid and endpoint ip.
	attacherDock, err := h.getAttacherDockByProvisioner(ctx, provisionerDock)
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
	h.volumeController.SetDock(attacherDock)
	if err := h.volumeController.DetachVolume(detachOpt); err != nil {
		log.Error("deatach failed,", err)
		return err
	}

	opt := &pb.DeleteAttachmentOpts{
		Id:       atm.Id,
		VolumeId: atm.VolumeId,
		HostInfo: &pb.HostInfo{
			Platform:  atm.Platform,
			OsType:    atm.OsType,
			Ip:        atm.Ip,
			Host:      atm.Host,
			Initiator: atm.Initiator,
		},
		Metadata:   utils.MergeStringMaps(atm.Metadata, vol.Metadata),
		DriverName: provisionerDock.DriverName,
		Context:    ctx.ToJson(),
	}

	h.volumeController.SetDock(provisionerDock)
	if err := h.volumeController.DeleteVolumeAttachment(opt); err != nil {
		log.Error("delete volume attachment failed, ", err)
		return err
	}
	db.C.DeleteVolumeAttachment(ctx, attachmentId)
	return nil
}

func (h *HostPairOperator) Create(ctx *c.Context, replica *ReplicationSpec, vol *VolumeSpec) (*ReplicationSpec, error) {
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

	atm, err := h.attach(ctx, vol, provisionerDock)
	if err != nil {
		return nil, err
	}
	data := map[string]string{
		"Mountpoint":   atm.Mountpoint,
		"AttachmentId": atm.Id,
		"HostName":     atm.Host,
		"HostIp":       atm.Ip,
	}

	volList, _ := db.C.ListVolumes(ctx)
	var volumeDataList []*pb.VolumeData
	for _, v := range volList {
		if v.ReplicationDriverData != nil && len(v.ReplicationDriverData) != 0 {
			v.ReplicationDriverData["VolumeId"] = v.Id
			volumeDataList = append(volumeDataList, &pb.VolumeData{Data: v.ReplicationDriverData})
		}
	}

	replica.VolumeDataList = volumeDataList
	if h.isPrimary {
		replica.PrimaryReplicationDriverData = utils.MergeStringMaps(replica.PrimaryReplicationDriverData, data)
	} else {
		// TODO: create replication pair in remote device.
		replica.SecondaryReplicationDriverData = utils.MergeStringMaps(replica.SecondaryReplicationDriverData, data)
	}

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
		PoolName:                       pool.Name,
		DockId:                         provisionerDock.Id,
		DriverName:                     pool.ReplicationDriverName,
		Context:                        ctx.ToJson(),
		IsPrimary:                      h.isPrimary,
		VolumeDataList:                 volumeDataList,
		Metadata:                       replica.Metadata,
	}
	h.volumeController.SetDock(provisionerDock)
	return h.volumeController.CreateReplication(opt)
}

func (h *HostPairOperator) Delete(ctx *c.Context, replica *ReplicationSpec, vol *VolumeSpec) error {
	pool, err := db.C.GetPool(ctx, vol.PoolId)
	if err != nil {
		log.Error("get pool failed", err)
		return err
	}

	provisionerDock, err := db.C.GetDockByPoolId(ctx, vol.PoolId)
	if err != nil {
		log.Error("When search dock in db by pool id", err)
		return err
	}

	opt := &pb.DeleteReplicationOpts{
		Id:                             replica.Id,
		Name:                           replica.Name,
		Description:                    replica.Description,
		PrimaryVolumeId:                replica.PrimaryVolumeId,
		SecondaryVolumeId:              replica.SecondaryVolumeId,
		PrimaryReplicationDriverData:   replica.PrimaryReplicationDriverData,
		SecondaryReplicationDriverData: replica.SecondaryReplicationDriverData,
		PoolName:                       pool.Name,
		DockId:                         provisionerDock.Id,
		DriverName:                     pool.ReplicationDriverName,
		Context:                        ctx.ToJson(),
		Metadata:                       replica.Metadata,
		IsPrimary:                      h.isPrimary,
	}

	// invoke both side
	h.volumeController.SetDock(provisionerDock)
	if err = h.volumeController.DeleteReplication(opt); err != nil {
		return err
	}

	var attachmentId string
	if h.isPrimary {
		attachmentId = replica.PrimaryReplicationDriverData["AttachmentId"]
	} else {
		attachmentId = replica.SecondaryReplicationDriverData["AttachmentId"]
	}

	return h.detach(ctx, attachmentId, vol, provisionerDock)
}

func NewArrayPairOperator(controller volume.Controller, isPrimary bool) *ArrayPairOperator {
	return &ArrayPairOperator{
		BaseOperator{
			volumeController: controller,
			isPrimary:        isPrimary,
		},
	}
}

type ArrayPairOperator struct {
	BaseOperator
}

func (a *ArrayPairOperator) Create(ctx *c.Context, replica *ReplicationSpec, vol *VolumeSpec) (*ReplicationSpec, error) {
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

	a.volumeController.SetDock(provisionerDock)
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
		PoolName:                       pool.Name,
		DockId:                         provisionerDock.Id,
		DriverName:                     provisionerDock.DriverName,
		Context:                        ctx.ToJson(),
		IsPrimary:                      a.isPrimary,
	}
	return a.volumeController.CreateReplication(opt)
}

func (a *ArrayPairOperator) Delete(ctx *c.Context, replica *ReplicationSpec, vol *VolumeSpec) error {
	pool, err := db.C.GetPool(ctx, vol.PoolId)
	if err != nil {
		log.Error("get pool failed", err)
		return err
	}

	provisionerDock, err := db.C.GetDockByPoolId(ctx, vol.PoolId)
	if err != nil {
		log.Error("When search dock in db by pool id", err)
		return err
	}

	a.volumeController.SetDock(provisionerDock)
	opt := &pb.DeleteReplicationOpts{
		Id:                             replica.Id,
		Name:                           replica.Name,
		Description:                    replica.Description,
		PrimaryVolumeId:                replica.PrimaryVolumeId,
		SecondaryVolumeId:              replica.SecondaryVolumeId,
		PrimaryReplicationDriverData:   replica.PrimaryReplicationDriverData,
		SecondaryReplicationDriverData: replica.SecondaryReplicationDriverData,
		PoolName:                       pool.Name,
		DockId:                         provisionerDock.Id,
		DriverName:                     provisionerDock.DriverName,
		Context:                        ctx.ToJson(),
		Metadata:                       replica.Metadata,
		IsPrimary:                      a.isPrimary,
	}
	return a.volumeController.DeleteReplication(opt)
}

type BaseOperator struct {
	volumeController volume.Controller
	isPrimary        bool
}

func (h *BaseOperator) Enable(ctx *c.Context, replica *ReplicationSpec, vol *VolumeSpec) error {
	pool, err := db.C.GetPool(ctx, vol.PoolId)
	if err != nil {
		log.Error("get pool failed", err)
		return err
	}

	provisionerDock, err := db.C.GetDockByPoolId(ctx, vol.PoolId)
	if err != nil {
		log.Error("When search dock in db by pool id", err)
		return err
	}

	opt := &pb.EnableReplicationOpts{
		Id:                             replica.Id,
		Name:                           replica.Name,
		Description:                    replica.Description,
		PrimaryVolumeId:                replica.PrimaryVolumeId,
		SecondaryVolumeId:              replica.SecondaryVolumeId,
		PrimaryReplicationDriverData:   replica.PrimaryReplicationDriverData,
		SecondaryReplicationDriverData: replica.SecondaryReplicationDriverData,
		PoolName:                       pool.Name,
		DockId:                         provisionerDock.Id,
		DriverName:                     pool.ReplicationDriverName,
		Context:                        ctx.ToJson(),
		Metadata:                       replica.Metadata,
		IsPrimary:                      h.isPrimary,
	}
	h.volumeController.SetDock(provisionerDock)
	return h.volumeController.EnableReplication(opt)
}

func (h *BaseOperator) Disable(ctx *c.Context, replica *ReplicationSpec, vol *VolumeSpec) error {
	pool, err := db.C.GetPool(ctx, vol.PoolId)
	if err != nil {
		log.Error("get pool failed", err)
		return err
	}

	provisionerDock, err := db.C.GetDockByPoolId(ctx, vol.PoolId)
	if err != nil {
		log.Error("When search dock in db by pool id", err)
		return err
	}

	opt := &pb.DisableReplicationOpts{
		Id:                             replica.Id,
		Name:                           replica.Name,
		Description:                    replica.Description,
		PrimaryVolumeId:                replica.PrimaryVolumeId,
		SecondaryVolumeId:              replica.SecondaryVolumeId,
		PrimaryReplicationDriverData:   replica.PrimaryReplicationDriverData,
		SecondaryReplicationDriverData: replica.SecondaryReplicationDriverData,
		PoolName:                       pool.Name,
		DockId:                         provisionerDock.Id,
		DriverName:                     pool.ReplicationDriverName,
		Context:                        ctx.ToJson(),
		Metadata:                       replica.Metadata,
		IsPrimary:                      h.isPrimary,
	}
	h.volumeController.SetDock(provisionerDock)
	return h.volumeController.DisableReplication(opt)
}

func (h *BaseOperator) Failover(ctx *c.Context, replica *ReplicationSpec, failover *FailoverReplicationSpec, vol *VolumeSpec) error {
	pool, err := db.C.GetPool(ctx, vol.PoolId)
	if err != nil {
		log.Error("get pool failed", err)
		return err
	}

	provisionerDock, err := db.C.GetDockByPoolId(ctx, vol.PoolId)
	if err != nil {
		log.Error("When search dock in db by pool id", err)
		return err
	}

	opt := &pb.FailoverReplicationOpts{
		Id:                             replica.Id,
		Name:                           replica.Name,
		Description:                    replica.Description,
		PrimaryVolumeId:                replica.PrimaryVolumeId,
		SecondaryVolumeId:              replica.SecondaryVolumeId,
		PrimaryReplicationDriverData:   replica.PrimaryReplicationDriverData,
		SecondaryReplicationDriverData: replica.SecondaryReplicationDriverData,
		PoolName:                       pool.Name,
		DockId:                         provisionerDock.Id,
		DriverName:                     pool.ReplicationDriverName,
		Context:                        ctx.ToJson(),
		Metadata:                       replica.Metadata,
		AllowAttachedVolume:            failover.AllowAttachedVolume,
		SecondaryBackendId:             failover.SecondaryBackendId,
		IsPrimary:                      h.isPrimary,
	}
	h.volumeController.SetDock(provisionerDock)
	return h.volumeController.FailoverReplication(opt)
}
