// Copyright 2017 The OpenSDS Authors.
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
	"github.com/opensds/opensds/pkg/db"
	"github.com/opensds/opensds/pkg/dock/discovery"
	pb "github.com/opensds/opensds/pkg/dock/proto"
	api "github.com/opensds/opensds/pkg/model"
	"github.com/opensds/opensds/pkg/utils"
)

var Brain *DockHub

// A reference to DockHub structure with fields that represent some required
// parameters for initializing and controlling the volume driver.
type DockHub struct {
	// Discoverer represents the mechanism of DockHub discovering the storage
	// capabilities from different backends.
	Discoverer *discovery.DockDiscoverer
	// Driver represents the specified backend resource. This field is used
	// for initializing the specified volume driver.
	Driver drivers.VolumeDriver
}

func NewDockHub() *DockHub {
	return &DockHub{
		Discoverer: discovery.NewDiscoverer(),
	}
}

func (d *DockHub) TriggerDiscovery() error {
	var err error

	if err = d.Discoverer.Init(); err != nil {
		return err
	}
	if err = d.Discoverer.Discover(d.Driver); err != nil {
		return err
	}
	if err = d.Discoverer.Store(); err != nil {
		return err
	}

	return err
}

func (d *DockHub) CreateVolume(opt *pb.CreateVolumeOpts) (*api.VolumeSpec, error) {
	//Get the storage drivers and do some initializations.
	d.Driver = drivers.Init(opt.GetDriverName())
	defer drivers.Clean(d.Driver)

	log.Info("Calling volume driver to create volume...")

	//Call function of StorageDrivers configured by storage drivers.
	vol, err := d.Driver.CreateVolume(opt)
	if err != nil {
		log.Error("When calling volume driver to create volume:", err)
		return nil, err
	}
	vol.PoolId, vol.ProfileId = opt.GetPoolId(), opt.GetProfileId()

	// Validate the data.
	if err = utils.ValidateData(vol, utils.S); err != nil {
		log.Error("When validate volume data:", err)
		return nil, err
	}

	// Store the volume data into database.
	if err = db.C.CreateVolume(vol); err != nil {
		log.Error("When create volume in db module:", err)
		return nil, err
	}

	return vol, nil
}

func (d *DockHub) DeleteVolume(opt *pb.DeleteVolumeOpts) error {
	var err error

	//Get the storage drivers and do some initializations.
	d.Driver = drivers.Init(opt.GetDriverName())
	defer drivers.Clean(d.Driver)

	log.Info("Calling volume driver to delete volume...")

	//Call function of StorageDrivers configured by storage drivers.
	if err = d.Driver.DeleteVolume(opt); err != nil {
		log.Error("When calling volume driver to delete volume:", err)
		return err
	}

	if err = db.C.DeleteVolume(opt.GetId()); err != nil {
		log.Error("Error occured in dock module when delete volume in db:", err)
		return err
	}

	return nil
}

func (d *DockHub) CreateVolumeAttachment(opt *pb.CreateAttachmentOpts) (*api.VolumeAttachmentSpec, error) {
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

	var atc = &api.VolumeAttachmentSpec{
		BaseModel: &api.BaseModel{},
		VolumeId:  opt.GetVolumeId(),
		HostInfo: &api.HostInfo{
			Platform:  opt.HostInfo.GetPlatform(),
			OsType:    opt.HostInfo.GetOsType(),
			Ip:        opt.HostInfo.GetIp(),
			Host:      opt.HostInfo.GetHost(),
			Initiator: opt.HostInfo.GetInitiator(),
		},
		ConnectionInfo: connInfo,
		Metadata:       opt.GetMetadata(),
	}

	// Validate the data.
	if err = utils.ValidateData(atc, utils.S); err != nil {
		log.Error("When validate volume attachment data:", err)
		return nil, err
	}

	result, err := db.C.CreateVolumeAttachment(atc)
	if err != nil {
		log.Error("Error occured in dock module when create volume attachment in db:", err)
		return nil, err
	}

	return result, nil
}

func (d *DockHub) DeleteVolumeAttachment(opt *pb.DeleteAttachmentOpts) error {
	//Get the storage drivers and do some initializations.
	d.Driver = drivers.Init(opt.GetDriverName())
	defer drivers.Clean(d.Driver)

	log.Info("Calling volume driver to terminate volume connection...")

	//Call function of StorageDrivers configured by storage drivers.
	if err := d.Driver.TerminateConnection(opt); err != nil {
		log.Error("Call driver to terminate volume connection failed:", err)
		return err
	}

	if err := db.C.DeleteVolumeAttachment(opt.GetId()); err != nil {
		log.Error("Error occured in dock module when delete volume attachment in db:", err)
		return err
	}

	return nil
}

func (d *DockHub) CreateSnapshot(opt *pb.CreateVolumeSnapshotOpts) (*api.VolumeSnapshotSpec, error) {
	//Get the storage drivers and do some initializations.
	d.Driver = drivers.Init(opt.GetDriverName())
	defer drivers.Clean(d.Driver)

	log.Info("Calling volume driver to create snapshot...")

	//Call function of StorageDrivers configured by storage drivers.
	snp, err := d.Driver.CreateSnapshot(opt)
	if err != nil {
		log.Error("Call driver to create volume snashot failed:", err)
		return nil, err
	}

	// Validate the data.
	if err = utils.ValidateData(snp, utils.S); err != nil {
		log.Error("When validate volume snapshot data:", err)
	}

	if err := db.C.CreateVolumeSnapshot(snp); err != nil {
		log.Error("Error occured in dock module when create volume snapshot in db:", err)
		return nil, err
	}

	return snp, nil
}

func (d *DockHub) DeleteSnapshot(opt *pb.DeleteVolumeSnapshotOpts) error {
	var err error

	//Get the storage drivers and do some initializations.
	d.Driver = drivers.Init(opt.GetDriverName())
	defer drivers.Clean(d.Driver)

	log.Info("Calling volume driver to delete snapshot...")

	//Call function of StorageDrivers configured by storage drivers.
	if err = d.Driver.DeleteSnapshot(opt); err != nil {
		log.Error("When calling volume driver to delete volume:", err)
		return err
	}

	if err = db.C.DeleteVolumeSnapshot(opt.GetId()); err != nil {
		log.Error("Error occured in dock module when delete volume snapshot in db:", err)
		return err
	}

	return nil
}
