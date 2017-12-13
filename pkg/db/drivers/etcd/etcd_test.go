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
	. "github.com/opensds/opensds/testutils/collection"
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
		resp = append(resp, StringSliceDocks[0])
	}
	if strings.Contains(req.Url, "pools") {
		resp = append(resp, StringSlicePools[0])
	}
	if strings.Contains(req.Url, "profiles") {
		resp = append(resp, StringSliceProfiles[0])
	}
	if strings.Contains(req.Url, "volumes") {
		resp = append(resp, StringSliceVolumes[0])
	}
	if strings.Contains(req.Url, "attachments") {
		resp = append(resp, StringSliceAttachments[0])
	}
	if strings.Contains(req.Url, "snapshots") {
		resp = append(resp, StringSliceSnapshots[0])
	}

	return &Response{
		Status:  "Success",
		Message: resp,
	}
}

func (*fakeClientCaller) List(req *Request) *Response {
	var resp []string

	if strings.Contains(req.Url, "docks") {
		resp = StringSliceDocks
	}
	if strings.Contains(req.Url, "pools") {
		resp = StringSlicePools
	}
	if strings.Contains(req.Url, "profiles") {
		resp = StringSliceProfiles
	}
	if strings.Contains(req.Url, "volumes") {
		resp = StringSliceVolumes
	}
	if strings.Contains(req.Url, "attachments") {
		resp = StringSliceAttachments
	}
	if strings.Contains(req.Url, "snapshots") {
		resp = StringSliceSnapshots
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
	if _, err := fc.CreateDock(&model.DockSpec{BaseModel: &model.BaseModel{}}); err != nil {
		t.Error("Create dock failed:", err)
	}
}

func TestCreatePool(t *testing.T) {
	if _, err := fc.CreatePool(&model.StoragePoolSpec{BaseModel: &model.BaseModel{}}); err != nil {
		t.Error("Create pool failed:", err)
	}
}

func TestCreateProfile(t *testing.T) {
	if _, err := fc.CreateProfile(&model.ProfileSpec{BaseModel: &model.BaseModel{}}); err != nil {
		t.Error("Create profile failed:", err)
	}
}

func TestCreateVolume(t *testing.T) {
	if _, err := fc.CreateVolume(&model.VolumeSpec{BaseModel: &model.BaseModel{}}); err != nil {
		t.Error("Create volume failed:", err)
	}
}

func TestCreateVolumeAttachment(t *testing.T) {
	if _, err := fc.CreateVolumeAttachment(&model.VolumeAttachmentSpec{BaseModel: &model.BaseModel{}}); err != nil {
		t.Error("Create volume attachment failed:", err)
	}
}

func TestCreateVolumeSnapshot(t *testing.T) {
	if _, err := fc.CreateVolumeSnapshot(&model.VolumeSnapshotSpec{BaseModel: &model.BaseModel{}}); err != nil {
		t.Error("Create volume snapshot failed:", err)
	}
}

func TestGetDock(t *testing.T) {
	dck, err := fc.GetDock("")
	if err != nil {
		t.Error("Get dock failed:", err)
	}

	var expected = &SampleDocks[0]
	if !reflect.DeepEqual(dck, expected) {
		t.Errorf("Expected %+v, got %+v\n", expected, dck)
	}
}

func TestGetPool(t *testing.T) {
	pol, err := fc.GetPool("")
	if err != nil {
		t.Error("Get pool failed:", err)
	}

	var expected = &SamplePools[0]
	if !reflect.DeepEqual(pol, expected) {
		t.Errorf("Expected %+v, got %+v\n", expected, pol)
	}
}

func TestGetProfile(t *testing.T) {
	prf, err := fc.GetProfile("")
	if err != nil {
		t.Error("Get profile failed:", err)
	}

	var expected = &SampleProfiles[0]
	if !reflect.DeepEqual(prf, expected) {
		t.Errorf("Expected %+v, got %+v\n", expected, prf)
	}
}

func TesGetVolume(t *testing.T) {
	vol, err := fc.GetVolume("")
	if err != nil {
		t.Error("Get volume failed:", err)
	}

	var expected = &SampleVolumes[0]
	if !reflect.DeepEqual(vol, expected) {
		t.Errorf("Expected %+v, got %+v\n", expected, vol)
	}
}

func TestGetVolumeAttachment(t *testing.T) {
	atc, err := fc.GetVolumeAttachment("")
	if err != nil {
		t.Error("Get volume attachment failed:", err)
	}

	var expected = &SampleAttachments[0]
	if !reflect.DeepEqual(atc, expected) {
		t.Errorf("Expected %+v, got %+v\n", expected, atc)
	}
}

func TestGetVolumeSnapshot(t *testing.T) {
	snp, err := fc.GetVolumeSnapshot("")
	if err != nil {
		t.Error("Get volume snapshot failed:", err)
	}

	var expected = &SampleSnapshots[0]
	if !reflect.DeepEqual(snp, expected) {
		t.Errorf("Expected %+v, got %+v\n", expected, snp)
	}
}

func TestListDocks(t *testing.T) {
	dcks, err := fc.ListDocks()
	if err != nil {
		t.Error("List docks failed:", err)
	}

	var expected []*model.DockSpec
	for i := range SampleDocks {
		expected = append(expected, &SampleDocks[i])
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

	var expected []*model.StoragePoolSpec
	for i := range SamplePools {
		expected = append(expected, &SamplePools[i])
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

	var expected []*model.ProfileSpec
	for i := range SampleProfiles {
		expected = append(expected, &SampleProfiles[i])
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

	var expected []*model.VolumeSpec
	for i := range SampleVolumes {
		expected = append(expected, &SampleVolumes[i])
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

	var expected []*model.VolumeAttachmentSpec
	for i := range SampleAttachments {
		expected = append(expected, &SampleAttachments[i])
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

	var expected []*model.VolumeSnapshotSpec
	for i := range SampleSnapshots {
		expected = append(expected, &SampleSnapshots[i])
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
