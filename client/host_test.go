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

var fakeHostMgr = &HostMgr{
	Receiver: NewFakeHostReceiver(),
}

func TestCreateHost(t *testing.T) {
	expected := &model.HostSpec{
		BaseModel: &model.BaseModel{
			Id:        "202964b5-8e73-46fd-b41b-a8e403f3c30b",
			CreatedAt: "2019-11-11T11:01:33",
		},
		TenantId:          "x",
		AccessMode:        "agentless",
		HostName:          "sap1",
		IP:                "192.168.56.12",
		AvailabilityZones: []string{"az1", "az2"},
		Initiators: []*model.Initiator{
			&model.Initiator{
				PortName: "20000024ff5bb888",
				Protocol: "iscsi",
			},
			&model.Initiator{
				PortName: "20000024ff5bc999",
				Protocol: "iscsi",
			},
		},
	}

	host, err := fakeHostMgr.CreateHost(&model.HostSpec{})
	if err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(host, expected) {
		t.Errorf("expected %+v, got %+v", expected, host)
		return
	}
}

func TestGetHost(t *testing.T) {
	hostID := "d2975ebe-d82c-430f-b28e-f373746a71ca"
	expected := &model.HostSpec{
		BaseModel: &model.BaseModel{
			Id:        "202964b5-8e73-46fd-b41b-a8e403f3c30b",
			CreatedAt: "2019-11-11T11:01:33",
		},
		TenantId:          "x",
		AccessMode:        "agentless",
		HostName:          "sap1",
		IP:                "192.168.56.12",
		AvailabilityZones: []string{"az1", "az2"},
		Initiators: []*model.Initiator{
			&model.Initiator{
				PortName: "20000024ff5bb888",
				Protocol: "iscsi",
			},
			&model.Initiator{
				PortName: "20000024ff5bc999",
				Protocol: "iscsi",
			},
		},
	}

	host, err := fakeHostMgr.GetHost(hostID)
	if err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(host, expected) {
		t.Errorf("expected %v, got %v", expected, host)
		return
	}
}

func TestListHosts(t *testing.T) {
	sampleHosts := []model.HostSpec{
		{
			BaseModel: &model.BaseModel{
				Id:        "202964b5-8e73-46fd-b41b-a8e403f3c30b",
				CreatedAt: "2019-11-11T11:01:33",
			},
			TenantId:          "x",
			AccessMode:        "agentless",
			HostName:          "sap1",
			IP:                "192.168.56.12",
			AvailabilityZones: []string{"az1", "az2"},
			Initiators: []*model.Initiator{
				&model.Initiator{
					PortName: "20000024ff5bb888",
					Protocol: "iscsi",
				},
				&model.Initiator{
					PortName: "20000024ff5bc999",
					Protocol: "iscsi",
				},
			},
		},
		{
			BaseModel: &model.BaseModel{
				Id:        "eb73e59a-8b0f-4517-8b95-023ec134aec9",
				CreatedAt: "2019-11-11T11:13:57",
			},
			TenantId:          "x",
			AccessMode:        "agentless",
			HostName:          "sap2",
			IP:                "192.168.56.13",
			AvailabilityZones: []string{"az1", "az2"},
			Initiators: []*model.Initiator{
				&model.Initiator{
					PortName: "20012324ff5ac132",
					Protocol: "iscsi",
				},
			},
		},
	}

	var expected []*model.HostSpec
	expected = append(expected, &sampleHosts[0])
	expected = append(expected, &sampleHosts[1])
	hosts, err := fakeHostMgr.ListHosts(map[string]string{"limit": "0", "offset": "10"})

	if err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(hosts, expected) {
		t.Errorf("expected %v, got %v", expected, hosts)
		return
	}
}

func TestUpdateHost(t *testing.T) {
	hostID := "202964b5-8e73-46fd-b41b-a8e403f3c30b"
	host := &model.HostSpec{
		HostName: "sap1-updated",
	}

	result, err := fakeHostMgr.UpdateHost(hostID, host)
	if err != nil {
		t.Error(err)
		return
	}

	expected := &model.HostSpec{
		BaseModel: &model.BaseModel{
			Id:        "202964b5-8e73-46fd-b41b-a8e403f3c30b",
			CreatedAt: "2019-11-11T11:01:33",
		},
		TenantId:          "x",
		AccessMode:        "agentless",
		HostName:          "sap1",
		IP:                "192.168.56.12",
		AvailabilityZones: []string{"az1", "az2"},
		Initiators: []*model.Initiator{
			&model.Initiator{
				PortName: "20000024ff5bb888",
				Protocol: "iscsi",
			},
			&model.Initiator{
				PortName: "20000024ff5bc999",
				Protocol: "iscsi",
			},
		},
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
		return
	}
}

func TestDeleteHost(t *testing.T) {
	var hostID = "d202964b5-8e73-46fd-b41b-a8e403f3c30b"

	if err := fakeHostMgr.DeleteHost(hostID); err != nil {
		t.Error(err)
		return
	}
}
