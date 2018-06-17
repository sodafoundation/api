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

// +build integration

package integration

import (
	"reflect"
	"testing"

	"github.com/opensds/opensds/client"
	"github.com/opensds/opensds/pkg/model"
	"github.com/opensds/opensds/pkg/utils/constants"
	. "github.com/opensds/opensds/testutils/collection"
)

var c = client.NewClient(&client.Config{
	Endpoint:    "http://localhost:50040",
	AuthOptions: client.NewNoauthOptions(constants.DefaultTenantId)})

func TestClientCreateProfile(t *testing.T) {
	var body = &model.ProfileSpec{
		Name:        "silver",
		Description: "silver policy",
		Extras: model.ExtraSpec{
			"diskType": "SAS",
		},
	}

	prf, err := c.CreateProfile(body)
	if err != nil {
		t.Error("create profile in client failed:", err)
		return
	}
	// If extras are not defined, create an empty one.
	if prf.Extras == nil {
		prf.Extras = model.ExtraSpec{}
	}

	var expected = &SampleProfiles[0]
	if !reflect.DeepEqual(prf, expected) {
		t.Errorf("expected %+v, got %+v\n", expected, prf)
	}
}

func TestClientGetProfile(t *testing.T) {
	var prfID = "2f9c0a04-66ef-11e7-ade2-43158893e017"

	prf, err := c.GetProfile(prfID)
	if err != nil {
		t.Error("get profile in client failed:", err)
		return
	}

	var expected = &SampleProfiles[1]
	if !reflect.DeepEqual(prf, expected) {
		t.Errorf("expected %+v, got %+v\n", expected, prf)
	}
}

func TestClientListProfiles(t *testing.T) {
	prfs, err := c.ListProfiles()
	if err != nil {
		t.Error("list profiles in client failed:", err)
		return
	}
	// If extras are not defined, create an empty one.
	for _, prf := range prfs {
		if prf.Extras == nil {
			prf.Extras = model.ExtraSpec{}
		}
	}

	var expected []*model.ProfileSpec
	for i := range SampleProfiles {
		expected = append(expected, &SampleProfiles[i])
	}
	if !reflect.DeepEqual(prfs, expected) {
		t.Errorf("expected %+v, got %+v\n", expected, prfs)
	}
}

func TestClientDeleteProfile(t *testing.T) {
	var prfID = "2f9c0a04-66ef-11e7-ade2-43158893e017"

	if err := c.DeleteProfile(prfID); err != nil {
		t.Error("delete profile in client failed:", err)
		return
	}

	t.Log("Delete profile success!")
}

func TestClientAddExtraProperty(t *testing.T) {
	var prfID = "2f9c0a04-66ef-11e7-ade2-43158893e017"
	var body = &model.ExtraSpec{
		"diskType": "SAS",
	}

	ext, err := c.AddExtraProperty(prfID, body)
	if err != nil {
		t.Error("add profile extra property in client failed:", err)
		return
	}

	var expected = &SampleProfiles[0].Extras
	if !reflect.DeepEqual(ext, expected) {
		t.Errorf("expected %+v, got %+v\n", expected, ext)
	}
}

func TestClientListExtraProperties(t *testing.T) {
	var prfID = "2f9c0a04-66ef-11e7-ade2-43158893e017"

	ext, err := c.ListExtraProperties(prfID)
	if err != nil {
		t.Error("list profile extra properties in client failed:", err)
		return
	}

	var expected = &SampleProfiles[0].Extras
	if !reflect.DeepEqual(ext, expected) {
		t.Errorf("expected %+v, got %+v\n", expected, ext)
	}
}

func TestClientRemoveExtraProperty(t *testing.T) {
	var prfID = "2f9c0a04-66ef-11e7-ade2-43158893e017"
	var extraKey = "iops"

	if err := c.RemoveExtraProperty(prfID, extraKey); err != nil {
		t.Error("remove profile extra property in client failed:", err)
		return
	}

	t.Log("Remove extra property success!")
}

