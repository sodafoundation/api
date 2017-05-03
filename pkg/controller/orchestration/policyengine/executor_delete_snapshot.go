// Copyright (c) 2016 Huawei Technologies Co., Ltd. All Rights Reserved.
//
//    Licensed under the Apache License, Version 2.0 (the "License"); you may
//    not use this file except in compliance with the License. You may obtain
//    a copy of the License at
//
//         http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//    WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//    License for the specific language governing permissions and limitations
//    under the License.

/*
This module implements the policy-based scheduling by parsing storage
profiles configured by admin.

*/

package policyengine

import (
	"encoding/json"
	"log"
	"time"

	api "github.com/opensds/opensds/pkg/api/v1"
	"github.com/opensds/opensds/pkg/grpc/dock/client"
	pb "github.com/opensds/opensds/pkg/grpc/opensds"
)

const (
	RETRY_INTERVAL = 3
	MAX_RETRY_TIME = 5
)

type DeleteSnapshotExecutor struct {
	Request *pb.VolumeRequest
}

func (dse *DeleteSnapshotExecutor) Init(in string) error {
	return nil
}

func (dse *DeleteSnapshotExecutor) Asynchronized() error {
	remainSnaps, err := findRemainingSnapshot(dse.Request)
	if err != nil {
		return err
	}

	for i, snapId := range remainSnaps {
		dse.Request.SnapshotId = snapId
		if _, err = client.DeleteVolumeSnapshot(dse.Request); err != nil {
			log.Printf("[Error] When %dth delete volume snapshot: %v\n", i+1, err)
			return err
		}
	}
	// Waiting for snapshots deleted
	for i := 0; i < MAX_RETRY_TIME; i++ {
		if CheckSnapshotDeleted(dse.Request) {
			break
		}
		time.Sleep(RETRY_INTERVAL * time.Second)
	}
	return nil
}

func CheckSnapshotDeleted(vr *pb.VolumeRequest) bool {
	snaps, err := findRemainingSnapshot(vr)
	if err != nil {
		return false
	}

	if len(snaps) == 0 {
		return true
	}
	return false
}

func findRemainingSnapshot(vr *pb.VolumeRequest) ([]string, error) {
	var remainingSnapshots = []string{}
	var snapshots []api.VolumeSnapshotResponse
	result, err := client.ListVolumeSnapshots(vr)
	if err != nil {
		log.Println("[Error] When list volume snapshots:", err)
		return remainingSnapshots, err
	}
	if err = json.Unmarshal([]byte(result.GetMessage()), &snapshots); err != nil {
		log.Println("[Error] When parse volume snapshots:", err)
		return remainingSnapshots, err
	}

	for _, snap := range snapshots {
		if snap.VolumeId != vr.GetId() {
			continue
		}
		remainingSnapshots = append(remainingSnapshots, snap.Id)
	}
	return remainingSnapshots, nil
}
