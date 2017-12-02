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
This module implements the etcd database operation of data structure
defined in api module.

*/

package etcd

import (
	"reflect"
	"strings"
	"testing"

	"github.com/opensds/opensds/pkg/model"
)

type fakeClientCaller struct{}

func (*fakeClientCaller) Create(req *Request) *Response {
	return &Response{
		Status: "Success",
	}
}

func (*fakeClientCaller) Get(req *Request) *Response {
	var resp []string

	if strings.Contains(req.Url, "docks") {
		resp = append(resp, sampleDocks[0])
	}
	if strings.Contains(req.Url, "pools") {
		resp = append(resp, samplePools[0])
	}
	if strings.Contains(req.Url, "profiles") {
		resp = append(resp, sampleProfiles[0])
	}
	if strings.Contains(req.Url, "volumes") {
		resp = append(resp, sampleVolumes[0])
	}
	if strings.Contains(req.Url, "attachments") {
		resp = append(resp, sampleAttachments[0])
	}
	if strings.Contains(req.Url, "snapshots") {
		resp = append(resp, sampleSnapshots[0])
	}

	return &Response{
		Status:  "Success",
		Message: resp,
	}
}

func (*fakeClientCaller) List(req *Request) *Response {
	var resp []string

	if strings.Contains(req.Url, "docks") {
		resp = sampleDocks
	}
	if strings.Contains(req.Url, "pools") {
		resp = samplePools
	}
	if strings.Contains(req.Url, "profiles") {
		resp = sampleProfiles
	}
	if strings.Contains(req.Url, "volumes") {
		resp = sampleVolumes
	}
	if strings.Contains(req.Url, "attachments") {
		resp = sampleAttachments
	}
	if strings.Contains(req.Url, "snapshots") {
		resp = sampleSnapshots
	}

	return &Response{
		Status:  "Success",
		Message: resp,
	}
}

func (*fakeClientCaller) Update(req *Request) *Response {
	return &Response{
		Status: "Success",
	}
}

func (*fakeClientCaller) Delete(req *Request) *Response {
	return &Response{
		Status: "Success",
	}
}

var fc = &Client{
	clientInterface: &fakeClientCaller{},
}

func TestCreateDock(t *testing.T) {
	if err := fc.CreateDock(&model.DockSpec{BaseModel: &model.BaseModel{}}); err != nil {
		t.Error("Create dock failed:", err)
	}
}

func TestCreatePool(t *testing.T) {
	if err := fc.CreatePool(&model.StoragePoolSpec{BaseModel: &model.BaseModel{}}); err != nil {
		t.Error("Create pool failed:", err)
	}
}

func TestCreateProfile(t *testing.T) {
	if err := fc.CreateProfile(&model.ProfileSpec{BaseModel: &model.BaseModel{}}); err != nil {
		t.Error("Create profile failed:", err)
	}
}

func TestCreateVolume(t *testing.T) {
	if err := fc.CreateVolume(&model.VolumeSpec{BaseModel: &model.BaseModel{}}); err != nil {
		t.Error("Create volume failed:", err)
	}
}

func TestCreateVolumeAttachment(t *testing.T) {
	if _, err := fc.CreateVolumeAttachment(&model.VolumeAttachmentSpec{BaseModel: &model.BaseModel{}}); err != nil {
		t.Error("Create volume attachment failed:", err)
	}
}

func TestCreateVolumeSnapshot(t *testing.T) {
	if err := fc.CreateVolumeSnapshot(&model.VolumeSnapshotSpec{BaseModel: &model.BaseModel{}}); err != nil {
		t.Error("Create volume snapshot failed:", err)
	}
}

func TestGetDock(t *testing.T) {
	dck, err := fc.GetDock("")
	if err != nil {
		t.Error("Get dock failed:", err)
	}

	var expected = &model.DockSpec{
		BaseModel: &model.BaseModel{
			Id: "b7602e18-771e-11e7-8f38-dbd6d291f4e0",
		},
		Name:        "sample",
		Description: "sample backend service",
		Endpoint:    "localhost:50050",
		DriverName:  "sample",
	}
	if !reflect.DeepEqual(dck, expected) {
		t.Errorf("Expected %+v, got %+v\n", expected, dck)
	}
}

func TestGetPool(t *testing.T) {
	pol, err := fc.GetPool("")
	if err != nil {
		t.Error("Get pool failed:", err)
	}

	var expected = &model.StoragePoolSpec{
		BaseModel: &model.BaseModel{
			Id: "084bf71e-a102-11e7-88a8-e31fe6d52248",
		},
		Name:             "sample-pool-01",
		Description:      "This is the first sample storage pool for testing",
		TotalCapacity:    int64(100),
		FreeCapacity:     int64(90),
		DockId:           "b7602e18-771e-11e7-8f38-dbd6d291f4e0",
		AvailabilityZone: "default",
		Parameters: map[string]interface{}{
			"diskType": "SSD",
			"thin":     true,
		},
	}
	if !reflect.DeepEqual(pol, expected) {
		t.Errorf("Expected %+v, got %+v\n", expected, pol)
	}
}

