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

package client

import (
	log "github.com/golang/glog"
	pb "github.com/opensds/opensds/pkg/dock/proto"
	"google.golang.org/grpc"
)

// Client interface provides an abstract description about how to interact
// with gRPC client. Besides some nested methods defined in pb.DockClient,
// Client also exposes two methods: Connect() and Close(), for which callers
// can easily open and close gRPC connection.
type Client interface {
	pb.ProvisionDockClient
	pb.AttachDockClient

	Connect(edp string) error

	Close()
}

// client structure is one implementation of Client interface and will be
// called in real environment. There would be more other kind of connection
// in the long run.
type client struct {
	pb.ProvisionDockClient
	pb.AttachDockClient
	*grpc.ClientConn
}

func NewClient() Client { return &client{} }

func (c *client) Connect(edp string) error {
	// Set up a connection to the Dock server.
	conn, err := grpc.Dial(edp, grpc.WithInsecure())
	if err != nil {
		log.Errorf("did not connect: %+v\n", err)
		return err
	}
	// Create dock client via the connection.
	c.ProvisionDockClient = pb.NewProvisionDockClient(conn)
	c.AttachDockClient = pb.NewAttachDockClient(conn)
	c.ClientConn = conn

	return nil
}

func (c *client) Close() {
	c.ClientConn.Close()
}
