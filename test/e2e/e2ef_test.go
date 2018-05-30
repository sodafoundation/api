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

// +build e2ef

package e2e

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/opensds/opensds/client"
	"github.com/opensds/opensds/pkg/model"
	"github.com/opensds/opensds/pkg/utils/constants"
)

var u = client.NewClient(&client.Config{
	Endpoint:    "http://localhost:50040",
	AuthOptions: client.NewNoauthOptions(constants.DefaultTenantId)})

//init Create Profile
func init() {
	var body = &model.ProfileSpec{
		Name:        "default",
		Description: "default policy",
	}
	cprf, err := u.CreateProfile(body)
	if err != nil {
		fmt.Println("create profile in client failed:", err)
		return
	}
	prfBody, _ := json.MarshalIndent(cprf, "", "    ")
	fmt.Println(string(prfBody))
}

//Test Case: List Profile
func TestListProfile(t *testing.T) {
	t.Log("Begin to List Profile")
	listRs, err := u.ListProfiles()
	if err != nil {
		t.Error("list profiles in client failed:", err)
		return
	}
	prflist, _ := json.MarshalIndent(listRs, "", "    ")
	t.Log(string(prflist))
	t.Log("List Profile Success")
}

//Test Case: Get Profile Detail
func TestGetProfileDetail(t *testing.T) {
	t.Log("Begin to Get Profile Detail")
	pro, err := PrepareProfile()
	if err != nil {
		t.Error("Prepare Profile Fail!", err)
		return
	}
	defer DeleteProfile(t, pro.Id)
	detail, err := u.GetProfile(pro.Id)
	if err != nil {
		t.Error("get profile in client failed:", err)
		return
	}
	prfdet, _ := json.MarshalIndent(detail, "", "    ")
	t.Log(string(prfdet))
	t.Log("Get Profile Detail Success")
}

//Test Case:Delete Profile
func TestDeleteProfile(t *testing.T) {
	t.Log("Begin to Delete Profile...")
	pro, err := PrepareProfile()
	if err != nil {
		t.Error("Prepare Profile Fail!", err)
		return
	}
	errdel := u.DeleteProfile(pro.Id)
	if errdel != nil {
		t.Error("delete profile in client failed:", err)
		return
	}
	t.Log("Delete Profile Success")
}

//Test Case:Create Volume
func TestCreateVolumeFlow(t *testing.T) {
	t.Log("Begin to Create Volume.....")
	var volbody = &model.VolumeSpec{
		Name:        "flowVolume",
		Description: "This a volume for test",
		Size:        int64(1),
	}
	create, err := u.CreateVolume(volbody)
	if err != nil {
		t.Error("Create Volume Fail", err)
		return
	}
	defer DeleteVolume(create.Id)
	vol, _ := json.MarshalIndent(create, "", "    ")
	t.Log(string(vol))
	t.Log("Create Volume Success")
}

//Test Case:Update Volume
func TestUpdateVolumeFlow(t *testing.T) {
	t.Log("Begin to Update Volume...")
	vol, _ := PrepareVolume()
	body := &model.VolumeSpec{
		Name:        "UpdateName",
		Description: "Update Description",
	}
	upvol, err := u.UpdateVolume(vol.Id, body)
	if err != nil {
		t.Error("Update Volume Fail", err)
		return
	}
	defer DeleteVolume(vol.Id)
	volrs, _ := json.MarshalIndent(upvol, "", "    ")
	t.Log(string(volrs))
	t.Log("Update Volume Success")
}

//Test Case:Extend Volume
func TestExtendVolumeFlow(t *testing.T) {
	t.Log("Begin to Extend Volume...")
	vol, _ := PrepareVolume()
	body := &model.ExtendVolumeSpec{
		NewSize: int64(2),
	}
	ext, err := u.ExtendVolume(vol.Id, body)
	t.Log("err:", err, "ext.Size", ext.Size)
	if err != nil {
		t.Error("Extend volume fail", err)
	}
	defer DeleteVolume(vol.Id)
	extrs, _ := json.MarshalIndent(ext, "", "    ")
	t.Log(string(extrs))
	t.Log("Extend Volume Success")

}

