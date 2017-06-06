// Copyright (c) 2016 Huawei Technologies Co., Ltd. All Rights Reserved.
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

/*
This module implements the entry into operations of storageDock module.

*/

package api

import (
	"encoding/json"
	"log"
	"strings"

	api "github.com/opensds/opensds/pkg/api/v1"
	db "github.com/opensds/opensds/pkg/db/api"
	dock "github.com/opensds/opensds/pkg/dock/volume"
	pb "github.com/opensds/opensds/pkg/grpc/opensds"
)

func CreateVolume(vr *pb.VolumeRequest) (*pb.Response, error) {
	var dck = &api.Dock{}
	if err := json.Unmarshal([]byte(vr.GetDockInfo()), dck); err != nil {
		log.Println("[Error] When parsing dock info:", err)
	}

	vol, err := dock.CreateVolume(dck.DriverName, vr.GetVolumeName(), vr.GetSize())
	if err != nil {
		log.Println("[Error] When create volume in dock module:", err)
		return &pb.Response{}, err
	}
	vol.PoolName = vr.GetPoolName()

	result, err := db.CreateVolume(vol)
	if err != nil {
		log.Println("[Error] When create volume in db module:", err)
		return &pb.Response{}, err
	}

	volBody, _ := json.Marshal(result)
	return &pb.Response{
		Status:  "Success",
		Message: string(volBody),
	}, nil
}

func GetVolume(vr *pb.VolumeRequest) (*pb.Response, error) {
	var dck = &api.Dock{}
	if err := json.Unmarshal([]byte(vr.GetDockInfo()), dck); err != nil {
		log.Println("[Error] When parsing dock info:", err)
	}

	result, err := dock.GetVolume(dck.DriverName, vr.GetVolumeId())
	if err != nil {
		log.Println("[Error] When get volume in dock module:", err)
		return &pb.Response{}, err
	}

	volBody, _ := json.Marshal(result)
	return &pb.Response{
		Status:  "Success",
		Message: string(volBody),
	}, nil
}

func DeleteVolume(vr *pb.VolumeRequest) (*pb.Response, error) {
	var dck = &api.Dock{}
	if err := json.Unmarshal([]byte(vr.GetDockInfo()), dck); err != nil {
		log.Println("[Error] When parsing dock info:", err)
	}

	if err := dock.DeleteVolume(dck.DriverName, vr.GetVolumeId()); err != nil {
		log.Println("Error occured in dock module when delete volume:", err)
		return &pb.Response{}, err
	}

	if err := db.DeleteVolume(vr.GetVolumeId()); err != nil {
		log.Println("Error occured in dock module when delete volume in db:", err)
		return &pb.Response{}, err
	}

	return &pb.Response{
		Status:  "Success",
		Message: "Delete volume success",
	}, nil
}

func CreateVolumeAttachment(vr *pb.VolumeRequest) (*pb.Response, error) {
	var dck, hostInfo = &api.Dock{}, &api.HostInfo{}
	if err := json.Unmarshal([]byte(vr.GetDockInfo()), dck); err != nil {
		log.Println("[Error] When parsing dock info:", err)
	}
	if err := json.Unmarshal([]byte(vr.GetHostInfo()), hostInfo); err != nil {
		log.Println("Error occured in dock module when parsing host info:", err)
		return &pb.Response{}, err
	}

	atc, err := dock.CreateVolumeAttachment(dck.DriverName, vr.GetVolumeId(), vr.GetDoLocalAttach(), vr.GetMultiPath(), hostInfo)
	if err != nil {
		log.Println("Error occured in dock module when create volume attachment:", err)
		return &pb.Response{}, err
	}

	result, err := db.CreateVolumeAttachment(vr.GetVolumeId(), atc)
	if err != nil {
		log.Println("Error occured in dock module when create volume attachment in db:", err)
		return &pb.Response{}, err
	}

	atcBody, _ := json.Marshal(result)
	return &pb.Response{
		Status:  "Success",
		Message: string(atcBody),
	}, nil
}

