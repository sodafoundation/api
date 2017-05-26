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
This module implements the database operation interface of data structure
defined in api module.

*/

package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/coreos/etcd/clientv3"
	api "github.com/opensds/opensds/pkg/api/v1"
	"github.com/opensds/opensds/pkg/db"
)

var (
	URL_PREFIX = "/api/v1/"
	dbCli      = &db.DbClient{}
)

func init() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		cli.Close()
		fmt.Errorf(err.Error())
	}

	dbCli.Client = cli
}

func CreateVolume(vol *api.VolumeResponse) (*api.VolumeResponse, error) {
	volBody, err := json.Marshal(vol)
	if err != nil {
		return &api.VolumeResponse{}, err
	}

	dbReq := &db.DbRequest{
		Url:     URL_PREFIX + "volumes/" + vol.Id,
		Content: string(volBody),
	}
	dbRes := dbCli.Create(dbReq)
	if dbRes.Status != "Success" {
		log.Println("[Error] When create volume in db:", dbRes.Error)
		return &api.VolumeResponse{}, errors.New(dbRes.Error)
	}

	return vol, nil
}

func GetVolume(volID string) (*api.VolumeResponse, error) {
	dbReq := &db.DbRequest{
		Url: URL_PREFIX + "volumes/" + volID,
	}
	dbRes := dbCli.Get(dbReq)
	if dbRes.Status != "Success" {
		log.Println("[Error] When get volume in db:", dbRes.Error)
		return &api.VolumeResponse{}, errors.New(dbRes.Error)
	}

	var vol = &api.VolumeResponse{}
	if err := json.Unmarshal([]byte(dbRes.Message[0]), vol); err != nil {
		log.Println("[Error] When parsing volume in db:", dbRes.Error)
		return &api.VolumeResponse{}, errors.New(dbRes.Error)
	}
	return vol, nil
}

func ListVolumes() (*[]api.VolumeResponse, error) {
	dbReq := &db.DbRequest{
		Url: URL_PREFIX + "volumes",
	}
	dbRes := dbCli.List(dbReq)
	if dbRes.Status != "Success" {
		log.Println("[Error] When list volumes in db:", dbRes.Error)
		return &[]api.VolumeResponse{}, errors.New(dbRes.Error)
	}

	var vols = []api.VolumeResponse{}
	if len(dbRes.Message) == 0 {
		return &vols, nil
	}
	for _, msg := range dbRes.Message {
		var vol = api.VolumeResponse{}
		if err := json.Unmarshal([]byte(msg), &vol); err != nil {
			log.Println("[Error] When parsing volume in db:", dbRes.Error)
			return &[]api.VolumeResponse{}, errors.New(dbRes.Error)
		}
		vols = append(vols, vol)
	}
	return &vols, nil
}

func DeleteVolume(volID string) error {
	dbReq := &db.DbRequest{
		Url: URL_PREFIX + "volumes/" + volID,
	}
	dbRes := dbCli.Delete(dbReq)
	if dbRes.Status != "Success" {
		log.Println("[Error] When delete volume in db:", dbRes.Error)
		return errors.New(dbRes.Error)
	}
	return nil
}

func CreateVolumeAttachment(volID string, atc *api.VolumeAttachment) (*api.VolumeAttachment, error) {
	atcBody, err := json.Marshal(atc)
	if err != nil {
		return &api.VolumeAttachment{}, err
	}

	dbReq := &db.DbRequest{
		Url:     URL_PREFIX + "volume/" + volID + "/attachments/" + atc.Id,
		Content: string(atcBody),
	}
	dbRes := dbCli.Create(dbReq)
	if dbRes.Status != "Success" {
		log.Println("[Error] When create volume attachment in db:", dbRes.Error)
		return &api.VolumeAttachment{}, errors.New(dbRes.Error)
	}

	return atc, nil
}

func GetVolumeAttachment(volID, attachmentID string) (*api.VolumeAttachment, error) {
	dbReq := &db.DbRequest{
		Url: URL_PREFIX + "volume/" + volID + "/attachments/" + attachmentID,
	}
	dbRes := dbCli.Get(dbReq)
	if dbRes.Status != "Success" {
		log.Println("[Error] When get volume attachment in db:", dbRes.Error)
		return &api.VolumeAttachment{}, errors.New(dbRes.Error)
	}

	var atc = &api.VolumeAttachment{}
	if err := json.Unmarshal([]byte(dbRes.Message[0]), atc); err != nil {
		log.Println("[Error] When parsing volume attachment in db:", dbRes.Error)
		return &api.VolumeAttachment{}, errors.New(dbRes.Error)
	}
	return atc, nil
}

