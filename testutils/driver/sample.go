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

func (*Driver) Setup() error { return nil }

func (*Driver) Unset() error { return nil }

func (*Driver) CreateVolume(opt *pb.CreateVolumeOpts) (*model.VolumeSpec, error) {
	return &sampleVolume, nil
}

func (*Driver) PullVolume(volIdentifier string) (*model.VolumeSpec, error) {
	if volIdentifier == sampleVolume.GetId() {
		return &sampleVolume, nil
	}

	return nil, errors.New("Can't find volume " + volIdentifier)
}

func (*Driver) DeleteVolume(opt *pb.DeleteVolumeOpts) error {
	return nil
}

func (*Driver) InitializeConnection(opt *pb.CreateAttachmentOpts) (*model.ConnectionInfo, error) {
	return &sampleConnection, nil
}

func (*Driver) TerminateConnection(opt *pb.DeleteAttachmentOpts) error { return nil }

func (*Driver) CreateSnapshot(opt *pb.CreateVolumeSnapshotOpts) (*model.VolumeSnapshotSpec, error) {
	return &sampleSnapshots[0], nil
}

func (*Driver) PullSnapshot(snapIdentifier string) (*model.VolumeSnapshotSpec, error) {
	for _, snapshot := range sampleSnapshots {
		if snapIdentifier == snapshot.GetId() {
			return &snapshot, nil
		}
	}

	return nil, errors.New("Can't find snapshot " + snapIdentifier)
}

func (*Driver) DeleteSnapshot(opt *pb.DeleteVolumeSnapshotOpts) error {
	return nil
}

func (*Driver) ListPools() ([]*model.StoragePoolSpec, error) {
	var pols []*model.StoragePoolSpec

	for i := range samplePools {
		pols = append(pols, &samplePools[i])
	}
	return pols, nil
}

var (
	samplePools = []model.StoragePoolSpec{
		{
			BaseModel: &model.BaseModel{
				Id: "084bf71e-a102-11e7-88a8-e31fe6d52248",
			},
			Name:             "sample-pool-01",
			Description:      "This is the first sample storage pool for testing",
			TotalCapacity:    int64(100),
			FreeCapacity:     int64(90),
			AvailabilityZone: "default",
			Parameters: map[string]interface{}{
				"diskType": "SSD",
				"thin":     true,
			},
		},
		{
			BaseModel: &model.BaseModel{
				Id: "a594b8ac-a103-11e7-985f-d723bcf01b5f",
			},
			Name:             "sample-pool-02",
			Description:      "This is the second sample storage pool for testing",
			TotalCapacity:    int64(200),
			FreeCapacity:     int64(170),
			AvailabilityZone: "default",
			Parameters: map[string]interface{}{
				"diskType": "SAS",
				"thin":     true,
			},
		},
	}

	sampleVolume = model.VolumeSpec{
		BaseModel: &model.BaseModel{
			Id: "bd5b12a8-a101-11e7-941e-d77981b584d8",
		},
		Name:        "sample-volume",
		Description: "This is a sample volume for testing",
		Size:        int64(1),
		Status:      "available",
		PoolId:      "084bf71e-a102-11e7-88a8-e31fe6d52248",
		ProfileId:   "1106b972-66ef-11e7-b172-db03f3689c9c",
	}

	sampleConnection = model.ConnectionInfo{
		DriverVolumeType: "iscsi",
		ConnectionData: map[string]interface{}{
			"targetDiscovered": true,
			"targetIqn":        "iqn.2017-10.io.opensds:volume:00000001",
			"targetPortal":     "127.0.0.0.1:3260",
			"discard":          false,
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