//Test Case:List Volume
func TestGetVolumeList(t *testing.T) {
	t.Log("Begin to List Volume....")
	vol, err := PrepareVolume()
	if err != nil {
		t.Error("Prepare Volume Fail", err)
		return
	}
	defer DeleteVolume(vol.Id)
	vols, err2 := u.ListVolumes()
	if err2 != nil {
		t.Error("List Volume Fail", err2)
		return
	}
	list, _ := json.MarshalIndent(vols, "", "    ")
	t.Log(string(list))
	t.Log("List Volume Success")
}

//Test Case:Get Volume Detail
func TestGetVolumeDetail(t *testing.T) {
	t.Log("Begin to Get Volume Detail....")
	vol, err := PrepareVolume()
	if err != nil {
		t.Error("Prepare Volume Fail", err)
		return
	}
	defer DeleteVolume(vol.Id)
	Info, err2 := u.GetVolume(vol.Id)
	if err2 != nil {
		t.Error("Get Volume Detail Fail", err2)
		return
	}
	detail, _ := json.MarshalIndent(Info, "", "    ")
	t.Log(string(detail))
	t.Log("Get Volume Detail Success")
}

//Test Case:Delete Volume
func TestDeleteVolume(t *testing.T) {
	t.Log("Begin to delete volume ....")
	vol, err := PrepareVolume()
	if err != nil {
		t.Error("Prepare Volume Fail", err)
		return
	}
	errd := u.DeleteVolume(vol.Id, nil)
	if errd != nil {
		t.Error("Delete Volume fail", errd)
		return
	}
	t.Log("Delete Volume Success!")
}

//Test Case:Create Snapshot
func TestCreateSnapshot(t *testing.T) {
	vol, err := PrepareVolume()
	if err != nil {
		t.Error("Prepare Volume Fail", err)
		return
	}
	defer DeleteVolume(vol.Id)
	t.Log("Start preparing volume snapshot...")
	t.Log("Volume info:")
	volInfo, _ := json.MarshalIndent(vol, "", "    ")
	t.Log(string(volInfo))
	var body = &model.VolumeSnapshotSpec{
		Name:        "new-snapshot",
		Description: "This is a snapshot test for new",
		VolumeId:    vol.Id,
	}
	snp, errs := u.CreateVolumeSnapshot(body)
	if errs != nil {
		t.Error("prepare volume snapshot failed:", errs)
		//cleanVolumeIfFailedOrFinished(t, snp.VolumeId)
		return
	}
	defer DeleteSnapshot(snp.Id)
	snap, _ := json.MarshalIndent(snp, "", "    ")
	t.Log(string(snap))
	t.Log("End preparing volume snapshot...")
}

//Test Case:Update Snapshot
func TestUpdateSnapshot(t *testing.T) {
	snap, err := PrepareSnapshot()
	if err != nil {
		t.Error("Prepare Snapshot Fail!", err)
		return
	}
	defer DeleteVolume(snap.VolumeId)
	defer DeleteSnapshot(snap.Id)

	body := &model.VolumeSnapshotSpec{
		Name:        "UpSnapshot",
		Description: "Update Snapshot Description",
	}
	upsnap, errup := u.UpdateVolumeSnapshot(snap.Id, body)
	if errup != nil {
		t.Error("update volume snapshot failed:", errup)
		return
	}
	upsnaprs, _ := json.MarshalIndent(upsnap, "", "    ")
	t.Log(string(upsnaprs))
	t.Log("-----Update Snapshot Success-----")
}

//Test Case:List Snapshot
func TestListSnapshot(t *testing.T) {
	snap, err := PrepareSnapshot()
	if err != nil {
		t.Error("Prepare Snapshot Fail!", err)
		return
	}
	defer DeleteVolume(snap.VolumeId)
	defer DeleteSnapshot(snap.Id)

	snpli, errli := u.ListVolumeSnapshots()
	if errli != nil {
		t.Error("List Snapshot Fail", errli)
		return
	}
	snapli, _ := json.MarshalIndent(snpli, "", "    ")
	t.Log(string(snapli))
	t.Log("-----Update Snapshot Success-----")
}

