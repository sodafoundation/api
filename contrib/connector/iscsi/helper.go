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

package iscsi

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/opensds/opensds/contrib/connector"
)

// IscsiConnectorInfo define
type IscsiConnectorInfo struct {
	AccessMode string `mapstructure:"accessMode"`
	AuthUser   string `mapstructure:"authUserName"`
	AuthPass   string `mapstructure:"authPassword"`
	AuthMethod string `mapstructure:"authMethod"`
	TgtDisco   bool   `mapstructure:"targetDiscovered"`
	TgtIQN     string `mapstructure:"targetIqn"`
	TgtPortal  string `mapstructure:"targetPortal"`
	VolumeID   string `mapstructure:"volumeId"`
	TgtLun     int    `mapstructure:"targetLun"`
	Encrypted  bool   `mapstructure:"encrypted"`
}

////////////////////////////////////////////////////////////////////////////////
//      Refer some codes from: https://github.com/j-griffith/csi-cinder       //
//      Refer some codes from: https://github.com/kubernetes/kubernetes       //
////////////////////////////////////////////////////////////////////////////////

const (
	//ISCSITranslateTCP tcp
	ISCSITranslateTCP = "tcp"
)

// statFunc define
type statFunc func(string) (os.FileInfo, error)

// globFunc define
type globFunc func(string) ([]string, error)

// waitForPathToExist scan the device path
func waitForPathToExist(devicePath *string, maxRetries int, deviceTransport string) bool {
	// This makes unit testing a lot easier
	return waitForPathToExistInternal(devicePath, maxRetries, deviceTransport, os.Stat, filepath.Glob)
}

// waitForPathToExistInternal scan the device path
func waitForPathToExistInternal(devicePath *string, maxRetries int, deviceTransport string, osStat statFunc, filepathGlob globFunc) bool {
	if devicePath == nil {
		return false
	}

	for i := 0; i < maxRetries; i++ {
		var err error
		if deviceTransport == ISCSITranslateTCP {
			_, err = osStat(*devicePath)
		} else {
			fpath, _ := filepathGlob(*devicePath)
			if fpath == nil {
				err = os.ErrNotExist
			} else {
				// There might be a case that fpath contains multiple device paths if
				// multiple PCI devices connect to same iscsi target. We handle this
				// case at subsequent logic. Pick up only first path here.
				*devicePath = fpath[0]
			}
		}
		if err == nil {
			return true
		}
		if !os.IsNotExist(err) {
			return false
		}
		if i == maxRetries-1 {
			break
		}
		time.Sleep(time.Second)
	}
	return false
}

// GetInitiator returns all the ISCSI Initiator Name
func GetInitiator() ([]string, error) {
	res, err := connector.ExecCmd("cat", "/etc/iscsi/initiatorname.iscsi")
	iqns := []string{}
	if err != nil {
		log.Printf("Error encountered gathering initiator names: %v", err)
		return iqns, nil
	}

	lines := strings.Split(string(res), "\n")
	for _, l := range lines {
		if strings.Contains(l, "InitiatorName=") {
			iqns = append(iqns, strings.Split(l, "=")[1])
		}
	}

	log.Printf("Found the following iqns: %s", iqns)
	return iqns, nil
}

// Login ISCSI Target
func SetAuth(portal string, targetiqn string, name string, passwd string) error {
	// Set UserName
	info, err := connector.ExecCmd("iscsiadm", "-m", "node", "-p", portal, "-T", targetiqn,
		"--op=update", "--name", "node.session.auth.username", "--value", name)
	if err != nil {
		log.Fatalf("Received error on set income username: %v, %v", err, info)
		return err
	}
	// Set Password
	info, err = connector.ExecCmd("iscsiadm", "-m", "node", "-p", portal, "-T", targetiqn,
		"--op=update", "--name", "node.session.auth.password", "--value", passwd)
	if err != nil {
		log.Fatalf("Received error on set income password: %v, %v", err, info)
		return err
	}
	return nil
}

// Discovery ISCSI Target
func Discovery(portal string) error {
	info, err := connector.ExecCmd("iscsiadm", "-m", "discovery", "-t", "sendtargets", "-p", portal)
	if err != nil {
		log.Println("Error encountered in sendtargets:", string(info), err)
		return err
	}
	return nil
}

