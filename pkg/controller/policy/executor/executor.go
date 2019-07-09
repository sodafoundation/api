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
	"errors"

	log "github.com/golang/glog"
	"github.com/opensds/opensds/pkg/model"
	pb "github.com/opensds/opensds/pkg/model/proto"
)

// AsynchronizedExecutor
type AsynchronizedExecutor interface {
	Init(in string) error
	Asynchronized() error
}

// AsynchronizedWorkflow
type AsynchronizedWorkflow map[string]AsynchronizedExecutor

// RegisterAsynchronizedWorkflow
func RegisterAsynchronizedWorkflow(
	req interface{},
	tags map[string]string,
	dockInfo *model.DockSpec,
	in string) (
	AsynchronizedWorkflow, error) {

	var asynWorkflow = AsynchronizedWorkflow{}
	for key := range tags {
		switch key {
		case "intervalSnapshot":
			ise := &IntervalSnapshotExecutor{
				Request:  req.(*pb.CreateVolumeSnapshotOpts),
				Interval: tags[key],
				DockInfo: dockInfo,
			}

			if err := ise.Init(in); err != nil {
				log.Errorf("When register async policy %s: %v\n", key, err)
				return asynWorkflow, err
			}
			asynWorkflow[key] = ise

		case "deleteSnapshotPolicy":
			ise := &DeleteSnapshotExecutor{
				Request:  req.(*pb.DeleteVolumeSnapshotOpts),
				DockInfo: dockInfo,
			}

			if err := ise.Init(in); err != nil {
				log.Errorf("When register async policy %s: %v\n", key, err)
				return asynWorkflow, err
			}
			asynWorkflow[key] = ise
		}
	}

	log.Info("Register asynchronized work flow success, awf =", asynWorkflow)
	return asynWorkflow, nil
}

// ExecuteAsynchronizedWorkflow
func ExecuteAsynchronizedWorkflow(asynWorkflow AsynchronizedWorkflow) error {
	for key := range asynWorkflow {
		if asynWorkflow[key] == nil {
			return errors.New("Could not execute the policy " + key)
		}
		return asynWorkflow[key].Asynchronized()
	}
	return nil
}

// SynchronizedExecutor
type SynchronizedExecutor interface {
	Init() error
	Synchronized() error
}

// SynchronizedWorkflow
type SynchronizedWorkflow map[string]SynchronizedExecutor

// RegisterSynchronizedWorkflow
func RegisterSynchronizedWorkflow(req interface{}, tags map[string]interface{}) (SynchronizedWorkflow, error) {
	return SynchronizedWorkflow{}, nil
}

// ExecuteSynchronizedWorkflow
func ExecuteSynchronizedWorkflow(synWorkflow SynchronizedWorkflow) error {
	for key := range synWorkflow {
		if synWorkflow[key] == nil {
			return errors.New("Could not execute the policy " + key)
		}
		if err := synWorkflow[key].Synchronized(); err != nil {
			return err
		}
	}
	return nil
}
