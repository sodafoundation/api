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
	"errors"
	"strconv"
	"strings"
	"time"

	log "github.com/golang/glog"
	"github.com/opensds/opensds/pkg/dock/client"
	"github.com/opensds/opensds/pkg/model"
	pb "github.com/opensds/opensds/pkg/model/proto"
)

type IntervalSnapshotExecutor struct {
	client.Client

	Request  *pb.CreateVolumeSnapshotOpts
	DockInfo *model.DockSpec
	Interval string
	TotalNum int
}

func (ise *IntervalSnapshotExecutor) Init(in string) (err error) {
	var volumeResponse model.VolumeSpec
	if err = json.Unmarshal([]byte(in), &volumeResponse); err != nil {
		return err
	}

	ise.Request.VolumeId = volumeResponse.Id
	ise.Request.Name = "snapshot-" + volumeResponse.Id
	ise.Request.Size = volumeResponse.Size
	ise.Client = client.NewClient()
	ise.Client.Connect(ise.DockInfo.Endpoint)

	return nil
}

func (ise *IntervalSnapshotExecutor) Asynchronized() error {
	if ise.TotalNum == 0 {
		ise.TotalNum = 3
	}
	num, err := ParseInterval(ise.Interval)
	if err != nil {
		log.Error("When parse snapshot interval:", err)
		return err
	}

	for i := 0; i < ise.TotalNum; i++ {
		// Sleep interval time
		for j := 0; j < num; j++ {
			time.Sleep(time.Second)
		}
		if _, err = ise.Client.CreateVolumeSnapshot(context.Background(), ise.Request); err != nil {
			log.Errorf("When %dth create volume snapshot: %v\n", i+1, err)
			return err
		}
	}
	return nil
}

func ParseInterval(interval string) (int, error) {
	var times int
	unit := strings.ToLower(interval[len(interval)-1:])
	if unit != "s" && unit != "m" && unit != "h" && unit != "d" {
		return 0, errors.New("interval unit is not correct")
	}
	switch unit {
	case "s":
		times = 1
	case "m":
		times = 60
	case "h":
		times = 60 * 60
	case "d":
		times = 60 * 60 * 24
	default:
		return 0, errors.New("interval unit is not correct")
	}

	num, err := strconv.Atoi(interval[0 : len(interval)-1])
	if err != nil {
		return 0, err
	}
	return num * times, nil
}
