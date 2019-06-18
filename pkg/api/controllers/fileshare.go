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

package controllers

import (
	"context"
	"encoding/json"
	"fmt"

	log "github.com/golang/glog"
	"github.com/opensds/opensds/pkg/api/policy"
	"github.com/opensds/opensds/pkg/api/util"
	c "github.com/opensds/opensds/pkg/context"
	"github.com/opensds/opensds/pkg/controller/client"
	"github.com/opensds/opensds/pkg/db"
	"github.com/opensds/opensds/pkg/model"
	pb "github.com/opensds/opensds/pkg/model/proto"
	"github.com/opensds/opensds/pkg/utils"
	. "github.com/opensds/opensds/pkg/utils/config"
	"github.com/opensds/opensds/pkg/utils/constants"
)

func NewFileSharePortal() *FileSharePortal {
	return &FileSharePortal{
		CtrClient: client.NewClient(),
	}
}

type FileSharePortal struct {
	BasePortal

	CtrClient client.Client
}

// Function to store Acl's related entry into databse
func (f *FileSharePortal) CreateFileShareAcl() {
	if !policy.Authorize(f.Ctx, "fileshareacl:create") {
		return
	}
	ctx := c.GetContext(f.Ctx)
	// Get profile
	var prf *model.ProfileSpec
	var fileshareacl = model.FileShareAclSpec{
		BaseModel: &model.BaseModel{},
	}
	// Unmarshal the request body
	if err := json.NewDecoder(f.Ctx.Request.Body).Decode(&fileshareacl); err != nil {
		reason := fmt.Sprintf("parse fileshare access rules request body failed: %s", err.Error())
		f.ErrorHandle(model.ErrorBadRequest, reason)
		log.Error(reason)
		return
	}
	result, err := util.CreateFileShareAclDBEntry(c.GetContext(f.Ctx), &fileshareacl)
	if err != nil {
		reason := fmt.Sprintf("createFileshareAcldbentry failed: %s", err.Error())
		f.ErrorHandle(model.ErrorBadRequest, reason)
		log.Error(reason)
		return
	}
	fileshare, err := db.C.GetFileShare(ctx, result.FileShareId)
	if err != nil {
		reason := fmt.Sprintf("getFileshare failed in createfileshare acl: %s", err.Error())
		f.ErrorHandle(model.ErrorBadRequest, reason)
		log.Error(reason)
		return
	}
	prf, err = db.C.GetProfile(ctx, fileshare.ProfileId)
	if err != nil {
		errMsg := fmt.Sprintf("get profile failed: %s", err.Error())
		f.ErrorHandle(model.ErrorBadRequest, errMsg)
		return
	}
	// Marshal the result.
	body, err := json.Marshal(result)
	if err != nil {
		reason := fmt.Sprintf("marshal fileshare access rules created result failed: %s", err.Error())
		f.ErrorHandle(model.ErrorBadRequest, reason)
		log.Error(reason)
		return
	}
	f.SuccessHandle(StatusAccepted, body)

	// FileShare acl access creation request is sent to dock and drivers
	if err := f.CtrClient.Connect(CONF.OsdsLet.ApiEndpoint); err != nil {
		log.Error("when connecting controller client:", err)
		return
	}
	defer f.CtrClient.Close()

	opt := &pb.CreateFileShareAclOpts{
		Id:               result.Id,
		FileshareId:      result.FileShareId,
		Description:      result.Description,
		Type:             result.Type,
		AccessCapability: result.AccessCapability,
		AccessTo:         result.AccessTo,
		Metadata:         fileshare.Metadata,
		Context:          ctx.ToJson(),
		Profile:          prf.ToJson(),
	}
	response, err := f.CtrClient.CreateFileShareAcl(context.Background(), opt)
	if err != nil {
		log.Error("create file share acl failed in controller service:", err)
		return
	}
	if errorMsg := response.GetError(); errorMsg != nil {
		log.Errorf("failed to create file share acl in controller, code: %v, message: %v",
			errorMsg.GetCode(), errorMsg.GetDescription())
		return
	}

	return
}

