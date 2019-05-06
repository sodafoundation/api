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
	"fmt"
	"net"
	"os"
	"runtime"
	"testing"

	"github.com/opensds/opensds/client"
	"github.com/opensds/opensds/pkg/model"
	"github.com/opensds/opensds/pkg/utils/constants"
)

var (
	c         *client.Client
	profileId string
)

const (
	iscsiProtocol = "iscsi"
	localIqn      = "iqn.2017-10.io.opensds:volume:00000001"
)

func init() {
	c, _ = client.NewClient(&client.Config{
		Endpoint:    "http://localhost:50040",
		AuthOptions: client.NewNoauthOptions(constants.DefaultTenantId),
	})

	prfs, _ := c.ListProfiles()
	if len(prfs) == 0 {
		fmt.Println("Start creating profile...")
		var body = &model.ProfileSpec{
			Name:        "default",
			Description: "default policy",
		}
		prf, err := c.CreateProfile(body)
		if err != nil {
			fmt.Printf("create profile failed: %v\n", err)
			return
		}
		prfs = append(prfs, prf)
	}
	profileId = prfs[0].Id
}

func TestListDocks(t *testing.T) {
	t.Log("Start listing docks...")
	if _, err := c.ListDocks(); err != nil {
		t.Error("list docks failed:", err)
		return
	}
	t.Log("list docks success!")
}

func TestListPools(t *testing.T) {
	t.Log("Start listing pools...")
	if _, err := c.ListPools(); err != nil {
		t.Error("list pools failed:", err)
		return
	}
	t.Log("list pools success!")
}

func TestCreateVolume(t *testing.T) {
	t.Log("Start creating volume...")
	var body = &model.VolumeSpec{
		Name:        "test",
		Description: "This is a test",
		Size:        int64(1),
	}
	vol, err := c.CreateVolume(body)
	if err != nil {
		t.Error("create volume failed:", err)
		return
	}

	defer cleanVolumeIfFailedOrFinished(t, vol.Id)

	// Check if the status of created volume is available.
	vol, _ = c.GetVolume(vol.Id)
	if vol.Status != model.VolumeAvailable {
		t.Errorf("status expected is %s, got %s\n", model.VolumeAvailable, vol.Status)
		return
	}
	t.Log("Create volume success!")
}

func TestGetVolume(t *testing.T) {
	vol, err := prepareVolume(t)
	if err != nil {
		t.Error("Failed to run volume prepare function:", err)
		return
	}
	defer cleanVolumeIfFailedOrFinished(t, vol.Id)

	t.Log("Start checking volume...")
	vol, err = c.GetVolume(vol.Id)
	if err != nil {
		t.Error("Check volume failed:", err)
		return
	}
	// Test the status of created volume.
	if vol.Status != model.VolumeAvailable {
		t.Error("The status of volume is not available!")
		return
	}
	t.Log("Check volume success!")
}

func TestListVolumes(t *testing.T) {
	vol, err := prepareVolume(t)
	if err != nil {
		t.Error("Failed to run volume prepare function:", err)
		return
	}
	defer cleanVolumeIfFailedOrFinished(t, vol.Id)

	t.Log("Start checking all volumes...")
	if _, err := c.ListVolumes(); err != nil {
		t.Error("Check all volumes failed:", err)
		return
	}
	t.Log("Check all volumes success!")
}

func TestUpdateVolume(t *testing.T) {
	vol, err := prepareVolume(t)
	if err != nil {
		t.Error("Failed to run volume prepare function:", err)
		return
	}
	defer cleanVolumeIfFailedOrFinished(t, vol.Id)

	t.Log("Start updating volume...")
	var body = &model.VolumeSpec{
		Name:        "Update Volume Name",
		Description: "Update Volume Description",
	}
	newVol, err := c.UpdateVolume(vol.Id, body)
	if err != nil {
		t.Error("update volume failed:", err)
		return
	}
	// Check if the status of updated volume is available.
	newVol, _ = c.GetVolume(newVol.Id)
	if newVol.Status != model.VolumeAvailable {
		t.Errorf("status expected is %s, got %s\n", model.VolumeAvailable, newVol.Status)
		return
	}
	t.Log("Update volume success!")
}

