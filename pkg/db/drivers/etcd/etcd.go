// Copyright (c) 2016 Huawei Technologies Co., Ltd. All Rights Reserved.
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
This module implements the etcd database operation of data structure
defined in api module.

*/

package etcd

import (
	"encoding/json"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/coreos/etcd/clientv3"
	api "github.com/opensds/opensds/pkg/model"
)

const (
	prefix  = "/api/v1alpha/block"
	timeOut = 3 * time.Second
)

var c = &client{}

func Init(edps []string) *client {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   edps,
		DialTimeout: timeOut,
	})
	if err != nil {
		cli.Close()
		panic(err)
	}

	c.cli = cli
	return c
}

type client struct {
	cli  *clientv3.Client
	lock sync.Mutex
}

func (c *client) CreateDock(dck *api.DockSpec) (*api.DockSpec, error) {
	dckBody, err := json.Marshal(dck)
	if err != nil {
		return &api.DockSpec{}, err
	}

	dbReq := &Request{
		Url:     GenerateUrl(prefix, "docks", dck.GetId()),
		Content: string(dckBody),
	}
	dbRes := c.Create(dbReq)
	if dbRes.Status != "Success" {
		log.Println("[Error] When create dock in db:", dbRes.Error)
		return &api.DockSpec{}, errors.New(dbRes.Error)
	}

	return dck, nil
}

func (c *client) GetDock(dckID string) (*api.DockSpec, error) {
	dbReq := &Request{
		Url: GenerateUrl(prefix, "docks", dckID),
	}
	dbRes := c.Get(dbReq)
	if dbRes.Status != "Success" {
		log.Println("[Error] When get dock in db:", dbRes.Error)
		return &api.DockSpec{}, errors.New(dbRes.Error)
	}

	var dck = &api.DockSpec{}
	if err := json.Unmarshal([]byte(dbRes.Message[0]), dck); err != nil {
		log.Println("[Error] When parsing dock in db:", dbRes.Error)
		return &api.DockSpec{}, errors.New(dbRes.Error)
	}
	return dck, nil
}

func (c *client) ListDocks() (*[]api.DockSpec, error) {
	dbReq := &Request{
		Url: GenerateUrl(prefix, "docks"),
	}
	dbRes := c.List(dbReq)
	if dbRes.Status != "Success" {
		log.Println("[Error] When list docks in db:", dbRes.Error)
		return &[]api.DockSpec{}, errors.New(dbRes.Error)
	}

	var dcks = []api.DockSpec{}
	if len(dbRes.Message) == 0 {
		return &dcks, nil
	}
	for _, msg := range dbRes.Message {
		var dck = api.DockSpec{}
		if err := json.Unmarshal([]byte(msg), &dck); err != nil {
			log.Println("[Error] When parsing dock in db:", dbRes.Error)
			return &[]api.DockSpec{}, errors.New(dbRes.Error)
		}
		dcks = append(dcks, dck)
	}
	return &dcks, nil
}

func (c *client) UpdateDock(dckID, name, desp string) (*api.DockSpec, error) {
	dck, err := c.GetDock(dckID)
	if err != nil {
		return &api.DockSpec{}, err
	}
	if name != "" {
		dck.Name = name
	}
	if desp != "" {
		dck.Description = desp
	}
	dckBody, err := json.Marshal(dck)
	if err != nil {
		return &api.DockSpec{}, err
	}

	dbReq := &Request{
		Url:        GenerateUrl(prefix, "docks", dckID),
		NewContent: string(dckBody),
	}
	dbRes := c.Update(dbReq)
	if dbRes.Status != "Success" {
		log.Println("[Error] When update dock in db:", dbRes.Error)
		return &api.DockSpec{}, errors.New(dbRes.Error)
	}
	return dck, nil
}

func (c *client) DeleteDock(dckID string) error {
	dbReq := &Request{
		Url: GenerateUrl(prefix, "docks", dckID),
	}
	dbRes := c.Delete(dbReq)
	if dbRes.Status != "Success" {
		log.Println("[Error] When delete dock in db:", dbRes.Error)
		return errors.New(dbRes.Error)
	}
	return nil
}

