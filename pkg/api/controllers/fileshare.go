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
	pb "github.com/opensds/opensds/pkg/model/fileshareproto"
	. "github.com/opensds/opensds/pkg/utils/config"
	"golang.org/x/net/context"
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
	ctx := c.GetContext(f.Ctx)
	var fileshare = model.FileShareSpec{
		BaseModel: &model.BaseModel{},
	}
	// Unmarshal the request body
	if err := json.NewDecoder(f.Ctx.Request.Body).Decode(&fileshare); err != nil {
		reason := fmt.Sprintf("Parse fileshare request body failed: %s", err.Error())
		//f.Ctx.Output.SetStatus(model.ErrorBadRequest)
		f.ErrorHandle(model.ErrorBadRequest, reason)
		//f.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		//log.Error(reason)
		return
	}

	// get profile
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

	// NOTE:It will create a file share entry into the database and initialize its status
	// as "creating". It will not wait for the real file share creation to complete
	// and will return result immediately.
	result, err := util.CreateFileShareDBEntry(c.GetContext(f.Ctx), &fileshare)
	if err != nil {
		reason := fmt.Sprintf("Create fileshare failed: %s", err.Error())
		//f.Ctx.Output.SetStatus(model.ErrorBadRequest)
		//f.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		f.ErrorHandle(model.ErrorBadRequest, reason)
		//log.Error(reason)
		return
	}
	// Marshal the result.
	body, _ := json.Marshal(result)
	f.SuccessHandle(StatusAccepted, body)

	// NOTE:The real volume creation process.
	// Volume creation request is sent to the Dock. Dock will update volume status to "available"
	// after volume creation is completed.
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
		// TODO: ProfileId will be removed later.
		ProfileId:         result.ProfileId,
		Profile:           prf.ToJson(),
		PoolId:            result.PoolId,
		//SnapshotId:        result.SnapshotId,
		Metadata:          result.Metadata,
		//SnapshotFromCloud: result.SnapshotFromCloud,
		Context:           ctx.ToJson(),
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
		if err := db.C.DeleteVolume(ctx, fileshare.Id); err != nil {
			errMsg := fmt.Sprintf("delete file share failed: %v", err.Error())
			f.ErrorHandle(model.ErrorInternalServer, errMsg)
			return
		}
		f.SuccessHandle(StatusAccepted, nil)
		return
	}

	// NOTE:It will update the the status of the fileshare waiting for deletion in
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

	// NOTE:The real fileshare deletion process.
	// File Share deletion request is sent to the Dock. Dock will delete file share from driver
	// and database or update file share status to "errorDeleting" if deletion from driver faild.
	if err := f.CtrClient.Connect(CONF.OsdsLet.ApiEndpoint); err != nil {
		log.Error("when connecting controller client:", err)
		return
	}
	defer f.CtrClient.Close()

	opt := &pb.DeleteFileShareOpts{
		Id:        fileshare.Id,
		ProfileId: fileshare.ProfileId,
		PoolId:    fileshare.PoolId,
		Metadata:  fileshare.Metadata,
		Context:   ctx.ToJson(),
		Profile:   prf.ToJson(),
	}
	if _, err = f.CtrClient.DeleteFileShare(context.Background(), opt); err != nil {
		log.Error("delete fileshare failed in controller service:", err)
		return
	}

	return
}


