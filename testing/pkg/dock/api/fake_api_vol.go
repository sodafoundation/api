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

func CreateVolume(vr *pb.VolumeRequest) (*pb.Response, error) {
	return &pb.Response{
		Status:  "Success",
		Message: sampleVolumeData,
	}, nil
}

func GetVolume(vr *pb.VolumeRequest) (*pb.Response, error) {
	return &pb.Response{
		Status:  "Success",
		Message: sampleVolumeDetailData,
	}, nil
}

func ListVolumes(vr *pb.VolumeRequest) (*pb.Response, error) {
	return &pb.Response{
		Status:  "Success",
		Message: sampleVolumesData,
	}, nil
}

func DeleteVolume(vr *pb.VolumeRequest) (*pb.Response, error) {
	return &pb.Response{
		Status:  "Success",
		Message: "Delete volume success!",
	}, nil
}

func AttachVolume(vr *pb.VolumeRequest) (*pb.Response, error) {
	return &pb.Response{
		Status:  "Success",
		Message: "Attach volume success!",
	}, nil
}

func DetachVolume(vr *pb.VolumeRequest) (*pb.Response, error) {
	return &pb.Response{
		Status:  "Success",
		Message: "Detach volume success!",
	}, nil
}

func MountVolume(vr *pb.VolumeRequest) (*pb.Response, error) {
	return &pb.Response{
		Status:  "Success",
		Message: "Mount volume success!",
	}, nil
}

func UnmountVolume(vr *pb.VolumeRequest) (*pb.Response, error) {
	return &pb.Response{
		Status:  "Success",
		Message: "Unmount volume success!",
	}, nil
}

var sampleVolumeData = `{
	"name":"myvol1",
	"id":"f5fc9874-fc89-4814-a358-23ba83a6115f",
	"status":"available",
	"size":2,
	"volume_type":"lvmdriver-1",
	"attachments":[]
}`

var sampleVolumeDetailData = `{
	"id":"30becf77-63fe-4f5e-9507-a0578ffe0949",
	"attachments":[
		{
			"attachment_id": "ddb2ac07-ed62-49eb-93da-73b258dd9bec",
			"host_name": "host_test",
			"volume_id": "30becf77-63fe-4f5e-9507-a0578ffe0949",
			"device": "/dev/vdb",
			"id": "30becf77-63fe-4f5e-9507-a0578ffe0949",
			"server_id": "0f081aae-1b0c-4b89-930c-5f2562460c72"
		}
	],
	"links":[
		{
			"href": "http://172.16.197.131:8776/v2/1d8837c5fcef4892951397df97661f97/volumes/30becf77-63fe-4f5e-9507-a0578ffe0949",
			"rel": "self"
		},
		{
			"href": "http://172.16.197.131:8776/1d8837c5fcef4892951397df97661f97/volumes/30becf77-63fe-4f5e-9507-a0578ffe0949",
			"rel": "bookmark"
		}
	],
	"metadata":{
		"readonly": "false",
		"attached_mode": "rw"
	},
	"protected":false,
	"status":"in-use",
	"migrationStatus":null,
	"user_id":"a971aa69-c61a-4a49-b392-b0e41609bc5d",
	"encrypted":false,
	"multiattach":false,
	"created_at":"2014-09-29T14:44:31",
	"description":"test volume",
	"volume_type":"test_type",
	"name":"test_volume",
	"source_volid":"4b58bbb8-3b00-4f87-8243-8c622707bbab",
	"snapshot_id":"cc488e4a-9649-4e5f-ad12-20ab37c683b5",
	"size":2,

	"availability_zone":"default_cluster",
	"replication_status":null,
	"consistencygroup_id":null
}`

var sampleVolumesData = `[
	{
		"name":"myvol1",
		"id":"f5fc9874-fc89-4814-a358-23ba83a6115f",
		"status":"in-use",
		"size":1,
		"volume_type":"lvmdriver-1",
		"attachments":[
			{
				"attached_at":"2017-02-11T14:08:17.000000",
				"attachment_id":"c7f84865-640c-44ea-94ab-379a27f0ff65",
				"device":"/dev/vdc",
				"host_name":"localhost",
				"id":"034af8c9-ef44-4855-8e70-d51dceed7fc4",
				"server_id":"",
				"volume_id":"034af8c9-ef44-4855-8e70-d51dceed7fc4"
			}
		]
	},
	{
		"name":"myvol2",
		"id":"60055a0a-2451-4d78-af9c-f2302150602f",
		"status":"available",
		"size":2,
		"volume_type":"lvmdriver-1",
		"attachments":[]
	}
]`
