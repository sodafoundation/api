// Copyright (c) 2014 Hewlett-Packard Development Company, L.P.

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

// volume.go
package v2

import (
	"errors"
	"net/http"
	"strings"
	"testing"

	"git.openstack.org/openstack/golang-client/openstack"
	"git.openstack.org/openstack/golang-client/testUtil"
	"git.openstack.org/openstack/golang-client/util"
)

var tokn = "ae5aebe5-6a5d-4a40-840a-9736a067aff4"

func TestListVolumes(t *testing.T) {
	anon := func(volumeService *Service) {
		volumes, err := volumeService.Volumes()
		if err != nil {
			t.Error(err)
		}

		if len(volumes) != 2 {
			t.Error(errors.New("Incorrect number of volumes found"))
		}
		expectedVolume := Response{
			Name: "volume_test1",
			ID:   "f5fc9874-fc89-4814-a358-23ba83a6115f",
			Links: []map[string]string{{"href": "http://172.16.197.131:8776/v2/1d8837c5fcef4892951397df97661f97/volumes/f5fc9874-fc89-4814-a358-23ba83a6115f", "rel": "self"},
				{"href": "http://172.16.197.131:8776/1d8837c5fcef4892951397df97661f97/volumes/f5fc9874-fc89-4814-a358-23ba83a6115f", "rel": "bookmark"}}}
		// Verify first one matches expected values
		testUtil.Equals(t, expectedVolume, volumes[0])
	}

	testVolumeServiceAction(t, "volumes", sampleVolumesData, anon)
}

func TestListVolumeDetails(t *testing.T) {
	anon := func(volumeService *Service) {
		volumes, err := volumeService.VolumesDetail()
		if err != nil {
			t.Error(err)
		}

		if len(volumes) != 2 {
			t.Error(errors.New("Incorrect number of volumes found"))
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

	testVolumeServiceAction(t, "volumes/detail", sampleVolumeDetailsData, anon)
}

func TestLimitFilterUrlProduced(t *testing.T) {
	testVolumeQueryParameter(t, "volumes?limit=2",
		QueryParameters{Limit: 2})
}

func TestAll_tenantFilterUrlProduced(t *testing.T) {
	testVolumeQueryParameter(t, "volumes?all_tenant=1",
		QueryParameters{All_tenant: 1})
}

func TestMarkerUrlProduced(t *testing.T) {
	testVolumeQueryParameter(t, "volumes?marker=1776335d-72f1-48c9-b0e7-74c62cb8fede",
		QueryParameters{Marker: "1776335d-72f1-48c9-b0e7-74c62cb8fede"})
}

func TestSortKeySortUrlProduced(t *testing.T) {
	testVolumeQueryParameter(t, "volumes?sort_key=id",
		QueryParameters{SortKey: "id"})
}

func TestSortDirSortUrlProduced(t *testing.T) {
	testVolumeQueryParameter(t, "volumes?sort_dir=asc",
		QueryParameters{SortDirection: Asc})
}

func testVolumeQueryParameter(t *testing.T, uriEndsWith string, queryParameters QueryParameters) {
	anon := func(volumeService *Service) {
		_, _ = volumeService.QueryVolumes(&queryParameters)
	}

	testVolumeServiceAction(t, uriEndsWith, sampleVolumesData, anon)
}

func testVolumeServiceAction(t *testing.T, uriEndsWith string, testData string, volumeServiceAction func(*Service)) {
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
	volumeService := Service{
		Session: *sess,
		URL:     apiServer.URL,
	}
	volumeServiceAction(&volumeService)
}

var sampleVolumesData = `{
   "volumes":[
	  {
		 "name":"volume_test1",
		 "id":"f5fc9874-fc89-4814-a358-23ba83a6115f",
		 "links":[{"href": "http://172.16.197.131:8776/v2/1d8837c5fcef4892951397df97661f97/volumes/f5fc9874-fc89-4814-a358-23ba83a6115f", "rel": "self"},
		 {"href": "http://172.16.197.131:8776/1d8837c5fcef4892951397df97661f97/volumes/f5fc9874-fc89-4814-a358-23ba83a6115f", "rel": "bookmark"}]
	  },
	  {
		 "name":"volume_test1",
		 "id":"60055a0a-2451-4d78-af9c-f2302150602f",
		 "links":[{"href": "http://172.16.197.131:8776/v2/1d8837c5fcef4892951397df97661f97/volumes/60055a0a-2451-4d78-af9c-f2302150602f", "rel": "self"},
		 {"href": "http://172.16.197.131:8776/1d8837c5fcef4892951397df97661f97/volumes/60055a0a-2451-4d78-af9c-f2302150602f", "rel": "bookmark"}]
	  }
   ]
}`

var sampleVolumeDetailsData = `{
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
