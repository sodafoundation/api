/*
 *
 * Copyright 2015, Google Inc.
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are
 * met:
 *
 *     * Redistributions of source code must retain the above copyright
 * notice, this list of conditions and the following disclaimer.
 *     * Redistributions in binary form must reproduce the above
 * copyright notice, this list of conditions and the following disclaimer
 * in the documentation and/or other materials provided with the
 * distribution.
 *     * Neither the name of Google Inc. nor the names of its
 * contributors may be used to endorse or promote products derived from
 * this software without specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
 * "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
 * LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
 * A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
 * OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
 * SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
 * LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
 * DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
 * THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
 * OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
 */

package client

import (
	"encoding/json"
	"log"

	pb "github.com/opensds/opensds/pkg/grpc/opensds"
	api "github.com/opensds/opensds/pkg/model"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	DOCK_PORT = ":50050"
)

func NewDockClient(dockInfo string) (pb.DockClient, *grpc.ClientConn, error) {
	// Get Dock endpoint from dock info.
	var dck = &api.DockSpec{}
	if err := json.Unmarshal([]byte(dockInfo), dck); err != nil {
		log.Println("[Error] When parsing dock info:", err)
	}

	// Generate Dock server address.
	address := dck.Endpoint + DOCK_PORT

	// Set up a connection to the Dock server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Printf("did not connect: %+v\n", err)
		return nil, nil, err
	}

	return pb.NewDockClient(conn), conn, nil
}

func CreateVolume(vr *pb.DockRequest) (*pb.DockResponse, error) {
	c, conn, err := NewDockClient(vr.GetDockInfo())
	if err != nil {
		log.Printf("get dock client failed: %+v\n", err)
		return &pb.DockResponse{}, err
	}
	defer conn.Close()

	resp, err := c.CreateVolume(context.Background(), vr)
	if err != nil {
		log.Printf("could not create: %+v\n", err)
		return &pb.DockResponse{}, err
	}

	log.Println("Dock client receive create volume response, vr =", resp)

	return resp, nil
}

func GetVolume(vr *pb.DockRequest) (*pb.DockResponse, error) {
	c, conn, err := NewDockClient(vr.GetDockInfo())
	if err != nil {
		log.Printf("get dock client failed: %+v\n", err)
		return &pb.DockResponse{}, err
	}
	defer conn.Close()

	resp, err := c.GetVolume(context.Background(), vr)
	if err != nil {
		log.Printf("could not get: %+v\n", err)
		return &pb.DockResponse{}, err
	}

	log.Println("Dock client receive get volume response, vr =", resp)

	return resp, nil
}

func DeleteVolume(vr *pb.DockRequest) (*pb.DockResponse, error) {
	c, conn, err := NewDockClient(vr.GetDockInfo())
	if err != nil {
		log.Printf("get dock client failed: %+v\n", err)
		return &pb.DockResponse{}, err
	}
	defer conn.Close()

	resp, err := c.DeleteVolume(context.Background(), vr)
	if err != nil {
		log.Printf("could not delete: %+v\n", err)
		return &pb.DockResponse{}, err
	}

	log.Println("Dock client receive delete volume response, vr =", resp)

	return resp, nil
}

func CreateVolumeAttachment(vr *pb.DockRequest) (*pb.DockResponse, error) {
	c, conn, err := NewDockClient(vr.GetDockInfo())
	if err != nil {
		log.Printf("get dock client failed: %+v\n", err)
		return &pb.DockResponse{}, err
	}
	defer conn.Close()

	resp, err := c.CreateVolumeAttachment(context.Background(), vr)
	if err != nil {
		log.Printf("could not create volume attachment: %+v\n", err)
		return &pb.DockResponse{}, err
	}

	log.Println("Dock client receive create volume attachment response, vr =", resp)

	return resp, nil
}

func UpdateVolumeAttachment(vr *pb.DockRequest) (*pb.DockResponse, error) {
	c, conn, err := NewDockClient(vr.GetDockInfo())
	if err != nil {
		log.Printf("get dock client failed: %+v\n", err)
		return &pb.DockResponse{}, err
	}
	defer conn.Close()

	resp, err := c.UpdateVolumeAttachment(context.Background(), vr)
	if err != nil {
		log.Printf("could not update volume attachment: %+v\n", err)
		return &pb.DockResponse{}, err
	}

	log.Println("Dock client receive update volume attachment response, vr =", resp)

	return resp, nil
}

func DeleteVolumeAttachment(vr *pb.DockRequest) (*pb.DockResponse, error) {
	c, conn, err := NewDockClient(vr.GetDockInfo())
	if err != nil {
		log.Printf("get dock client failed: %+v\n", err)
		return &pb.DockResponse{}, err
	}
	defer conn.Close()

	resp, err := c.DeleteVolumeAttachment(context.Background(), vr)
	if err != nil {
		log.Printf("could not delete volume attachment: %+v\n", err)
		return &pb.DockResponse{}, err
	}

	log.Println("Dock client receive delete volume attachment response, vr =", resp)

	return resp, nil
}

func CreateVolumeSnapshot(vr *pb.DockRequest) (*pb.DockResponse, error) {
	c, conn, err := NewDockClient(vr.GetDockInfo())
	if err != nil {
		log.Printf("get dock client failed: %+v\n", err)
		return &pb.DockResponse{}, err
	}
	defer conn.Close()

	resp, err := c.CreateVolumeSnapshot(context.Background(), vr)
	if err != nil {
		log.Printf("could not create: %+v\n", err)
		return &pb.DockResponse{}, err
	}

	log.Println("Dock client receive create volume snapshot response, vr =", resp)

	return resp, nil
}

func GetVolumeSnapshot(vr *pb.DockRequest) (*pb.DockResponse, error) {
	c, conn, err := NewDockClient(vr.GetDockInfo())
	if err != nil {
		log.Printf("get dock client failed: %+v\n", err)
		return &pb.DockResponse{}, err
	}
	defer conn.Close()

	resp, err := c.GetVolumeSnapshot(context.Background(), vr)
	if err != nil {
		log.Printf("could not get: %+v\n", err)
		return &pb.DockResponse{}, err
	}

	log.Println("Dock client receive get volume snapshot response, vr =", resp)

	return resp, nil
}

func DeleteVolumeSnapshot(vr *pb.DockRequest) (*pb.DockResponse, error) {
	c, conn, err := NewDockClient(vr.GetDockInfo())
	if err != nil {
		log.Printf("get dock client failed: %+v\n", err)
		return &pb.DockResponse{}, err
	}
	defer conn.Close()

	resp, err := c.DeleteVolumeSnapshot(context.Background(), vr)
	if err != nil {
		log.Printf("could not delete: %+v\n", err)
		return &pb.DockResponse{}, err
	}

	log.Println("Dock client receive delete volume snapshot response, vr =", resp)

	return resp, nil
}