func (c *client) CreatePool(pol *api.StoragePoolSpec) (*api.StoragePoolSpec, error) {
	polBody, err := json.Marshal(pol)
	if err != nil {
		return &api.StoragePoolSpec{}, err
	}

	dbReq := &Request{
		Url:     GenerateUrl(prefix, "pools", pol.GetId()),
		Content: string(polBody),
	}
	dbRes := c.Create(dbReq)
	if dbRes.Status != "Success" {
		log.Println("[Error] When create pol in db:", dbRes.Error)
		return &api.StoragePoolSpec{}, errors.New(dbRes.Error)
	}

	return pol, nil
}

func (c *client) GetPool(polID string) (*api.StoragePoolSpec, error) {
	dbReq := &Request{
		Url: GenerateUrl(prefix, "pools", polID),
	}
	dbRes := c.Get(dbReq)
	if dbRes.Status != "Success" {
		log.Println("[Error] When get pool in db:", dbRes.Error)
		return &api.StoragePoolSpec{}, errors.New(dbRes.Error)
	}

	var pol = &api.StoragePoolSpec{}
	if err := json.Unmarshal([]byte(dbRes.Message[0]), pol); err != nil {
		log.Println("[Error] When parsing pool in db:", dbRes.Error)
		return &api.StoragePoolSpec{}, errors.New(dbRes.Error)
	}
	return pol, nil
}

func (c *client) ListPools() (*[]api.StoragePoolSpec, error) {
	dbReq := &Request{
		Url: GenerateUrl(prefix, "pools"),
	}
	dbRes := c.List(dbReq)
	if dbRes.Status != "Success" {
		log.Println("[Error] When list pools in db:", dbRes.Error)
		return &[]api.StoragePoolSpec{}, errors.New(dbRes.Error)
	}

	var pols = []api.StoragePoolSpec{}
	if len(dbRes.Message) == 0 {
		return &pols, nil
	}
	for _, msg := range dbRes.Message {
		var pol = api.StoragePoolSpec{}
		if err := json.Unmarshal([]byte(msg), &pol); err != nil {
			log.Println("[Error] When parsing pool in db:", dbRes.Error)
			return &[]api.StoragePoolSpec{}, errors.New(dbRes.Error)
		}
		pols = append(pols, pol)
	}
	return &pols, nil
}

func (c *client) UpdatePool(polID, name, desp string, usedCapacity int64, used bool) (*api.StoragePoolSpec, error) {
	pol, err := c.GetPool(polID)
	if err != nil {
		return &api.StoragePoolSpec{}, err
	}
	if name != "" {
		pol.Name = name
	}
	if desp != "" {
		pol.Description = desp
	}
	polBody, err := json.Marshal(pol)
	if err != nil {
		return &api.StoragePoolSpec{}, err
	}

	dbReq := &Request{
		Url:        GenerateUrl(prefix, "pools", polID),
		NewContent: string(polBody),
	}
	dbRes := c.Update(dbReq)
	if dbRes.Status != "Success" {
		log.Println("[Error] When update pool in db:", dbRes.Error)
		return &api.StoragePoolSpec{}, errors.New(dbRes.Error)
	}
	return pol, nil
}

func (c *client) DeletePool(polID string) error {
	dbReq := &Request{
		Url: GenerateUrl(prefix, "pools", polID),
	}
	dbRes := c.Delete(dbReq)
	if dbRes.Status != "Success" {
		log.Println("[Error] When delete pool in db:", dbRes.Error)
		return errors.New(dbRes.Error)
	}
	return nil
}

func (c *client) CreateProfile(prf *api.ProfileSpec) (*api.ProfileSpec, error) {
	prfBody, err := json.Marshal(prf)
	if err != nil {
		return &api.ProfileSpec{}, err
	}

	dbReq := &Request{
		Url:     GenerateUrl(prefix, "profiles", prf.GetId()),
		Content: string(prfBody),
	}
	dbRes := c.Create(dbReq)
	if dbRes.Status != "Success" {
		log.Println("[Error] When create profile in db:", dbRes.Error)
		return &api.ProfileSpec{}, errors.New(dbRes.Error)
	}

	return prf, nil
}

func (c *client) GetProfile(prfID string) (*api.ProfileSpec, error) {
	dbReq := &Request{
		Url: GenerateUrl(prefix, "profiles", prfID),
	}
	dbRes := c.Get(dbReq)
	if dbRes.Status != "Success" {
		log.Println("[Error] When get profile in db:", dbRes.Error)
		return &api.ProfileSpec{}, errors.New(dbRes.Error)
	}

	var prf = &api.ProfileSpec{}
	if err := json.Unmarshal([]byte(dbRes.Message[0]), prf); err != nil {
		log.Println("[Error] When parsing profile in db:", dbRes.Error)
		return &api.ProfileSpec{}, errors.New(dbRes.Error)
	}
	return prf, nil
}

