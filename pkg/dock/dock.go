// Copyright (c) 2017 Huawei Technologies Co., Ltd. All Rights Reserved.
//
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

/*
This module implements a standard SouthBound interface of resources to
storage plugins.

*/

package dock

import (
	log "github.com/golang/glog"
	"github.com/opensds/opensds/contrib/drivers"
	c "github.com/opensds/opensds/pkg/context"
	"github.com/opensds/opensds/pkg/db"
	"github.com/opensds/opensds/pkg/dock/discovery"
	pb "github.com/opensds/opensds/pkg/dock/proto"
	"github.com/opensds/opensds/pkg/model"
)

// Brain is a global variable that controls the dock module.
var Brain *DockHub

// DockHub is a reference structure with fields that represent some required
// parameters for initializing and controlling the volume driver.
type DockHub struct {
	// Discoverer represents the mechanism of DockHub discovering the storage
	// capabilities from different backends.
	Discoverer *discovery.DockDiscoverer
	// Driver represents the specified backend resource. This field is used
	// for initializing the specified volume driver.
	Driver drivers.VolumeDriver
}

// NewDockHub method creates a new DockHub and returns its pointer.
func NewDockHub() *DockHub {
	return &DockHub{
		Discoverer: discovery.NewDiscoverer(),
	}
}

// TriggerDiscovery
func (d *DockHub) TriggerDiscovery() error {
	var err error

	if err = d.Discoverer.Init(); err != nil {
		return err
	}

	ctx := &discovery.Context{
		StopChan: make(chan bool),
		ErrChan:  make(chan error),
		MetaChan: make(chan string),
	}
	go discovery.DiscoveryAndReport(d.Discoverer, ctx)
	go func(ctx *discovery.Context) {
		if err = <-ctx.ErrChan; err != nil {
			log.Error("When calling capabilty report method:", err)
			ctx.StopChan <- true
		}
	}(ctx)

	return nil
}

// CreateVolume
func (d *DockHub) CreateVolume(opt *pb.CreateVolumeOpts) (*model.VolumeSpec, error) {
	//Get the storage drivers and do some initializations.
	d.Driver = drivers.Init(opt.GetDriverName())
	defer drivers.Clean(d.Driver)

	log.Info("Calling volume driver to create volume...")

	//Call function of StorageDrivers configured by storage drivers.
	vol, err := d.Driver.CreateVolume(opt)
	if err != nil {
		//Change the status of the volume to error when the creation faild
		err = d.UpdateVolumeStatus(c.NewContextFormJson(opt.GetContext()), opt.GetId(), model.VOLUME_ERROR)
		if err != nil {
			log.Error("When update volume in db:", err)
			return nil, err
		}
		log.Error("When calling volume driver to create volume:", err)
		return nil, err
	}

	vol.PoolId, vol.ProfileId = opt.GetPoolId(), opt.GetProfileId()

	// Update the volume data in database.
	vol.Status = model.VOLUME_AVAILABLE

	result, err := db.C.UpdateVolume(c.NewContextFormJson(opt.GetContext()), vol)
	if err != nil {
		log.Error("When update volume in db:", err)
		return nil, err
	}

	return result, nil
}

func (d *DockHub) UpdateVolumeStatus(ctx *c.Context, VolumeId string, status string) error {
	var vol = &model.VolumeSpec{
		Status: status,
	}
	_, err := db.C.UpdateVolume(ctx, vol)
	if err != nil {
		log.Error("When update volume in db:", err)
		return err
	}
	return nil
}

// DeleteVolume
func (d *DockHub) DeleteVolume(opt *pb.DeleteVolumeOpts) error {
	var err error

	//Get the storage drivers and do some initializations.
	d.Driver = drivers.Init(opt.GetDriverName())
	defer drivers.Clean(d.Driver)

	log.Info("Calling volume driver to delete volume...")

	//Call function of StorageDrivers configured by storage drivers.
	if err = d.Driver.DeleteVolume(opt); err != nil {
		log.Error("When calling volume driver to delete volume:", err)
		errdel := d.UpdateVolumeStatus(c.NewContextFormJson(opt.GetContext()), opt.GetId(), model.VOLUEM_ERROR_DELETING)
		if errdel != nil {
			return errdel
		}
		return err
	}

	if err = db.C.DeleteVolume(c.NewContextFormJson(opt.GetContext()), opt.GetId()); err != nil {
		log.Error("Error occurred in dock module when delete volume in db:", err)
		return err
	}

	return nil
}

// ExtendVolume ...
func (d *DockHub) ExtendVolume(opt *pb.ExtendVolumeOpts) (*model.VolumeSpec, error) {
	//Get the storage drivers and do some initializations.
	d.Driver = drivers.Init(opt.GetDriverName())
	defer drivers.Clean(d.Driver)

	log.Info("Calling volume driver to extend volume...")

	//Call function of StorageDrivers configured by storage drivers.
	vol, err := d.Driver.ExtendVolume(opt)
	if err != nil {
		vol.Status = model.VOLUME_ERROR
		_, err := db.C.UpdateVolume(c.NewContextFormJson(opt.GetContext()), vol)
		if err != nil {
			log.Error("When update volume in db:", err)
			return nil, err
		}
		log.Error("When calling volume driver to extend volume:", err)
		return nil, err
	}

	vol.PoolId, vol.ProfileId = opt.GetPoolId(), opt.GetProfileId()
	// Update the volume data in database.
	vol.Status = model.VOLUME_AVAILABLE
	result, err := db.C.UpdateVolume(c.NewContextFormJson(opt.GetContext()), vol)
	if err != nil {
		log.Error("When update volume in db:", err)
		return nil, err
	}

	return result, nil
}

