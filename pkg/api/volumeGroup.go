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

func (v *VolumeGroupPortal) CreateVolumeGroup() {
	if !policy.Authorize(v.Ctx, "volume_group:create") {
		return
	}
	ctx := c.GetContext(v.Ctx)

	var volumeGroup = &model.VolumeGroupSpec{
		BaseModel: &model.BaseModel{},
	}

	// Unmarshal the request body
	if err := json.NewDecoder(v.Ctx.Request.Body).Decode(&volumeGroup); err != nil {
		v.ErrorHandle("Parse volume group request body failed", model.ErrorBadRequest, err)
		return
	}
	// NOTE:It will create a volume group entry into the database and initialize its status
	// as "creating". It will not wait for the real volume group process creation to complete
	// and will return result immediately.
	result, err := CreateVolumeGroupDBEntry(ctx, volumeGroup)
	if err != nil {
		v.ErrorHandle("Create volume group failed",
			model.ErrorInternalServer, err)
		return
	}

	// Marshal the result.
	body, err := json.Marshal(result)
	if err != nil {
		v.ErrorHandle("Marshal volume group created result failed",
			model.ErrorInternalServer, err)
		return
	}
	v.SuccessHandle(StatusOK, body)

	// NOTE:The real volume group creation process.
	// Volume group creation request is sent to the Dock. Dock will set
	// volume group status to 'available' after volume group creation operation
	// is completed.
	if err = v.CtrClient.Connect(CONF.OsdsLet.ApiEndpoint); err != nil {
		log.Error("When connecting controller client:", err)
		return
	}
	defer v.CtrClient.Close()

	opt := &pb.CreateVolumeGroupOpts{
		Message: string(body),
		Context: ctx.ToJson(),
	}
	if _, err = v.CtrClient.CreateVolumeGroup(context.Background(), opt); err != nil {
		log.Error("Create volume group failed in controller service:", err)
		return
	}

	return
}

func (v *VolumeGroupPortal) UpdateVolumeGroup() {
	if !policy.Authorize(v.Ctx, "volume_group:update") {
		return
	}
	ctx := c.GetContext(v.Ctx)
	var vg = &model.VolumeGroupSpec{
		BaseModel: &model.BaseModel{},
	}

	id := v.Ctx.Input.Param(":groupId")
	if err := json.NewDecoder(v.Ctx.Request.Body).Decode(&vg); err != nil {
		v.ErrorHandle("Parse volume group request body failed", model.ErrorBadRequest, err)
		return
	}

	vg.Id = id
	result, err := UpdateVolumeGroupDBEntry(ctx, vg)
	if err != nil {
		v.ErrorHandle("Update volume group failed", model.ErrorInternalServer, err)
		return
	}
	// Marshal the result.
	body, err := json.Marshal(result)
	if err != nil {
		v.ErrorHandle("Marshal volume group updated result failed",
			model.ErrorInternalServer, err)
		return
	}
	v.SuccessHandle(StatusOK, body)

	// NOTE:The real volume group creation process.
	// Volume group creation request is sent to the Dock. Dock will set
	// volume group status to 'available' after volume group creation operation
	// is completed.
	if err = v.CtrClient.Connect(CONF.OsdsLet.ApiEndpoint); err != nil {
		log.Error("When connecting controller client:", err)
		return
	}
	defer v.CtrClient.Close()

	opt := &pb.CreateVolumeGroupOpts{
		Message: string(body),
		Context: ctx.ToJson(),
	}
	if _, err = v.CtrClient.CreateVolumeGroup(context.Background(), opt); err != nil {
		log.Error("Create volume group failed in controller service:", err)
		return
	}

	return
}

func (v *VolumeGroupPortal) DeleteVolumeGroup() {
	if !policy.Authorize(v.Ctx, "volume_group:delete") {
		return
	}
	ctx := c.GetContext(v.Ctx)

	id := v.Ctx.Input.Param(":groupId")
	vg, err := db.C.GetVolumeGroup(ctx, id)
	if err != nil {
		v.ErrorHandle("Delete volume group failed",
			model.ErrorBadRequest, err)
		return
	}

	if err = DeleteVolumeGroupDBEntry(c.GetContext(v.Ctx), id); err != nil {
		v.ErrorHandle("Delete volume group failed",
			model.ErrorInternalServer, err)
		return
	}

	v.SuccessHandle(StatusAccepted, nil)

	// NOTE:The real volume group deletion process.
	// Volume group deletion request is sent to the Dock. Dock will remove
	// volume group record after volume group deletion operation is completed.
	if err = v.CtrClient.Connect(CONF.OsdsLet.ApiEndpoint); err != nil {
		log.Error("When connecting controller client:", err)
		return
	}
	defer v.CtrClient.Close()

	body, _ := json.Marshal(vg)
	opt := &pb.DeleteVolumeGroupOpts{
		Message: string(body),
		Context: ctx.ToJson(),
	}
	if _, err = v.CtrClient.DeleteVolumeGroup(context.Background(), opt); err != nil {
		log.Error("Delete volume group failed in controller service:", err)
		return
	}

	return
}

func (v *VolumeGroupPortal) GetVolumeGroup() {
	if !policy.Authorize(v.Ctx, "volume_group:get") {
		return
	}

	// Call db api module to handle get volume request.
	result, err := db.C.GetVolumeGroup(c.GetContext(v.Ctx), v.Ctx.Input.Param(":groupId"))
	if err != nil {
		v.ErrorHandle("Get volume group failed", model.ErrorBadRequest, err)
		return
	}

	// Marshal the result.
	body, err := json.Marshal(result)
	if err != nil {
		v.ErrorHandle("Marshal volume group showed result failed", model.ErrorInternalServer, err)
		return
	}

	v.SuccessHandle(StatusOK, body)
	return
}

func (v *VolumeGroupPortal) ListVolumeGroups() {
	if !policy.Authorize(v.Ctx, "volume_group:get") {
		return
	}

	m, err := v.GetParameters()
	if err != nil {
		v.ErrorHandle("List volume groups failed", model.ErrorBadRequest, err)
		return
	}

	result, err := db.C.ListVolumeGroupsWithFilter(c.GetContext(v.Ctx), m)
	if err != nil {
		v.ErrorHandle("List volume groups failed", model.ErrorBadRequest, err)
		return
	}

	// Marshal the result.
	body, err := json.Marshal(result)
	if err != nil {
		v.ErrorHandle("Marshal volume groups listed result failed", model.ErrorInternalServer, err)
		return
	}

	v.SuccessHandle(StatusOK, body)
	return
}
