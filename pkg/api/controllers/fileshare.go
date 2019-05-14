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
	"github.com/opensds/opensds/pkg/api/util"
	c "github.com/opensds/opensds/pkg/context"
	"github.com/opensds/opensds/pkg/controller/client"
	"github.com/opensds/opensds/pkg/db"
	"github.com/opensds/opensds/pkg/model"
	pb "github.com/opensds/opensds/pkg/model/proto"
	. "github.com/opensds/opensds/pkg/utils/config"
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
		reason := fmt.Sprintf("create access rules for fileshare failed: %s", err.Error())
		f.ErrorHandle(model.ErrorBadRequest, reason)
		log.Error(reason)
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
	return
}

func (f *FileSharePortal) ListFileSharesAcl() {
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

// Function to store filesahre related entry into databse
func (f *FileSharePortal) CreateFileShare() {
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
		prf, err = db.C.GetDefaultProfile(ctx)
		fileshare.ProfileId = prf.Id
	} else {
		prf, err = db.C.GetProfile(ctx, fileshare.ProfileId)
	}
	if err != nil {
		errMsg := fmt.Sprintf("get profile failed: %s", err.Error())
		f.ErrorHandle(model.ErrorBadRequest, errMsg)
		return
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
		Metadata:         result.Metadata,
		Context:          ctx.ToJson(),
	}
	if _, err = f.CtrClient.CreateFileShare(context.Background(), opt); err != nil {
		log.Error("create file share failed in controller service:", err)
		return
	}

	return
}

func (f *FileSharePortal) ListFileShares() {
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
	ctx := c.GetContext(f.Ctx)

	var err error
	id := f.Ctx.Input.Param(":aclId")
	acl, err := db.C.GetFileShareAcl(ctx, id)
	if err != nil {
		errMsg := fmt.Sprintf("fileshare acl %s not found: %s", id, err.Error())
		f.ErrorHandle(model.ErrorNotFound, errMsg)
		return
	}

	if err := db.C.DeleteFileShareAcl(ctx, acl.Id); err != nil {
		errMsg := fmt.Sprintf("delete fileshare acl failed: %v", err.Error())
		f.ErrorHandle(model.ErrorInternalServer, errMsg)
		return
	}
	f.SuccessHandle(StatusAccepted, nil)
	return
}

func (f *FileSharePortal) DeleteFileShare() {
	ctx := c.GetContext(f.Ctx)

	var err error
	id := f.Ctx.Input.Param(":fileshareId")
	fileshare, err := db.C.GetFileShare(ctx, id)
	if err != nil {
		errMsg := fmt.Sprintf("fileshare %s not found: %s", id, err.Error())
		f.ErrorHandle(model.ErrorNotFound, errMsg)
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
	prf, err := db.C.GetProfile(ctx, fileshare.ProfileId)
	if err != nil {
		errMsg := fmt.Sprintf("delete file share failed: %v", err.Error())
		f.ErrorHandle(model.ErrorInternalServer, errMsg)
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
	if _, err = f.CtrClient.DeleteFileShare(context.Background(), opt); err != nil {
		log.Error("delete fileshare failed in controller service:", err)
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
	ctx := c.GetContext(f.Ctx)
	var snapshot = model.FileShareSnapshotSpec{
		BaseModel: &model.BaseModel{},
	}

	if err := json.NewDecoder(f.Ctx.Request.Body).Decode(&snapshot); err != nil {
		errMsg := fmt.Sprintf("parse fileshare snapshot request body failed: %s", err.Error())
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

	return
}

func (f *FileShareSnapshotPortal) ListFileShareSnapshots() {
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
	ctx := c.GetContext(f.Ctx)
	id := f.Ctx.Input.Param(":snapshotId")

	snapshot, err := db.C.GetFileShareSnapshot(ctx, id)
	if err != nil {
		errMsg := fmt.Sprintf("fileshare snapshot %s not found: %s", id, err.Error())
		f.ErrorHandle(model.ErrorNotFound, errMsg)
		return
	}

	// NOTE: It will update the the status of the file share snapshot waiting for deletion in
	// the database to "deleting" and return the result immediately.
	err = util.DeleteFileShareSnapshotDBEntry(ctx, snapshot)
	if err != nil {
		errMsg := fmt.Sprintf("delete file share snapshot failed: %v", err.Error())
		f.ErrorHandle(model.ErrorBadRequest, errMsg)
		return
	}
	f.Ctx.Output.SetStatus(StatusAccepted)
	return

}
