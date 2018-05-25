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
	"os"
	"os/exec"

	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func connectVolume(conn map[string]interface{}) (map[string]string, error) {
	hbas, err := getFChbasInfo()
	if err != nil {
		return nil, err
	}
	volPaths := getVolumePaths(conn, hbas)
	if len(volPaths) == 0 {
		errMsg := fmt.Sprintf("No fc devices found.")
		log.Println(errMsg)
		return nil, errors.New(errMsg)
	}

	var tries = 3
	devicePath, deviceName := volPathDiscovery(volPaths, tries, conn["target_wwn"].(string), hbas)
	if devicePath != "" && deviceName != "" {
		log.Printf("Found fibre channel volume name, devicePath is %s, deviceName is %s", devicePath, deviceName)
	}

	deviceWWN, err := getSCSIWWN(devicePath)
	if err != nil {
		return nil, err
	}

	return map[string]string{"scsi_wwn": deviceWWN, "path": devicePath}, nil
}

func getSCSIWWN(devicePath string) (string, error) {
	out, err := exec.Command("/lib/udev/scsi_id", "--page", "0x83", "--whitelisted", devicePath).CombinedOutput()
	outString := string(out)
	if err != nil {
		errMsg := fmt.Sprintf("Error occurred when get device wwn:", outString, err)
		log.Println(errMsg)
		return "", errors.New(errMsg)
	}
	return strings.TrimSpace(outString), nil
}

func volPathDiscovery(volPaths []string, tries int, tgtWWN string, hbas []map[string]string) (string, string) {
	for i := 0; i < tries; i++ {
		for _, path := range volPaths {
			if pathExists(path) {
				deviceName, _ := filepath.Abs(path)
				return path, deviceName
			} else {
				rescanHosts(tgtWWN, hbas)
			}
		}
		time.Sleep(2 * time.Second)
	}
	return "", ""
}

func rescanHosts(tgtWWN string, hbas []map[string]string) error {
	for _, hba := range hbas {
		chlAndSCSitgt, err := getHbaChanAndSCSITgt(hba)
		if err != nil {
			return nil
		}
		if len(chlAndSCSitgt) == 0 {
			chlAndSCSitgt = []string{"-", "-"}
		}

		//		for hbaChl,targ
		//		cmd:="tee -a /sys/class/scsi_host/%s/scan",hba["host_device"],
	}
	return nil
}

func getHbaChanAndSCSITgt(hba map[string]string) ([]string, error) {
	// Get HBA channel and scsi target for an HBA
	hostDevice := hba["host_device"]
	if hostDevice != "" && len(hostDevice) > 4 {
		hostDevice = hostDevice[4:]
	}

	path := fmt.Sprintf("/sys/class/fc_transport/target%s:", hostDevice)
	cmd := fmt.Sprintf("grep %s %s*/node_name", hba["node_name"], path)

	out, err := exec.Command(cmd).CombinedOutput()
	outString := string(out)

	if err != nil {
		errMsg := fmt.Sprintf("Error occurred when get HBA Channel and SCSI target:", outString, err)
		log.Println(errMsg)
		return nil, errors.New(errMsg)
	}

	for _, line := range strings.Split(outString, "\n") {
		if strings.Contains(line, path) {
			out1 := strings.Split(line, "/")[4]
			return strings.Split(out1, ":")[1:], nil
		}
	}

	return []string{}, nil
}

func getVolumePaths(conn map[string]interface{}, hbas []map[string]string) []string {
	wwnports := conn["target_wwn"].([]string)
	devices := getDevices(hbas, wwnports)
	lun := conn["target_lun"].(string)
	hostPaths := getHostDevices(devices, lun)
	return hostPaths
}

func getHostDevices(devices []map[string]string, lun string) []string {
	var hostDevices []string
	for _, device := range devices {
		var hostDevice string
		for pciNum, tgtWWN := range device {
			hostDevice = fmt.Sprintf("/dev/disk/by-path/pci-%s-fc-%s-lun-%s", pciNum, tgtWWN, processLunId(lun))
		}
		hostDevices = append(hostDevices, hostDevice)
	}
	return hostDevices
}

