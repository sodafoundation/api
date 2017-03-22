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
Package share implements a client library for accessing OpenStack Share service

The CRUD operation of shares can be retrieved using the api.

*/

package share

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

func TestCreateShare(t *testing.T) {
	anon := func(shareService Service) {
		requestBody := RequestBody{}
		requestBody.Name = "My_share"
		requestBody.Size = 1
		requestBody.Share_proto = "NFS"
		requestBody.Share_type = "g-nfs"
		body := CreateBody{requestBody}
		result, err := shareService.Create(&body)
		if err != nil {
			t.Error(err)
		}

		expectedShare := Response{
			ID:   "d94a8548-2079-4be0-b21c-0a887acd31ca",
			Name: "My_share",
			Links: []map[string]string{{"href": "http://172.18.198.54:8786/v2/16e1ab15c35a457e9c2b2aa189f544e1/shares/d94a8548-2079-4be0-b21c-0a887acd31ca", "rel": "self"},
				{"href": "http://172.18.198.54:8786/16e1ab15c35a457e9c2b2aa189f544e1/shares/d94a8548-2079-4be0-b21c-0a887acd31ca", "rel": "bookmark"}}}
		testUtil.Equals(t, expectedShare, result)
	}

	testCreateShareServiceAction(t, tokn, sampleShareData, "shares", sampleRequestBody, anon)
}

func testCreateShareServiceAction(t *testing.T, tokn string, testData string, uriEndsWith string, sampleRequestBody string, shareServiceAction func(Service)) {
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
	shareService, _ := NewService(*sess, *http.DefaultClient, apiServer.URL)
	shareServiceAction(shareService)
}

func TestGetShare(t *testing.T) {
	anon := func(shareService Service) {
		volID := "d94a8548-2079-4be0-b21c-0a887acd31ca"
		result, err := shareService.Show(volID)
		if err != nil {
			t.Error(err)
		}

		createdAt, _ := util.NewDateTime(`"2015-09-18T10:25:24.000000"`)
		expectedShareDetail := DetailResponse{
			Links: []map[string]string{{"href": "http://172.18.198.54:8786/v2/16e1ab15c35a457e9c2b2aa189f544e1/shares/d94a8548-2079-4be0-b21c-0a887acd31ca", "rel": "self"},
				{"href": "http://172.18.198.54:8786/16e1ab15c35a457e9c2b2aa189f544e1/shares/d94a8548-2079-4be0-b21c-0a887acd31ca", "rel": "bookmark"}},
			Availability_zone:           "nova",
			Share_network_id:            "713df749-aac0-4a54-af52-10f6c991e80c",
			Export_locations:            []string{},
			Share_server_id:             "e268f4aa-d571-43dd-9ab3-f49ad06ffaef",
			Snapshot_id:                 "",
			ID:                          "d94a8548-2079-4be0-b21c-0a887acd31ca",
			Size:                        1,
			Share_type:                  "25747776-08e5-494f-ab40-a64b9d20d8f7",
			Share_type_name:             "default",
			Export_location:             "",
			Consistency_group_id:        "9397c191-8427-4661-a2e8-b23820dc01d4",
			Project_id:                  "16e1ab15c35a457e9c2b2aa189f544e1",
			Metadata:                    map[string]string{"project": "my_app", "aim": "doc"},
			Status:                      "available",
			Description:                 "My custom share London",
			Host:                        "manila2@generic1#GENERIC1",
			Access_rules_status:         "active",
			Has_replicas:                false,
			Replication_type:            "",
			Task_state:                  "",
			Is_public:                   true,
			Snapshot_support:            true,
			Name:                        "My_share",
			Created_at:                  createdAt,
			Share_proto:                 "NFS",
			Volume_type:                 "default",
			Source_cgsnapshot_member_id: ""}
		testUtil.Equals(t, expectedShareDetail, result)
	}

	testGetShareServiceAction(t, "d94a8548-2079-4be0-b21c-0a887acd31ca", sampleShareDetailData, anon)
}

