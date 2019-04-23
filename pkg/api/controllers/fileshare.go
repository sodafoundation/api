package controllers

import (
	"encoding/json"
	"fmt"

	log "github.com/golang/glog"
	"github.com/opensds/opensds/pkg/api/util"
	c "github.com/opensds/opensds/pkg/context"
	"github.com/opensds/opensds/pkg/filesharecontroller/client"
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

func (f *FileSharePortal) CreateFileShare() {
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
	result, err := util.CreateFileShareDBEntry(c.GetContext(f.Ctx), &fileshare)
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
	if err = util.DeleteFileShareDBEntry(ctx, fshare); err != nil {
		errMsg := fmt.Sprintf("delete fileshare failed: %v", err.Error())
		f.ErrorHandle(model.ErrorBadRequest, errMsg)
		return
	}

	f.SuccessHandle(StatusAccepted, nil)

	return
}