func TestExtendVolume(t *testing.T) {
	vol, err := prepareVolume(t)
	if err != nil {
		t.Error("Failed to run volume prepare function:", err)
		return
	}
	defer cleanVolumeIfFailedOrFinished(t, vol.Id)

	t.Log("Start extending volume...")
	body := &model.ExtendVolumeSpec{
		NewSize: int64(vol.Size + 1),
	}

	newVol, err := c.ExtendVolume(vol.Id, body)
	if err != nil {
		t.Error("extend volume failed:", err)
		return
	}
	// Check if the status of extended volume is available.
	newVol, _ = c.GetVolume(newVol.Id)
	if newVol.Status != model.VolumeAvailable {
		t.Errorf("status expected is %s, got %s\n", model.VolumeAvailable, newVol.Status)
		return
	}
	t.Log("Extend volume success!")
}

func TestDeleteVolume(t *testing.T) {
	vol, err := prepareVolume(t)
	if err != nil {
		t.Error("Failed to run volume prepare function:", err)
		return
	}

	t.Log("Start deleting volume...")
	if err := c.DeleteVolume(vol.Id, nil); err != nil {
		t.Error("delete volume failed:", err)
		return
	}
	t.Log("Delete volume success!")
}

func TestCreateVolumeAttachment(t *testing.T) {
	vol, err := prepareVolume(t)
	if err != nil {
		t.Error("Failed to run volume prepare function:", err)
		return
	}
	defer cleanVolumeIfFailedOrFinished(t, vol.Id)

	t.Log("Start creating volume attachment...")
	host, _ := os.Hostname()
	var body = &model.VolumeAttachmentSpec{
		VolumeId: vol.Id,
		HostInfo: model.HostInfo{
			Host:      host,
			Platform:  runtime.GOARCH,
			OsType:    runtime.GOOS,
			Ip:        getHostIp(),
			Initiator: localIqn,
		},
		AccessProtocol: iscsiProtocol,
	}
	atc, err := c.CreateVolumeAttachment(body)
	if err != nil {
		t.Error("create volume attachment failed:", err)
		return
	}
	// Check if the status of created volume attachment is available.
	atc, _ = c.GetVolumeAttachment(atc.Id)
	if atc.Status != model.VolumeAttachAvailable {
		t.Errorf("status expected is %s, got %s\n", model.VolumeAttachAvailable, atc.Status)
		return
	}
	t.Log("create volume attachment success!")

	t.Log("Start cleaning volume attachment...")
	if err := c.DeleteVolumeAttachment(atc.Id, nil); err != nil {
		t.Error("Clean volume attachment failed:", err)
		return
	}
	t.Log("End cleaning volume attachment...")
}

func TestGetVolumeAttachment(t *testing.T) {
	atc, err := prepareVolumeAttachment(t)
	if err != nil {
		t.Error("Failed to run volume attachment prepare function:", err)
		return
	}
	defer cleanVolumeAndAttachmentIfFailedOrFinished(t, atc.VolumeId, atc.Id)

	t.Log("Start checking volume attachment...")
	atc, err = c.GetVolumeAttachment(atc.Id)
	if err != nil {
		t.Error("Check volume attachment failed:", err)
		return
	}
	if atc.Status != model.VolumeAttachAvailable {
		t.Errorf("status expected is %s, got %s\n", model.VolumeAttachAvailable, atc.Status)
		return
	}
	t.Log("Check volume attachment success!")
}

func TestListVolumeAttachments(t *testing.T) {
	atc, err := prepareVolumeAttachment(t)
	if err != nil {
		t.Error("Failed to run volume attachment prepare function:", err)
		return
	}
	defer cleanVolumeAndAttachmentIfFailedOrFinished(t, atc.VolumeId, atc.Id)

	t.Log("Start checking all volume attachments...")
	if _, err := c.ListVolumeAttachments(); err != nil {
		t.Error("Check all volume attachments failed:", err)
		return
	}
	t.Log("list volume attachments success!")
}

