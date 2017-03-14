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
Package volume implements a client library for accessing OpenStack Volume service

The CRUD operation of volumes can be retrieved using the api.

*/

package volume

import (
	"errors"
	"net/http"
	"strings"
	"testing"

	"git.openstack.org/openstack/golang-client.git/openstack"
	"git.openstack.org/openstack/golang-client.git/testUtil"
	"git.openstack.org/openstack/golang-client.git/util"
)

var tokn = "ae5aebe5-6a5d-4a40-840a-9736a067aff4"

func TestCreateVolume(t *testing.T) {
	anon := func(volumeService Service) {
		requestBody := RequestBody{}
		requestBody.Name = "myvol1"
		requestBody.Size = 2
		body := CreateBody{requestBody}
		result, err := volumeService.Create(&body)
		if err != nil {
			t.Error(err)
		}

		expectedVolume := Response{
			Name:        "myvol1",
			ID:          "f5fc9874-fc89-4814-a358-23ba83a6115f",
			Status:      "available",
			Size:        2,
			Volume_type: "lvmdriver-1",
			Attachments: []map[string]string{}}
		testUtil.Equals(t, expectedVolume, result)
	}

	testCreateVolumeServiceAction(t, tokn, sampleVolumeData, "volumes", sampleRequestBody, anon)
}

func testCreateVolumeServiceAction(t *testing.T, tokn string, testData string, uriEndsWith string, sampleRequestBody string, volumeServiceAction func(Service)) {
	apiServer := testUtil.CreatePostJSONTestRequestServer(t, tokn, testData, uriEndsWith, sampleRequestBody)
	defer apiServer.Close()

	auth := openstack.AuthToken{
		Access: openstack.AccessType{
			Token: openstack.Token{
				ID: tokn,
			},
		},
	}
	sess, _ := openstack.NewSession(http.DefaultClient, auth, nil)
	volumeService, _ := NewService(*sess, *http.DefaultClient, apiServer.URL)
	volumeServiceAction(volumeService)
}

func TestGetVolume(t *testing.T) {
	anon := func(volumeService Service) {
		volID := "30becf77-63fe-4f5e-9507-a0578ffe0949"
		result, err := volumeService.Show(volID)
		if err != nil {
			t.Error(err)
		}

		createdAt, _ := util.NewDateTime(`"2014-09-29T14:44:31"`)
		expectedVolume := DetailResponse{
			ID:          "30becf77-63fe-4f5e-9507-a0578ffe0949",
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
			Volume_type:     "test_type",
			Name:            "test_volume",
			Source_volid:    "4b58bbb8-3b00-4f87-8243-8c622707bbab",
			Snapshot_id:     "cc488e4a-9649-4e5f-ad12-20ab37c683b5",
			Size:            2,

			Aavailability_zone:  "default_cluster",
			Rreplication_status: "",
			Consistencygroup_id: ""}
		testUtil.Equals(t, expectedVolume, result)
	}

	testGetVolumeServiceAction(t, "30becf77-63fe-4f5e-9507-a0578ffe0949", sampleVolumeDetailData, anon)
}

func testGetVolumeServiceAction(t *testing.T, uriEndsWith string, testData string, volumeServiceAction func(Service)) {
	anon := func(req *http.Request) {
		reqURL := req.URL.String()
		if !strings.HasSuffix(reqURL, uriEndsWith) {
			t.Error(errors.New("Incorrect url created, expected:" + uriEndsWith + " at the end, actual url:" + reqURL))
		}
	}
	apiServer := testUtil.CreateGetJSONTestRequestServer(t, tokn, testData, anon)
	defer apiServer.Close()

	auth := openstack.AuthToken{
		Access: openstack.AccessType{
			Token: openstack.Token{
				ID: tokn,
			},
		},
	}
	sess, _ := openstack.NewSession(http.DefaultClient, auth, nil)
	volumeService, _ := NewService(*sess, *http.DefaultClient, apiServer.URL)
	volumeServiceAction(volumeService)
}

func TestGetAllVolumes(t *testing.T) {
	anon := func(volumeService Service) {
		volumes, err := volumeService.List()
		if err != nil {
			t.Error(err)
		}

		expectedVolume := Response{
			Name:        "myvol1",
			ID:          "f5fc9874-fc89-4814-a358-23ba83a6115f",
			Status:      "in-use",
			Size:        1,
			Volume_type: "lvmdriver-1",
			Attachments: []map[string]string{{"attached_at": "2017-02-11T14:08:17.000000", "attachment_id": "c7f84865-640c-44ea-94ab-379a27f0ff65", "device": "/dev/vdc", "host_name": "localhost", "id": "034af8c9-ef44-4855-8e70-d51dceed7fc4", "server_id": "", "volume_id": "034af8c9-ef44-4855-8e70-d51dceed7fc4"}}}
		testUtil.Equals(t, expectedVolume, volumes[0])
	}

	testDetailAllVolumesServiceAction(t, "volumes/detail", sampleVolumesData, anon)
}

