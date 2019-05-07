// Copyright (c) 2018 Huawei Technologies Co., Ltd. All Rights Reserved.
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

const (
	nvmepool       = "opensds-volumes-nvme"
	defaultgroup   = "opensds-volumes-default"
	localIqn       = "iqn.2017-10.io.opensds:volume:00000001"
	nvmeofProtocol = "nvmeof"
	iscsiProtocol  = "iscsi"
)

var u *client.Client

//init Create Profile
func init() {
	u, _ = client.NewClient(&client.Config{
		Endpoint:    "http://localhost:50040",
		AuthOptions: client.NewNoauthOptions(constants.DefaultTenantId)})

	var body = &model.ProfileSpec{
		Name:        "default",
		Description: "default policy",
		StorageType: "block",
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
	err = u.DeleteProfile(pro.Id)
	if err != nil {
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
	_, err := u.ExtendVolume(vol.Id, body)
	if err != nil {
		t.Error("Extend volume fail", err)
	}
	defer DeleteVolume(vol.Id)
	time.Sleep(3 * 1e9)
	info, _ := u.GetVolume(vol.Id)
	t.Log("SIZE:", info.Size)
	if info.Size != 2 {
		t.Error("Extend Volume Size is wrong!")
		return
	}
	t.Log("Extend Size is Right!")
	extrs, _ := json.MarshalIndent(info, "", "    ")
	t.Log(string(extrs))

	t.Log("Creates a volume snapshot...")
	snpBody := &model.VolumeSnapshotSpec{
		Name:        "new-snapshot",
		Description: "This is a snapshot test for new",
		VolumeId:    vol.Id,
	}

	snp, err := u.CreateVolumeSnapshot(snpBody)
	if err != nil {
		t.Error("prepare volume snapshot failed:", err)
		return
	}

	defer DeleteSnapshot(snp.Id)

	body = &model.ExtendVolumeSpec{
		NewSize: int64(3),
	}
	time.Sleep(3 * 1e9)
	_, err = u.ExtendVolume(vol.Id, body)
	if err != nil {
		t.Error("Extend volume fail", err)
	}

	time.Sleep(3 * 1e9)
	info, _ = u.GetVolume(vol.Id)
	t.Log("SIZE:", info.Size)
	if info.Size != 3 {
		t.Error("Extend Volume Size is wrong!")
		return
	}

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
	vols, err := u.ListVolumes()
	if err != nil {
		t.Error("List Volume Fail", err)
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
	Info, err := u.GetVolume(vol.Id)
	if err != nil {
		t.Error("Get Volume Detail Fail", err)
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
	err = u.DeleteVolume(vol.Id, nil)
	if err != nil {
		t.Error("Delete Volume fail", err)
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
	snp, err := u.CreateVolumeSnapshot(body)
	if err != nil {
		t.Error("prepare volume snapshot failed:", err)
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
	upsnap, err := u.UpdateVolumeSnapshot(snap.Id, body)
	if err != nil {
		t.Error("update volume snapshot failed:", err)
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

	snpli, err := u.ListVolumeSnapshots()
	if err != nil {
		t.Error("List Snapshot Fail", err)
		return
	}
	snapli, _ := json.MarshalIndent(snpli, "", "    ")
	t.Log(string(snapli))
	t.Log("-----List Snapshot Success-----")
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

	snapinfo, err := u.GetVolumeSnapshot(snap.Id)
	if err != nil {
		t.Error("Get Snapshot Detail Fail!", err)
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
	err = u.DeleteVolumeSnapshot(snap.Id, nil)
	if err != nil {
		t.Error("Delete Snapshot Fail!", err)
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
		VolumeId:       vol.Id,
		HostInfo:       model.HostInfo{},
		AccessProtocol: iscsiProtocol,
	}
	attc, err := u.CreateVolumeAttachment(body)
	if err != nil {
		t.Error("create volume attachment failed:", err)
		return
	}
	defer DeleteAttachment(attc.Id)
	attrs, _ := json.MarshalIndent(attc, "", "    ")
	t.Log(string(attrs))
	t.Log("Create Volume Attachment Success")
}

//Test Case:List Attachment
func TestListAttachment(t *testing.T) {
	attc, err := PrepareAttachment(t)
	if err != nil {
		t.Error("Prepare Attachment Fail!", err)
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
	attc, err := PrepareAttachment(t)
	if err != nil {
		t.Error("Prepare Attachment Fail!", err)
		return
	}
	defer DeleteVolume(attc.VolumeId)
	defer DeleteAttachment(attc.Id)

	getatt, err := u.GetVolumeAttachment(attc.Id)
	if err != nil || getatt.Status != "available" {
		t.Error("Get Volume Attachment Detail Fail!", err)
		return
	}
	t.Log("Get Volume Attachment Detail Success")
}

//Test Case:Volume Attach
func TestVolumeAttach(t *testing.T) {
	attc, err := PrepareAttachment(t)
	if err != nil {
		t.Error("Prepare Attachment Fail:", err)
		return
	}
	defer DeleteVolume(attc.VolumeId)
	defer DeleteAttachment(attc.Id)

	getatt, err := u.GetVolumeAttachment(attc.Id)
	if err != nil || getatt.Status != "available" {
		t.Errorf("attachment(%s) is not available: %v", attc.Id, err)
		return
	}

	t.Log("Begin to Scan Volume:")
	t.Log("getatt.Metadata", getatt.ConnectionData)

	//execute bin file
	conn, err := json.Marshal(&getatt.ConnectionData)
	if err != nil {
		t.Error("Failed to marshal connection data:", err)
		return
	}
	accPro := getatt.AccessProtocol
	output, err := execCmd("sudo", "./volume-connector",
		"attach", string(conn), accPro)
	if err != nil {
		t.Error("Failed to attach volume:", output, err)
		return
	}
	t.Log(output)
	t.Log("Volume attach success!")
}

//Test Case:Delete Attachment
func TestDeleteAttach(t *testing.T) {
	attc, err := PrepareAttachment(t)
	if err != nil {
		t.Error("Prepare Attachment Fail!", err)
		return
	}
	defer DeleteVolume(attc.VolumeId)
	err = u.DeleteVolumeAttachment(attc.Id, nil)
	if err != nil {
		t.Error("Delete Attachment Fail", err)
		return
	}
	_, err = u.GetVolumeAttachment(attc.Id)
	t.Log("err:", err)
	if strings.Contains(err.Error(), "can't find") {
		t.Log("Delete attachment Success")
		return
	} else {
		t.Error("Delete Attachment Fail!", err)
	}
}

//Test Case:Delete Attachment
func TestVolumeDetach(t *testing.T) {
	attc, err := PrepareAttachment(t)
	if err != nil {
		t.Error("Prepare Attachment Fail!", err)
		return
	}

	out, _ := execCmd("/bin/bash", "-c", "iscsiadm -m session")
	fmt.Println("session is ", out)
	defer DeleteVolume(attc.VolumeId)
	defer DeleteAttachment(attc.Id)

	getatt, err := u.GetVolumeAttachment(attc.Id)
	if err != nil || getatt.Status != "available" {
		t.Errorf("attachment(%s) is not available: %v", attc.Id, err)
		return
	}

	t.Log("Begin to Scan volume:")
	t.Log("getatt.Accessprotocol", getatt.AccessProtocol)
	t.Log("getatt.Metadata", getatt.ConnectionData)

	//execute bin file
	conn, err := json.Marshal(&getatt.ConnectionData)
	if err != nil {
		t.Error("Failed to marshal connection data:", err)
		return
	}

	// attach first, then detach
	accPro := getatt.AccessProtocol
	output, err := execCmd("sudo", "./volume-connector",
		"attach", string(conn), accPro)
	if err != nil {
		t.Error("Failed to attach volume:", output, err)
		return
	}

	t.Log(output)

	output, err = execCmd("sudo", "./volume-connector",
		"detach", string(conn), accPro)
	if err != nil {
		t.Error("Failed to detach volume:", output, err)
		return
	}
	t.Log(output)
	t.Log("Volume Detach Success!")
}

//Test for nvmeof  connection
func TestNvmeofAttachIssues(t *testing.T) {
	// pool list get nvme pool
	pols, err := u.ListPools()
	if err != nil {
		t.Error("list pools failed:", err)
		return
	}
	polId := ""
	for _, p := range pols {
		if p.Name == nvmepool {
			polId = p.Id
			t.Log("nvme pool id is: ", polId)
			break
		}
	}
	if polId == "" {
		t.Log("no nvme pool ")
		return
	}
	//PrepareNvmeVolume()
	err = CreateNvmeofAttach(t)
	if err != nil {
		t.Error("create nvmeof attachment fail", err)
		return
	}
	err = ListNvmeofAttachment(t)
	if err != nil {
		t.Error("list nvmeof attachment fail", err)
		return
	}
	err = ShowNvmeofAttachDetail(t)
	if err != nil {
		t.Error("show nvmeof attachment fail", err)
		return
	}
	err = NvmeofVolumeAttach(t)
	if err != nil {
		t.Error("connect nvmeof attachment fail", err)
		return
	}

	err = DeleteNvmeofAttach(t)
	if err != nil {
		t.Error("delete nvmeof attachment fail", err)
		return
	}

	t.Log("nvmeof attach issues success")
}

func CreateNvmeofAttach(t *testing.T) error {
	vol, err := PrepareNvmeVolume()
	if err != nil {
		t.Error("Prepare nvme Volume  Fail", err)
		return err
	}
	defer DeleteVolume(vol.Id)
	var body = &model.VolumeAttachmentSpec{
		VolumeId: vol.Id,
		HostInfo: model.HostInfo{},
	}
	attc, err := u.CreateVolumeAttachment(body)
	if err != nil {
		t.Error("create nvmeof volume attachment failed:", err)
		return err
	}
	defer DeleteAttachment(attc.Id)
	attrs, _ := json.MarshalIndent(attc, "", "    ")
	t.Log(string(attrs))
	t.Log("Create nvmeof Volume Attachment Success")
	return nil
}

func ListNvmeofAttachment(t *testing.T) error {
	attc, err := PrepareNvmeofAttachment(t)
	if err != nil {
		t.Error("Prepare nvmeof Attachment Fail!", err)
		return err
	}
	defer DeleteVolume(attc.VolumeId)
	defer DeleteAttachment(attc.Id)
	atts, err := u.ListVolumeAttachments()
	if err != nil {
		t.Error("List nvmeof Attachment Error!", err)
		return err
	}
	attli, _ := json.MarshalIndent(atts, "", "    ")
	t.Log(string(attli))
	t.Log("List nvmeof Attachment Success!")
	return nil
}

func ShowNvmeofAttachDetail(t *testing.T) error {
	attc, err := PrepareNvmeofAttachment(t)
	if err != nil {
		t.Error("Prepare Attachment Fail!", err)
		return err
	}
	defer DeleteVolume(attc.VolumeId)
	defer DeleteAttachment(attc.Id)

	getatt, err := u.GetVolumeAttachment(attc.Id)
	if err != nil || getatt.Status != "available" {
		t.Error("Get Volume Attachment Detail Fail!", err)
		return err
	}
	t.Log("Get Volume Attachment Detail Success")
	return nil
}

func DeleteNvmeofAttach(t *testing.T) error {
	attc, err := PrepareNvmeofAttachment(t)
	if err != nil {
		t.Error("Prepare Attachment Fail!", err)
		return err
	}
	defer DeleteVolume(attc.VolumeId)
	err = u.DeleteVolumeAttachment(attc.Id, nil)
	if err != nil {
		t.Error("Delete nvme Attachment Fail", err)
		return err
	}
	_, err = u.GetVolumeAttachment(attc.Id)
	t.Log("err:", err)
	if strings.Contains(err.Error(), "can't find") {
		t.Log("Delete attachment Success")
		return nil
	} else {
		t.Error("Delete Attachment Fail!", err)
		return err
	}
}

//Test Case:Nvmeof Volume Attach
func NvmeofVolumeAttach(t *testing.T) error {
	attc, err := PrepareNvmeofAttachment(t)
	if err != nil {
		t.Error("Prepare Attachment Fail:", err)
		return err
	}
	defer DeleteVolume(attc.VolumeId)
	defer DeleteAttachment(attc.Id)

	getatt, err := u.GetVolumeAttachment(attc.Id)
	if err != nil || getatt.Status != "available" {
		t.Errorf("attachment(%s) is not available: %v", attc.Id, err)
		return err
	}

	t.Log("Begin to Scan Volume:")
	t.Log("getatt.AccessProtocol", getatt.AccessProtocol)
	t.Log("getatt.Metadata", getatt.ConnectionData)

	//execute bin file
	conn, err := json.Marshal(&getatt.ConnectionData)
	if err != nil {
		t.Error("Failed to marshal connection data:", err)
		return err
	}
	accPro := getatt.AccessProtocol
	output, err := execCmd("sudo", "./volume-connector",
		"attach", string(conn), accPro)
	if err != nil {
		t.Error("Failed to attach volume:", output, err)
		return err
	}
	t.Log(output)
	t.Log("Nvmeof Volume attach yoyo success!")
	// detach it
	err = NvmeofVolumeDetach(t, attc)
	if err != nil {
		t.Error("detach failed")
		return err
	}
	return nil
}

//Test Case:Delete Attachment
func NvmeofVolumeDetach(t *testing.T, attc *model.VolumeAttachmentSpec) error {
	defer DeleteVolume(attc.VolumeId)
	defer DeleteAttachment(attc.Id)

	getatt, err := u.GetVolumeAttachment(attc.Id)
	if err != nil || getatt.Status != "available" {
		t.Errorf("attachment(%s) is not available: %v", attc.Id, err)
		return err
	}

	t.Log("Begin to Scan volume:")
	t.Log("getatt.AccessProtocol", getatt.AccessProtocol)
	t.Log("getatt.Metadata", getatt.ConnectionData)

	//execute bin file
	conn, err := json.Marshal(&getatt.ConnectionData)
	if err != nil {
		t.Error("Failed to marshal connection data:", err)
		return err
	}
	accPro := getatt.AccessProtocol
	output, err := execCmd("sudo", "./volume-connector",
		"detach", string(conn), accPro)
	if err != nil {
		t.Error("Failed to detach volume:", attc.VolumeId, output, err)
		return err
	}
	t.Log(output)
	t.Log("Volume Detach Success!")
	return nil
}

// cmd
func execCmd(name string, arg ...string) (string, error) {
	fmt.Printf("Command: %s %s:\n", name, strings.Join(arg, " "))
	info, err := exec.Command(name, arg...).CombinedOutput()
	return string(info), err
}

//prepare attachment
func PrepareAttachment(t *testing.T) (*model.VolumeAttachmentSpec, error) {
	vol, err := PrepareVolume()
	if err != nil {
		t.Error("Prepare Volume Fail", err)
		return nil, err
	}

	//host, _ := os.Hostname()

	var body = &model.VolumeAttachmentSpec{
		VolumeId:       vol.Id,
		HostInfo:       model.HostInfo{},
		AccessProtocol: iscsiProtocol,
	}
	attc, err := u.CreateVolumeAttachment(body)
	if err != nil {
		t.Error("prepare volume attachment failed:", err)
		return nil, err
	}

	t.Log("prepare Volume Attachment Success")
	return attc, nil
}

// prepare nvmeof attachment
func PrepareNvmeofAttachment(t *testing.T) (*model.VolumeAttachmentSpec, error) {
	vol, err := PrepareNvmeVolume()
	if err != nil {
		t.Error("Prepare nvmeof  Volume Fail", err)
		return nil, err
	}

	var body = &model.VolumeAttachmentSpec{
		VolumeId:       vol.Id,
		HostInfo:       model.HostInfo{},
		AccessProtocol: nvmeofProtocol,
	}
	attc, err := u.CreateVolumeAttachment(body)
	if err != nil {
		t.Error("prepare volume attachment failed:", err)
		return nil, err
	}

	t.Log("prepare nvmeof Volume Attachment Success")
	return attc, nil
}

//delete attachment
func DeleteAttachment(attId string) error {
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
		return err
	}
	fmt.Println("Delete Attachment Success!")
	return nil
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
	snp, err := u.CreateVolumeSnapshot(body)
	if err != nil {
		fmt.Println("prepare volume snapshot failed:", err)
		//cleanVolumeIfFailedOrFinished(t, snp.VolumeId)
		return nil, err
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

	fmt.Println("Prepare Volume Success")
	return create, nil
}

//nvme volume is essential for nvmeof attachment ,so the volume should be created in nvme pool
func PrepareNvmeVolume() (*model.VolumeSpec, error) {
	// get poolid
	pols, err := u.ListPools()
	if err != nil {
		return nil, err
	}
	polId := ""
	for _, p := range pols {
		if p.Name == nvmepool {
			polId = p.Id
			break
		}
	}
	if polId == "" {
		return nil, nil
	}

	//create volume in specified nvme pool
	var volbody = &model.VolumeSpec{
		Name:        "nvme flowTest",
		Description: "This a test for nvme flow",
		Size:        int64(1),
		PoolId:      polId,
	}
	create, err := u.CreateVolume(volbody)
	if err != nil {
		return nil, err
	}
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
		StorageType: "block",
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