func TestDeleteVolumeAttachment(t *testing.T) {
	atc, err := prepareVolumeAttachment(t)
	if err != nil {
		t.Error("Failed to run volume attachment prepare function:", err)
		return
	}
	defer cleanVolumeIfFailedOrFinished(t, atc.VolumeId)

	t.Log("Start deleting volume attachment...")
	if err := c.DeleteVolumeAttachment(atc.Id, nil); err != nil {
		t.Error("delete volume attachment failed:", err)
		return
	}
	t.Log("Delete volume attachment success!")
}

func TestCreateVolumeSnapshot(t *testing.T) {
	vol, err := prepareVolume(t)
	if err != nil {
		t.Error("Failed to run volume prepare function:", err)
		return
	}
	defer cleanVolumeIfFailedOrFinished(t, vol.Id)

	t.Log("Start creating volume snapshot...")
	var body = &model.VolumeSnapshotSpec{
		Name:        "test-snapshot",
		Description: "This is a snapshot test",
		VolumeId:    vol.Id,
	}
	snp, err := c.CreateVolumeSnapshot(body)
	if err != nil {
		t.Error("create volume snapshot failed:", err)
		return
	}
	// Check if the status of created volume snapshot is available.
	snp, _ = c.GetVolumeSnapshot(snp.Id)
	if snp.Status != model.VolumeSnapAvailable {
		t.Errorf("status expected is %s, got %s\n", model.VolumeSnapAvailable, snp.Status)
	}
	t.Log("create volume snapshot success!")

	t.Log("Start cleaing volume snapshot...")
	if err := c.DeleteVolumeSnapshot(snp.Id, nil); err != nil {
		t.Error("Clean volume snapshot failed:", err)
		return
	}
	t.Log("End cleaing volume snapshot...")
}

func TestGetVolumeSnapshot(t *testing.T) {
	snp, err := prepareVolumeSnapshot(t)
	if err != nil {
		t.Error("Failed to run volume snapshot prepare function:", err)
		return
	}
	defer cleanVolumeAndSnapshotIfFailedOrFinished(t, snp.VolumeId, snp.Id)

	t.Log("Start checking volume snapshot...")
	snp, err = c.GetVolumeSnapshot(snp.Id)
	if err != nil {
		t.Error("Check volume snapshot failed:", err)
		return
	}
	// Test the status of created volume snapshot.
	if snp.Status != model.VolumeSnapAvailable {
		t.Error("The status of volume snapshot is not available!")
		return
	}
	t.Log("Check volume snapshot success!")
}

func TestListVolumeSnapshots(t *testing.T) {
	snp, err := prepareVolumeSnapshot(t)
	if err != nil {
		t.Error("Failed to run volume snapshot prepare function:", err)
		return
	}
	defer cleanVolumeAndSnapshotIfFailedOrFinished(t, snp.VolumeId, snp.Id)

	t.Log("Start checking all volume snapshots...")
	if _, err := c.ListVolumeSnapshots(); err != nil {
		t.Error("list volume snapshots failed:", err)
		return
	}
	t.Log("Check all volume snapshots success!")
}

func TestDeleteVolumeSnapshot(t *testing.T) {
	snp, err := prepareVolumeSnapshot(t)
	if err != nil {
		t.Error("Failed to run volume snapshot prepare function:", err)
		return
	}
	defer cleanVolumeIfFailedOrFinished(t, snp.VolumeId)

	t.Log("Start deleting volume snapshot...")
	if err := c.DeleteVolumeSnapshot(snp.Id, nil); err != nil {
		t.Error("delete volume snapshot failed:", err)
		return
	}
	t.Log("Delete volume snapshot success!")
}

