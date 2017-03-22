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

	Attach(id string, reqBody *AttachBody) error

	Detach(id string, reqBody *DetachBody) error
}

type volumeService struct {
	Session openstack.Session
	Client  http.Client
	URL     string
}

func NewService(
	session openstack.Session,
	client http.Client,
	url string) (*volumeService, error) {

	vs := &volumeService{
		Session: session,
		Client:  client,
		URL:     url,
	}

	return vs, nil
}

type RequestBody struct {
	Name         string `json:"name"`
	VolumeType   string `json:"volume_type"`
	Size         int32  `json:"size"`
	HostName     string `json:"host_name"`
	Mountpoint   string `json:"mountpoint"`
	AttachmentID string `json:"attachment_id"`
}

type CreateBody struct {
	VolumeBody RequestBody `json:"volume"`
}

type AttachBody struct {
	VolumeBody RequestBody `json:"os-attach"`
}

type DetachBody struct {
	VolumeBody RequestBody `json:"os-detach"`
}

// Response is a structure for all properties of
// an volume for a non detailed query
type Response struct {
	ID          string              `json:"id"`
	Name        string              `json:"name"`
	Status      string              `json:"status"`
	Size        int                 `json:"size"`
	Volume_type string              `json:"volume_type"`
	Attachments []map[string]string `json:"attachments"`
}

// DetailResponse is a structure for all properties of
// an volume for a detailed query
type DetailResponse struct {
	ID              string               `json:"id"`
	Attachments     []map[string]string  `json:"attachments"`
	Links           []map[string]string  `json:"links"`
	Metadata        map[string]string    `json:"metadata"`
	Protected       bool                 `json:"protected"`
	Status          string               `json:"status"`
	MigrationStatus string               `json:"migration_status"`
	UserID          string               `json:"user_id"`
	Encrypted       bool                 `json:"encrypted"`
	Multiattach     bool                 `json:"multiattach"`
	CreatedAt       util.RFC8601DateTime `json:"created_at"`
	Description     string               `json:"description"`
	Volume_type     string               `json:"volume_type"`
	Name            string               `json:"name"`
	Source_volid    string               `json:"source_volid"`
	Snapshot_id     string               `json:"snapshot_id"`
	Size            int                  `json:"size"`

	Aavailability_zone  string `json:"availability_zone"`
	Rreplication_status string `json:"replication_status"`
	Consistencygroup_id string `json:"consistencygroup_id"`
}

type VolumeResponse struct {
	Volume Response `json:"volume"`
}

type VolumesResponse struct {
	Volumes []Response `json:"volumes"`
}

type DetailVolumeResponse struct {
	DetailVolume DetailResponse `json:"volume"`
}

type DetailVolumesResponse struct {
	DetailVolumes []DetailResponse `json:"volumes"`
}

func (vs volumeService) Create(reqBody *CreateBody) (Response, error) {
	return createVolume(vs, reqBody)
}

func createVolume(vs volumeService, reqBody *CreateBody) (Response, error) {
	nullResponse := Response{}

	reqURL, err := url.Parse(vs.URL)
	if err != nil {
		log.Println("Parse URL error:", err)
		return nullResponse, err
	}
	urlPostFix := "/volumes"
	reqURL.Path += urlPostFix

	var headers http.Header = http.Header{}
	headers.Set("Content-Type", "application/json")
	body, _ := json.Marshal(reqBody)
	log.Printf("Start POST request to create volume, url = %s, body = %s\n",
		reqURL.String(), body)
	resp, err := vs.Session.Post(reqURL.String(), nil, &headers, &body)
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

	volumeResponse := new(VolumeResponse)
	if err = json.Unmarshal(rbody, volumeResponse); err != nil {
		return nullResponse, err
	}
	return volumeResponse.Volume, nil
}

func (vs volumeService) Show(id string) (DetailResponse, error) {
	return getVolume(vs, id)
}

func getVolume(vs volumeService, id string) (DetailResponse, error) {
	nullResponse := DetailResponse{}

	reqURL, err := url.Parse(vs.URL)
	if err != nil {
		log.Println("Parse URL error:", err)
		return nullResponse, err
	}
	urlPostFix := "/volumes" + "/" + id
	reqURL.Path += urlPostFix

	var headers http.Header = http.Header{}
	headers.Set("Content-Type", "application/json")
	log.Println("Start GET request to get volume, url =", reqURL.String())
	resp, err := vs.Session.Get(reqURL.String(), nil, &headers)
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

	detailVolumeResponse := new(DetailVolumeResponse)
	if err = json.Unmarshal(rbody, detailVolumeResponse); err != nil {
		return nullResponse, err
	}
	return detailVolumeResponse.DetailVolume, nil
}

