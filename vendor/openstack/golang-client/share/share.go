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
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"git.openstack.org/openstack/golang-client.git/openstack"
	"git.openstack.org/openstack/golang-client.git/util"
)

type Service interface {
	Create(reqBody *CreateBody) (Response, error)

	Show(volID string) (DetailResponse, error)

	List() ([]Response, error)

	Detail() ([]DetailResponse, error)

	Delete(id string) error
}

type shareService struct {
	Session openstack.Session
	Client  http.Client
	URL     string
}

func NewService(
	session openstack.Session,
	client http.Client,
	url string) (*shareService, error) {

	ss := &shareService{
		Session: session,
		Client:  client,
		URL:     url,
	}

	return ss, nil
}

type RequestBody struct {
	Name        string `json:"name"`
	Size        int32  `json:"size"`
	Share_proto string `json:"share_proto"`
	Share_type  string `json:"share_type"`
}

type CreateBody struct {
	ShareBody RequestBody `json:"share"`
}

// Response is a structure for all properties of
// an share for a non detailed query
type Response struct {
	ID    string              `json:"id"`
	Name  string              `json:"name"`
	Links []map[string]string `json:"links"`
}

// DetailResponse is a structure for all properties of
// an share for a detailed query
type DetailResponse struct {
	Links                       []map[string]string  `json:"links"`
	Availability_zone           string               `json:"availability_zone"`
	Share_network_id            string               `json:"share_network_id"`
	Export_locations            []string             `json:"export_locations"`
	Share_server_id             string               `json:"share_server_id"`
	Snapshot_id                 string               `json:"snapshot_id"`
	ID                          string               `json:"id"`
	Size                        int                  `json:"size"`
	Share_type                  string               `json:"share_type"`
	Share_type_name             string               `json:"share_type_name"`
	Export_location             string               `json:"export_location"`
	Consistency_group_id        string               `json:"consistency_group_id"`
	Project_id                  string               `json:"project_id"`
	Metadata                    map[string]string    `json:"metadata"`
	Status                      string               `json:"status"`
	Access_rules_status         string               `json:"access_rules_status"`
	Description                 string               `json:"description"`
	Host                        string               `json:"host"`
	Task_state                  string               `json:"task_state"`
	Is_public                   bool                 `json:"is_public"`
	Snapshot_support            bool                 `json:"snapshot_support"`
	Name                        string               `json:"name"`
	Has_replicas                bool                 `json:"has_replicas"`
	Replication_type            string               `json:"replication_type"`
	Created_at                  util.RFC8601DateTime `json:"created_at"`
	Share_proto                 string               `json:"share_proto"`
	Volume_type                 string               `json:"volume_type"`
	Source_cgsnapshot_member_id string               `json:"source_cgsnapshot_member_id"`
}

type ShareResponse struct {
	Share Response `json:"share"`
}

type SharesResponse struct {
	Shares []Response `json:"shares"`
}

type DetailShareResponse struct {
	DetailShare DetailResponse `json:"share"`
}

type DetailSharesResponse struct {
	DetailShares []DetailResponse `json:"shares"`
}

func (ss shareService) Create(reqBody *CreateBody) (Response, error) {
	return createShare(ss, reqBody)
}

func createShare(ss shareService, reqBody *CreateBody) (Response, error) {
	nullResponse := Response{}

	reqURL, err := url.Parse(ss.URL)
	if err != nil {
		log.Println("Parse URL error:", err)
		return nullResponse, err
	}
	urlPostFix := "/shares"
	reqURL.Path += urlPostFix

	var headers http.Header = http.Header{}
	headers.Set("Content-Type", "application/json")
	body, _ := json.Marshal(reqBody)
	log.Printf("Start POST request to create share, url = %s, body = %s\n",
		reqURL.String(), body)
	resp, err := ss.Session.Post(reqURL.String(), nil, &headers, &body)
	if err != nil {
		log.Println("POST response error:", err)
		return nullResponse, err
	}

	err = util.CheckHTTPResponseStatusCode(resp)
	if err != nil {
		return nullResponse, err
	}

	rbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Read response body failed:", err)
		return nullResponse, err
	}

	shareResponse := new(ShareResponse)
	if err = json.Unmarshal(rbody, shareResponse); err != nil {
		return nullResponse, err
	}
	return shareResponse.Share, nil
}

