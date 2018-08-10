// Copyright (c) 2018 Huawei Technologies Co., Ltd. All Rights Reserved.
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

package fc

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

var (
	tries = 3
)

type fibreChannel struct {
	helper *linuxfc
}

func (f *fibreChannel) connectVolume(conn map[string]interface{}) (map[string]string, error) {
	hbas, err := f.getFChbasInfo()
	if err != nil {
		return nil, err
	}
	volPaths := f.getVolumePaths(conn, hbas)
	if len(volPaths) == 0 {
		errMsg := fmt.Sprintf("No FC devices found.")
		log.Println(errMsg)
		return nil, errors.New(errMsg)
	}

	devicePath, deviceName := f.volPathDiscovery(volPaths, tries, conn["target_wwn"].([]string), hbas)
	if devicePath != "" && deviceName != "" {
		log.Printf("Found Fibre Channel volume name, devicePath is %s, deviceName is %s", devicePath, deviceName)
	}

	deviceWWN, err := f.helper.getSCSIWWN(devicePath)
	if err != nil {
		return nil, err
	}

	return map[string]string{"scsi_wwn": deviceWWN, "path": devicePath}, nil
}

func (f *fibreChannel) getVolumePaths(conn map[string]interface{}, hbas []map[string]string) []string {
	wwnports := conn["target_wwn"].([]string)
	devices := f.getDevices(hbas, wwnports)
	lun := conn["target_lun"].(string)
	hostPaths := f.getHostDevices(devices, lun)
	return hostPaths
}

func (f *fibreChannel) volPathDiscovery(volPaths []string, tries int, tgtWWN []string, hbas []map[string]string) (string, string) {
	for i := 0; i < tries; i++ {
		for _, path := range volPaths {
			if f.helper.pathExists(path) {
				deviceName := f.helper.getContentfromSymboliclink(path)
				return path, deviceName
			} else {
				f.helper.rescanHosts(tgtWWN, hbas)
			}
		}

		time.Sleep(2 * time.Second)
	}
	return "", ""
}

func (f *fibreChannel) getHostDevices(devices []map[string]string, lun string) []string {
	var hostDevices []string
	for _, device := range devices {
		var hostDevice string
		for pciNum, tgtWWN := range device {
			hostDevice = fmt.Sprintf("/dev/disk/by-path/pci-%s-fc-%s-lun-%s", pciNum, tgtWWN, f.processLunId(lun))
		}
		hostDevices = append(hostDevices, hostDevice)
	}
	return hostDevices
}

func (f *fibreChannel) disconnectVolume(conn map[string]interface{}) error {
	volPaths, err := f.getVolumePathsForDetach(conn)
	if err != nil {
		return err
	}

	var devices []map[string]string
	for _, path := range volPaths {
		realPath := f.helper.getContentfromSymboliclink(path)
		deviceInfo, _ := f.helper.getDeviceInfo(realPath)
		devices = append(devices, deviceInfo)
	}

	return f.removeDevices(devices)
}

func (f *fibreChannel) removeDevices(devices []map[string]string) error {
	for _, device := range devices {
		path := fmt.Sprintf("/sys/block/%s/device/delete", strings.Replace(device["device"], "/dev/", "", -1))
		if f.helper.pathExists(path) {
			if err := f.helper.flushDeviceIO(device["device"]); err != nil {
				return err
			}

			if err := f.helper.removeSCSIDevice(path); err != nil {
				return err
			}
		}
	}
	return nil
}

func (f *fibreChannel) getPciNum(hba map[string]string) string {
	for k, v := range hba {
		if k == "device_path" {
			path := strings.Split(v, "/")
			for idx, u := range path {
				if strings.Contains(u, "net") || strings.Contains(u, "host") {
					return path[idx-1]
				}
			}
		}
	}
	return ""
}

func (f *fibreChannel) getVolumePathsForDetach(conn map[string]interface{}) ([]string, error) {
	var volPaths []string
	hbas, err := f.getFChbasInfo()
	if err != nil {
		return nil, err
	}

	devicePaths := f.getVolumePaths(conn, hbas)
	for _, path := range devicePaths {
		if f.helper.pathExists(path) {
			volPaths = append(volPaths, path)
		}
	}
	return volPaths, nil
}

func (f *fibreChannel) getDevices(hbas []map[string]string, wwnports []string) []map[string]string {
	var device []map[string]string
	for _, hba := range hbas {
		pciNum := f.getPciNum(hba)
		if pciNum != "" {
			for _, wwn := range wwnports {
				tgtWWN := map[string]string{pciNum: "0x" + wwn}
				device = append(device, tgtWWN)
			}
		}
	}
	return device
}

func (f *fibreChannel) processLunId(lunId string) string {
	lunIdInt, _ := strconv.Atoi(lunId)
	if lunIdInt < 256 {
		return lunId
	}
	return fmt.Sprintf("0x%04x%04x00000000", lunIdInt&0xffff, lunIdInt>>16&0xffff)
}

func (f *fibreChannel) getFChbasInfo() ([]map[string]string, error) {
	// Get Fibre Channel WWNs and device paths from the system.
	hbas, err := f.helper.getFChbas()
	if err != nil {
		return nil, err
	}
	var hbasInfos []map[string]string
	for _, hba := range hbas {
		wwpn := strings.Replace(hba["port_name"], "0x", "", -1)
		wwnn := strings.Replace(hba["node_name"], "0x", "", -1)
		devicePath := hba["ClassDevicepath"]
		device := hba["ClassDevice"]

		hbasInfo := map[string]string{"port_name": wwpn, "node_name": wwnn, "host_device": device, "device_path": devicePath}

		hbasInfos = append(hbasInfos, hbasInfo)
	}

	return hbasInfos, nil
}