func TestUpdateVolumeSnapshot(t *testing.T) {
	snp, err := prepareVolumeSnapshot(t)
	if err != nil {
		t.Error("Failed to run volume snapshot prepare function:", err)
		return
	}
	defer cleanVolumeAndSnapshotIfFailedOrFinished(t, snp.VolumeId, snp.Id)

	t.Log("Start updating volume snapshot...")
	var body = &model.VolumeSnapshotSpec{
		Name:        "Update Volume Snapshot Name",
		Description: "Update Volume Snapshot Description",
	}
	newSnp, err := c.UpdateVolumeSnapshot(snp.Id, body)
	if err != nil {
		t.Error("update volume snapshot failed:", err)
		return
	}
	// Check if the status of updated volume snapshot is available.
	newSnp, _ = c.GetVolumeSnapshot(newSnp.Id)
	if newSnp.Status != model.VolumeSnapAvailable {
		t.Errorf("status expected is %s, got %s\n", model.VolumeSnapAvailable, newSnp.Status)
	}
	t.Log("Update volume snapshot success!")
}

func prepareVolume(t *testing.T) (*model.VolumeSpec, error) {
	t.Log("Start preparing volume...")
	var body = &model.VolumeSpec{
		Name:        "test",
		Description: "This is a test",
		Size:        int64(1),
	}
	vol, err := c.CreateVolume(body)
	if err != nil {
		t.Error("prepare volume failed:", err)
		return nil, err
	}
	if vol, _ = c.GetVolume(vol.Id); vol.Status != model.VolumeAvailable {
		return nil, fmt.Errorf("The status of volume is not available!")
	}

	t.Log("End preparing volume...")
	return vol, nil
}

func prepareVolumeAttachment(t *testing.T) (*model.VolumeAttachmentSpec, error) {
	vol, err := prepareVolume(t)
	if err != nil {
		t.Error("Failed to run volume prepare function:", err)
		return nil, err
	}

	t.Log("Start preparing volume attachment...")
	host, _ := os.Hostname()
	var body = &model.VolumeAttachmentSpec{
		VolumeId: vol.Id,
		HostInfo: model.HostInfo{
			Host:      host,
			Platform:  runtime.GOARCH,
			OsType:    runtime.GOOS,
			Ip:        getHostIp(),
			Initiator: localIqn,
		},
		AccessProtocol: iscsiProtocol,
	}
	atc, err := c.CreateVolumeAttachment(body)
	if err != nil {
		t.Error("prepare volume attachment failed:", err)
		// Run volume clean function if failed to prepare volume attachment.
		cleanVolumeIfFailedOrFinished(t, atc.VolumeId)
		return nil, err
	}
	atc, _ = c.GetVolumeAttachment(atc.Id)
	if atc.Status != model.VolumeAttachAvailable {
		// Run volume clean function if failed to prepare volume attachment.
		cleanVolumeIfFailedOrFinished(t, atc.VolumeId)
		return nil, fmt.Errorf("The status of volume attachment is not available!")
	}

	t.Log("End preparing volume attachment...")
	return atc, nil
}

func prepareVolumeSnapshot(t *testing.T) (*model.VolumeSnapshotSpec, error) {
	vol, err := prepareVolume(t)
	if err != nil {
		t.Error("Failed to run volume prepare function:", err)
		return nil, err
	}

	t.Log("Start preparing volume snapshot...")
	var body = &model.VolumeSnapshotSpec{
		Name:        "test-snapshot",
		Description: "This is a snapshot test",
		VolumeId:    vol.Id,
	}
	snp, err := c.CreateVolumeSnapshot(body)
	if err != nil {
		t.Error("prepare volume snapshot failed:", err)
		// Run volume clean function if failed to prepare volume snapshot.
		cleanVolumeIfFailedOrFinished(t, snp.VolumeId)
		return nil, err
	}
	if snp, _ = c.GetVolumeSnapshot(snp.Id); snp.Status != model.VolumeSnapAvailable {
		// Run volume clean function if failed to prepare volume snapshot.
		cleanVolumeIfFailedOrFinished(t, snp.VolumeId)
		return nil, fmt.Errorf("The status of volume snapshot is not available!")
	}

	t.Log("End preparing volume snapshot...")
	return snp, nil
}