func testGetAllVolumesServiceAction(t *testing.T, uriEndsWith string, testData string, volumeServiceAction func(Service)) {
	anon := func(req *http.Request) {
		reqURL := req.URL.String()
		if !strings.HasSuffix(reqURL, uriEndsWith) {
			t.Error(errors.New("Incorrect url created, expected:" + uriEndsWith + " at the end, actual url:" + reqURL))
		}
	}
	apiServer := testUtil.CreateGetJSONTestRequestServer(t, tokn, testData, anon)
	defer apiServer.Close()

	auth := openstack.AuthToken{
		Access: openstack.AccessType{
			Token: openstack.Token{
				ID: tokn,
			},
		},
	}
	sess, _ := openstack.NewSession(http.DefaultClient, auth, nil)
	volumeService, _ := NewService(*sess, *http.DefaultClient, apiServer.URL)
	volumeServiceAction(volumeService)
}

func TestTailAllVolumes(t *testing.T) {
	anon := func(volumeService Service) {
		volumes, err := volumeService.Detail()
		if err != nil {
			t.Error(err)
		}

		createdAt, _ := util.NewDateTime(`"2014-09-29T14:44:31"`)
		expectedVolumeDetail := DetailResponse{
			ID:          "30becf77-63fe-4f5e-9507-a0578ffe0949",
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
			Volume_type:     "test_type",
			Name:            "test_volume",
			Source_volid:    "4b58bbb8-3b00-4f87-8243-8c622707bbab",
			Snapshot_id:     "cc488e4a-9649-4e5f-ad12-20ab37c683b5",
			Size:            2,

			Aavailability_zone:  "default_cluster",
			Rreplication_status: "",
			Consistencygroup_id: ""}

		testUtil.Equals(t, expectedVolumeDetail, volumes[0])
	}

	testDetailAllVolumesServiceAction(t, "volumes/detail", sampleVolumesDetailData, anon)
}

func testDetailAllVolumesServiceAction(t *testing.T, uriEndsWith string, testData string, volumeServiceAction func(Service)) {
	anon := func(req *http.Request) {
		reqURL := req.URL.String()
		if !strings.HasSuffix(reqURL, uriEndsWith) {
			t.Error(errors.New("Incorrect url created, expected:" + uriEndsWith + " at the end, actual url:" + reqURL))
		}
	}
	apiServer := testUtil.CreateGetJSONTestRequestServer(t, tokn, testData, anon)
	defer apiServer.Close()

	auth := openstack.AuthToken{
		Access: openstack.AccessType{
			Token: openstack.Token{
				ID: tokn,
			},
		},
	}
	sess, _ := openstack.NewSession(http.DefaultClient, auth, nil)
	volumeService, _ := NewService(*sess, *http.DefaultClient, apiServer.URL)
	volumeServiceAction(volumeService)
}

func TestDeleteVolume(t *testing.T) {
	anon := func(volumeService Service) {
		volID := "30becf77-63fe-4f5e-9507-a0578ffe0949"
		err := volumeService.Delete(volID)
		if err != nil {
			t.Error(err)
		}
	}

	testDeleteVolumeServiceAction(t, "volumes/30becf77-63fe-4f5e-9507-a0578ffe0949", anon)
}

func testDeleteVolumeServiceAction(t *testing.T, uriEndsWith string, volumeServiceAction func(Service)) {
	apiServer := testUtil.CreateDeleteTestRequestServer(t, tokn, uriEndsWith)
	defer apiServer.Close()

	auth := openstack.AuthToken{
		Access: openstack.AccessType{
			Token: openstack.Token{
				ID: tokn,
			},
		},
	}
	sess, _ := openstack.NewSession(http.DefaultClient, auth, nil)
	volumeService, _ := NewService(*sess, *http.DefaultClient, apiServer.URL)
	volumeServiceAction(volumeService)
}

var sampleRequestBody = `{"volume":{"name":"myvol1","size":2,"host_name":"","device":"","attachment_id":""}}`

