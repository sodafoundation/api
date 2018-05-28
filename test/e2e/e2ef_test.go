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

// +build e2e

package e2e

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/opensds/opensds/client"
	"github.com/opensds/opensds/pkg/model"
	"github.com/opensds/opensds/pkg/utils/constants"
)

var VOLNAME = "flowvolume"
var VOLDESC = "description"
var UPDATENAME = "Update Name"
var UPDATEDESC = "Update Desc"
var EXTENDSIZE int64 = 2
var SNAPNAME = "flowsnap"
var SNAPDESC = "snapdesc"
var UPDATESNAPNAME = "UpdateSnapName"
var UPDATESNAPDESC = "UpdateSnapDesc"

//Get Profile ID
func GetProfileID() string {
	proList, _ := u.ListProfiles()
	projs, _ := json.Marshal(proList)
	var pros []model.ProfileSpec
	var prfID string
	json.Unmarshal([]byte(projs), &pros)
	for _, v := range pros {
		if v.Name == "default" {
			prfID = v.Id
		}
	}
	return prfID
}

//Get volume ID & status
func GetVolumeID(volName string) []string {
	volList, err := u.ListVolumes()
	if err != nil {
		fmt.Println("Can't list volume!")
		return []string{}
	}
	voljs, _ := json.MarshalIndent(volList, "", " ")
	var vols []model.VolumeSpec
	var volID string
	var status string
	var name string
	var desc string
	var size string
	json.Unmarshal([]byte(voljs), &vols)
	for _, v := range vols {
		if v.Name == volName {
			volID = v.Id
			status = v.Status
			name = v.Name
			desc = v.Description
			size = string(v.Size)
		}
	}
	return []string{volID, status, name, desc, size}

}

//Get volume Attachment ID & status
func GetVolAttaID() []string {
	attList, err := u.ListVolumeAttachments()
	if err != nil {
		return []string{}
	}
	attjs, _ := json.Marshal(attList)
	volID := GetVolumeID(UPDATENAME)[0]
	var atts []model.VolumeAttachmentSpec
	var attID string
	var status string
	json.Unmarshal([]byte(attjs), &atts)
	for _, v := range atts {
		if v.VolumeId == volID {
			attID = v.Id
			status = v.Status
		}
	}

	return []string{attID, status}
}

//Get volume snapInfo

