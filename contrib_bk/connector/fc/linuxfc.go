// Copyright (c) 2019 Huawei Technologies Co., Ltd. All Rights Reserved.
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
	"strings"

	"github.com/opensds/opensds/contrib/connector"
)

func getSCSIWWN(devicePath string) (string, error) {
	out, err := connector.ExecCmd("/lib/udev/scsi_id", "--page", "0x83", "--whitelisted", devicePath)
	if err != nil {
		errMsg := fmt.Sprintf("Error occurred when get device wwn: %s, %v\n", out, err)
		log.Printf(errMsg)
		return "", errors.New(errMsg)
	}
	return strings.TrimSpace(out), nil
}

func getContentfromSymboliclink(symboliclink string) string {
	out, _ := connector.ExecCmd("readlink", "-f", symboliclink)
	return strings.TrimSuffix(out, "\n")
}

func rescanHosts(tgtWWN []string, hbas []map[string]string) error {
	for _, hba := range hbas {
		cmd := fmt.Sprintf("echo \"- - -\" > /sys/class/scsi_host/%s/scan", hba["host_device"])
		out, err := connector.ExecCmd("/bin/bash", "-c", cmd)
		if err != nil {
			errMsg := fmt.Sprintf("Error occurred when rescan hosts: %s, %v\n", out, err)
			log.Printf(errMsg)
			return errors.New(errMsg)
		}
	}
	return nil
}

func getFChbas() ([]map[string]string, error) {
	if !fcSupport() {
		errMsg := fmt.Sprintf("No Fibre Channel support detected.\n")
		log.Printf(errMsg)
		return nil, errors.New(errMsg)
	}

	out, err := connector.ExecCmd("systool", "-c", "fc_host", "-v")
	if err != nil {
		errMsg := fmt.Sprintf("Error occurred when get FC hbas info: systool is not installed: %s, %v\n", out, err)
		log.Printf(errMsg)
		return nil, errors.New(errMsg)
	}

	if out == "" {
		errMsg := fmt.Sprintf("No Fibre Channel support detected.\n")
		log.Printf(errMsg)
		return nil, errors.New(errMsg)
	}

	lines := strings.Split(out, "\n")
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
				key = trimDoubleQuotesInText(key)

				val := strings.Replace(val[1], " ", "", -1)
				val = trimDoubleQuotesInText(val)

				hba[key] = val
			}
		}
		lastline = line
	}
	return hbas, nil
}

func removeSCSIDevice(path string) error {
	cmd := "echo 1 >" + path
	out, err := connector.ExecCmd("/bin/bash", "-c", cmd)
	if err != nil {
		errMsg := fmt.Sprintf("Error occurred when remove scsi device: %s, %v\n", out, err)
		log.Printf(errMsg)
		return errors.New(errMsg)
	}
	return nil
}

func flushDeviceIO(device string) error {
	cmd := "blockdev --flushbufs " + device
	out, err := connector.ExecCmd("/bin/bash", "-c", cmd)
	if err != nil {
		errMsg := fmt.Sprintf("Error occurred when get device info when detach volume: %s, %v\n", out, err)
		log.Printf(errMsg)
		return errors.New(errMsg)
	}
	return nil
}

func getDeviceInfo(devicePath string) (map[string]string, error) {
	cmd := "sg_scan " + devicePath
	out, err := connector.ExecCmd("/bin/bash", "-c", cmd)
	if err != nil {
		errMsg := fmt.Sprintf("Error occurred when get device info: %s, %v\n", out, err)
		log.Printf(errMsg)
		return nil, errors.New(errMsg)
	}

	devInfo := make(map[string]string)
	devInfo["device"] = devicePath

	line := strings.TrimSpace(out)

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

func fcSupport() bool {
	var FcHostSYSFcPATH = "/sys/class/fc_host"
	return pathExists(FcHostSYSFcPATH)
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

func trimDoubleQuotesInText(str string) string {
	if strings.HasPrefix(str, "\"") && strings.HasSuffix(str, "\"") {
		return str[1 : len(str)-1]
	}
	return str
}