func TestClientGetDock(t *testing.T) {
	var dckID = "b7602e18-771e-11e7-8f38-dbd6d291f4e0"

	dck, err := c.GetDock(dckID)
	if err != nil {
		t.Error("get dock in client failed:", err)
		return
	}

	var expected = &SampleDocks[0]
	if !reflect.DeepEqual(dck, expected) {
		t.Errorf("expected %+v, got %+v\n", expected, dck)
	}
}

func TestClientListDocks(t *testing.T) {
	dcks, err := c.ListDocks()
	if err != nil {
		t.Error("list docks in client failed:", err)
		return
	}

	var expected []*model.DockSpec
	for i := range SampleDocks {
		expected = append(expected, &SampleDocks[i])
	}
	if !reflect.DeepEqual(dcks, expected) {
		t.Errorf("expected %+v, got %+v\n", expected, dcks)
	}
}

func TestClientGetPool(t *testing.T) {
	var polID = "084bf71e-a102-11e7-88a8-e31fe6d52248"

	pol, err := c.GetPool(polID)
	if err != nil {
		t.Error("get pool in client failed:", err)
		return
	}

	var expected = &SamplePools[0]
	if !reflect.DeepEqual(pol, expected) {
		t.Errorf("expected %+v, got %+v\n", expected, pol)
	}
}

func TestClientListPools(t *testing.T) {
	pols, err := c.ListPools()
	if err != nil {
		t.Error("list pools in client failed:", err)
		return
	}

	var expected []*model.StoragePoolSpec
	for i := range SamplePools {
		expected = append(expected, &SamplePools[i])
	}
	if !reflect.DeepEqual(pols, expected) {
		t.Errorf("expected %+v, got %+v\n", expected, pols)
	}
}

func TestClientCreateVolume(t *testing.T) {
	var body = &model.VolumeSpec{
		Name:        "test",
		Description: "This is a test",
		Size:        int64(1),
	}

	if _, err := c.CreateVolume(body); err != nil {
		t.Error("create volume in client failed:", err)
		return
	}

	t.Log("Create volume success!")
}

func TestClientGetVolume(t *testing.T) {
	var volID = "bd5b12a8-a101-11e7-941e-d77981b584d8"

	vol, err := c.GetVolume(volID)
	if err != nil {
		t.Error("get volume in client failed:", err)
		return
	}

	var expected = &SampleVolumes[0]
	if !reflect.DeepEqual(vol, expected) {
		t.Errorf("expected %+v, got %+v\n", expected, vol)
	}
}

func TestClientListVolumes(t *testing.T) {
	vols, err := c.ListVolumes()
	if err != nil {
		t.Error("list volumes in client failed:", err)
		return
	}

	var expected []*model.VolumeSpec
	for i := range SampleVolumes {
		expected = append(expected, &SampleVolumes[i])
	}
	if !reflect.DeepEqual(vols, expected) {
		t.Errorf("expected %+v, got %+v\n", expected, vols)
	}
}

func TestClientUpdateVolume(t *testing.T) {
	var volID = "bd5b12a8-a101-11e7-941e-d77981b584d8"
	body := &model.VolumeSpec{
		Name:        "sample-volume",
		Description: "This is a sample volume for testing",
	}

	vol, err := c.UpdateVolume(volID, body)
	if err != nil {
		t.Error("update volume in client failed:", err)
		return
	}

	var expected = &SampleVolumes[0]
	if !reflect.DeepEqual(vol, expected) {
		t.Errorf("expected %+v, got %+v\n", expected, vol)
	}
}

func TestClientExtendVolume(t *testing.T) {
	var volID = "bd5b12a8-a101-11e7-941e-d77981b584d8"

	oldVol, err := c.GetVolume(volID)
	if err != nil {
		t.Error("get volume in client failed:", err)
		return
	}

	body := &model.ExtendVolumeSpec{
		NewSize: int64(oldVol.Size + 1),
	}
	if _, err := c.ExtendVolume(volID, body); err != nil {
		t.Error("extend volume in client failed:", err)
		return
	}

	t.Log("Extend volume success!")
}