func testGetShareServiceAction(t *testing.T, uriEndsWith string, testData string, shareServiceAction func(Service)) {
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
	shareService, _ := NewService(*sess, *http.DefaultClient, apiServer.URL)
	shareServiceAction(shareService)
}

func TestGetAllShares(t *testing.T) {
	anon := func(shareService Service) {
		shares, err := shareService.List()
		if err != nil {
			t.Error(err)
		}

		expectedShare := Response{
			ID: "d94a8548-2079-4be0-b21c-0a887acd31ca",
			Links: []map[string]string{{"href": "http://172.18.198.54:8786/v2/16e1ab15c35a457e9c2b2aa189f544e1/shares/d94a8548-2079-4be0-b21c-0a887acd31ca", "rel": "self"},
				{"href": "http://172.18.198.54:8786/16e1ab15c35a457e9c2b2aa189f544e1/shares/d94a8548-2079-4be0-b21c-0a887acd31ca", "rel": "bookmark"}},
			Name: "My_share"}
		testUtil.Equals(t, expectedShare, shares[0])
	}

	testGetAllSharesServiceAction(t, "shares", sampleSharesData, anon)
}

func testGetAllSharesServiceAction(t *testing.T, uriEndsWith string, testData string, shareServiceAction func(Service)) {
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
	shareService, _ := NewService(*sess, *http.DefaultClient, apiServer.URL)
	shareServiceAction(shareService)
}

func TestTailAllShares(t *testing.T) {
	anon := func(shareService Service) {
		shares, err := shareService.Detail()
		if err != nil {
			t.Error(err)
		}

		createdAt, _ := util.NewDateTime(`"2015-09-16T18:19:50.000000"`)
		expectedShareDetail := DetailResponse{
			Links: []map[string]string{{"href": "http://172.18.198.54:8786/v2/16e1ab15c35a457e9c2b2aa189f544e1/shares/d94a8548-2079-4be0-b21c-0a887acd31ca", "rel": "self"},
				{"href": "http://172.18.198.54:8786/16e1ab15c35a457e9c2b2aa189f544e1/shares/d94a8548-2079-4be0-b21c-0a887acd31ca", "rel": "bookmark"}},
			Availability_zone:           "nova",
			Share_network_id:            "f9b2e754-ac01-4466-86e1-5c569424754e",
			Export_locations:            []string{},
			Share_server_id:             "87d8943a-f5da-47a4-b2f2-ddfa6794aa82",
			Snapshot_id:                 "",
			ID:                          "d94a8548-2079-4be0-b21c-0a887acd31ca",
			Size:                        1,
			Share_type:                  "25747776-08e5-494f-ab40-a64b9d20d8f7",
			Share_type_name:             "default",
			Export_location:             "",
			Consistency_group_id:        "9397c191-8427-4661-a2e8-b23820dc01d4",
			Project_id:                  "16e1ab15c35a457e9c2b2aa189f544e1",
			Metadata:                    map[string]string{},
			Status:                      "error",
			Access_rules_status:         "active",
			Description:                 "There is a share description.",
			Host:                        "manila2@generic1#GENERIC1",
			Task_state:                  "",
			Is_public:                   true,
			Snapshot_support:            true,
			Name:                        "My_share",
			Has_replicas:                false,
			Replication_type:            "",
			Created_at:                  createdAt,
			Share_proto:                 "NFS",
			Volume_type:                 "default",
			Source_cgsnapshot_member_id: ""}
		testUtil.Equals(t, expectedShareDetail, shares[0])
	}

	testDetailAllSharesServiceAction(t, "shares/detail", sampleSharesDetailData, anon)
}

func testDetailAllSharesServiceAction(t *testing.T, uriEndsWith string, testData string, shareServiceAction func(Service)) {
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
	shareService, _ := NewService(*sess, *http.DefaultClient, apiServer.URL)
	shareServiceAction(shareService)
}

func TestDeleteShare(t *testing.T) {
	anon := func(shareService Service) {
		volID := "d94a8548-2079-4be0-b21c-0a887acd31ca"
		err := shareService.Delete(volID)
		if err != nil {
			t.Error(err)
		}
	}

	testDeleteShareServiceAction(t, "shares/d94a8548-2079-4be0-b21c-0a887acd31ca", anon)
}

