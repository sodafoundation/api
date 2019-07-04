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

package scms

import (
	"path/filepath"

	"github.com/golang/glog"
	"github.com/opensds/opensds/pkg/model"
	pb "github.com/opensds/opensds/pkg/model/proto"
)

// Replication driver
type ReplicationDriver struct{}

// Setup
func (r *ReplicationDriver) Setup() error { return nil }

// Unset
func (r *ReplicationDriver) Unset() error { return nil }

// Create and start replication
func (r *ReplicationDriver) CreateReplication(opt *pb.CreateReplicationOpts) (*model.ReplicationSpec, error) {
	glog.Infoln("CMS create migration task ...")

	replica := &model.ReplicationSpec{
		// TODO: return additional important information
		PrimaryReplicationDriverData:   make(map[string]string),
		SecondaryReplicationDriverData: make(map[string]string),
	}

	isPrimary := opt.GetIsPrimary()
	if !isPrimary {
		// on CMSProxy, do nothing
		return replica, nil
	}

	bandwidth := opt.GetReplicationBandwidth()
	primaryData := opt.GetPrimaryReplicationDriverData()
	secondaryData := opt.GetSecondaryReplicationDriverData()
	primaryVolId := opt.GetPrimaryVolumeId()
	secondaryVolId := opt.GetSecondaryVolumeId()
	path, _ := filepath.EvalSymlinks(primaryData["Mountpoint"])
	primaryBackingDevice, _ := filepath.Abs(path)

	path, _ = secondaryData["Mountpoint"]
	secondaryBackingDevice, _ := filepath.Abs(path)
	glog.Infof("%s:%s\n", primaryBackingDevice, secondaryBackingDevice)
	sourceVol := CmsVolume{VolumeId: primaryVolId, VolumeName: primaryBackingDevice}
	targetVol := CmsVolume{VolumeId: secondaryVolId, VolumeName: secondaryBackingDevice}

	task := NewCmsTask(bandwidth, false)
	if err := task.AddVolume(sourceVol, targetVol); err != nil {
		return nil, err
	}

	cmsadm := NewCmsAdm()
	if _, err := cmsadm.CreateTask(task); err != nil {
		return nil, err
	}

	if _, err := cmsadm.Up(); err != nil {
                return nil, err
        }

	return replica, nil
}

// Delete replication
func (r *ReplicationDriver) DeleteReplication(opt *pb.DeleteReplicationOpts) error {
	glog.Infoln("CMS delete migration task ...")

	isPrimary := opt.GetIsPrimary()
	if !isPrimary {
		return nil
	}

	cmsadm := NewCmsAdm()
	_, err := cmsadm.DeleteTask()
	return err
}

// Start replication
func (r *ReplicationDriver) EnableReplication(opt *pb.EnableReplicationOpts) error {
	glog.Infoln("CMS start migration task ....")

	isPrimary := opt.GetIsPrimary()
	if !isPrimary {
		return nil
	}

	cmsadm := NewCmsAdm()
	_, err := cmsadm.Up()
	return err
}

// Stop replication
func (r *ReplicationDriver) DisableReplication(opt *pb.DisableReplicationOpts) error {
	glog.Infoln("CMS stop migration task")

	isPrimary := opt.GetIsPrimary()
	if !isPrimary {
		return nil
	}

	cmsadm := NewCmsAdm()

	_, err := cmsadm.Down()
	return err
}

// Failover Replication
func (r *ReplicationDriver) FailoverReplication(opt *pb.FailoverReplicationOpts) error {
	glog.Infoln("CMS failover ....")
	// Nothing to do here. Failover is executed automatically by CMS plugin.
	return nil
}
