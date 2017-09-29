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
This module implements the policy-based scheduling by parsing storage
profiles configured by admin.
*/

package controller

import (
	"errors"

	log "github.com/golang/glog"

	"github.com/opensds/opensds/pkg/db"
	pb "github.com/opensds/opensds/pkg/dock/proto"
	"github.com/opensds/opensds/pkg/model"
)

func SearchProfile(prfId string, dbCli db.Client) (*model.ProfileSpec, error) {
	// If a user doesn't specify profile id, then a default profile will be
	// automatically assigned.
	if prfId == "" {
		prfs, err := dbCli.ListProfiles()
		if err != nil {
			log.Error("When list profiles:", err)
			return nil, err
		}

		for _, prf := range *prfs {
			if prf.GetName() == "default" {
				return &prf, nil
			}
		}

		return nil, errors.New("Can not find default profile in db!")
	}

	return dbCli.GetProfile(prfId)
}

func GenerateCreateVolumeOpts(vol *model.VolumeSpec) (*pb.CreateVolumeOpts, error) {
	if vol.GetName() == "" {
		return nil, errors.New("Name field needed for creating a volume!")
	}
	if vol.GetSize() == int64(0) {
		return nil, errors.New("Size field needed for creating a volume!")
	}

	return &pb.CreateVolumeOpts{
		Id:               vol.GetId(),
		Name:             vol.GetName(),
		Description:      vol.GetDescription(),
		Size:             vol.GetSize(),
		AvailabilityZone: vol.GetAvailabilityZone(),
		ProfileId:        vol.GetProfileId(),
	}, nil
}

func GenerateDeleteVolumeOpts(vol *model.VolumeSpec) (*pb.DeleteVolumeOpts, error) {
	if vol.GetId() == "" {
		return nil, errors.New("Id field needed for deleting a volume!")
	}

	return &pb.DeleteVolumeOpts{
		Id: vol.GetId(),
	}, nil
}

func GenerateCreateAttachmentOpts(atc *model.VolumeAttachmentSpec) (*pb.CreateAttachmentOpts, error) {
	if atc.GetVolumeId() == "" {
		return nil, errors.New("Volume id field needed for creating a volume attachment!")
	}

	return &pb.CreateAttachmentOpts{
		Id:       atc.GetId(),
		VolumeId: atc.GetVolumeId(),
	}, nil
}

func GenerateCreateVolumeSnapshotOpts(snp *model.VolumeSnapshotSpec) (*pb.CreateVolumeSnapshotOpts, error) {
	if snp.GetName() == "" {
		return nil, errors.New("Name field needed for creating a volume snapshot!")
	}
	if snp.GetVolumeId() == "" {
		return nil, errors.New("Volume id field needed for creating a volume snapshot!")
	}

	return &pb.CreateVolumeSnapshotOpts{
		Id:          snp.GetId(),
		Name:        snp.GetName(),
		Description: snp.GetDescription(),
		Size:        snp.GetSize(),
		VolumeId:    snp.GetVolumeId(),
	}, nil
}

func GenerateDeleteVolumeSnapshotOpts(snp *model.VolumeSnapshotSpec) (*pb.DeleteVolumeSnapshotOpts, error) {
	if snp.GetId() == "" {
		return nil, errors.New("Id field needed for deleting a volume snapshot!")
	}

	return &pb.DeleteVolumeSnapshotOpts{
		Id: snp.GetId(),
	}, nil
}
