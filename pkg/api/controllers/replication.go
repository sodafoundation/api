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
	. "github.com/opensds/opensds/pkg/utils/config"
)

func NewReplicationPortal() *ReplicationPortal {
	return &ReplicationPortal{
		CtrClient: client.NewClient(),
	}
}

type ReplicationPortal struct {
	BasePortal

	CtrClient client.Client
}

var whiteListSimple = []string{"Id", "Name", "ReplicationStatus"}
var whiteList = []string{"Id", "CreatedAt", "UpdatedAt", "Name", "Description", "AvailabilityZone", "ReplicationStatus",
	"PrimaryVolumeId", "SecondaryVolumeId", "PrimaryReplicationDriverData", "SecondaryReplicationDriverData",
	"ReplicationMode", "ReplicationPeriod", "ProfileId", "Metadata"}

func (r *ReplicationPortal) CreateReplication() {
	if !policy.Authorize(r.Ctx, "replication:create") {
		return
	}
	ctx := c.GetContext(r.Ctx)
	var replication = &model.ReplicationSpec{
		BaseModel: &model.BaseModel{},
	}

	// Unmarshal the request body
	if err := json.NewDecoder(r.Ctx.Request.Body).Decode(replication); err != nil {
		errMsg := fmt.Sprintf("parse replication request body failed: %s", err.Error())
		r.ErrorHandle(model.ErrorBadRequest, errMsg)
		return
	}

	result, err := util.CreateReplicationDBEntry(ctx, replication)
	if err != nil {
		errMsg := fmt.Sprintf("create volume replication failed: %s", err.Error())
		r.ErrorHandle(model.ErrorBadRequest, errMsg)
		return
	}

	// Marshal the result.
	body, err := json.Marshal(result)
	if err != nil {
		errMsg := fmt.Sprintf("marshal replication created result failed: %s", err.Error())
		r.ErrorHandle(model.ErrorInternalServer, errMsg)
		return
	}
	r.SuccessHandle(StatusAccepted, body)

	// NOTE:The real volume replication creation process.
	// Volume replication creation request is sent to the Dock. Dock will update volume status to "available"
	// after volume replication creation is completed.
	if err = r.CtrClient.Connect(CONF.OsdsLet.ApiEndpoint); err != nil {
		log.Error("when connecting controller client:", err)
		return
	}
	defer r.CtrClient.Close()

	opt := &pb.CreateReplicationOpts{
		Id:                result.Id,
		Name:              result.Name,
		Description:       result.Description,
		PrimaryVolumeId:   result.PrimaryVolumeId,
		SecondaryVolumeId: result.SecondaryVolumeId,
		AvailabilityZone:  result.AvailabilityZone,
		ProfileId:         result.ProfileId,
		Context:           ctx.ToJson(),
	}
	response, err := r.CtrClient.CreateReplication(context.Background(), opt)
	if err != nil {
		log.Error("create volume replication failed in controller service:", err)
		return
	}
	if errorMsg := response.GetError(); errorMsg != nil {
		log.Errorf("failed to create volume replication in controller, code: %v, message: %v",
			errorMsg.GetCode(), errorMsg.GetDescription())
		return
	}

	return
}

func (r *ReplicationPortal) ListReplications() {
	if !policy.Authorize(r.Ctx, "replication:list") {
		return
	}

	// Call db api module to handle list replications request.
	params, err := r.GetParameters()
	if err != nil {
		errMsg := fmt.Sprintf("list replications failed: %s", err.Error())
		r.ErrorHandle(model.ErrorBadRequest, errMsg)
		return
	}

	result, err := db.C.ListReplicationWithFilter(c.GetContext(r.Ctx), params)
	if err != nil {
		errMsg := fmt.Sprintf("list replications failed: %s", err.Error())
		r.ErrorHandle(model.ErrorInternalServer, errMsg)
		return
	}

	// Marshal the result.
	body, err := json.Marshal(r.outputFilter(result, whiteListSimple))
	if err != nil {
		errMsg := fmt.Sprintf("marshal replications listed result failed: %s", err.Error())
		r.ErrorHandle(model.ErrorInternalServer, errMsg)
		return
	}

	r.SuccessHandle(StatusOK, body)
	return

}

