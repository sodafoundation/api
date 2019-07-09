// Copyright 2019 The OpenSDS Authors.
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
	"reflect"
	"testing"

	"github.com/opensds/opensds/pkg/model"
	. "github.com/opensds/opensds/testutils/collection"
)

var fakeShareMgr = &FileShareMgr{
	Receiver: NewFakeFileShareReceiver(),
}

func TestCreateFileShare(t *testing.T) {
	expected := &model.FileShareSpec{
		BaseModel: &model.BaseModel{
			Id: "bd5b12a8-a101-11e7-941e-d77981b584d8",
		},
		Name:        "sample-fileshare",
		Description: "This is a sample fileshare for testing",
		Size:        int64(1),
		Status:      "available",
		PoolId:      "084bf71e-a102-11e7-88a8-e31fe6d52248",
		ProfileId:   "1106b972-66ef-11e7-b172-db03f3689c9c",
	}

	share, err := fakeShareMgr.CreateFileShare(&model.FileShareSpec{})
	if err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(share, expected) {
		t.Errorf("expected %+v, got %+v", expected, share)
		return
	}
}

func TestGetFileShare(t *testing.T) {
	var shareID = "d2975ebe-d82c-430f-b28e-f373746a71ca"
	expected := &model.FileShareSpec{
		BaseModel: &model.BaseModel{
			Id: "bd5b12a8-a101-11e7-941e-d77981b584d8",
		},
		Name:        "sample-fileshare",
		Description: "This is a sample fileshare for testing",
		Size:        int64(1),
		Status:      "available",
		PoolId:      "084bf71e-a102-11e7-88a8-e31fe6d52248",
		ProfileId:   "1106b972-66ef-11e7-b172-db03f3689c9c",
	}

	share, err := fakeShareMgr.GetFileShare(shareID)
	if err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(share, expected) {
		t.Errorf("expected %v, got %v", expected, share)
		return
	}
}

func TestListFileShares(t *testing.T) {
	sampleFileShares := []model.FileShareSpec{
		{
			BaseModel: &model.BaseModel{
				Id: "d2975ebe-d82c-430f-b28e-f373746a71ca",
			},
			Name:             "sample-fileshare-01",
			Description:      "This is first sample fileshare for testing",
			Size:             int64(1),
			Status:           "available",
			PoolId:           "a5965ebe-dg2c-434t-b28e-f373746a71ca",
			ProfileId:        "b3585ebe-c42c-120g-b28e-f373746a71ca",
			SnapshotId:       "b7602e18-771e-11e7-8f38-dbd6d291f4eg",
			AvailabilityZone: "default",
			ExportLocations:  []string{"192.168.100.100"},
		},
		{
			BaseModel: &model.BaseModel{
				Id: "1e643aca-4922-4b1a-bb98-4245054aeff4",
			},
			Name:             "sample-fileshare-2",
			Description:      "This is second sample fileshare for testing",
			Size:             int64(1),
			Status:           "available",
			PoolId:           "d5f65ebe-ag2c-341s-a25e-f373746a71dr",
			ProfileId:        "1e643aca-4922-4b1a-bb98-4245054aeff4",
			SnapshotId:       "a5965ebe-dg2c-434t-b28e-f373746a71ca",
			AvailabilityZone: "default",
			ExportLocations:  []string{"192.168.100.101"},
		},
	}

	var expected []*model.FileShareSpec
	expected = append(expected, &sampleFileShares[0])
	expected = append(expected, &sampleFileShares[1])
	shares, err := fakeShareMgr.ListFileShares(map[string]string{"limit": "3", "offset": "4"})

	if err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(shares, expected) {
		t.Errorf("expected %v, got %v", expected, shares)
		return
	}
}

func TestDeleteFileShare(t *testing.T) {
	var shareID = "d2975ebe-d82c-430f-b28e-f373746a71ca"

	if err := fakeShareMgr.DeleteFileShare(shareID); err != nil {
		t.Error(err)
		return
	}
}

func TestUpdateFileShare(t *testing.T) {
	var shareID = "d2975ebe-d82c-430f-b28e-f373746a71ca"
	share := &model.FileShareSpec{
		Name:        "sample-share",
		Description: "This is a sample share for testing",
	}

	result, err := fakeShareMgr.UpdateFileShare(shareID, share)
	if err != nil {
		t.Error(err)
		return
	}

	expected := &model.FileShareSpec{
		BaseModel: &model.BaseModel{
			Id: "bd5b12a8-a101-11e7-941e-d77981b584d8",
		},
		Name:        "sample-fileshare",
		Description: "This is a sample fileshare for testing",
		Size:        int64(1),
		Status:      "available",
		PoolId:      "084bf71e-a102-11e7-88a8-e31fe6d52248",
		ProfileId:   "1106b972-66ef-11e7-b172-db03f3689c9c",
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
		return
	}
}

func TestCreateFileShareSnapshot(t *testing.T) {
	expected := &model.FileShareSnapshotSpec{
		BaseModel: &model.BaseModel{
			Id: "3769855c-a102-11e7-b772-17b880d2f537",
		},
		Name:        "sample-snapshot-01",
		Description: "This is the first sample snapshot for testing",
		Status:      "available",
		ProfileId:   "1106b972-66ef-11e7-b172-db03f3689c9c",
		FileShareId: "bd5b12a8-a101-11e7-941e-d77981b584d8",
		ShareSize:   1,
	}

	shareSnapshot, err := fakeShareMgr.CreateFileShareSnapshot(&model.FileShareSnapshotSpec{})
	if err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(shareSnapshot, expected) {
		t.Errorf("expected %+v, got %+v", expected, shareSnapshot)
		return
	}
}

