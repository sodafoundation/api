// Copyright (c) 2017 Huawei Technologies Co., Ltd. All Rights Reserved.
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

package integration

import (
	"encoding/json"
	"testing"

	"github.com/opensds/opensds/client"
	"github.com/opensds/opensds/pkg/model"
)

var c = client.NewClient(&client.Config{"http://localhost:50040"})

func TestClientCreateProfile(t *testing.T) {
	var body = &model.ProfileSpec{
		Name:        "silver",
		Description: "silver policy",
		Extra: model.ExtraSpec{
			"diskType":  "SAS",
			"iops":      300,
			"bandwidth": 500,
		},
	}

	prf, err := c.CreateProfile(body)
	if err != nil {
		t.Error("create profile in client failed:", err)
		return
	}

	prfBody, _ := json.MarshalIndent(prf, "", "	")
	t.Log(string(prfBody))
}

func TestClientGetProfile(t *testing.T) {
	var prfID = "2f9c0a04-66ef-11e7-ade2-43158893e017"

	prf, err := c.GetProfile(prfID)
	if err != nil {
		t.Error("get profile in client failed:", err)
		return
	}

	prfBody, _ := json.MarshalIndent(prf, "", "	")
	t.Log(string(prfBody))
}

func TestClientListProfiles(t *testing.T) {
	prfs, err := c.ListProfiles()
	if err != nil {
		t.Error("list profiles in client failed:", err)
		return
	}

	prfsBody, _ := json.MarshalIndent(prfs, "", "	")
	t.Log(string(prfsBody))
}

func TestClientDeleteProfile(t *testing.T) {
	var prfID = "2f9c0a04-66ef-11e7-ade2-43158893e017"

	res := c.DeleteProfile(prfID)
	if err := res.ToError(); err != nil {
		t.Error("delete profile in client failed:", err)
		return
	}

	resBody, _ := json.MarshalIndent(res, "", "	")
	t.Log(string(resBody))
}

func TestClientAddExtraProperty(t *testing.T) {
	var prfID = "2f9c0a04-66ef-11e7-ade2-43158893e017"
	var body = &model.ExtraSpec{
		"diskType":  "SAS",
		"iops":      300,
		"bandwidth": 500,
	}

	ext, err := c.AddExtraProperty(body)
	if err != nil {
		t.Error("add profile extra property in client failed:", err)
		return
	}

	extBody, _ := json.MarshalIndent(ext, "", "	")
	t.Log(string(extBody))
}

func TestClientListExtraProperties(t *testing.T) {
	var prfID = "2f9c0a04-66ef-11e7-ade2-43158893e017"

	ext, err := c.ListExtraProperties(prfID)
	if err != nil {
		t.Error("list profile extra properties in client failed:", err)
		return
	}

	extBody, _ := json.MarshalIndent(ext, "", "	")
	t.Log(string(extBody))
}

func TestClientRemoveExtraProperty(t *testing.T) {
	var prfID = "2f9c0a04-66ef-11e7-ade2-43158893e017"
	var extraKey = "iops"

	res := c.RemoveExtraProperty(prfID, extraKey)
	if err := res.ToError(); err != nil {
		t.Error("remove profile extra property in client failed:", err)
		return
	}

	resBody, _ := json.MarshalIndent(res, "", "	")
	t.Log(string(resBody))
}

func TestClientGetDock(t *testing.T) {
	var dckID = "b7602e18-771e-11e7-8f38-dbd6d291f4e0"

	dck, err := c.GetDock(dckID)
	if err != nil {
		t.Error("get dock in client failed:", err)
		return
	}

	dckBody, _ := json.MarshalIndent(dck, "", "	")
	t.Log(string(dckBody))
}

func TestClientListDocks(t *testing.T) {
	dcks, err := c.ListDocks()
	if err != nil {
		t.Error("list docks in client failed:", err)
		return
	}

	dcksBody, _ := json.MarshalIndent(dcks, "", "	")
	t.Log(string(dcksBody))
}

