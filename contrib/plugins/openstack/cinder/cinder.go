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

/*
This module implements cinder driver for OpenSDS. Cinder driver will pass
these operation requests about volume to gophercloud which is an OpenStack
Go SDK.

*/

package cinder

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"

	api "github.com/opensds/opensds/pkg/model"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/extensions/schedulerstats"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/extensions/volumeactions"
	snapshotsv2 "github.com/gophercloud/gophercloud/openstack/blockstorage/v2/snapshots"
	volumesv2 "github.com/gophercloud/gophercloud/openstack/blockstorage/v2/volumes"
)

var conf = CinderConfig{}

type Driver struct {
	// Current block storage version
	blockStoragev2 *gophercloud.ServiceClient
	blockStoragev3 *gophercloud.ServiceClient

	config CinderConfig
}

type CinderConfig struct {
	IdentityEndpoint string `json:"endpoint,omitempty"`
	DomainID         string `json:"-"`
	DomainName       string `json:"name,omitempty"`
	Username         string `json:"username,omitempty"`
	Password         string `json:"password,omitempty"`
	TenantID         string `json:"tenantId,omitempty"`
	TenantName       string `json:"tenantName,omitempty"`
}

func (d *Driver) Setup() {
	// Read cinder config file
	userJSON, err := ioutil.ReadFile("/etc/opensds/config.json")
	if err != nil {
		panic(err)
	}

	// Marshal the result
	if err = json.Unmarshal(userJSON, &conf); err != nil {
		panic(err)
	}

	d.config = conf

	opts := gophercloud.AuthOptions{
		IdentityEndpoint: d.config.IdentityEndpoint,
		Username:         d.config.Username,
		Password:         d.config.Password,
		TenantID:         d.config.TenantID,
		TenantName:       d.config.TenantName,
	}

	provider, err := openstack.AuthenticatedClient(opts)
	if err != nil {
		panic(err)
	}

	d.blockStoragev2, err = openstack.NewBlockStorageV2(provider, gophercloud.EndpointOpts{})
	if err != nil {
		panic(err)
	}

	return
}

func (d *Driver) Unset() { return }

func (d *Driver) CreateVolume(name string, size int64) (*api.VolumeSpec, error) {
	//Configure create request body.
	opts := &volumesv2.CreateOpts{
		Name: name,
		Size: int(size),
	}

	vol, err := volumesv2.Create(d.blockStoragev2, opts).Extract()
	if err != nil {
		log.Println("[Error] Cannot create volume:", err)
		return new(api.VolumeSpec), err
	}

	return &api.VolumeSpec{
		BaseModel: &api.BaseModel{
			Id: vol.ID,
		},
		Name:             vol.Name,
		Description:      vol.Description,
		Size:             int64(vol.Size),
		AvailabilityZone: vol.AvailabilityZone,
		Status:           vol.Status,
	}, nil
}

func (d *Driver) GetVolume(volID string) (*api.VolumeSpec, error) {
	vol, err := volumesv2.Get(d.blockStoragev2, volID).Extract()
	if err != nil {
		log.Println("[Error] Cannot get volume:", err)
		return new(api.VolumeSpec), err
	}

	return &api.VolumeSpec{
		BaseModel: &api.BaseModel{
			Id: vol.ID,
		},
		Name:             vol.Name,
		Description:      vol.Description,
		Size:             int64(vol.Size),
		AvailabilityZone: vol.AvailabilityZone,
		Status:           vol.Status,
	}, nil
}

func (d *Driver) DeleteVolume(volID string) error {
	if err := volumesv2.Delete(d.blockStoragev2, volID).ExtractErr(); err != nil {
		log.Println("[Error] Cannot delete volume:", err)
		return err
	}

	return nil
}

func (d *Driver) InitializeConnection(volID string, doLocalAttach, multiPath bool, hostInfo *api.HostInfo) (*api.ConnectionInfo, error) {
	opts := &volumeactions.InitializeConnectionOpts{
		IP:        hostInfo.GetIp(),
		Host:      hostInfo.GetHost(),
		Initiator: hostInfo.GetInitiator(),
		Platform:  hostInfo.GetPlatform(),
		OSType:    hostInfo.GetOsType(),
		Multipath: &multiPath,
	}

	conn, err := volumeactions.InitializeConnection(d.blockStoragev2, volID, opts).Extract()
	if err != nil {
		return new(api.ConnectionInfo), err
	}

	return &api.ConnectionInfo{
		DriverVolumeType: "iscsi",
		ConnectionData:   conn,
	}, nil
}