func ListVolumeAttachments(volID string) (*[]api.VolumeAttachment, error) {
	dbReq := &db.DbRequest{
		Url: URL_PREFIX + "volume/" + volID + "/attachments",
	}
	dbRes := dbCli.List(dbReq)
	if dbRes.Status != "Success" {
		log.Println("[Error] When list volume attachments in db:", dbRes.Error)
		return &[]api.VolumeAttachment{}, errors.New(dbRes.Error)
	}

	var atcs = []api.VolumeAttachment{}
	if len(dbRes.Message) == 0 {
		return &atcs, nil
	}
	for _, msg := range dbRes.Message {
		var atc = api.VolumeAttachment{}
		if err := json.Unmarshal([]byte(msg), &atc); err != nil {
			log.Println("[Error] When parsing volume attachment in db:", dbRes.Error)
			return &[]api.VolumeAttachment{}, errors.New(dbRes.Error)
		}
		atcs = append(atcs, atc)
	}
	return &atcs, nil
}

func UpdateVolumeAttachment(volID, attachmentID, mountpoint string, hostInfo *api.HostInfo) (*api.VolumeAttachment, error) {
	atc, err := GetVolumeAttachment(volID, attachmentID)
	if err != nil {
		return &api.VolumeAttachment{}, err
	}

	atc.HostInfo = *hostInfo
	atc.Mountpoint = mountpoint
	atcBody, err := json.Marshal(atc)
	if err != nil {
		return &api.VolumeAttachment{}, err
	}

	dbReq := &db.DbRequest{
		Url:        URL_PREFIX + "volume/" + volID + "/attachments/" + attachmentID,
		NewContent: string(atcBody),
	}
	dbRes := dbCli.Update(dbReq)
	if dbRes.Status != "Success" {
		log.Println("[Error] When update volume attachment in db:", dbRes.Error)
		return &api.VolumeAttachment{}, errors.New(dbRes.Error)
	}
	return atc, nil
}

func DeleteVolumeAttachment(volID, attachmentID string) error {
	dbReq := &db.DbRequest{
		Url: URL_PREFIX + "volume/" + volID + "/attachments/" + attachmentID,
	}
	dbRes := dbCli.Delete(dbReq)
	if dbRes.Status != "Success" {
		log.Println("[Error] When delete volume attachment in db:", dbRes.Error)
		return errors.New(dbRes.Error)
	}
	return nil
}

func CreateVolumeSnapshot(vs *api.VolumeSnapshot) (*api.VolumeSnapshot, error) {
	vsBody, err := json.Marshal(vs)
	if err != nil {
		return &api.VolumeSnapshot{}, err
	}

	dbReq := &db.DbRequest{
		Url:     URL_PREFIX + "volume/snapshots/" + vs.Id,
		Content: string(vsBody),
	}
	dbRes := dbCli.Create(dbReq)
	if dbRes.Status != "Success" {
		log.Println("[Error] When create volume snapshot in db:", dbRes.Error)
		return &api.VolumeSnapshot{}, errors.New(dbRes.Error)
	}

	return vs, nil
}

func GetVolumeSnapshot(snapshotID string) (*api.VolumeSnapshot, error) {
	dbReq := &db.DbRequest{
		Url: URL_PREFIX + "volume/snapshots/" + snapshotID,
	}
	dbRes := dbCli.Get(dbReq)
	if dbRes.Status != "Success" {
		log.Println("[Error] When get volume attachment in db:", dbRes.Error)
		return &api.VolumeSnapshot{}, errors.New(dbRes.Error)
	}

	var vs = &api.VolumeSnapshot{}
	if err := json.Unmarshal([]byte(dbRes.Message[0]), vs); err != nil {
		log.Println("[Error] When parsing volume snapshot in db:", dbRes.Error)
		return &api.VolumeSnapshot{}, errors.New(dbRes.Error)
	}
	return vs, nil
}

func ListVolumeSnapshots() (*[]api.VolumeSnapshot, error) {
	dbReq := &db.DbRequest{
		Url: URL_PREFIX + "volume/snapshots",
	}
	dbRes := dbCli.List(dbReq)
	if dbRes.Status != "Success" {
		log.Println("[Error] When list volume snapshots in db:", dbRes.Error)
		return &[]api.VolumeSnapshot{}, errors.New(dbRes.Error)
	}

	var vss = []api.VolumeSnapshot{}
	if len(dbRes.Message) == 0 {
		return &vss, nil
	}
	for _, msg := range dbRes.Message {
		var vs = api.VolumeSnapshot{}
		if err := json.Unmarshal([]byte(msg), &vs); err != nil {
			log.Println("[Error] When parsing volume snapshot in db:", dbRes.Error)
			return &[]api.VolumeSnapshot{}, errors.New(dbRes.Error)
		}
		vss = append(vss, vs)
	}
	return &vss, nil
}

func DeleteVolumeSnapshot(snapshotID string) error {
	dbReq := &db.DbRequest{
		Url: URL_PREFIX + "volume/snapshots/" + snapshotID,
	}
	dbRes := dbCli.Delete(dbReq)
	if dbRes.Status != "Success" {
		log.Println("[Error] When delete volume snapshot in db:", dbRes.Error)
		return errors.New(dbRes.Error)
	}
	return nil
}
