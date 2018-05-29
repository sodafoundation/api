// Copyright (c) 2018 Huawei Technologies Co., Ltd. All Rights Reserved.
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

package api

import (
	"encoding/json"
	"errors"
	"os"
	"strings"

	c "github.com/opensds/opensds/client"
	"github.com/opensds/opensds/pkg/model"
	. "github.com/opensds/opensds/testutils/collection"
)

var (
	IsFakeClient = false
	TestEp       = "TestEndPoint"
)

func NewFakeClient(config *c.Config) *c.Client {
	os.Setenv("OPENSDS_ENDPOINT", config.Endpoint)
	IsFakeClient = true

	return &c.Client{
		ProfileMgr: &c.ProfileMgr{
			Receiver: NewFakeProfileReceiver(),
			Endpoint: config.Endpoint,
		},
		DockMgr: &c.DockMgr{
			Receiver: NewFakeDockReceiver(),
			Endpoint: config.Endpoint,
		},
		PoolMgr: &c.PoolMgr{
			Receiver: NewFakePoolReceiver(),
			Endpoint: config.Endpoint,
		},
		VolumeMgr: &c.VolumeMgr{
			Receiver: NewFakeVolumeReceiver(),
			Endpoint: config.Endpoint,
		},
		VersionMgr: &c.VersionMgr{
			Receiver: NewFakeVersionReceiver(),
			Endpoint: config.Endpoint,
		},
	}
}

func NewFakeDockReceiver() c.Receiver {
	return &fakeDockReceiver{}
}

type fakeDockReceiver struct{}

func (*fakeDockReceiver) Recv(
	string,
	method string,
	in interface{},
	out interface{},
) error {
	if strings.ToUpper(method) != "GET" {
		return errors.New("method not supported")
	}

	switch out.(type) {
	case *model.DockSpec:
		if err := json.Unmarshal([]byte(ByteDock), out); err != nil {
			return err
		}
		break
	case *[]*model.DockSpec:
		if err := json.Unmarshal([]byte(ByteDocks), out); err != nil {
			return err
		}
		break
	default:
		return errors.New("output format not supported")
	}

	return nil
}

func NewFakePoolReceiver() c.Receiver {
	return &fakePoolReceiver{}
}

type fakePoolReceiver struct{}

func (*fakePoolReceiver) Recv(
	string,
	method string,
	in interface{},
	out interface{},
) error {
	if strings.ToUpper(method) != "GET" {
		return errors.New("method not supported")
	}

	switch out.(type) {
	case *model.StoragePoolSpec:
		if err := json.Unmarshal([]byte(BytePool), out); err != nil {
			return err
		}
		break
	case *[]*model.StoragePoolSpec:
		if err := json.Unmarshal([]byte(BytePools), out); err != nil {
			return err
		}
		break
	default:
		return errors.New("output format not supported")
	}

	return nil
}

func NewFakeProfileReceiver() c.Receiver {
	return &fakeProfileReceiver{}
}

type fakeProfileReceiver struct{}

func (*fakeProfileReceiver) Recv(
	string,
	method string,
	in interface{},
	out interface{},
) error {
	switch strings.ToUpper(method) {
	case "POST", "PUT":
		switch out.(type) {
		case *model.ProfileSpec:
			if err := json.Unmarshal([]byte(ByteProfile), out); err != nil {
				return err
			}
			break
		case *model.ExtraSpec:
			if err := json.Unmarshal([]byte(ByteExtras), out); err != nil {
				return err
			}
			break
		default:
			return errors.New("output format not supported")
		}
		break
	case "GET":
		switch out.(type) {
		case *model.ProfileSpec:
			if err := json.Unmarshal([]byte(ByteProfile), out); err != nil {
				return err
			}
			break
		case *[]*model.ProfileSpec:
			if err := json.Unmarshal([]byte(ByteProfiles), out); err != nil {
				return err
			}
			break
		case *model.ExtraSpec:
			if err := json.Unmarshal([]byte(ByteExtras), out); err != nil {
				return err
			}
			break
		default:
			return errors.New("output format not supported")
		}
		break
	case "DELETE":
		break
	default:
		return errors.New("inputed method format not supported")
	}

	return nil
}

func NewFakeVolumeReceiver() c.Receiver {
	return &fakeVolumeReceiver{}
}

type fakeVolumeReceiver struct{}

func (*fakeVolumeReceiver) Recv(
	string,
	method string,
	in interface{},
	out interface{},
) error {
	switch strings.ToUpper(method) {
	case "POST", "PUT":
		switch out.(type) {
		case *model.VolumeSpec:
			if err := json.Unmarshal([]byte(ByteVolume), out); err != nil {
				return err
			}
			break
		case *model.VolumeAttachmentSpec:
			if err := json.Unmarshal([]byte(ByteAttachment), out); err != nil {
				return err
			}
			break
		case *model.VolumeSnapshotSpec:
			if err := json.Unmarshal([]byte(ByteSnapshot), out); err != nil {
				return err
			}
			break
		default:
			return errors.New("output format not supported")
		}
		break
	case "GET":
		switch out.(type) {
		case *model.VolumeSpec:
			if err := json.Unmarshal([]byte(ByteVolume), out); err != nil {
				return err
			}
			break
		case *[]*model.VolumeSpec:
			if err := json.Unmarshal([]byte(ByteVolumes), out); err != nil {
				return err
			}
			break
		case *model.VolumeAttachmentSpec:
			if err := json.Unmarshal([]byte(ByteAttachment), out); err != nil {
				return err
			}
			break
		case *[]*model.VolumeAttachmentSpec:
			if err := json.Unmarshal([]byte(ByteAttachments), out); err != nil {
				return err
			}
			break
		case *model.VolumeSnapshotSpec:
			if err := json.Unmarshal([]byte(ByteSnapshot), out); err != nil {
				return err
			}
			break
		case *[]*model.VolumeSnapshotSpec:
			if err := json.Unmarshal([]byte(ByteSnapshots), out); err != nil {
				return err
			}
			break
		default:
			return errors.New("output format not supported")
		}
		break
	case "DELETE":
		break
	default:
		return errors.New("inputed method format not supported")
	}

	return nil
}

func NewFakeVersionReceiver() c.Receiver {
	return &fakeVersionReceiver{}
}

type fakeVersionReceiver struct{}

func (*fakeVersionReceiver) Recv(
	string,
	method string,
	in interface{},
	out interface{},
) error {
	switch strings.ToUpper(method) {
	case "GET":
		switch out.(type) {
		case *model.VersionSpec:
			if err := json.Unmarshal([]byte(ByteVersion), out); err != nil {
				return err
			}
			break
		case *[]*model.VersionSpec:
			if err := json.Unmarshal([]byte(ByteVersions), out); err != nil {
				return err
			}
			break
		default:
			return errors.New("output format not supported")
		}
		break
	case "DELETE":
		break
	default:
		return errors.New("inputed method format not supported")
	}

	return nil
}

// ErrorSpec describes Detailed HTTP error response, which consists of a HTTP
// status code, and a custom error message unique for each failure case.
type ErrorSpec struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}
