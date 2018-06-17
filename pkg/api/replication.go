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

	"github.com/opensds/opensds/pkg/api/policy"
	c "github.com/opensds/opensds/pkg/context"
	"github.com/opensds/opensds/pkg/controller"
	"github.com/opensds/opensds/pkg/db"
	"github.com/opensds/opensds/pkg/model"
)

func NewReplicationPortal() *ReplicationPortal {
	return &ReplicationPortal{}
}

type ReplicationPortal struct {
	BasePortal
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

	replication.ReplicationStatus = model.ReplicationCreating
	replication, err = db.C.CreateReplication(ctx, replication)
	if err != nil {
		model.HttpError(this.Ctx, http.StatusInternalServerError,
			"create replication in db failed")
		return
	}
	// Call global controller variable to handle create replication request.
	result, err := controller.Brain.CreateReplication(ctx, replication)
	if err != nil {
		model.HttpError(this.Ctx, http.StatusBadRequest,
			"create replication failed: %s", err.Error())
		return
	}

	// Marshal the result.
	body, err := json.Marshal(result)
	if err != nil {
		model.HttpError(this.Ctx, http.StatusBadRequest,
			"marshal replication created result failed: %s", err.Error())
		return
	}

	this.Ctx.Output.SetStatus(StatusAccepted)
	this.Ctx.Output.Body(body)
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

	id := this.Ctx.Input.Param(":replicationId")
	r, err := db.C.GetReplication(c.GetContext(this.Ctx), id)
	if err != nil {
		model.HttpError(this.Ctx, http.StatusBadRequest,
			"get replication failed: %s", err.Error())
		return
	}

	if err := DeleteReplicationDBEntry(c.GetContext(this.Ctx), r); err != nil {
		model.HttpError(this.Ctx, http.StatusBadRequest, err.Error())
		return
	}
	// Call global controller variable to handle delete replication request.
	err = controller.Brain.DeleteReplication(c.GetContext(this.Ctx), r)
	if err != nil {
		model.HttpError(this.Ctx, http.StatusBadRequest,
			"delete replication failed: %v", err.Error())
		return
	}

	this.Ctx.Output.SetStatus(StatusAccepted)
}

func (this *ReplicationPortal) EnableReplication() {
	ctx := this.Ctx
	if !policy.Authorize(ctx, "replication:enable") {
		return
	}

	id := ctx.Input.Param(":replicationId")
	r, err := db.C.GetReplication(c.GetContext(ctx), id)
	if err != nil {
		model.HttpError(ctx, http.StatusBadRequest,
			"get replication failed: %s", err.Error())
		return
	}

	if err := EnableReplicationDBEntry(c.GetContext(ctx), r); err != nil {
		model.HttpError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	// Call global controller variable to handle delete replication request.
	err = controller.Brain.EnableReplication(c.GetContext(ctx), r)
	if err != nil {
		model.HttpError(ctx, http.StatusBadRequest,
			"enable replication failed: %v", err.Error())
		return
	}

	ctx.Output.SetStatus(StatusAccepted)
}

func (this *ReplicationPortal) DisableReplication() {
	ctx := this.Ctx
	if !policy.Authorize(ctx, "replication:disable") {
		return
	}

	id := ctx.Input.Param(":replicationId")
	r, err := db.C.GetReplication(c.GetContext(ctx), id)
	if err != nil {
		model.HttpError(ctx, http.StatusBadRequest,
			"get replication failed: %s", err.Error())
		return
	}

	if err := DisableReplicationDBEntry(c.GetContext(ctx), r); err != nil {
		model.HttpError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	// Call global controller variable to handle delete r request.
	err = controller.Brain.DisableReplication(c.GetContext(ctx), r)
	if err != nil {
		model.HttpError(ctx, http.StatusBadRequest,
			"enable replication failed: %v", err.Error())
		return
	}

	ctx.Output.SetStatus(StatusAccepted)
}

func (this *ReplicationPortal) FailoverReplication() {
	ctx := this.Ctx
	if !policy.Authorize(ctx, "replication:failover") {
		return
	}

	var failover = model.FailoverReplicationSpec{}
	if err := json.NewDecoder(this.Ctx.Request.Body).Decode(&failover); err != nil {
		model.HttpError(ctx, http.StatusBadRequest,
			"parse replication request body failed: %s", err.Error())
		return
	}

	id := ctx.Input.Param(":replicationId")
	r, err := db.C.GetReplication(c.GetContext(ctx), id)
	if err != nil {
		model.HttpError(ctx, http.StatusBadRequest,
			"get replication failed: %s", err.Error())
		return
	}

	if err := FailoverReplicationDBEntry(c.GetContext(ctx), r, failover.SecondaryBackendId); err != nil {
		model.HttpError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	// Call global controller variable to handle delete r request.
	err = controller.Brain.FailoverReplication(c.GetContext(ctx), r, &failover)
	if err != nil {
		model.HttpError(ctx, http.StatusBadRequest,
			"failover replication failed: %v", err.Error())
		return
	}

	ctx.Output.SetStatus(StatusAccepted)
}