func (vs volumeService) List() ([]Response, error) {
	return getAllVolumes(vs)
}

func getAllVolumes(vs volumeService) ([]Response, error) {
	nullResponses := make([]Response, 0)

	reqURL, err := url.Parse(vs.URL)
	if err != nil {
		log.Println("Parse URL error:", err)
		return nullResponses, err
	}
	urlPostFix := "/volumes/detail"
	reqURL.Path += urlPostFix

	var headers http.Header = http.Header{}
	headers.Set("Content-Type", "application/json")
	log.Println("Start GET request to get all volumes, url =", reqURL.String())
	resp, err := vs.Session.Get(reqURL.String(), nil, &headers)
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

	volumesResponse := new(VolumesResponse)
	if err = json.Unmarshal(rbody, volumesResponse); err != nil {
		return nullResponses, err
	}
	return volumesResponse.Volumes, nil
}

func (vs volumeService) Detail() ([]DetailResponse, error) {
	return detailAllVolumes(vs)
}

func detailAllVolumes(vs volumeService) ([]DetailResponse, error) {
	nullResponses := make([]DetailResponse, 0)

	reqURL, err := url.Parse(vs.URL)
	if err != nil {
		log.Println("Parse URL error:", err)
		return nullResponses, err
	}
	urlPostFix := "/volumes/detail"
	reqURL.Path += urlPostFix

	var headers http.Header = http.Header{}
	headers.Set("Content-Type", "application/json")
	log.Println("Start GET request to detail all volumes, url =", reqURL.String())
	resp, err := vs.Session.Get(reqURL.String(), nil, &headers)
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

	detailVolumesResponse := new(DetailVolumesResponse)
	if err = json.Unmarshal(rbody, detailVolumesResponse); err != nil {
		return nullResponses, err
	}
	return detailVolumesResponse.DetailVolumes, nil
}

func (vs volumeService) Delete(id string) error {
	return deleteVolume(vs, id)
}

func deleteVolume(vs volumeService, id string) error {
	reqURL, err := url.Parse(vs.URL)
	if err != nil {
		log.Println("Parse URL error:", err)
		return err
	}
	urlPostFix := "/volumes" + "/" + id
	reqURL.Path += urlPostFix

	var headers http.Header = http.Header{}
	headers.Set("Content-Type", "application/json")
	log.Println("Start DELETE request to delete volume, url =", reqURL.String())
	resp, err := vs.Session.Delete(reqURL.String(), nil, &headers)
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

func (vs volumeService) Attach(id string, reqBody *AttachBody) error {
	return attachVolume(vs, id, reqBody)
}

func attachVolume(vs volumeService, id string, reqBody *AttachBody) error {
	reqURL, err := url.Parse(vs.URL)
	if err != nil {
		log.Println("Parse URL error:", err)
		return err
	}
	urlPostFix := "/volumes" + "/" + id + "/action"
	reqURL.Path += urlPostFix

	var headers http.Header = http.Header{}
	headers.Set("Content-Type", "application/json")
	body, _ := json.Marshal(reqBody)
	log.Printf("Start POST request to attach volume, url = %s, body = %s\n",
		reqURL.String(), body)
	resp, err := vs.Session.Post(reqURL.String(), nil, &headers, &body)
	if err != nil {
		log.Println("POST response error:", err)
		return err
	}

	err = util.CheckHTTPResponseStatusCode(resp)
	if err != nil {
		return err
	}

	return nil
}

func (vs volumeService) Detach(id string, reqBody *DetachBody) error {
	return detachVolume(vs, id, reqBody)
}

func detachVolume(vs volumeService, id string, reqBody *DetachBody) error {
	reqURL, err := url.Parse(vs.URL)
	if err != nil {
		log.Println("Parse URL error:", err)
		return err
	}
	urlPostFix := "/volumes" + "/" + id + "/action"
	reqURL.Path += urlPostFix

	var headers http.Header = http.Header{}
	headers.Set("Content-Type", "application/json")
	body, _ := json.Marshal(reqBody)
	log.Printf("Start POST request to detach volume, url = %s, body = %s\n",
		reqURL.String(), body)
	resp, err := vs.Session.Post(reqURL.String(), nil, &headers, &body)
	if err != nil {
		log.Println("PUT response error:", err)
		return err
	}

	err = util.CheckHTTPResponseStatusCode(resp)
	if err != nil {
		return err
	}

	return nil
}
