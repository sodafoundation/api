// Copyright (c) 2018 Huawei Technologies Co., Ltd. All Rights Reserved.
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

package api

import (
	"encoding/json"
	"net/http"

	log "github.com/golang/glog"
	"github.com/opensds/opensds/pkg/api/policy"
	c "github.com/opensds/opensds/pkg/context"
	"github.com/opensds/opensds/pkg/controller/client"
	pb "github.com/opensds/opensds/pkg/controller/proto"
	"github.com/opensds/opensds/pkg/db"
	"github.com/opensds/opensds/pkg/model"
	. "github.com/opensds/opensds/pkg/utils/config"
	"golang.org/x/net/context"
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

	var replication = &model.ReplicationSpec{
		BaseModel: &model.BaseModel{},
	}

	// Unmarshal the request body
	if err := json.NewDecoder(r.Ctx.Request.Body).Decode(replication); err != nil {
		model.HttpError(r.Ctx, http.StatusBadRequest,
			"parse replication request body failed: %s", err.Error())
		return
	}

	// Body check
	ctx := c.GetContext(r.Ctx)
	_, err := db.C.GetVolume(ctx, replication.PrimaryVolumeId)
	if err != nil {
		model.HttpError(r.Ctx, http.StatusBadRequest,
			"can't find the specified primary volume(%s)", replication.PrimaryVolumeId)
		return
	}
	_, err = db.C.GetVolume(ctx, replication.SecondaryVolumeId)
	if err != nil {
		model.HttpError(r.Ctx, http.StatusBadRequest,
			"can't find the specified secondary volume(%s)", replication.PrimaryVolumeId)
		return
	}

	// check if specified volume has already been used in other replication.
	v, err := db.C.GetReplicationByVolumeId(ctx, replication.PrimaryVolumeId)
	if err != nil {
		if _, ok := err.(*model.NotFoundError); !ok {
			model.HttpError(r.Ctx, http.StatusBadRequest,
				"get replication by volume id %s failed", replication.PrimaryVolumeId)
			return
		}
	}
	if v != nil {
		model.HttpError(r.Ctx, http.StatusBadRequest,
			"specified primary volume(%s) has already been used in replication(%s) ",
			replication.PrimaryVolumeId, v.Id)
		return
	}

	// check if specified volume has already been used in other replication.
	v, err = db.C.GetReplicationByVolumeId(ctx, replication.SecondaryVolumeId)
	if err != nil {
		if _, ok := err.(*model.NotFoundError); !ok {
			model.HttpError(r.Ctx, http.StatusBadRequest,
				"get replication by volume id %s failed", replication.SecondaryVolumeId)
			return
		}
	}
	if v != nil {
		model.HttpError(r.Ctx, http.StatusBadRequest,
			"specified secondary volume(%s) has already been used in replication(%s) ",
			replication.SecondaryVolumeId, v.Id)
		return
	}

	replication.ReplicationStatus = model.ReplicationCreating
	replication, err = db.C.CreateReplication(ctx, replication)
	if err != nil {
		model.HttpError(r.Ctx, http.StatusInternalServerError,
			"create replication in db failed")
		return
	}

	// Marshal the result.
	body, err := json.Marshal(replication)
	if err != nil {
		model.HttpError(r.Ctx, http.StatusBadRequest,
			"marshal replication created result failed: %s", err.Error())
		return
	}
	r.Ctx.Output.SetStatus(StatusAccepted)
	r.Ctx.Output.Body(body)

	// NOTE:The real volume replication creation process.
	// Volume replication creation request is sent to the Dock. Dock will update volume status to "available"
	// after volume replication creation is completed.
	if err = r.CtrClient.Connect(CONF.OsdsLet.ApiEndpoint); err != nil {
		log.Error("When connecting controller client:", err)
		return
	}
	defer r.CtrClient.Close()

	opt := &pb.CreateReplicationOpts{
		Message: string(body),
		Context: ctx.ToJson(),
	}
	if _, err = r.CtrClient.CreateReplication(context.Background(), opt); err != nil {
		log.Error("Create volume replication failed in controller service:", err)
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
		model.HttpError(r.Ctx, http.StatusBadRequest,
			"list replications failed: %s", err.Error())
		return
	}

	result, err := db.C.ListReplicationWithFilter(c.GetContext(r.Ctx), params)
	if err != nil {
		model.HttpError(r.Ctx, http.StatusBadRequest,
			"list replications failed: %s", err.Error())
		return
	}

	// Marshal the result.
	body, err := json.Marshal(r.outputFilter(result, whiteListSimple))
	if err != nil {
		model.HttpError(r.Ctx, http.StatusBadRequest,
			"marshal replications listed result failed: %s", err.Error())
		return
	}

	r.Ctx.Output.SetStatus(StatusOK)
	r.Ctx.Output.Body(body)
	return

}

