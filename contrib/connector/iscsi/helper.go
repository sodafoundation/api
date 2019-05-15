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
	"fmt"
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
	AccessMode string   `mapstructure:"accessMode"`
	AuthUser   string   `mapstructure:"authUserName"`
	AuthPass   string   `mapstructure:"authPassword"`
	AuthMethod string   `mapstructure:"authMethod"`
	TgtDisco   bool     `mapstructure:"targetDiscovered"`
	TgtIQN     []string `mapstructure:"targetIQN"`
	TgtPortal  []string `mapstructure:"targetPortal"`
	VolumeID   string   `mapstructure:"volumeId"`
	TgtLun     int      `mapstructure:"targetLun"`
	Encrypted  bool     `mapstructure:"encrypted"`
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
func getInitiator() ([]string, error) {
	res, err := connector.ExecCmd("cat", "/etc/iscsi/initiatorname.iscsi")
	iqns := []string{}
	if err != nil {
		log.Printf("Error encountered gathering initiator names: %v\n", err)
		return iqns, nil
	}

	lines := strings.Split(string(res), "\n")
	for _, l := range lines {
		if strings.Contains(l, "InitiatorName=") {
			iqns = append(iqns, strings.Split(l, "=")[1])
		}
	}

	log.Printf("Found the following iqns: %s\n", iqns)
	return iqns, nil
}

// Discovery ISCSI Target
func discovery(portal string) (string, error) {
	log.Printf("Discovery portal: %s\n", portal)
	result, err := connector.ExecCmd("iscsiadm", "-m", "discovery", "-t", "sendtargets", "-p", portal)
	if err != nil {
		log.Printf("Error encountered in sendtargets: %v\n", err)
		return "", err
	}
	return strings.Replace(result, "\n", "", -1), nil
}

// Login ISCSI Target
func setAuth(portal string, targetiqn string, name string, passwd string) error {
	// Set UserName
	info, err := connector.ExecCmd("iscsiadm", "-m", "node", "-p", portal, "-T", targetiqn,
		"--op=update", "--name", "node.session.auth.username", "--value", name)
	if err != nil {
		log.Printf("Received error on set income username: %v, %v\n", err, info)
		return err
	}
	// Set Password
	info, err = connector.ExecCmd("iscsiadm", "-m", "node", "-p", portal, "-T", targetiqn,
		"--op=update", "--name", "node.session.auth.password", "--value", passwd)
	if err != nil {
		log.Printf("Received error on set income password: %v, %v\n", err, info)
		return err
	}
	return nil
}

// Login ISCSI Target
func login(portal string, targetiqn string) error {
	log.Printf("Login portal: %s targetiqn: %s\n", portal, targetiqn)
	// Do not login again if there is an active session.
	cmd := "iscsiadm -m session |grep -w " + portal + "|grep -w " + targetiqn
	_, err := connector.ExecCmd("/bin/bash", "-c", cmd)
	if err == nil {
		log.Printf("there is an active session\n")
		_, err := connector.ExecCmd("iscsiadm", "-m", "session", "-R")
		if err == nil {
			log.Printf("rescan iscsi session success.\n")
		}
		return nil
	}

	info, err := connector.ExecCmd("iscsiadm", "-m", "node", "-p", portal, "-T", targetiqn, "--login")
	if err != nil {
		log.Printf("Received error on login attempt: %v, %s\n", err, info)
		return err
	}
	return nil
}

// Logout ISCSI Target
func logout(portal string, targetiqn string) error {
	log.Printf("Logout portal: %s targetiqn: %s\n", portal, targetiqn)
	info, err := connector.ExecCmd("iscsiadm", "-m", "node", "-p", portal, "-T", targetiqn, "--logout")
	if err != nil {
		log.Println("Received error on logout attempt", err, info)
		return err
	}
	return nil
}

// Delete ISCSI Node
func delete(targetiqn string) (err error) {
	log.Printf("Delete targetiqn: %s\n", targetiqn)
	_, err = connector.ExecCmd("iscsiadm", "-m", "node", "-o", "delete", "-T", targetiqn)
	if err != nil {
		log.Printf("Received error on Delete attempt: %v\n", err)
		return err
	}
	return nil
}

// ParseIscsiConnectInfo decode
func parseIscsiConnectInfo(connectInfo map[string]interface{}) (*IscsiConnectorInfo, int, error) {
	var con IscsiConnectorInfo
	mapstructure.Decode(connectInfo, &con)

	fmt.Printf("iscsi target portal: %s, target iqn: %s, target lun: %d\n", con.TgtPortal, con.TgtIQN, con.TgtLun)
	if len(con.TgtPortal) == 0 || con.TgtLun == 0 {
		return nil, -1, errors.New("iscsi connection data invalid.")
	}

	var index int

	log.Printf("TgtPortal:%v\n", con.TgtPortal)
	for i, portal := range con.TgtPortal {
		strs := strings.Split(portal, ":")
		ip := strs[0]
		cmd := "ping -c 2 " + ip
		res, err := connector.ExecCmd("/bin/bash", "-c", cmd)
		log.Printf("ping result:%v\n", res)
		if err != nil {
			log.Printf("ping error:%v\n", err)
			if i == len(con.TgtPortal)-1 {
				return nil, -1, errors.New("no available iscsi portal.")
			}
			continue
		}
		index = i
		break
	}

	return &con, index, nil
}