func (f *FileSharePortal) ListFileSharesAcl() {
	if !policy.Authorize(f.Ctx, "fileshareacl:list") {
		return
	}
	m, err := f.GetParameters()
	if err != nil {
		errMsg := fmt.Sprintf("list fileshares failed: %s", err.Error())
		f.ErrorHandle(model.ErrorBadRequest, errMsg)
		return
	}
	result, err := db.C.ListFileSharesAclWithFilter(c.GetContext(f.Ctx), m)
	if err != nil {
		errMsg := fmt.Sprintf("list fileshares failed: %s", err.Error())
		f.ErrorHandle(model.ErrorInternalServer, errMsg)
		return
	}
	// Marshal the result.
	body, _ := json.Marshal(result)
	f.SuccessHandle(StatusOK, body)

	return
}

// Function to store fileshare related entry into databse
func (f *FileSharePortal) CreateFileShare() {
	if !policy.Authorize(f.Ctx, "fileshare:create") {
		return
	}
	ctx := c.GetContext(f.Ctx)
	var fileshare = model.FileShareSpec{
		BaseModel: &model.BaseModel{},
	}

	// Unmarshal the request body
	if err := json.NewDecoder(f.Ctx.Request.Body).Decode(&fileshare); err != nil {
		reason := fmt.Sprintf("parse fileshare request body failed: %s", err.Error())
		f.ErrorHandle(model.ErrorBadRequest, reason)
		log.Error(reason)
		return
	}

	// Get profile
	var prf *model.ProfileSpec
	var err error
	if fileshare.ProfileId == "" {
		log.Warning("Use default profile when user doesn't specify profile.")
		prf, err = db.C.GetDefaultProfileFileShare(ctx)
		if err != nil {
			errMsg := fmt.Sprintf("get profile failed: %s", err.Error())
			f.ErrorHandle(model.ErrorBadRequest, errMsg)
			return
		}
		fileshare.ProfileId = prf.Id
	} else {
		prf, err = db.C.GetProfile(ctx, fileshare.ProfileId)
		if err != nil {
			errMsg := fmt.Sprintf("get profile failed: %s", err.Error())
			f.ErrorHandle(model.ErrorBadRequest, errMsg)
			return
		}
		if prf.StorageType != constants.File {
			errMsg := fmt.Sprintf("storageType should be only file. Currently it is: %s", prf.StorageType)
			log.Error(errMsg)
			f.ErrorHandle(model.ErrorBadRequest, errMsg)
			return
		}
	}

	// NOTE: It will create a file share entry into the database and initialize its status
	// as "creating". It will not wait for the real file share creation to complete
	// and will return result immediately.
	result, err := util.CreateFileShareDBEntry(c.GetContext(f.Ctx), &fileshare)
	if err != nil {
		reason := fmt.Sprintf("create fileshare failed: %s", err.Error())
		f.ErrorHandle(model.ErrorBadRequest, reason)
		log.Error(reason)
		return
	}
	// Marshal the result.
	body, err := json.Marshal(result)
	if err != nil {
		reason := fmt.Sprintf("marshal fileshare created result failed: %s", err.Error())
		f.ErrorHandle(model.ErrorBadRequest, reason)
		log.Error(reason)
		return
	}
	f.SuccessHandle(StatusAccepted, body)

	// NOTE: The real file share creation process.
	// FileShare creation request is sent to the Dock. Dock will update file share status to "available"
	// after file share creation is completed.
	if err := f.CtrClient.Connect(CONF.OsdsLet.ApiEndpoint); err != nil {
		log.Error("when connecting controller client:", err)
		return
	}
	defer f.CtrClient.Close()

	opt := &pb.CreateFileShareOpts{
		Id:               result.Id,
		Name:             result.Name,
		Description:      result.Description,
		Size:             result.Size,
		AvailabilityZone: result.AvailabilityZone,
		Profile:          prf.ToJson(),
		PoolId:           result.PoolId,
		ExportLocations:  result.ExportLocations,
		Metadata:         result.Metadata,
		Context:          ctx.ToJson(),
	}
	response, err := f.CtrClient.CreateFileShare(context.Background(), opt)
	if err != nil {
		log.Error("create file share failed in controller service:", err)
		return
	}
	if errorMsg := response.GetError(); errorMsg != nil {
		log.Errorf("failed to create file share in controller, code: %v, message: %v",
			errorMsg.GetCode(), errorMsg.GetDescription())
		return
	}

	return
}

