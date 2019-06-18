// Copyright 2017 The OpenSDS Authors.
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
This module implements the policy-based scheduling by parsing storage
profiles configured by admin.

*/

package executor

import (
	"context"
	"encoding/json"
	"time"

	log "github.com/golang/glog"
	c "github.com/opensds/opensds/pkg/context"
	"github.com/opensds/opensds/pkg/db"
	"github.com/opensds/opensds/pkg/dock/client"
	"github.com/opensds/opensds/pkg/model"
	pb "github.com/opensds/opensds/pkg/model/proto"
)

const (
	RETRY_INTERVAL = 5
	MAX_RETRY_TIME = 10
)

type DeleteSnapshotExecutor struct {
	client.Client

	VolumeId string
	Request  *pb.DeleteVolumeSnapshotOpts
	DockInfo *model.DockSpec
}

func (dse *DeleteSnapshotExecutor) Init(in string) (err error) {
	var volumeResponse model.VolumeSpec
	if err = json.Unmarshal([]byte(in), &volumeResponse); err != nil {
		return err
	}
	dse.VolumeId = volumeResponse.Id
	dse.Client = client.NewClient()
	dse.Client.Connect(dse.DockInfo.Endpoint)

	return nil
}

func (dse *DeleteSnapshotExecutor) Asynchronized() error {
	remainSnaps, err := findRemainingSnapshot(dse.VolumeId)
	if err != nil {
		return err
	}

	for i, snapId := range remainSnaps {
		dse.Request.Id = snapId
		if _, err = dse.Client.DeleteVolumeSnapshot(context.Background(), dse.Request); err != nil {
			log.Errorf("When %dth delete volume snapshot: %v\n", i+1, err)
			return err
		}
	}
	// Waiting for snapshots deleted
	for i := 0; i < MAX_RETRY_TIME; i++ {
		if CheckSnapshotDeleted(dse.VolumeId) {
			break
		}
		time.Sleep(RETRY_INTERVAL * time.Second)
	}
	return nil
}

func CheckSnapshotDeleted(volumeId string) bool {
	snaps, err := findRemainingSnapshot(volumeId)
	if err != nil {
		return false
	}

	if len(snaps) == 0 {
		return true
	}
	return false
}

func findRemainingSnapshot(volumeId string) ([]string, error) {
	var remainingSnapshots = []string{}
	snapshots, err := db.C.ListVolumeSnapshots(c.NewAdminContext())
	if err != nil {
		log.Error("When list volume snapshots:", err)
		return remainingSnapshots, err
	}

	for _, snap := range snapshots {
		if snap.VolumeId != volumeId {
			continue
		}
		remainingSnapshots = append(remainingSnapshots, snap.Id)
	}
	return remainingSnapshots, nil
}
