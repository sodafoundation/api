// Copyright (c) 2019 The OpenSDS Authors.
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
	"strings"

	"github.com/opensds/opensds/pkg/model"
	"github.com/opensds/opensds/pkg/utils/urls"
)

// FileShareBuilder contains request body of handling a fileshare request.
// Currently it's assigned as the pointer of FileShareSpec struct, but it
// could be discussed if it's better to define an interface.
type FileShareBuilder *model.FileShareSpec

// FileShareSnapshotBuilder contains request body of handling a fileshare snapshot request.
// Currently it's assigned as the pointer of FileShareSnapshotSpec struct, but it
// could be discussed if it's better to define an interface.
type FileShareSnapshotBuilder *model.FileShareSnapshotSpec

// FileShareAclBuilder contains request body of handling a fileshare acl request.
// Currently it's assigned as the pointer of FileShareAclSpec struct, but it
// could be discussed if it's better to define an interface.
type FileShareAclBuilder *model.FileShareAclSpec

// NewFileShareMgr implementation
func NewFileShareMgr(r Receiver, edp string, tenantID string) *FileShareMgr {
	return &FileShareMgr{
		Receiver: r,
		Endpoint: edp,
		TenantID: tenantID,
	}
}

// FileShareMgr implementation
type FileShareMgr struct {
	Receiver
	Endpoint string
	TenantID string
}

// CreateFileShare implementation
func (v *FileShareMgr) CreateFileShare(body FileShareBuilder) (*model.FileShareSpec, error) {
	var res model.FileShareSpec

	url := strings.Join([]string{
		v.Endpoint,
		urls.GenerateFileShareURL(urls.Client, v.TenantID)}, "/")

	if err := v.Recv(url, "POST", body, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// DeleteFileShare implementation
func (v *FileShareMgr) DeleteFileShare(ID string) error {
	url := strings.Join([]string{
		v.Endpoint,
		urls.GenerateFileShareURL(urls.Client, v.TenantID, ID)}, "/")

	return v.Recv(url, "DELETE", nil, nil)
}

// GetFileShare implementation
func (v *FileShareMgr) GetFileShare(ID string) (*model.FileShareSpec, error) {
	var res model.FileShareSpec
	url := strings.Join([]string{
		v.Endpoint,
		urls.GenerateFileShareURL(urls.Client, v.TenantID, ID)}, "/")

	if err := v.Recv(url, "GET", nil, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// ListFileShares implementation
func (v *FileShareMgr) ListFileShares(args ...interface{}) ([]*model.FileShareSpec, error) {
	url := strings.Join([]string{
		v.Endpoint,
		urls.GenerateFileShareURL(urls.Client, v.TenantID)}, "/")

	param, err := processListParam(args)
	if err != nil {
		return nil, err
	}

	if param != "" {
		url += "?" + param
	}

	var res []*model.FileShareSpec
	if err := v.Recv(url, "GET", nil, &res); err != nil {
		return nil, err
	}
	return res, nil
}

// UpdateFileShare implementation
func (v *FileShareMgr) UpdateFileShare(ID string, body FileShareBuilder) (*model.FileShareSpec, error) {
	var res model.FileShareSpec
	url := strings.Join([]string{
		v.Endpoint,
		urls.GenerateFileShareURL(urls.Client, v.TenantID, ID)}, "/")

	if err := v.Recv(url, "PUT", body, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// CreateFileShareSnapshot implementation
func (v *FileShareMgr) CreateFileShareSnapshot(body FileShareSnapshotBuilder) (*model.FileShareSnapshotSpec, error) {
	var res model.FileShareSnapshotSpec

	url := strings.Join([]string{
		v.Endpoint,
		urls.GenerateFileShareSnapshotURL(urls.Client, v.TenantID)}, "/")

	if err := v.Recv(url, "POST", body, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// DeleteFileShareSnapshot implementation
func (v *FileShareMgr) DeleteFileShareSnapshot(ID string) error {
	url := strings.Join([]string{
		v.Endpoint,
		urls.GenerateFileShareSnapshotURL(urls.Client, v.TenantID, ID)}, "/")

	return v.Recv(url, "DELETE", nil, nil)
}

// GetFileShareSnapshot implementation
func (v *FileShareMgr) GetFileShareSnapshot(ID string) (*model.FileShareSnapshotSpec, error) {
	var res model.FileShareSnapshotSpec
	url := strings.Join([]string{
		v.Endpoint,
		urls.GenerateFileShareSnapshotURL(urls.Client, v.TenantID, ID)}, "/")

	if err := v.Recv(url, "GET", nil, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// ListFileShareSnapshots implementation
func (v *FileShareMgr) ListFileShareSnapshots(args ...interface{}) ([]*model.FileShareSnapshotSpec, error) {
	url := strings.Join([]string{
		v.Endpoint,
		urls.GenerateFileShareSnapshotURL(urls.Client, v.TenantID)}, "/")

	param, err := processListParam(args)
	if err != nil {
		return nil, err
	}

	if param != "" {
		url += "?" + param
	}

	var res []*model.FileShareSnapshotSpec
	if err := v.Recv(url, "GET", nil, &res); err != nil {
		return nil, err
	}
	return res, nil
}

// UpdateFileShareSnapshot implementation
func (v *FileShareMgr) UpdateFileShareSnapshot(ID string, body FileShareSnapshotBuilder) (*model.FileShareSnapshotSpec, error) {
	var res model.FileShareSnapshotSpec
	url := strings.Join([]string{
		v.Endpoint,
		urls.GenerateFileShareSnapshotURL(urls.Client, v.TenantID, ID)}, "/")

	if err := v.Recv(url, "PUT", body, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// CreateFileShareAcl implementation
func (v *FileShareMgr) CreateFileShareAcl(body FileShareAclBuilder) (*model.FileShareAclSpec, error) {
	var res model.FileShareAclSpec

	url := strings.Join([]string{
		v.Endpoint,
		urls.GenerateFileShareAclURL(urls.Client, v.TenantID)}, "/")

	if err := v.Recv(url, "POST", body, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// DeleteFileShareAcl implementation
func (v *FileShareMgr) DeleteFileShareAcl(ID string) error {
	url := strings.Join([]string{
		v.Endpoint,
		urls.GenerateFileShareAclURL(urls.Client, v.TenantID, ID)}, "/")

	return v.Recv(url, "DELETE", nil, nil)
}

// GetFileShareAcl implementation
func (v *FileShareMgr) GetFileShareAcl(ID string) (*model.FileShareAclSpec, error) {
	var res model.FileShareAclSpec
	url := strings.Join([]string{
		v.Endpoint,
		urls.GenerateFileShareAclURL(urls.Client, v.TenantID, ID)}, "/")

	if err := v.Recv(url, "GET", nil, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// ListFileSharesAcl implementation
func (v *FileShareMgr) ListFileSharesAcl(args ...interface{}) ([]*model.FileShareAclSpec, error) {
	url := strings.Join([]string{
		v.Endpoint,
		urls.GenerateFileShareAclURL(urls.Client, v.TenantID)}, "/")

	param, err := processListParam(args)
	if err != nil {
		return nil, err
	}

	if param != "" {
		url += "?" + param
	}

	var res []*model.FileShareAclSpec
	if err := v.Recv(url, "GET", nil, &res); err != nil {
		return nil, err
	}
	return res, nil
}
