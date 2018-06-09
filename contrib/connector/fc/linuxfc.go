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
	"strings"
)

type linuxfc struct{}

func (l *linuxfc) getSCSIWWN(devicePath string) (string, error) {
	out, err := exec.Command("/lib/udev/scsi_id", "--page", "0x83", "--whitelisted", devicePath).CombinedOutput()
	outString := string(out)
	if err != nil {
		errMsg := fmt.Sprintf("Error occurred when get device wwn:", outString, err)
		log.Println(errMsg)
		return "", errors.New(errMsg)
	}
	return strings.TrimSpace(outString), nil
}

func (l *linuxfc) getContentfromSymboliclink(symboliclink string) string {
	out, _ := exec.Command("readlink", "-f", symboliclink).CombinedOutput()
	return strings.TrimSuffix(string(out), "\n")
}

func (l *linuxfc) rescanHosts(tgtWWN []string, hbas []map[string]string) error {
	for _, hba := range hbas {
		cmd := fmt.Sprintf("echo \"- - -\" > /sys/class/scsi_host/%s/scan", hba["host_device"])
		out, err := exec.Command("/bin/bash", "-c", cmd).CombinedOutput()

		outString := string(out)
		if err != nil {
			errMsg := fmt.Sprintf("Error occurred when rescan hosts", outString, err)
			log.Println(errMsg)
			return errors.New(errMsg)
		}
	}
	return nil
}

func (l *linuxfc) getFChbas() ([]map[string]string, error) {
	if !l.fcSupport() {
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
			if len(hba) > 0 {
				hbas = append(hbas, hba)
				hba = make(map[string]string)
			}
		} else {
			val := strings.Split(line, "=")
			if len(val) == 2 {
				key := strings.Replace(val[0], " ", "", -1)
				key = l.trimDoubleQuotesInText(key)

				val := strings.Replace(val[1], " ", "", -1)
				val = l.trimDoubleQuotesInText(val)

				hba[key] = val
			}
		}
		lastline = line
	}
	return hbas, nil
}

func (l *linuxfc) removeSCSIDevice(path string) error {
	cmd := "echo 1 >" + path
	out, err := exec.Command("/bin/bash", "-c", cmd).CombinedOutput()
	outString := string(out)
	if err != nil {
		errMsg := fmt.Sprintf("Error occurred when remove scsi device:", outString, err)
		log.Println(errMsg)
		return errors.New(errMsg)
	}
	return nil
}

func (l *linuxfc) flushDeviceIO(device string) error {
	cmd := "blockdev --flushbufs " + device
	out, err := exec.Command("/bin/bash", "-c", cmd).CombinedOutput()
	outString := string(out)
	if err != nil {
		errMsg := fmt.Sprintf("Error occurred when get device info when detach volume:", outString, err)
		log.Println(errMsg)
		return errors.New(errMsg)
	}
	return nil
}

func (l *linuxfc) getDeviceInfo(devicePath string) (map[string]string, error) {
	cmd := "sg_scan " + devicePath
	out, err := exec.Command("/bin/bash", "-c", cmd).CombinedOutput()
	outString := string(out)
	if err != nil {
		errMsg := fmt.Sprintf("Error occurred when get device info:", outString, err)
		log.Println(errMsg)
		return nil, errors.New(errMsg)
	}

	devInfo := make(map[string]string)
	devInfo["device"] = devicePath

	line := strings.TrimSpace(outString)

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

func (l *linuxfc) fcSupport() bool {
	var FC_HOST_SYSFC_PATH = "/sys/class/fc_host"
	return l.pathExists(FC_HOST_SYSFC_PATH)
}

func (l *linuxfc) pathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func (l *linuxfc) trimDoubleQuotesInText(str string) string {
	if strings.HasPrefix(str, "\"") && strings.HasSuffix(str, "\"") {
		return str[1 : len(str)-1]
	}
	return str
}
