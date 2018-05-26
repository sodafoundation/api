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

/*
This module implements a entry into the OpenSDS northbound service.

*/

package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/astaxie/beego"
	log "github.com/golang/glog"
	"github.com/opensds/opensds/contrib/cindercompatibleapi/converter"
	"github.com/opensds/opensds/pkg/model"
)

// SnapshotPortal ...
type SnapshotPortal struct {
	beego.Controller
}

// ListSnapshotsDetails ...
func (portal *SnapshotPortal) ListSnapshotsDetails() {
	snapshots, err := client.ListVolumeSnapshots()
	if err != nil {
		reason := fmt.Sprintf("List snapshots and details failed: %v", err)
		portal.Ctx.Output.SetStatus(model.ErrorInternalServer)
		portal.Ctx.Output.Body(model.ErrorInternalServerStatus(reason))
		log.Error(reason)
		return
	}

	result := converter.ListSnapshotsDetailsResp(snapshots)
	body, err := json.Marshal(result)
	if err != nil {
		reason := fmt.Sprintf("List snapshots and details, marshal result failed: %v", err)
		portal.Ctx.Output.SetStatus(model.ErrorInternalServer)
		portal.Ctx.Output.Body(model.ErrorInternalServerStatus(reason))
		log.Error(reason)
		return
	}

	portal.Ctx.Output.SetStatus(http.StatusOK)
	portal.Ctx.Output.Body(body)
	return
}

// CreateSnapshot ...
func (portal *SnapshotPortal) CreateSnapshot() {
	var cinderReq = converter.CreateSnapshotReqSpec{}
	if err := json.NewDecoder(portal.Ctx.Request.Body).Decode(&cinderReq); err != nil {
		reason := fmt.Sprintf("Create a snapshot, parse request body failed: %s", err.Error())
		portal.Ctx.Output.SetStatus(model.ErrorBadRequest)
		portal.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	snapshot, err := converter.CreateSnapshotReq(&cinderReq)
	if err != nil {
		reason := fmt.Sprintf("Create a snapshot failed: %s", err.Error())
		portal.Ctx.Output.SetStatus(model.ErrorBadRequest)
		portal.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	snapshot, err = client.CreateVolumeSnapshot(snapshot)
	if err != nil {
		reason := fmt.Sprintf("Create a snapshot failed: %s", err.Error())
		portal.Ctx.Output.SetStatus(model.ErrorInternalServer)
		portal.Ctx.Output.Body(model.ErrorInternalServerStatus(reason))
		log.Error(reason)
		return
	}

	result := converter.CreateSnapshotResp(snapshot)
	body, err := json.Marshal(result)
	if err != nil {
		reason := fmt.Sprintf("Create a snapshot, marshal result failed: %s", err.Error())
		portal.Ctx.Output.SetStatus(model.ErrorInternalServer)
		portal.Ctx.Output.Body(model.ErrorInternalServerStatus(reason))
		log.Error(reason)
		return
	}

	portal.Ctx.Output.SetStatus(http.StatusAccepted)
	portal.Ctx.Output.Body(body)
	return
}

// ListSnapshots ...
func (portal *SnapshotPortal) ListSnapshots() {
	snapshots, err := client.ListVolumeSnapshots()
	if err != nil {
		reason := fmt.Sprintf("List accessible snapshots failed: %v", err)
		portal.Ctx.Output.SetStatus(model.ErrorInternalServer)
		portal.Ctx.Output.Body(model.ErrorInternalServerStatus(reason))
		log.Error(reason)
		return
	}

	result := converter.ListSnapshotsResp(snapshots)
	body, err := json.Marshal(result)
	if err != nil {
		reason := fmt.Sprintf("List accessible snapshots, marshal result failed: %v", err)
		portal.Ctx.Output.SetStatus(model.ErrorInternalServer)
		portal.Ctx.Output.Body(model.ErrorInternalServerStatus(reason))
		log.Error(reason)
		return
	}

	portal.Ctx.Output.SetStatus(http.StatusOK)
	portal.Ctx.Output.Body(body)
	return
}

// GetSnapshot ...
func (portal *SnapshotPortal) GetSnapshot() {
	id := portal.Ctx.Input.Param(":snapshotId")
	snapshot, err := client.GetVolumeSnapshot(id)

	if err != nil {
		reason := fmt.Sprintf("Show a snapshot's details failed: %v", err)
		portal.Ctx.Output.SetStatus(model.ErrorInternalServer)
		portal.Ctx.Output.Body(model.ErrorInternalServerStatus(reason))
		log.Error(reason)
		return
	}

	result := converter.ShowSnapshotDetailsResp(snapshot)
	body, err := json.Marshal(result)
	if err != nil {
		reason := fmt.Sprintf("Show a snapshot's details, marshal result failed: %v", err)
		portal.Ctx.Output.SetStatus(model.ErrorInternalServer)
		portal.Ctx.Output.Body(model.ErrorInternalServerStatus(reason))
		log.Error(reason)
		return
	}

	portal.Ctx.Output.SetStatus(http.StatusOK)
	portal.Ctx.Output.Body(body)
	return
}

// UpdateSnapshot ...
func (portal *SnapshotPortal) UpdateSnapshot() {
	id := portal.Ctx.Input.Param(":snapshotId")
	var cinderUpdateReq = converter.UpdateSnapshotReqSpec{}

	if err := json.NewDecoder(portal.Ctx.Request.Body).Decode(&cinderUpdateReq); err != nil {
		reason := fmt.Sprintf("Update a snapshot, parse request body failed: %s", err.Error())
		portal.Ctx.Output.SetStatus(model.ErrorBadRequest)
		portal.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	snapshot := converter.UpdateSnapshotReq(&cinderUpdateReq)
	snapshot, err := client.UpdateVolumeSnapshot(id, snapshot)

	if err != nil {
		reason := fmt.Sprintf("Update a snapshot failed: %s", err.Error())
		portal.Ctx.Output.SetStatus(model.ErrorInternalServer)
		portal.Ctx.Output.Body(model.ErrorInternalServerStatus(reason))
		log.Error(reason)
		return
	}

	result := converter.UpdateSnapshotResp(snapshot)
	body, err := json.Marshal(result)
	if err != nil {
		reason := fmt.Sprintf("Update a snapshot, marshal result failed: %s", err.Error())
		portal.Ctx.Output.SetStatus(model.ErrorInternalServer)
		portal.Ctx.Output.Body(model.ErrorInternalServerStatus(reason))
		log.Error(reason)
		return
	}

	portal.Ctx.Output.SetStatus(http.StatusOK)
	portal.Ctx.Output.Body(body)
	return
}

// DeleteSnapshot ...
func (portal *SnapshotPortal) DeleteSnapshot() {
	id := portal.Ctx.Input.Param(":snapshotId")
	err := client.DeleteVolumeSnapshot(id, nil)

	if err != nil {
		reason := fmt.Sprintf("Delete a snapshot failed: %v", err)
		portal.Ctx.Output.SetStatus(model.ErrorInternalServer)
		portal.Ctx.Output.Body(model.ErrorInternalServerStatus(reason))
		log.Error(reason)
		return
	}

	portal.Ctx.Output.SetStatus(http.StatusAccepted)
	return
}
