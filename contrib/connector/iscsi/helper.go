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
	res, err := exec.Command("cat", "/etc/iscsi/initiatorname.iscsi").CombinedOutput()
	log.Printf("result from cat: %s", res)
	iqns := []string{}
	if err != nil {
		log.Printf("Error encountered gathering initiator names: %v", err)
		return iqns, nil
	}

	lines := strings.Split(string(res), "\n")
	for _, l := range lines {
		log.Printf("Inspect line: %s", l)
		if strings.Contains(l, "InitiatorName=") {
			iqns = append(iqns, strings.Split(l, "=")[1])
		}
	}

	log.Printf("Found the following iqns: %s", iqns)
	return iqns, nil
}

// Discovery ISCSI Target
func Discovery(portal string) error {
	log.Printf("Discovery portal: %s", portal)
	_, err := exec.Command("iscsiadm", "-m", "discovery", "-t", "sendtargets", "-p", portal).CombinedOutput()
	if err != nil {
		log.Fatalf("Error encountered in sendtargets: %v", err)
		return err
	}
	return nil
}

// Login ISCSI Target
func Login(portal string, targetiqn string) error {
	log.Printf("Login portal: %s targetiqn: %s", portal, targetiqn)
	_, err := exec.Command("iscsiadm", "-m", "node", "-p", portal, "-T", targetiqn, "--login").CombinedOutput()
	if err != nil {
		log.Fatalf("Received error on login attempt: %v", err)
		return err
	}
	return nil
}

// Logout ISCSI Target
func Logout(portal string, targetiqn string) error {
	log.Printf("Logout portal: %s targetiqn: %s", portal, targetiqn)
	_, err := exec.Command("iscsiadm", "-m", "node", "-p", portal, "-T", targetiqn, "--logout").CombinedOutput()
	if err != nil {
		log.Fatalf("Received error on logout attempt: %v", err)
		return err
	}
	return nil
}

// Delete ISCSI Node
func Delete(targetiqn string) (err error) {
	log.Printf("Delete targetiqn: %s", targetiqn)
	_, err = exec.Command("iscsiadm", "-m", "node", "-o", "delete", "-T", targetiqn).CombinedOutput()
	if err != nil {
		log.Fatalf("Received error on Delete attempt: %v", err)
		return err
	}
	return nil
}

// Connect ISCSI Target
func Connect(portal string, targetiqn string, targetlun string) (string, error) {
	log.Printf("Connect portal: %s targetiqn: %s targetlun: %s", portal, targetiqn, targetlun)
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

// Disconnect ISCSI Target
func Disconnect(portal string, targetiqn string) error {
	log.Printf("Disconnect portal: %s targetiqn: %s", portal, targetiqn)

	// Logout
	err := Logout(portal, targetiqn)
	if err != nil {
		return err
	}

	//Delete
	err = Delete(targetiqn)
	if err != nil {
		return err
	}

	return nil
}

// GetFSType returns the File System Type of device
func GetFSType(device string) string {
	log.Printf("GetFSType: %s", device)
	fsType := ""
	res, err := exec.Command("blkid", device).CombinedOutput()
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
		_, err := exec.Command("mkfs", "-t", fstype, "-F", device).CombinedOutput()
		if err != nil {
			log.Fatalf("failed to Format: %v", err)
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

	_, err := exec.Command("mkdir", "-p", mountpoint).CombinedOutput()
	if err != nil {
		log.Fatalf("failed to mkdir: %v", err)
	}
	_, err = exec.Command("mount", device, mountpoint).CombinedOutput()
	if err != nil {
		log.Fatalf("failed to mount: %v", err)
		return err
	}
	return nil
}

// FormatandMount device
func FormatandMount(device string, fstype string, mountpoint string) error {
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

	_, err := exec.Command("umount", mountpoint).CombinedOutput()
	if err != nil {
		log.Fatalf("failed to Umount: %v", err)
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