func (r *ReplicationPortal) ListReplicationsDetail() {
	if !policy.Authorize(r.Ctx, "replication:list_detail") {
		return
	}

	// Call db api module to handle list replications request.
	params, err := r.GetParameters()
	if err != nil {
		errMsg := fmt.Sprintf("list replications detail failed: %s", err.Error())
		r.ErrorHandle(model.ErrorBadRequest, errMsg)
		return
	}

	result, err := db.C.ListReplicationWithFilter(c.GetContext(r.Ctx), params)
	if err != nil {
		errMsg := fmt.Sprintf("list replications detail failed: %s", err.Error())
		r.ErrorHandle(model.ErrorInternalServer, errMsg)
		return
	}

	// Marshal the result.
	body, err := json.Marshal(result)
	if err != nil {
		errMsg := fmt.Sprintf("marshal replications detail listed result failed: %s", err.Error())
		r.ErrorHandle(model.ErrorInternalServer, errMsg)
		return
	}

	r.SuccessHandle(StatusOK, body)
	return
}

func (r *ReplicationPortal) GetReplication() {
	if !policy.Authorize(r.Ctx, "replication:get") {
		return
	}

	id := r.Ctx.Input.Param(":replicationId")
	// Call db api module to handle get volume request.
	result, err := db.C.GetReplication(c.GetContext(r.Ctx), id)
	if err != nil {
		errMsg := fmt.Sprintf("get replication failed: %s", err.Error())
		r.ErrorHandle(model.ErrorNotFound, errMsg)
		return
	}

	// Marshal the result.
	body, err := json.Marshal(r.outputFilter(result, whiteList))
	if err != nil {
		errMsg := fmt.Sprintf("marshal replication showed result failed: %s", err.Error())
		r.ErrorHandle(model.ErrorInternalServer, errMsg)
		return
	}

	r.SuccessHandle(StatusOK, body)
	return
}

func (r *ReplicationPortal) UpdateReplication() {
	if !policy.Authorize(r.Ctx, "replication:update") {
		return
	}
	var mr = model.ReplicationSpec{
		BaseModel: &model.BaseModel{},
	}

	id := r.Ctx.Input.Param(":replicationId")
	if err := json.NewDecoder(r.Ctx.Request.Body).Decode(&mr); err != nil {
		errMsg := fmt.Sprintf("parse replication request body failed: %s", err.Error())
		r.ErrorHandle(model.ErrorBadRequest, errMsg)
		return
	}

	if mr.ProfileId != "" {
		if _, err := db.C.GetProfile(c.GetContext(r.Ctx), mr.ProfileId); err != nil {
			errMsg := fmt.Sprintf("get profile failed: %s", err.Error())
			r.ErrorHandle(model.ErrorNotFound, errMsg)
			return
		}
		// TODO:compare with the original profile_id to get the differences
	}

	result, err := db.C.UpdateReplication(c.GetContext(r.Ctx), id, &mr)
	if err != nil {
		errMsg := fmt.Sprintf("update replication failed: %s", err.Error())
		r.ErrorHandle(model.ErrorInternalServer, errMsg)
		return
	}

	// Marshal the result.
	body, err := json.Marshal(result)
	if err != nil {
		errMsg := fmt.Sprintf("marshal replication updated result failed: %s", err.Error())
		r.ErrorHandle(model.ErrorInternalServer, errMsg)
		return
	}

	r.SuccessHandle(StatusOK, body)
	return
}

