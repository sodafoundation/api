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

// +build e2ef

package e2e

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"testing"

	"github.com/opensds/opensds/client"
	"github.com/opensds/opensds/pkg/model"
	"github.com/opensds/opensds/pkg/utils/constants"
)

const (
	nvmepool       = "opensds-volumes-nvme"
	defaultgroup   = "opensds-volumes-default"
	localIqn       = "iqn.2017-10.io.opensds:volume:00000001"
	localNqn       = "nqn.ini.1a2bc3d4-27c5-448f-ac84-5bf7fc154321"
	nvmeofProtocol = "nvmeof"
	iscsiProtocol  = "iscsi"
)

var u *SyncClient

// init Create Profile
func init() {
	u, _ = NewSyncClient(&client.Config{
		Endpoint:    constants.DefaultOpensdsEndpoint,
		AuthOptions: client.NewNoauthOptions(constants.DefaultTenantId),
	})

	prfs, _ := u.ListProfiles()
	if len(prfs) == 0 {
		fmt.Println("Start creating profile...")
		var body = &model.ProfileSpec{
			Name:        "default",
			Description: "default policy",
			StorageType: "block",
		}
		prf, err := u.CreateProfile(body)
		if err != nil {
			fmt.Printf("create profile failed: %v\n", err)
			return
		}
		prfs = append(prfs, prf)
	}
}

// Test Case: volume operation flow
func TestVolumeOperationFlow(t *testing.T) {
	t.Log("Begin to create volume...")
	vol, err := u.CreateVolume(&model.VolumeSpec{
		Name:        "flowVolume",
		Description: "This a volume for test",
		Size:        int64(1),
	})
	if err != nil {
		t.Error("create volume failed:", err)
		return
	}
	defer cleanVolumeForTest(t, vol.Id)
	// Check if the status of created volume is available.
	if vol, _ = u.GetVolume(vol.Id); vol.Status != model.VolumeAvailable {
		t.Errorf("status expected %s, got %s\n", model.VolumeAvailable, vol.Status)
		return
	}
	t.Log("Create volume success!")

	t.Log("Start to update volume...")
	vol, err = u.UpdateVolume(vol.Id, &model.VolumeSpec{
		Name:        "Update Volume Name",
		Description: "Update Volume Description",
	})
	if err != nil {
		t.Error("update Volume failed:", err)
		return
	}
	// Check if the status of updated volume is available.
	if vol, _ = u.GetVolume(vol.Id); vol.Status != model.VolumeAvailable {
		t.Errorf("status expected %s, got %s\n", model.VolumeAvailable, vol.Status)
		return
	}
	t.Log("Update volume success!")

	t.Log("Begin to extend volume...")
	vol, err = u.ExtendVolume(vol.Id, &model.ExtendVolumeSpec{
		NewSize: int64(2),
	})
	if err != nil {
		t.Error("extend volume failed:", err)
	}
	// Check if the status of extended volume is available.
	vol, _ = u.GetVolume(vol.Id)
	if vol.Status != model.VolumeAvailable {
		t.Errorf("status expected %s, got %s\n", model.VolumeAvailable, vol.Status)
		return
	}
	// Check if the size of extended volume is 2GB.
	if vol.Size != int64(2) {
		t.Errorf("size expected %d, got %d\n", int64(2), vol.Size)
		return
	}
	t.Log("Extend volume success!")

	t.Log("Begin to delete volume...")
	if err = u.DeleteVolume(vol.Id, nil); err != nil {
		t.Error("delete volume failed:", err)
	}
	t.Log("Delete volume success!")
}

// Test Case: volume snapshot operation flow
func TestVolumeSnapshotOperationFlow(t *testing.T) {
	vol, err := prepareVolumeForTest(t)
	if err != nil {
		t.Error("failed to run volume prepare function:", err)
		return
	}

	t.Log("Begin to create volume snapshot...")
	snp, err := u.CreateVolumeSnapshot(&model.VolumeSnapshotSpec{
		Name:        "new-snapshot",
		Description: "This is a snapshot test for new",
		VolumeId:    vol.Id,
	})
	if err != nil {
		t.Error("create volume snapshot failed:", err)
		return
	}
	defer cleanVolumeAndSnapshotForTest(t, vol.Id, snp.Id)
	// Check if the status of created volume snapshot is available.
	if snp, _ = u.GetVolumeSnapshot(snp.Id); snp.Status != model.VolumeSnapAvailable {
		t.Errorf("status expected is %s, got %s\n", model.VolumeSnapAvailable, snp.Status)
	}
	t.Log("Create volume snapshot success!")

	t.Log("Begin to update volume snapshot...")
	snp, err = u.UpdateVolumeSnapshot(snp.Id, &model.VolumeSnapshotSpec{
		Name:        "Update Volume Snapshot Name",
		Description: "Update Volume Snapshot Description",
	})
	if err != nil {
		t.Error("update volume snapshot failed:", err)
		return
	}
	// Check if the status of updated volume snapshot is available.
	if snp, _ = u.GetVolumeSnapshot(snp.Id); snp.Status != model.VolumeSnapAvailable {
		t.Errorf("status expected is %s, got %s\n", model.VolumeSnapAvailable, snp.Status)
	}
	t.Log("Update volume snapshot success!")

	t.Log("Begin to delete volume snapshot...")
	if err := u.DeleteVolumeSnapshot(snp.Id, nil); err != nil {
		t.Error("delete volume snapshot failed:", err)
		return
	}
	t.Log("Delete volume snapshot success!")
}