func (r *ReplicationPortal) ListReplicationsDetail() {
	if !policy.Authorize(r.Ctx, "replication:list_detail") {
		return
	}

	// Call db api module to handle list replications request.
	params, err := r.GetParameters()
	if err != nil {
		model.HttpError(r.Ctx, http.StatusBadRequest,
			"list replications detail failed: %s", err.Error())
		return
	}

	result, err := db.C.ListReplicationWithFilter(c.GetContext(r.Ctx), params)
	if err != nil {
		model.HttpError(r.Ctx, http.StatusBadRequest,
			"list replications detail failed: %s", err.Error())
		return
	}

	// Marshal the result.
	body, err := json.Marshal(result)
	if err != nil {
		model.HttpError(r.Ctx, http.StatusInternalServerError,
			"marshal replications detail listed result failed: %s", err.Error())
		return
	}

	r.Ctx.Output.SetStatus(StatusOK)
	r.Ctx.Output.Body(body)
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
		model.HttpError(r.Ctx, http.StatusBadRequest,
			"get replication failed: %s", err.Error())
		return
	}

	// Marshal the result.
	body, err := json.Marshal(r.outputFilter(result, whiteList))
	if err != nil {
		model.HttpError(r.Ctx, http.StatusInternalServerError,
			"marshal replication showed result failed: %s", err.Error())
		return
	}

	r.Ctx.Output.SetStatus(StatusOK)
	r.Ctx.Output.Body(body)
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
		model.HttpError(r.Ctx, http.StatusBadRequest,
			"parse replication request body failed: %s", err.Error())
		return
	}

	if mr.ProfileId != "" {
		_, err := db.C.GetProfile(c.GetContext(r.Ctx), mr.ProfileId)
		if err != nil {
			model.HttpError(r.Ctx, http.StatusBadRequest,
				"get profile failed: %s", err.Error())
			return
		}
		// TODO:compare with the original profile_id to get the differences
	}

	result, err := db.C.UpdateReplication(c.GetContext(r.Ctx), id, &mr)
	if err != nil {
		model.HttpError(r.Ctx, http.StatusBadRequest,
			"update replication failed: %s", err.Error())
		return
	}

	// Marshal the result.
	body, err := json.Marshal(result)
	if err != nil {
		model.HttpError(r.Ctx, http.StatusInternalServerError,
			"marshal replication updated result failed: %s", err.Error())
		return
	}

	r.Ctx.Output.SetStatus(StatusOK)
	r.Ctx.Output.Body(body)
}