func (r *ReplicationPortal) DeleteReplication() {
	if !policy.Authorize(r.Ctx, "replication:delete") {
		return
	}
	ctx := c.GetContext(r.Ctx)

	id := r.Ctx.Input.Param(":replicationId")
	rep, err := db.C.GetReplication(ctx, id)
	if err != nil {
		errMsg := fmt.Sprintf("get replication failed: %s", err.Error())
		r.ErrorHandle(model.ErrorNotFound, errMsg)
		return
	}

	if err := util.DeleteReplicationDBEntry(ctx, rep); err != nil {
		r.ErrorHandle(model.ErrorBadRequest, err.Error())
		return
	}
	r.SuccessHandle(StatusAccepted, nil)

	// NOTE:The real volume replication deletion process.
	// Volume replication deletion request is sent to the Dock. Dock will remove
	// replicaiton record after volume replication creation is completed.
	if err = r.CtrClient.Connect(CONF.OsdsLet.ApiEndpoint); err != nil {
		log.Error("when connecting controller client:", err)
		return
	}
	defer r.CtrClient.Close()

	opt := &pb.DeleteReplicationOpts{
		Id:                rep.Id,
		PrimaryVolumeId:   rep.PrimaryVolumeId,
		SecondaryVolumeId: rep.SecondaryVolumeId,
		AvailabilityZone:  rep.AvailabilityZone,
		ProfileId:         rep.ProfileId,
		Metadata:          rep.Metadata,
		Context:           ctx.ToJson(),
	}
	response, err := r.CtrClient.DeleteReplication(context.Background(), opt)
	if err != nil {
		log.Error("delete volume replication failed in controller service:", err)
		return
	}
	if errorMsg := response.GetError(); errorMsg != nil {
		log.Errorf("failed to delete volume replication in controller, code: %v, message: %v",
			errorMsg.GetCode(), errorMsg.GetDescription())
		return
	}

	return
}

func (r *ReplicationPortal) EnableReplication() {
	if !policy.Authorize(r.Ctx, "replication:enable") {
		return
	}
	ctx := c.GetContext(r.Ctx)

	id := r.Ctx.Input.Param(":replicationId")
	rep, err := db.C.GetReplication(ctx, id)
	if err != nil {
		errMsg := fmt.Sprintf("get replication failed: %s", err.Error())
		r.ErrorHandle(model.ErrorNotFound, errMsg)
		return
	}

	if err := util.EnableReplicationDBEntry(ctx, rep); err != nil {
		r.ErrorHandle(model.ErrorBadRequest, err.Error())
		return
	}
	r.SuccessHandle(StatusAccepted, nil)

	// NOTE:The real volume replication enable process.
	// Volume replication enable request is sent to the Dock. Dock will set
	// volume replication status to 'available' after volume replication enable
	// operation is completed.
	if err = r.CtrClient.Connect(CONF.OsdsLet.ApiEndpoint); err != nil {
		log.Error("when connecting controller client:", err)
		return
	}
	defer r.CtrClient.Close()

	opt := &pb.EnableReplicationOpts{
		Id:                rep.Id,
		PrimaryVolumeId:   rep.PrimaryVolumeId,
		SecondaryVolumeId: rep.SecondaryVolumeId,
		AvailabilityZone:  rep.AvailabilityZone,
		ProfileId:         rep.ProfileId,
		Metadata:          rep.Metadata,
		Context:           ctx.ToJson(),
	}
	response, err := r.CtrClient.EnableReplication(context.Background(), opt)
	if err != nil {
		log.Error("enable volume replication failed in controller service:", err)
		return
	}
	if errorMsg := response.GetError(); errorMsg != nil {
		log.Errorf("failed to enable volume replication in controller, code: %v, message: %v",
			errorMsg.GetCode(), errorMsg.GetDescription())
		return
	}

	return
}