func (c *client) ListProfiles() (*[]api.ProfileSpec, error) {
	dbReq := &Request{
		Url: GenerateUrl(prefix, "profiles"),
	}
	dbRes := c.List(dbReq)
	if dbRes.Status != "Success" {
		log.Println("[Error] When list profiles in db:", dbRes.Error)
		return &[]api.ProfileSpec{}, errors.New(dbRes.Error)
	}

	var prfs = []api.ProfileSpec{}
	if len(dbRes.Message) == 0 {
		return &prfs, nil
	}
	for _, msg := range dbRes.Message {
		var prf = api.ProfileSpec{}
		if err := json.Unmarshal([]byte(msg), &prf); err != nil {
			log.Println("[Error] When parsing profile in db:", dbRes.Error)
			return &[]api.ProfileSpec{}, errors.New(dbRes.Error)
		}
		prfs = append(prfs, prf)
	}
	return &prfs, nil
}

func (c *client) UpdateProfile(prfID string, input *api.ProfileSpec) (*api.ProfileSpec, error) {
	prf, err := c.GetProfile(prfID)
	if err != nil {
		return &api.ProfileSpec{}, err
	}
	if name := input.GetName(); name != "" {
		prf.Name = name
	}
	if desp := input.GetDescription(); desp != "" {
		prf.Description = desp
	}
	if props := input.Extra; len(props) != 0 {
		return &api.ProfileSpec{}, errors.New("Failed to update extra properties!")
	}

	prfBody, err := json.Marshal(prf)
	if err != nil {
		return &api.ProfileSpec{}, err
	}

	dbReq := &Request{
		Url:        GenerateUrl(prefix, "profiles", prfID),
		NewContent: string(prfBody),
	}
	dbRes := c.Update(dbReq)
	if dbRes.Status != "Success" {
		log.Println("[Error] When update profile in db:", dbRes.Error)
		return &api.ProfileSpec{}, errors.New(dbRes.Error)
	}
	return prf, nil
}

func (c *client) DeleteProfile(prfID string) error {
	dbReq := &Request{
		Url: GenerateUrl(prefix, "profiles", prfID),
	}
	dbRes := c.Delete(dbReq)
	if dbRes.Status != "Success" {
		log.Println("[Error] When delete profile in db:", dbRes.Error)
		return errors.New(dbRes.Error)
	}
	return nil
}

func (c *client) AddExtraProperty(prfID string, ext api.ExtraSpec) (*api.ExtraSpec, error) {
	prf, err := c.GetProfile(prfID)
	if err != nil {
		return &api.ExtraSpec{}, err
	}

	for k, v := range ext {
		prf.Extra[k] = v
	}

	prf, err = c.CreateProfile(prf)
	if err != nil {
		return &api.ExtraSpec{}, err
	}
	return &prf.Extra, nil
}

func (c *client) ListExtraProperties(prfID string) (*api.ExtraSpec, error) {
	prf, err := c.GetProfile(prfID)
	if err != nil {
		return &api.ExtraSpec{}, err
	}
	return &prf.Extra, nil
}

func (c *client) RemoveExtraProperty(prfID, extraKey string) error {
	prf, err := c.GetProfile(prfID)
	if err != nil {
		return err
	}

	delete(prf.Extra, extraKey)
	prf, err = c.CreateProfile(prf)
	if err != nil {
		return err
	}
	return nil
}

func (c *client) CreateVolume(vol *api.VolumeSpec) (*api.VolumeSpec, error) {
	volBody, err := json.Marshal(vol)
	if err != nil {
		return &api.VolumeSpec{}, err
	}

	dbReq := &Request{
		Url:     GenerateUrl(prefix, "volumes", vol.GetId()),
		Content: string(volBody),
	}
	dbRes := c.Create(dbReq)
	if dbRes.Status != "Success" {
		log.Println("[Error] When create volume in db:", dbRes.Error)
		return &api.VolumeSpec{}, errors.New(dbRes.Error)
	}

	return vol, nil
}

