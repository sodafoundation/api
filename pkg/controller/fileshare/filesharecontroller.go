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

/*
This module implements a entry into the OpenSDS file share controller service.

*/
package fileshare

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	log "github.com/golang/glog"
	"github.com/opensds/opensds/pkg/dock/client"
	"github.com/opensds/opensds/pkg/model"
	pb "github.com/opensds/opensds/pkg/model/proto"
)

// Controller is an interface for exposing some operations of different file share
// controllers.
type Controller interface {
	SetDock(dockInfo *model.DockSpec)

	CreateFileShare(opt *pb.CreateFileShareOpts) (*model.FileShareSpec, error)

	CreateFileShareAcl(opt *pb.CreateFileShareAclOpts) (*model.FileShareAclSpec, error)

	DeleteFileShareAcl(opt *pb.DeleteFileShareAclOpts) error

	DeleteFileShare(opt *pb.DeleteFileShareOpts) error

	CreateFileShareSnapshot(opt *pb.CreateFileShareSnapshotOpts) (*model.FileShareSnapshotSpec, error)

	DeleteFileShareSnapshot(opts *pb.DeleteFileShareSnapshotOpts) error
}

// NewController method creates a controller structure and expose its pointer.
func NewController() Controller {
	return &controller{
		Client: client.NewClient(),
	}
}

type controller struct {
	client.Client
	DockInfo *model.DockSpec
}

func (c *controller) CreateFileShareAcl(opt *pb.CreateFileShareAclOpts) (*model.FileShareAclSpec, error) {
	if err := c.Client.Connect(c.DockInfo.Endpoint); err != nil {
		log.Error("when connecting dock client:", err)
		return nil, err
	}

	response, err := c.Client.CreateFileShareAcl(context.Background(), opt)
	if err != nil {
		log.Error("create file share acl failed in file share controller:", err)
		return nil, err
	}
	defer c.Client.Close()

	if errorMsg := response.GetError(); errorMsg != nil {
		return nil,
			fmt.Errorf("failed to create file share acl in file share controller, code: %v, message: %v",
				errorMsg.GetCode(), errorMsg.GetDescription())
	}

	var fileshare = &model.FileShareAclSpec{}
	if err = json.Unmarshal([]byte(response.GetResult().GetMessage()), fileshare); err != nil {
		log.Error("create file share acl failed in file share controller:", err)
		return nil, err
	}

	return fileshare, nil

}

func (c *controller) DeleteFileShareAcl(opt *pb.DeleteFileShareAclOpts) error {
	if err := c.Client.Connect(c.DockInfo.Endpoint); err != nil {
		log.Error("when connecting dock client:", err)
		return err
	}

	response, err := c.Client.DeleteFileShareAcl(context.Background(), opt)
	if err != nil {
		log.Error("delete file share acl failed in file share controller:", err)
		return err
	}
	defer c.Client.Close()

	if errorMsg := response.GetError(); errorMsg != nil {
		return fmt.Errorf("failed to delete file share acl in file share controller, code: %v, message: %v",
			errorMsg.GetCode(), errorMsg.GetDescription())
	}

	return nil
}

func (c *controller) CreateFileShare(opt *pb.CreateFileShareOpts) (*model.FileShareSpec, error) {
	if err := c.Client.Connect(c.DockInfo.Endpoint); err != nil {
		log.Error("when connecting dock client:", err)
		return nil, err
	}

	log.V(5).Infof("dock create fleshare:  connected to dock endpoint : %+v", c.DockInfo.Endpoint)

	response, err := c.Client.CreateFileShare(context.Background(), (opt))
	if err != nil {
		log.Error("create file share failed in file share controller:", err)
		return nil, err
	}
	defer c.Client.Close()

	log.V(5).Infof("dock create fleshare:  Sent to driver : %+v", c.DockInfo.DriverName)

	if errorMsg := response.GetError(); errorMsg != nil {
		return nil,
			fmt.Errorf("failed to create file share in file share controller, code: %v, message: %v",
				errorMsg.GetCode(), errorMsg.GetDescription())
	}

	var fileshare = &model.FileShareSpec{}
	if err = json.Unmarshal([]byte(response.GetResult().GetMessage()), fileshare); err != nil {
		log.Error("create file share failed in file share controller:", err)
		return nil, err
	}

	return fileshare, nil

}

func (c *controller) DeleteFileShare(opt *pb.DeleteFileShareOpts) error {
	if err := c.Client.Connect(c.DockInfo.Endpoint); err != nil {
		log.Error("when connecting dock client:", err)
		return err
	}

	response, err := c.Client.DeleteFileShare(context.Background(), opt)
	if err != nil {
		log.Error("delete file share failed in file share controller:", err)
		return err
	}
	defer c.Client.Close()

	if errorMsg := response.GetError(); errorMsg != nil {
		return errors.New(errorMsg.GetDescription())
	}

	return nil
}

func (c *controller) SetDock(dockInfo *model.DockSpec) {
	c.DockInfo = dockInfo
}

func (c *controller) CreateFileShareSnapshot(opt *pb.CreateFileShareSnapshotOpts) (*model.FileShareSnapshotSpec, error) {
	if err := c.Client.Connect(c.DockInfo.Endpoint); err != nil {
		log.Error("when connecting dock client:", err)
		return nil, err
	}

	response, err := c.Client.CreateFileShareSnapshot(context.Background(), opt)
	if err != nil {
		log.Error("create file share snapshot failed in file share controller:", err)
		return nil, err
	}
	defer c.Client.Close()

	if errorMsg := response.GetError(); errorMsg != nil {
		return nil,
			fmt.Errorf("failed to create file share snapshot in file share controller, code: %v, message: %v",
				errorMsg.GetCode(), errorMsg.GetDescription())
	}

	var snp = &model.FileShareSnapshotSpec{}
	if err = json.Unmarshal([]byte(response.GetResult().GetMessage()), snp); err != nil {
		log.Error("create file share snapshot failed in file share controller:", err)
		return nil, err
	}

	return snp, nil
}

func (c *controller) DeleteFileShareSnapshot(opt *pb.DeleteFileShareSnapshotOpts) error {
	if err := c.Client.Connect(c.DockInfo.Endpoint); err != nil {
		log.Error("when connecting dock client:", err)
		return err
	}

	response, err := c.Client.DeleteFileShareSnapshot(context.Background(), opt)
	if err != nil {
		log.Error("delete file share snapshot failed in file share controller:", err)
		return err
	}
	defer c.Client.Close()

	if errorMsg := response.GetError(); errorMsg != nil {
		return errors.New(errorMsg.GetDescription())
	}

	return nil
}
