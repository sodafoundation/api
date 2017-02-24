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

package volumes

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"openstack/golang-client/util"

	"github.com/opensds/opensds/pkg/api"
)

type fakeVolumeRequest struct {
	ResourceType string `json:"resourcetType,omitempty"`
	Id           string `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	Size         int    `json:"size"`
	AllowDetails bool   `json:"allowDetails"`

	ActionType string `json:"actionType,omitempty"`
	Host       string `json:"host,omitempty"`
	Device     string `json:"device,omitempty"`
	Attachment string `json:"attachment,omitempty"`
	MountDir   string `json:"mountDir,omitempty"`
	FsType     string `json:"fsType,omitempty"`
}

func (fvr fakeVolumeRequest) createVolume() (string, error) {
	return sampleVolumeData, nil
}

func (fvr fakeVolumeRequest) getVolume() (string, error) {
	return sampleVolumeDetailData, nil
}

func (fvr fakeVolumeRequest) getAllVolumes() (string, error) {
	return sampleVolumesData, nil
}

func (fvr fakeVolumeRequest) updateVolume() (string, error) {
	return sampleModifiedVolumeData, nil
}

func (fvr fakeVolumeRequest) deleteVolume() (string, error) {
	return "Delete volume success!", nil
}

func (fvr fakeVolumeRequest) attachVolume() (string, error) {
	return "Attach volume success!", nil
}

func (fvr fakeVolumeRequest) detachVolume() (string, error) {
	return "Detach volume success!", nil
}

func (fvr fakeVolumeRequest) mountVolume() (string, error) {
	return "Mount volume success!", nil
}

func (fvr fakeVolumeRequest) unmountVolume() (string, error) {
	return "Unmount volume success!", nil
}

func TestCreate(t *testing.T) {
	var fvr fakeVolumeRequest

	err := json.Unmarshal([]byte(sampleVolumeCreateRequest), &fvr)
	if err != nil {
		t.Fatal(err)
	}

	volume, err := Create(fvr)
	if err != nil {
		t.Fatal(err)
	}

	expectedVolume := api.VolumeResponse{
		Name:        "myvol1",
		ID:          "f5fc9874-fc89-4814-a358-23ba83a6115f",
		Status:      "available",
		Size:        2,
		VolumeType:  "lvmdriver-1",
		Attachments: []map[string]string{}}
	if !reflect.DeepEqual(expectedVolume, volume) {
		t.Fatalf("Expected\n%#v\ngot\n%#v", expectedVolume, volume)
	}
	if !reflect.DeepEqual(fvr.Name, volume.Name) {
		t.Fatalf("Expected\n%#v\ngot\n%#v", fvr.Name, volume.Name)
	}
	if !reflect.DeepEqual(fvr.Size, volume.Size) {
		t.Fatalf("Expected\n%#v\ngot\n%#v", fvr.Size, volume.Size)
	}
}

func TestGet(t *testing.T) {
	var fvr fakeVolumeRequest

	err := json.Unmarshal([]byte(sampleVolumeGetRequest), &fvr)
	if err != nil {
		t.Fatal(err)
	}

	volume, err := Show(fvr)
	if err != nil {
		t.Fatal(err)
	}

	createdAt, _ := util.NewDateTime(`"2014-09-29T14:44:31"`)
	expectedVolume := api.VolumeDetailResponse{
		Id:          "30becf77-63fe-4f5e-9507-a0578ffe0949",
		Attachments: []map[string]string{{"attachment_id": "ddb2ac07-ed62-49eb-93da-73b258dd9bec", "host_name": "host_test", "volume_id": "30becf77-63fe-4f5e-9507-a0578ffe0949", "device": "/dev/vdb", "id": "30becf77-63fe-4f5e-9507-a0578ffe0949", "server_id": "0f081aae-1b0c-4b89-930c-5f2562460c72"}},
		Links: []map[string]string{{"href": "http://172.16.197.131:8776/v2/1d8837c5fcef4892951397df97661f97/volumes/30becf77-63fe-4f5e-9507-a0578ffe0949", "rel": "self"},
			{"href": "http://172.16.197.131:8776/1d8837c5fcef4892951397df97661f97/volumes/30becf77-63fe-4f5e-9507-a0578ffe0949", "rel": "bookmark"}},
		Metadata:        map[string]string{"readonly": "false", "attached_mode": "rw"},
		Protected:       false,
		Status:          "available",
		MigrationStatus: "",
		UserID:          "a971aa69-c61a-4a49-b392-b0e41609bc5d",
		Encrypted:       false,
		Multiattach:     false,
		CreatedAt:       createdAt,
		Description:     "test volume",
		VolumeType:      "test_type",
		Name:            "test_volume",
		SourceVolid:     "4b58bbb8-3b00-4f87-8243-8c622707bbab",
		SnapshotId:      "cc488e4a-9649-4e5f-ad12-20ab37c683b5",
		Size:            2,

		AvailabilityZone:   "default_cluster",
		ReplicationStatus:  "",
		ConsistencygroupId: ""}
	if !reflect.DeepEqual(expectedVolume, volume) {
		t.Fatalf("Expected\n%#v\ngot\n%#v", expectedVolume, volume)
	}
	if !reflect.DeepEqual(fvr.Id, volume.Id) {
		t.Fatalf("Expected\n%#v\ngot\n%#v", fvr.Id, volume.Id)
	}
}

func TestList(t *testing.T) {
	var fvr fakeVolumeRequest

	err := json.Unmarshal([]byte(sampleVolumeListRequest), &fvr)
	if err != nil {
		t.Fatal(err)
	}

	volumes, err := List(fvr)
	if err != nil {
		t.Fatal(err)
	}

	expectedVolume := api.VolumeResponse{
		Name:        "myvol1",
		ID:          "f5fc9874-fc89-4814-a358-23ba83a6115f",
		Status:      "in-use",
		Size:        1,
		VolumeType:  "lvmdriver-1",
		Attachments: []map[string]string{{"attached_at": "2017-02-11T14:08:17.000000", "attachment_id": "c7f84865-640c-44ea-94ab-379a27f0ff65", "device": "/dev/vdc", "host_name": "localhost", "id": "034af8c9-ef44-4855-8e70-d51dceed7fc4", "server_id": "", "volume_id": "034af8c9-ef44-4855-8e70-d51dceed7fc4"}}}
	if !reflect.DeepEqual(expectedVolume, volumes[0]) {
		t.Fatalf("Expected\n%#v\ngot\n%#v", expectedVolume, volumes[0])
	}
}

func TestUpdate(t *testing.T) {
	var fvr fakeVolumeRequest

	err := json.Unmarshal([]byte(sampleVolumeUpdateRequest), &fvr)
	if err != nil {
		t.Fatal(err)
	}

	volume, err := Update(fvr)
	if err != nil {
		t.Fatal(err)
	}

	expectedVolume := api.VolumeResponse{
		Name:        "newvol",
		ID:          "f5fc9874-fc89-4814-a358-23ba83a6115f",
		Status:      "available",
		Size:        2,
		VolumeType:  "lvmdriver-1",
		Attachments: []map[string]string{}}
	if !reflect.DeepEqual(expectedVolume, volume) {
		t.Fatalf("Expected\n%#v\ngot\n%#v", expectedVolume, volume)
	}
	if !reflect.DeepEqual(fvr.Name, volume.Name) {
		t.Fatalf("Expected\n%#v\ngot\n%#v", fvr.Name, volume.Name)
	}
}

func TestDelete(t *testing.T) {
	var fvr fakeVolumeRequest

	err := json.Unmarshal([]byte(sampleVolumeDeleteRequest), &fvr)
	if err != nil {
		t.Fatal(err)
	}

	result, err := Delete(fvr)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(result, "success") {
		t.Fatal("Delete volume failed!")
	}
}

func TestAttach(t *testing.T) {
	var fvr fakeVolumeRequest

	err := json.Unmarshal([]byte(sampleVolumeAttachRequest), &fvr)
	if err != nil {
		t.Fatal(err)
	}

	result, err := Attach(fvr)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(result, "success") {
		t.Fatal("Attach volume failed!")
	}
}

func TestDetach(t *testing.T) {
	var fvr fakeVolumeRequest

	err := json.Unmarshal([]byte(sampleVolumeDetachRequest), &fvr)
	if err != nil {
		t.Fatal(err)
	}

	result, err := Detach(fvr)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(result, "success") {
		t.Fatal("Detach volume failed!")
	}
}

func TestMount(t *testing.T) {
	var fvr fakeVolumeRequest

	err := json.Unmarshal([]byte(sampleVolumeMountRequest), &fvr)
	if err != nil {
		t.Fatal(err)
	}

	result, err := Mount(fvr)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(result, "success") {
		t.Fatal("Mount volume failed!")
	}
}

func TestUnmount(t *testing.T) {
	var fvr fakeVolumeRequest

	err := json.Unmarshal([]byte(sampleVolumeUnmountRequest), &fvr)
	if err != nil {
		t.Fatal(err)
	}

	result, err := Unmount(fvr)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(result, "success") {
		t.Fatal("Delete volume failed!")
	}
}

var sampleVolumeCreateRequest = `{
	"resourceType":"cinder",
	"name":"myvol1",
	"size":2
}`

var sampleVolumeGetRequest = `{
	"resourceType":"cinder",
	"id":"30becf77-63fe-4f5e-9507-a0578ffe0949"
}`

var sampleVolumeListRequest = `{
	"resourceType":"cinder",
	"allowDetails":false
}`

var sampleVolumeUpdateRequest = `{
	"resourceType":"cinder",
	"id":"f5fc9874-fc89-4814-a358-23ba83a6115f",
	"name":"newvol"
}`

var sampleVolumeDeleteRequest = `{
	"resourceType":"cinder",
	"id":"f5fc9874-fc89-4814-a358-23ba83a6115f"
}`

var sampleVolumeAttachRequest = `{
	"resourceType":"cinder",
	"id":"f5fc9874-fc89-4814-a358-23ba83a6115f",
	"host":"localhost",
	"device":"/dev/vdc"
}`

var sampleVolumeDetachRequest = `{
	"resourceType":"cinder",
	"id":"f5fc9874-fc89-4814-a358-23ba83a6115f",
	"attachment":"ddb2ac07-ed62-49eb-93da-73b258dd9bec"
}`

var sampleVolumeMountRequest = `{
	"mountDir":"/mnt",
	"device":"/dev/vdc",
	"id":"f5fc9874-fc89-4814-a358-23ba83a6115f",
	"fsType":"ext4"
}`

var sampleVolumeUnmountRequest = `{
	"mountDir":"/mnt"
}`

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
	"status":"available",
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

var sampleModifiedVolumeData = `{
	"name":"newvol",
	"id":"f5fc9874-fc89-4814-a358-23ba83a6115f",
	"status":"available",
	"size":2,
	"volume_type":"lvmdriver-1",
	"attachments":[]
}`
