// Copyright (c) 2017 OpenSDS Authors.
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

package csi

import (
	csipb "github.com/container-storage-interface/spec/lib/go/csi"
	log "github.com/golang/glog"
	"google.golang.org/grpc"
)

type Client interface {
	csipb.ControllerClient

	Close()
}

type client struct {
	csipb.ControllerClient

	*grpc.ClientConn
}

func NewClient(edp string) Client {
	// Set up a connection to the Controller server.
	conn, err := GetCSIClientConn(edp)
	if err != nil {
		log.Errorf("did not connect: %+v\n", err)
	}
	// Create a controller client via the connection.
	dc := csipb.NewControllerClient(conn)

	return &client{
		ControllerClient: dc,
		ClientConn:       conn,
	}
}

func (c *client) Close() {
	c.ClientConn.Close()
}