func processLunId(lunId string) string {
	lunIdInt, _ := strconv.Atoi(lunId)
	if lunIdInt < 256 {
		return lunId
	} else {
		return fmt.Sprintf("0x%04x%04x00000000", lunIdInt&0xffff, lunIdInt>>16&0xffff)
	}

	return ""
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

func getFChbasInfo() ([]map[string]string, error) {
	// Get Fibre Channel WWNs and device paths from the system.
	hbas, err := getFChbas()
	if err != nil {
		return nil, err
	}
	var hbasInfos []map[string]string
	for _, hba := range hbas {
		wwpn := strings.Replace(hba["prot_name"], "0x", "", -1)
		wwnn := strings.Replace(hba["node_name"], "0x", "", -1)
		devicePath := hba["ClassDevicepath"]
		device := hba["ClassDevice"]

		hbasInfo := map[string]string{"port_name": wwpn, "node_name": wwnn, "host_device": device, "device_path": devicePath}

		hbasInfos = append(hbasInfos, hbasInfo)
	}
	return hbasInfos, nil
}

func getFChbas() ([]map[string]string, error) {
	if !fcSupport() {
		errMsg := fmt.Sprintf("No Fibre Channel support detected.")
		log.Printf(errMsg)
		return nil, errors.New(errMsg)
	}

	out, err := exec.Command("systool", "-c", "fc_host", "-v").CombinedOutput()
	outString := string(out)

	if err != nil {
		errMsg := fmt.Sprintf("Error occurred when get FC hbas info: systool is not installed", outString, err)
		log.Println(errMsg)
		return nil, errors.New(errMsg)
	}

	if outString == "" {
		errMsg := fmt.Sprintf("No Fibre Channel support detected.")
		log.Printf(errMsg)
		return nil, errors.New(errMsg)
	}

	lines := strings.Split(outString, "\n")
	lines = lines[2:]
	hba := make(map[string]string)
	hbas := []map[string]string{}
	lastline := ""

	for _, line := range lines {
		line = strings.TrimSpace(line)
		// 2 newlines denotes a new hba port
		if line == "" && lastline == "" {
			if len(hbas) > 0 {
				hbas = append(hbas, hba)
				hba = make(map[string]string)
			}
		} else {
			val := strings.Split(line, "=")
			if len(val) == 2 {
				key := strings.Replace(val[0], " ", "", -1)
				val := strings.Replace(val[1], " ", "", -1)
				hba[key] = val
			}
		}
		lastline = line
	}

	return hbas, nil
}

func fcSupport() bool {
	var FC_HOST_SYSFC_PATH = "/sys/class/fc_host"
	return pathExists(FC_HOST_SYSFC_PATH)
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func getVolumePathsForDetach(conn map[string]interface{}) ([]string, error) {
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

func disconnectVolume(conn map[string]interface{}) error {
	volPaths, err := getVolumePathsForDetach(conn)
	if err != nil {
		return err
	}

	var devices []map[string]string
	for _, path := range volPaths {
		realPath, _ := filepath.Abs(path)
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

func removeSCSIDevice(path string) error {
	out, err := exec.Command("echo", "1", "|", "tee -a", path).CombinedOutput()
	outString := string(out)
	if err != nil {
		errMsg := fmt.Sprintf("Error occurred when get device info when detach volume:", outString, err)
		log.Println(errMsg)
		return errors.New(errMsg)
	}
	return nil
}

func flushDeviceIO(device string) error {
	out, err := exec.Command("blockdev ", "--flushbufs", device).CombinedOutput()
	outString := string(out)
	if err != nil {
		errMsg := fmt.Sprintf("Error occurred when get device info when detach volume:", outString, err)
		log.Println(errMsg)
		return errors.New(errMsg)
	}
	return nil
}

func getDeviceInfo(devicePath string) (map[string]string, error) {
	out, err := exec.Command("sg_scan", devicePath).CombinedOutput()
	outString := string(out)
	if err != nil {
		errMsg := fmt.Sprintf("Error occurred when get device info when detach volume:", outString, err)
		log.Println(errMsg)
		return nil, errors.New(errMsg)
	}

	devInfo := make(map[string]string)
	line := strings.TrimSpace(outString)
	line = strings.Replace(line, devicePath+": ", "", -1)
	info := strings.Split(line, " ")

	for _, v := range info {
		if strings.Contains(v, "=") {
			pair := strings.Split(v, "=")
			devInfo[pair[0]] = pair[1]
		}
		if strings.Contains(v, "scsi") {
			devInfo["host"] = strings.Replace(v, "scsi", "", -1)
		}
	}

	return devInfo, nil
}
