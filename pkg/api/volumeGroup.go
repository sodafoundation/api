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

	c "github.com/opensds/opensds/pkg/context"
	//"github.com/opensds/opensds/pkg/controller"
	"github.com/opensds/opensds/pkg/api/policy"
	"github.com/opensds/opensds/pkg/db"
	"github.com/opensds/opensds/pkg/model"
)

type VolumeGroupPortal struct {
	BasePortal
}

func (v *VolumeGroupPortal) CreateVolumeGroup() {
	if !policy.Authorize(v.Ctx, "volume_group:create") {
		return
	}

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
	result, err := CreateVolumeGroupDBEntry(c.GetContext(v.Ctx), volumeGroup)
	if err != nil {
		v.ErrorHandle("Create volume group failed", model.ErrorInternalServer, err)
		return
	}

	// Marshal the result.
	body, err := json.Marshal(result)
	if err != nil {
		v.ErrorHandle("Marshal profile created result failed", model.ErrorBadRequest, err)
		return
	}

	v.SuccessHandle(StatusOK, body)
	return
}

func (v *VolumeGroupPortal) UpdateVolumeGroup() {
	if !policy.Authorize(v.Ctx, "volume_group:update") {
		return
	}
	var vg = &model.VolumeGroupSpec{
		BaseModel: &model.BaseModel{},
	}

	id := v.Ctx.Input.Param(":groupId")
	if err := json.NewDecoder(v.Ctx.Request.Body).Decode(&vg); err != nil {
		v.ErrorHandle("Parse volume group request body failed", model.ErrorBadRequest, err)
		return
	}

	vg.Id = id

	result, err := UpdateVolumeGroupDBEntry(c.GetContext(v.Ctx), vg)
	if err != nil {
		v.ErrorHandle("Update volume group failed", model.ErrorInternalServer, err)
		return
	}
	// Marshal the result.
	body, err := json.Marshal(result)
	if err != nil {
		v.ErrorHandle("Marshal volume group updated result failed", model.ErrorInternalServer, err)
		return
	}

	v.SuccessHandle(StatusOK, body)
	return
}

func (v *VolumeGroupPortal) DeleteVolumeGroup() {
	if !policy.Authorize(v.Ctx, "volume_group:delete") {
		return
	}

	err := DeleteVolumeGroupDBEntry(c.GetContext(v.Ctx), v.Ctx.Input.Param(":groupId"))
	if err != nil {
		v.ErrorHandle("Delete volume group failed", model.ErrorInternalServer, err)
		return
	}
	v.SuccessHandle(StatusAccepted, nil)
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