func TestClientGetPool(t *testing.T) {
	var polID = "084bf71e-a102-11e7-88a8-e31fe6d52248"

	pol, err := c.GetPool(polID)
	if err != nil {
		t.Error("get pool in client failed:", err)
		return
	}

	polBody, _ := json.MarshalIndent(pol, "", "	")
	t.Log(string(polBody))
}

func TestClientListPools(t *testing.T) {
	pols, err := c.ListPools()
	if err != nil {
		t.Error("list pools in client failed:", err)
		return
	}

	polsBody, _ := json.MarshalIndent(pols, "", "	")
	t.Log(string(polsBody))
}

func TestClientCreateVolume(t *testing.T) {
	var body = &model.VolumeSpec{
		Name:        "test",
		Description: "This is a test",
		Size:        int64(1),
	}

	vol, err := c.CreateVolume(body)
	if err != nil {
		t.Error("create volume in client failed:", err)
		return
	}

	volBody, _ := json.MarshalIndent(vol, "", "	")
	t.Log(string(volBody))
}

func TestClientGetVolume(t *testing.T) {
	var volID = "bd5b12a8-a101-11e7-941e-d77981b584d8"

	vol, err := c.GetVolume(volID)
	if err != nil {
		t.Error("get volume in client failed:", err)
		return
	}

	volBody, _ := json.MarshalIndent(vol, "", "	")
	t.Log(string(volBody))
}

func TestClientListVolumes(t *testing.T) {
	vols, err := c.ListVolumes()
	if err != nil {
		t.Error("list volumes in client failed:", err)
		return
	}

	volsBody, _ := json.MarshalIndent(vols, "", "	")
	t.Log(string(volsBody))
}

func TestClientDeleteVolume(t *testing.T) {
	var volID = "bd5b12a8-a101-11e7-941e-d77981b584d8"
	body := &model.VolumeSpec{}

	res := c.DeleteVolume(volID, body)
	if err := res.ToError(); err != nil {
		t.Error("delete volume in client failed:", err)
		return
	}

	resBody, _ := json.MarshalIndent(res, "", "	")
	t.Log(string(resBody))
}

func TestClientCreateVolumeSnapshot(t *testing.T) {
	var body = &model.VolumeSnapshotSpec{
		Name:        "test",
		Description: "This is a test",
		VolumeId:    "bd5b12a8-a101-11e7-941e-d77981b584d8",
	}

	snp, err := c.CreateVolumeSnapshot(body)
	if err != nil {
		t.Error("create volume snapshot in client failed:", err)
		return
	}

	snpBody, _ := json.MarshalIndent(snp, "", "	")
	t.Log(string(snpBody))
}

func TestClientGetVolumeSnapshot(t *testing.T) {
	var snpID = "3769855c-a102-11e7-b772-17b880d2f537"

	snp, err := c.GetVolumeSnapshot(snpID)
	if err != nil {
		t.Error("get volume snapshot in client failed:", err)
		return
	}

	snpBody, _ := json.MarshalIndent(snp, "", "	")
	t.Log(string(snpBody))
}

func TestClientListVolumeSnapshots(t *testing.T) {
	snps, err := c.ListVolumeSnapshots()
	if err != nil {
		t.Error("list volume snapshots in client failed:", err)
		return
	}

	snpsBody, _ := json.MarshalIndent(snps, "", "	")
	t.Log(string(snpsBody))
}

func TestClientDeleteVolumeSnapshot(t *testing.T) {
	var snpID = "3769855c-a102-11e7-b772-17b880d2f537"
	body := &model.VolumeSnapshotSpec{
		VolumeId: "bd5b12a8-a101-11e7-941e-d77981b584d8",
	}

	res := c.DeleteVolumeSnapshot(snpID, body)
	if err := res.ToError(); err != nil {
		t.Error("delete volume snapshot in client failed:", err)
		return
	}

	resBody, _ := json.MarshalIndent(res, "", "	")
	t.Log(string(resBody))
}
