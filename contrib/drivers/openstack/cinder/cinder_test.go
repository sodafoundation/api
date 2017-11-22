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
This module implements cinder driver for OpenSDS. Cinder driver will pass
these operation requests about volume to gophercloud which is an OpenStack
Go SDK.

*/

package cinder

import (
	"encoding/json"
	"testing"

	"github.com/bouk/monkey"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/extensions/schedulerstats"
	snapshotsv2 "github.com/gophercloud/gophercloud/openstack/blockstorage/v2/snapshots"
	volumesv2 "github.com/gophercloud/gophercloud/openstack/blockstorage/v2/volumes"
	"github.com/gophercloud/gophercloud/pagination"
	pb "github.com/opensds/opensds/pkg/dock/proto"
	"github.com/opensds/opensds/pkg/utils/config"
)

func TestSetup(t *testing.T) {
	var opt gophercloud.AuthOptions
	defer monkey.UnpatchAll()
	monkey.Patch(openstack.AuthenticatedClient,
		func(options gophercloud.AuthOptions) (*gophercloud.ProviderClient, error) {
			opt = options
			return &gophercloud.ProviderClient{}, nil
		})
	monkey.Patch(openstack.NewBlockStorageV2,
		func(client *gophercloud.ProviderClient, eo gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error) {
			return &gophercloud.ServiceClient{}, nil
		})

	config.CONF.OsdsDock.Backends.Cinder.ConfigPath = "testdata/cinder.yaml"
	d := Driver{}
	d.Setup()
	if opt.IdentityEndpoint != "http://192.168.56.104/identity" {
		t.Error("IdentityEndpoint error.")
	}
	if opt.DomainID != "Default" {
		t.Error("DomainID error.")
	}
	if opt.DomainName != "Default" {
		t.Error("DomainName error.")
	}
	if opt.Username != "admin" {
		t.Error("Username error.")
	}
	if opt.Password != "admin" {
		t.Error("Password error.")
	}
	if opt.TenantID != "04154b841eb644a3947506c54fa73c76" {
		t.Error("TenantID error.")
	}
	if opt.TenantName != "admin" {
		t.Error("TenantName error.")
	}

	if d.conf.Pool["pool1"].DiskType != "SSD" {
		t.Error("Test config pool1 DiskType error")
	}
	if d.conf.Pool["pool1"].IOPS != 1000 {
		t.Error("Test config pool1 IOPS error")
	}
	if d.conf.Pool["pool1"].BandWidth != 1000 {
		t.Error("Test config pool1 BandWidth error")
	}
	if d.conf.Pool["pool1"].AZ != "nova-01" {
		t.Error("Test config pool1 AZ error")
	}

	if d.conf.Pool["pool2"].DiskType != "SAS" {
		t.Error("Test config pool2 DiskType error")
	}
	if d.conf.Pool["pool2"].IOPS != 800 {
		t.Error("Test config pool2 IOPS error")
	}
	if d.conf.Pool["pool2"].BandWidth != 800 {
		t.Error("Test config pool2 BandWidth error")
	}
	if d.conf.Pool["pool2"].AZ != "nova-02" {
		t.Error("Test config pool2 AZ error")
	}
}

var volumeResp = `
{
    "volume": {
        "status": "creating",
        "migration_status": null,
        "user_id": "0eea4eabcf184061a3b6db1e0daaf010",
        "attachments": [],
        "links": [
            {
                "href": "http://23.253.248.171:8776/v2/bab7d5c60cd041a0a36f7c4b6e1dd978/volumes/6edbc2f4-1507-44f8-ac0d-eed1d2608d38",
                "rel": "self"
            },
            {
                "href": "http://23.253.248.171:8776/bab7d5c60cd041a0a36f7c4b6e1dd978/volumes/6edbc2f4-1507-44f8-ac0d-eed1d2608d38",
                "rel": "bookmark"
            }
        ],
        "availability_zone": "nova",
        "bootable": "false",
        "encrypted": false,
        "created_at": "2015-11-29T03:01:44.000000",
        "description": "OpenSDS testing volume",
        "updated_at": null,
        "volume_type": "lvmdriver-1",
        "name": "test1",
        "replication_status": "disabled",
        "consistencygroup_id": null,
        "source_volid": null,
        "snapshot_id": null,
        "multiattach": false,
        "metadata": {},
        "id": "6edbc2f4-1507-44f8-ac0d-eed1d2608d38",
        "size": 2
    }
}`