func UpdateVolumeAttachment(vr *pb.VolumeRequest) (*pb.Response, error) {
	var dck, hostInfo = &api.Dock{}, &api.HostInfo{}
	if err := json.Unmarshal([]byte(vr.GetDockInfo()), dck); err != nil {
		log.Println("[Error] When parsing dock info:", err)
	}
	if err := json.Unmarshal([]byte(vr.GetHostInfo()), hostInfo); err != nil {
		log.Println("Error occured in dock module when parsing host info:", err)
		return &pb.Response{}, err
	}

	if err := dock.UpdateVolumeAttachment(dck.DriverName, vr.GetVolumeId(), hostInfo.Host, vr.GetMountpoint()); err != nil {
		log.Println("Error occured in dock module when update volume attachment:", err)
		return &pb.Response{}, err
	}

	result, err := db.UpdateVolumeAttachment(vr.GetVolumeId(), vr.GetAttachmentId(), vr.GetMountpoint(), hostInfo)
	if err != nil {
		log.Println("Error occured in dock module when update volume attachment in db:", err)
		return &pb.Response{}, err
	}

	atcBody, _ := json.Marshal(result)
	return &pb.Response{
		Status:  "Success",
		Message: string(atcBody),
	}, nil
}

func DeleteVolumeAttachment(vr *pb.VolumeRequest) (*pb.Response, error) {
	var dck = &api.Dock{}
	if err := json.Unmarshal([]byte(vr.GetDockInfo()), dck); err != nil {
		log.Println("[Error] When parsing dock info:", err)
	}

	if err := dock.DeleteVolumeAttachment(dck.DriverName, vr.GetVolumeId()); err != nil {
		log.Println("Error occured in dock module when delete volume attachment:", err)
		if strings.Contains(err.Error(), "The status of volume is not in-use") {
			if err = db.DeleteVolumeAttachment(vr.GetVolumeId(), vr.GetAttachmentId()); err != nil {
				log.Println("Error occured in dock module when delete volume attachment in db:", err)
				return &pb.Response{}, err
			}
		}
		return &pb.Response{}, nil
	}

	if err := db.DeleteVolumeAttachment(vr.GetVolumeId(), vr.GetAttachmentId()); err != nil {
		log.Println("Error occured in dock module when delete volume attachment in db:", err)
		return &pb.Response{}, err
	}

	return &pb.Response{
		Status:  "Success",
		Message: "Delete volume attachment success",
	}, nil
}

func CreateVolumeSnapshot(vr *pb.VolumeRequest) (*pb.Response, error) {
	var dck = &api.Dock{}
	if err := json.Unmarshal([]byte(vr.GetDockInfo()), dck); err != nil {
		log.Println("[Error] When parsing dock info:", err)
	}

	snp, err := dock.CreateSnapshot(dck.DriverName,
		vr.GetSnapshotName(),
		vr.GetVolumeId(),
		vr.GetSnapshotDescription())
	if err != nil {
		log.Println("Error occured in dock module when create snapshot:", err)
		return &pb.Response{}, err
	}

	result, err := db.CreateVolumeSnapshot(snp)
	if err != nil {
		log.Println("Error occured in dock module when create volume snapshot in db:", err)
		return &pb.Response{}, err
	}

	snpBody, _ := json.Marshal(result)
	return &pb.Response{
		Status:  "Success",
		Message: string(snpBody),
	}, nil
}

func GetVolumeSnapshot(vr *pb.VolumeRequest) (*pb.Response, error) {
	var dck = &api.Dock{}
	if err := json.Unmarshal([]byte(vr.GetDockInfo()), dck); err != nil {
		log.Println("[Error] When parsing dock info:", err)
	}

	result, err := dock.GetSnapshot(dck.DriverName, vr.GetSnapshotId())
	if err != nil {
		log.Println("Error occured in dock module when get snapshot:", err)
		return &pb.Response{}, err
	}

	snpBody, _ := json.Marshal(result)
	return &pb.Response{
		Status:  "Success",
		Message: string(snpBody),
	}, nil
}

func DeleteVolumeSnapshot(vr *pb.VolumeRequest) (*pb.Response, error) {
	var dck = &api.Dock{}
	if err := json.Unmarshal([]byte(vr.GetDockInfo()), dck); err != nil {
		log.Println("[Error] When parsing dock info:", err)
	}

	if err := dock.DeleteSnapshot(dck.DriverName, vr.GetSnapshotId()); err != nil {
		log.Println("Error occured in dock module when delete snapshot:", err)
		return &pb.Response{}, err
	}

	if err := db.DeleteVolumeSnapshot(vr.GetSnapshotId()); err != nil {
		log.Println("Error occured in dock module when delete volume snapshot in db:", err)
		return &pb.Response{}, err
	}

	return &pb.Response{
		Status:  "Success",
		Message: "Delete snapshot success",
	}, nil
}