// Test Case: volume attachment operation flow
func TestVolumeAttachmentOperationFlow(t *testing.T) {
	vol, err := prepareVolumeForTest(t)
	if err != nil {
		t.Error("failed to run volume prepare function:", err)
		return
	}

	t.Log("Begin to create volume attachment...")
	host, _ := os.Hostname()
	atc, err := u.CreateVolumeAttachment(&model.VolumeAttachmentSpec{
		VolumeId: vol.Id,
		HostInfo: model.HostInfo{
			Host:      host,
			Platform:  runtime.GOARCH,
			OsType:    runtime.GOOS,
			Ip:        "127.0.0.1",
			Initiator: localIqn,
		},
		AccessProtocol: iscsiProtocol,
	})
	if err != nil {
		t.Error("create volume attachment failed:", err)
		return
	}
	defer cleanVolumeAndAttachmentForTest(t, vol.Id, atc.Id)
	// Check if the status of created volume attachment is available.
	if atc, _ = u.GetVolumeAttachment(atc.Id); atc.Status != model.VolumeAttachAvailable {
		t.Errorf("status expected %s, got %s\n", model.VolumeAttachAvailable, atc.Status)
		return
	}
	t.Log("Create volume attachment success!")

	t.Log("Begin to delete volume attachment...")
	if err := u.DeleteVolumeAttachment(atc.Id, nil); err != nil {
		t.Error("delete volume attachment failed:", err)
		return
	}
	t.Log("Delete volume attachment success!")
}

// Test Case: volume attach/detach operation flow
func TestVolumeAttachOperationFlow(t *testing.T) {
	atc, err := prepareVolumeAttachmentForTest(t)
	if err != nil {
		t.Error("prepare volume attachment failed:", err)
		return
	}
	defer cleanVolumeAndAttachmentForTest(t, atc.VolumeId, atc.Id)

	t.Log("Begin to attach volume...")
	t.Log("atc.ConnectionData is:", atc.ConnectionData)
	// run volume-connector tool to execute attach operation
	conn, err := json.Marshal(&atc.ConnectionData)
	if err != nil {
		t.Error("failed to marshal connection data:", err)
		return
	}
	accPro := atc.AccessProtocol
	output, err := execCmd("sudo", "./volume-connector",
		"attach", string(conn), accPro)
	if err != nil {
		t.Error("failed to attach volume:", output, err)
		return
	}
	t.Log(output)
	t.Log("Volume attach success!")

	out, _ := execCmd("/bin/bash", "-c", "iscsiadm -m session")
	t.Log("Session is:", out)

	t.Log("Begin to detach volume...")
	// run volume-connector tool to execute detach operation
	output, err = execCmd("sudo", "./volume-connector",
		"detach", string(conn), accPro)
	if err != nil {
		t.Error("failed to detach volume:", output, err)
		return
	}
	t.Log(output)
	t.Log("Volume detach Success!")
}

