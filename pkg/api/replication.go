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

func (this *ReplicationPortal) CreateReplication() {
	if !policy.Authorize(this.Ctx, "replication:create") {
		return
	}

	var replication = &model.ReplicationSpec{
		BaseModel: &model.BaseModel{},
	}

	// Unmarshal the request body
	if err := json.NewDecoder(this.Ctx.Request.Body).Decode(replication); err != nil {
		model.HttpError(this.Ctx, http.StatusBadRequest,
			"parse replication request body failed: %s", err.Error())
		return
	}

	// Body check
	ctx := c.GetContext(this.Ctx)
	_, err := db.C.GetVolume(ctx, replication.PrimaryVolumeId)
	if err != nil {
		model.HttpError(this.Ctx, http.StatusBadRequest,
			"can't find the specified primary volume(%s)", replication.PrimaryVolumeId)
		return
	}
	_, err = db.C.GetVolume(ctx, replication.SecondaryVolumeId)
	if err != nil {
		model.HttpError(this.Ctx, http.StatusBadRequest,
			"can't find the specified secondary volume(%s)", replication.PrimaryVolumeId)
		return
	}

	// check if specified volume has already been used in other replication.
	v, err := db.C.GetReplicationByVolumeId(ctx, replication.PrimaryVolumeId)
	if err != nil {
		if _, ok := err.(*model.NotFoundError); !ok {
			model.HttpError(this.Ctx, http.StatusBadRequest,
				"get replication by volume id %s failed", replication.PrimaryVolumeId)
			return
		}
	}
	if v != nil {
		model.HttpError(this.Ctx, http.StatusBadRequest,
			"specified primary volume(%s) has already been used in replication(%s) ",
			replication.PrimaryVolumeId, v.Id)
		return
	}

	// check if specified volume has already been used in other replication.
	v, err = db.C.GetReplicationByVolumeId(ctx, replication.SecondaryVolumeId)
	if err != nil {
		if _, ok := err.(*model.NotFoundError); !ok {
			model.HttpError(this.Ctx, http.StatusBadRequest,
				"get replication by volume id %s failed", replication.SecondaryVolumeId)
			return
		}
	}
	if v != nil {
		model.HttpError(this.Ctx, http.StatusBadRequest,
			"specified secondary volume(%s) has already been used in replication(%s) ",
			replication.SecondaryVolumeId, v.Id)
		return
	}

	replication.ReplicationStatus = model.ReplicationCreating
	replication, err = db.C.CreateReplication(ctx, replication)
	if err != nil {
		model.HttpError(this.Ctx, http.StatusInternalServerError,
			"create replication in db failed")
		return
	}

	// Marshal the result.
	body, err := json.Marshal(replication)
	if err != nil {
		model.HttpError(this.Ctx, http.StatusBadRequest,
			"marshal replication created result failed: %s", err.Error())
		return
	}
	this.Ctx.Output.SetStatus(StatusAccepted)
	this.Ctx.Output.Body(body)

	// NOTE:The real volume replication creation process.
	// Volume replication creation request is sent to the Dock. Dock will update volume status to "available"
	// after volume replication creation is completed.
	if err = this.CtrClient.Connect(CONF.OsdsLet.ApiEndpoint); err != nil {
		log.Error("When connecting controller client:", err)
		return
	}
	defer this.CtrClient.Close()

	opt := &pb.CreateReplicationOpts{
		Message: string(body),
		Context: ctx.ToJson(),
	}
	if _, err = this.CtrClient.CreateReplication(context.Background(), opt); err != nil {
		log.Error("Create volume replication failed in controller service:", err)
		return
	}

	return
}

func (this *ReplicationPortal) ListReplications() {
	if !policy.Authorize(this.Ctx, "replication:list") {
		return
	}

	// Call db api module to handle list replications request.
	params, err := this.GetParameters()
	if err != nil {
		model.HttpError(this.Ctx, http.StatusBadRequest,
			"list replications failed: %s", err.Error())
		return
	}

	result, err := db.C.ListReplicationWithFilter(c.GetContext(this.Ctx), params)
	if err != nil {
		model.HttpError(this.Ctx, http.StatusBadRequest,
			"list replications failed: %s", err.Error())
		return
	}

	// Marshal the result.
	body, err := json.Marshal(this.outputFilter(result, whiteListSimple))
	if err != nil {
		model.HttpError(this.Ctx, http.StatusBadRequest,
			"marshal replications listed result failed: %s", err.Error())
		return
	}

	this.Ctx.Output.SetStatus(StatusOK)
	this.Ctx.Output.Body(body)
	return

}