func (f *FileSharePortal) ListFileShares() {
	if !policy.Authorize(f.Ctx, "fileshare:list") {
		return
	}
	m, err := f.GetParameters()
	if err != nil {
		errMsg := fmt.Sprintf("list fileshares failed: %s", err.Error())
		f.ErrorHandle(model.ErrorBadRequest, errMsg)
		return
	}
	result, err := db.C.ListFileSharesWithFilter(c.GetContext(f.Ctx), m)
	if err != nil {
		errMsg := fmt.Sprintf("list fileshares failed: %s", err.Error())
		f.ErrorHandle(model.ErrorInternalServer, errMsg)
		return
	}
	// Marshal the result.
	body, _ := json.Marshal(result)
	f.SuccessHandle(StatusOK, body)

	return
}

func (f *FileSharePortal) GetFileShareAcl() {
	if !policy.Authorize(f.Ctx, "fileshareacl:get") {
		return
	}
	id := f.Ctx.Input.Param(":aclId")

	// Call db api module to handle get fileshare request.
	result, err := db.C.GetFileShareAcl(c.GetContext(f.Ctx), id)
	if err != nil {
		errMsg := fmt.Sprintf("fileshare acl %s not found: %s", id, err.Error())
		f.ErrorHandle(model.ErrorNotFound, errMsg)
		return
	}

	// Marshal the result.
	body, _ := json.Marshal(result)
	if err != nil {
		reason := fmt.Sprintf("marshal fileshare acl list result failed: %s", err.Error())
		f.ErrorHandle(model.ErrorBadRequest, reason)
		log.Error(reason)
		return
	}
	f.SuccessHandle(StatusOK, body)

	return
}

func (f *FileSharePortal) GetFileShare() {
	if !policy.Authorize(f.Ctx, "fileshare:get") {
		return
	}
	id := f.Ctx.Input.Param(":fileshareId")

	// Call db api module to handle get file share request.
	result, err := db.C.GetFileShare(c.GetContext(f.Ctx), id)
	if err != nil {
		errMsg := fmt.Sprintf("fileshare %s not found: %s", id, err.Error())
		f.ErrorHandle(model.ErrorNotFound, errMsg)
		return
	}

	// Marshal the result.
	body, _ := json.Marshal(result)
	if err != nil {
		reason := fmt.Sprintf("marshal fileshare list result failed: %s", err.Error())
		f.ErrorHandle(model.ErrorBadRequest, reason)
		log.Error(reason)
		return
	}
	f.SuccessHandle(StatusOK, body)

	return
}

func (f *FileSharePortal) UpdateFileShare() {
	if !policy.Authorize(f.Ctx, "fileshare:update") {
		return
	}
	var fshare = model.FileShareSpec{
		BaseModel: &model.BaseModel{},
	}

	id := f.Ctx.Input.Param(":fileshareId")
	if err := json.NewDecoder(f.Ctx.Request.Body).Decode(&fshare); err != nil {
		errMsg := fmt.Sprintf("parse fileshare request body failed: %s", err.Error())
		f.ErrorHandle(model.ErrorBadRequest, errMsg)
		return
	}

	fshare.Id = id
	result, err := db.C.UpdateFileShare(c.GetContext(f.Ctx), &fshare)
	if err != nil {
		errMsg := fmt.Sprintf("update fileshare failed: %s", err.Error())
		f.ErrorHandle(model.ErrorInternalServer, errMsg)
		return
	}

	// Marshal the result.
	body, _ := json.Marshal(result)
	f.SuccessHandle(StatusOK, body)

	return
}

