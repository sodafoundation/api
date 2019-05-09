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
	"encoding/json"
	"fmt"

	log "github.com/golang/glog"
	"github.com/opensds/opensds/pkg/api/util"
	c "github.com/opensds/opensds/pkg/context"
	"github.com/opensds/opensds/pkg/controller/client"
	"github.com/opensds/opensds/pkg/db"
	"github.com/opensds/opensds/pkg/model"
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
		reason := fmt.Sprintf("Parse fileshare access rules request body failed: %s", err.Error())
		f.ErrorHandle(model.ErrorBadRequest, reason)
		log.Error(reason)
		return
	}
	result, err := util.CreateFileShareAclDBEntry(c.GetContext(f.Ctx), &fileshareacl)
	if err != nil {
		reason := fmt.Sprintf("Create access rules for fileshare failed: %s", err.Error())
		f.ErrorHandle(model.ErrorBadRequest, reason)
		log.Error(reason)
		return
	}
	// Marshal the result.
	body, err := json.Marshal(result)
	if err != nil {
		reason := fmt.Sprintf("Marshal fileshare access rules created result failed: %s", err.Error())
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
	var fileshare = model.FileShareSpec{
		BaseModel: &model.BaseModel{},
	}

	// Unmarshal the request body
	if err := json.NewDecoder(f.Ctx.Request.Body).Decode(&fileshare); err != nil {
		reason := fmt.Sprintf("Parse fileshare request body failed: %s", err.Error())
		f.ErrorHandle(model.ErrorBadRequest, reason)
		log.Error(reason)
		return
	}
	result, err := util.CreateFileShareDBEntry(c.GetContext(f.Ctx), &fileshare)
	if err != nil {
		reason := fmt.Sprintf("Create fileshare failed: %s", err.Error())
		f.ErrorHandle(model.ErrorBadRequest, reason)
		log.Error(reason)
		return
	}
	// Marshal the result.
	body, err := json.Marshal(result)
	if err != nil {
		reason := fmt.Sprintf("Marshal fileshare created result failed: %s", err.Error())
		f.ErrorHandle(model.ErrorBadRequest, reason)
		log.Error(reason)
		return
	}

	f.SuccessHandle(StatusAccepted, body)
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
		reason := fmt.Sprintf("Marshal fileshare acl list result failed: %s", err.Error())
		f.ErrorHandle(model.ErrorBadRequest, reason)
		log.Error(reason)
		return
	}
	f.SuccessHandle(StatusOK, body)

	return
}

func (f *FileSharePortal) GetFileShare() {
	id := f.Ctx.Input.Param(":fileshareId")

	// Call db api module to handle get fileshare request.
	result, err := db.C.GetFileShare(c.GetContext(f.Ctx), id)
	if err != nil {
		errMsg := fmt.Sprintf("fileshare %s not found: %s", id, err.Error())
		f.ErrorHandle(model.ErrorNotFound, errMsg)
		return
	}

	// Marshal the result.
	body, _ := json.Marshal(result)
	if err != nil {
		reason := fmt.Sprintf("Marshal fileshare list result failed: %s", err.Error())
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
	fshare, err := db.C.GetFileShare(ctx, id)
	if err != nil {
		errMsg := fmt.Sprintf("fileshare %s not found: %s", id, err.Error())
		f.ErrorHandle(model.ErrorNotFound, errMsg)
		return
	}

	if err := db.C.DeleteFileShare(ctx, fshare.Id); err != nil {
		errMsg := fmt.Sprintf("delete fileshare failed: %v", err.Error())
		f.ErrorHandle(model.ErrorInternalServer, errMsg)
		return
	}
	f.SuccessHandle(StatusAccepted, nil)
	return

	// NOTE:It will update the the status of the fileshare waiting for deletion in
	// the database to "deleting" and return the result immediately.
	if err = util.DeleteFileShareDBEntry(ctx, fshare); err != nil {
		errMsg := fmt.Sprintf("delete fileshare failed: %v", err.Error())
		f.ErrorHandle(model.ErrorBadRequest, errMsg)
		return
	}

	f.SuccessHandle(StatusAccepted, nil)

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

	// NOTE:It will update the the status of the volume snapshot waiting for deletion in
	// the database to "deleting" and return the result immediately.
	err = util.DeleteFileShareSnapshotDBEntry(ctx, snapshot)
	if err != nil {
		errMsg := fmt.Sprintf("delete volume snapshot failed: %v", err.Error())
		f.ErrorHandle(model.ErrorBadRequest, errMsg)
		return
	}
	f.Ctx.Output.SetStatus(StatusAccepted)
	return
}
