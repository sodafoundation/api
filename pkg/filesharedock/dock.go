// Copyright (c) 2019 OpenSDS Authors. All Rights Reserved.
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
This module implements the entry into operations of storageDock module.
*/

package dock

import (
	_ "encoding/json"
	"github.com/opensds/opensds/contrib/drivers/filesharedrivers"
	"net"

	log "github.com/golang/glog"

	"github.com/opensds/opensds/pkg/filesharedock/discovery"
	pb "github.com/opensds/opensds/pkg/model/fileshareproto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

)

// dockServer is used to implement pb.DockServer
type dockServer struct {
	Port string
	// Discoverer represents the mechanism of DockHub discovering the storage
	// capabilities from different backends.
	Discoverer discovery.DockDiscoverer
	// Driver represents the specified backend resource. This field is used
	// for initializing the specified file share driver.

	Driver filesharedrivers.FileShareDriver
}

// NewDockServer returns a dockServer instance.
func NewDockServer(dockType, port string) *dockServer {
	return &dockServer{
		Port:       port,
		Discoverer: discovery.NewDockDiscoverer(dockType),
	}
}

// Run method would automatically discover dock and pool resources from
// backends, and then start the listen mechanism of dock module.
func (ds *dockServer) Run() error {
	// New Grpc Server
	s := grpc.NewServer()
	// Register dock service.
	pb.RegisterFProvisionDockServer(s, ds)

	// Trigger the discovery and report loop so that the dock service would
	// update the capabilities from backends automatically.
	if err := func() error {
		var err error
		if err = ds.Discoverer.Init(); err != nil {
			return err
		}
		ctx := &discovery.Context{
			StopChan: make(chan bool),
			ErrChan:  make(chan error),
			MetaChan: make(chan string),
		}
		go discovery.DiscoveryAndReport(ds.Discoverer, ctx)
		go func(ctx *discovery.Context) {
			if err = <-ctx.ErrChan; err != nil {
				log.Error("when calling capabilty report method:", err)
				ctx.StopChan <- true
			}
		}(ctx)
		return err
	}(); err != nil {
		return err
	}

	// Listen the dock server port.
	lis, err := net.Listen("tcp", ds.Port)
	if err != nil {
		log.Fatalf("failed to listen: %+v", err)
		return err
	}

	log.Info("Dock server initialized! Start listening on port:", lis.Addr())

	// Start dock server watching loop.
	defer s.Stop()
	return s.Serve(lis)
}

// CreateFileShare implements pb.DockServer.CreateFileShare
func (ds *dockServer) CreateFileShare(ctx context.Context, opt *pb.CreateFileShareOpts) (*pb.GenericResponse, error) {
	// Get the storage drivers and do some initializations.
	ds.Driver = filesharedrivers.Init(opt.GetDriverName())
	defer filesharedrivers.Clean(ds.Driver)

	log.Info("Dock server receive create file share request, vr =", opt)

	fileshare, err := ds.Driver.CreateFileShare(opt)
	if err != nil {
		log.Error("when create file share in dock module:", err)
		return pb.GenericResponseError(err), err
	}


	// TODO: maybe need to update status in DB.
	return pb.GenericResponseResult(fileshare), nil
}