func cleanVolumeIfFailedOrFinished(t *testing.T, volID string) error {
	t.Log("Start cleaning volume...")
	if err := c.DeleteVolume(volID, nil); err != nil {
		t.Error("Clean volume failed:", err)
		return err
	}
	t.Log("End cleaning volume...")
	return nil
}

func cleanVolumeAndAttachmentIfFailedOrFinished(t *testing.T, volID, atcID string) error {
	t.Log("Start cleaning volume attachment...")

	if err := c.DeleteVolumeAttachment(atcID, nil); err != nil {
		t.Error("Clean volume attachment failed:", err)
		return err
	}
	t.Log("End cleaning volume attachment...")

	t.Log("Start cleaning volume...")
	if err := c.DeleteVolume(volID, nil); err != nil {
		t.Error("Clean volume failed:", err)
		return err
	}
	t.Log("End cleaning volume...")
	return nil
}

func cleanVolumeAndSnapshotIfFailedOrFinished(t *testing.T, volID, snpID string) error {
	t.Log("Start cleaing volume snapshot...")
	if err := c.DeleteVolumeSnapshot(snpID, nil); err != nil {
		t.Error("Clean volume snapshot failed:", err)
		return err
	}
	t.Log("End cleaning volume snapshot...")

	t.Log("Start cleaning volume...")
	if err := c.DeleteVolume(volID, nil); err != nil {
		t.Error("Clean volume failed:", err)
		return err
	}
	t.Log("End cleaning volume...")
	return nil
}

// getHostIp return Host IP
func getHostIp() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "127.0.0.1"
	}

	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			return ipnet.IP.String()
		}
	}

	return "127.0.0.1"
}

func prepareVolumeGroup(t *testing.T) (*model.VolumeGroupSpec, error) {
	vol, err := prepareVolume(t)
	if err != nil {
		t.Error("Failed to run volume prepare function: ", err)
		return nil, err
	}
	vol1, err := prepareVolume(t)
	if err != nil {
		t.Error("Failed to run volume prepare function: ", err)
		return nil, err
	}

	t.Log("starting to prepare VolumeGroupSpec......")
	var body = &model.VolumeGroupSpec{
		Name:        "VGStest",
		Description: "This is a test",
		AddVolumes:  []string{vol.Id, vol1.Id},
		Profiles:    []string{profileId},
	}
	vg, err := c.CreateVolumeGroup(body)
	if err != nil {
		t.Error("prepare volume group failed: ", err)
		cleanVolumeIfFailedOrFinished(t, vol.Id)
		cleanVolumeIfFailedOrFinished(t, vol1.Id)
		return nil, err
	}
	vg, _ = c.GetVolumeGroup(vg.Id)
	if vg.Status != model.VolumeGroupAvailable {
		cleanVolumeIfFailedOrFinished(t, vol.Id)
		cleanVolumeIfFailedOrFinished(t, vol1.Id)
		return nil, fmt.Errorf("The status of volume group is not available!")
	}

	t.Log("End preparing volume group...")
	return vg, nil
}

func cleanVolumeAndGroupIfFailedOrFinished(t *testing.T, vgId string, body *model.VolumeGroupSpec) error {
	t.Log("Start cleaning volume group...")
	if err := c.DeleteVolumeGroup(vgId, body); err != nil {
		t.Error("Clean volume group failed:", err)
		return err
	}
	t.Log("End cleaning volume group...")

	//t.Log("Start cleaning volume...")
	//for i := range body.AddVolumes {
	//	if err := c.DeleteVolume(body.AddVolumes[i], nil); err != nil {
	//		t.Error("Clean volume failed: ", err)
	//		return err
	//	}
	//}
	//t.Log("End cleaning volume...")
	return nil
}

