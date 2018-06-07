// Copyright (c) 2017 Huawei Technologies Co., Ltd. All Rights Reserved.
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
This module implements a entry into the OpenSDS volume controller service.

*/

package volume

import (
	"encoding/json"
	"errors"
	"fmt"

	log "github.com/golang/glog"
	"github.com/opensds/opensds/pkg/dock/client"
	pb "github.com/opensds/opensds/pkg/dock/proto"
	"github.com/opensds/opensds/pkg/model"
	"golang.org/x/net/context"
)

// Controller is an interface for exposing some operations of different volume
// controllers.
type Controller interface {
	CreateVolume(opt *pb.CreateVolumeOpts) (*model.VolumeSpec, error)

	DeleteVolume(opt *pb.DeleteVolumeOpts) error

	ExtendVolume(opt *pb.ExtendVolumeOpts) (*model.VolumeSpec, error)

	CreateVolumeAttachment(opt *pb.CreateAttachmentOpts) (*model.VolumeAttachmentSpec, error)

	DeleteVolumeAttachment(opt *pb.DeleteAttachmentOpts) error

	CreateVolumeSnapshot(opt *pb.CreateVolumeSnapshotOpts) (*model.VolumeSnapshotSpec, error)

	DeleteVolumeSnapshot(opt *pb.DeleteVolumeSnapshotOpts) error

	CreateReplication(opt *pb.CreateReplicationOpts) (*model.ReplicationSpec, error)

	DeleteReplication(opt *pb.DeleteReplicationOpts) error

	EnableReplication(opt *pb.EnableReplicationOpts) error

	DisableReplication(opt *pb.DisableReplicationOpts) error

	FailoverReplication(opt *pb.FailoverReplicationOpts) error

	AttachVolume(opt *pb.AttachVolumeOpts) (string, error)

	DetachVolume(opt *pb.DetachVolumeOpts) error

	CreateVolumeGroup(*pb.CreateVolumeGroupOpts) (*model.VolumeGroupSpec, error)

	UpdateVolumeGroup(*pb.UpdateVolumeGroupOpts) error

	DeleteVolumeGroup(*pb.DeleteVolumeGroupOpts) error

	SetDock(dockInfo *model.DockSpec)
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

func (c *controller) CreateVolume(opt *pb.CreateVolumeOpts) (*model.VolumeSpec, error) {
	if err := c.Client.Connect(c.DockInfo.Endpoint); err != nil {
		log.Error("When connecting dock client:", err)
		return nil, err
	}

	response, err := c.Client.CreateVolume(context.Background(), opt)
	if err != nil {
		log.Error("create volume failed in volume controller:", err)
		return nil, err
	}
	defer c.Client.Close()

	if errorMsg := response.GetError(); errorMsg != nil {
		return nil,
			fmt.Errorf("failed to create volume in volume controller, code: %v, message: %v",
				errorMsg.GetCode(), errorMsg.GetDescription())
	}

	var vol = &model.VolumeSpec{}
	if err = json.Unmarshal([]byte(response.GetResult().GetMessage()), vol); err != nil {
		log.Error("create volume failed in volume controller:", err)
		return nil, err
	}

	return vol, nil

}

func (c *controller) DeleteVolume(opt *pb.DeleteVolumeOpts) error {
	if err := c.Client.Connect(c.DockInfo.Endpoint); err != nil {
		log.Error("When connecting dock client:", err)
		return err
	}

	response, err := c.Client.DeleteVolume(context.Background(), opt)
	if err != nil {
		log.Error("Delete volume failed in volume controller:", err)
		return err
	}
	defer c.Client.Close()

	if errorMsg := response.GetError(); errorMsg != nil {
		return errors.New(errorMsg.GetDescription())
	}

	return nil
}

func (c *controller) ExtendVolume(opt *pb.ExtendVolumeOpts) (*model.VolumeSpec, error) {
	if err := c.Client.Connect(c.DockInfo.Endpoint); err != nil {
		log.Error("When connecting dock client:", err)
		return nil, err
	}

	response, err := c.Client.ExtendVolume(context.Background(), opt)
	if err != nil {
		log.Error("extend volume failed in volume controller:", err)
		return nil, err
	}
	defer c.Client.Close()

	if errorMsg := response.GetError(); errorMsg != nil {
		return nil,
			fmt.Errorf("failed to extend volume in volume controller, code: %v, message: %v",
				errorMsg.GetCode(), errorMsg.GetDescription())
	}

	var vol = &model.VolumeSpec{}
	if err = json.Unmarshal([]byte(response.GetResult().GetMessage()), vol); err != nil {
		log.Error("extend volume failed in volume controller:", err)
		return nil, err
	}

	return vol, nil
}

func (c *controller) CreateVolumeAttachment(opt *pb.CreateAttachmentOpts) (*model.VolumeAttachmentSpec, error) {
	if err := c.Client.Connect(c.DockInfo.Endpoint); err != nil {
		log.Error("When connecting dock client:", err)
		return nil, err
	}

	response, err := c.Client.CreateAttachment(context.Background(), opt)
	if err != nil {
		log.Error("Create volume attachment failed in volume controller:", err)
		return nil, err
	}
	defer c.Client.Close()

	if errorMsg := response.GetError(); errorMsg != nil {
		return nil,
			fmt.Errorf("failed to create volume attachment in volume controller, code: %v, message: %v",
				errorMsg.GetCode(), errorMsg.GetDescription())
	}

	var atc = &model.VolumeAttachmentSpec{}
	if err = json.Unmarshal([]byte(response.GetResult().GetMessage()), atc); err != nil {
		log.Error("create volume attachment failed in volume controller:", err)
		return nil, err
	}

	return atc, nil
}

func (c *controller) DeleteVolumeAttachment(opt *pb.DeleteAttachmentOpts) error {
	if err := c.Client.Connect(c.DockInfo.Endpoint); err != nil {
		log.Error("When connecting dock client:", err)
		return err
	}

	response, err := c.Client.DeleteAttachment(context.Background(), opt)
	if err != nil {
		log.Error("Delete volume attachment failed in volume controller:", err)
		return err
	}
	defer c.Client.Close()

	if errorMsg := response.GetError(); errorMsg != nil {
		return errors.New(errorMsg.GetDescription())
	}

	return nil
}

func (c *controller) CreateVolumeSnapshot(opt *pb.CreateVolumeSnapshotOpts) (*model.VolumeSnapshotSpec, error) {
	if err := c.Client.Connect(c.DockInfo.Endpoint); err != nil {
		log.Error("When connecting dock client:", err)
		return nil, err
	}

	response, err := c.Client.CreateVolumeSnapshot(context.Background(), opt)
	if err != nil {
		log.Error("Create volume snapshot failed in volume controller:", err)
		return nil, err
	}
	defer c.Client.Close()

	if errorMsg := response.GetError(); errorMsg != nil {
		return nil,
			fmt.Errorf("failed to create volume snapshot in volume controller, code: %v, message: %v",
				errorMsg.GetCode(), errorMsg.GetDescription())
	}

	var snp = &model.VolumeSnapshotSpec{}
	if err = json.Unmarshal([]byte(response.GetResult().GetMessage()), snp); err != nil {
		log.Error("create volume snapshot failed in volume controller:", err)
		return nil, err
	}

	return snp, nil
}

func (c *controller) DeleteVolumeSnapshot(opt *pb.DeleteVolumeSnapshotOpts) error {
	if err := c.Client.Connect(c.DockInfo.Endpoint); err != nil {
		log.Error("When connecting dock client:", err)
		return err
	}

	response, err := c.Client.DeleteVolumeSnapshot(context.Background(), opt)
	if err != nil {
		log.Error("Delete volume snapshot failed in volume controller:", err)
		return err
	}
	defer c.Client.Close()

	if errorMsg := response.GetError(); errorMsg != nil {
		return errors.New(errorMsg.GetDescription())
	}

	return nil
}

func (c *controller) CreateReplication(opt *pb.CreateReplicationOpts) (*model.ReplicationSpec, error) {
	if err := c.Client.Connect(c.DockInfo.Endpoint); err != nil {
		log.Error("When connecting dock client:", err)
		return nil, err
	}

	response, err := c.Client.CreateReplication(context.Background(), opt)
	if err != nil {
		log.Error("Create replication failed in volume controller:", err)
		return nil, err
	}
	defer c.Client.Close()

	if errorMsg := response.GetError(); errorMsg != nil {
		return nil,
			fmt.Errorf("failed to create volume snapshot in volume controller, code: %v, message: %v",
				errorMsg.GetCode(), errorMsg.GetDescription())
	}

	var snp = &model.ReplicationSpec{}
	if err = json.Unmarshal([]byte(response.GetResult().GetMessage()), snp); err != nil {
		log.Error("create volume snapshot failed in volume controller:", err)
		return nil, err
	}

	return snp, nil
}

func (c *controller) DeleteReplication(opt *pb.DeleteReplicationOpts) error {
	if err := c.Client.Connect(c.DockInfo.Endpoint); err != nil {
		log.Error("When connecting dock client:", err)
		return err
	}

	response, err := c.Client.DeleteReplication(context.Background(), opt)
	if err != nil {
		log.Error("Delete replication failed in volume controller:", err)
		return err
	}
	defer c.Client.Close()
	
	if errorMsg := response.GetError(); errorMsg != nil {
		return errors.New(errorMsg.GetDescription())
	}

	return nil
}

func (c *controller) EnableReplication(opt *pb.EnableReplicationOpts) error {
	if err := c.Client.Connect(c.DockInfo.Endpoint); err != nil {
		log.Error("When connecting dock client:", err)
		return err
	}

	response, err := c.Client.EnableReplication(context.Background(), opt)
	if err != nil {
		log.Error("Enable replication failed in volume controller:", err)
		return err
	}
	defer c.Client.Close()

	if errorMsg := response.GetError(); errorMsg != nil {
		return errors.New(errorMsg.GetDescription())
	}

	return nil
}

func (c *controller) DisableReplication(opt *pb.DisableReplicationOpts) error {
	if err := c.Client.Connect(c.DockInfo.Endpoint); err != nil {
		log.Error("When connecting dock client:", err)
		return err
	}

	response, err := c.Client.DisableReplication(context.Background(), opt)
	if err != nil {
		log.Error("Disable replication failed in volume controller:", err)
		return err
	}
	defer c.Client.Close()

	if errorMsg := response.GetError(); errorMsg != nil {
		return errors.New(errorMsg.GetDescription())
	}

	return nil
}

func (c *controller) FailoverReplication(opt *pb.FailoverReplicationOpts) error {
	if err := c.Client.Connect(c.DockInfo.Endpoint); err != nil {
		log.Error("When connecting dock client:", err)
		return err
	}

	response, err := c.Client.FailoverReplication(context.Background(), opt)
	if err != nil {
		log.Error("Failover replication failed in volume controller:", err)
		return err
	}
	defer c.Client.Close()

	if errorMsg := response.GetError(); errorMsg != nil {
		return errors.New(errorMsg.GetDescription())
	}

	return nil
}

func (c *controller) AttachVolume(opt *pb.AttachVolumeOpts) (string, error) {
	fmt.Println(c.Client)
	fmt.Println(c.DockInfo)
	if err := c.Client.Connect(c.DockInfo.Endpoint); err != nil {
		log.Error("When connecting dock client:", err)
		return "", err
	}

	response, err := c.Client.AttachVolume(context.Background(), opt)
	if err != nil {
		log.Error("Attach volume failed in volume controller:", err)
		return "", err
	}
	defer c.Client.Close()

	if errorMsg := response.GetError(); errorMsg != nil {
		return "",
			fmt.Errorf("Failed to attach volume in volume controller, code: %v, message: %v",
				errorMsg.GetCode(), errorMsg.GetDescription())
	}

	return response.GetResult().GetMessage(), nil
}

func (c *controller) DetachVolume(opt *pb.DetachVolumeOpts) error {
	if err := c.Client.Connect(c.DockInfo.Endpoint); err != nil {
		log.Error("When connecting dock client:", err)
		return err
	}
	response, err := c.Client.DetachVolume(context.Background(), opt)
	if err != nil {
		log.Error("Detach volume failed in volume controller:", err)
		return err
	}
	defer c.Client.Close()

	if errorMsg := response.GetError(); errorMsg != nil {
		return errors.New(errorMsg.GetDescription())
	}

	return nil
}

func (c *controller) CreateVolumeGroup(opt *pb.CreateVolumeGroupOpts) (*model.VolumeGroupSpec, error) {
	if err := c.Client.Connect(c.DockInfo.Endpoint); err != nil {
		log.Error("When connecting dock client:", err)
		return nil, err
	}

	response, err := c.Client.CreateVolumeGroup(context.Background(), opt)
	if err != nil {
		log.Error("Create volume group failed in volume controller:", err)
		return nil, err
	}
	defer c.Client.Close()

	if errorMsg := response.GetError(); errorMsg != nil {
		return nil,
			fmt.Errorf("failed to create volume group in volume controller, code: %v, message: %v",
				errorMsg.GetCode(), errorMsg.GetDescription())
	}

	var vg = &model.VolumeGroupSpec{}
	if err = json.Unmarshal([]byte(response.GetResult().GetMessage()), vg); err != nil {
		log.Error("create volume group failed in volume controller:", err)
		return nil, err
	}

	return vg, nil
}

func (c *controller) UpdateVolumeGroup(opt *pb.UpdateVolumeGroupOpts) error {
	if err := c.Client.Connect(c.DockInfo.Endpoint); err != nil {
		log.Error("When connecting dock client:", err)
		return err
	}

	response, err := c.Client.UpdateVolumeGroup(context.Background(), opt)
	if err != nil {
		log.Error("Update volume group failed in volume controller:", err)
		return err
	}
	defer c.Client.Close()

	if errorMsg := response.GetError(); errorMsg != nil {
		return fmt.Errorf("Failed to update volume group in volume controller, code: %v, message: %v",
			errorMsg.GetCode(), errorMsg.GetDescription())
	}
	return nil
}

func (c *controller) DeleteVolumeGroup(opt *pb.DeleteVolumeGroupOpts) error {
	if err := c.Client.Connect(c.DockInfo.Endpoint); err != nil {
		log.Error("When connecting dock client:", err)
		return err
	}

	response, err := c.Client.DeleteVolumeGroup(context.Background(), opt)
	if err != nil {
		log.Error("Delete volume group failed in volume controller:", err)
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
