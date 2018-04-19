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

/*
This module implements the etcd database operation of data structure
defined in api module.

*/

package etcd

import (
	"encoding/json"
	"errors"
	"time"

	"fmt"

	log "github.com/golang/glog"
	c "github.com/opensds/opensds/pkg/context"
	"github.com/opensds/opensds/pkg/model"
	"github.com/opensds/opensds/pkg/utils/constants"
	"github.com/opensds/opensds/pkg/utils/urls"
	"github.com/satori/go.uuid"
)

func IsAdminContext(ctx *c.Context) bool {
	return ctx.IsAdmin
}

func AuthorizeProjectContext(ctx *c.Context, tenantId string) bool {
	return ctx.TenantId == tenantId
}

// NewClient
func NewClient(edps []string) *Client {
	return &Client{
		clientInterface: Init(edps),
	}
}

// Client
type Client struct {
	clientInterface
}

// CreateDock
func (c *Client) CreateDock(ctx *c.Context, dck *model.DockSpec) (*model.DockSpec, error) {
	if dck.Id == "" {
		dck.Id = uuid.NewV4().String()
	}

	if dck.CreatedAt == "" {
		dck.CreatedAt = time.Now().Format(constants.TimeFormat)
	}

	dckBody, err := json.Marshal(dck)
	if err != nil {
		return nil, err
	}

	dbReq := &Request{
		Url:     urls.GenerateDockURL(urls.Etcd, "", dck.Id),
		Content: string(dckBody),
	}
	dbRes := c.Create(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When create dock in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}

	return dck, nil
}

