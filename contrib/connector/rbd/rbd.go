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

const (
	rbdDriver = "rbd"
)

var (
	rbdBusPath    = "/sys/bus/rbd"
	rbdDevicePath = path.Join(rbdBusPath, "devices")
	rbdDev        = "/dev/rbd"
)

type RBD struct{}

var _ connector.Connector = &RBD{}

func init() {
	connector.RegisterConnector(rbdDriver, &RBD{})
}

func (*RBD) Attach(conn map[string]interface{}) (string, error) {
	if _, ok := conn["name"]; !ok {
		return "", os.ErrInvalid
	}

	name := conn["name"].(string)
	fields := strings.Split(name, "/")
	if len(fields) != 2 {
		return "", os.ErrInvalid
	}

	if _, ok := conn["hosts"].([]interface{}); !ok {
		return "", os.ErrInvalid
	}
	hosts := conn["hosts"].([]interface{})

	if _, ok := conn["ports"].([]interface{}); !ok {
		return "", os.ErrInvalid
	}
	ports := conn["ports"].([]interface{})

	poolName, imageName := fields[0], fields[1]
	device, err := mapDevice(poolName, imageName, hosts, ports)
	if err != nil {
		return "", err
	}

	return device, nil
}

func (*RBD) Detach(conn map[string]interface{}) error {
	if _, ok := conn["name"]; !ok {
		return os.ErrInvalid
	}

	name := conn["name"].(string)
	fields := strings.Split(name, "/")
	if len(fields) != 2 {
		return os.ErrInvalid
	}

	poolName, imageName := fields[0], fields[1]
	device, err := findDevice(poolName, imageName, 1)
	if err != nil {
		return err
	}

	_, err = exec.Command("rbd", "unmap", device).CombinedOutput()
	return err
}

func mapDevice(poolName, imageName string, hosts, ports []interface{}) (string, error) {
	devName, err := findDevice(poolName, imageName, 1)
	if err == nil {
		return devName, nil
	}

	// modprobe
	exec.Command("modprobe", "rbd").CombinedOutput()

	for i := 0; i < len(hosts); i++ {
		_, err = exec.Command("rbd", "map", imageName, "--pool", poolName).CombinedOutput()
		if err == nil {
			break
		}
	}

	devName, err = findDevice(poolName, imageName, 10)
	if err != nil {
		return "", err
	}

	return devName, nil
}

func findDevice(poolName, imageName string, retries int) (string, error) {
	for i := 0; i < retries; i++ {
		if name, err := findDeviceTree(poolName, imageName); err == nil {
			if _, err := os.Stat(rbdDev + name); err != nil {
				return "", err
			}

			return rbdDev + name, nil
		}

		time.Sleep(time.Second)
	}

	return "", os.ErrNotExist
}

func findDeviceTree(poolName, imageName string) (string, error) {
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

		if strings.TrimSpace(string(content)) == imageName {
			poolPath := filepath.Join(rbdDevicePath, f.Name(), "pool")
			content, err := ioutil.ReadFile(poolPath)
			if err != nil {
				return "", err
			}

			if strings.TrimSpace(string(content)) == poolName {
				return f.Name(), err
			}
		}
	}

	return "", os.ErrNotExist
}
