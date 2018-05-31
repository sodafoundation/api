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
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
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
	res, err := execCmd("cat", "/etc/iscsi/initiatorname.iscsi")
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
	info, err := execCmd("iscsiadm", "-m", "node", "-p", portal, "-T", targetiqn,
		"--op=update", "--name", "node.session.auth.username", "--value", name)
	if err != nil {
		log.Fatalf("Received error on set income username: %v, %v", err, info)
		return err
	}
	// Set Password
	info, err = execCmd("iscsiadm", "-m", "node", "-p", portal, "-T", targetiqn,
		"--op=update", "--name", "node.session.auth.password", "--value", passwd)
	if err != nil {
		log.Fatalf("Received error on set income password: %v, %v", err, info)
		return err
	}
	return nil
}

// Discovery ISCSI Target
func Discovery(portal string) error {
	info, err := execCmd("iscsiadm", "-m", "discovery", "-t", "sendtargets", "-p", portal)
	if err != nil {
		log.Println("Error encountered in sendtargets:", string(info), err)
		return err
	}
	return nil
}

// Login ISCSI Target
func Login(portal string, targetiqn string) error {
	info, err := execCmd("iscsiadm", "-m", "node", "-p", portal, "-T", targetiqn, "--login")
	if err != nil {
		log.Println("Received error on login attempt:", string(info), err)
		return err
	}
	return nil
}

// Logout ISCSI Target
func Logout(portal string, targetiqn string) error {
	info, err := execCmd("iscsiadm", "-m", "node", "-p", portal, "-T", targetiqn, "--logout")
	if err != nil {
		log.Println("Received error on logout attempt:", string(info), err)
		return err
	}
	return nil
}

// Delete ISCSI Node
func Delete(targetiqn string) error {
	info, err := execCmd("iscsiadm", "-m", "node", "-o", "delete", "-T", targetiqn)
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
	info, err := execCmd("iscsiadm", "-m", "session", "-s")
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
	_, err := execCmd("iscsiadm", "-m", "node", "-o", "show",
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

// GetFSType returns the File System Type of device
func GetFSType(device string) string {
	fsType := ""
	res, err := execCmd("blkid", device)
	if err != nil {
		log.Printf("failed to GetFSType: %v", err)
		return fsType
	}

	if strings.Contains(string(res), "TYPE=") {
		for _, v := range strings.Split(string(res), " ") {
			if strings.Contains(v, "TYPE=") {
				fsType = strings.Split(v, "=")[1]
				fsType = strings.Replace(fsType, "\"", "", -1)
			}
		}
	}
	return fsType
}

// Format device by File System Type
func Format(device string, fstype string) error {
	log.Printf("Format device: %s fstype: %s", device, fstype)

	// Get current File System Type
	curFSType := GetFSType(device)
	if curFSType == "" {
		// Default File Sysem Type is ext4
		if fstype == "" {
			fstype = "ext4"
		}
		_, err := execCmd("mkfs", "-t", fstype, "-F", device)
		if err != nil {
			log.Printf("failed to Format: %v", err)
			return err
		}
	} else {
		log.Printf("Device: %s has been formatted yet. fsType: %s", device, curFSType)
	}
	return nil
}

// Mount device into mount point
func Mount(device string, mountpoint string) error {
	log.Printf("Mount device: %s mountpoint: %s", device, mountpoint)

	_, err := execCmd("mkdir", "-p", mountpoint)
	if err != nil {
		log.Printf("failed to mkdir: %v", err)
	}
	_, err = execCmd("mount", device, mountpoint)
	if err != nil {
		log.Printf("failed to mount: %v", err)
		return err
	}
	return nil
}

// FormatAndMount device
func FormatAndMount(device string, fstype string, mountpoint string) error {
	log.Printf("FormatandMount device: %s fstype: %s mountpoint: %s", device, fstype, mountpoint)

	// Format
	err := Format(device, fstype)
	if err != nil {
		return err
	}

	// Mount
	err = Mount(device, mountpoint)
	if err != nil {
		return err
	}

	return nil
}

// Umount from mountpoint
func Umount(mountpoint string) error {
	log.Printf("Umount mountpoint: %s", mountpoint)

	_, err := execCmd("umount", mountpoint)
	if err != nil {
		log.Printf("failed to Umount: %v", err)
		return err
	}
	return nil
}

// ParseIscsiConnectInfo decode
func ParseIscsiConnectInfo(connectInfo map[string]interface{}) *IscsiConnectorInfo {
	var con IscsiConnectorInfo
	mapstructure.Decode(connectInfo, &con)
	return &con
}

// GetHostIp return Host IP
func GetHostIp() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "127.0.0.1"
	}

	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			return ipnet.IP.String()
		}
	}

	return "127.0.0.1"
}

func execCmd(name string, arg ...string) (string, error) {
	log.Printf("Command: %s %s:\n", name, strings.Join(arg, " "))
	info, err := exec.Command(name, arg...).CombinedOutput()
	return string(info), err
}
