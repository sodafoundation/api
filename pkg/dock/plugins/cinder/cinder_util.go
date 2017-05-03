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
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	// "time"

	"openstack/golang-client/volume"

	"github.com/opensds/opensds/pkg/dock/plugins/connector"
)

func AttachVolumeToHost(plugin *CinderPlugin, volID string) (string, error) {
	conn, err := getConnectionInfo(plugin, volID)
	if err != nil {
		return "", err
	}

	log.Printf("Receive connection info: %+v\n", conn)

	devPath, err := conn.ConnectVolume()
	if err != nil {
		return "", err
	}

	host, err := os.Hostname()
	if err != nil {
		return "", err
	}

	if err = sendAttachRequest(plugin, volID, host, devPath); err != nil {
		return "", err
	} else {
		return devPath, nil
	}
}

func sendAttachRequest(plugin *CinderPlugin, volID, host, device string) error {
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
			Mountpoint: device,
		},
	}

	err = volumeService.AttachVolume(volID, body)
	if err != nil {
		log.Println("Cannot attach volume:", err)
		return err
	}

	return nil
}

func DetachVolumeFromHost(plugin *CinderPlugin, device string) (string, error) {
	ind := strings.Index(device, "by-id/")
	if ind <= 0 {
		return "", fmt.Errorf("Detach disk: no volume id in %s", device)
	}

	var volID = device[ind+6 : len(device)]

	conn, err := getConnectionInfo(plugin, volID)
	if err != nil {
		return "", err
	}

	log.Printf("Receive connection info: %+v\n", conn)

	_, err = conn.DisconnectVolume()
	if err != nil {
		return "", err
	}

	err = sendDetachRequest(plugin, volID)
	if err != nil {
		return "", err
	} else {
		return "Detach volume success!", nil
	}
}

func sendDetachRequest(plugin *CinderPlugin, volID string) error {
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

func getConnectionInfo(plugin *CinderPlugin, volID string) (*connector.Connector, error) {
	isMultipath := false
	properties, err := connector.GetConnectorProperties(isMultipath)
	if err != nil {
		return &connector.Connector{}, err
	}

	//Get the certified volume service.
	volumeService, err := plugin.getVolumeService()
	if err != nil {
		log.Println("Cannot access volume service:", err)
		return &connector.Connector{}, err
	}

	body := &volume.InitializeBody{
		Connector: volume.Connector{
			ConnectorProperties: *properties,
		},
	}

	connInfo, err := volumeService.InitializeConnection(volID, body)
	if err != nil {
		log.Println("Cannot initialize volume connection:", err)
		return &connector.Connector{}, err
	}

	conn := &connector.Connector{
		ConnInfo: *connInfo,
	}
	return conn, nil
}