func TestCreateVolume(t *testing.T) {
	defer monkey.UnpatchAll()
	monkey.Patch(volumesv2.Create,
		func(client *gophercloud.ServiceClient, opts volumesv2.CreateOptsBuilder) (r volumesv2.CreateResult) {
			json.Unmarshal([]byte(volumeResp), &r.Body)
			return
		})

	opt := &pb.CreateVolumeOpts{
		Name:             "test1",
		Description:      "OpenSDS testing volume",
		AvailabilityZone: "nova",
		Size:             2,
	}
	d := Driver{}
	resp, err := d.CreateVolume(opt)

	if err != nil {
		t.Error("Create volume error")
	}
	if resp.Id != "6edbc2f4-1507-44f8-ac0d-eed1d2608d38" {
		t.Error("Create volume Id error.")
	}
	if resp.Name != "test1" {
		t.Error("Create volume Name error.")
	}
	if resp.Description != "OpenSDS testing volume" {
		t.Error("Create volume Description error.")
	}
	if resp.Size != 2 {
		t.Error("Create volume Size error.")
	}
	if resp.AvailabilityZone != "nova" {
		t.Error("Create volume AvailabilityZone error.")
	}
	if resp.Status != "creating" {
		t.Error("Create volume Status error.")
	}
}

func TestPullVolume(t *testing.T) {
	defer monkey.UnpatchAll()
	monkey.Patch(volumesv2.Get,
		func(client *gophercloud.ServiceClient, id string) (r volumesv2.GetResult) {
			json.Unmarshal([]byte(volumeResp), &r.Body)
			return
		})
	d := Driver{}
	resp, err := d.PullVolume("6edbc2f4-1507-44f8-ac0d-eed1d2608d38")
	if err != nil {
		t.Error("Get volume error")
	}
	if resp.Id != "6edbc2f4-1507-44f8-ac0d-eed1d2608d38" {
		t.Error("Get volume Id error.")
	}
	if resp.Name != "test1" {
		t.Error("Get volume Name error.")
	}
	if resp.Description != "OpenSDS testing volume" {
		t.Error("Get volume Description error.")
	}
	if resp.Size != 2 {
		t.Error("Get volume Size error.")
	}
	if resp.AvailabilityZone != "nova" {
		t.Error("Get volume AvailabilityZone error.")
	}
	if resp.Status != "creating" {
		t.Error("Get volume Status error.")
	}
}

func TestDeleteVolume(t *testing.T) {
	defer monkey.UnpatchAll()
	monkey.Patch(volumesv2.Delete,
		func(client *gophercloud.ServiceClient, id string) (r volumesv2.DeleteResult) {
			json.Unmarshal([]byte(volumeResp), &r.Body)
			r.Err = nil
			return
		})
	opt := &pb.DeleteVolumeOpts{
		Id: "6edbc2f4-1507-44f8-ac0d-eed1d2608d38",
	}
	d := Driver{}
	err := d.DeleteVolume(opt)
	if err != nil {
		t.Error("Delete volume error")
	}
}

var snapshotResp = `
{
    "snapshot": {
        "status": "available",
        "os-extended-snapshot-attributes:progress": "100%",
        "description": "OpenSDS testing snapshot",
        "created_at": "2013-02-25T04:13:17.000000",
        "metadata": {},
        "volume_id": "5aa119a8-d25b-45a7-8d1b-88e127885635",
        "os-extended-snapshot-attributes:project_id": "0c2eba2c5af04d3f9e9d0d410b371fde",
        "size": 1,
        "id": "2bb856e1-b3d8-4432-a858-09e4ce939389",
        "name": "test1"
    }
}`

func TestCreateSnapshot(t *testing.T) {
	defer monkey.UnpatchAll()
	monkey.Patch(snapshotsv2.Create,
		func(client *gophercloud.ServiceClient, opts snapshotsv2.CreateOptsBuilder) (r snapshotsv2.CreateResult) {
			json.Unmarshal([]byte(snapshotResp), &r.Body)
			return
		})
	opt := &pb.CreateVolumeSnapshotOpts{
		Id:          "2bb856e1-b3d8-4432-a858-09e4ce939389",
		Name:        "test1",
		Description: "OpenSDS testing snapshot",
		VolumeId:    "5aa119a8-d25b-45a7-8d1b-88e127885635",
	}
	d := Driver{}
	resp, err := d.CreateSnapshot(opt)
	if err != nil {
		t.Error("Create volume snapshot error")
	}
	if resp.Id != "2bb856e1-b3d8-4432-a858-09e4ce939389" {
		t.Error("Create volume Id error.")
	}
	if resp.Name != "test1" {
		t.Error("Create volume snapshot Name error.")
	}
	if resp.Description != "OpenSDS testing snapshot" {
		t.Error("Create volume snapshot Description error.")
	}
	if resp.VolumeId != "5aa119a8-d25b-45a7-8d1b-88e127885635" {
		t.Error("Create volume Id error.")
	}

	if resp.Status != "available" {
		t.Error("Create volume snapshot Status error.")
	}
}