func (r *ReplicationPortal) DeleteReplication() {
	if !policy.Authorize(r.Ctx, "replication:delete") {
		return
	}
	ctx := c.GetContext(r.Ctx)

	id := r.Ctx.Input.Param(":replicationId")
	rep, err := db.C.GetReplication(ctx, id)
	if err != nil {
		model.HttpError(r.Ctx, http.StatusBadRequest,
			"get replication failed: %s", err.Error())
		return
	}

	if err := DeleteReplicationDBEntry(ctx, rep); err != nil {
		model.HttpError(r.Ctx, http.StatusBadRequest, err.Error())
		return
	}
	r.Ctx.Output.SetStatus(StatusAccepted)

	// NOTE:The real volume replication deletion process.
	// Volume replication deletion request is sent to the Dock. Dock will remove
	// replicaiton record after volume replication creation is completed.
	if err = r.CtrClient.Connect(CONF.OsdsLet.ApiEndpoint); err != nil {
		log.Error("When connecting controller client:", err)
		return
	}
	defer r.CtrClient.Close()

	body, _ := json.Marshal(rep)
	opt := &pb.DeleteReplicationOpts{
		Message: string(body),
		Context: ctx.ToJson(),
	}
	if _, err = r.CtrClient.DeleteReplication(context.Background(), opt); err != nil {
		log.Error("Delete volume replication failed in controller service:", err)
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
		model.HttpError(r.Ctx, http.StatusBadRequest,
			"get replication failed: %s", err.Error())
		return
	}

	if err := EnableReplicationDBEntry(ctx, rep); err != nil {
		model.HttpError(r.Ctx, http.StatusBadRequest, err.Error())
		return
	}

			r.Ctx.Output.SetStatus(StatusAccepted)

	// NOTE:The real volume replication enable process.
	// Volume replication enable request is sent to the Dock. Dock will set
	// volume replication status to 'available' after volume replication enable
	// operation is completed.
	if err = r.CtrClient.Connect(CONF.OsdsLet.ApiEndpoint); err != nil {
		log.Error("When connecting controller client:", err)
		return
	}
	defer r.CtrClient.Close()

	body, _ := json.Marshal(rep)
	opt := &pb.EnableReplicationOpts{
		Message: string(body),
		Context: ctx.ToJson(),
	}
	if _, err = r.CtrClient.EnableReplication(context.Background(), opt); err != nil {
		log.Error("Enable volume replication failed in controller service:", err)
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
		model.HttpError(r.Ctx, http.StatusBadRequest,
			"get replication failed: %s", err.Error())
		return
	}

	if err := DisableReplicationDBEntry(ctx, rep); err != nil {
		model.HttpError(r.Ctx, http.StatusBadRequest, err.Error())
		return
	}
	r.Ctx.Output.SetStatus(StatusAccepted)

	// NOTE:The real volume replication disable process.
	// Volume replication diable request is sent to the Dock. Dock will set
	// volume replication status to 'available' after volume replication disable
	// operation is completed.
	if err = r.CtrClient.Connect(CONF.OsdsLet.ApiEndpoint); err != nil {
		log.Error("When connecting controller client:", err)
		return
	}
	defer r.CtrClient.Close()

	body, _ := json.Marshal(rep)
	opt := &pb.DisableReplicationOpts{
		Message: string(body),
		Context: ctx.ToJson(),
	}
	if _, err = r.CtrClient.DisableReplication(context.Background(), opt); err != nil {
		log.Error("Disable volume replication failed in controller service:", err)
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
		model.HttpError(r.Ctx, http.StatusBadRequest,
			"parse replication request body failed: %s", err.Error())
		return
	}

	id := r.Ctx.Input.Param(":replicationId")
	rep, err := db.C.GetReplication(ctx, id)
	if err != nil {
		model.HttpError(r.Ctx, http.StatusBadRequest,
			"get replication failed: %s", err.Error())
		return
	}

	if err := FailoverReplicationDBEntry(ctx, rep, failover.SecondaryBackendId); err != nil {
		model.HttpError(r.Ctx, http.StatusBadRequest, err.Error())
		return
	}
	r.Ctx.Output.SetStatus(StatusAccepted)

	// NOTE:The real volume replication failover process.
	// Volume replication failover request is sent to the Dock. Dock will set
	// volume replication status to 'available' after volume replication failover
	// operation is completed.
	if err = r.CtrClient.Connect(CONF.OsdsLet.ApiEndpoint); err != nil {
		log.Error("When connecting controller client:", err)
		return
	}
	defer r.CtrClient.Close()

	body, _ := json.Marshal(rep)
	foBody, _ := json.Marshal(&failover)
	opt := &pb.FailoverReplicationOpts{
		Message:         string(body),
		FailoverMessage: string(foBody),
		Context: ctx.ToJson(),
	}
	if _, err = r.CtrClient.FailoverReplication(context.Background(), opt); err != nil {
		log.Error("Failover volume replication failed in controller service:", err)
		return
	}

	return
}
