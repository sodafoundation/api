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

	"git.openstack.org/openstack/golang-client/util"
)

type SnapshotRequestBody struct {
	Name            string `json:"name"`
	VolumeID        string `json:"volume_id"`
	Description     string `json:"description"`
	ForceSnapshoted bool   `json:"force"`
}

type SnapshotBody struct {
	SnapshotRequestBody `json:"snapshot"`
}

// SnapshotResponse is a structure for all properties of
// a volume snapshot for a non detailed query
type SnapshotResponse struct {
	ID        string               `json:"id"`
	Name      string               `json:"name"`
	Status    string               `json:"status"`
	CreatedAt util.RFC8601DateTime `json:"created_at"`
	Volume_id string               `json:"volume_id"`
	Size      int                  `json:"size"`
}

type VolumeSnapshotResponse struct {
	Snapshot SnapshotResponse `json:"snapshot"`
}

type VolumeSnapshotsResponse struct {
	Snapshots []SnapshotResponse `json:"snapshots"`
}

func (vs volumeService) CreateSnapshot(reqBody *SnapshotBody) (SnapshotResponse, error) {
	return createSnapshot(vs, reqBody)
}

func createSnapshot(vs volumeService, reqBody *SnapshotBody) (SnapshotResponse, error) {
	reqURL, err := url.Parse(vs.URL)
	if err != nil {
		log.Println("Parse URL error:", err)
		return SnapshotResponse{}, err
	}
	urlPostFix := "/snapshots"
	reqURL.Path += urlPostFix

	var headers http.Header = http.Header{}
	headers.Set("Content-Type", "application/json")
	body, _ := json.Marshal(reqBody)

	log.Printf("Start POST request to create snapshot, url = %s, body = %s\n",
		reqURL.String(), body)

	resp, err := vs.Session.Post(reqURL.String(), nil, &headers, &body)
	if err != nil {
		log.Println("POST response error:", err)
		return SnapshotResponse{}, err
	}

	err = util.CheckHTTPResponseStatusCode(resp)
	if err != nil {
		return SnapshotResponse{}, err
	}

	rbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Read response body failed:", err)
		return SnapshotResponse{}, err
	}

	snapshotResponse := new(VolumeSnapshotResponse)
	if err = json.Unmarshal(rbody, snapshotResponse); err != nil {
		return SnapshotResponse{}, err
	}
	return snapshotResponse.Snapshot, nil
}

func (vs volumeService) ShowSnapshot(id string) (SnapshotResponse, error) {
	return getSnapshot(vs, id)
}

func getSnapshot(vs volumeService, id string) (SnapshotResponse, error) {
	reqURL, err := url.Parse(vs.URL)
	if err != nil {
		log.Println("Parse URL error:", err)
		return SnapshotResponse{}, err
	}
	urlPostFix := "/snapshots" + "/" + id
	reqURL.Path += urlPostFix

	var headers http.Header = http.Header{}
	headers.Set("Content-Type", "application/json")

	log.Println("Start GET request to get snapshot, url =", reqURL.String())

	resp, err := vs.Session.Get(reqURL.String(), nil, &headers)
	if err != nil {
		log.Println("GET response error:", err)
		return SnapshotResponse{}, err
	}

	err = util.CheckHTTPResponseStatusCode(resp)
	if err != nil {
		return SnapshotResponse{}, err
	}

	rbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Read response body failed:", err)
		return SnapshotResponse{}, err
	}

	snapshotResponse := new(VolumeSnapshotResponse)
	if err = json.Unmarshal(rbody, snapshotResponse); err != nil {
		return SnapshotResponse{}, err
	}
	return snapshotResponse.Snapshot, nil
}

func (vs volumeService) ListSnapshots() ([]SnapshotResponse, error) {
	return getAllSnapshots(vs)
}

func getAllSnapshots(vs volumeService) ([]SnapshotResponse, error) {
	nullResponses := make([]SnapshotResponse, 0)

	reqURL, err := url.Parse(vs.URL)
	if err != nil {
		log.Println("Parse URL error:", err)
		return nullResponses, err
	}
	urlPostFix := "/snapshots"
	reqURL.Path += urlPostFix

	var headers http.Header = http.Header{}
	headers.Set("Content-Type", "application/json")

	log.Println("Start GET request to get all snapshots, url =", reqURL.String())

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

	snapshotsResponse := new(VolumeSnapshotsResponse)
	if err = json.Unmarshal(rbody, snapshotsResponse); err != nil {
		return nullResponses, err
	}
	return snapshotsResponse.Snapshots, nil
}

func (vs volumeService) DeleteSnapshot(id string) error {
	return deleteSnapshot(vs, id)
}

func deleteSnapshot(vs volumeService, id string) error {
	reqURL, err := url.Parse(vs.URL)
	if err != nil {
		log.Println("Parse URL error:", err)
		return err
	}
	urlPostFix := "/snapshots" + "/" + id
	reqURL.Path += urlPostFix

	var headers http.Header = http.Header{}
	headers.Set("Content-Type", "application/json")

	log.Println("Start DELETE request to delete snapshot, url =", reqURL.String())

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
