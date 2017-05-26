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

	"git.openstack.org/openstack/golang-client/openstack"
	"git.openstack.org/openstack/golang-client/util"
)

type Service interface {
	InitializeConnection(id string, reqBody *InitializeBody) (*ConnectionInfo, error)

	CreateVolume(reqBody *VolumeCreateBody) (Response, error)

	ShowVolume(volID string) (DetailResponse, error)

	ListVolumes() ([]Response, error)

	DetailVolumes() ([]DetailResponse, error)

	DeleteVolume(id string) error

	AttachVolume(id string, reqBody *VolumeAttachBody) error

	DetachVolume(id string, reqBody *VolumeDetachBody) error

	CreateSnapshot(reqBody *SnapshotBody) (SnapshotResponse, error)

	ShowSnapshot(volID string) (SnapshotResponse, error)

	ListSnapshots() ([]SnapshotResponse, error)

	DeleteSnapshot(id string) error
}

type volumeService struct {
	Session *openstack.Session
	Client  *http.Client
	URL     string
}

func NewService(
	session *openstack.Session,
	client *http.Client,
	url string) *volumeService {

	return &volumeService{
		Session: session,
		Client:  client,
		URL:     url,
	}
}

type VolumeRequestBody struct {
	Name         string `json:"name"`
	VolumeType   string `json:"volume_type"`
	Size         int32  `json:"size"`
	HostName     string `json:"host_name"`
	Mountpoint   string `json:"mountpoint"`
	AttachmentID string `json:"attachment_id"`
}

type VolumeCreateBody struct {
	VolumeRequestBody `json:"volume"`
}

type VolumeAttachBody struct {
	VolumeRequestBody `json:"os-attach"`
}

type VolumeDetachBody struct {
	VolumeRequestBody `json:"os-detach"`
}

// Response is a structure for all properties of
// a volume for a non detailed query
type Response struct {
	Id                string `json:"id"`
	Name              string `json:"name"`
	Description       string `json:"description"`
	Status            string `json:"status"`
	Size              int    `json:"size"`
	Availability_zone string `json:"availability_zone"`
}

// DetailResponse is a structure for all properties of
// a volume for a detailed query
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

func (vs volumeService) CreateVolume(reqBody *VolumeCreateBody) (Response, error) {
	return createVolume(vs, reqBody)
}

func createVolume(vs volumeService, reqBody *VolumeCreateBody) (Response, error) {
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

func (vs volumeService) ShowVolume(id string) (DetailResponse, error) {
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

func (vs volumeService) ListVolumes() ([]Response, error) {
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

func (vs volumeService) DetailVolumes() ([]DetailResponse, error) {
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

func (vs volumeService) DeleteVolume(id string) error {
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

func (vs volumeService) AttachVolume(id string, reqBody *VolumeAttachBody) error {
	return attachVolume(vs, id, reqBody)
}

func attachVolume(vs volumeService, id string, reqBody *VolumeAttachBody) error {
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

func (vs volumeService) DetachVolume(id string, reqBody *VolumeDetachBody) error {
	return detachVolume(vs, id, reqBody)
}

func detachVolume(vs volumeService, id string, reqBody *VolumeDetachBody) error {
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

type ConnectorProperties struct {
	DoLocalAttach bool   `json:"do_local_attach"`
	Platform      string `json:"platform"`
	OsType        string `json:"os_type"`
	Ip            string `json:"ip"`
	Host          string `json:"host"`
	MultiPath     bool   `json:"multipath"`
	Initiator     string `json:"initiator"`
}

type Connector struct {
	ConnectorProperties `json:"connector"`
}

type InitializeBody struct {
	Connector `json:"os-initialize_connection"`
}

type TerminateBody struct {
	Connector `json:"os-terminateconnection"`
}

type ConnectionInfo struct {
	DriverVolumeType string                 `json:"driver_volume_type"`
	ConnectionData   map[string]interface{} `json:"data"`
}

type ConnectionResponse struct {
	Info ConnectionInfo `json:"connection_info"`
}

func (vs volumeService) InitializeConnection(
	id string,
	reqBody *InitializeBody,
) (*ConnectionInfo, error) {
	return initializeConnection(vs, id, reqBody)
}

func initializeConnection(
	vs volumeService,
	id string,
	reqBody *InitializeBody,
) (*ConnectionInfo, error) {
	var fakeCon = &ConnectionInfo{}

	reqURL, err := url.Parse(vs.URL)
	if err != nil {
		log.Println("Parse URL error:", err)
		return fakeCon, err
	}
	urlPostFix := "/volumes" + "/" + id + "/action"
	reqURL.Path += urlPostFix

	var headers = &http.Header{}
	headers.Set("Content-Type", "application/json")
	body, _ := json.Marshal(reqBody)

	log.Printf("Start POST request to initialize connection, url = %s, body = %s\n",
		reqURL.String(), body)

	resp, err := vs.Session.Post(reqURL.String(), nil, headers, &body)
	if err != nil {
		log.Println("PUT response error:", err)
		return fakeCon, err
	}

	err = util.CheckHTTPResponseStatusCode(resp)
	if err != nil {
		return fakeCon, err
	}

	rbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Read response body failed:", err)
		return fakeCon, err
	}

	var conResp = &ConnectionResponse{}
	if err = json.Unmarshal(rbody, conResp); err != nil {
		return fakeCon, err
	}
	return &conResp.Info, nil
}