var sampleVolumeData = `{
	"volume":{
		"name":"myvol1",
		"id":"f5fc9874-fc89-4814-a358-23ba83a6115f",
		"status":"available",
		"size":2,
		"volume_type":"lvmdriver-1",
		"attachments":[]
	}
}`

var sampleVolumesData = `{
   "volumes":[
		{
			"name":"myvol1",
			"id":"f5fc9874-fc89-4814-a358-23ba83a6115f",
			"status":"in-use",
			"size":1,
			"volume_type":"lvmdriver-1",
			"attachments":[
				{"attached_at":"2017-02-11T14:08:17.000000","attachment_id":"c7f84865-640c-44ea-94ab-379a27f0ff65","device":"/dev/vdc","host_name":"localhost","id":"034af8c9-ef44-4855-8e70-d51dceed7fc4","server_id":"","volume_id":"034af8c9-ef44-4855-8e70-d51dceed7fc4"}
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
   	]
}`

var sampleVolumeDetailData = `{
   "volume": {
		"id":"30becf77-63fe-4f5e-9507-a0578ffe0949",
		"attachments":[{"attachment_id": "ddb2ac07-ed62-49eb-93da-73b258dd9bec", "host_name": "host_test", "volume_id": "30becf77-63fe-4f5e-9507-a0578ffe0949", "device": "/dev/vdb", "id": "30becf77-63fe-4f5e-9507-a0578ffe0949", "server_id": "0f081aae-1b0c-4b89-930c-5f2562460c72"}],
		"links":[{"href": "http://172.16.197.131:8776/v2/1d8837c5fcef4892951397df97661f97/volumes/30becf77-63fe-4f5e-9507-a0578ffe0949", "rel": "self"},
				{"href": "http://172.16.197.131:8776/1d8837c5fcef4892951397df97661f97/volumes/30becf77-63fe-4f5e-9507-a0578ffe0949", "rel": "bookmark"}],
		"metadata":{"readonly": "false", "attached_mode": "rw"},
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
	}
}`

var sampleVolumesDetailData = `{
   "volumes":[
	  {
		"id":"30becf77-63fe-4f5e-9507-a0578ffe0949",
		"attachments":[{"attachment_id": "ddb2ac07-ed62-49eb-93da-73b258dd9bec", "host_name": "host_test", "volume_id": "30becf77-63fe-4f5e-9507-a0578ffe0949", "device": "/dev/vdb", "id": "30becf77-63fe-4f5e-9507-a0578ffe0949", "server_id": "0f081aae-1b0c-4b89-930c-5f2562460c72"}],
		"links":[{"href": "http://172.16.197.131:8776/v2/1d8837c5fcef4892951397df97661f97/volumes/30becf77-63fe-4f5e-9507-a0578ffe0949", "rel": "self"},
				{"href": "http://172.16.197.131:8776/1d8837c5fcef4892951397df97661f97/volumes/30becf77-63fe-4f5e-9507-a0578ffe0949", "rel": "bookmark"}],
		"metadata":{"readonly": "false", "attached_mode": "rw"},
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
	  },
	  {
		"id":"242d3d14-2efd-4c63-9a6b-ef6bc8eed756",
		"attachments":[{"attachment_id": "9d4fb045-f957-489b-9e7d-f6f156002c04", "host_name": "host_test2", "volume_id": "242d3d14-2efd-4c63-9a6b-ef6bc8eed756", "device": "/dev/vdb", "id": "242d3d14-2efd-4c63-9a6b-ef6bc8eed756", "server_id": "9f47bd1c-c596-424d-abbe-e5e1a7a65fdc"}],
		"links":[{"href": "http://172.16.197.131:8776/v2/1d8837c5fcef4892951397df97661f97/volumes/242d3d14-2efd-4c63-9a6b-ef6bc8eed756", "rel": "self"},
				{"href": "http://172.16.197.131:8776/1d8837c5fcef4892951397df97661f97/volumes/242d3d14-2efd-4c63-9a6b-ef6bc8eed756", "rel": "bookmark"}],
		"metadata":{"readonly": "false", "attached_mode": "rw"},
		"protected":false,
		"status":"available",
		"migrationStatus":null,
		"user_id":"a971aa69-c61a-4a49-b392-b0e41609bc5d",
		"encrypted":false,
		"multiattach":false,
		"created_at":"2014-09-29T14:44:35",
		"description":"test volume 2",
		"volume_type":"test_type",
		"name":"test_volume2",
		"source_volid":null,
		"snapshot_id":null,
		"size":2,

		"availability_zone":"default_cluster",
		"replication_status":null,
		"consistencygroup_id":null
	  }
   ]
}`