// Connect ISCSI Target
func connect(connMap map[string]interface{}) (string, error) {
	conn, index, err := parseIscsiConnectInfo(connMap)
	if err != nil {
		return "", err
	}
	log.Println("connmap info: ", connMap)
	log.Println("conn info is: ", conn)
	portal := conn.TgtPortal[index]

	var targetiqn string
	var targetiqnIdx = 1
	if len(conn.TgtIQN) == 0 {
		content, _ := discovery(portal)
		targetiqn = strings.Split(content, " ")[targetiqnIdx]
	} else {
		targetiqn = conn.TgtIQN[index]
	}

	targetlun := strconv.Itoa(conn.TgtLun)

	cmd := "pgrep -f /sbin/iscsid"
	_, err = connector.ExecCmd("/bin/bash", "-c", cmd)

	if err != nil {
		cmd = "/sbin/iscsid"
		_, errExec := connector.ExecCmd("/bin/bash", "-c", cmd)
		if errExec != nil {
			return "", fmt.Errorf("Please stop the iscsi process outside the container first: %v", errExec)
		}
	}

	log.Printf("Connect portal: %s targetiqn: %s targetlun: %s\n", portal, targetiqn, targetlun)
	devicePath := strings.Join([]string{
		"/dev/disk/by-path/ip",
		portal,
		"iscsi",
		targetiqn,
		"lun",
		targetlun}, "-")

	log.Println("devicepath is ", devicePath)

	// Discovery
	_, err = discovery(portal)
	if err != nil {
		return "", err
	}
	if len(conn.AuthMethod) != 0 {
		setAuth(portal, targetiqn, conn.AuthUser, conn.AuthPass)
	}
	//Login
	err = login(portal, targetiqn)
	if err != nil {
		return "", err
	}

	isexist := waitForPathToExist(&devicePath, 10, ISCSITranslateTCP)

	if !isexist {
		return "", errors.New("Could not connect volume: Timeout after 10s")
	}

	return devicePath, nil
}

// Disconnect ISCSI Target
func disconnect(conn map[string]interface{}) error {
	iscsiCon, index, err := parseIscsiConnectInfo(conn)
	if err != nil {
		return err
	}
	portal := iscsiCon.TgtPortal[index]

	var targetiqn string
	if len(iscsiCon.TgtIQN) == 0 {
		content, _ := discovery(portal)
		targetiqn = strings.Split(content, " ")[1]
	} else {
		targetiqn = iscsiCon.TgtIQN[index]
	}

	cmd := "ls /dev/disk/by-path/ |grep -w " + portal + "|grep -w " + targetiqn + "|wc -l |awk '{if($1>1) print 1; else print 0}'"
	logoutFlag, err := connector.ExecCmd("/bin/bash", "-c", cmd)
	if err != nil {
		log.Printf("Disconnect iscsi target failed, %v\n", err)
		return err
	}

	logoutFlag = strings.Replace(logoutFlag, "\n", "", -1)
	if logoutFlag == "0" {
		log.Printf("Disconnect portal: %s targetiqn: %s\n", portal, targetiqn)
		// Logout
		err = logout(portal, targetiqn)
		if err != nil {
			return err
		}

		//Delete
		err = delete(targetiqn)
		if err != nil {
			return err
		}

		return nil
	}
	log.Println("logoutFlag: ", logoutFlag)
	return nil
}

func getTgtPortalAndTgtIQN() (string, string, error) {
	log.Println("GetTgtPortalAndTgtIQN")
	var targetiqn, targetportal string
	out, err := connector.ExecCmd("iscsiadm", "-m", "session")
	if err != nil {
		errGetPortalAndIQN := fmt.Errorf("Get targetportal And targetiqn failed: %v", err)
		log.Println("Get targetportal And targetiqn failed: ", errGetPortalAndIQN)
		return "", "", errGetPortalAndIQN
	}

	lines := strings.Split(string(out), "\n")

	for _, line := range lines {
		if strings.Contains(line, "tcp") {
			lineSplit := strings.Split(line, " ")
			targetportalTemp := lineSplit[2]
			targetportal = strings.Split(targetportalTemp, ",")[0]
			targetiqn = lineSplit[3]
		}
	}

	if targetiqn != "" && targetportal != "" {
		return targetiqn, targetportal, nil
	}

	msg := "targetportal And targetiqn not found"
	log.Println(msg)
	return "", "", errors.New(msg)

}

func getInitiatorInfo() (string, error) {
	initiators, err := getInitiator()
	if err != nil {
		return "", err
	}

	if len(initiators) == 0 {
		return "", errors.New("No iqn found")
	}

	if len(initiators) > 1 {
		return "", errors.New("the number of iqn is wrong")
	}

	return initiators[0], nil
}