func (this *ReplicationPortal) ListReplicationsDetail() {
	if !policy.Authorize(this.Ctx, "replication:list_detail") {
		return
	}

	// Call db api module to handle list replications request.
	params, err := this.GetParameters()
	if err != nil {
		model.HttpError(this.Ctx, http.StatusBadRequest,
			"list replications detail failed: %s", err.Error())
		return
	}

	result, err := db.C.ListReplicationWithFilter(c.GetContext(this.Ctx), params)
	if err != nil {
		model.HttpError(this.Ctx, http.StatusBadRequest,
			"list replications detail failed: %s", err.Error())
		return
	}

	// Marshal the result.
	body, err := json.Marshal(result)
	if err != nil {
		model.HttpError(this.Ctx, http.StatusInternalServerError,
			"marshal replications detail listed result failed: %s", err.Error())
		return
	}

	this.Ctx.Output.SetStatus(StatusOK)
	this.Ctx.Output.Body(body)
	return
}

func (this *ReplicationPortal) GetReplication() {
	if !policy.Authorize(this.Ctx, "replication:get") {
		return
	}

	id := this.Ctx.Input.Param(":replicationId")
	// Call db api module to handle get volume request.
	result, err := db.C.GetReplication(c.GetContext(this.Ctx), id)
	if err != nil {
		model.HttpError(this.Ctx, http.StatusBadRequest,
			"get replication failed: %s", err.Error())
		return
	}

	// Marshal the result.
	body, err := json.Marshal(this.outputFilter(result, whiteList))
	if err != nil {
		model.HttpError(this.Ctx, http.StatusInternalServerError,
			"marshal replication showed result failed: %s", err.Error())
		return
	}

	this.Ctx.Output.SetStatus(StatusOK)
	this.Ctx.Output.Body(body)
}

func (this *ReplicationPortal) UpdateReplication() {
	if !policy.Authorize(this.Ctx, "replication:update") {
		return
	}
	var r = model.ReplicationSpec{
		BaseModel: &model.BaseModel{},
	}

	id := this.Ctx.Input.Param(":replicationId")
	if err := json.NewDecoder(this.Ctx.Request.Body).Decode(&r); err != nil {
		model.HttpError(this.Ctx, http.StatusBadRequest,
			"parse replication request body failed: %s", err.Error())
		return
	}

	if r.ProfileId != "" {
		_, err := db.C.GetProfile(c.GetContext(this.Ctx), r.ProfileId)
		if err != nil {
			model.HttpError(this.Ctx, http.StatusBadRequest,
				"get profile failed: %s", err.Error())
			return
		}
		// TODO:compare with the original profile_id to get the differences
	}

	result, err := db.C.UpdateReplication(c.GetContext(this.Ctx), id, &r)
	if err != nil {
		model.HttpError(this.Ctx, http.StatusBadRequest,
			"update replication failed: %s", err.Error())
		return
	}

	// Marshal the result.
	body, err := json.Marshal(result)
	if err != nil {
		model.HttpError(this.Ctx, http.StatusInternalServerError,
			"marshal replication updated result failed: %s", err.Error())
		return
	}

	this.Ctx.Output.SetStatus(StatusOK)
	this.Ctx.Output.Body(body)
}

func (this *ReplicationPortal) DeleteReplication() {
	if !policy.Authorize(this.Ctx, "replication:delete") {
		return
	}
	ctx := c.GetContext(this.Ctx)

	id := this.Ctx.Input.Param(":replicationId")
	r, err := db.C.GetReplication(ctx, id)
	if err != nil {
		model.HttpError(this.Ctx, http.StatusBadRequest,
			"get replication failed: %s", err.Error())
		return
	}

	if err := DeleteReplicationDBEntry(ctx, r); err != nil {
		model.HttpError(this.Ctx, http.StatusBadRequest, err.Error())
		return
	}

	this.Ctx.Output.SetStatus(StatusAccepted)

	// NOTE:The real volume replication deletion process.
	// Volume replication deletion request is sent to the Dock. Dock will remove
	// replicaiton record after volume replication creation is completed.
	if err = this.CtrClient.Connect(CONF.OsdsLet.ApiEndpoint); err != nil {
		log.Error("When connecting controller client:", err)
		return
	}
	defer this.CtrClient.Close()

	body, _ := json.Marshal(r)
	opt := &pb.DeleteReplicationOpts{
		Message: string(body),
		Context: ctx.ToJson(),
	}
	if _, err = this.CtrClient.DeleteReplication(context.Background(), opt); err != nil {
		log.Error("Delete volume replication failed in controller service:", err)
		return
	}

	return
}

