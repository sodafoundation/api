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
This module implements cinder plugin for OpenSDS. Cinder plugin will pass these
operation requests about volume to OpenStack go-client module.

*/

package cinder

import (
	"crypto/tls"
	"errors"
	"log"
	"net/http"
	// "time"

	"openstack/golang-client/auth"
	"openstack/golang-client/volume"

	"git.openstack.org/openstack/golang-client/openstack"
	api "github.com/opensds/opensds/pkg/api/v1"
)

type CinderPlugin struct {
	Host        string
	Methods     []string
	Username    string
	Password    string
	ProjectId   string
	ProjectName string
}

func (plugin *CinderPlugin) Setup() {

}

func (plugin *CinderPlugin) Unset() {

}

func (plugin *CinderPlugin) CreateVolume(name string, size int32) (*api.VolumeResponse, error) {
	//Get the certified volume service.
	volumeService, err := plugin.getVolumeService()
	if err != nil {
		log.Println("Cannot access volume service:", err)
		return &api.VolumeResponse{}, err
	}

	//Configure create request body, the body is defined in volume package.
	body := &volume.VolumeCreateBody{
		VolumeRequestBody: volume.VolumeRequestBody{
			Name:       name,
			VolumeType: "",
			Size:       size,
		},
	}

	vol, err := volumeService.CreateVolume(body)
	if err != nil {
		log.Println("Cannot create volume:", err)
		return &api.VolumeResponse{}, err
	}

	log.Println("Create volume success, dls =", vol)
	return &api.VolumeResponse{
		Id:               vol.Id,
		Name:             vol.Name,
		Description:      vol.Description,
		Size:             vol.Size,
		Status:           vol.Status,
		AvailabilityZone: vol.Availability_zone,
	}, nil
}

func (plugin *CinderPlugin) GetVolume(volID string) (*api.VolumeResponse, error) {
	volumeService, err := plugin.getVolumeService()
	if err != nil {
		log.Println("Cannot access volume service:", err)
		return &api.VolumeResponse{}, err
	}

	vol, err := volumeService.ShowVolume(volID)
	if err != nil {
		log.Println("Cannot show volume:", err)
		return &api.VolumeResponse{}, err
	}

	log.Println("Get volume success, dls =", vol)
	return &api.VolumeResponse{
		Id:               vol.ID,
		Name:             vol.Name,
		Description:      vol.Description,
		Size:             vol.Size,
		Status:           vol.Status,
		AvailabilityZone: vol.Aavailability_zone,
	}, nil
}

func (plugin *CinderPlugin) DeleteVolume(volID string) error {
	volumeService, err := plugin.getVolumeService()
	if err != nil {
		log.Println("Cannot access volume service:", err)
		return err
	}

	if _, err = volumeService.ShowVolume(volID); err != nil {
		log.Println("Cannot get volume:", err)
		return err
	}

	if err = volumeService.DeleteVolume(volID); err != nil {
		log.Println("Cannot delete volume:", err)
		return err
	}
	return nil
}

func (plugin *CinderPlugin) InitializeConnection(volID string, doLocalAttach, multiPath bool, hostInfo *api.HostInfo) (*api.ConnectionInfo, error) {
	//Get the certified volume service.
	volumeService, err := plugin.getVolumeService()
	if err != nil {
		log.Println("Cannot access volume service:", err)
		return &api.ConnectionInfo{}, err
	}

	//Configure initialize request body, the body is defined in volume package.
	body := &volume.InitializeBody{
		Connector: volume.Connector{
			ConnectorProperties: volume.ConnectorProperties{
				DoLocalAttach: doLocalAttach,
				MultiPath:     multiPath,
				Platform:      hostInfo.Platform,
				OsType:        hostInfo.OsType,
				Ip:            hostInfo.Ip,
				Host:          hostInfo.Host,
				Initiator:     hostInfo.Initiator,
			},
		},
	}

	connInfo, err := volumeService.InitializeConnection(volID, body)
	if err != nil {
		log.Println("Cannot initialize volume connection:", err)
		return &api.ConnectionInfo{}, err
	}
	log.Println("Initialize volume connection success, dls =", connInfo)
	return &api.ConnectionInfo{
		DriverVolumeType: connInfo.DriverVolumeType,
		ConnectionData:   connInfo.ConnectionData,
	}, nil
}

func (plugin *CinderPlugin) AttachVolume(volID, host, mountpoint string) error {
	volumeService, err := plugin.getVolumeService()
	if err != nil {
		log.Println("Cannot access volume service:", err)
		return err
	}

	vol, err := volumeService.ShowVolume(volID)
	if err != nil {
		log.Println("Cannot get volume:", err)
		return err
	}
	if vol.Status != "available" && !vol.Multiattach {
		err = errors.New("The status of volume is not available!")
		log.Println("Cannot attach volume:", err)
		return err
	}

	//Configure attach request body, the body is defined in volume package.
	body := &volume.VolumeAttachBody{
		VolumeRequestBody: volume.VolumeRequestBody{
			HostName:   host,
			Mountpoint: mountpoint,
		},
	}

	err = volumeService.AttachVolume(volID, body)
	if err != nil {
		log.Println("Cannot attach volume:", err)
		return err
	}

	return nil
}

