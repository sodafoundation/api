package api

import (
	"encoding/json"
	"fmt"
	"github.com/opensds/opensds/pkg/db"
	"github.com/opensds/opensds/pkg/model"
	log "github.com/golang/glog"
	c "github.com/opensds/opensds/pkg/context"
	"github.com/opensds/opensds/pkg/filesharecontroller/client"
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

func (f *FileSharePortal) CreateFileShare() {
	//if !policy.Authorize(f.Ctx, "fileshare:create") {
	//      return
	//}
	//ctx := c.GetContext(f.Ctx)
	var fileshare = model.FileShareSpec{
		BaseModel: &model.BaseModel{},
	}

	// Unmarshal the request body
	if err := json.NewDecoder(f.Ctx.Request.Body).Decode(&fileshare); err != nil {
		reason := fmt.Sprintf("Parse fileshare request body failed: %s", err.Error())
		f.Ctx.Output.SetStatus(model.ErrorBadRequest)
		f.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}
	result, err := CreateFileShareDBEntry(c.GetContext(f.Ctx), &fileshare)
	if err != nil {
		reason := fmt.Sprintf("Create fileshare failed: %s", err.Error())
		f.Ctx.Output.SetStatus(model.ErrorBadRequest)
		f.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}
	// Marshal the result.
	body, err := json.Marshal(result)
	if err != nil {
		reason := fmt.Sprintf("Marshal fileshare created result failed: %s", err.Error())
		f.Ctx.Output.SetStatus(model.ErrorBadRequest)
		f.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	f.Ctx.Output.SetStatus(StatusAccepted)
	f.Ctx.Output.Body(body)
	return
}

func (f *FileSharePortal) ListFileShares() {
	//if !policy.Authorize(f.Ctx, "fileshare:list") {
	//        return
	//}
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

func (f *FileSharePortal) GetFileShare() {
	//if !policy.Authorize(f.Ctx, "fileshare:get") {
	//	return
	//}
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
	f.SuccessHandle(StatusOK, body)

	return
}

func (f *FileSharePortal) UpdateFileShare() {
	//if !policy.Authorize(f.Ctx, "fileshare:update") {
	//	return
	//}
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

func (f *FileSharePortal) DeleteFileShare() {
	//if !policy.Authorize(f.Ctx, "fileshare:delete") {
	//	return
	//}
	ctx := c.GetContext(f.Ctx)

	var err error
	id := f.Ctx.Input.Param(":fileshareId")
	fshare, err := db.C.GetFileShare(ctx, id)
	if err != nil {
		errMsg := fmt.Sprintf("fileshare %s not found: %s", id, err.Error())
		f.ErrorHandle(model.ErrorNotFound, errMsg)
		return
	}

	// NOTE:It will update the the status of the fileshare waiting for deletion in
	// the database to "deleting" and return the result immediately.
	if err = DeleteFileShareDBEntry(ctx, fshare); err != nil {
		errMsg := fmt.Sprintf("delete fileshare failed: %v", err.Error())
		f.ErrorHandle(model.ErrorBadRequest, errMsg)
		return
	}

	f.SuccessHandle(StatusAccepted, nil)

	return
}

/*
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
	//if !policy.Authorize(f.Ctx, "snapshot:create") {
	//	return
	//}
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
	result, err := CreateFileShareSnapshotDBEntry(ctx, &snapshot)
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
	//if !policy.Authorize(f.Ctx, "snapshot:list") {
	//	return
	//}
	fmt.Println("I m here")
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
	//if !policy.Authorize(f.Ctx, "snapshot:get") {
	//	return
	//}
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
	//if !policy.Authorize(f.Ctx, "snapshot:update") {
	//	return
	//}
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
	//if !policy.Authorize(f.Ctx, "snapshot:delete") {
	//	return
	//}
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
	err = DeleteFileShareSnapshotDBEntry(ctx, snapshot)
	if err != nil {
		errMsg := fmt.Sprintf("delete volume snapshot failed: %v", err.Error())
		f.ErrorHandle(model.ErrorBadRequest, errMsg)
		return
	}

	f.Ctx.Output.SetStatus(StatusAccepted)
	return
}
*/
