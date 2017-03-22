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

package fake_client

import (
	"log"

	pb "github.com/opensds/opensds/testing/pkg/grpc/fake_opensds"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	address = "localhost:50050"
)

func CreateVolume(vr *pb.VolumeRequest) (*pb.Response, error) {
	// Set up a connection to the orchestration server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Printf("did not connect: %v\n", err)
		return &pb.Response{}, err
	}
	defer conn.Close()

	c := pb.NewDockClient(conn)
	resp, err := c.CreateVolume(context.Background(), vr)
	if err != nil {
		log.Printf("could not create: %v\n", err)
		return &pb.Response{}, err
	}

	return resp, nil
}

func GetVolume(vr *pb.VolumeRequest) (*pb.Response, error) {
	// Set up a connection to the orchestration server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Printf("did not connect: %v\n", err)
		return &pb.Response{}, err
	}
	defer conn.Close()

	c := pb.NewDockClient(conn)
	resp, err := c.GetVolume(context.Background(), vr)
	if err != nil {
		log.Printf("could not get: %v\n", err)
		return &pb.Response{}, err
	}

	return resp, nil
}

func ListVolumes(vr *pb.VolumeRequest) (*pb.Response, error) {
	// Set up a connection to the orchestration server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Printf("did not connect: %v\n", err)
		return &pb.Response{}, err
	}
	defer conn.Close()

	c := pb.NewDockClient(conn)
	resp, err := c.ListVolumes(context.Background(), vr)
	if err != nil {
		log.Printf("could not list: %v\n", err)
		return &pb.Response{}, err
	}

	return resp, nil
}

func DeleteVolume(vr *pb.VolumeRequest) (*pb.Response, error) {
	// Set up a connection to the orchestration server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Printf("did not connect: %v\n", err)
		return &pb.Response{}, err
	}
	defer conn.Close()

	c := pb.NewDockClient(conn)
	resp, err := c.DeleteVolume(context.Background(), vr)
	if err != nil {
		log.Printf("could not delete: %v\n", err)
		return &pb.Response{}, err
	}

	return resp, nil
}

func AttachVolume(vr *pb.VolumeRequest) (*pb.Response, error) {
	// Set up a connection to the orchestration server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Printf("did not connect: %v\n", err)
		return &pb.Response{}, err
	}
	defer conn.Close()

	c := pb.NewDockClient(conn)
	resp, err := c.AttachVolume(context.Background(), vr)
	if err != nil {
		log.Printf("could not attach: %v\n", err)
		return &pb.Response{}, err
	}

	return resp, nil
}

func DetachVolume(vr *pb.VolumeRequest) (*pb.Response, error) {
	// Set up a connection to the orchestration server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Printf("did not connect: %v\n", err)
		return &pb.Response{}, err
	}
	defer conn.Close()

	c := pb.NewDockClient(conn)
	resp, err := c.DetachVolume(context.Background(), vr)
	if err != nil {
		log.Printf("could not detach: %v\n", err)
		return &pb.Response{}, err
	}

	return resp, nil
}

func MountVolume(vr *pb.VolumeRequest) (*pb.Response, error) {
	// Set up a connection to the orchestration server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Printf("did not connect: %v\n", err)
		return &pb.Response{}, err
	}
	defer conn.Close()

	c := pb.NewDockClient(conn)
	resp, err := c.MountVolume(context.Background(), vr)
	if err != nil {
		log.Printf("could not mount: %v\n", err)
		return &pb.Response{}, err
	}

	return resp, nil
}

func UnmountVolume(vr *pb.VolumeRequest) (*pb.Response, error) {
	// Set up a connection to the orchestration server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Printf("did not connect: %v\n", err)
		return &pb.Response{}, err
	}
	defer conn.Close()

	c := pb.NewDockClient(conn)
	resp, err := c.UnmountVolume(context.Background(), vr)
	if err != nil {
		log.Printf("could not unmount: %v\n", err)
		return &pb.Response{}, err
	}

	return resp, nil
}

func CreateShare(sr *pb.ShareRequest) (*pb.Response, error) {
	// Set up a connection to the orchestration server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Printf("did not connect: %v\n", err)
		return &pb.Response{}, err
	}
	defer conn.Close()

	c := pb.NewDockClient(conn)
	resp, err := c.CreateShare(context.Background(), sr)
	if err != nil {
		log.Printf("could not create: %v\n", err)
		return &pb.Response{}, err
	}

	return resp, nil
}

func GetShare(sr *pb.ShareRequest) (*pb.Response, error) {
	// Set up a connection to the orchestration server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Printf("did not connect: %v\n", err)
		return &pb.Response{}, err
	}
	defer conn.Close()

	c := pb.NewDockClient(conn)
	resp, err := c.GetShare(context.Background(), sr)
	if err != nil {
		log.Printf("could not get: %v\n", err)
		return &pb.Response{}, err
	}

	return resp, nil
}

func ListShares(sr *pb.ShareRequest) (*pb.Response, error) {
	// Set up a connection to the orchestration server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Printf("did not connect: %v\n", err)
		return &pb.Response{}, err
	}
	defer conn.Close()

	c := pb.NewDockClient(conn)
	resp, err := c.ListShares(context.Background(), sr)
	if err != nil {
		log.Printf("could not list: %v\n", err)
		return &pb.Response{}, err
	}

	return resp, nil
}

func DeleteShare(sr *pb.ShareRequest) (*pb.Response, error) {
	// Set up a connection to the orchestration server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Printf("did not connect: %v\n", err)
		return &pb.Response{}, err
	}
	defer conn.Close()

	c := pb.NewDockClient(conn)
	resp, err := c.DeleteShare(context.Background(), sr)
	if err != nil {
		log.Printf("could not delete: %v\n", err)
		return &pb.Response{}, err
	}

	return resp, nil
}