func testDeleteShareServiceAction(t *testing.T, uriEndsWith string, shareServiceAction func(Service)) {
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
	shareService, _ := NewService(*sess, *http.DefaultClient, apiServer.URL)
	shareServiceAction(shareService)
}

var sampleRequestBody = `{"share":{"name":"My_share","size":1,"share_proto":"NFS","share_type":"g-nfs"}}`

var sampleShareData = `{
	"share": {
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
    }
}`

var sampleSharesData = `{
	"shares": [
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
    ]
}`

var sampleShareDetailData = `{
   "share": {
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
    }
}`

var sampleSharesDetailData = `{
	"shares":[
        {
            "links":[{"href": "http://172.18.198.54:8786/v2/16e1ab15c35a457e9c2b2aa189f544e1/shares/d94a8548-2079-4be0-b21c-0a887acd31ca", "rel": "self"},
               	{"href": "http://172.18.198.54:8786/16e1ab15c35a457e9c2b2aa189f544e1/shares/d94a8548-2079-4be0-b21c-0a887acd31ca", "rel": "bookmark"}],
            "availability_zone":"nova",
            "share_network_id":"f9b2e754-ac01-4466-86e1-5c569424754e",
            "export_locations":[],
            "share_server_id":"87d8943a-f5da-47a4-b2f2-ddfa6794aa82",
            "snapshot_id":null,
            "id":"d94a8548-2079-4be0-b21c-0a887acd31ca",
            "size":1,
            "share_type":"25747776-08e5-494f-ab40-a64b9d20d8f7",
            "share_type_name":"default",
            "export_location":null,
            "consistency_group_id":"9397c191-8427-4661-a2e8-b23820dc01d4",
            "project_id":"16e1ab15c35a457e9c2b2aa189f544e1",
            "metadata":{},
            "status":"error",
            "access_rules_status":"active",
            "description":"There is a share description.",
            "host":"manila2@generic1#GENERIC1",
            "task_state":null,
            "is_public":true,
            "snapshot_support":true,
            "name":"My_share",
            "has_replicas":false,
            "replication_type":null,
            "created_at":"2015-09-16T18:19:50.000000",
            "share_proto":"NFS",
            "volume_type":"default",
            "source_cgsnapshot_member_id":null
        },
        {
            "links":[{"href": "http://172.18.198.54:8786/v2/16e1ab15c35a457e9c2b2aa189f544e1/shares/406ea93b-32e9-4907-a117-148b3945749f", "rel": "self"},
               	{"href": "http://172.18.198.54:8786/16e1ab15c35a457e9c2b2aa189f544e1/shares/406ea93b-32e9-4907-a117-148b3945749f", "rel": "bookmark"}],
            "availability_zone":"nova",
            "share_network_id":"f9b2e754-ac01-4466-86e1-5c569424754e",
            "export_locations":["10.254.0.5:/shares/share-50ad5e7b-f6f1-4b78-a651-0812cef2bb67"],
            "share_server_id":"87d8943a-f5da-47a4-b2f2-ddfa6794aa82",
            "snapshot_id":null,
            "id":"406ea93b-32e9-4907-a117-148b3945749f",
            "size":1,
            "share_type":"25747776-08e5-494f-ab40-a64b9d20d8f7",
            "share_type_name":"default",
            "export_location":"10.254.0.5:/shares/share-50ad5e7b-f6f1-4b78-a651-0812cef2bb67",
            "consistency_group_id":"9397c191-8427-4661-a2e8-b23820dc01d4",
            "project_id":"16e1ab15c35a457e9c2b2aa189f544e1",
            "metadata":{},
            "status":"available",
            "access_rules_status":"active",
            "description":"Changed description.",
            "host":"manila2@generic1#GENERIC1",
            "task_state":null,
            "is_public":true,
            "snapshot_support":true,
            "name":"Share1",
            "has_replicas":false,
            "replication_type":null,
            "created_at":"2015-09-16T17:26:28.000000",
            "share_proto":"NFS",
            "volume_type":"default",
            "source_cgsnapshot_member_id":null
        }
    ]
}`