func TestGetFileShareSnapshot(t *testing.T) {
	var shareID = "3769855c-a102-11e7-b772-17b880d2f537"
	expected := &model.FileShareSnapshotSpec{
		BaseModel: &model.BaseModel{
			Id: "3769855c-a102-11e7-b772-17b880d2f537",
		},
		Name:        "sample-snapshot-01",
		Description: "This is the first sample snapshot for testing",
		Status:      "available",
		ProfileId:   "1106b972-66ef-11e7-b172-db03f3689c9c",
		FileShareId: "bd5b12a8-a101-11e7-941e-d77981b584d8",
		ShareSize:   1,
	}

	shareSnapshot, err := fakeShareMgr.GetFileShareSnapshot(shareID)
	if err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(shareSnapshot, expected) {
		t.Errorf("expected %v, got %v", expected, shareSnapshot)
		return
	}
}

func TestListFileShareSnapshots(t *testing.T) {
	SampleFileShareSnapshots = []model.FileShareSnapshotSpec{
		{
			BaseModel: &model.BaseModel{
				Id: "3769855c-a102-11e7-b772-17b880d2f537",
			},
			Name:         "sample-snapshot-01",
			Description:  "This is the first sample snapshot for testing",
			SnapshotSize: int64(1),
			Status:       "available",
		},
		{
			BaseModel: &model.BaseModel{
				Id: "3bfaf2cc-a102-11e7-8ecb-63aea739d755",
			},
			Name:         "sample-snapshot-02",
			Description:  "This is the second sample snapshot for testing",
			SnapshotSize: int64(1),
			Status:       "available",
		},
	}
	var expected []*model.FileShareSnapshotSpec
	expected = append(expected, &SampleFileShareSnapshots[0])
	expected = append(expected, &SampleFileShareSnapshots[1])
	shareSnapshots, err := fakeShareMgr.ListFileShareSnapshots(map[string]string{"limit": "3", "offset": "4"})

	if err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(shareSnapshots, expected) {
		t.Errorf("expected %v, got %v", expected, shareSnapshots)
		return
	}
}

func TestDeleteFileShareSnapshot(t *testing.T) {
	var shareSnapshotID = "3769855c-a102-11e7-b772-17b880d2f537"

	if err := fakeShareMgr.DeleteFileShareSnapshot(shareSnapshotID); err != nil {
		t.Error(err)
		return
	}
}

func TestUpdateFileShareSnapshot(t *testing.T) {
	var shareSnapshotID = "3769855c-a102-11e7-b772-17b880d2f537"
	shareSnapshot := &model.FileShareSnapshotSpec{
		Name:        "sample-share",
		Description: "This is a sample share for testing",
	}

	result, err := fakeShareMgr.UpdateFileShareSnapshot(shareSnapshotID, shareSnapshot)
	if err != nil {
		t.Error(err)
		return
	}

	expected := &model.FileShareSnapshotSpec{
		BaseModel: &model.BaseModel{
			Id: "3769855c-a102-11e7-b772-17b880d2f537",
		},
		Name:        "sample-snapshot-01",
		Description: "This is the first sample snapshot for testing",
		Status:      "available",
		ProfileId:   "1106b972-66ef-11e7-b172-db03f3689c9c",
		FileShareId: "bd5b12a8-a101-11e7-941e-d77981b584d8",
		ShareSize:   1,
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
		return
	}
}

func TestCreateFileShareAcl(t *testing.T) {
	fileShareAcl, err := fakeShareMgr.CreateFileShareAcl(&model.FileShareAclSpec{})
	if err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(fileShareAcl, &SampleFileSharesAcl[0]) {
		t.Errorf("expected %+v, got %+v", &SampleFileSharesAcl[0], fileShareAcl)
		return
	}
}

func TestDeleteFileShareAcl(t *testing.T) {
	var ShareAclID = "d2975ebe-d82c-430f-b28e-f373746a71ca"
	err := fakeShareMgr.DeleteFileShareAcl(ShareAclID)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestGetFileShareAcl(t *testing.T) {
	var ShareAclID = "d2975ebe-d82c-430f-b28e-f373746a71ca"
	expected := &model.FileShareAclSpec{
		BaseModel: &model.BaseModel{
			Id: "d2975ebe-d82c-430f-b28e-f373746a71ca",
		},
		Description: "This is a sample Acl for testing",
	}

	shareAcl, err := fakeShareMgr.GetFileShareAcl(ShareAclID)
	if err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(shareAcl, expected) {
		t.Errorf("expected %v, got %v", expected, shareAcl)
		return
	}
}

func TestListFileShareAcl(t *testing.T) {
	SampleFileSharesAcl = []model.FileShareAclSpec{
		{
			BaseModel: &model.BaseModel{
				Id: "d2975ebe-d82c-430f-b28e-f373746a71ca",
			},
			Description: "This is a sample Acl for testing",
		},
		{
			BaseModel: &model.BaseModel{
				Id: "1e643aca-4922-4b1a-bb98-4245054aeff4",
			},
			Description: "This is a sample Acl for testing",
		},
	}

	var expected []*model.FileShareAclSpec
	expected = append(expected, &SampleFileSharesAcl[0])
	expected = append(expected, &SampleFileSharesAcl[1])
	sharesAcl, err := fakeShareMgr.ListFileSharesAcl(map[string]string{"limit": "3", "offset": "4"})

	if err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(sharesAcl, expected) {
		t.Errorf("expected %v, got %v", expected, sharesAcl)
		return
	}
}