// Test for nvmeof connection
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

	vol, err := PrepareNvmeVolume()
	if err != nil {
		t.Error("prepare nvme volume failed:", err)
		return
	}

	host, _ := os.Hostname()
	atc, err := u.CreateVolumeAttachment(&model.VolumeAttachmentSpec{
		VolumeId: vol.Id,
		HostInfo: model.HostInfo{
			Host:      host,
			Platform:  runtime.GOARCH,
			OsType:    runtime.GOOS,
			Ip:        "127.0.0.1",
			Initiator: localNqn,
		},
		AccessProtocol: nvmeofProtocol,
	})
	if err != nil {
		t.Error("create volume attachment failed:", err)
		return
	}
	defer cleanVolumeAndAttachmentForTest(t, vol.Id, atc.Id)
	// Check if the status of created volume attachment is available.
	if atc, _ = u.GetVolumeAttachment(atc.Id); atc.Status != model.VolumeAttachAvailable {
		t.Errorf("status expected %s, got %s\n", model.VolumeAttachAvailable, atc.Status)
		return
	}
	t.Log("Create volume attachment success!")

	t.Log("Begin to delete volume attachment...")
	if err := u.DeleteVolumeAttachment(atc.Id, nil); err != nil {
		t.Error("delete volume attachment failed:", err)
		return
	}
	t.Log("Delete volume attachment success!")

	err = NvmeofVolumeAttachHost(t)
	if err != nil {
		t.Error("connect nvmeof attachment fail", err)
		return
	}

	t.Log("nvmeof attach issues success")
}

func CreateNvmeofAttach(t *testing.T) error {
	vol, err := PrepareNvmeVolume()
	if err != nil {
		t.Error("prepare nvme volume failed:", err)
		return err
	}
	attc, err := u.CreateVolumeAttachment(&model.VolumeAttachmentSpec{
		VolumeId: vol.Id,
		HostInfo: model.HostInfo{},
	})
	if err != nil {
		t.Error("create nvmeof volume attachment failed:", err)
		return err
	}
	defer cleanVolumeAndAttachmentForTest(t, vol.Id, attc.Id)

	attrs, _ := json.MarshalIndent(attc, "", "    ")
	t.Log(string(attrs))
	t.Log("Create nvmeof Volume Attachment Success!")
	return nil
}

func ListNvmeofAttachment(t *testing.T) error {
	attc, err := PrepareNvmeofAttachment(t)
	if err != nil {
		t.Error("prepare nvmeof attachment failed:", err)
		return err
	}
	defer cleanVolumeAndAttachmentForTest(t, attc.VolumeId, attc.Id)

	atts, err := u.ListVolumeAttachments()
	if err != nil {
		t.Error("list nvmeof attachment failed:", err)
		return err
	}
	attli, _ := json.MarshalIndent(atts, "", "    ")
	t.Log(string(attli))
	t.Log("List nvmeof attachments success!")
	return nil
}

func ShowNvmeofAttachDetail(t *testing.T) error {
	attc, err := PrepareNvmeofAttachment(t)
	if err != nil {
		t.Error("prepare attachment failed:", err)
		return err
	}
	defer cleanVolumeAndAttachmentForTest(t, attc.VolumeId, attc.Id)

	getatt, err := u.GetVolumeAttachment(attc.Id)
	if err != nil || getatt.Status != "available" {
		t.Error("get volume attachment detail failed:", err)
		return err
	}
	t.Log("Get Volume Attachment Detail Success")
	return nil
}

func DeleteNvmeofAttach(t *testing.T) error {
	attc, err := PrepareNvmeofAttachment(t)
	if err != nil {
		t.Error("prepare attachment failed:", err)
		return err
	}
	defer cleanVolumeForTest(t, attc.VolumeId)

	err = u.DeleteVolumeAttachment(attc.Id, nil)
	if err != nil {
		t.Error("delete nvmeof attachment failed:", err)
		return err
	}
	_, err = u.GetVolumeAttachment(attc.Id)
	t.Log("err:", err)
	if strings.Contains(err.Error(), "can't find") {
		t.Log("Delete attachment Success")
		return nil
	} else {
		t.Error("delete attachment failed:", err)
		return err
	}
}