func (f *FileSharePortal) DeleteFileShareAcl() {
	if !policy.Authorize(f.Ctx, "fileshareacl:delete") {
		return
	}
	ctx := c.GetContext(f.Ctx)

	id := f.Ctx.Input.Param(":aclId")
	acl, err := db.C.GetFileShareAcl(ctx, id)
	if err != nil {
		errMsg := fmt.Sprintf("fileshare acl %s not found: %s", id, err.Error())
		f.ErrorHandle(model.ErrorNotFound, errMsg)
		return
	}
	fileshare, err := db.C.GetFileShare(ctx, acl.FileShareId)
	if err != nil {
		errMsg := fmt.Sprintf("fileshare for the acl %s not found: %s", id, err.Error())
		f.ErrorHandle(model.ErrorNotFound, errMsg)
		return
	}
	prf, err := db.C.GetProfile(ctx, fileshare.ProfileId)
	if err != nil {
		errMsg := fmt.Sprintf("get profile failed: %s", err.Error())
		f.ErrorHandle(model.ErrorBadRequest, errMsg)
		return
	}

	// NOTE: It will update the the status of the file share acl waiting for deletion
	// in the database to "deleting" and return the result immediately.
	if err = util.DeleteFileShareAclDBEntry(ctx, acl); err != nil {
		errMsg := fmt.Sprintf("delete fileshare acl failed: %v", err.Error())
		f.ErrorHandle(model.ErrorBadRequest, errMsg)
		return
	}

	f.SuccessHandle(StatusAccepted, nil)

	// NOTE: The real file share deletion process.
	// File Share deletion request is sent to the Dock. Dock will delete file share from driver
	// and database or update file share status to "errorDeleting" if deletion from driver failed.
	if err := f.CtrClient.Connect(CONF.OsdsLet.ApiEndpoint); err != nil {
		log.Error("when connecting controller client:", err)
		return
	}
	defer f.CtrClient.Close()

	opt := &pb.DeleteFileShareAclOpts{
		Id:               acl.Id,
		FileshareId:      acl.FileShareId,
		Description:      acl.Description,
		Type:             acl.Type,
		AccessCapability: acl.AccessCapability,
		AccessTo:         acl.AccessTo,
		Metadata:         utils.MergeStringMaps(fileshare.Metadata, acl.Metadata),
		Context:          ctx.ToJson(),
		Profile:          prf.ToJson(),
	}
	response, err := f.CtrClient.DeleteFileShareAcl(context.Background(), opt)
	if err != nil {
		log.Error("delete fileshare acl failed in controller service:", err)
		return
	}
	if errorMsg := response.GetError(); errorMsg != nil {
		log.Errorf("failed to delete fileshare acl in controller, code: %v, message: %v",
			errorMsg.GetCode(), errorMsg.GetDescription())
		return
	}

	return
}

func (f *FileSharePortal) DeleteFileShare() {
	if !policy.Authorize(f.Ctx, "fileshare:delete") {
		return
	}
	ctx := c.GetContext(f.Ctx)

	id := f.Ctx.Input.Param(":fileshareId")
	fileshare, err := db.C.GetFileShare(ctx, id)
	if err != nil {
		errMsg := fmt.Sprintf("fileshare %s not found: %s", id, err.Error())
		f.ErrorHandle(model.ErrorNotFound, errMsg)
		return
	}
	prf, err := db.C.GetProfile(ctx, fileshare.ProfileId)
	if err != nil {
		errMsg := fmt.Sprintf("delete file share failed: %v", err.Error())
		f.ErrorHandle(model.ErrorInternalServer, errMsg)
		return
	}

	// If profileId or poolId of the file share doesn't exist, it would mean that
	// the file share provisioning operation failed before the create method in
	// storage driver was called, therefore the file share entry should be deleted
	// from db directly.
	if fileshare.ProfileId == "" || fileshare.PoolId == "" {
		if err := db.C.DeleteFileShare(ctx, fileshare.Id); err != nil {
			errMsg := fmt.Sprintf("delete file share failed: %v", err.Error())
			f.ErrorHandle(model.ErrorInternalServer, errMsg)
			return
		}
		f.SuccessHandle(StatusAccepted, nil)
		return
	}

	// NOTE: It will update the the status of the file share waiting for deletion in
	// the database to "deleting" and return the result immediately.
	if err = util.DeleteFileShareDBEntry(ctx, fileshare); err != nil {
		errMsg := fmt.Sprintf("delete fileshare failed: %v", err.Error())
		f.ErrorHandle(model.ErrorBadRequest, errMsg)
		return
	}

	f.SuccessHandle(StatusAccepted, nil)

	// NOTE: The real file share deletion process.
	// File Share deletion request is sent to the Dock. Dock will delete file share from driver
	// and database or update file share status to "errorDeleting" if deletion from driver failed.
	if err := f.CtrClient.Connect(CONF.OsdsLet.ApiEndpoint); err != nil {
		log.Error("when connecting controller client:", err)
		return
	}
	defer f.CtrClient.Close()

	opt := &pb.DeleteFileShareOpts{
		Id:       fileshare.Id,
		PoolId:   fileshare.PoolId,
		Metadata: fileshare.Metadata,
		Context:  ctx.ToJson(),
		Profile:  prf.ToJson(),
	}
	response, err := f.CtrClient.DeleteFileShare(context.Background(), opt)
	if err != nil {
		log.Error("delete fileshare failed in controller service:", err)
		return
	}
	if errorMsg := response.GetError(); errorMsg != nil {
		log.Errorf("failed to delete fileshare in controller, code: %v, message: %v",
			errorMsg.GetCode(), errorMsg.GetDescription())
		return
	}

	return
}

