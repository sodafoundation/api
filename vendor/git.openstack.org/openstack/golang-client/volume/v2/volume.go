// Copyright (c) 2014 Hewlett-Packard Development Company, L.P.
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

Volumes and VolumeDetails can be retrieved using the api.

In addition more complex filtering and sort queries can by using the VolumeQueryParameters.

*/
package v2

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"git.openstack.org/openstack/golang-client/openstack"
	"git.openstack.org/openstack/golang-client/util"
)

// Service is a client service that can make
// requests against a OpenStack volume service.
// Below is an example on creating an volume service and getting volumes:
// 	volumeService := volume.VolumeService{Client: *http.DefaultClient, TokenId: tokenId, Url: "http://volumeservicelocation"}
//  volumes:= volumeService.Volumes()
type Service struct {
	Session openstack.Session
	Client  http.Client
	URL     string
}

// Response is a structure for all properties of
// an volume for a non detailed query
type Response struct {
	ID    string              `json:"id"`
	Links []map[string]string `json:"links"`
	Name  string              `json:"name"`
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

// QueryParameters is a structure that
// contains the filter, sort, and paging parameters for
// an volume or volumedetail query.
type QueryParameters struct {
	All_tenant    int64
	Marker        string
	Limit         int64
	SortKey       string
	SortDirection SortDirection
}

// SortDirection of the sort, ascending or descending.
type SortDirection string

const (
	// Desc specifies the sort direction to be descending.
	Desc SortDirection = "desc"
	// Asc specifies the sort direction to be ascending.
	Asc SortDirection = "asc"
)

// Volumes will issue a get request to OpenStack to retrieve the list of volumes.
func (volumeService Service) Volumes() (volume []Response, err error) {
	return volumeService.QueryVolumes(nil)
}

// VolumesDetail will issue a get request to OpenStack to retrieve the list of volumes complete with
// additional details.
func (volumeService Service) VolumesDetail() (volume []DetailResponse, err error) {
	return volumeService.QueryVolumesDetail(nil)
}

// QueryVolumes will issue a get request with the specified VolumeQueryParameters to retrieve the list of
// volumes.
func (volumeService Service) QueryVolumes(queryParameters *QueryParameters) ([]Response, error) {
	volumesContainer := volumesResponse{}
	err := volumeService.queryVolumes(false /*includeDetails*/, &volumesContainer, queryParameters)
	if err != nil {
		return nil, err
	}

	return volumesContainer.Volumes, nil
}

// QueryVolumesDetail will issue a get request with the specified QueryParameters to retrieve the list of
// volumes with additional details.
func (volumeService Service) QueryVolumesDetail(queryParameters *QueryParameters) ([]DetailResponse, error) {
	volumesDetailContainer := volumesDetailResponse{}
	err := volumeService.queryVolumes(true /*includeDetails*/, &volumesDetailContainer, queryParameters)
	if err != nil {
		return nil, err
	}

	return volumesDetailContainer.Volumes, nil
}

func (volumeService Service) queryVolumes(includeDetails bool, volumesResponseContainer interface{}, queryParameters *QueryParameters) error {
	urlPostFix := "/volumes"
	if includeDetails {
		urlPostFix = urlPostFix + "/detail"
	}

	reqURL, err := buildQueryURL(volumeService, queryParameters, urlPostFix)
	if err != nil {
		return err
	}

	var headers http.Header = http.Header{}
	headers.Set("Accept", "application/json")
	resp, err := volumeService.Session.Get(reqURL.String(), nil, &headers)
	if err != nil {
		return err
	}

	err = util.CheckHTTPResponseStatusCode(resp)
	if err != nil {
		return err
	}

	rbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.New("aaa")
	}
	if err = json.Unmarshal(rbody, &volumesResponseContainer); err != nil {
		return err
	}
	return nil
}

func buildQueryURL(volumeService Service, queryParameters *QueryParameters, volumePartialURL string) (*url.URL, error) {
	reqURL, err := url.Parse(volumeService.URL)
	if err != nil {
		return nil, err
	}

	if queryParameters != nil {
		values := url.Values{}
		if queryParameters.All_tenant != 0 {
			values.Set("all_tenant", fmt.Sprintf("%d", queryParameters.All_tenant))
		}
		if queryParameters.Limit != 0 {
			values.Set("limit", fmt.Sprintf("%d", queryParameters.Limit))
		}
		if queryParameters.Marker != "" {
			values.Set("marker", queryParameters.Marker)
		}
		if queryParameters.SortKey != "" {
			values.Set("sort_key", queryParameters.SortKey)
		}
		if queryParameters.SortDirection != "" {
			values.Set("sort_dir", string(queryParameters.SortDirection))
		}

		if len(values) > 0 {
			reqURL.RawQuery = values.Encode()
		}
	}
	reqURL.Path += volumePartialURL

	return reqURL, nil
}

type volumesDetailResponse struct {
	Volumes []DetailResponse `json:"volumes"`
}

type volumesResponse struct {
	Volumes []Response `json:"volumes"`
}
