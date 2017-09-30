// Copyright (c) 2016 Huawei Technologies Co., Ltd. All Rights Reserved.
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

/*
This module implements a sample driver for OpenSDS. This driver will handle all
operations of volume and return a fake value.

*/

package sample

import (
	"errors"

	pb "github.com/opensds/opensds/pkg/dock/proto"
	"github.com/opensds/opensds/pkg/model"
)

type Driver struct{}

func (d *Driver) Setup() error { return nil }

func (d *Driver) Unset() error { return nil }

func (d *Driver) CreateVolume(opt *pb.CreateVolumeOpts) (*model.VolumeSpec, error) {
	return &sampleVolume, nil
}

func (d *Driver) PullVolume(volIdentifier string) (*model.VolumeSpec, error) {
	if volIdentifier == sampleVolume.GetId() {
		return &sampleVolume, nil
	}

	return nil, errors.New("Can't find volume " + volIdentifier)
}

func (d *Driver) DeleteVolume(opt *pb.DeleteVolumeOpts) error {
	return nil
}

func (d *Driver) InitializeConnection(opt *pb.CreateAttachmentOpts) (*model.ConnectionInfo, error) {
	return &sampleConnection, nil
}

func (d *Driver) CreateSnapshot(opt *pb.CreateVolumeSnapshotOpts) (*model.VolumeSnapshotSpec, error) {
	return &sampleSnapshots[0], nil
}

func (d *Driver) PullSnapshot(snapIdentifier string) (*model.VolumeSnapshotSpec, error) {
	for _, snapshot := range sampleSnapshots {
		if snapIdentifier == snapshot.GetId() {
			return &snapshot, nil
		}
	}

	return nil, errors.New("Can't find snapshot " + snapIdentifier)
}

func (d *Driver) DeleteSnapshot(opt *pb.DeleteVolumeSnapshotOpts) error {
	return nil
}

func (d *Driver) ListPools() (*[]model.StoragePoolSpec, error) {
	return &samplePools, nil
}

var (
	samplePools = []model.StoragePoolSpec{
		{
			BaseModel: &model.BaseModel{
				Id: "084bf71e-a102-11e7-88a8-e31fe6d52248",
			},
			Name:             "sample-pool-01",
			Description:      "This is the first sample storage pool for testing",
			AvailabilityZone: "nova",
			TotalCapacity:    int64(10),
			FreeCapacity:     int64(9),
			Parameters: map[string]interface{}{
				"iops":      1000,
				"disk-type": "ssd",
			},
		},
		{
			BaseModel: &model.BaseModel{
				Id: "a594b8ac-a103-11e7-985f-d723bcf01b5f",
			},
			Name:             "sample-pool-02",
			Description:      "This is the second sample storage pool for testing",
			AvailabilityZone: "nova",
			TotalCapacity:    int64(20),
			FreeCapacity:     int64(17),
			Parameters: map[string]interface{}{
				"disk-type":    "hdd",
				"replica-sets": 3,
			},
		},
	}

	sampleVolume = model.VolumeSpec{
		BaseModel: &model.BaseModel{
			Id: "bd5b12a8-a101-11e7-941e-d77981b584d8",
		},
		Name:             "sample-volume",
		Description:      "This is a sample volume for testing",
		Size:             int64(1),
		AvailabilityZone: "nova",
		Status:           "available",
		PoolId:           "084bf71e-a102-11e7-88a8-e31fe6d52248",
		ProfileId:        "gold",
	}

	sampleConnection = model.ConnectionInfo{
		DriverVolumeType: "iscsi",
		ConnectionData: map[string]interface{}{
			"target_discovered": true,
			"target_iqn":        "iqn.2010-10.org.openstack:volume-00000001",
			"target_portal":     "127.0.0.0.1:3260",
			"volume_id":         "9a0d35d0-175a-11e4-8c21-0800200c9a66",
			"discard":           false,
		},
	}

	sampleSnapshots = []model.VolumeSnapshotSpec{
		{
			BaseModel: &model.BaseModel{
				Id: "3769855c-a102-11e7-b772-17b880d2f537",
			},
			Name:        "sample-snapshot-01",
			Description: "This is the first sample snapshot for testing",
			Size:        int64(1),
			Status:      "created",
			VolumeId:    "bd5b12a8-a101-11e7-941e-d77981b584d8",
		},
		{
			BaseModel: &model.BaseModel{
				Id: "3bfaf2cc-a102-11e7-8ecb-63aea739d755",
			},
			Name:        "sample-snapshot-02",
			Description: "This is the second sample snapshot for testing",
			Size:        int64(1),
			Status:      "created",
			VolumeId:    "bd5b12a8-a101-11e7-941e-d77981b584d8",
		},
	}
)
