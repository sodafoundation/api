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

package drbd

import (
	log "github.com/golang/glog"
	pb "github.com/opensds/opensds/pkg/dock/proto"
	"github.com/opensds/opensds/pkg/model"
)

// ReplicationDriver
type ReplicationDriver struct{}

// Setup
func (r *ReplicationDriver) Setup() error { return nil }

// Unset
func (r *ReplicationDriver) Unset() error { return nil }

// CreateReplication
func (r *ReplicationDriver) CreateReplication(opt *pb.CreateReplicationOpts) (*model.ReplicationSpec, error) {
	log.Infof("DRBD create replication ....")
	return &model.ReplicationSpec{
		PrimaryReplicationDriverData:   map[string]string{"primary-key1": "test1"},
		SecondaryReplicationDriverData: map[string]string{"secondary-key1": "test2"},
		Metadata:                       map[string]string{"meta-key1": "test2"},
	}, nil
}

func (r *ReplicationDriver) DeleteReplication(opt *pb.DeleteReplicationOpts) error {
	log.Infof("DRBD delete replication ....")
	return nil
}

func (r *ReplicationDriver) EnableReplication(opt *pb.EnableReplicationOpts) error {
	log.Infof("DRBD enable replication ....")
	return nil
}

func (r *ReplicationDriver) DisableReplication(opt *pb.DisableReplicationOpts) error {
	log.Infof("DRBD disable replication ....")
	return nil
}

func (r *ReplicationDriver) FailoverReplication(opt *pb.FailoverReplicationOpts) error {
	log.Infof("DRBD failover replication ....")
	return nil
}