// GetDock
func (c *Client) GetDock(ctx *c.Context, dckID string) (*model.DockSpec, error) {
	dbReq := &Request{
		Url: urls.GenerateDockURL(urls.Etcd, "", dckID),
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

// GetDockByPoolId
func (c *Client) GetDockByPoolId(ctx *c.Context, poolId string) (*model.DockSpec, error) {
	pool, err := c.GetPool(ctx, poolId)
	if err != nil {
		log.Error("Get pool failed in db: ", err)
		return nil, err
	}

	docks, err := c.ListDocks(ctx)
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

// ListDocks
func (c *Client) ListDocks(ctx *c.Context) ([]*model.DockSpec, error) {
	dbReq := &Request{
		Url: urls.GenerateDockURL(urls.Etcd, ""),
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

// UpdateDock
func (c *Client) UpdateDock(ctx *c.Context, dckID, name, desp string) (*model.DockSpec, error) {
	dck, err := c.GetDock(ctx, dckID)
	if err != nil {
		return nil, err
	}
	if name != "" {
		dck.Name = name
	}
	if desp != "" {
		dck.Description = desp
	}
	dck.UpdatedAt = time.Now().Format(constants.TimeFormat)

	dckBody, err := json.Marshal(dck)
	if err != nil {
		return nil, err
	}

	dbReq := &Request{
		Url:        urls.GenerateDockURL(urls.Etcd, "", dckID),
		NewContent: string(dckBody),
	}
	dbRes := c.Update(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When update dock in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}
	return dck, nil
}

// DeleteDock
func (c *Client) DeleteDock(ctx *c.Context, dckID string) error {
	dbReq := &Request{
		Url: urls.GenerateDockURL(urls.Etcd, "", dckID),
	}
	dbRes := c.Delete(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When delete dock in db:", dbRes.Error)
		return errors.New(dbRes.Error)
	}
	return nil
}

// CreatePool
func (c *Client) CreatePool(ctx *c.Context, pol *model.StoragePoolSpec) (*model.StoragePoolSpec, error) {
	if pol.Id == "" {
		pol.Id = uuid.NewV4().String()
	}

	if pol.CreatedAt == "" {
		pol.CreatedAt = time.Now().Format(constants.TimeFormat)
	}
	polBody, err := json.Marshal(pol)
	if err != nil {
		return nil, err
	}

	dbReq := &Request{
		Url:     urls.GeneratePoolURL(urls.Etcd, "", pol.Id),
		Content: string(polBody),
	}
	dbRes := c.Create(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When create pol in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}

	return pol, nil
}

// GetPool
func (c *Client) GetPool(ctx *c.Context, polID string) (*model.StoragePoolSpec, error) {
	dbReq := &Request{
		Url: urls.GeneratePoolURL(urls.Etcd, "", polID),
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

// ListPools
func (c *Client) ListPools(ctx *c.Context) ([]*model.StoragePoolSpec, error) {
	dbReq := &Request{
		Url: urls.GeneratePoolURL(urls.Etcd, ""),
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

// UpdatePool
func (c *Client) UpdatePool(ctx *c.Context, polID, name, desp string, usedCapacity int64, used bool) (*model.StoragePoolSpec, error) {
	pol, err := c.GetPool(ctx, polID)
	if err != nil {
		return nil, err
	}
	if name != "" {
		pol.Name = name
	}
	if desp != "" {
		pol.Description = desp
	}
	pol.UpdatedAt = time.Now().Format(constants.TimeFormat)

	polBody, err := json.Marshal(pol)
	if err != nil {
		return nil, err
	}

	dbReq := &Request{
		Url:        urls.GeneratePoolURL(urls.Etcd, "", polID),
		NewContent: string(polBody),
	}
	dbRes := c.Update(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When update pool in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}
	return pol, nil
}

// DeletePool
func (c *Client) DeletePool(ctx *c.Context, polID string) error {
	dbReq := &Request{
		Url: urls.GeneratePoolURL(urls.Etcd, "", polID),
	}
	dbRes := c.Delete(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When delete pool in db:", dbRes.Error)
		return errors.New(dbRes.Error)
	}
	return nil
}

// CreateProfile
func (c *Client) CreateProfile(ctx *c.Context, prf *model.ProfileSpec) (*model.ProfileSpec, error) {

	if prf.Id == "" {
		prf.Id = uuid.NewV4().String()
	}

	if prf.CreatedAt == "" {
		prf.CreatedAt = time.Now().Format(constants.TimeFormat)
	}
	prf.TenantId = ctx.TenantId
	prfBody, err := json.Marshal(prf)
	if err != nil {
		return nil, err
	}

	dbReq := &Request{
		Url:     urls.GenerateProfileURL(urls.Etcd, "", prf.Id),
		Content: string(prfBody),
	}
	dbRes := c.Create(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When create profile in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}

	return prf, nil
}

// GetProfile
func (c *Client) GetProfile(ctx *c.Context, prfID string) (*model.ProfileSpec, error) {
	dbReq := &Request{
		Url: urls.GenerateProfileURL(urls.Etcd, "", prfID),
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

// GetDefaultProfile
func (c *Client) GetDefaultProfile(ctx *c.Context) (*model.ProfileSpec, error) {
	profiles, err := c.ListProfiles(ctx)
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

// ListProfiles
func (c *Client) ListProfiles(ctx *c.Context) ([]*model.ProfileSpec, error) {
	dbReq := &Request{
		Url: urls.GenerateProfileURL(urls.Etcd, ""),
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

// UpdateProfile
func (c *Client) UpdateProfile(ctx *c.Context, prfID string, input *model.ProfileSpec) (*model.ProfileSpec, error) {
	prf, err := c.GetProfile(ctx, prfID)
	if err != nil {
		return nil, err
	}
	if name := input.Name; name != "" {
		prf.Name = name
	}
	if desp := input.Description; desp != "" {
		prf.Description = desp
	}
	prf.UpdatedAt = time.Now().Format(constants.TimeFormat)

	if props := input.Extras; len(props) != 0 {
		if prf.Extras == nil {
			prf.Extras = make(map[string]interface{})
		}
		for k, v := range props {
			prf.Extras[k] = v
		}
	}

	prf.UpdatedAt = time.Now().Format(constants.TimeFormat)

	prfBody, err := json.Marshal(prf)
	if err != nil {
		return nil, err
	}

	dbReq := &Request{
		Url:        urls.GenerateProfileURL(urls.Etcd, "", prfID),
		NewContent: string(prfBody),
	}
	dbRes := c.Update(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When update profile in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}
	return prf, nil
}

// DeleteProfile
func (c *Client) DeleteProfile(ctx *c.Context, prfID string) error {
	dbReq := &Request{
		Url: urls.GenerateProfileURL(urls.Etcd, "", prfID),
	}
	dbRes := c.Delete(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When delete profile in db:", dbRes.Error)
		return errors.New(dbRes.Error)
	}
	return nil
}

// AddExtraProperty
func (c *Client) AddExtraProperty(ctx *c.Context, prfID string, ext model.ExtraSpec) (*model.ExtraSpec, error) {
	prf, err := c.GetProfile(ctx, prfID)
	if err != nil {
		return nil, err
	}

	if prf.Extras == nil {
		prf.Extras = make(map[string]interface{})
	}

	for k, v := range ext {
		prf.Extras[k] = v
	}

	prf.UpdatedAt = time.Now().Format(constants.TimeFormat)

	if _, err = c.CreateProfile(nil, prf); err != nil {
		return nil, err
	}
	return &prf.Extras, nil
}

// ListExtraProperties
func (c *Client) ListExtraProperties(ctx *c.Context, prfID string) (*model.ExtraSpec, error) {
	prf, err := c.GetProfile(ctx, prfID)
	if err != nil {
		return nil, err
	}
	return &prf.Extras, nil
}

// RemoveExtraProperty
func (c *Client) RemoveExtraProperty(ctx *c.Context, prfID, extraKey string) error {
	prf, err := c.GetProfile(ctx, prfID)
	if err != nil {
		return err
	}

	delete(prf.Extras, extraKey)
	if _, err = c.CreateProfile(nil, prf); err != nil {
		return err
	}
	return nil
}

// CreateVolume
func (c *Client) CreateVolume(ctx *c.Context, vol *model.VolumeSpec) (*model.VolumeSpec, error) {
	vol.TenantId = ctx.TenantId
	volBody, err := json.Marshal(vol)
	if err != nil {
		return nil, err
	}

	dbReq := &Request{
		Url:     urls.GenerateVolumeURL(urls.Etcd, ctx.TenantId, vol.Id),
		Content: string(volBody),
	}
	dbRes := c.Create(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When create volume in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}

	return vol, nil
}

// GetVolume
func (c *Client) GetVolume(ctx *c.Context, volID string) (*model.VolumeSpec, error) {
	vol, err := c.getVolume(ctx, volID)
	if !IsAdminContext(ctx) || err == nil {
		return vol, err
	}
	vols, err := c.ListVolumes(ctx)
	if err != nil {
		return nil, err
	}
	for _, v := range vols {
		if v.Id == volID {
			return v, nil
		}
	}
	return nil, fmt.Errorf("specified volume(%s) can't find", volID)
}

func (c *Client) getVolume(ctx *c.Context, volID string) (*model.VolumeSpec, error) {
	dbReq := &Request{
		Url: urls.GenerateVolumeURL(urls.Etcd, ctx.TenantId, volID),
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

// ListVolumes
func (c *Client) ListVolumes(ctx *c.Context) ([]*model.VolumeSpec, error) {
	dbReq := &Request{
		Url: urls.GenerateVolumeURL(urls.Etcd, ctx.TenantId),
	}

	// Admin user should get all volumes including the volumes whose tenant is not admin.
	if IsAdminContext(ctx) {
		dbReq.Url = urls.GenerateVolumeURL(urls.Etcd, "")
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

// UpdateVolume ...
func (c *Client) UpdateVolume(ctx *c.Context, vol *model.VolumeSpec) (*model.VolumeSpec, error) {
	result, err := c.GetVolume(ctx, vol.Id)
	if err != nil {
		return nil, err
	}
	if vol.Name != "" {
		result.Name = vol.Name
	}
	if vol.AvailabilityZone != "" {
		result.AvailabilityZone = vol.AvailabilityZone
	}
	if vol.Description != "" {
		result.Description = vol.Description
	}
	if vol.Metadata != nil {
		result.Metadata = vol.Metadata
	}
	if vol.PoolId != "" {
		result.PoolId = vol.PoolId
	}
	if vol.ProfileId != "" {
		result.ProfileId = vol.ProfileId
	}
	if vol.Size != 0 {
		result.Size = vol.Size
	}
	if vol.Status != "" {
		result.Status = vol.Status
	}

	result.GroupId = vol.GroupId

	// Set update time
	result.UpdatedAt = time.Now().Format(constants.TimeFormat)

	body, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	// If an admin want to access other tenant's resource just fake other's tenantId.
	if !IsAdminContext(ctx) && !AuthorizeProjectContext(ctx, result.TenantId) {
		return nil, fmt.Errorf("opertaion is not permitted")
	}

	dbReq := &Request{
		Url:        urls.GenerateVolumeURL(urls.Etcd, result.TenantId, vol.Id),
		NewContent: string(body),
	}

	dbRes := c.Update(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When update volume in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}
	return result, nil
}

func (c *Client) VolumesToUpdate(ctx *c.Context, volumeList []*model.VolumeSpec) ([]*model.VolumeSpec, error) {
	var volumeRefs []*model.VolumeSpec
	for _, values := range volumeList {
		v, err := c.UpdateVolume(ctx, values)
		if err != nil {
			return nil, err
		}
		volumeRefs = append(volumeRefs, v)
	}
	return volumeRefs, nil
}

// DeleteVolume
func (c *Client) DeleteVolume(ctx *c.Context, volID string) error {
	// If an admin want to access other tenant's resource just fake other's tenantId.
	tenantId := ctx.TenantId
	if IsAdminContext(ctx) {
		vol, err := c.GetVolume(ctx, volID)
		if err != nil {
			log.Error(err)
			return err
		}
		tenantId = vol.TenantId
	}
	dbReq := &Request{
		Url: urls.GenerateVolumeURL(urls.Etcd, tenantId, volID),
	}

	dbRes := c.Delete(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When delete volume in db:", dbRes.Error)
		return errors.New(dbRes.Error)
	}
	return nil
}

// ExtendVolume ...
func (c *Client) ExtendVolume(ctx *c.Context, vol *model.VolumeSpec) (*model.VolumeSpec, error) {
	result, err := c.GetVolume(ctx, vol.Id)
	if err != nil {
		return nil, err
	}

	if vol.Size > 0 {
		result.Size = vol.Size
	}
	result.Status = vol.Status
	// Set update time
	result.UpdatedAt = time.Now().Format(constants.TimeFormat)

	body, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	dbReq := &Request{
		Url:        urls.GenerateVolumeURL(urls.Etcd, ctx.TenantId, vol.Id),
		NewContent: string(body),
	}

	dbRes := c.Update(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When extend volume in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}
	return result, nil
}

// CreateVolumeAttachment
func (c *Client) CreateVolumeAttachment(ctx *c.Context, attachment *model.VolumeAttachmentSpec) (*model.VolumeAttachmentSpec, error) {
	attachment.TenantId = ctx.TenantId

	atcBody, err := json.Marshal(attachment)
	if err != nil {
		return nil, err
	}
	dbReq := &Request{
		Url:     urls.GenerateAttachmentURL(urls.Etcd, ctx.TenantId, attachment.Id),
		Content: string(atcBody),
	}
	dbRes := c.Create(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When create volume attachment in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}

	return attachment, nil
}

func (c *Client) GetVolumeAttachment(ctx *c.Context, attachmentId string) (*model.VolumeAttachmentSpec, error) {
	attach, err := c.getVolumeAttachment(ctx, attachmentId)
	if !IsAdminContext(ctx) || err == nil {
		return attach, err
	}
	attachs, err := c.ListVolumeAttachments(ctx, "")
	if err != nil {
		return nil, err
	}
	for _, v := range attachs {
		if v.Id == attachmentId {
			return v, nil
		}
	}
	return nil, fmt.Errorf("specified volume attachment(%s) can't find", attachmentId)
}

// GetVolumeAttachment
func (c *Client) getVolumeAttachment(ctx *c.Context, attachmentId string) (*model.VolumeAttachmentSpec, error) {
	dbReq := &Request{
		Url: urls.GenerateAttachmentURL(urls.Etcd, ctx.TenantId, attachmentId),
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

// ListVolumeAttachments
func (c *Client) ListVolumeAttachments(ctx *c.Context, volumeId string) ([]*model.VolumeAttachmentSpec, error) {
	dbReq := &Request{
		Url: urls.GenerateAttachmentURL(urls.Etcd, ctx.TenantId),
	}
	if IsAdminContext(ctx) {
		dbReq.Url = urls.GenerateAttachmentURL(urls.Etcd, "")
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

		if len(volumeId) == 0 || atc.VolumeId == volumeId {
			atcs = append(atcs, atc)
		}
	}
	return atcs, nil

}

// UpdateVolumeAttachment
func (c *Client) UpdateVolumeAttachment(ctx *c.Context, attachmentId string, attachment *model.VolumeAttachmentSpec) (*model.VolumeAttachmentSpec, error) {
	result, err := c.GetVolumeAttachment(ctx, attachmentId)
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
	if result.Metadata == nil {
		result.Metadata = make(map[string]string)
	}
	for k, v := range attachment.Metadata {
		result.Metadata[k] = v
	}
	// Update onnectionData
	for k, v := range attachment.ConnectionData {
		result.ConnectionData[k] = v
	}
	// Set update time
	result.UpdatedAt = time.Now().Format(constants.TimeFormat)

	atcBody, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	// If an admin want to access other tenant's resource just fake other's tenantId.
	if !IsAdminContext(ctx) && !AuthorizeProjectContext(ctx, result.TenantId) {
		return nil, fmt.Errorf("opertaion is not permitted")
	}

	dbReq := &Request{
		Url:        urls.GenerateAttachmentURL(urls.Etcd, result.TenantId, attachmentId),
		NewContent: string(atcBody),
	}

	dbRes := c.Update(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When update volume attachment in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}
	return result, nil
}

// DeleteVolumeAttachment
func (c *Client) DeleteVolumeAttachment(ctx *c.Context, attachmentId string) error {

	// If an admin want to access other tenant's resource just fake other's tenantId.
	tenantId := ctx.TenantId
	if IsAdminContext(ctx) {
		attach, err := c.GetVolumeAttachment(ctx, attachmentId)
		if err != nil {
			log.Error(err)
			return err
		}
		tenantId = attach.TenantId
	}
	dbReq := &Request{
		Url: urls.GenerateAttachmentURL(urls.Etcd, tenantId, attachmentId),
	}

	dbRes := c.Delete(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When delete volume attachment in db:", dbRes.Error)
		return errors.New(dbRes.Error)
	}
	return nil
}

// CreateVolumeSnapshot
func (c *Client) CreateVolumeSnapshot(ctx *c.Context, snp *model.VolumeSnapshotSpec) (*model.VolumeSnapshotSpec, error) {
	snp.TenantId = ctx.TenantId
	snpBody, err := json.Marshal(snp)
	if err != nil {
		return nil, err
	}

	dbReq := &Request{
		Url:     urls.GenerateSnapshotURL(urls.Etcd, ctx.TenantId, snp.Id),
		Content: string(snpBody),
	}
	dbRes := c.Create(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When create volume snapshot in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}

	return snp, nil
}

func (c *Client) GetVolumeSnapshot(ctx *c.Context, snpID string) (*model.VolumeSnapshotSpec, error) {
	snap, err := c.getVolumeSnapshot(ctx, snpID)
	if !IsAdminContext(ctx) || err == nil {
		return snap, err
	}
	snaps, err := c.ListVolumeSnapshots(ctx)
	if err != nil {
		return nil, err
	}
	for _, v := range snaps {
		if v.Id == snpID {
			return v, nil
		}
	}
	return nil, fmt.Errorf("specified volume snapshot(%s) can't find", snpID)
}

// GetVolumeSnapshot
func (c *Client) getVolumeSnapshot(ctx *c.Context, snpID string) (*model.VolumeSnapshotSpec, error) {
	dbReq := &Request{
		Url: urls.GenerateSnapshotURL(urls.Etcd, ctx.TenantId, snpID),
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

// ListVolumeSnapshots
func (c *Client) ListVolumeSnapshots(ctx *c.Context) ([]*model.VolumeSnapshotSpec, error) {
	dbReq := &Request{
		Url: urls.GenerateSnapshotURL(urls.Etcd, ctx.TenantId),
	}
	if IsAdminContext(ctx) {
		dbReq.Url = urls.GenerateSnapshotURL(urls.Etcd, "")
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

// UpdateVolumeSnapshot
func (c *Client) UpdateVolumeSnapshot(ctx *c.Context, snpID string, snp *model.VolumeSnapshotSpec) (*model.VolumeSnapshotSpec, error) {
	result, err := c.GetVolumeSnapshot(ctx, snpID)
	if err != nil {
		return nil, err
	}
	if snp.Name != "" {
		result.Name = snp.Name
	}
	if snp.Metadata != nil {
		result.Metadata = snp.Metadata
	}
	if snp.Size > 0 {
		result.Size = snp.Size
	}
	if snp.VolumeId != "" {
		result.VolumeId = snp.VolumeId
	}
	if snp.Description != "" {
		result.Description = snp.Description
	}
	if snp.Status != "" {
		result.Status = snp.Status
	}
	// Set update time
	result.UpdatedAt = time.Now().Format(constants.TimeFormat)

	atcBody, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	// If an admin want to access other tenant's resource just fake other's tenantId.
	if !IsAdminContext(ctx) && !AuthorizeProjectContext(ctx, result.TenantId) {
		return nil, fmt.Errorf("opertaion is not permitted")
	}

	dbReq := &Request{
		Url:        urls.GenerateSnapshotURL(urls.Etcd, result.TenantId, snpID),
		NewContent: string(atcBody),
	}

	dbRes := c.Update(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When update volume snapshot in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}
	return result, nil
}

// DeleteVolumeSnapshot
func (c *Client) DeleteVolumeSnapshot(ctx *c.Context, snpID string) error {
	// If an admin want to access other tenant's resource just fake other's tenantId.
	tenantId := ctx.TenantId
	if IsAdminContext(ctx) {
		snap, err := c.GetVolumeSnapshot(ctx, snpID)
		if err != nil {
			log.Error(err)
			return err
		}
		tenantId = snap.TenantId
	}
	dbReq := &Request{
		Url: urls.GenerateSnapshotURL(urls.Etcd, tenantId, snpID),
	}

	dbRes := c.Delete(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When delete volume snapshot in db:", dbRes.Error)
		return errors.New(dbRes.Error)
	}
	return nil
}

func (c *Client) CreateVolumeGroup(ctx *c.Context, vg *model.VolumeGroupSpec) (*model.VolumeGroupSpec, error) {
	vgBody, err := json.Marshal(vg)
	if err != nil {
		return nil, err
	}

	dbReq := &Request{
		Url:     urls.GenerateVolumeGroupURL(urls.Etcd, ctx.TenantId, vg.Id),
		Content: string(vgBody),
	}
	dbRes := c.Create(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When create volume group in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}

	return vg, nil
}

func (c *Client) GetVolumeGroup(ctx *c.Context, vgId string) (*model.VolumeGroupSpec, error) {
	dbReq := &Request{
		Url: urls.GenerateVolumeGroupURL(urls.Etcd, ctx.TenantId, vgId),
	}
	dbRes := c.Get(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When get volume group in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}

	var vg = &model.VolumeGroupSpec{}
	if err := json.Unmarshal([]byte(dbRes.Message[0]), vg); err != nil {
		log.Error("When parsing volume group in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}
	return vg, nil
}

func (c *Client) UpdateVolumeGroup(ctx *c.Context, volumeGroup *model.VolumeGroupSpec) (*model.VolumeGroupSpec, error) {
	group, err := c.GetVolumeGroup(ctx, volumeGroup.Id)
	if err != nil {
		return nil, err
	}
	if volumeGroup.Name != "" && volumeGroup.Name != group.Name {
		group.Name = volumeGroup.Name
	}
	if volumeGroup.AvailabilityZone != "" && volumeGroup.AvailabilityZone != group.AvailabilityZone {
		group.AvailabilityZone = volumeGroup.AvailabilityZone
	}
	if volumeGroup.Description != "" && volumeGroup.Description != group.Description {
		group.Description = volumeGroup.Description
	}
	if volumeGroup.PoolId != "" && volumeGroup.PoolId != group.PoolId {
		group.PoolId = volumeGroup.PoolId
	}
	if volumeGroup.Status != "" && volumeGroup.Status != group.Status {
		group.Status = volumeGroup.Status
	}
	if volumeGroup.PoolId != "" && volumeGroup.PoolId != group.PoolId {
		group.PoolId = volumeGroup.PoolId
	}
	if volumeGroup.CreatedAt != "" && volumeGroup.CreatedAt != group.CreatedAt {
		group.CreatedAt = volumeGroup.CreatedAt
	}
	if volumeGroup.UpdatedAt != "" && volumeGroup.UpdatedAt != group.UpdatedAt {
		group.UpdatedAt = volumeGroup.UpdatedAt
	}

	vgBody, err := json.Marshal(group)
	if err != nil {
		return nil, err
	}

	dbReq := &Request{
		Url:        urls.GenerateVolumeGroupURL(urls.Etcd, ctx.TenantId, volumeGroup.Id),
		NewContent: string(vgBody),
	}
	dbRes := c.Update(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When update volume group in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}
	return group, nil
}

func (c *Client) UpdateStatus(ctx *c.Context, in interface{}, status string) error {
	switch in.(type) {

	case *model.VolumeSnapshotSpec:
		snap := in.(*model.VolumeSnapshotSpec)
		snap.Status = status
		if _, errUpdate := c.UpdateVolumeSnapshot(ctx, snap.Id, snap); errUpdate != nil {
			log.Error("Error occurs when update volume snapshot status in db:", errUpdate.Error())
			return errUpdate
		}

	case *model.VolumeAttachmentSpec:
		attm := in.(*model.VolumeAttachmentSpec)
		attm.Status = status
		if _, errUpdate := c.UpdateVolumeAttachment(ctx, attm.Id, attm); errUpdate != nil {
			log.Error("Error occurred in dock module when update volume attachment status in db:", errUpdate)
			return errUpdate
		}

	case *model.VolumeSpec:
		volume := in.(*model.VolumeSpec)
		volume.Status = status
		if _, errUpdate := c.UpdateVolume(ctx, volume); errUpdate != nil {
			log.Error("When update volume status in db:", errUpdate.Error())
			return errUpdate
		}

	case *model.VolumeGroupSpec:
		vg := in.(*model.VolumeGroupSpec)
		vg.Status = status
		if _, errUpdate := c.UpdateVolumeGroup(ctx, vg); errUpdate != nil {
			log.Error("When update volume status in db:", errUpdate.Error())
			return errUpdate
		}

	case []*model.VolumeSpec:
		vols := in.([]*model.VolumeSpec)
		if _, errUpdate := c.VolumesToUpdate(ctx, vols); errUpdate != nil {
			return errUpdate
		}
	}
	return nil
}

func (c *Client) ListVolumesByGroupId(ctx *c.Context, groupId string) ([]*model.VolumeSpec, error) {
	volumes, err := c.ListVolumes(ctx)

	if err != nil {
		return nil, err
	}

	var volumesInSameGroup []*model.VolumeSpec
	for _, v := range volumes {
		if v.GroupId == groupId {
			volumesInSameGroup = append(volumesInSameGroup, v)
		}
	}

	return volumesInSameGroup, nil
}

// ListVolumes
func (c *Client) ListVolumeGroups(ctx *c.Context) ([]*model.VolumeGroupSpec, error) {
	dbReq := &Request{
		Url: urls.GenerateVolumeGroupURL(urls.Etcd, ctx.TenantId),
	}

	dbRes := c.List(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When list volume groups in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}

	var groups []*model.VolumeGroupSpec
	if len(dbRes.Message) == 0 {
		return groups, nil
	}
	for _, msg := range dbRes.Message {
		var group = &model.VolumeGroupSpec{}
		if err := json.Unmarshal([]byte(msg), group); err != nil {
			log.Error("When parsing volume group in db:", dbRes.Error)
			return nil, errors.New(dbRes.Error)
		}
		groups = append(groups, group)
	}
	return groups, nil
}

func (c *Client) DeleteVolumeGroup(ctx *c.Context, volumeGroupId string) error {
	// If an admin want to access other tenant's resource just fake other's tenantId.
	tenantId := ctx.TenantId
	if IsAdminContext(ctx) {
		group, err := c.GetVolumeGroup(ctx, volumeGroupId)
		if err != nil {
			log.Error(err)
			return err
		}
		tenantId = group.TenantId
	}
	dbReq := &Request{
		Url: urls.GenerateVolumeGroupURL(urls.Etcd, tenantId, volumeGroupId),
	}

	dbRes := c.Delete(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When delete volume group in db:", dbRes.Error)
		return errors.New(dbRes.Error)
	}
	return nil
}

func (c *Client) ListSnapshotsByVolumeId(ctx *c.Context, volumeId string) ([]*model.VolumeSnapshotSpec, error) {
	snaps, err := c.ListVolumeSnapshots(ctx)
	if err != nil {
		return nil, err
	}

	var snapList []*model.VolumeSnapshotSpec
	for _, snap := range snaps {
		if snap.VolumeId == volumeId {
			snapList = append(snapList, snap)
		}
	}
	return snapList, nil
}
