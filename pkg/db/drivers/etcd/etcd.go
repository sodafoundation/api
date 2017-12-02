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
This module implements the etcd database operation of data structure
defined in api module.

*/

package etcd

import (
	"encoding/json"
	"errors"
	"time"

	log "github.com/golang/glog"
	"github.com/opensds/opensds/pkg/model"
	"github.com/opensds/opensds/pkg/utils"
)

func NewClient(edps []string) *Client {
	return &Client{
		clientInterface: Init(edps),
	}
}

type Client struct {
	clientInterface
}

func (c *Client) CreateDock(dck *model.DockSpec) error {
	dckBody, err := json.Marshal(dck)
	if err != nil {
		return err
	}

	dbReq := &Request{
		Url:     GenerateDockURL(dck.GetId()),
		Content: string(dckBody),
	}
	dbRes := c.Create(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When create dock in db:", dbRes.Error)
		return errors.New(dbRes.Error)
	}

	return nil
}

func (c *Client) GetDock(dckID string) (*model.DockSpec, error) {
	dbReq := &Request{
		Url: GenerateDockURL(dckID),
	}
	dbRes := c.Get(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When get dock in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}

	var dck = &model.DockSpec{}
	if err := json.Unmarshal([]byte(dbRes.Message[0]), dck); err != nil {
		log.Error("When parsing dock in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}
	return dck, nil
}

func (c *Client) GetDockByPoolId(poolId string) (*model.DockSpec, error) {
	pool, err := c.GetPool(poolId)
	if err != nil {
		log.Error("Get pool failed in db: ", err)
		return nil, err
	}

	docks, err := c.ListDocks()
	if err != nil {
		log.Error("List docks failed failed in db: ", err)
		return nil, err
	}
	for _, dock := range docks {
		if pool.DockId == dock.Id {
			return dock, nil
		}
	}
	return nil, errors.New("Get dock failed by pool id: " + poolId)
}

func (c *Client) ListDocks() ([]*model.DockSpec, error) {
	dbReq := &Request{
		Url: GenerateDockURL(),
	}
	dbRes := c.List(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When list docks in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}

	var dcks = []*model.DockSpec{}
	if len(dbRes.Message) == 0 {
		return dcks, nil
	}
	for _, msg := range dbRes.Message {
		var dck = &model.DockSpec{}
		if err := json.Unmarshal([]byte(msg), dck); err != nil {
			log.Error("When parsing dock in db:", dbRes.Error)
			return nil, errors.New(dbRes.Error)
		}
		dcks = append(dcks, dck)
	}
	return dcks, nil
}

func (c *Client) UpdateDock(dckID, name, desp string) (*model.DockSpec, error) {
	dck, err := c.GetDock(dckID)
	if err != nil {
		return nil, err
	}
	if name != "" {
		dck.Name = name
	}
	if desp != "" {
		dck.Description = desp
	}
	dckBody, err := json.Marshal(dck)
	if err != nil {
		return nil, err
	}

	dbReq := &Request{
		Url:        GenerateDockURL(dckID),
		NewContent: string(dckBody),
	}
	dbRes := c.Update(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When update dock in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}
	return dck, nil
}

func (c *Client) DeleteDock(dckID string) error {
	dbReq := &Request{
		Url: GenerateDockURL(dckID),
	}
	dbRes := c.Delete(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When delete dock in db:", dbRes.Error)
		return errors.New(dbRes.Error)
	}
	return nil
}

func (c *Client) CreatePool(pol *model.StoragePoolSpec) error {
	polBody, err := json.Marshal(pol)
	if err != nil {
		return err
	}

	dbReq := &Request{
		Url:     GeneratePoolURL(pol.GetId()),
		Content: string(polBody),
	}
	dbRes := c.Create(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When create pol in db:", dbRes.Error)
		return errors.New(dbRes.Error)
	}

	return nil
}

func (c *Client) GetPool(polID string) (*model.StoragePoolSpec, error) {
	dbReq := &Request{
		Url: GeneratePoolURL(polID),
	}
	dbRes := c.Get(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When get pool in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}

	var pol = &model.StoragePoolSpec{}
	if err := json.Unmarshal([]byte(dbRes.Message[0]), pol); err != nil {
		log.Error("When parsing pool in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}
	return pol, nil
}

func (c *Client) ListPools() ([]*model.StoragePoolSpec, error) {
	dbReq := &Request{
		Url: GeneratePoolURL(),
	}
	dbRes := c.List(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When list pools in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}

	var pols = []*model.StoragePoolSpec{}
	if len(dbRes.Message) == 0 {
		return pols, nil
	}
	for _, msg := range dbRes.Message {
		var pol = &model.StoragePoolSpec{}
		if err := json.Unmarshal([]byte(msg), pol); err != nil {
			log.Error("When parsing pool in db:", dbRes.Error)
			return nil, errors.New(dbRes.Error)
		}
		pols = append(pols, pol)
	}
	return pols, nil
}

func (c *Client) UpdatePool(polID, name, desp string, usedCapacity int64, used bool) (*model.StoragePoolSpec, error) {
	pol, err := c.GetPool(polID)
	if err != nil {
		return nil, err
	}
	if name != "" {
		pol.Name = name
	}
	if desp != "" {
		pol.Description = desp
	}
	polBody, err := json.Marshal(pol)
	if err != nil {
		return nil, err
	}

	dbReq := &Request{
		Url:        GeneratePoolURL(polID),
		NewContent: string(polBody),
	}
	dbRes := c.Update(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When update pool in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}
	return pol, nil
}

func (c *Client) DeletePool(polID string) error {
	dbReq := &Request{
		Url: GeneratePoolURL(polID),
	}
	dbRes := c.Delete(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When delete pool in db:", dbRes.Error)
		return errors.New(dbRes.Error)
	}
	return nil
}

func (c *Client) CreateProfile(prf *model.ProfileSpec) error {
	prfBody, err := json.Marshal(prf)
	if err != nil {
		return err
	}

	dbReq := &Request{
		Url:     GenerateProfileURL(prf.GetId()),
		Content: string(prfBody),
	}
	dbRes := c.Create(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When create profile in db:", dbRes.Error)
		return errors.New(dbRes.Error)
	}

	return nil
}

func (c *Client) GetProfile(prfID string) (*model.ProfileSpec, error) {
	dbReq := &Request{
		Url: GenerateProfileURL(prfID),
	}
	dbRes := c.Get(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When get profile in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}

	var prf = &model.ProfileSpec{}
	if err := json.Unmarshal([]byte(dbRes.Message[0]), prf); err != nil {
		log.Error("When parsing profile in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}
	return prf, nil
}

func (c *Client) GetDefaultProfile() (*model.ProfileSpec, error) {
	profiles, err := c.ListProfiles()
	if err != nil {
		log.Error("Get default profile failed in db: ", err)
		return nil, err
	}

	for _, profile := range profiles {
		if profile.Name == "default" {
			return profile, nil
		}
	}
	return nil, errors.New("No default profile in db.")
}

func (c *Client) ListProfiles() ([]*model.ProfileSpec, error) {
	dbReq := &Request{
		Url: GenerateProfileURL(),
	}
	dbRes := c.List(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When list profiles in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}

	var prfs = []*model.ProfileSpec{}
	if len(dbRes.Message) == 0 {
		return prfs, nil
	}
	for _, msg := range dbRes.Message {
		var prf = &model.ProfileSpec{}
		if err := json.Unmarshal([]byte(msg), prf); err != nil {
			log.Error("When parsing profile in db:", dbRes.Error)
			return nil, errors.New(dbRes.Error)
		}
		prfs = append(prfs, prf)
	}
	return prfs, nil
}

func (c *Client) UpdateProfile(prfID string, input *model.ProfileSpec) (*model.ProfileSpec, error) {
	prf, err := c.GetProfile(prfID)
	if err != nil {
		return nil, err
	}
	if name := input.GetName(); name != "" {
		prf.Name = name
	}
	if desp := input.GetDescription(); desp != "" {
		prf.Description = desp
	}

	if props := input.Extra; len(props) != 0 {
		if prf.Extra == nil {
			prf.Extra = make(map[string]interface{})
		}
		for k, v := range props {
			prf.Extra[k] = v
		}
	}

	prf.UpdatedAt = time.Now().Format(utils.TimeFormat)

	prfBody, err := json.Marshal(prf)
	if err != nil {
		return nil, err
	}

	dbReq := &Request{
		Url:        GenerateProfileURL(prfID),
		NewContent: string(prfBody),
	}
	dbRes := c.Update(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When update profile in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}
	return prf, nil
}

func (c *Client) DeleteProfile(prfID string) error {
	dbReq := &Request{
		Url: GenerateProfileURL(prfID),
	}
	dbRes := c.Delete(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When delete profile in db:", dbRes.Error)
		return errors.New(dbRes.Error)
	}
	return nil
}

func (c *Client) AddExtraProperty(prfID string, ext model.ExtraSpec) (*model.ExtraSpec, error) {
	prf, err := c.GetProfile(prfID)
	if err != nil {
		return nil, err
	}

	if prf.Extra == nil {
		prf.Extra = make(map[string]interface{})
	}

	for k, v := range ext {
		prf.Extra[k] = v
	}

	prf.UpdatedAt = time.Now().Format(utils.TimeFormat)

	if err = c.CreateProfile(prf); err != nil {
		return nil, err
	}
	return &prf.Extra, nil
}

func (c *Client) ListExtraProperties(prfID string) (*model.ExtraSpec, error) {
	prf, err := c.GetProfile(prfID)
	if err != nil {
		return nil, err
	}
	return &prf.Extra, nil
}

func (c *Client) RemoveExtraProperty(prfID, extraKey string) error {
	prf, err := c.GetProfile(prfID)
	if err != nil {
		return err
	}

	delete(prf.Extra, extraKey)
	if err = c.CreateProfile(prf); err != nil {
		return err
	}
	return nil
}

func (c *Client) CreateVolume(vol *model.VolumeSpec) error {
	volBody, err := json.Marshal(vol)
	if err != nil {
		return err
	}

	dbReq := &Request{
		Url:     GenerateVolumeURL(vol.GetId()),
		Content: string(volBody),
	}
	dbRes := c.Create(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When create volume in db:", dbRes.Error)
		return errors.New(dbRes.Error)
	}

	return nil
}

func (c *Client) GetVolume(volID string) (*model.VolumeSpec, error) {
	dbReq := &Request{
		Url: GenerateVolumeURL(volID),
	}
	dbRes := c.Get(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When get volume in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}

	var vol = &model.VolumeSpec{}
	if err := json.Unmarshal([]byte(dbRes.Message[0]), vol); err != nil {
		log.Error("When parsing volume in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}
	return vol, nil
}

func (c *Client) ListVolumes() ([]*model.VolumeSpec, error) {
	dbReq := &Request{
		Url: GenerateVolumeURL(),
	}
	dbRes := c.List(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When list volumes in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}

	var vols = []*model.VolumeSpec{}
	if len(dbRes.Message) == 0 {
		return vols, nil
	}
	for _, msg := range dbRes.Message {
		var vol = &model.VolumeSpec{}
		if err := json.Unmarshal([]byte(msg), vol); err != nil {
			log.Error("When parsing volume in db:", dbRes.Error)
			return nil, errors.New(dbRes.Error)
		}
		vols = append(vols, vol)
	}
	return vols, nil
}

func (c *Client) DeleteVolume(volID string) error {
	dbReq := &Request{
		Url: GenerateVolumeURL(volID),
	}
	dbRes := c.Delete(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When delete volume in db:", dbRes.Error)
		return errors.New(dbRes.Error)
	}
	return nil
}

func (c *Client) CreateVolumeAttachment(attachment *model.VolumeAttachmentSpec) (*model.VolumeAttachmentSpec, error) {
	atcBody, err := json.Marshal(attachment)
	if err != nil {
		return nil, err
	}

	dbReq := &Request{
		Url:     GenerateAttachmentURL(attachment.Id),
		Content: string(atcBody),
	}
	dbRes := c.Create(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When create volume attachment in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}

	return attachment, nil
}

func (c *Client) GetVolumeAttachment(attachmentId string) (*model.VolumeAttachmentSpec, error) {
	dbReq := &Request{
		Url: GenerateAttachmentURL(attachmentId),
	}
	dbRes := c.Get(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When get volume attachment in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}

	var atc = &model.VolumeAttachmentSpec{}
	if err := json.Unmarshal([]byte(dbRes.Message[0]), atc); err != nil {
		log.Error("When parsing volume attachment in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}
	return atc, nil
}

func (c *Client) ListVolumeAttachments(volumeId string) ([]*model.VolumeAttachmentSpec, error) {
	dbReq := &Request{
		Url: GenerateAttachmentURL(),
	}
	dbRes := c.List(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When list volume attachments in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}

	var atcs = []*model.VolumeAttachmentSpec{}
	for _, msg := range dbRes.Message {
		var atc = &model.VolumeAttachmentSpec{}
		if err := json.Unmarshal([]byte(msg), atc); err != nil {
			log.Error("When parsing volume attachment in db:", dbRes.Error)
			return nil, errors.New(dbRes.Error)
		}

		if len(volumeId) == 0 || atc.Id == volumeId {
			atcs = append(atcs, atc)
		}
	}
	return atcs, nil

}

func (c *Client) UpdateVolumeAttachment(attachmentId string, attachment *model.VolumeAttachmentSpec) (*model.VolumeAttachmentSpec, error) {
	result, err := c.GetVolumeAttachment(attachmentId)
	if err != nil {
		return nil, err
	}
	if len(attachment.Mountpoint) > 0 {
		result.Mountpoint = attachment.Mountpoint
	}
	if len(attachment.Status) > 0 {
		result.Status = attachment.Status
	}
	if len(attachment.Platform) > 0 {
		result.Platform = attachment.Platform
	}
	if len(attachment.OsType) > 0 {
		result.OsType = attachment.OsType
	}
	if len(attachment.Ip) > 0 {
		result.Ip = attachment.Ip
	}
	if len(attachment.Host) > 0 {
		result.Host = attachment.Host
	}
	if len(attachment.Initiator) > 0 {
		result.Initiator = attachment.Initiator
	}
	if len(attachment.DriverVolumeType) > 0 {
		result.DriverVolumeType = attachment.DriverVolumeType
	}
	// Update metadata
	for k, v := range attachment.Metadata {
		result.Metadata[k] = v
	}
	// Update onnectionData
	for k, v := range attachment.ConnectionData {
		result.ConnectionData[k] = v
	}
	// Set update time
	result.UpdatedAt = time.Now().Format(utils.TimeFormat)

	atcBody, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	dbReq := &Request{
		Url:        GenerateAttachmentURL(attachmentId),
		NewContent: string(atcBody),
	}
	dbRes := c.Update(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When update volume attachment in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}
	return result, nil
}

func (c *Client) DeleteVolumeAttachment(attachmentId string) error {
	dbReq := &Request{
		Url: GenerateAttachmentURL(attachmentId),
	}
	dbRes := c.Delete(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When delete volume attachment in db:", dbRes.Error)
		return errors.New(dbRes.Error)
	}
	return nil
}

func (c *Client) CreateVolumeSnapshot(snp *model.VolumeSnapshotSpec) error {
	snpBody, err := json.Marshal(snp)
	if err != nil {
		return err
	}

	dbReq := &Request{
		Url:     GenerateSnapshotURL(snp.GetId()),
		Content: string(snpBody),
	}
	dbRes := c.Create(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When create volume snapshot in db:", dbRes.Error)
		return errors.New(dbRes.Error)
	}

	return nil
}

func (c *Client) GetVolumeSnapshot(snpID string) (*model.VolumeSnapshotSpec, error) {
	dbReq := &Request{
		Url: GenerateSnapshotURL(snpID),
	}
	dbRes := c.Get(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When get volume attachment in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}

	var vs = &model.VolumeSnapshotSpec{}
	if err := json.Unmarshal([]byte(dbRes.Message[0]), vs); err != nil {
		log.Error("When parsing volume snapshot in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}
	return vs, nil
}

func (c *Client) ListVolumeSnapshots() ([]*model.VolumeSnapshotSpec, error) {
	dbReq := &Request{
		Url: GenerateSnapshotURL(),
	}
	dbRes := c.List(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When list volume snapshots in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}

	var vss = []*model.VolumeSnapshotSpec{}
	if len(dbRes.Message) == 0 {
		return vss, nil
	}
	for _, msg := range dbRes.Message {
		var vs = &model.VolumeSnapshotSpec{}
		if err := json.Unmarshal([]byte(msg), vs); err != nil {
			log.Error("When parsing volume snapshot in db:", dbRes.Error)
			return nil, errors.New(dbRes.Error)
		}
		vss = append(vss, vs)
	}
	return vss, nil
}

func (c *Client) DeleteVolumeSnapshot(snpID string) error {
	dbReq := &Request{
		Url: GenerateSnapshotURL(snpID),
	}
	dbRes := c.Delete(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When delete volume snapshot in db:", dbRes.Error)
		return errors.New(dbRes.Error)
	}
	return nil
}