func TestClientCreateVolumeAttachment(t *testing.T) {
	var body = &model.VolumeAttachmentSpec{
		VolumeId: "bd5b12a8-a101-11e7-941e-d77981b584d8",
		HostInfo: model.HostInfo{},
	}

	if _, err := c.CreateVolumeAttachment(body); err != nil {
		t.Error("create volume attachment in client failed:", err)
		return
	}

	t.Log("Create volume attachment success!")
}

func TestClientGetVolumeAttachment(t *testing.T) {
	var atcID = "f2dda3d2-bf79-11e7-8665-f750b088f63e"

	atc, err := c.GetVolumeAttachment(atcID)
	if err != nil {
		t.Error("get volume attachment in client failed:", err)
		return
	}

	var expected = &SampleAttachments[0]
	if !reflect.DeepEqual(atc, expected) {
		t.Errorf("expected %+v, got %+v\n", expected, atc)
	}
}

func TestClientListVolumeAttachments(t *testing.T) {
	atcs, err := c.ListVolumeAttachments()
	if err != nil {
		t.Error("list volume attachments in client failed:", err)
		return
	}

	var expected []*model.VolumeAttachmentSpec
	for i := range SampleAttachments {
		expected = append(expected, &SampleAttachments[i])
	}
	if !reflect.DeepEqual(atcs, expected) {
		t.Errorf("expected %+v, got %+v\n", expected, atcs)
	}
}

func TestClientDeleteVolumeAttachment(t *testing.T) {
	var atcID = "f2dda3d2-bf79-11e7-8665-f750b088f63e"

	if err := c.DeleteVolumeAttachment(atcID, nil); err != nil {
		t.Error("delete volume attachment in client failed:", err)
		return
	}

	t.Log("Delete volume attachment success!")
}

func TestClientCreateVolumeSnapshot(t *testing.T) {
	var body = &model.VolumeSnapshotSpec{
		Name:        "test",
		Description: "This is a test",
		VolumeId:    "bd5b12a8-a101-11e7-941e-d77981b584d8",
	}

	if _, err := c.CreateVolumeSnapshot(body); err != nil {
		t.Error("create volume snapshot in client failed:", err)
		return
	}

	t.Log("Create volume snapshot success!")
}

func TestClientGetVolumeSnapshot(t *testing.T) {
	var snpID = "3769855c-a102-11e7-b772-17b880d2f537"

	snp, err := c.GetVolumeSnapshot(snpID)
	if err != nil {
		t.Error("get volume snapshot in client failed:", err)
		return
	}

	var expected = &SampleSnapshots[0]
	if !reflect.DeepEqual(snp, expected) {
		t.Errorf("expected %+v, got %+v\n", expected, snp)
	}
}

func TestClientListVolumeSnapshots(t *testing.T) {
	snps, err := c.ListVolumeSnapshots()
	if err != nil {
		t.Error("list volume snapshots in client failed:", err)
		return
	}

	var expected []*model.VolumeSnapshotSpec
	for i := range SampleSnapshots {
		expected = append(expected, &SampleSnapshots[i])
	}
	if !reflect.DeepEqual(snps, expected) {
		t.Errorf("expected %+v, got %+v\n", expected, snps)
	}
}

func TestClientDeleteVolumeSnapshot(t *testing.T) {
	var snpID = "3769855c-a102-11e7-b772-17b880d2f537"

	if err := c.DeleteVolumeSnapshot(snpID, nil); err != nil {
		t.Error("delete volume snapshot in client failed:", err)
		return
	}

	t.Log("Delete volume snapshot success!")
}

