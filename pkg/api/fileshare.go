package api

import (
	"encoding/json"
	"fmt"
	"github.com/opensds/opensds/pkg/db"
	"github.com/opensds/opensds/pkg/model"

	log "github.com/golang/glog"
	"github.com/opensds/opensds/pkg/api/policy"
	c "github.com/opensds/opensds/pkg/context"
	"github.com/opensds/opensds/pkg/filesharecontroller/client"

	pb "github.com/opensds/opensds/pkg/model/fileshareproto"
	. "github.com/opensds/opensds/pkg/utils/config"
	"golang.org/x/net/context"
)

func NewFileSharePortal() *FileSharePortal {
	fmt.Sprintf("entered to initiate file share client")
	return &FileSharePortal{
		CtrClient: client.NewClient(),
	}
}

type FileSharePortal struct {
	BasePortal

	CtrClient client.Client
}

func (v *FileSharePortal) CreateFileShare() {
	if !policy.Authorize(v.Ctx, "fileshare:create") {
		return
	}
	ctx := c.GetContext(v.Ctx)
	fmt.Sprintf("Getting fileshare specs which is in model")
	var fileshare = model.FileShareSpec{
		BaseModel: &model.BaseModel{},
	}

	// Unmarshal the request body
	if err := json.NewDecoder(v.Ctx.Request.Body).Decode(&fileshare); err != nil {
		errMsg := fmt.Sprintf("parse file share request body failed: %s", err.Error())
		v.ErrorHandle(model.ErrorBadRequest, errMsg)
		return
	}
	// NOTE:It will create a file share entry into the database and initialize its status
	// as "creating". It will not wait for the real file share creation to complete
	// and will return result immediately.
	result, err := CreateFileShareDBEntry(ctx, &fileshare)
	if err != nil {
		errMsg := fmt.Sprintf("create file share failed: %s", err.Error())
		v.ErrorHandle(model.ErrorBadRequest, errMsg)
		return
	}

	// Marshal the result.
	body, _ := json.Marshal(result)
	v.SuccessHandle(StatusAccepted, body)

	// NOTE:The real file share creation process.
	// File Share creation request is sent to the Dock. Dock will update file share status to "available"
	// after file share creation is completed.
	if err := v.CtrClient.Connect(CONF.OsdsLet.ApiEndpoint); err != nil {
		log.Error("when connecting controller client:", err)
		return
	}
	defer v.CtrClient.Close()

	opt := &pb.CreateFileShareOpts{
		Id:                result.Id,
		Name:              result.Name,
		Description:       result.Description,
		Size:              result.Size,
		AvailabilityZone:  result.AvailabilityZone,
		ProfileId:         result.ProfileId,
		PoolId:            result.PoolId,
		Metadata:          result.Metadata,
		Context:           ctx.ToJson(),
	}
	if _, err = v.CtrClient.CreateFileShare(context.Background(), opt); err != nil {
		log.Error("create file share failed in controller service:", err)
		return
	}

	return
}

func (v *FileSharePortal) ListFileShare() {
	//if !policy.Authorize(v.Ctx, "fileshare:list") {
	//        return
	//}
	m, err := v.GetParameters()
	if err != nil {
		errMsg := fmt.Sprintf("list fileshares failed: %s", err.Error())
		v.ErrorHandle(model.ErrorBadRequest, errMsg)
		return
	}
	result, err := db.C.ListFileSharesWithFilter(c.GetContext(v.Ctx), m)
	if err != nil {
		errMsg := fmt.Sprintf("list fileshares failed: %s", err.Error())
		v.ErrorHandle(model.ErrorInternalServer, errMsg)
		return
	}
	// Marshal the result.
	body, _ := json.Marshal(result)
	v.SuccessHandle(StatusOK, body)

	return
}

func (v *FileSharePortal) GetFileShare() {
	if !policy.Authorize(v.Ctx, "fileshare:get") {
		return
	}
	id := v.Ctx.Input.Param(":fileshareId")

	// Call db api module to handle get fileshare request.
	result, err := db.C.GetFileShare(c.GetContext(v.Ctx), id)
	if err != nil {
		errMsg := fmt.Sprintf("file share %s not found: %s", id, err.Error())
		v.ErrorHandle(model.ErrorNotFound, errMsg)
		return
	}

	// Marshal the result.
	body, _ := json.Marshal(result)
	v.SuccessHandle(StatusOK, body)

	return
}

func (v *FileSharePortal) UpdateFileShare() {
	if !policy.Authorize(v.Ctx, "fileshare:update") {
		return
	}
	var fileshare = model.FileShareSpec{
		BaseModel: &model.BaseModel{},
	}

	id := v.Ctx.Input.Param(":fileshareId")
	if err := json.NewDecoder(v.Ctx.Request.Body).Decode(&fileshare); err != nil {
		errMsg := fmt.Sprintf("parse file share request body failed: %s", err.Error())
		v.ErrorHandle(model.ErrorBadRequest, errMsg)
		return
	}

	fileshare.Id = id
	result, err := db.C.UpdateFileShare(c.GetContext(v.Ctx), &fileshare)
	if err != nil {
		errMsg := fmt.Sprintf("update file share failed: %s", err.Error())
		v.ErrorHandle(model.ErrorInternalServer, errMsg)
		return
	}

	// Marshal the result.
	body, _ := json.Marshal(result)
	v.SuccessHandle(StatusOK, body)

	return
}

func (v *FileSharePortal) DeleteFileShare() {
	if !policy.Authorize(v.Ctx, "fileshare:delete") {
		return
	}
	ctx := c.GetContext(v.Ctx)

	var err error
	id := v.Ctx.Input.Param(":fileshareId")
	fileshare, err := db.C.GetFileShare(ctx, id)
	if err != nil {
		errMsg := fmt.Sprintf("fileshare %s not found: %s", id, err.Error())
		v.ErrorHandle(model.ErrorNotFound, errMsg)
		return
	}

	// NOTE:It will update the the status of the fileshare waiting for deletion in
	// the database to "deleting" and return the result immediately.
	if err = DeleteFileShareDBEntry(ctx, fileshare); err != nil {
		errMsg := fmt.Sprintf("delete fileshare failed: %v", err.Error())
		v.ErrorHandle(model.ErrorBadRequest, errMsg)
		return
	}
	v.SuccessHandle(StatusAccepted, nil)

	// NOTE:The real fileshare deletion process.
	// FileShare deletion request is sent to the Dock. Dock will delete volume from driver
	// and database or update fileshare status to "errorDeleting" if deletion from driver faild.
	if err := v.CtrClient.Connect(CONF.OsdsLet.ApiEndpoint); err != nil {
		log.Error("when connecting controller client:", err)
		return
	}
	defer v.CtrClient.Close()

	opt := &pb.DeleteFileShareOpts{
		Id:        fileshare.Id,
		ProfileId: fileshare.ProfileId,
		PoolId:    fileshare.PoolId,
		Metadata:  fileshare.Metadata,
		Context:   ctx.ToJson(),
	}
	if _, err = v.CtrClient.DeleteFileShare(context.Background(), opt); err != nil {
		log.Error("delete fileshare failed in controller service:", err)
		return
	}

	return
}
