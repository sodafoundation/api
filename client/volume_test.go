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

package client

import (
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/opensds/opensds/pkg/model"
)

var fv = &VolumeMgr{
	Receiver: NewFakeVolumeReceiver(),
}

func NewFakeVolumeReceiver() Receiver {
	return &fakeVolumeReceiver{}
}

type fakeVolumeReceiver struct{}

func (*fakeVolumeReceiver) Recv(
	f reqFunc,
	string,
	method string,
	in interface{},
	out interface{},
) error {
	switch strings.ToUpper(method) {
	case "POST", "PUT":
		switch out.(type) {
		case *model.VolumeSpec:
			if err := json.Unmarshal([]byte(sampleVolume), out); err != nil {
				return err
			}
			break
		case *model.VolumeAttachmentSpec:
			if err := json.Unmarshal([]byte(sampleAttachment), out); err != nil {
				return err
			}
			break
		case *model.VolumeSnapshotSpec:
			if err := json.Unmarshal([]byte(sampleSnapshot), out); err != nil {
				return err
			}
			break
		default:
			return errors.New("output format not supported!")
		}
		break
	case "GET":
		switch out.(type) {
		case *model.VolumeSpec:
			if err := json.Unmarshal([]byte(sampleVolume), out); err != nil {
				return err
			}
			break
		case *[]*model.VolumeSpec:
			if err := json.Unmarshal([]byte(sampleVolumes), out); err != nil {
				return err
			}
			break
		case *model.VolumeAttachmentSpec:
			if err := json.Unmarshal([]byte(sampleAttachment), out); err != nil {
				return err
			}
			break
		case *[]*model.VolumeAttachmentSpec:
			if err := json.Unmarshal([]byte(sampleAttachments), out); err != nil {
				return err
			}
			break
		case *model.VolumeSnapshotSpec:
			if err := json.Unmarshal([]byte(sampleSnapshot), out); err != nil {
				return err
			}
			break
		case *[]*model.VolumeSnapshotSpec:
			if err := json.Unmarshal([]byte(sampleSnapshots), out); err != nil {
				return err
			}
			break
		default:
			return errors.New("output format not supported!")
		}
		break
	case "DELETE":
		break
	default:
		return errors.New("inputed method format not supported!")
	}

	return nil
}