// CreateVolumeAttachment
func (d *DockHub) CreateVolumeAttachment(opt *pb.CreateAttachmentOpts) (*model.VolumeAttachmentSpec, error) {
	//Get the storage drivers and do some initializations.
	d.Driver = drivers.Init(opt.GetDriverName())
	defer drivers.Clean(d.Driver)

	log.Info("Calling volume driver to initialize volume connection...")

	//Call function of StorageDrivers configured by storage drivers.
	connInfo, err := d.Driver.InitializeConnection(opt)
	if err != nil {
		log.Error("Call driver to initialize volume connection failed:", err)
		return nil, err
	}

	var atc = &model.VolumeAttachmentSpec{
		BaseModel: &model.BaseModel{
			Id: opt.GetId(),
		},
		Status:   model.VOLUMEATM_AVAILABLE,
		VolumeId: opt.GetVolumeId(),
		HostInfo: model.HostInfo{
			Platform:  opt.HostInfo.GetPlatform(),
			OsType:    opt.HostInfo.GetOsType(),
			Ip:        opt.HostInfo.GetIp(),
			Host:      opt.HostInfo.GetHost(),
			Initiator: opt.HostInfo.GetInitiator(),
		},
		ConnectionInfo: *connInfo,
		Metadata:       opt.GetMetadata(),
	}
	result, err := db.C.UpdateVolumeAttachment(c.NewContextFormJson(opt.GetContext()), atc.Id, atc)
	if err != nil {
		log.Error("Error occurred in dock module when update volume attachment in db:", err)
		return nil, err
	}

	return result, nil
}

// DeleteVolumeAttachment
func (d *DockHub) DeleteVolumeAttachment(opt *pb.DeleteAttachmentOpts) error {
	//Get the storage drivers and do some initializations.
	d.Driver = drivers.Init(opt.GetDriverName())
	defer drivers.Clean(d.Driver)

	log.Info("Calling volume driver to terminate volume connection...")

	//Call function of StorageDrivers configured by storage drivers.
	if err := d.Driver.TerminateConnection(opt); err != nil {
		var atc = &model.VolumeAttachmentSpec{
			BaseModel: &model.BaseModel{
				Id: opt.GetId(),
			},
			Status: model.VOLUMEATM_ERROR_DELETING,
		}
		_, err := db.C.UpdateVolumeAttachment(c.NewContextFormJson(opt.GetContext()), atc.Id, atc)
		if err != nil {
			log.Error("Error occurred in dock module when update volume attachment in db:", err)
			return err
		}
		log.Error("Call driver to terminate volume connection failed:", err)
		return err
	}

	if err := db.C.DeleteVolumeAttachment(c.NewContextFormJson(opt.GetContext()), opt.GetId()); err != nil {
		log.Error("Error occurred in dock module when delete volume attachment in db:", err)
		return err
	}

	return nil
}

// CreateSnapshot
func (d *DockHub) CreateSnapshot(opt *pb.CreateVolumeSnapshotOpts) (*model.VolumeSnapshotSpec, error) {
	//Get the storage drivers and do some initializations.
	d.Driver = drivers.Init(opt.GetDriverName())
	defer drivers.Clean(d.Driver)

	log.Info("Calling volume driver to create snapshot...")

	//Call function of StorageDrivers configured by storage drivers.
	snp, err := d.Driver.CreateSnapshot(opt)
	if err != nil {
		err = d.UpdateSnapshotStatus(c.NewContextFormJson(opt.GetContext()), opt.GetId(), model.VOLUMESNAP_ERROR)
		if err != nil {
			log.Error("When update volume snapshot in db:", err)
			return nil, err
		}
		log.Error("Call driver to create volume snashot failed:", err)
		return nil, err
	}
	snp.Status = model.VOLUMESNAP_AVAILABLE
	result, err := db.C.UpdateVolumeSnapshot(c.NewContextFormJson(opt.GetContext()), opt.GetId(), snp)
	if err != nil {
		log.Error("When update volume snapshot in db:", err)
		return nil, err
	}

	return result, nil
}

func (d *DockHub) UpdateSnapshotStatus(ctx *c.Context, snapId string, status string) error {
	var snap = &model.VolumeSnapshotSpec{
		Status: status,
	}
	_, err := db.C.UpdateVolumeSnapshot(ctx, snapId, snap)
	if err != nil {
		return err
	}
	return nil
}

// DeleteSnapshot
func (d *DockHub) DeleteSnapshot(opt *pb.DeleteVolumeSnapshotOpts) error {
	var err error

	//Get the storage drivers and do some initializations.
	d.Driver = drivers.Init(opt.GetDriverName())
	defer drivers.Clean(d.Driver)

	log.Info("Calling volume driver to delete snapshot...")

	//Call function of StorageDrivers configured by storage drivers.
	if err = d.Driver.DeleteSnapshot(opt); err != nil {
		errDel := d.UpdateSnapshotStatus(c.NewContextFormJson(opt.GetContext()), opt.GetId(), model.VOLUMESNAP_ERROR_DELETING)
		if errDel != nil {
			log.Error("Error occurs when deleting volume snapshot from driver:", err)
			return errDel
		}
		log.Error("When calling volume driver to delete volume:", err)
		return err
	}

	if err = db.C.DeleteVolumeSnapshot(c.NewContextFormJson(opt.GetContext()), opt.GetId()); err != nil {
		log.Error("Error occurred in dock module when delete volume snapshot in db:", err)
		return err
	}

	return nil
}