// Login ISCSI Target
func Login(portal string, targetiqn string) error {
	info, err := connector.ExecCmd("iscsiadm", "-m", "node", "-p", portal, "-T", targetiqn, "--login")
	if err != nil {
		log.Println("Received error on login attempt:", string(info), err)
		return err
	}
	return nil
}

// Logout ISCSI Target
func Logout(portal string, targetiqn string) error {
	info, err := connector.ExecCmd("iscsiadm", "-m", "node", "-p", portal, "-T", targetiqn, "--logout")
	if err != nil {
		log.Println("Received error on logout attempt:", string(info), err)
		return err
	}
	return nil
}

// Delete ISCSI Node
func Delete(targetiqn string) error {
	info, err := connector.ExecCmd("iscsiadm", "-m", "node", "-o", "delete", "-T", targetiqn)
	if err != nil {
		log.Println("Received error on Delete attempt:", string(info), err)
		return err
	}
	return nil
}

// Connect ISCSI Target
func Connect(connMap map[string]interface{}) (string, error) {
	conn := ParseIscsiConnectInfo(connMap)
	portal := conn.TgtPortal
	targetiqn := conn.TgtIQN
	targetlun := strconv.Itoa(conn.TgtLun)

	devicePath := strings.Join([]string{
		"/dev/disk/by-path/ip",
		portal,
		"iscsi",
		targetiqn,
		"lun",
		targetlun}, "-")

	isexist := waitForPathToExist(&devicePath, 1, ISCSITranslateTCP)
	if !isexist {

		// Discovery
		err := Discovery(portal)
		if err != nil {
			return "", err
		}
		// Set authentication messages,if is has.
		if len(conn.AuthMethod) != 0 {
			SetAuth(portal, targetiqn, conn.AuthUser, conn.AuthPass)
		}
		//Login
		err = Login(portal, targetiqn)
		if err != nil {
			return "", err
		}

		isexist = waitForPathToExist(&devicePath, 10, ISCSITranslateTCP)
		if !isexist {
			return "", errors.New("Could not connect volume: Timeout after 10s")
		}

	}
	return devicePath, nil
}

func sessionExists(portal string, tgtIqn string) bool {
	info, err := connector.ExecCmd("iscsiadm", "-m", "session", "-s")
	if err != nil {
		log.Println("Warning: get session failed,", string(info))
		return false
	}
	portal = strings.Replace(portal, ":", ",", -1)
	for _, line := range strings.Split(string(info), "\n") {
		if strings.Contains(line, tgtIqn) && strings.Contains(line, portal) {
			return true
		}
	}
	return false
}

func recordExists(portal string, tgtIqn string) bool {
	_, err := connector.ExecCmd("iscsiadm", "-m", "node", "-o", "show",
		"-T", tgtIqn, "-p", portal)
	return err == nil
}

// Disconnect ISCSI Target
func Disconnect(portal string, targetiqn string) error {
	log.Printf("Disconnect portal: %s targetiqn: %s", portal, targetiqn)
	if sessionExists(portal, targetiqn) {
		if err := Logout(portal, targetiqn); err != nil {
			return err
		}
	}

	if recordExists(portal, targetiqn) {
		return Delete(targetiqn)
	}
	return nil
}

// ParseIscsiConnectInfo decode
func ParseIscsiConnectInfo(connectInfo map[string]interface{}) *IscsiConnectorInfo {
	var con IscsiConnectorInfo
	mapstructure.Decode(connectInfo, &con)
	return &con
}

// getInitiatorInfo implementation
func getInitiatorInfo() (connector.InitiatorInfo, error) {
	var initiatorInfo connector.InitiatorInfo

	initiators, err := GetInitiator()
	if err != nil {
		return initiatorInfo, err
	}

	if len(initiators) == 0 {
		return initiatorInfo, errors.New("The number of iqn is wrong")
	}

	initiatorInfo.InitiatorData = make(map[string]interface{})
	initiatorInfo.InitiatorData[Iqn] = initiators[0]

	hostName, err := connector.GetHostName()
	if err != nil {
		return initiatorInfo, err
	}

	initiatorInfo.HostName = hostName
	log.Printf("getFChbasInfo success: protocol=%v, initiatorInfo=%v",
		iscsiDriver, initiatorInfo)

	return initiatorInfo, nil
}