func TestCreateVolumeGroup(t *testing.T) {
	vol, err := prepareVolume(t)
	if err != nil {
		t.Error("Failed to run volume prepare function : ", err)
		return
	}
	vol1, err := prepareVolume(t)
	if err != nil {
		t.Error("Failed to run volume prepare function : ", err)
		return
	}

	t.Log("Start creating volume group...")
	var body = &model.VolumeGroupSpec{
		Name:        "testvolumegroup",
		Description: "This is a volume group test",
		AddVolumes:  []string{vol.Id, vol1.Id},
		Profiles:    []string{profileId},
	}
	vg, err := c.CreateVolumeGroup(body)
	if err != nil {
		t.Error("create volume group failed : ", err)
		cleanVolumeIfFailedOrFinished(t, vol.Id)
		cleanVolumeIfFailedOrFinished(t, vol1.Id)
		return
	}
	// Check if the status of created volume group is available.
	vg, _ = c.GetVolumeGroup(vg.Id)
	if vg.Status != model.VolumeGroupAvailable {
		t.Errorf("status expected is %s, got %s\n", model.VolumeGroupAvailable, vg.Status)
		cleanVolumeIfFailedOrFinished(t, vol.Id)
		cleanVolumeIfFailedOrFinished(t, vol.Id)
		return
	}
	t.Log("create volume group success!")

	t.Log("Starting cleaning volume group...")
	if err := c.DeleteVolumeGroup(vg.Id, vg); err != nil {
		t.Error("Clean volume group failed : ", err)
		cleanVolumeIfFailedOrFinished(t, vol.Id)
		cleanVolumeIfFailedOrFinished(t, vol.Id)
		return
	}
	t.Log("End cleaning volume group...")
}

func TestGetVolumeGroup(t *testing.T) {
	vg, err := prepareVolumeGroup(t)
	if err != nil {
		t.Error("prepare Volume Group failed :", err)
		return
	}
	defer cleanVolumeAndGroupIfFailedOrFinished(t, vg.Id, vg)

	t.Log("Start checking volume group...")
	vg, err = c.GetVolumeGroup(vg.Id)
	if err != nil {
		t.Error("Check volume group failed:", err)
		return
	}
	if vg.Status != model.VolumeGroupAvailable {
		t.Error("The status of volume group is not available!")
		return
	}
	t.Log("Check volume group success!")
}

func TestListVolumeGroups(t *testing.T) {
	vg, err := prepareVolumeGroup(t)
	if err != nil {
		t.Error("Failed to run volume group prepare function :", err)
		return
	}
	defer cleanVolumeAndGroupIfFailedOrFinished(t, vg.Id, vg)

	t.Log("Start checking all volume group...")
	if _, err := c.ListVolumeGroups(); err != nil {
		t.Error("list volume groups failed :", err)
		return
	}
	t.Log("Check all volume groups success!")
}

func TestDeleteVolumeGroup(t *testing.T) {
	vg, err := prepareVolumeGroup(t)
	if err != nil {
		t.Error("Failed to run volume group prepare function: ", err)
		return
	}

	t.Log("Start delete volume group...")
	if err := c.DeleteVolumeGroup(vg.Id, vg); err != nil {
		t.Error("delete volume group failed :", err)
		cleanVolumeAndGroupIfFailedOrFinished(t, vg.Id, vg)
		return
	}
	t.Log("Delete volume group success")
	//for i := range vg.AddVolumes {
	//	cleanVolumeIfFailedOrFinished(t, vg.AddVolumes[i])
	//}
}

func TestUpdateVolumeGroup(t *testing.T) {
	vg, err := prepareVolumeGroup(t)
	if err != nil {
		t.Error("Failed to run volume group prepare function: ", err)
		return
	}
	defer cleanVolumeAndGroupIfFailedOrFinished(t, vg.Id, vg)

	t.Log("Start updating volume group...")
	var body = &model.VolumeGroupSpec{
		Name:        "Update Volume Group Name",
		Description: "Update Volume Group Description",
		AddVolumes:  []string{},
	}
	newVg, err := c.UpdateVolumeGroup(vg.Id, body)
	if err != nil {
		t.Error("Update volume group failed: ", err)
		return
	}
	// Check if the status of created volume group is available.
	newVg, _ = c.GetVolumeGroup(newVg.Id)
	if newVg.Status != model.VolumeGroupAvailable {
		t.Errorf("status expected is %s, got %s\n", model.VolumeGroupAvailable, newVg.Status)
	}
	t.Log("Update volume group success!")
}
