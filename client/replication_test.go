// Copyright 2019 The OpenSDS Authors.
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

package client

import (
	"reflect"
	"testing"

	"github.com/opensds/opensds/pkg/model"
)

var fr = &ReplicationMgr{
	Receiver: NewFakeReplicationReceiver(),
}

func TestCreateReplication(t *testing.T) {
	expected := &model.ReplicationSpec{
		BaseModel: &model.BaseModel{
			Id: "c299a978-4f3e-11e8-8a5c-977218a83359",
		},
		PrimaryVolumeId:   "bd5b12a8-a101-11e7-941e-d77981b584d8",
		SecondaryVolumeId: "bd5b12a8-a101-11e7-941e-d77981b584d8",
		Name:              "sample-replication-01",
		Description:       "This is a sample replication for testing",
		PoolId:            "084bf71e-a102-11e7-88a8-e31fe6d52248",
		ProfileId:         "1106b972-66ef-11e7-b172-db03f3689c9c",
	}

	replica, err := fr.CreateReplication(&model.ReplicationSpec{})
	if err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(replica, expected) {
		t.Errorf("expected %v, got %v", expected, replica)
		return
	}
}

func TestGetReplication(t *testing.T) {
	var replicaID = "c299a978-4f3e-11e8-8a5c-977218a83359"
	expected := &model.ReplicationSpec{
		BaseModel: &model.BaseModel{
			Id: "c299a978-4f3e-11e8-8a5c-977218a83359",
		},
		PrimaryVolumeId:   "bd5b12a8-a101-11e7-941e-d77981b584d8",
		SecondaryVolumeId: "bd5b12a8-a101-11e7-941e-d77981b584d8",
		Name:              "sample-replication-01",
		Description:       "This is a sample replication for testing",
		PoolId:            "084bf71e-a102-11e7-88a8-e31fe6d52248",
		ProfileId:         "1106b972-66ef-11e7-b172-db03f3689c9c",
	}

	replica, err := fr.GetReplication(replicaID)
	if err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(replica, expected) {
		t.Errorf("expected %v, got %v", expected, replica)
		return
	}
}

func TestListReplications(t *testing.T) {
	expected := []*model.ReplicationSpec{
		{
			BaseModel: &model.BaseModel{
				Id: "c299a978-4f3e-11e8-8a5c-977218a83359",
			},
			PrimaryVolumeId:   "bd5b12a8-a101-11e7-941e-d77981b584d8",
			SecondaryVolumeId: "bd5b12a8-a101-11e7-941e-d77981b584d8",
			Name:              "sample-replication-01",
			Description:       "This is a sample replication for testing",
			PoolId:            "084bf71e-a102-11e7-88a8-e31fe6d52248",
			ProfileId:         "1106b972-66ef-11e7-b172-db03f3689c9c",
		},
		{
			BaseModel: &model.BaseModel{
				Id: "73bfdd58-4f3f-11e8-91c0-d39a05f391ee",
			},
			PrimaryVolumeId:   "bd5b12a8-a101-11e7-941e-d77981b584d8",
			SecondaryVolumeId: "bd5b12a8-a101-11e7-941e-d77981b584d8",
			Name:              "sample-replication-02",
			Description:       "This is a sample replication for testing",
			PoolId:            "084bf71e-a102-11e7-88a8-e31fe6d52248",
			ProfileId:         "1106b972-66ef-11e7-b172-db03f3689c9c",
		},
	}

	replicas, err := fr.ListReplications(map[string]string{"limit": "3", "offset": "4"})
	if err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(replicas, expected) {
		t.Errorf("expected %v, got %v", expected, replicas)
		return
	}
}

func TestDeleteReplication(t *testing.T) {
	var replicaID = "c299a978-4f3e-11e8-8a5c-977218a83359"

	if err := fr.DeleteReplication(replicaID, &model.ReplicationSpec{}); err != nil {
		t.Error(err)
		return
	}
}

func TestUpdateReplication(t *testing.T) {
	var replicaID = "c299a978-4f3e-11e8-8a5c-977218a83359"
	replica := &model.ReplicationSpec{
		Name:        "sample-replication-03",
		Description: "This is a sample replication for testing",
	}

	result, err := fr.UpdateReplication(replicaID, replica)
	if err != nil {
		t.Error(err)
		return
	}

	expected := &model.ReplicationSpec{
		BaseModel: &model.BaseModel{
			Id: "c299a978-4f3e-11e8-8a5c-977218a83359",
		},
		PrimaryVolumeId:   "bd5b12a8-a101-11e7-941e-d77981b584d8",
		SecondaryVolumeId: "bd5b12a8-a101-11e7-941e-d77981b584d8",
		Name:              "sample-replication-01",
		Description:       "This is a sample replication for testing",
		PoolId:            "084bf71e-a102-11e7-88a8-e31fe6d52248",
		ProfileId:         "1106b972-66ef-11e7-b172-db03f3689c9c",
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
		return
	}
}