func TestPullSnapshot(t *testing.T) {
	defer monkey.UnpatchAll()
	monkey.Patch(snapshotsv2.Get,
		func(client *gophercloud.ServiceClient, id string) (r snapshotsv2.GetResult) {
			json.Unmarshal([]byte(snapshotResp), &r.Body)
			return
		})
	d := Driver{}
	resp, err := d.PullSnapshot("2bb856e1-b3d8-4432-a858-09e4ce939389")
	if err != nil {
		t.Error("Get volume snapshot error")
	}
	if resp.Id != "2bb856e1-b3d8-4432-a858-09e4ce939389" {
		t.Error("Get volume Id error.")
	}
	if resp.Name != "test1" {
		t.Error("Get volume snapshot Name error.")
	}
	if resp.Description != "OpenSDS testing snapshot" {
		t.Error("Get volume snapshot Description error.")
	}
	if resp.VolumeId != "5aa119a8-d25b-45a7-8d1b-88e127885635" {
		t.Error("Get volume Id error.")
	}

	if resp.Status != "available" {
		t.Error("Get volume snapshot Status error.")
	}
}

func TestDeleteSnapshot(t *testing.T) {
	defer monkey.UnpatchAll()
	monkey.Patch(snapshotsv2.Delete,
		func(client *gophercloud.ServiceClient, id string) (r snapshotsv2.DeleteResult) {
			json.Unmarshal([]byte(snapshotResp), &r.Body)
			return
		})
	opt := &pb.DeleteVolumeSnapshotOpts{Id: "2bb856e1-b3d8-4432-a858-09e4ce939389"}
	d := Driver{}
	err := d.DeleteSnapshot(opt)
	if err != nil {
		t.Error("Delete volume snapshot error")
	}
}

func TestListPools(t *testing.T) {
	defer monkey.UnpatchAll()
	monkey.Patch(openstack.AuthenticatedClient,
		func(options gophercloud.AuthOptions) (*gophercloud.ProviderClient, error) {
			return &gophercloud.ProviderClient{}, nil
		})
	monkey.Patch(openstack.NewBlockStorageV2,
		func(client *gophercloud.ProviderClient, eo gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error) {
			return &gophercloud.ServiceClient{}, nil
		})
	monkey.Patch(schedulerstats.List,
		func(client *gophercloud.ServiceClient, opts schedulerstats.ListOptsBuilder) pagination.Pager {
			return pagination.Pager{}
		})
	monkey.Patch(pagination.Pager.AllPages,
		func(p pagination.Pager) (pagination.Page, error) {
			return schedulerstats.StoragePoolPage{}, nil
		})
	monkey.Patch(schedulerstats.ExtractStoragePools,
		func(p pagination.Page) ([]schedulerstats.StoragePool, error) {
			pools := []schedulerstats.StoragePool{
				{
					Name: "pool1",
					Capabilities: schedulerstats.Capabilities{
						TotalCapacityGB: 100.0,
						FreeCapacityGB:  50.0,
					},
				},
				{
					Name: "pool2",
					Capabilities: schedulerstats.Capabilities{
						TotalCapacityGB: 1000.0,
						FreeCapacityGB:  500.0,
					},
				},
				{
					Name: "ShouldBeFilterd",
					Capabilities: schedulerstats.Capabilities{
						TotalCapacityGB: 1000.0,
						FreeCapacityGB:  500.0,
					},
				},
			}
			return pools, nil
		})
	config.CONF.OsdsDock.Backends.Cinder.ConfigPath = "testdata/cinder.yaml"
	d := Driver{}
	d.Setup()
	resp, err := d.ListPools()
	if err != nil {
		t.Error("Delete volume snapshot error")
	}
	if resp[0].Name != "pool1" {
		t.Error("List pool name error.")
	}
	if resp[0].TotalCapacity != 100 {
		t.Error("List pool TotalCapacity error.")
	}
	if resp[0].FreeCapacity != 50 {
		t.Error("List pool TotalCapacity error.")
	}
	if resp[0].Parameters["diskType"] != "SSD" {
		t.Error("List pool Parameters diskType error.")
	}
	if resp[0].Parameters["iops"].(int64) != 1000 {
		t.Error("List pool Parameters iops error.")
	}
	if resp[0].Parameters["bandwidth"].(int64) != 1000 {
		t.Error("List pool Parameters bandwidth error.")
	}
	if resp[1].Name != "pool2" {
		t.Error("List pool name error.")
	}
	if resp[1].TotalCapacity != 1000 {
		t.Error("List pool TotalCapacity error.")
	}
	if resp[1].FreeCapacity != 500 {
		t.Error("List pool TotalCapacity error.")
	}
	if len(resp) != 2 {
		t.Error("List pool number error")
	}
}