func GetVolSnapInfo() []string {
	snapList, err := u.ListVolumeSnapshots()
	if err != nil {
		return []string{}
	}
	snapjs, _ := json.MarshalIndent(snapList, "", " ")

	var snaps []model.VolumeSnapshotSpec
	var snapID string
	var snapName string
	var snapDesc string
	json.Unmarshal([]byte(snapjs), &snaps)
	volID := GetVolumeID(UPDATENAME)[0]
	for _, v := range snaps {
		if v.VolumeId == volID {
			snapID = v.Id
			snapName = v.Name
			snapDesc = v.Description
			return []string{snapID, snapName, snapDesc}
		}
	}
	return []string{}
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

//clear log method
func CleanLog() {
	cmd := exec.Command("/bin/bash", "./clscanlog.sh")
	_, err := cmd.Output()
	if err != nil {
		fmt.Println("cmd.Output: ", err)
		return
	}
}

//Check if Disk Log contain /dev/sd&& 2 GiB
func DiskChk(log string, str string) bool {
	dLog, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer dLog.Close()
	disk := bufio.NewReader(dLog)
	count := 0
	for {
		line, _ := disk.ReadString('\n')
		if strings.Index(line, str) != -1 {
			return true
		} else {
			count++
			if count > 25 {
				break
			}
			continue
		}
		if err != nil || io.EOF == err {
			return false
		}
	}
	return false
}

var u = client.NewClient(&client.Config{
	Endpoint:    "http://localhost:50040",
	AuthOptions: client.NewNoauthOptions(constants.DefaultTenantId)})

func TestCreateProfile(t *testing.T) {
	CleanLog()
	var body = &model.ProfileSpec{
		Name:        "default",
		Description: "default policy",
	}
	cprf, err := u.CreateProfile(body)
	if err != nil {
		t.Error("create profile in client failed:", err)
		return
	}
	prfBody, _ := json.MarshalIndent(cprf, "", "    ")
	t.Log(string(prfBody))
	t.Log("Create Profile Success")
}

func TestGetProfileList(t *testing.T) {
	listRs, err := u.ListProfiles()
	if err != nil {
		t.Error("list profiles in client failed:", err)
		return
	}
	prflist, _ := json.MarshalIndent(listRs, "", "    ")
	t.Log(string(prflist))
	t.Log("Get Profile List Success")
}

func TestGetProfileDetail(t *testing.T) {
	//get create profile id
	prfID := GetProfileID()
	//Get ProfileDetail
	detail, err := u.GetProfile(prfID)
	if err != nil {
		t.Error("get profile in client failed:", err)
		return
	}
	prfdet, _ := json.MarshalIndent(detail, "", "    ")
	t.Log(string(prfdet))
	t.Log("Get Profile Detail Success")
}

func TestCreateVolumeFlow(t *testing.T) {
	var volbody = &model.VolumeSpec{
		Name:        VOLNAME,
		Description: VOLDESC,
		Size:        int64(1),
	}
	create, err := u.CreateVolume(volbody)
	if err != nil {
		t.Log("Create Volume Fail")
	}
	volrs, _ := json.MarshalIndent(create, "", "    ")
	t.Log(string(volrs))
	t.Log("Create Volume Success")
}

//update volume
func TestUpdateVolumeFlow(t *testing.T) {
	volID := GetVolumeID(VOLNAME)[0]
	body := &model.VolumeSpec{
		Name:        UPDATENAME,
		Description: UPDATEDESC,
	}
	upvol, err := u.UpdateVolume(volID, body)
	volary := GetVolumeID(UPDATENAME)
	name := volary[2]
	desc := volary[3]
	if err != nil || name != UPDATENAME || desc != UPDATEDESC {
		t.Error("update volume fail!", err)
		return
	}
	volrs, _ := json.MarshalIndent(upvol, "", "    ")
	t.Log(string(volrs))
	t.Log("Update Volume Success")
}

//extend volume
func TestExtendVolumeFlow(t *testing.T) {
	volID := GetVolumeID(UPDATENAME)[0]

	body := &model.ExtendVolumeSpec{
		Extend: model.ExtendSpec{NewSize: EXTENDSIZE},
	}

	ext, err := u.ExtendVolume(volID, body)
	volext := GetVolumeID(UPDATENAME)
	size := volext[4]
	if err != nil || size != string(EXTENDSIZE) {
		t.Error("Extend volume fail", err)
	}
	extrs, _ := json.MarshalIndent(ext, "", "    ")
	t.Log(string(extrs))
	t.Log("Extend Volume Success")
}

//get volume detail
func TestGetVolumeFlow(t *testing.T) {
	time.Sleep(5 * 1e9)
	volary := GetVolumeID(UPDATENAME)
	volID := volary[0]
	status := volary[1]

	volinfo, err := u.GetVolume(volID)
	if err != nil || status != "available" {
		t.Error("get volume in client failed:", err)
		return
	}
	detail, _ := json.MarshalIndent(volinfo, "", "    ")
	t.Log(string(detail))
	t.Log("Get volume detail Success")
}

//Create Volume Snapshot
func TestCreateSnapshot(t *testing.T) {
	volID := GetVolumeID(UPDATENAME)[0]
	var body = &model.VolumeSnapshotSpec{
		Name:        SNAPNAME,
		Description: SNAPDESC,
		VolumeId:    volID,
	}
	snap, err := u.CreateVolumeSnapshot(body)
	if err != nil {
		t.Error("create volume snapshot in client failed:", err)
		return
	}
	snapcre, _ := json.MarshalIndent(snap, "", "    ")
	t.Log(string(snapcre))
	t.Log("---Create Volume Snapshot Success----")
}

//Update Volume Snapshot
func TestUpdateSnapshot(t *testing.T) {
	snpID := GetVolSnapInfo()[0]
	body := &model.VolumeSnapshotSpec{
		Name:        UPDATESNAPNAME,
		Description: UPDATESNAPDESC,
	}
	upsnap, err := u.UpdateVolumeSnapshot(snpID, body)
	if err != nil {
		t.Error("update volume snapshot in client failed:", err)
		return
	}
	upsnaprs, _ := json.MarshalIndent(upsnap, "", "    ")
	t.Log(string(upsnaprs))
	t.Log("-----Update Snapshot Success-----")
}

//Get Volume snapshot Detail
func TestGetSnapshotDetail(t *testing.T) {
	snpID := GetVolSnapInfo()[0]
	snapdet, err := u.GetVolumeSnapshot(snpID)
	if err != nil || GetVolSnapInfo()[1] != UPDATESNAPNAME || GetVolSnapInfo()[2] != UPDATESNAPDESC {
		t.Error("get volume snapshot in client failed:", err)
		return
	}
	detail, _ := json.MarshalIndent(snapdet, "", "    ")
	t.Log(string(detail))
	t.Log("---Get Snapshot Detail Success----")
}

//Delete Volume Snapshot
func TestDelSnapshot(t *testing.T) {
	snapId := GetVolSnapInfo()[0]
	if err := u.DeleteVolumeSnapshot(snapId, nil); err != nil {
		t.Error("Delete Snapshot Fail")
		return
	}
	t.Log("Delete Snapshot Success")
}

//Create Vloume Attachement
func TestCreateVolAttch(t *testing.T) {
	volID := GetVolumeID(UPDATENAME)[0]
	var body = &model.VolumeAttachmentSpec{
		VolumeId: volID,
		HostInfo: model.HostInfo{},
	}
	attc, err := u.CreateVolumeAttachment(body)
	if err != nil {
		t.Error("create volume attachment in client failed:", err)
		return
	}
	attrs, _ := json.MarshalIndent(attc, "", "    ")
	t.Log(string(attrs))
	t.Log("Create Volume Attachment Success")
}

//show attachment detail
func TestShowAttacDetail(t *testing.T) {
	attID := GetVolAttaID()[0]
	attsta := GetVolAttaID()[1]
	getatt, err := u.GetVolumeAttachment(attID)
	//scan volume
	out := ScanVolume()
	//read Dsik.log
	dev := DiskChk(out, "/dev/sd")
	ca := DiskChk(out, "2 GiB")
	if err != nil || attsta != "available" || dev == false || ca == false {
		t.Log("Volume attachment detail check fail", err)
		return
	}
	detail, _ := json.MarshalIndent(getatt, "", "    ")
	t.Log(string(detail))
	t.Log("Volume attach detail Check Success!")
}

//delete volume attachment
func TestDeleteVolAttach(t *testing.T) {
	attID := GetVolAttaID()[0]
	err := u.DeleteVolumeAttachment(attID, nil)
	//check if attachment is exist after
	for i := 0; i < 5; i++ {
		chk, _ := u.GetVolumeAttachment(attID)
		if chk == nil {
			break
		} else {
			time.Sleep(1 * 1e9)
		}
	}
	//Chk volume scan
	ScanVolume()
	b := DiskChk(DISKLOG, "/dev/sdb")
	chk2, _ := u.GetVolumeAttachment(attID)
	if err != nil || chk2 != nil || b == true {
		t.Log("Delete Attachment Fail", err)
	}
	t.Log("Delete attachment Success")
}

//delete volume
func TestClientDeleteVolume(t *testing.T) {
	volID := GetVolumeID(UPDATENAME)[0]
	body := &model.VolumeSpec{}
	err := u.DeleteVolume(volID, body)
	//check if attachment is exist
	for i := 0; i < 10; i++ {
		chk, _ := u.GetVolume(volID)
		if chk == nil {
			break
		} else {
			time.Sleep(1 * 1e9)
		}
	}
	chk2, _ := u.GetVolume(volID)
	if err != nil || chk2 != nil {
		t.Error("delete volume in client failed:", err)
	}

	t.Log("Delete volume success!")
}

//delete profile
func TestDeleteProfile(t *testing.T) {
	prfID := GetProfileID()
	err := u.DeleteProfile(prfID)
	if err != nil {
		t.Error("delete profile in client failed:", err)
		return
	}
	t.Log("Delete Profile Success")
}
