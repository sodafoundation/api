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

func NewVolumeGroupPortal() *VolumeGroupPortal {
	return &VolumeGroupPortal{
		CtrClient: client.NewClient(),
	}
}

type VolumeGroupPortal struct {
	BasePortal

	CtrClient client.Client
}

func (this *VolumeGroupPortal) CreateVolumeGroup() {
	if !policy.Authorize(this.Ctx, "volume_group:create") {
		return
	}
	ctx := c.GetContext(this.Ctx)

	var volumeGroup = &model.VolumeGroupSpec{
		BaseModel: &model.BaseModel{},
	}

	// Unmarshal the request body
	if err := json.NewDecoder(this.Ctx.Request.Body).Decode(&volumeGroup); err != nil {
		this.ErrorHandle("Parse volume group request body failed", model.ErrorBadRequest, err)
		return
	}
	// NOTE:It will create a volume group entry into the database and initialize its status
	// as "creating". It will not wait for the real volume group process creation to complete
	// and will return result immediately.
	result, err := CreateVolumeGroupDBEntry(ctx, volumeGroup)
	if err != nil {
		this.ErrorHandle("Create volume group failed",
			model.ErrorInternalServer, err)
		return
	}

	// Marshal the result.
	body, err := json.Marshal(result)
	if err != nil {
		this.ErrorHandle("Marshal volume group created result failed",
			model.ErrorInternalServer, err)
		return
	}
	this.SuccessHandle(StatusOK, body)

	// NOTE:The real volume group creation process.
	// Volume group creation request is sent to the Dock. Dock will set
	// volume group status to 'available' after volume group creation operation
	// is completed.
	if err = this.CtrClient.Connect(CONF.OsdsLet.ApiEndpoint); err != nil {
		log.Error("When connecting controller client:", err)
		return
	}
	defer this.CtrClient.Close()

	opt := &pb.CreateVolumeGroupOpts{
		Message: string(body),
		Context: ctx.ToJson(),
	}
	if _, err = this.CtrClient.CreateVolumeGroup(context.Background(), opt); err != nil {
		log.Error("Create volume group failed in controller service:", err)
		return
	}

	return
}

func (this *VolumeGroupPortal) UpdateVolumeGroup() {
	if !policy.Authorize(this.Ctx, "volume_group:update") {
		return
	}
	ctx := c.GetContext(this.Ctx)
	var vg = &model.VolumeGroupSpec{
		BaseModel: &model.BaseModel{},
	}

	id := this.Ctx.Input.Param(":groupId")
	if err := json.NewDecoder(this.Ctx.Request.Body).Decode(&vg); err != nil {
		this.ErrorHandle("Parse volume group request body failed", model.ErrorBadRequest, err)
		return
	}

	vg.Id = id

	result, err := UpdateVolumeGroupDBEntry(ctx, vg)
	if err != nil {
		this.ErrorHandle("Update volume group failed", model.ErrorInternalServer, err)
		return
	}
	// Marshal the result.
	body, err := json.Marshal(result)
	if err != nil {
		this.ErrorHandle("Marshal volume group updated result failed",
			model.ErrorInternalServer, err)
		return
	}
	this.SuccessHandle(StatusOK, body)

	// NOTE:The real volume group creation process.
	// Volume group creation request is sent to the Dock. Dock will set
	// volume group status to 'available' after volume group creation operation
	// is completed.
	if err = this.CtrClient.Connect(CONF.OsdsLet.ApiEndpoint); err != nil {
		log.Error("When connecting controller client:", err)
		return
	}
	defer this.CtrClient.Close()

	opt := &pb.CreateVolumeGroupOpts{
		Message: string(body),
		Context: ctx.ToJson(),
	}
	if _, err = this.CtrClient.CreateVolumeGroup(context.Background(), opt); err != nil {
		log.Error("Create volume group failed in controller service:", err)
		return
	}

	return
}

func (this *VolumeGroupPortal) DeleteVolumeGroup() {
	if !policy.Authorize(this.Ctx, "volume_group:delete") {
		return
	}
	ctx := c.GetContext(this.Ctx)

	id := this.Ctx.Input.Param(":groupId")
	vg, err := db.C.GetVolumeGroup(ctx, id)
	if err != nil {
		this.ErrorHandle("Delete volume group failed",
			model.ErrorBadRequest, err)
		return
	}

	if err = DeleteVolumeGroupDBEntry(c.GetContext(this.Ctx), id); err != nil {
		this.ErrorHandle("Delete volume group failed",
			model.ErrorInternalServer, err)
		return
	}

	this.SuccessHandle(StatusAccepted, nil)

	// NOTE:The real volume group deletion process.
	// Volume group deletion request is sent to the Dock. Dock will remove
	// volume group record after volume group deletion operation is completed.
	if err = this.CtrClient.Connect(CONF.OsdsLet.ApiEndpoint); err != nil {
		log.Error("When connecting controller client:", err)
		return
	}
	defer this.CtrClient.Close()

	body, _ := json.Marshal(vg)
	opt := &pb.DeleteVolumeGroupOpts{
		Message: string(body),
		Context: ctx.ToJson(),
	}
	if _, err = this.CtrClient.DeleteVolumeGroup(context.Background(), opt); err != nil {
		log.Error("Delete volume group failed in controller service:", err)
		return
	}

	return
}

func (this *VolumeGroupPortal) GetVolumeGroup() {
	if !policy.Authorize(this.Ctx, "volume_group:get") {
		return
	}

	// Call db api module to handle get volume request.
	result, err := db.C.GetVolumeGroup(c.GetContext(this.Ctx), this.Ctx.Input.Param(":groupId"))
	if err != nil {
		this.ErrorHandle("Get volume group failed", model.ErrorBadRequest, err)
		return
	}

	// Marshal the result.
	body, err := json.Marshal(result)
	if err != nil {
		this.ErrorHandle("Marshal volume group showed result failed", model.ErrorInternalServer, err)
		return
	}

	this.SuccessHandle(StatusOK, body)
	return
}

func (this *VolumeGroupPortal) ListVolumeGroups() {
	if !policy.Authorize(this.Ctx, "volume_group:get") {
		return
	}

	m, err := this.GetParameters()
	if err != nil {
		this.ErrorHandle("List volume groups failed", model.ErrorBadRequest, err)
		return
	}

	result, err := db.C.ListVolumeGroupsWithFilter(c.GetContext(this.Ctx), m)
	if err != nil {
		this.ErrorHandle("List volume groups failed", model.ErrorBadRequest, err)
		return
	}

	// Marshal the result.
	body, err := json.Marshal(result)
	if err != nil {
		this.ErrorHandle("Marshal volume groups listed result failed", model.ErrorInternalServer, err)
		return
	}

	this.SuccessHandle(StatusOK, body)
	return
}
