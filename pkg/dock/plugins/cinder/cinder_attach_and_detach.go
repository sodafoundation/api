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
	"log"
	"os"
	"os/exec"
	"strings"
	// "time"

	"openstack/golang-client/volume"
)

const (
	CEPH_POOL_NAME     string = "volumes"
	CEPH_LINK_PREFIX   string = "volumes/volume-"
	CINDER_LINK_PREFIX string = "/dev/mapper/cinder--volumes-volume--"
	DEVICE_PREFIX      string = "/dev/"
)

func AttachVolumeToHost(plugin *CinderPlugin, volID, volType string) (string, error) {
	var devPath string
	switch volType {
	case "lvm":
		dev, err := getLvmDevicePath(volID)
		if err != nil {
			return "", err
		}

		devPath = dev
	case "ceph":
		dev, err := getCephDevicePath(volID)
		if err != nil {
			return "", err
		}

		devPath = dev
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

	vol, err := volumeService.Show(volID)
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
	body := &volume.AttachBody{
		VolumeBody: volume.RequestBody{
			HostName:   host,
			Mountpoint: device,
		},
	}

	err = volumeService.Attach(volID, body)
	if err != nil {
		log.Println("Cannot attach volume:", err)
		return err
	}

	return nil
}

func DetachVolumeFromHost(plugin *CinderPlugin, device string) (string, error) {
	var volID string

	if strings.HasPrefix(device, CINDER_LINK_PREFIX) {
		volumeId := device[len(CINDER_LINK_PREFIX):len(device)]
		volID = strings.Replace(volumeId, "--", "-", 4)
	} else {
		if strings.HasPrefix(device, "/dev/rbd") {
			image, err := parseCephDevicePath(device)
			if err != nil {
				return "", err
			}

			volID = image[7:len(image)]

			unmapCmd := exec.Command("rbd", "unmap", device)
			_, err = unmapCmd.CombinedOutput()
			if err != nil {
				return "", err
			}
		} else {
			err := errors.New("Unexpect device prefix: " + device)
			return "", err
		}
	}

	err := sendDetachRequest(plugin, volID)
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

	vol, err := volumeService.Show(volID)
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
	body := &volume.DetachBody{
		VolumeBody: volume.RequestBody{
			AttachmentID: vol.Attachments[0]["attachment_id"],
		},
	}

	err = volumeService.Detach(volID, body)
	if err != nil {
		log.Println("Cannot detach volume:", err)
		return err
	}

	return nil
}

func getLvmDevicePath(volId string) (string, error) {
	link := CINDER_LINK_PREFIX + strings.Replace(volId, "-", "--", 4)

	path, err := os.Readlink(link)
	if err != nil {
		err = errors.New("Can't find device path!")
		return "", err
	}

	slice := strings.Split(path, "/")
	device := "/dev/" + slice[1]
	return device, nil
}

func getCephDevicePath(volId string) (string, error) {
	imagesCmd := exec.Command("rbd", "showmapped")
	images, err := imagesCmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	if strings.Contains(string(images), volId) {
		err := errors.New("This volume have been attached!")
		return "", err
	}

	pathCmd := exec.Command("rbd", "map", CEPH_LINK_PREFIX+volId)
	device, err := pathCmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	devSlice := strings.Split(string(device), "\n")
	return devSlice[0], nil
}

func parseCephDevicePath(device string) (string, error) {
	linksCmd := exec.Command("rbd", "showmapped")
	links, err := linksCmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	linkSlice := strings.Split(string(links), "\n")
	for _, i := range linkSlice {
		if strings.Contains(i, device) {
			imageSlice := strings.Fields(i)
			return imageSlice[2], nil
		}
	}

	err = errors.New("Can't parse device path!")
	return "", err
}
