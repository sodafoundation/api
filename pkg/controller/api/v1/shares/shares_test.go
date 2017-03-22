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
This module implements the entry into CRUD operation of volumes.

*/

package shares

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/opensds/opensds/pkg/controller/api"
	pb "github.com/opensds/opensds/pkg/grpc/opensds"
)

type FakeShareRequest struct {
	DockId       string `json:"dockId,omitempty"`
	ResourceType string `json:"resourceType,omitempty"`
	Id           string `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	Size         int    `json:"size"`
	ShareType    string `json:"shareType,omitempty"`
	ShareProto   string `json:"shareProto,omitempty"`
	AllowDetails bool   `json:"allowDetails"`
}

func (fsr FakeShareRequest) createShare() *pb.Response {
	return &pb.Response{
		Status:  "Success",
		Message: sampleShareData,
	}
}

func (fsr FakeShareRequest) getShare() *pb.Response {
	return &pb.Response{
		Status:  "Success",
		Message: sampleShareDetailData,
	}
}

func (fsr FakeShareRequest) listShares() *pb.Response {
	return &pb.Response{
		Status:  "Success",
		Message: sampleSharesData,
	}
}

func (fsr FakeShareRequest) deleteShare() *pb.Response {
	return &pb.Response{
		Status:  "Success",
		Message: "Delete share success!",
	}
}

func TestCreateShare(t *testing.T) {
	var fsr FakeShareRequest

	err := json.Unmarshal([]byte(sampleShareCreateRequest), &fsr)
	if err != nil {
		t.Fatal(err)
	}

	share, err := CreateShare(fsr)
	if err != nil {
		t.Fatal(err)
	}

	expectedShare := api.ShareResponse{
		Id:   "d94a8548-2079-4be0-b21c-0a887acd31ca",
		Name: "My_share",
		Links: []map[string]string{{"href": "http://172.18.198.54:8786/v2/16e1ab15c35a457e9c2b2aa189f544e1/shares/d94a8548-2079-4be0-b21c-0a887acd31ca", "rel": "self"},
			{"href": "http://172.18.198.54:8786/16e1ab15c35a457e9c2b2aa189f544e1/shares/d94a8548-2079-4be0-b21c-0a887acd31ca", "rel": "bookmark"}}}
	if !reflect.DeepEqual(expectedShare, share) {
		t.Fatalf("Expected\n%#v\ngot\n%#v", expectedShare, share)
	}
	if !reflect.DeepEqual(fsr.Name, share.Name) {
		t.Fatalf("Expected\n%#v\ngot\n%#v", fsr.Name, share.Name)
	}
}

func TestGetShare(t *testing.T) {
	var fsr FakeShareRequest

	err := json.Unmarshal([]byte(sampleShareGetRequest), &fsr)
	if err != nil {
		t.Fatal(err)
	}

	share, err := GetShare(fsr)
	if err != nil {
		t.Fatal(err)
	}

	expectedShare := api.ShareDetailResponse{
		Links: []map[string]string{{"href": "http://172.18.198.54:8786/v2/16e1ab15c35a457e9c2b2aa189f544e1/shares/d94a8548-2079-4be0-b21c-0a887acd31ca", "rel": "self"},
			{"href": "http://172.18.198.54:8786/16e1ab15c35a457e9c2b2aa189f544e1/shares/d94a8548-2079-4be0-b21c-0a887acd31ca", "rel": "bookmark"}},
		AvailabilityZone:         "nova",
		ShareNetworkId:           "713df749-aac0-4a54-af52-10f6c991e80c",
		ExportLocations:          []string{},
		ShareServerId:            "e268f4aa-d571-43dd-9ab3-f49ad06ffaef",
		SnapshotId:               "",
		Id:                       "d94a8548-2079-4be0-b21c-0a887acd31ca",
		Size:                     1,
		ShareType:                "25747776-08e5-494f-ab40-a64b9d20d8f7",
		ShareTypeName:            "default",
		ExportLocation:           "",
		ConsistencyGroupId:       "9397c191-8427-4661-a2e8-b23820dc01d4",
		ProjectId:                "16e1ab15c35a457e9c2b2aa189f544e1",
		Metadata:                 map[string]string{"project": "my_app", "aim": "doc"},
		Status:                   "available",
		Description:              "My custom share London",
		Host:                     "manila2@generic1#GENERIC1",
		AccessRulesStatus:        "active",
		HasReplicas:              false,
		ReplicationType:          "",
		TaskState:                "",
		IsPublic:                 true,
		SnapshotSupport:          true,
		Name:                     "My_share",
		ShareProto:               "NFS",
		VolumeType:               "default",
		SourceCgsnapshotMemberId: ""}
	if !reflect.DeepEqual(expectedShare, share) {
		t.Fatalf("Expected\n%#v\ngot\n%#v", expectedShare, share)
	}
	if !reflect.DeepEqual(fsr.Id, share.Id) {
		t.Fatalf("Expected\n%#v\ngot\n%#v", fsr.Id, share.Id)
	}
}

func TestListShares(t *testing.T) {
	var fsr FakeShareRequest

	err := json.Unmarshal([]byte(sampleShareListRequest), &fsr)
	if err != nil {
		t.Fatal(err)
	}

	shares, err := ListShares(fsr)
	if err != nil {
		t.Fatal(err)
	}

	expectedShare := api.ShareResponse{
		Id: "d94a8548-2079-4be0-b21c-0a887acd31ca",
		Links: []map[string]string{{"href": "http://172.18.198.54:8786/v2/16e1ab15c35a457e9c2b2aa189f544e1/shares/d94a8548-2079-4be0-b21c-0a887acd31ca", "rel": "self"},
			{"href": "http://172.18.198.54:8786/16e1ab15c35a457e9c2b2aa189f544e1/shares/d94a8548-2079-4be0-b21c-0a887acd31ca", "rel": "bookmark"}},
		Name: "My_share"}
	if !reflect.DeepEqual(expectedShare, shares[0]) {
		t.Fatalf("Expected\n%#v\ngot\n%#v", expectedShare, shares[0])
	}
}

func TestDeleteShare(t *testing.T) {
	var fsr FakeShareRequest

	err := json.Unmarshal([]byte(sampleShareDeleteRequest), &fsr)
	if err != nil {
		t.Fatal(err)
	}

	result := DeleteShare(fsr)
	if result.Status != "Success" {
		t.Fatal(result.Error)
	}
}

var sampleShareCreateRequest = `{
	"resourceType":"manila",
	"name":"My_share",
	"shareType":"25747776-08e5-494f-ab40-a64b9d20d8f7",
	"shareProto":"NFS",
	"size":2
}`

var sampleShareGetRequest = `{
	"resourceType":"manila",
	"id":"d94a8548-2079-4be0-b21c-0a887acd31ca"
}`

var sampleShareListRequest = `{
	"resourceType":"manila",
	"allowDetails":false
}`

var sampleShareDeleteRequest = `{
	"resourceType":"manila",
	"id":"d94a8548-2079-4be0-b21c-0a887acd31ca"
}`

var sampleShareData = `{
    "id": "d94a8548-2079-4be0-b21c-0a887acd31ca",
    "links": [
		{
			"href": "http://172.18.198.54:8786/v2/16e1ab15c35a457e9c2b2aa189f544e1/shares/d94a8548-2079-4be0-b21c-0a887acd31ca",
			"rel": "self"
		},
		{
            "href": "http://172.18.198.54:8786/16e1ab15c35a457e9c2b2aa189f544e1/shares/d94a8548-2079-4be0-b21c-0a887acd31ca",
			"rel": "bookmark"
        }
    ],
    "name": "My_share"
}`

var sampleShareDetailData = `{
    "links": [
        {
            "href": "http://172.18.198.54:8786/v2/16e1ab15c35a457e9c2b2aa189f544e1/shares/d94a8548-2079-4be0-b21c-0a887acd31ca",
            "rel": "self"
        },
        {
            "href": "http://172.18.198.54:8786/16e1ab15c35a457e9c2b2aa189f544e1/shares/d94a8548-2079-4be0-b21c-0a887acd31ca",
            "rel": "bookmark"
        }
    ],
    "availability_zone": "nova",
    "share_network_id": "713df749-aac0-4a54-af52-10f6c991e80c",
    "export_locations": [],
    "share_server_id": "e268f4aa-d571-43dd-9ab3-f49ad06ffaef",
    "snapshot_id": null,
    "id": "d94a8548-2079-4be0-b21c-0a887acd31ca",
    "size": 1,
    "share_type": "25747776-08e5-494f-ab40-a64b9d20d8f7",
    "share_type_name": "default",
    "export_location": null,
    "consistency_group_id": "9397c191-8427-4661-a2e8-b23820dc01d4",
    "project_id": "16e1ab15c35a457e9c2b2aa189f544e1",
    "metadata": {
        "project": "my_app",
        "aim": "doc"
    },
    "status": "available",
    "description": "My custom share London",
    "host": "manila2@generic1#GENERIC1",
    "access_rules_status": "active",
    "has_replicas": false,
    "replication_type": null,
    "task_state": null,
    "is_public": true,
    "snapshot_support": true,
    "name": "My_share",
    "created_at": "2015-09-18T10:25:24.000000",
    "share_proto": "NFS",
    "volume_type": "default",
    "source_cgsnapshot_member_id": null
}`

var sampleSharesData = `[
    {
        "id": "d94a8548-2079-4be0-b21c-0a887acd31ca",
        "links": [
            {
                "href": "http://172.18.198.54:8786/v2/16e1ab15c35a457e9c2b2aa189f544e1/shares/d94a8548-2079-4be0-b21c-0a887acd31ca",
                "rel": "self"
            },
            {
                "href": "http://172.18.198.54:8786/16e1ab15c35a457e9c2b2aa189f544e1/shares/d94a8548-2079-4be0-b21c-0a887acd31ca",
                "rel": "bookmark"
            }
        ],
        "name": "My_share"
    },
    {
        "id": "406ea93b-32e9-4907-a117-148b3945749f",
        "links": [
            {
                "href": "http://172.18.198.54:8786/v2/16e1ab15c35a457e9c2b2aa189f544e1/shares/406ea93b-32e9-4907-a117-148b3945749f",
                "rel": "self"
            },
            {
                "href": "http://172.18.198.54:8786/16e1ab15c35a457e9c2b2aa189f544e1/shares/406ea93b-32e9-4907-a117-148b3945749f",
                "rel": "bookmark"
            }
        ],
        "name": "Share1"
    }
]`