func (c *client) GetVolume(volID string) (*api.VolumeSpec, error) {
	dbReq := &Request{
		Url: GenerateUrl(prefix, "volumes", volID),
	}
	dbRes := c.Get(dbReq)
	if dbRes.Status != "Success" {
		log.Println("[Error] When get volume in db:", dbRes.Error)
		return &api.VolumeSpec{}, errors.New(dbRes.Error)
	}

	var vol = &api.VolumeSpec{}
	if err := json.Unmarshal([]byte(dbRes.Message[0]), vol); err != nil {
		log.Println("[Error] When parsing volume in db:", dbRes.Error)
		return &api.VolumeSpec{}, errors.New(dbRes.Error)
	}
	return vol, nil
}

func (c *client) ListVolumes() (*[]api.VolumeSpec, error) {
	dbReq := &Request{
		Url: GenerateUrl(prefix, "volumes"),
	}
	dbRes := c.List(dbReq)
	if dbRes.Status != "Success" {
		log.Println("[Error] When list volumes in db:", dbRes.Error)
		return &[]api.VolumeSpec{}, errors.New(dbRes.Error)
	}

	var vols = []api.VolumeSpec{}
	if len(dbRes.Message) == 0 {
		return &vols, nil
	}
	for _, msg := range dbRes.Message {
		var vol = api.VolumeSpec{}
		if err := json.Unmarshal([]byte(msg), &vol); err != nil {
			log.Println("[Error] When parsing volume in db:", dbRes.Error)
			return &[]api.VolumeSpec{}, errors.New(dbRes.Error)
		}
		vols = append(vols, vol)
	}
	return &vols, nil
}

func (c *client) DeleteVolume(volID string) error {
	dbReq := &Request{
		Url: GenerateUrl(prefix, "volumes", volID),
	}
	dbRes := c.Delete(dbReq)
	if dbRes.Status != "Success" {
		log.Println("[Error] When delete volume in db:", dbRes.Error)
		return errors.New(dbRes.Error)
	}
	return nil
}

func (c *client) CreateVolumeAttachment(volID string, atc *api.VolumeAttachmentSpec) (*api.VolumeAttachmentSpec, error) {
	atcBody, err := json.Marshal(atc)
	if err != nil {
		return &api.VolumeAttachmentSpec{}, err
	}

	dbReq := &Request{
		Url:     GenerateUrl(prefix, "volume", volID, "attachments", atc.GetId()),
		Content: string(atcBody),
	}
	dbRes := c.Create(dbReq)
	if dbRes.Status != "Success" {
		log.Println("[Error] When create volume attachment in db:", dbRes.Error)
		return &api.VolumeAttachmentSpec{}, errors.New(dbRes.Error)
	}

	return atc, nil
}

func (c *client) GetVolumeAttachment(volID, atcID string) (*api.VolumeAttachmentSpec, error) {
	dbReq := &Request{
		Url: GenerateUrl(prefix, "volume", volID, "attachments", atcID),
	}
	dbRes := c.Get(dbReq)
	if dbRes.Status != "Success" {
		log.Println("[Error] When get volume attachment in db:", dbRes.Error)
		return &api.VolumeAttachmentSpec{}, errors.New(dbRes.Error)
	}

	var atc = &api.VolumeAttachmentSpec{}
	if err := json.Unmarshal([]byte(dbRes.Message[0]), atc); err != nil {
		log.Println("[Error] When parsing volume attachment in db:", dbRes.Error)
		return &api.VolumeAttachmentSpec{}, errors.New(dbRes.Error)
	}
	return atc, nil
}

func (c *client) ListVolumeAttachments(volID string) (*[]api.VolumeAttachmentSpec, error) {
	dbReq := &Request{
		Url: GenerateUrl(prefix, "volume", volID, "attachments"),
	}
	dbRes := c.List(dbReq)
	if dbRes.Status != "Success" {
		log.Println("[Error] When list volume attachments in db:", dbRes.Error)
		return &[]api.VolumeAttachmentSpec{}, errors.New(dbRes.Error)
	}

	var atcs = []api.VolumeAttachmentSpec{}
	if len(dbRes.Message) == 0 {
		return &atcs, nil
	}
	for _, msg := range dbRes.Message {
		var atc = api.VolumeAttachmentSpec{}
		if err := json.Unmarshal([]byte(msg), &atc); err != nil {
			log.Println("[Error] When parsing volume attachment in db:", dbRes.Error)
			return &[]api.VolumeAttachmentSpec{}, errors.New(dbRes.Error)
		}
		atcs = append(atcs, atc)
	}
	return &atcs, nil
}