func TestGetProfile(t *testing.T) {
	prf, err := fc.GetProfile("")
	if err != nil {
		t.Error("Get profile failed:", err)
	}

	var expected = &model.ProfileSpec{
		BaseModel: &model.BaseModel{
			Id: "1106b972-66ef-11e7-b172-db03f3689c9c",
		},
		Name:        "default",
		Description: "default policy",
		Extra:       model.ExtraSpec{},
	}
	if !reflect.DeepEqual(prf, expected) {
		t.Errorf("Expected %+v, got %+v\n", expected, prf)
	}
}

func TesGetVolume(t *testing.T) {
	vol, err := fc.GetVolume("")
	if err != nil {
		t.Error("Get volume failed:", err)
	}

	var expected = &model.VolumeSpec{
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
	if !reflect.DeepEqual(vol, expected) {
		t.Errorf("Expected %+v, got %+v\n", expected, vol)
	}
}

func TestGetVolumeAttachment(t *testing.T) {
	atc, err := fc.GetVolumeAttachment("")
	if err != nil {
		t.Error("Get volume attachment failed:", err)
	}

	var expected = &model.VolumeAttachmentSpec{
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
	if !reflect.DeepEqual(atc, expected) {
		t.Errorf("Expected %+v, got %+v\n", expected, atc)
	}
}

func TestGetVolumeSnapshot(t *testing.T) {
	snp, err := fc.GetVolumeSnapshot("")
	if err != nil {
		t.Error("Get volume snapshot failed:", err)
	}

	var expected = &model.VolumeSnapshotSpec{
		BaseModel: &model.BaseModel{
			Id: "3769855c-a102-11e7-b772-17b880d2f537",
		},
		Name:        "sample-snapshot-01",
		Description: "This is the first sample snapshot for testing",
		Size:        int64(1),
		Status:      "created",
		VolumeId:    "bd5b12a8-a101-11e7-941e-d77981b584d8",
	}
	if !reflect.DeepEqual(snp, expected) {
		t.Errorf("Expected %+v, got %+v\n", expected, snp)
	}
}

func TestListDocks(t *testing.T) {
	dcks, err := fc.ListDocks()
	if err != nil {
		t.Error("List docks failed:", err)
	}

	var expected = []*model.DockSpec{
		{
			BaseModel: &model.BaseModel{
				Id: "b7602e18-771e-11e7-8f38-dbd6d291f4e0",
			},
			Name:        "sample",
			Description: "sample backend service",
			Endpoint:    "localhost:50050",
			DriverName:  "sample",
		},
	}
	if !reflect.DeepEqual(dcks, expected) {
		t.Errorf("Expected %+v, got %+v\n", expected, dcks)
	}
}

func TestListPools(t *testing.T) {
	pols, err := fc.ListPools()
	if err != nil {
		t.Error("List pools failed:", err)
	}

	var expected = []*model.StoragePoolSpec{
		{
			BaseModel: &model.BaseModel{
				Id: "084bf71e-a102-11e7-88a8-e31fe6d52248",
			},
			Name:             "sample-pool-01",
			Description:      "This is the first sample storage pool for testing",
			TotalCapacity:    int64(100),
			FreeCapacity:     int64(90),
			DockId:           "b7602e18-771e-11e7-8f38-dbd6d291f4e0",
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
			DockId:           "b7602e18-771e-11e7-8f38-dbd6d291f4e0",
			Parameters: map[string]interface{}{
				"diskType": "SAS",
				"thin":     true,
			},
		},
	}
	if !reflect.DeepEqual(pols, expected) {
		t.Errorf("Expected %+v, got %+v\n", expected, pols)
	}
}

func TestListProfiles(t *testing.T) {
	prfs, err := fc.ListProfiles()
	if err != nil {
		t.Error("List profiles failed:", err)
	}

	var expected = []*model.ProfileSpec{
		{
			BaseModel: &model.BaseModel{
				Id: "1106b972-66ef-11e7-b172-db03f3689c9c",
			},
			Name:        "default",
			Description: "default policy",
			Extra:       model.ExtraSpec{},
		},
		{
			BaseModel: &model.BaseModel{
				Id: "2f9c0a04-66ef-11e7-ade2-43158893e017",
			},
			Name:        "silver",
			Description: "silver policy",
			Extra: model.ExtraSpec{
				"diskType": "SAS",
				"thin":     true,
			},
		},
	}
	if !reflect.DeepEqual(prfs, expected) {
		t.Errorf("Expected %+v, got %+v\n", expected, prfs)
	}
}

func TestListVolumes(t *testing.T) {
	vols, err := fc.ListVolumes()
	if err != nil {
		t.Error("List volumes failed:", err)
	}

	var expected = []*model.VolumeSpec{
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
	if !reflect.DeepEqual(vols, expected) {
		t.Errorf("Expected %+v, got %+v\n", expected, vols)
	}
}

func TestListVolumeAttachments(t *testing.T) {
	atcs, err := fc.ListVolumeAttachments("")
	if err != nil {
		t.Error("List volume attachments failed:", err)
	}

	var expected = []*model.VolumeAttachmentSpec{
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
	if !reflect.DeepEqual(atcs, expected) {
		t.Errorf("Expected %+v, got %+v\n", expected, atcs)
	}
}

func TestListVolumeSnapshots(t *testing.T) {
	snps, err := fc.ListVolumeSnapshots()
	if err != nil {
		t.Error("List volume snapshots failed:", err)
	}

	var expected = []*model.VolumeSnapshotSpec{
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
	if !reflect.DeepEqual(snps, expected) {
		t.Errorf("Expected %+v, got %+v\n", expected, snps)
	}
}

func TestDeleteDock(t *testing.T) {
	if err := fc.DeleteDock(""); err != nil {
		t.Error("Delete dock failed:", err)
	}
}

func TestDeletePool(t *testing.T) {
	if err := fc.DeletePool(""); err != nil {
		t.Error("Delete pool failed:", err)
	}
}

func TestDeleteProfile(t *testing.T) {
	if err := fc.DeleteProfile(""); err != nil {
		t.Error("Delete profile failed:", err)
	}
}

func TestDeleteVolume(t *testing.T) {
	if err := fc.DeleteVolume(""); err != nil {
		t.Error("Delete volume failed:", err)
	}
}

func TestDeleteVolumeAttachment(t *testing.T) {
	if err := fc.DeleteVolumeAttachment(""); err != nil {
		t.Error("Delete volume attachment failed:", err)
	}
}

func TestDeleteVolumeSnapshot(t *testing.T) {
	if err := fc.DeleteVolumeSnapshot(""); err != nil {
		t.Error("Delete volume snapshot failed:", err)
	}
}

var (
	sampleProfiles = []string{
		`{
			"id": "1106b972-66ef-11e7-b172-db03f3689c9c",
			"name":        "default",
			"description": "default policy",
			"extras": {}
		}`,
		`{
			"id": "2f9c0a04-66ef-11e7-ade2-43158893e017",
			"name":        "silver",
			"description": "silver policy",
			"extras": {
				"diskType": "SAS",
				"thin":     true
			}
		}`,
	}

	sampleDocks = []string{
		`{
			"id": "b7602e18-771e-11e7-8f38-dbd6d291f4e0",
			"name":        "sample",
			"description": "sample backend service",
			"endpoint":    "localhost:50050",
			"driverName":  "sample"
		}`,
	}

	samplePools = []string{
		`{
			"id": "084bf71e-a102-11e7-88a8-e31fe6d52248",
			"name":             "sample-pool-01",
			"description":      "This is the first sample storage pool for testing",
			"totalCapacity":    100,
			"freeCapacity":     90,
			"dockId":           "b7602e18-771e-11e7-8f38-dbd6d291f4e0",
			"availabilityZone": "default",
			"extras": {
				"diskType": "SSD",
				"thin":     true
			}
		}`,
		`{
			"id": "a594b8ac-a103-11e7-985f-d723bcf01b5f",
			"name":             "sample-pool-02",
			"description":      "This is the second sample storage pool for testing",
			"totalCapacity":    200,
			"freeCapacity":     170,
			"availabilityZone": "default",
			"dockId":           "b7602e18-771e-11e7-8f38-dbd6d291f4e0",
			"extras": {
				"diskType": "SAS",
				"thin":     true
			}
		}`,
	}

	sampleVolumes = []string{
		`{
			"id": "bd5b12a8-a101-11e7-941e-d77981b584d8",
			"name":        "sample-volume",
			"description": "This is a sample volume for testing",
			"size":        1,
			"status":      "available",
			"poolId":      "084bf71e-a102-11e7-88a8-e31fe6d52248",
			"profileId":   "1106b972-66ef-11e7-b172-db03f3689c9c"
		}`,
	}

	sampleAttachments = []string{
		`{
			"id": "f2dda3d2-bf79-11e7-8665-f750b088f63e",
			"status":   "available",
			"volumeId": "bd5b12a8-a101-11e7-941e-d77981b584d8",
			"hostInfo": {},
			"connectionInfo": {
				"driverVolumeType": "iscsi",
				"data": {
					"targetDiscovered": true,
					"targetIqn":        "iqn.2017-10.io.opensds:volume:00000001",
					"targetPortal":     "127.0.0.0.1:3260",
					"discard":          false
				}
			}
		}`,
	}

	sampleSnapshots = []string{
		`{
			"id": "3769855c-a102-11e7-b772-17b880d2f537",
			"name":        "sample-snapshot-01",
			"description": "This is the first sample snapshot for testing",
			"size":        1,
			"status":      "created",
			"volumeId":    "bd5b12a8-a101-11e7-941e-d77981b584d8"
		}`,
		`{
			"id": "3bfaf2cc-a102-11e7-8ecb-63aea739d755",
			"name":        "sample-snapshot-02",
			"description": "This is the second sample snapshot for testing",
			"size":        1,
			"status":      "created",
			"volumeId":    "bd5b12a8-a101-11e7-941e-d77981b584d8"
		}`,
	}
)
