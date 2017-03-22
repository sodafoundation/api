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

The CRUD operation of volumes can be retrieved using the api. Right now only

Show and List methods can work.

*/

package v3

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"

	"git.openstack.org/openstack/golang-client/openstack"
	"git.openstack.org/openstack/golang-client/util"
)

type Service struct {
	Session openstack.Session
	Client  http.Client
	URL     string
}

type RequestBody struct {
	// The volume name [OPTIONAL]
	Name string `json:"name"`
	// The size of the volume, in gibibytes (GiB) [REQUIRED]
	Size int `json:"size"`
}

type Body struct {
	VolumeBody RequestBody `json:"volume"`
}

// Response is a structure for all properties of
// an volume for a non detailed query
type Response struct {
	ID   string `json:"id"`
	Name string `json:"name"`

	Consistencygroup_id string `json:"consistencygroup_id"`
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
	Size            int64                `json:"size"`

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

func (volumeService Service) Create(reqBody *Body) (Response, error) {
	return volumeService.createVolume(reqBody)
}

func (volumeService Service) createVolume(reqBody *Body) (Response, error) {
	nullResponse := Response{}

	reqURL, err := url.Parse(volumeService.URL)
	if err != nil {
		return nullResponse, err
	}
	urlPostFix := "/volumes"
	reqURL.Path += urlPostFix

	var headers http.Header = http.Header{}
	headers.Set("Content-Type", "application/json")
	body, _ := json.Marshal(reqBody)
	resp, err := volumeService.Session.Post(reqURL.String(), nil, &headers, &body)
	if err != nil {
		return nullResponse, err
	}

	err = util.CheckHTTPResponseStatusCode(resp)
	if err != nil {
		return nullResponse, err
	}

	rbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nullResponse, errors.New("aaa")
	}

	volumeResponse := new(VolumeResponse)
	if err = json.Unmarshal(rbody, volumeResponse); err != nil {
		return nullResponse, err
	}
	return volumeResponse.Volume, nil
}

func (volumeService Service) Show(id string) (Response, error) {
	return volumeService.getVolume(id)
}

func (volumeService Service) getVolume(id string) (Response, error) {
	nullResponse := Response{}

	reqURL, err := url.Parse(volumeService.URL)
	if err != nil {
		return nullResponse, err
	}
	urlPostFix := "/volumes" + "/" + id
	reqURL.Path += urlPostFix

	var headers http.Header = http.Header{}
	headers.Set("Content-Type", "application/json")
	resp, err := volumeService.Session.Get(reqURL.String(), nil, &headers)
	if err != nil {
		return nullResponse, err
	}

	err = util.CheckHTTPResponseStatusCode(resp)
	if err != nil {
		return nullResponse, err
	}

	rbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nullResponse, errors.New("aaa")
	}

	volumeResponse := new(VolumeResponse)
	if err = json.Unmarshal(rbody, volumeResponse); err != nil {
		return nullResponse, err
	}
	return volumeResponse.Volume, nil
}

func (volumeService Service) List() ([]Response, error) {
	return volumeService.getAllVolumes()
}

func (volumeService Service) getAllVolumes() ([]Response, error) {
	nullResponses := make([]Response, 0)

	reqURL, err := url.Parse(volumeService.URL)
	if err != nil {
		return nullResponses, err
	}
	urlPostFix := "/volumes"
	reqURL.Path += urlPostFix

	var headers http.Header = http.Header{}
	headers.Set("Content-Type", "application/json")
	resp, err := volumeService.Session.Get(reqURL.String(), nil, &headers)
	if err != nil {
		return nullResponses, err
	}

	err = util.CheckHTTPResponseStatusCode(resp)
	if err != nil {
		return nullResponses, err
	}

	rbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nullResponses, errors.New("aaa")
	}

	volumesResponse := new(VolumesResponse)
	if err = json.Unmarshal(rbody, volumesResponse); err != nil {
		return nullResponses, err
	}
	return volumesResponse.Volumes, nil
}

func (volumeService Service) Update(id string, reqBody *Body) (Response, error) {
	return volumeService.updateVolume(id, reqBody)
}

func (volumeService Service) updateVolume(id string, reqBody *Body) (Response, error) {
	nullResponse := Response{}

	reqURL, err := url.Parse(volumeService.URL)
	if err != nil {
		return nullResponse, err
	}
	urlPostFix := "/volumes" + "/" + id
	reqURL.Path += urlPostFix

	var headers http.Header = http.Header{}
	headers.Set("Content-Type", "application/json")
	body, _ := json.Marshal(reqBody)
	resp, err := volumeService.Session.Put(reqURL.String(), nil, &headers, &body)
	if err != nil {
		return nullResponse, err
	}

	err = util.CheckHTTPResponseStatusCode(resp)
	if err != nil {
		return nullResponse, err
	}

	rbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nullResponse, errors.New("aaa")
	}

	volumeResponse := new(VolumeResponse)
	if err = json.Unmarshal(rbody, volumeResponse); err != nil {
		return nullResponse, err
	}
	return volumeResponse.Volume, nil
}

func (volumeService Service) Delete(id string) error {
	return volumeService.deleteVolume(id)
}

func (volumeService Service) deleteVolume(id string) error {
	reqURL, err := url.Parse(volumeService.URL)
	if err != nil {
		return err
	}
	urlPostFix := "/volumes" + "/" + id
	reqURL.Path += urlPostFix

	var headers http.Header = http.Header{}
	headers.Set("Content-Type", "application/json")
	resp, err := volumeService.Session.Delete(reqURL.String(), nil, &headers)
	if err != nil {
		return err
	}

	err = util.CheckHTTPResponseStatusCode(resp)
	if err != nil {
		return err
	}

	return nil
}