func (r *ReplicationPortal) DisableReplication() {
	if !policy.Authorize(r.Ctx, "replication:disable") {
		return
	}
	ctx := c.GetContext(r.Ctx)

	id := r.Ctx.Input.Param(":replicationId")
	rep, err := db.C.GetReplication(ctx, id)
	if err != nil {
		errMsg := fmt.Sprintf("get replication failed: %s", err.Error())
		r.ErrorHandle(model.ErrorNotFound, errMsg)
		return
	}

	if err := util.DisableReplicationDBEntry(ctx, rep); err != nil {
		r.ErrorHandle(model.ErrorBadRequest, err.Error())
		return
	}
	r.SuccessHandle(StatusAccepted, nil)

	// NOTE:The real volume replication disable process.
	// Volume replication diable request is sent to the Dock. Dock will set
	// volume replication status to 'available' after volume replication disable
	// operation is completed.
	if err = r.CtrClient.Connect(CONF.OsdsLet.ApiEndpoint); err != nil {
		log.Error("when connecting controller client:", err)
		return
	}
	defer r.CtrClient.Close()

	opt := &pb.DisableReplicationOpts{
		Id:                rep.Id,
		PrimaryVolumeId:   rep.PrimaryVolumeId,
		SecondaryVolumeId: rep.SecondaryVolumeId,
		AvailabilityZone:  rep.AvailabilityZone,
		ProfileId:         rep.ProfileId,
		Metadata:          rep.Metadata,
		Context:           ctx.ToJson(),
	}
	response, err := r.CtrClient.DisableReplication(context.Background(), opt)
	if err != nil {
		log.Error("disable volume replication failed in controller service:", err)
		return
	}
	if errorMsg := response.GetError(); errorMsg != nil {
		log.Errorf("failed to disable volume replication in controller, code: %v, message: %v",
			errorMsg.GetCode(), errorMsg.GetDescription())
		return
	}

	return
}

func (r *ReplicationPortal) FailoverReplication() {
	if !policy.Authorize(r.Ctx, "replication:failover") {
		return
	}
	ctx := c.GetContext(r.Ctx)

	var failover = model.FailoverReplicationSpec{}
	if err := json.NewDecoder(r.Ctx.Request.Body).Decode(&failover); err != nil {
		errMsg := fmt.Sprintf("parse replication request body failed: %s", err.Error())
		r.ErrorHandle(model.ErrorBadRequest, errMsg)
		return
	}

	id := r.Ctx.Input.Param(":replicationId")
	rep, err := db.C.GetReplication(ctx, id)
	if err != nil {
		errMsg := fmt.Sprintf("get replication failed: %s", err.Error())
		r.ErrorHandle(model.ErrorNotFound, errMsg)
		return
	}

	if err := util.FailoverReplicationDBEntry(ctx, rep, failover.SecondaryBackendId); err != nil {
		r.ErrorHandle(model.ErrorBadRequest, err.Error())
		return
	}
	r.SuccessHandle(StatusAccepted, nil)

	// NOTE:The real volume replication failover process.
	// Volume replication failover request is sent to the Dock. Dock will set
	// volume replication status to 'available' after volume replication failover
	// operation is completed.
	if err = r.CtrClient.Connect(CONF.OsdsLet.ApiEndpoint); err != nil {
		log.Error("when connecting controller client:", err)
		return
	}
	defer r.CtrClient.Close()

	opt := &pb.FailoverReplicationOpts{
		Id:                  rep.Id,
		PrimaryVolumeId:     rep.PrimaryVolumeId,
		SecondaryVolumeId:   rep.SecondaryVolumeId,
		AvailabilityZone:    rep.AvailabilityZone,
		ProfileId:           rep.ProfileId,
		Metadata:            rep.Metadata,
		AllowAttachedVolume: failover.AllowAttachedVolume,
		SecondaryBackendId:  failover.SecondaryBackendId,
		Context:             ctx.ToJson(),
	}
	response, err := r.CtrClient.FailoverReplication(context.Background(), opt)
	if err != nil {
		log.Error("failover volume replication failed in controller service:", err)
		return
	}
	if errorMsg := response.GetError(); errorMsg != nil {
		log.Errorf("failed to failover volume replication in controller, code: %v, message: %v",
			errorMsg.GetCode(), errorMsg.GetDescription())
		return
	}

	return
}