//Test Case:Nvmeof Volume Attach to specific host
func NvmeofVolumeAttachHost(t *testing.T) error {
	attc, err := PrepareNvmeofAttachmentHost(t)
	if err != nil {
		t.Error("Prepare Attachment Fail:", err)
		return err
	}
	defer cleanVolumeAndAttachmentForTest(t, attc.VolumeId, attc.Id)
	getatt, err := u.GetVolumeAttachment(attc.Id)
	if err != nil || getatt.Status != "available" {
		t.Errorf("attachment(%s) is not available: %v", attc.Id, err)
		return err
	}

	t.Log("Begin to Scan Volume:")
	t.Log("getatt.AccessProtocol", getatt.AccessProtocol)
	t.Log("getatt.Metadata", getatt.ConnectionData)

	output, _ := execCmd("/bin/bash", "-c", "ps -ef")
	t.Log(output)
	//execute bin file
	conn, err := json.Marshal(&getatt.ConnectionData)
	if err != nil {
		t.Error("Failed to marshal connection data:", err)
		return err
	}
	accPro := getatt.AccessProtocol
	output, err = execCmd("sudo", "./volume-connector",
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

//Test Case:Nvmeof Volume Attach
func NvmeofVolumeAttach(t *testing.T) error {
	attc, err := PrepareNvmeofAttachment(t)
	if err != nil {
		t.Error("prepare attachment failed:", err)
		return err
	}
	defer cleanVolumeAndAttachmentForTest(t, attc.VolumeId, attc.Id)

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
	defer cleanVolumeAndAttachmentForTest(t, attc.VolumeId, attc.Id)

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

// execCmd operation
func execCmd(name string, arg ...string) (string, error) {
	fmt.Printf("Command: %s %s:\n", name, strings.Join(arg, " "))
	info, err := exec.Command(name, arg...).CombinedOutput()
	return string(info), err
}

// prepare volume for test
func prepareVolumeForTest(t *testing.T) (*model.VolumeSpec, error) {
	t.Log("Start preparing volume...")
	// get poolid
	pols, err := u.ListPools()
	if err != nil {
		return nil, err
	}
	polId := ""
	for _, p := range pols {
		if p.Name == defaultgroup {
			polId = p.Id
			break
		}
	}
	if polId == "" {
		return nil, nil
	}

	// create volume in default pool
	vol, err := u.CreateVolume(&model.VolumeSpec{
		Name:        "test",
		Description: "This is a test",
		Size:        int64(1),
		PoolId:      polId,
	})
	if err != nil {
		t.Error("prepare volume failed:", err)
		return nil, err
	}
	if vol, _ = u.GetVolume(vol.Id); vol.Status != model.VolumeAvailable {
		return nil, fmt.Errorf("the status of volume is not available!")
	}

	t.Log("End preparing volume...")
	return vol, nil
}

// prepare volume attachment for test
func prepareVolumeAttachmentForTest(t *testing.T) (*model.VolumeAttachmentSpec, error) {
	vol, err := prepareVolumeForTest(t)
	if err != nil {
		t.Error("failed to run volume prepare function:", err)
		return nil, err
	}

	t.Log("Start preparing volume attachment...")
	atc, err := u.CreateVolumeAttachment(&model.VolumeAttachmentSpec{
		VolumeId:       vol.Id,
		HostInfo:       model.HostInfo{},
		AccessProtocol: iscsiProtocol,
	})
	if err != nil {
		t.Error("prepare volume attachment failed:", err)
		// Run volume clean function if failed to prepare volume attachment.
		cleanVolumeForTest(t, atc.VolumeId)
		return nil, err
	}
	if atc, _ = u.GetVolumeAttachment(atc.Id); atc.Status != model.VolumeAttachAvailable {
		// Run volume clean function if failed to prepare volume attachment.
		cleanVolumeForTest(t, atc.VolumeId)
		return nil, fmt.Errorf("the status of volume attachment is not available!")
	}

	t.Log("End preparing volume attachment...")
	return atc, nil
}

// nvme volume is essential for nvmeof attachment, so the volume should be
// created in nvme pool
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

	// create volume in specified nvme pool
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

// prepare nvmeof attachment availible to specific host
func PrepareNvmeofAttachmentHost(t *testing.T) (*model.VolumeAttachmentSpec, error) {
	vol, err := PrepareNvmeVolume()
	if err != nil {
		t.Error("Prepare nvmeof  Volume Fail", err)
		return nil, err
	}

	var body = &model.VolumeAttachmentSpec{
		VolumeId: vol.Id,
		HostInfo: model.HostInfo{
			Initiator: "nqn.ini.1A2B3C4D5E6F7G8H",
		},
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

func cleanVolumeForTest(t *testing.T, volID string) {
	t.Log("Start cleaning volume...")
	u.DeleteVolume(volID, nil)
	t.Log("End cleaning volume...")
}

func cleanVolumeAndAttachmentForTest(t *testing.T, volID, atcID string) {
	t.Log("Start cleaning volume attachment...")
	u.DeleteVolumeAttachment(atcID, nil)
	t.Log("End cleaning volume attachment...")

	t.Log("Start cleaning volume...")
	u.DeleteVolume(volID, nil)
	t.Log("End cleaning volume...")
}

func cleanVolumeAndSnapshotForTest(t *testing.T, volID, snpID string) {
	t.Log("Start cleaing volume snapshot...")
	u.DeleteVolumeSnapshot(snpID, nil)
	t.Log("End cleaning volume snapshot...")

	t.Log("Start cleaning volume...")
	u.DeleteVolume(volID, nil)
	t.Log("End cleaning volume...")
}
