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

package rbd

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/opensds/opensds/contrib/connector"
)

var (
	rbdBusPath    = "/sys/bus/rbd"
	rbdDevicePath = path.Join(rbdBusPath, "devices")
	rbdDev        = "/dev/rbd"
)

type RBD struct{}

var _ connector.Connector = &RBD{}

func init() {
	connector.RegisterConnector(connector.RbdDriver, &RBD{})
}

func (*RBD) Attach(conn map[string]interface{}) (string, error) {
	if _, ok := conn["name"]; !ok {
		return "", fmt.Errorf("cann't get name in connection")
	}

	name, ok := conn["name"].(string)
	if !ok {
		return "", fmt.Errorf("invalid connection name %v", conn["name"])
	}
	hosts, ok := conn["hosts"].([]string)
	if !ok {
		return "", fmt.Errorf("invalid connection hosts %v", conn["hosts"])
	}

	ports, ok := conn["ports"].([]string)
	if !ok {
		return "", fmt.Errorf("invalid connection ports %v", conn["hosts"])
	}

	device, err := mapDevice(name, hosts, ports)
	if err != nil {
		return "", err
	}

	return device, nil
}

func (*RBD) Detach(conn map[string]interface{}) error {
	if _, ok := conn["name"]; !ok {
		return os.ErrInvalid
	}
	name, ok := conn["name"].(string)
	if !ok {
		return fmt.Errorf("invalid connection name %v", conn["name"])
	}
	device, err := findDevice(name, 1)
	if err != nil {
		return err
	}

	_, err = exec.Command("rbd", "unmap", device).CombinedOutput()
	return err
}

// GetInitiatorInfo implementation
func (*RBD) GetInitiatorInfo() (string, error) {
	hostName, err := connector.GetHostName()

	if err != nil {
		return "", err
	}

	return hostName, nil
}

func parseName(name string) (poolName, imageName, snapName string, err error) {
	fields := strings.Split(name, "/")
	if len(fields) != 2 {
		err = fmt.Errorf("invalid connection name %s", name)
		return
	}
	poolName, imageName, snapName = fields[0], fields[1], "-"

	imgAndSnap := strings.Split(fields[1], "@")
	if len(imgAndSnap) == 2 {
		imageName, snapName = imgAndSnap[0], imgAndSnap[1]
	}
	return
}

func mapDevice(name string, hosts, ports []string) (string, error) {
	devName, err := findDevice(name, 1)
	if err == nil {
		return devName, nil
	}

	// modprobe
	exec.Command("modprobe", "rbd").CombinedOutput()

	for i := 0; i < len(hosts); i++ {
		_, err = exec.Command("rbd", "map", name).CombinedOutput()
		if err == nil {
			break
		}
	}

	devName, err = findDevice(name, 10)
	if err != nil {
		return "", err
	}

	return devName, nil
}

func findDevice(name string, retries int) (string, error) {
	poolName, imageName, snapName, err := parseName(name)
	if err != nil {
		return "", err
	}

	for i := 0; i < retries; i++ {
		if name, err := findDeviceTree(poolName, imageName, snapName); err == nil {
			if _, err := os.Stat(rbdDev + name); err != nil {
				return "", err
			}

			return rbdDev + name, nil
		}

		time.Sleep(time.Second)
	}

	return "", os.ErrNotExist
}

func findDeviceTree(poolName, imageName, snapName string) (string, error) {
	fi, err := ioutil.ReadDir(rbdDevicePath)
	if err != nil && err != os.ErrNotExist {
		return "", err
	} else if err == os.ErrNotExist {
		return "", fmt.Errorf("Could not locate devices directory")
	}

	for _, f := range fi {
		namePath := filepath.Join(rbdDevicePath, f.Name(), "name")
		content, err := ioutil.ReadFile(namePath)
		if err != nil {
			return "", err
		}
		if strings.TrimSpace(string(content)) != imageName {
			continue
		}

		poolPath := filepath.Join(rbdDevicePath, f.Name(), "pool")
		content, err = ioutil.ReadFile(poolPath)
		if err != nil {
			return "", err
		}
		if strings.TrimSpace(string(content)) != poolName {
			continue
		}

		snapPath := filepath.Join(rbdDevicePath, f.Name(), "current_snap")
		content, err = ioutil.ReadFile(snapPath)
		if err != nil {
			return "", err
		}
		if strings.TrimSpace(string(content)) == snapName {
			return f.Name(), nil
		}
	}

	return "", os.ErrNotExist
}