func (ss shareService) Show(id string) (DetailResponse, error) {
	return getShare(ss, id)
}

func getShare(ss shareService, id string) (DetailResponse, error) {
	nullResponse := DetailResponse{}

	reqURL, err := url.Parse(ss.URL)
	if err != nil {
		log.Println("Parse URL error:", err)
		return nullResponse, err
	}
	urlPostFix := "/shares" + "/" + id
	reqURL.Path += urlPostFix

	var headers http.Header = http.Header{}
	headers.Set("Content-Type", "application/json")
	log.Println("Start GET request to get share, url =", reqURL.String())
	resp, err := ss.Session.Get(reqURL.String(), nil, &headers)
	if err != nil {
		log.Println("GET response error:", err)
		return nullResponse, err
	}

	err = util.CheckHTTPResponseStatusCode(resp)
	if err != nil {
		return nullResponse, err
	}

	rbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Read response body failed:", err)
		return nullResponse, err
	}

	detailShareResponse := new(DetailShareResponse)
	if err = json.Unmarshal(rbody, detailShareResponse); err != nil {
		return nullResponse, err
	}
	return detailShareResponse.DetailShare, nil
}

func (ss shareService) List() ([]Response, error) {
	return getAllShares(ss)
}

func getAllShares(ss shareService) ([]Response, error) {
	nullResponses := make([]Response, 0)

	reqURL, err := url.Parse(ss.URL)
	if err != nil {
		log.Println("Parse URL error:", err)
		return nullResponses, err
	}
	urlPostFix := "/shares"
	reqURL.Path += urlPostFix

	var headers http.Header = http.Header{}
	headers.Set("Content-Type", "application/json")
	log.Println("Start GET request to get all shares, url =", reqURL.String())
	resp, err := ss.Session.Get(reqURL.String(), nil, &headers)
	if err != nil {
		log.Println("GET response error:", err)
		return nullResponses, err
	}

	err = util.CheckHTTPResponseStatusCode(resp)
	if err != nil {
		return nullResponses, err
	}

	rbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Read response body failed:", err)
		return nullResponses, err
	}

	sharesResponse := new(SharesResponse)
	if err = json.Unmarshal(rbody, sharesResponse); err != nil {
		return nullResponses, err
	}
	return sharesResponse.Shares, nil
}

func (ss shareService) Detail() ([]DetailResponse, error) {
	return detailAllShares(ss)
}

func detailAllShares(ss shareService) ([]DetailResponse, error) {
	nullResponses := make([]DetailResponse, 0)

	reqURL, err := url.Parse(ss.URL)
	if err != nil {
		log.Println("Parse URL error:", err)
		return nullResponses, err
	}
	urlPostFix := "/shares/detail"
	reqURL.Path += urlPostFix

	var headers http.Header = http.Header{}
	headers.Set("Content-Type", "application/json")
	log.Println("Start GET request to detail all shares, url =", reqURL.String())
	resp, err := ss.Session.Get(reqURL.String(), nil, &headers)
	if err != nil {
		log.Println("GET response error:", err)
		return nullResponses, err
	}

	err = util.CheckHTTPResponseStatusCode(resp)
	if err != nil {
		return nullResponses, err
	}

	rbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Read response body failed:", err)
		return nullResponses, err
	}

	detailSharesResponse := new(DetailSharesResponse)
	if err = json.Unmarshal(rbody, detailSharesResponse); err != nil {
		return nullResponses, err
	}
	return detailSharesResponse.DetailShares, nil
}

func (ss shareService) Delete(id string) error {
	return deleteShare(ss, id)
}

func deleteShare(ss shareService, id string) error {
	reqURL, err := url.Parse(ss.URL)
	if err != nil {
		log.Println("Parse URL error:", err)
		return err
	}
	urlPostFix := "/shares" + "/" + id
	reqURL.Path += urlPostFix

	var headers http.Header = http.Header{}
	headers.Set("Content-Type", "application/json")
	log.Println("Start DELETE request to delete share, url =", reqURL.String())
	resp, err := ss.Session.Delete(reqURL.String(), nil, &headers)
	if err != nil {
		log.Println("DELETE response error:", err)
		return err
	}

	err = util.CheckHTTPResponseStatusCode(resp)
	if err != nil {
		return err
	}

	return nil
}
