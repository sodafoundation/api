// Copyright 2021 The SODA Foundation Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package client

import (
    "github.com/sodafoundation/api/pkg/model"
    "github.com/sodafoundation/api/pkg/utils/urls"
    "strings"
)

// AkSkBuilder contains request body of handling a AkSk request.
type AkSkBuilder *model.AkSkSpec

// NewAkSkMgr implementation
func NewAkSkMgr(r Receiver, edp string, tenantID string) *AkSkMgr {
    return &AkSkMgr{
        Receiver: r,
        Endpoint: edp,
        TenantID: tenantID,
    }
}

// AkSk implementation
type AkSkMgr struct {
    Receiver
    Endpoint string
    TenantID string
}

/*type credentials struct {
    Blob   int       `json:"blob"`
    ProjectId string `json:"project_Id"`
    CredentialsType string `json:"type"`
    UserId string `json:"user_Id"`
}


func createAkSk(param,options){
}

func deleteAkSk(id,options){
}

func downloadAkSk(id,options){
}


func getAkSkList(){
}

func addKey(projectId string, userId string){
}
*/


// CreateAkSk implementation
func (h *AkSkMgr) CreateAkSk(body AkSkBuilder) (*model.AkSkSpec, error) {
    var res model.AkSkSpec

    url := strings.Join([]string{
        h.Endpoint,
        urls.GenerateAkSkURL(urls.Client, h.TenantID)}, "/")

    if err := h.Recv(url, "POST", body, &res);
        err != nil {
        return nil, err
    }

    return &res, nil
}

// GetAkSk implementation
func (h *AkSkMgr) GetAkSk(ID string) (*model.AkSkSpec, error) {
    var res model.AkSkSpec
    url := strings.Join([]string{
        h.Endpoint,
        urls.GenerateAkSkURL(urls.Client, h.TenantID, ID)}, "/")

    if err := h.Recv(url, "GET", nil, &res); err != nil {
        return nil, err
    }

    return &res, nil
}

// ListAkSks implementation
func (h *AkSkMgr) ListAkSks(args ...interface{}) ([]*model.AkSkSpec, error) {
    url := strings.Join([]string{
        h.Endpoint,
        urls.GenerateAkSkURL(urls.Client, h.TenantID)}, "/")

    param, err := processListParam(args)
    if err != nil {
        return nil, err
    }

    if param != "" {
        url += "?" + param
    }

    var res []*model.AkSkSpec
    if err := h.Recv(url, "GET", nil, &res); err != nil {
        return nil, err
    }
    return res, nil
}


// DeleteAkSk implementation
func (h *AkSkMgr) DeleteAkSk(ID string) error {
    url := strings.Join([]string{
        h.Endpoint,
        urls.GenerateAkSkURL(urls.Client, h.TenantID, ID)}, "/")

    return h.Recv(url, "DELETE", nil, nil)
}

