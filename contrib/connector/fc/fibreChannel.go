// Copyright 2019 The OpenSDS Authors.
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

	"github.com/mitchellh/mapstructure"
	"github.com/opensds/opensds/contrib/connector"
)

var (
	tries = 3
)

// ConnectorInfo define
type ConnectorInfo struct {
	AccessMode string   `mapstructure:"accessMode"`
	AuthUser   string   `mapstructure:"authUserName"`
	AuthPass   string   `mapstructure:"authPassword"`
	AuthMethod string   `mapstructure:"authMethod"`
	TgtDisco   bool     `mapstructure:"targetDiscovered"`
	TargetWWN  []string `mapstructure:"targetWWN"`
	VolumeID   string   `mapstructure:"volumeId"`
	TgtLun     string   `mapstructure:"targetLun"`
	Encrypted  bool     `mapstructure:"encrypted"`
}

// ParseIscsiConnectInfo decode
func parseFCConnectInfo(connectInfo map[string]interface{}) (*ConnectorInfo, error) {
	var con ConnectorInfo
	mapstructure.Decode(connectInfo, &con)

	if len(con.TargetWWN) == 0 || con.TgtLun == "0" {
		return nil, errors.New("fibrechannel connection data invalid.")
	}

	return &con, nil
}

func connectVolume(connMap map[string]interface{}) (map[string]string, error) {
	conn, err := parseFCConnectInfo(connMap)
	if err != nil {
		return nil, err
	}
	hbas, err := getFChbasInfo()
	if err != nil {
		return nil, err
	}
	volPaths := getVolumePaths(conn, hbas)
	if len(volPaths) == 0 {
		errMsg := fmt.Sprintf("No FC devices found.\n")
		log.Printf(errMsg)
		return nil, errors.New(errMsg)
	}

	devicePath, deviceName := volPathDiscovery(volPaths, tries, conn.TargetWWN, hbas)
	if devicePath != "" && deviceName != "" {
		log.Printf("Found Fibre Channel volume name, devicePath is %s, deviceName is %s\n", devicePath, deviceName)
	}

	deviceWWN, err := getSCSIWWN(devicePath)
	if err != nil {
		return nil, err
	}

	return map[string]string{"scsi_wwn": deviceWWN, "path": devicePath}, nil
}

func getVolumePaths(conn *ConnectorInfo, hbas []map[string]string) []string {
	wwnports := conn.TargetWWN
	devices := getDevices(hbas, wwnports)
	lun := conn.TgtLun
	hostPaths := getHostDevices(devices, lun)
	return hostPaths
}

func volPathDiscovery(volPaths []string, tries int, tgtWWN []string, hbas []map[string]string) (string, string) {
	for i := 0; i < tries; i++ {
		for _, path := range volPaths {
			if pathExists(path) {
				deviceName := getContentfromSymboliclink(path)
				return path, deviceName
			}
			rescanHosts(tgtWWN, hbas)
		}

		time.Sleep(2 * time.Second)
	}
	return "", ""
}

func getHostDevices(devices []map[string]string, lun string) []string {
	var hostDevices []string
	for _, device := range devices {
		var hostDevice string
		for pciNum, tgtWWN := range device {
			hostDevice = fmt.Sprintf("/dev/disk/by-path/pci-%s-fc-%s-lun-%s", pciNum, tgtWWN, processLunID(lun))
		}
		hostDevices = append(hostDevices, hostDevice)
	}
	return hostDevices
}

func disconnectVolume(connMap map[string]interface{}) error {
	conn, err := parseFCConnectInfo(connMap)
	if err != nil {
		return err
	}
	volPaths, err := getVolumePathsForDetach(conn)
	if err != nil {
		return err
	}

	var devices []map[string]string
	for _, path := range volPaths {
		realPath := getContentfromSymboliclink(path)
		deviceInfo, _ := getDeviceInfo(realPath)
		devices = append(devices, deviceInfo)
	}

	return removeDevices(devices)
}

func removeDevices(devices []map[string]string) error {
	for _, device := range devices {
		path := fmt.Sprintf("/sys/block/%s/device/delete", strings.Replace(device["device"], "/dev/", "", -1))
		if pathExists(path) {
			if err := flushDeviceIO(device["device"]); err != nil {
				return err
			}

			if err := removeSCSIDevice(path); err != nil {
				return err
			}
		}
	}
	return nil
}

func getPciNum(hba map[string]string) string {
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

func getVolumePathsForDetach(conn *ConnectorInfo) ([]string, error) {
	var volPaths []string
	hbas, err := getFChbasInfo()
	if err != nil {
		return nil, err
	}

	devicePaths := getVolumePaths(conn, hbas)
	for _, path := range devicePaths {
		if pathExists(path) {
			volPaths = append(volPaths, path)
		}
	}
	return volPaths, nil
}

func getDevices(hbas []map[string]string, wwnports []string) []map[string]string {
	var device []map[string]string
	for _, hba := range hbas {
		pciNum := getPciNum(hba)
		if pciNum != "" {
			for _, wwn := range wwnports {
				tgtWWN := map[string]string{pciNum: "0x" + wwn}
				device = append(device, tgtWWN)
			}
		}
	}
	return device
}

func processLunID(lunID string) string {
	lunIDInt, _ := strconv.Atoi(lunID)
	if lunIDInt < 256 {
		return lunID
	}
	return fmt.Sprintf("0x%04x%04x00000000", lunIDInt&0xffff, lunIDInt>>16&0xffff)
}

func getFChbasInfo() ([]map[string]string, error) {
	// Get Fibre Channel WWNs and device paths from the system.
	hbas, err := getFChbas()
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

func getInitiatorInfo() (string, error) {
	hbas, err := getFChbasInfo()
	if err != nil {
		return "", err
	}

	var initiatorInfo []string

	for _, hba := range hbas {
		if v, ok := hba[connector.PortName]; ok {
			initiatorInfo = append(initiatorInfo, "port_name:"+v)
		}
		if v, ok := hba[connector.NodeName]; ok {
			initiatorInfo = append(initiatorInfo, "node_name:"+v)
		}
	}

	return strings.Join(initiatorInfo, ","), nil
}