func TestCreateVolume(t *testing.T) {
	expected := &model.VolumeSpec{
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

	vol, err := fv.CreateVolume(&model.VolumeSpec{})
	if err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(vol, expected) {
		t.Errorf("Expected %v, got %v", expected, vol)
		return
	}
}

func TestGetVolume(t *testing.T) {
	var prfID = "bd5b12a8-a101-11e7-941e-d77981b584d8"
	expected := &model.VolumeSpec{
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

	vol, err := fv.GetVolume(prfID)
	if err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(vol, expected) {
		t.Errorf("Expected %v, got %v", expected, vol)
		return
	}
}

func TestListVolumes(t *testing.T) {
	expected := []*model.VolumeSpec{
		{
			BaseModel: &model.BaseModel{
				Id: "bd5b12a8-a101-11e7-941e-d77981b584d8",
			},
			Name:        "sample-volume",
			Description: "This is a sample volume for testing",
			Size:        int64(1),
			Status:      "available",
			PoolId:      "084bf71e-a102-11e7-88a8-e31fe6d52248",
			ProfileId:   "1106b972-66ef-11e7-b172-db03f3689c9c",
		},
	}

	vols, err := fv.ListVolumes()
	if err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(vols, expected) {
		t.Errorf("Expected %v, got %v", expected, vols)
		return
	}
}

func TestDeleteVolume(t *testing.T) {
	var volID = "bd5b12a8-a101-11e7-941e-d77981b584d8"

	if err := fv.DeleteVolume(volID, &model.VolumeSpec{}); err != nil {
		t.Error(err)
		return
	}
}

func TestCreateVolumeAttachment(t *testing.T) {
	var volID = "bd5b12a8-a101-11e7-941e-d77981b584d8"
	expected := &model.VolumeAttachmentSpec{
		BaseModel: &model.BaseModel{
			Id: "f2dda3d2-bf79-11e7-8665-f750b088f63e",
		},
		Status:   "available",
		VolumeId: "bd5b12a8-a101-11e7-941e-d77981b584d8",
		HostInfo: &model.HostInfo{},
		ConnectionInfo: &model.ConnectionInfo{
			DriverVolumeType: "iscsi",
			ConnectionData: map[string]interface{}{
				"targetDiscovered": true,
				"targetIqn":        "iqn.2017-10.io.opensds:volume:00000001",
				"targetPortal":     "127.0.0.0.1:3260",
				"discard":          false,
			},
		},
	}

	atc, err := fv.CreateVolumeAttachment(&model.VolumeAttachmentSpec{
		VolumeId: volID,
		HostInfo: &model.HostInfo{},
	})
	if err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(atc, expected) {
		t.Errorf("Expected %v, got %v", expected, atc)
		return
	}
}

func TestUpdateVolumeAttachment(t *testing.T) {
	var volID = "bd5b12a8-a101-11e7-941e-d77981b584d8"
	expected := &model.VolumeAttachmentSpec{
		BaseModel: &model.BaseModel{
			Id: "f2dda3d2-bf79-11e7-8665-f750b088f63e",
		},
		Status:   "available",
		VolumeId: "bd5b12a8-a101-11e7-941e-d77981b584d8",
		HostInfo: &model.HostInfo{},
		ConnectionInfo: &model.ConnectionInfo{
			DriverVolumeType: "iscsi",
			ConnectionData: map[string]interface{}{
				"targetDiscovered": true,
				"targetIqn":        "iqn.2017-10.io.opensds:volume:00000001",
				"targetPortal":     "127.0.0.0.1:3260",
				"discard":          false,
			},
		},
	}

	atc, err := fv.UpdateVolumeAttachment("f2dda3d2-bf79-11e7-8665-f750b088f63e", &model.VolumeAttachmentSpec{
		VolumeId: volID,
		HostInfo: &model.HostInfo{},
	})
	if err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(atc, expected) {
		t.Errorf("Expected %v, got %v", expected, atc)
		return
	}
}

func TestGetVolumeAttachment(t *testing.T) {
	var atcID = "f2dda3d2-bf79-11e7-8665-f750b088f63e"
	expected := &model.VolumeAttachmentSpec{
		BaseModel: &model.BaseModel{
			Id: "f2dda3d2-bf79-11e7-8665-f750b088f63e",
		},
		Status:   "available",
		VolumeId: "bd5b12a8-a101-11e7-941e-d77981b584d8",
		HostInfo: &model.HostInfo{},
		ConnectionInfo: &model.ConnectionInfo{
			DriverVolumeType: "iscsi",
			ConnectionData: map[string]interface{}{
				"targetDiscovered": true,
				"targetIqn":        "iqn.2017-10.io.opensds:volume:00000001",
				"targetPortal":     "127.0.0.0.1:3260",
				"discard":          false,
			},
		},
	}

	atc, err := fv.GetVolumeAttachment(atcID)
	if err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(atc, expected) {
		t.Errorf("Expected %v, got %v", expected, atc)
		return
	}
}

func TestListVolumeAttachments(t *testing.T) {
	expected := []*model.VolumeAttachmentSpec{
		{
			BaseModel: &model.BaseModel{
				Id: "f2dda3d2-bf79-11e7-8665-f750b088f63e",
			},
			Status:   "available",
			VolumeId: "bd5b12a8-a101-11e7-941e-d77981b584d8",
			HostInfo: &model.HostInfo{},
			ConnectionInfo: &model.ConnectionInfo{
				DriverVolumeType: "iscsi",
				ConnectionData: map[string]interface{}{
					"targetDiscovered": true,
					"targetIqn":        "iqn.2017-10.io.opensds:volume:00000001",
					"targetPortal":     "127.0.0.0.1:3260",
					"discard":          false,
				},
			},
		},
	}

	atcs, err := fv.ListVolumeAttachments()
	if err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(atcs, expected) {
		t.Errorf("Expected %v, got %v", expected, atcs)
		return
	}
}

func TestDeleteVolumeAttachment(t *testing.T) {
	var atcID = "f2dda3d2-bf79-11e7-8665-f750b088f63e"

	if err := fv.DeleteVolumeAttachment(atcID, &model.VolumeAttachmentSpec{
		VolumeId: "bd5b12a8-a101-11e7-941e-d77981b584d8",
	}); err != nil {
		t.Error(err)
		return
	}
}

func TestCreateVolumeSnapshot(t *testing.T) {
	expected := &model.VolumeSnapshotSpec{
		BaseModel: &model.BaseModel{
			Id: "3769855c-a102-11e7-b772-17b880d2f537",
		},
		Name:        "sample-snapshot-01",
		Description: "This is the first sample snapshot for testing",
		Size:        int64(1),
		Status:      "created",
		VolumeId:    "bd5b12a8-a101-11e7-941e-d77981b584d8",
	}

	snp, err := fv.CreateVolumeSnapshot(&model.VolumeSnapshotSpec{
		VolumeId: "bd5b12a8-a101-11e7-941e-d77981b584d8",
	})
	if err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(snp, expected) {
		t.Errorf("Expected %v, got %v", expected, snp)
		return
	}
}

func TestGetVolumeSnapshot(t *testing.T) {
	var snpID = "3769855c-a102-11e7-b772-17b880d2f537"
	expected := &model.VolumeSnapshotSpec{
		BaseModel: &model.BaseModel{
			Id: "3769855c-a102-11e7-b772-17b880d2f537",
		},
		Name:        "sample-snapshot-01",
		Description: "This is the first sample snapshot for testing",
		Size:        int64(1),
		Status:      "created",
		VolumeId:    "bd5b12a8-a101-11e7-941e-d77981b584d8",
	}

	snp, err := fv.GetVolumeSnapshot(snpID)
	if err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(snp, expected) {
		t.Errorf("Expected %v, got %v", expected, snp)
		return
	}
}

func TestListVolumeSnapshots(t *testing.T) {
	expected := []*model.VolumeSnapshotSpec{
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

	snps, err := fv.ListVolumeSnapshots()
	if err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(snps, expected) {
		t.Errorf("Expected %v, got %v", expected, snps)
		return
	}
}

func TestDeleteVolumeSnapshot(t *testing.T) {
	var snpID = "3769855c-a102-11e7-b772-17b880d2f537"

	if err := fv.DeleteVolumeSnapshot(snpID, &model.VolumeSnapshotSpec{
		VolumeId: "bd5b12a8-a101-11e7-941e-d77981b584d8",
	}); err != nil {
		t.Error(err)
		return
	}
}

var (
	sampleVolume = `{
		"id": "bd5b12a8-a101-11e7-941e-d77981b584d8",
		"name": "sample-volume",
		"description": "This is a sample volume for testing",
		"size": 1,
		"status": "available",
		"poolId": "084bf71e-a102-11e7-88a8-e31fe6d52248",
		"profileId": "1106b972-66ef-11e7-b172-db03f3689c9c"
	}`

	sampleVolumes = `[
		{
			"id": "bd5b12a8-a101-11e7-941e-d77981b584d8",
			"name": "sample-volume",
			"description": "This is a sample volume for testing",
			"size": 1,
			"status": "available",
			"poolId": "084bf71e-a102-11e7-88a8-e31fe6d52248",
			"profileId": "1106b972-66ef-11e7-b172-db03f3689c9c"
		}
	]`

	sampleAttachment = `{
		"id": "f2dda3d2-bf79-11e7-8665-f750b088f63e",
		"name": "sample-volume-attachment",
		"description": "This is a sample volume attachment for testing",
		"status": "available",
		"volumeId": "bd5b12a8-a101-11e7-941e-d77981b584d8",
		"hostInfo": {},
		"connectionInfo": {
			"driverVolumeType": "iscsi",
			"data": {
				"targetDiscovered": true,
				"targetIqn": "iqn.2017-10.io.opensds:volume:00000001",
				"targetPortal": "127.0.0.0.1:3260",
				"discard": false
			}
		}
	}`

	sampleAttachments = `[
		{
			"id": "f2dda3d2-bf79-11e7-8665-f750b088f63e",
			"name": "sample-volume-attachment",
			"description": "This is a sample volume attachment for testing",
			"status": "available",
			"volumeId": "bd5b12a8-a101-11e7-941e-d77981b584d8",
			"hostInfo": {},
			"connectionInfo": {
				"driverVolumeType": "iscsi",
				"data": {
					"targetDiscovered": true,
					"targetIqn": "iqn.2017-10.io.opensds:volume:00000001",
					"targetPortal": "127.0.0.0.1:3260",
					"discard": false
				}
			}
		}
	]`

	sampleSnapshot = `{
		"id": "3769855c-a102-11e7-b772-17b880d2f537",
		"name": "sample-snapshot-01",
		"description": "This is the first sample snapshot for testing",
		"size": 1,
		"status": "created",
		"volumeId": "bd5b12a8-a101-11e7-941e-d77981b584d8"		
	}`

	sampleSnapshots = `[
		{
			"id": "3769855c-a102-11e7-b772-17b880d2f537",
			"name": "sample-snapshot-01",
			"description": "This is the first sample snapshot for testing",
			"size": 1,
			"status": "created",
			"volumeId": "bd5b12a8-a101-11e7-941e-d77981b584d8"	
		},
		{
			"id": "3bfaf2cc-a102-11e7-8ecb-63aea739d755",
			"name": "sample-snapshot-02",
			"description": "This is the second sample snapshot for testing",
			"size": 1,
			"status": "created",
			"volumeId": "bd5b12a8-a101-11e7-941e-d77981b584d8"	
		}
	]`
)