func (c *client) UpdateVolumeAttachment(volID, atcID, mountpoint string, hostInfo *api.HostInfo) (*api.VolumeAttachmentSpec, error) {
	atc, err := c.GetVolumeAttachment(volID, atcID)
	if err != nil {
		return &api.VolumeAttachmentSpec{}, err
	}

	atc.HostInfo = hostInfo
	atc.Mountpoint = mountpoint
	atcBody, err := json.Marshal(atc)
	if err != nil {
		return &api.VolumeAttachmentSpec{}, err
	}

	dbReq := &Request{
		Url:        GenerateUrl(prefix, "volume", volID, "attachments", atcID),
		NewContent: string(atcBody),
	}
	dbRes := c.Update(dbReq)
	if dbRes.Status != "Success" {
		log.Println("[Error] When update volume attachment in db:", dbRes.Error)
		return &api.VolumeAttachmentSpec{}, errors.New(dbRes.Error)
	}
	return atc, nil
}

func (c *client) DeleteVolumeAttachment(volID, atcID string) error {
	dbReq := &Request{
		Url: GenerateUrl(prefix, "volume", volID, "attachments", atcID),
	}
	dbRes := c.Delete(dbReq)
	if dbRes.Status != "Success" {
		log.Println("[Error] When delete volume attachment in db:", dbRes.Error)
		return errors.New(dbRes.Error)
	}
	return nil
}

func (c *client) CreateVolumeSnapshot(snp *api.VolumeSnapshotSpec) (*api.VolumeSnapshotSpec, error) {
	snpBody, err := json.Marshal(snp)
	if err != nil {
		return &api.VolumeSnapshotSpec{}, err
	}

	dbReq := &Request{
		Url:     GenerateUrl(prefix, "volume", "snapshots", snp.GetId()),
		Content: string(snpBody),
	}
	dbRes := c.Create(dbReq)
	if dbRes.Status != "Success" {
		log.Println("[Error] When create volume snapshot in db:", dbRes.Error)
		return &api.VolumeSnapshotSpec{}, errors.New(dbRes.Error)
	}

	return snp, nil
}

func (c *client) GetVolumeSnapshot(snpID string) (*api.VolumeSnapshotSpec, error) {
	dbReq := &Request{
		Url: GenerateUrl(prefix, "volume", "snapshots", snpID),
	}
	dbRes := c.Get(dbReq)
	if dbRes.Status != "Success" {
		log.Println("[Error] When get volume attachment in db:", dbRes.Error)
		return &api.VolumeSnapshotSpec{}, errors.New(dbRes.Error)
	}

	var vs = &api.VolumeSnapshotSpec{}
	if err := json.Unmarshal([]byte(dbRes.Message[0]), vs); err != nil {
		log.Println("[Error] When parsing volume snapshot in db:", dbRes.Error)
		return &api.VolumeSnapshotSpec{}, errors.New(dbRes.Error)
	}
	return vs, nil
}

func (c *client) ListVolumeSnapshots() (*[]api.VolumeSnapshotSpec, error) {
	dbReq := &Request{
		Url: GenerateUrl(prefix, "volume", "snapshots"),
	}
	dbRes := c.List(dbReq)
	if dbRes.Status != "Success" {
		log.Println("[Error] When list volume snapshots in db:", dbRes.Error)
		return &[]api.VolumeSnapshotSpec{}, errors.New(dbRes.Error)
	}

	var vss = []api.VolumeSnapshotSpec{}
	if len(dbRes.Message) == 0 {
		return &vss, nil
	}
	for _, msg := range dbRes.Message {
		var vs = api.VolumeSnapshotSpec{}
		if err := json.Unmarshal([]byte(msg), &vs); err != nil {
			log.Println("[Error] When parsing volume snapshot in db:", dbRes.Error)
			return &[]api.VolumeSnapshotSpec{}, errors.New(dbRes.Error)
		}
		vss = append(vss, vs)
	}
	return &vss, nil
}

func (c *client) DeleteVolumeSnapshot(snpID string) error {
	dbReq := &Request{
		Url: GenerateUrl(prefix, "volume", "snapshots", snpID),
	}
	dbRes := c.Delete(dbReq)
	if dbRes.Status != "Success" {
		log.Println("[Error] When delete volume snapshot in db:", dbRes.Error)
		return errors.New(dbRes.Error)
	}
	return nil
}
