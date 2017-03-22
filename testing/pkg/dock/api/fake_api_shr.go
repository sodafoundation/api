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
This module implements the entry into operations of storageDock module.

*/

package api

import (
	pb "github.com/opensds/opensds/testing/pkg/grpc/fake_opensds"
)

func CreateShare(sr *pb.ShareRequest) (*pb.Response, error) {
	return &pb.Response{
		Status:  "Success",
		Message: sampleShareData,
	}, nil
}

func GetShare(sr *pb.ShareRequest) (*pb.Response, error) {
	return &pb.Response{
		Status:  "Success",
		Message: sampleShareDetailData,
	}, nil
}

func ListShares(sr *pb.ShareRequest) (*pb.Response, error) {
	return &pb.Response{
		Status:  "Success",
		Message: sampleSharesData,
	}, nil
}

func DeleteShare(sr *pb.ShareRequest) (*pb.Response, error) {
	return &pb.Response{
		Status:  "Success",
		Message: "Delete share success!",
	}, nil
}

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