func (plugin *CinderPlugin) DetachVolume(volID string) error {
	volumeService, err := plugin.getVolumeService()
	if err != nil {
		log.Println("Cannot access volume service:", err)
		return err
	}

	vol, err := volumeService.ShowVolume(volID)
	if err != nil {
		log.Println("Cannot get volume:", err)
		return err
	}
	if vol.Status != "in-use" {
		err = errors.New("The status of volume is not in-use!")
		log.Println("Cannot detach volume:", err)
		return err
	}

	//Configure detach request body, the body is defined in volume package.
	body := &volume.VolumeDetachBody{
		VolumeRequestBody: volume.VolumeRequestBody{
			AttachmentID: vol.Attachments[0]["attachment_id"],
		},
	}

	err = volumeService.DetachVolume(volID, body)
	if err != nil {
		log.Println("Cannot detach volume:", err)
		return err
	}

	return nil
}

func (plugin *CinderPlugin) CreateSnapshot(name, volID, description string) (*api.VolumeSnapshot, error) {
	//Get the certified volume service.
	volumeService, err := plugin.getVolumeService()
	if err != nil {
		log.Println("Cannot access volume service:", err)
		return &api.VolumeSnapshot{}, err
	}

	//Configure snapshot request body, the body is defined in volume package.
	body := &volume.SnapshotBody{
		SnapshotRequestBody: volume.SnapshotRequestBody{
			Name:            name,
			VolumeID:        volID,
			Description:     description,
			ForceSnapshoted: true,
		},
	}

	snapshot, err := volumeService.CreateSnapshot(body)
	if err != nil {
		log.Println("Cannot create snapshot:", err)
		return &api.VolumeSnapshot{}, err
	}

	log.Println("Create snapshot success, dls =", snapshot)
	return &api.VolumeSnapshot{
		Id:       snapshot.ID,
		Name:     snapshot.Name,
		Status:   snapshot.Status,
		VolumeId: volID,
		Size:     snapshot.Size,
	}, nil
}

func (plugin *CinderPlugin) GetSnapshot(snapID string) (*api.VolumeSnapshot, error) {
	volumeService, err := plugin.getVolumeService()
	if err != nil {
		log.Println("Cannot access volume service:", err)
		return &api.VolumeSnapshot{}, err
	}

	snapshot, err := volumeService.ShowSnapshot(snapID)
	if err != nil {
		log.Println("Cannot show snapshot:", err)
		return &api.VolumeSnapshot{}, err
	}

	log.Println("Get snapshot success, dls =", snapshot)
	return &api.VolumeSnapshot{
		Id:     snapshot.ID,
		Name:   snapshot.Name,
		Status: snapshot.Status,
		Size:   snapshot.Size,
	}, nil
}

func (plugin *CinderPlugin) DeleteSnapshot(snapID string) error {
	volumeService, err := plugin.getVolumeService()
	if err != nil {
		log.Println("Cannot access volume service:", err)
		return err
	}

	_, err = volumeService.ShowSnapshot(snapID)
	if err != nil {
		log.Println("Cannot get snapshot:", err)
		return err
	}

	err = volumeService.DeleteSnapshot(snapID)
	if err != nil {
		log.Println("Cannot delete snapshot:", err)
		return err
	}
	return nil
}

/*
There is some touble now in getVolumeService(). After setting up OpenSDS

service, this process would dump if any credential works don't work. And

we thought it could be solved by make this function a goroutine.

*/
func (plugin *CinderPlugin) getVolumeService() (volume.Service, error) {
	creds := auth.AuthOpts{
		AuthUrl:     plugin.Host,
		Methods:     plugin.Methods,
		Username:    plugin.Username,
		Password:    plugin.Password,
		ProjectId:   plugin.ProjectId,
		ProjectName: plugin.ProjectName,
	}
	auth, err := auth.DoAuthRequestV3(creds)
	if err != nil {
		log.Println("There was an error authenticating:", err)
	}
	/*
		if !auth.GetExpiration().After(time.Now()) {
			log.Fatalln("There was an error. The auth token has an invalid expiration.")
		}
	*/

	// Find the endpoint for the volume v2 service.
	url, err := auth.GetEndpoint("volumev2", "")
	if url == "" || err != nil {
		log.Println("Volume service url not found during authentication.")
	}

	// Make a new client with these creds, here configure InsecureSkipVerify
	// in tls.Config to skip the certificate verification.
	tls := &tls.Config{
		InsecureSkipVerify: true,
	}

	sess, err := openstack.NewSession(nil, auth, tls)
	if err != nil {
		log.Println("Error creating new Session:", err)
	}

	volumeService := volume.NewService(sess, http.DefaultClient, url)
	return volumeService, nil
}