func NewFileShareSnapshotPortal() *FileShareSnapshotPortal {
	return &FileShareSnapshotPortal{
		CtrClient: client.NewClient(),
	}
}

type FileShareSnapshotPortal struct {
	BasePortal

	CtrClient client.Client
}

func (f *FileShareSnapshotPortal) CreateFileShareSnapshot() {
	if !policy.Authorize(f.Ctx, "snapshot:create") {
		return
	}
	ctx := c.GetContext(f.Ctx)
	var snapshot = model.FileShareSnapshotSpec{
		BaseModel: &model.BaseModel{},
	}

	if err := json.NewDecoder(f.Ctx.Request.Body).Decode(&snapshot); err != nil {
		errMsg := fmt.Sprintf("parse fileshare snapshot request body failed: %s", err.Error())
		f.ErrorHandle(model.ErrorBadRequest, errMsg)
		return
	}

	fileshare, err := db.C.GetFileShare(ctx, snapshot.FileShareId)
	if err != nil {
		errMsg := fmt.Sprintf("fileshare %s not found: %s", snapshot.FileShareId, err.Error())
		f.ErrorHandle(model.ErrorNotFound, errMsg)
		return
	}
	snapshot.ShareSize = fileshare.Size
	// Usually snapshot.SnapshotSize and fileshare.Size are equal, even if they
	// are not equal, then snapshot.SnapshotSize will be updated to the correct value.
	snapshot.SnapshotSize = fileshare.Size

	if len(snapshot.ProfileId) == 0 {
		log.Warning("User doesn't specified profile id, using profile derived form fileshare")
		snapshot.ProfileId = fileshare.ProfileId
	}

	// Get profile
	prf, err := db.C.GetProfile(ctx, snapshot.ProfileId)
	if err != nil {
		errMsg := fmt.Sprintf("get profile failed: %s", err.Error())
		f.ErrorHandle(model.ErrorBadRequest, errMsg)
		return
	}

	// NOTE:It will create a fileshare snapshot entry into the database and initialize its status
	// as "creating". It will not wait for the real fileshare snapshot creation to complete
	// and will return result immediately.
	result, err := util.CreateFileShareSnapshotDBEntry(ctx, &snapshot)
	if err != nil {
		errMsg := fmt.Sprintf("create fileshare snapshot failed: %s", err.Error())
		f.ErrorHandle(model.ErrorBadRequest, errMsg)
		return
	}

	// Marshal the result.
	body, _ := json.Marshal(result)
	f.SuccessHandle(StatusAccepted, body)

	// NOTE:The real file share snapshot creation process.
	// FileShare snapshot creation request is sent to the Dock. Dock will update file share snapshot status to "available"
	// after file share snapshot creation complete.
	if err := f.CtrClient.Connect(CONF.OsdsLet.ApiEndpoint); err != nil {
		log.Error("when connecting controller client:", err)
		return
	}
	defer f.CtrClient.Close()

	opt := &pb.CreateFileShareSnapshotOpts{
		Id:          result.Id,
		Name:        result.Name,
		Description: result.Description,
		FileshareId: result.FileShareId,
		Size:        result.ShareSize,
		Context:     ctx.ToJson(),
		Metadata:    result.Metadata,
		Profile:     prf.ToJson(),
	}
	response, err := f.CtrClient.CreateFileShareSnapshot(context.Background(), opt)
	if err != nil {
		log.Error("create file share snapthot failed in controller service:", err)
		return
	}
	if errorMsg := response.GetError(); errorMsg != nil {
		log.Errorf("failed to create file share snapshot in controller, code: %v, message: %v",
			errorMsg.GetCode(), errorMsg.GetDescription())
		return
	}

	return
}

func (f *FileShareSnapshotPortal) ListFileShareSnapshots() {
	if !policy.Authorize(f.Ctx, "snapshot:list") {
		return
	}
	m, err := f.GetParameters()
	if err != nil {
		errMsg := fmt.Sprintf("list fileshare snapshots failed: %s", err.Error())
		f.ErrorHandle(model.ErrorBadRequest, errMsg)
		return
	}

	result, err := db.C.ListFileShareSnapshotsWithFilter(c.GetContext(f.Ctx), m)
	if err != nil {
		errMsg := fmt.Sprintf("list fileshare snapshots failed: %s", err.Error())
		f.ErrorHandle(model.ErrorInternalServer, errMsg)
		return
	}

	// Marshal the result.
	body, _ := json.Marshal(result)
	f.SuccessHandle(StatusOK, body)
	return
}