//Test Case:Get Snapshot detail
func TestGetSnapDetail(t *testing.T) {
	snap, err := PrepareSnapshot()
	if err != nil {
		t.Error("Prepare Snapshot Fail!", err)
		return
	}
	defer DeleteVolume(snap.VolumeId)
	defer DeleteSnapshot(snap.Id)

	snapinfo, err2 := u.GetVolumeSnapshot(snap.Id)
	if err2 != nil {
		t.Error("Get Snapshot Detail Fail!", err2)
		return
	}
	sndetail, _ := json.MarshalIndent(snapinfo, "", "    ")
	t.Log(string(sndetail))
	t.Log("-----Get Snapshot Detail Success-----")
}

//Test Case:Delete Snapshot
func TestDeleteSnapshot(t *testing.T) {
	snap, err := PrepareSnapshot()
	if err != nil {
		t.Error("Prepare Snapshot Fail!", err)
		return
	}
	defer DeleteVolume(snap.VolumeId)
	err2 := u.DeleteVolumeSnapshot(snap.Id, nil)
	if err2 != nil {
		t.Error("Delete Snapshot Fail!", err2)
		return
	}
	t.Log("Delete Snapshot Success!")
}

//Test Case:Create Volume Attachment
func TestCreateAttach(t *testing.T) {
	vol, err := PrepareVolume()
	if err != nil {
		t.Error("Prepare Volume Fail", err)
		return
	}
	defer DeleteVolume(vol.Id)
	var body = &model.VolumeAttachmentSpec{
		VolumeId: vol.Id,
		HostInfo: model.HostInfo{},
	}
	attc, err2 := u.CreateVolumeAttachment(body)
	if err2 != nil {
		t.Error("create volume attachment failed:", err2)
		return
	}
	attrs, _ := json.MarshalIndent(attc, "", "    ")
	t.Log(string(attrs))
	t.Log("Create Volume Attachment Success")
}

//Test Case:List Attachment
func TestListAttachment(t *testing.T) {
	attc, errp := PrepareAttachment(t)
	if errp != nil {
		t.Error("Prepare Attachment Fail!", errp)
		return
	}
	defer DeleteVolume(attc.VolumeId)
	defer DeleteAttachment(attc.Id)

	atts, err := u.ListVolumeAttachments()
	if err != nil {
		t.Error("List Attachment Error!", err)
		return
	}
	attli, _ := json.MarshalIndent(atts, "", "    ")
	t.Log(string(attli))
	t.Log("List Attachment Success!")
}

//Test Case:Get Attachment Detail
func TestShowAttachDetail(t *testing.T) {
	attc, errp := PrepareAttachment(t)
	if errp != nil {
		t.Error("Prepare Attachment Fail!", errp)
		return
	}
	defer DeleteVolume(attc.VolumeId)
	defer DeleteAttachment(attc.Id)
	getatt, err := u.GetVolumeAttachment(attc.Id)
	t.Log("Begin to Scan Volume:")
	out := ScanVolume()
	t.Log(string(out))
	t.Log("Scan Volume End!")
	//read Dsik.log
	dev := DiskChk(out, "/dev/sd")
	ca := DiskChk(out, "1 GiB")
	t.Log(err)
	t.Log(dev)
	t.Log(ca)
	if err != nil || dev == false || ca == false {
		t.Error("Volume attachment detail check fail", err)
		return
	}
	detail, _ := json.MarshalIndent(getatt, "", "    ")
	t.Log(string(detail))
	t.Log("Volume attach detail Check Success!")
}

//Test Case:Delete Attachment
func TestDeleteAttach(t *testing.T) {
	attc, errp := PrepareAttachment(t)
	if errp != nil {
		t.Error("Prepare Attachment Fail!", errp)
		return
	}
	defer DeleteVolume(attc.VolumeId)
	err := u.DeleteVolumeAttachment(attc.Id, nil)
	if err != nil {
		t.Error("Delete Attachment Fail", err)
		return
	}
	t.Log("Begin to Scan Volume:")
	out := ScanVolume()
	t.Log(out)
	t.Log("Scan Volume End!")
	tar := DiskChk(out, "login target: -p")
	//dev := DiskChk(out, "/dev/sd")
	//ca := DiskChk(out, "1 GiB")
	chk2, _ := u.GetVolumeAttachment(attc.Id)
	if err != nil || chk2 != nil || tar != true {
		t.Error("Delete Attachment Fail", err)
	}
	t.Log("Delete attachment Success")
}