func (this *ReplicationPortal) EnableReplication() {
	if !policy.Authorize(this.Ctx, "replication:enable") {
		return
	}
	ctx := c.GetContext(this.Ctx)

	id := this.Ctx.Input.Param(":replicationId")
	r, err := db.C.GetReplication(ctx, id)
	if err != nil {
		model.HttpError(this.Ctx, http.StatusBadRequest,
			"get replication failed: %s", err.Error())
		return
	}

	if err := EnableReplicationDBEntry(ctx, r); err != nil {
		model.HttpError(this.Ctx, http.StatusBadRequest, err.Error())
		return
	}

	this.Ctx.Output.SetStatus(StatusAccepted)

	// NOTE:The real volume replication enable process.
	// Volume replication enable request is sent to the Dock. Dock will set
	// volume replication status to 'available' after volume replication enable
	// operation is completed.
	if err = this.CtrClient.Connect(CONF.OsdsLet.ApiEndpoint); err != nil {
		log.Error("When connecting controller client:", err)
		return
	}
	defer this.CtrClient.Close()

	body, _ := json.Marshal(r)
	opt := &pb.EnableReplicationOpts{
		Message: string(body),
		Context: ctx.ToJson(),
	}
	if _, err = this.CtrClient.EnableReplication(context.Background(), opt); err != nil {
		log.Error("Enable volume replication failed in controller service:", err)
		return
	}

	return
}

func (this *ReplicationPortal) DisableReplication() {
	if !policy.Authorize(this.Ctx, "replication:disable") {
		return
	}
	ctx := c.GetContext(this.Ctx)

	id := this.Ctx.Input.Param(":replicationId")
	r, err := db.C.GetReplication(ctx, id)
	if err != nil {
		model.HttpError(this.Ctx, http.StatusBadRequest,
			"get replication failed: %s", err.Error())
		return
	}

	if err := DisableReplicationDBEntry(ctx, r); err != nil {
		model.HttpError(this.Ctx, http.StatusBadRequest, err.Error())
		return
	}

	this.Ctx.Output.SetStatus(StatusAccepted)

	// NOTE:The real volume replication disable process.
	// Volume replication diable request is sent to the Dock. Dock will set
	// volume replication status to 'available' after volume replication disable
	// operation is completed.
	if err = this.CtrClient.Connect(CONF.OsdsLet.ApiEndpoint); err != nil {
		log.Error("When connecting controller client:", err)
		return
	}
	defer this.CtrClient.Close()

	body, _ := json.Marshal(r)
	opt := &pb.DisableReplicationOpts{
		Message: string(body),
		Context: ctx.ToJson(),
	}
	if _, err = this.CtrClient.DisableReplication(context.Background(), opt); err != nil {
		log.Error("Disable volume replication failed in controller service:", err)
		return
	}

	return
}

func (this *ReplicationPortal) FailoverReplication() {
	if !policy.Authorize(this.Ctx, "replication:failover") {
		return
	}
	ctx := c.GetContext(this.Ctx)

	var failover = model.FailoverReplicationSpec{}
	if err := json.NewDecoder(this.Ctx.Request.Body).Decode(&failover); err != nil {
		model.HttpError(this.Ctx, http.StatusBadRequest,
			"parse replication request body failed: %s", err.Error())
		return
	}

	id := this.Ctx.Input.Param(":replicationId")
	r, err := db.C.GetReplication(ctx, id)
	if err != nil {
		model.HttpError(this.Ctx, http.StatusBadRequest,
			"get replication failed: %s", err.Error())
		return
	}

	if err := FailoverReplicationDBEntry(ctx, r, failover.SecondaryBackendId); err != nil {
		model.HttpError(this.Ctx, http.StatusBadRequest, err.Error())
		return
	}

	this.Ctx.Output.SetStatus(StatusAccepted)

	// NOTE:The real volume replication failover process.
	// Volume replication failover request is sent to the Dock. Dock will set
	// volume replication status to 'available' after volume replication failover
	// operation is completed.
	if err = this.CtrClient.Connect(CONF.OsdsLet.ApiEndpoint); err != nil {
		log.Error("When connecting controller client:", err)
		return
	}
	defer this.CtrClient.Close()

	body, _ := json.Marshal(r)
	foBody, _ := json.Marshal(&failover)
	opt := &pb.FailoverReplicationOpts{
		Message:         string(body),
		FailoverMessage: string(foBody),
		Context: ctx.ToJson(),
	}
	if _, err = this.CtrClient.FailoverReplication(context.Background(), opt); err != nil {
		log.Error("Failover volume replication failed in controller service:", err)
		return
	}

	return
}
