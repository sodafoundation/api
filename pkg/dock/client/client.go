// Copyright 2017 The OpenSDS Authors.
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
	"github.com/opensds/opensds/pkg/model"
	"google.golang.org/grpc"
)

type Client interface {
	pb.DockClient
	Update(dockInfo *model.DockSpec) error
	Close()
}

type client struct {
	pb.DockClient
	*grpc.ClientConn

	TargetPlace string
}

func NewClient() Client { return &client{} }

func (c *client) Update(dockInfo *model.DockSpec) error {

	// Set up a connection to the Dock server.
	conn, err := grpc.Dial(dockInfo.GetEndpoint(), grpc.WithInsecure())
	if err != nil {
		log.Errorf("did not connect: %+v\n", err)
		return err
	}
	// Create a dock client via the connection.
	dc := pb.NewDockClient(conn)

	c.DockClient = dc
	c.ClientConn = conn
	c.TargetPlace = dockInfo.GetEndpoint()

	return nil
}

func (c *client) Close() {
	c.ClientConn.Close()
}