//check Attachemnt By scan volume
func ScanVolume() string {
	cmd := exec.Command("/bin/bash", "./scanvolume.sh")
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("cmd.Output: ", err)
		return ""
	}
	return string(out)
}

//Check if Disk Log contain /dev/sd&& 2 GiB
func DiskChk(log string, str string) bool {
	if strings.Index(log, str) != -1 {
		return true
	}
	return false
}

//prepare attachment
func PrepareAttachment(t *testing.T) (*model.VolumeAttachmentSpec, error) {
	vol, err := PrepareVolume()
	if err != nil {
		t.Log("Prepare Volume Fail", err)
		return nil, err
	}
	defer DeleteVolume(vol.Id)
	var body = &model.VolumeAttachmentSpec{
		VolumeId: vol.Id,
		HostInfo: model.HostInfo{},
	}
	attc, err2 := u.CreateVolumeAttachment(body)
	if err2 != nil {
		t.Error("prepare volume attachment failed:", err2)
		return nil, err2
	}
	attrs, _ := json.MarshalIndent(attc, "", "    ")
	t.Log(string(attrs))
	t.Log("prepare Volume Attachment Success")
	return attc, nil
}

//delete attachment
func DeleteAttachment(attId string) {
	err := u.DeleteVolumeAttachment(attId, nil)
	//check if attachment is exist after
	for i := 0; i < 5; i++ {
		chk, _ := u.GetVolumeAttachment(attId)
		if chk == nil {
			break
		} else {
			time.Sleep(1 * 1e9)
		}
	}
	if err != nil {
		fmt.Println("Delete Attachment Fail", err)
	}
	fmt.Println("Delete Attachment Success!")
}

//prepare sanpshot
func PrepareSnapshot() (*model.VolumeSnapshotSpec, error) {
	vol, err := PrepareVolume()
	if err != nil {
		fmt.Println("Prepare Volume Fail", err)
		return nil, err
	}
	fmt.Println("Start preparing volume snapshot...")
	var body = &model.VolumeSnapshotSpec{
		Name:        "flow-snapshot",
		Description: "This is a snapshot test for flow",
		VolumeId:    vol.Id,
	}
	snp, errs := u.CreateVolumeSnapshot(body)
	if errs != nil {
		fmt.Println("prepare volume snapshot failed:", errs)
		//cleanVolumeIfFailedOrFinished(t, snp.VolumeId)
		return nil, errs
	}
	fmt.Println("End preparing volume snapshot...")
	return snp, nil
}

//delete snapshot
func DeleteSnapshot(snpId string) {
	if err := u.DeleteVolumeSnapshot(snpId, nil); err != nil {
		fmt.Println("Delete Snapshot Fail")
		return
	}
	fmt.Println("Delete Snapshot Success")
}

//prepare volume for test
func PrepareVolume() (*model.VolumeSpec, error) {
	var volbody = &model.VolumeSpec{
		Name:        "flowTest",
		Description: "This a test for flow",
		Size:        int64(1),
	}
	create, err := u.CreateVolume(volbody)
	if err != nil {
		fmt.Println("Prepare Volume Fail")
		return nil, err
	}
	volrs, _ := json.MarshalIndent(create, "", "    ")
	fmt.Println(string(volrs))
	fmt.Println("Prepare Volume Success")
	return create, nil
}

//Delete volume after test
func DeleteVolume(volId string) error {
	err := u.DeleteVolume(volId, nil)
	if err != nil {
		fmt.Println("Delete Volume fail")
		return err
	}
	fmt.Println("Delete Volume Success!")
	return nil
}

//prepare profile for Test
func PrepareProfile() (*model.ProfileSpec, error) {
	var body = &model.ProfileSpec{
		Name:        "policy",
		Description: "test policy",
	}
	cprf, err := u.CreateProfile(body)
	if err != nil {
		return nil, err
	}
	return cprf, nil
}

//Delete Profile
func DeleteProfile(t *testing.T, proId string) error {
	err := u.DeleteProfile(proId)
	if err != nil {
		fmt.Println("delete profile in client failed:", err)
		return err
	}
	fmt.Println("Delete Profile Success")
	return nil
}
