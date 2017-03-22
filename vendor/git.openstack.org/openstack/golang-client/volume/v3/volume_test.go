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

package v3

import (
	"errors"
	"net/http"
	"strings"
	"testing"

	"git.openstack.org/openstack/golang-client/openstack"
	"git.openstack.org/openstack/golang-client/testUtil"
)

var tokn = "ae5aebe5-6a5d-4a40-840a-9736a067aff4"

/*
func TestCreateVolume(t *testing.T) {
	anon := func(volumeService *v3.Service) {
		requestBody := v3.RequestBody{100, "myvol1"}
		volume, err := volumeService.Create(&requestBody)
		if err != nil {
			t.Error(err)
		}

		expectedVolume := v3.Response{
			Name: "myvol1",
			ID:   "f5fc9874-fc89-4814-a358-23ba83a6115f",
			Links: []map[string]string{{"href": "http://172.16.197.131:8776/v2/1d8837c5fcef4892951397df97661f97/volumes/f5fc9874-fc89-4814-a358-23ba83a6115f", "rel": "self"},
				{"href": "http://172.16.197.131:8776/1d8837c5fcef4892951397df97661f97/volumes/f5fc9874-fc89-4814-a358-23ba83a6115f", "rel": "bookmark"}}}
		testUtil.Equals(t, expectedVolume, volume)
	}

	//testCreateVolumeServiceAction(t, "volumes", sampleVolumesData, anon)
}
*/

// TODO(dtroyer): skipping due to job failure for now, this must be fixed
// func TestGetVolume(t *testing.T) {
// 	anon := func(volumeService *Service) {
// 		volID := "f5fc9874-fc89-4814-a358-23ba83a6115f"
// 		volume, err := volumeService.Show(volID)
// 		if err != nil {
// 			t.Error(err)
// 		}

// 		expectedVolume := Response{
// 			Name: "myvol1",
// 			ID:   "f5fc9874-fc89-4814-a358-23ba83a6115f",
// 			// Links: []map[string]string{{"href": "http://172.16.197.131:8776/v2/1d8837c5fcef4892951397df97661f97/volumes/f5fc9874-fc89-4814-a358-23ba83a6115f", "rel": "self"},
// 			// 	{"href": "http://172.16.197.131:8776/1d8837c5fcef4892951397df97661f97/volumes/f5fc9874-fc89-4814-a358-23ba83a6115f", "rel": "bookmark"}}
// 		}
// 		testUtil.Equals(t, expectedVolume, volume)
// 	}

// 	testGetVolumeServiceAction(t, "f5fc9874-fc89-4814-a358-23ba83a6115f", sampleVolumeData, anon)
// }

func testGetVolumeServiceAction(t *testing.T, uriEndsWith string, testData string, volumeServiceAction func(*Service)) {
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

func TestGetAllVolumes(t *testing.T) {
	anon := func(volumeService *Service) {
		volumes, err := volumeService.List()
		if err != nil {
			t.Error(err)
		}

		expectedVolume := Response{
			Name: "myvol1",
			ID:   "f5fc9874-fc89-4814-a358-23ba83a6115f",
			// Links: []map[string]string{{"href": "http://172.16.197.131:8776/v2/1d8837c5fcef4892951397df97661f97/volumes/f5fc9874-fc89-4814-a358-23ba83a6115f", "rel": "self"},
			// 	{"href": "http://172.16.197.131:8776/1d8837c5fcef4892951397df97661f97/volumes/f5fc9874-fc89-4814-a358-23ba83a6115f", "rel": "bookmark"}}
		}
		testUtil.Equals(t, expectedVolume, volumes[0])
	}

	testGetAllVolumesServiceAction(t, "volumes", sampleVolumesData, anon)
}

func testGetAllVolumesServiceAction(t *testing.T, uriEndsWith string, testData string, volumeServiceAction func(*Service)) {
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

/*
func TestDeleteVolume(t *testing.T) {
	anon := func(volumeService *v3.Service) {
		volID := "f5fc9874-fc89-4814-a358-23ba83a6115f"
		_, err := volumeService.Delete(volID)
		if err != nil {
			t.Error(err)
		}
		volume, err := volumeService.Show(volID)

		expectedVolume := v3.Response{}
		testUtil.Equals(t, expectedVolume, volume)
	}

	testDeleteVolumeServiceAction(t, "f5fc9874-fc89-4814-a358-23ba83a6115f", anon)
}

func testDeleteVolumeServiceAction(t *testing.T, uriEndsWith string, volumeServiceAction func(*v3.Service)) {
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
	volumeService := v3.Service{
		Session: *sess,
		URL:     apiServer.URL,
	}
	volumeServiceAction(&volumeService)
}
*/

var sampleVolumeData = `{
		 "name":"myvol1",
		 "id":"f5fc9874-fc89-4814-a358-23ba83a6115f",
		 "links":[{"href": "http://172.16.197.131:8776/v2/1d8837c5fcef4892951397df97661f97/volumes/f5fc9874-fc89-4814-a358-23ba83a6115f", "rel": "self"},
		 {"href": "http://172.16.197.131:8776/1d8837c5fcef4892951397df97661f97/volumes/f5fc9874-fc89-4814-a358-23ba83a6115f", "rel": "bookmark"}]
}`

var sampleVolumesData = `{
   "volumes":[
	  {
		 "name":"myvol1",
		 "id":"f5fc9874-fc89-4814-a358-23ba83a6115f",
		 "links":[{"href": "http://172.16.197.131:8776/v2/1d8837c5fcef4892951397df97661f97/volumes/f5fc9874-fc89-4814-a358-23ba83a6115f", "rel": "self"},
		 {"href": "http://172.16.197.131:8776/1d8837c5fcef4892951397df97661f97/volumes/f5fc9874-fc89-4814-a358-23ba83a6115f", "rel": "bookmark"}]
	  },
	  {
		 "name":"myvol2",
		 "id":"60055a0a-2451-4d78-af9c-f2302150602f",
		 "links":[{"href": "http://172.16.197.131:8776/v2/1d8837c5fcef4892951397df97661f97/volumes/60055a0a-2451-4d78-af9c-f2302150602f", "rel": "self"},
		 {"href": "http://172.16.197.131:8776/1d8837c5fcef4892951397df97661f97/volumes/60055a0a-2451-4d78-af9c-f2302150602f", "rel": "bookmark"}]
	  }
   ]
}`