func (f *FileShareSnapshotPortal) GetFileShareSnapshot() {
	if !policy.Authorize(f.Ctx, "snapshot:get") {
		return
	}
	id := f.Ctx.Input.Param(":snapshotId")

	result, err := db.C.GetFileShareSnapshot(c.GetContext(f.Ctx), id)
	if err != nil {
		errMsg := fmt.Sprintf("fileshare snapshot %s not found: %s", id, err.Error())
		f.ErrorHandle(model.ErrorNotFound, errMsg)
		return
	}

	// Marshal the result.
	body, _ := json.Marshal(result)
	f.SuccessHandle(StatusOK, body)

	return
}

func (f *FileShareSnapshotPortal) UpdateFileShareSnapshot() {
	if !policy.Authorize(f.Ctx, "snapshot:update") {
		return
	}
	var snapshot = model.FileShareSnapshotSpec{
		BaseModel: &model.BaseModel{},
	}

	id := f.Ctx.Input.Param(":snapshotId")
	if err := json.NewDecoder(f.Ctx.Request.Body).Decode(&snapshot); err != nil {
		errMsg := fmt.Sprintf("parse fileshare snapshot request body failed: %s", err.Error())
		f.ErrorHandle(model.ErrorBadRequest, errMsg)
		return
	}
	snapshot.Id = id

	result, err := db.C.UpdateFileShareSnapshot(c.GetContext(f.Ctx), id, &snapshot)
	if err != nil {
		errMsg := fmt.Sprintf("update fileshare snapshot failed: %s", err.Error())
		f.ErrorHandle(model.ErrorInternalServer, errMsg)
		return
	}

	// Marshal the result.
	body, _ := json.Marshal(result)
	f.SuccessHandle(StatusOK, body)

	return
}

func (f *FileShareSnapshotPortal) DeleteFileShareSnapshot() {
	if !policy.Authorize(f.Ctx, "snapshot:delete") {
		return
	}
	ctx := c.GetContext(f.Ctx)
	id := f.Ctx.Input.Param(":snapshotId")

	snapshot, err := db.C.GetFileShareSnapshot(ctx, id)
	if err != nil {
		errMsg := fmt.Sprintf("fileshare snapshot %s not found: %s", id, err.Error())
		f.ErrorHandle(model.ErrorNotFound, errMsg)
		return
	}

	prf, err := db.C.GetProfile(ctx, snapshot.ProfileId)
	if err != nil {
		errMsg := fmt.Sprintf("profile (%s) not found: %v", snapshot.ProfileId, err.Error())
		f.ErrorHandle(model.ErrorBadRequest, errMsg)
		return
	}

	// NOTE: It will update the the status of the file share snapshot waiting for deletion in
	// the database to "deleting" and return the result immediately.
	err = util.DeleteFileShareSnapshotDBEntry(ctx, snapshot)
	if err != nil {
		errMsg := fmt.Sprintf("delete file share snapshot in db failed: %v", err.Error())
		f.ErrorHandle(model.ErrorBadRequest, errMsg)
		return
	}

	f.Ctx.Output.SetStatus(StatusAccepted)

	// NOTE:The real file share snapshot deletion process.
	// FileShare snapshot deletion request is sent to the Dock. Dock will delete file share snapshot from driver and
	// database or update its status to "errorDeleting" if file share snapshot deletion from driver failed.
	if err := f.CtrClient.Connect(CONF.OsdsLet.ApiEndpoint); err != nil {
		log.Error("when connecting controller client:", err)
		return
	}
	defer f.CtrClient.Close()

	opt := &pb.DeleteFileShareSnapshotOpts{
		Id:          snapshot.Id,
		FileshareId: snapshot.FileShareId,
		Context:     ctx.ToJson(),
		Profile:     prf.ToJson(),
		Metadata:    snapshot.Metadata,
	}
	response, err := f.CtrClient.DeleteFileShareSnapshot(context.Background(), opt)
	if err != nil {
		log.Error("delete file share snapshot failed in controller service:", err)
		return
	}
	if errorMsg := response.GetError(); errorMsg != nil {
		log.Errorf("failed to delete file share snapshot in controller, code: %v, message: %v",
			errorMsg.GetCode(), errorMsg.GetDescription())
		return
	}

	return
}