func (d *Driver) AttachVolume(volID, host, mountpoint string) error {
	vol, err := volumesv2.Get(d.blockStoragev2, volID).Extract()
	if err != nil {
		return err
	}

	if vol.Status != "available" && !vol.Multiattach {
		err = errors.New("The status of volume is not available!")
		log.Println("Cannot attach volume:", err)
		return err
	}

	opts := &volumeactions.AttachOpts{
		HostName:   host,
		MountPoint: mountpoint,
	}

	if err = volumeactions.Attach(d.blockStoragev2, volID, opts).ExtractErr(); err != nil {
		log.Println("Cannot attach volume:", err)
		return err
	}

	return nil
}

func (d *Driver) DetachVolume(volID string) error {
	vol, err := volumesv2.Get(d.blockStoragev2, volID).Extract()
	if err != nil {
		return err
	}

	if vol.Status != "in-use" {
		err = errors.New("The status of volume is not in-use!")
		log.Println("Cannot detach volume:", err)
		return err
	}

	opts := &volumeactions.DetachOpts{
		AttachmentID: vol.Attachments[0].ID,
	}

	if err = volumeactions.Detach(d.blockStoragev2, volID, opts).ExtractErr(); err != nil {
		log.Println("Cannot detach volume:", err)
		return err
	}

	return nil
}

func (d *Driver) CreateSnapshot(name, volID, description string) (*api.VolumeSnapshotSpec, error) {
	opts := &snapshotsv2.CreateOpts{
		VolumeID:    volID,
		Name:        name,
		Description: description,
	}

	snp, err := snapshotsv2.Create(d.blockStoragev2, opts).Extract()
	if err != nil {
		log.Println("[Error] Cannot create snapshot:", err)
		return new(api.VolumeSnapshotSpec), err
	}

	return &api.VolumeSnapshotSpec{
		BaseModel: &api.BaseModel{
			Id: snp.ID,
		},
		Name:        snp.Name,
		Description: snp.Description,
		Size:        int64(snp.Size),
		Status:      snp.Status,
		VolumeId:    volID,
	}, nil
}

func (d *Driver) GetSnapshot(snapID string) (*api.VolumeSnapshotSpec, error) {
	snp, err := snapshotsv2.Get(d.blockStoragev2, snapID).Extract()
	if err != nil {
		log.Println("[Error] Cannot get snapshot:", err)
		return new(api.VolumeSnapshotSpec), err
	}

	return &api.VolumeSnapshotSpec{
		BaseModel: &api.BaseModel{
			Id: snp.ID,
		},
		Name:        snp.Name,
		Description: snp.Description,
		Size:        int64(snp.Size),
		Status:      snp.Status,
		VolumeId:    snp.VolumeID,
	}, nil
}

func (d *Driver) DeleteSnapshot(snapID string) error {
	if err := snapshotsv2.Delete(d.blockStoragev2, snapID).ExtractErr(); err != nil {
		log.Println("[Error] Cannot delete snapshot:", err)
		return err
	}

	return nil
}

func (d *Driver) ListPools() (*[]api.StoragePoolSpec, error) {
	opts := &schedulerstats.ListOpts{}

	pages, err := schedulerstats.List(d.blockStoragev2, opts).AllPages()
	if err != nil {
		log.Println("[Error] Cannot list storage pools:", err)
		return new([]api.StoragePoolSpec), err
	}

	polpages, err := schedulerstats.ExtractStoragePools(pages)
	if err != nil {
		log.Println("[Error] Cannot extract storage pools:", err)
		return new([]api.StoragePoolSpec), err
	}

	var pols []api.StoragePoolSpec
	for _, page := range polpages {
		pol := api.StoragePoolSpec{
			Name:          page.Name,
			TotalCapacity: int64(page.Capabilities.TotalCapacityGB),
			FreeCapacity:  int64(page.Capabilities.FreeCapacityGB),
		}

		pols = append(pols, pol)
	}
	return &pols, nil
}
