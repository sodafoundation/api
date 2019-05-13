// Copyright 2019 The OpenSDS Authors.
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

package nfs

import (
	"github.com/opensds/opensds/pkg/model"
	pb "github.com/opensds/opensds/pkg/model/proto"
	data "github.com/opensds/opensds/testutils/collection"
)

// NFSDriver
type NFSDriver struct{}

// Setup
func (*NFSDriver) Setup() error { return nil }

// Unset
func (*NFSDriver) Unset() error { return nil }

// ListPools
func (*NFSDriver) ListPools() ([]*model.StoragePoolSpec, error) {
	var pols []*model.StoragePoolSpec

	for i := range data.SamplePools {
		pols = append(pols, &data.SamplePools[i])
	}
	return pols, nil
}

// CreateFileShare
func (d *NFSDriver) CreateFileShare(opt *pb.CreateFileShareOpts) (*model.FileShareSpec, error) {
	return &data.SampleFileShares[0], nil
}

// DeleteFileShare
func (d *NFSDriver) DeleteFileShare(opt *pb.DeleteFileShareOpts) (*model.FileShareSpec, error) {
	return nil, nil
}
