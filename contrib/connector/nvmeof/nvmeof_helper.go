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

package nvmeof

import (
	"errors"
	"log"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/opensds/opensds/contrib/connector"
)

const (
	iniNvmePrefix = "nqn.ini."
)

// ConnectorInfo define
type ConnectorInfo struct {
	Nqn       string `mapstructure:"targetNQN"`    //NVMe subsystem name to the volume to be connected
	TgtPort   string `mapstructure:"targetPort"`   //NVMe target port that hosts the nqn sybsystem
	TgtPortal string `mapstructure:"targetIP"`     //NVMe target ip that hosts the nqn sybsystem
	TranType  string `mapstructure:"transporType"` // Nvme transport type
	HostNqn   string `mapstructure:"hostNqn"`      // host nqn
}

//////////////////////////////////////////////////////////////////////////////////////////
//      Refer some codes from: https://github.intel.com/yingxinc/cinder-rsd-os-brick    //
//////////////////////////////////////////////////////////////////////////////////////////

// GetInitiator returns all the Nvmeof UUID
func GetInitiator() ([]string, error) {
	res, err := connector.ExecCmd("dmidecode")
	nqns := []string{}
	if err != nil {
		log.Printf("Unable to execute dmidecode,Error encountered gathering Nvmeof UUID: %v\n", err)
		return nqns, nil
	}

	lines := strings.Split(string(res), "\n")
	for _, l := range lines {
		if strings.Contains(l, "UUID: ") {
			tmp := iniNvmePrefix + strings.Split(l, ": ")[1]
			nqns = append(nqns, tmp)
			log.Printf("Found the following nqns: %s", nqns)
			return nqns, nil
		}
	}
	log.Println("can not find any nqn initiator")
	return nqns, errors.New("can not find any nqn initiator")
}

func getInitiatorInfo() (string, error) {

	initiators, err := GetInitiator()
	if err != nil {
		return "", err
	}

	if len(initiators) == 0 {
		return "", errors.New("no nqn found")
	}

	if len(initiators) > 1 {
		return "", errors.New("the number of nqn is wrong")
	}

	hostName, err := connector.GetHostName()
	if err != nil {
		return "", errors.New("can not get hostname")
	}

	info := initiators[0] + "." + hostName
	return info, nil
}

// GetNvmeDevice get all the nvme devices
func GetNvmeDevice() (map[string]int, error) {
	nvmeDevice := make(map[string]int)
	pattern := "/dev/nvme"
	Npath, err := connector.ExecCmd("nvme", "list")
	if err != nil {
		return nvmeDevice, err
	}
	log.Println("nvme list succeed")
	lines := strings.Split(string(Npath), "\n")
	for _, l := range lines {
		if strings.Contains(l, pattern) {
			name := strings.Split(l, " ")[0]
			nvmeDevice[name] = 1
		}
	}
	return nvmeDevice, err
}

// GetNvmeSubsystems :list connected target name
func GetNvmeSubsystems() (map[string]int, error) {
	nqn := make(map[string]int)
	res, err := connector.ExecCmd("nvme", "list-subsys")
	if err != nil {
		return nqn, err
	}

	lines := strings.Split(string(res), "\n")
	for _, l := range lines {
		if strings.Contains(l, "NQN=") {
			name := strings.Split(l, "NQN=")[1]
			nqn[name] = 1
		}
	}

	log.Printf("Found the following NQN: %s", res)
	return nqn, nil
}

// Discovery NVMe-OF target
func Discovery(connMap map[string]interface{}) error {
	conn := ParseNvmeofConnectInfo(connMap)
	targetip := conn.TgtPortal
	targetport := conn.TgtPort
	info, err := connector.ExecCmd("nvme", "discover", "-t", "tcp", "-a", targetip, "-s", targetport)
	if err != nil {
		log.Printf("Error encountered in send targets:%v, %v\n", err, info)
		return err
	}
	return nil
}

// Connect NVMe-OF Target ,return the new target device path in this node
func Connect(connMap map[string]interface{}) (string, error) {
	CurrentNvmeDevice, _ := GetNvmeDevice()
	conn := ParseNvmeofConnectInfo(connMap)
	connNqn := conn.Nqn
	targetPortal := conn.TgtPortal
	port := conn.TgtPort
	nvmeTransportType := conn.TranType
	hostName := conn.HostNqn

	cmd := "nvme connect -t " + nvmeTransportType + " -n " + connNqn + " -s " + port + " -a " + targetPortal
	if hostName != "ALL" {
		cmd += " -q " + hostName
	}
	//log.Printf("conn information:%s, %s, %s ", connNqn, targetPortal, port)

	_, err := connector.ExecCmd("/bin/bash", "-c", cmd)

	if err != nil {
		log.Println("Failed to connect to NVMe nqn :", connNqn)
		return "", err
	}

	for retry := 0; retry < 10; retry++ {
		allNvmeDevices, _ := GetNvmeDevice()
		for p, _ := range allNvmeDevices {
			if _, ok := CurrentNvmeDevice[p]; !ok {
				log.Printf("NVMe device to be connected to is : %v", p)
				return p, nil
			}
			time.Sleep(time.Second)
		}
	}
	return "", errors.New("could not connect volume: Timeout after 10s")
}

// DisConnect nvme device by name
func DisConnect(nqn string) error {
	currentNvmeNames, err := GetNvmeSubsystems()
	if err != nil {
		log.Println("can not get nvme device")
		return err
	}
	if _, ok := currentNvmeNames[nqn]; !ok {
		log.Println("Trying to disconnect nqn" + nqn +
			"is not connected.")
		return errors.New("device path not found ")
	}

	_, err = connector.ExecCmd("nvme", "disconnect", "-n", nqn)
	if err != nil {
		log.Println("could not disconnect nvme nqn ： ", nqn)
		return err
	}
	log.Println(" disconnect nvme nqn ： ", nqn)
	return nil
}

// ParseNvmeofConnectInfo decode
func ParseNvmeofConnectInfo(connectInfo map[string]interface{}) *ConnectorInfo {
	var con ConnectorInfo
	mapstructure.Decode(connectInfo, &con)
	return &con
}