func TestClientUpdateVolumeSnapshot(t *testing.T) {
	var snpID = "3769855c-a102-11e7-b772-17b880d2f537"
	body := &model.VolumeSnapshotSpec{
		Name:        "sample-snapshot-01",
		Description: "This is the first sample snapshot for testing",
	}

	snp, err := c.UpdateVolumeSnapshot(snpID, body)
	if err != nil {
		t.Error("update volume snapshot in client failed:", err)
		return
	}

	var expected = &SampleSnapshots[0]
	if !reflect.DeepEqual(snp, expected) {
		t.Errorf("expected %+v, got %+v\n", expected, snp)
	}
}

func TestClientDeleteVolume(t *testing.T) {
	var volID = "bd5b12a8-a101-11e7-941e-d77981b584d8"
	body := &model.VolumeSpec{}

	if err := c.DeleteVolume(volID, body); err != nil {
		t.Error("delete volume in client failed:", err)
		return
	}

	t.Log("Delete volume success!")
}

// TODO: There are some deployment issues when testing Replicaiton operation,
// so these test cases would be hidden until we fix the bug.
/*
func TestClientCreateReplication(t *testing.T) {
	var body = &model.ReplicationSpec{
		Name:              "sample-replication-01",
		Description:       "This is a sample replication for testing",
		PrimaryVolumeId:   "bd5b12a8-a101-11e7-941e-d77981b584d8",
		SecondaryVolumeId: "bd5b12a8-a101-11e7-941e-d77981b584d8",
		ReplicationMode:   model.ReplicationModeSync,
	}

	replica, err := c.CreateReplication(body)
	if err != nil {
		t.Error("create volume replication in client failed:", err)
		return
	}

	replicaBody, _ := json.MarshalIndent(replica, "", "	")
	t.Log(string(replicaBody))
}

func TestClientGetReplication(t *testing.T) {
	var replicaID = "c299a978-4f3e-11e8-8a5c-977218a83359"

	replica, err := c.GetReplication(replicaID)
	if err != nil {
		t.Error("get volume replication in client failed:", err)
		return
	}

	replicaBody, _ := json.MarshalIndent(replica, "", "	")
	t.Log(string(replicaBody))
}

func TestClientListReplications(t *testing.T) {
	replicas, err := c.ListReplications()
	if err != nil {
		t.Error("list volume replications in client failed:", err)
		return
	}

	replicasBody, _ := json.MarshalIndent(replicas, "", "	")
	t.Log(string(replicasBody))
}

func TestClientUpdateReplication(t *testing.T) {
	var replicaID = "c299a978-4f3e-11e8-8a5c-977218a83359"
	body := &model.ReplicationSpec{
		Name:        "sample-replication-02",
		Description: "This is a super-cool replication for testing",
	}

	replica, err := c.UpdateReplication(replicaID, body)
	if err != nil {
		t.Error("update volume replication in client failed:", err)
		return
	}

	replicaBody, _ := json.MarshalIndent(replica, "", "	")
	t.Log(string(replicaBody))
}

func TestClientDeleteReplication(t *testing.T) {
	var replicaID = "c299a978-4f3e-11e8-8a5c-977218a83359"

	if err := c.DeleteReplication(replicaID, nil); err != nil {
		t.Error("delete volume replicaiton in client failed:", err)
		return
	}

	t.Log("Delete volume replication success!")
}

func TestClientEnableReplication(t *testing.T) {
	var replicaID = "c299a978-4f3e-11e8-8a5c-977218a83359"

	if err := c.EnableReplication(replicaID); err != nil {
		t.Error("enable volume replication in client failed:", err)
		return
	}

	t.Log("Enable volume replication success!")
}

func TestClientDisableReplication(t *testing.T) {
	var replicaID = "c299a978-4f3e-11e8-8a5c-977218a83359"

	if err := c.DisableReplication(replicaID); err != nil {
		t.Error("disable volume replicaiton in client failed:", err)
		return
	}

	t.Log("Disable volume replication success!")
}

func TestClientFailoverReplication(t *testing.T) {
	// TODO Add TestClientFailoverRelication method.

	t.Log("Disable volume replication not ready!")
}
*/
