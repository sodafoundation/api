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

package connector

import (
	"log"
	"net"
	"os/exec"
	"strings"
)

// execCmd ...
func ExecCmd(name string, arg ...string) (string, error) {
	log.Printf("Command: %s %s:\n", name, strings.Join(arg, " "))
	info, err := exec.Command(name, arg...).CombinedOutput()
	return string(info), err
}

// GetFSType returns the File System Type of device
func GetFSType(device string) string {
	fsType := ""
	res, err := ExecCmd("blkid", device)
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
		_, err := ExecCmd("mkfs", "-t", fstype, "-F", device)
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

	_, err := ExecCmd("mkdir", "-p", mountpoint)
	if err != nil {
		log.Printf("failed to mkdir: %v", err)
	}
	_, err = ExecCmd("mount", device, mountpoint)
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

	_, err := ExecCmd("umount", mountpoint)
	if err != nil {
		log.Printf("failed to Umount: %v", err)
		return err
	}
	return nil
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

// GetHostName return Host Name
func GetHostName() (string, error) {
	hostName, err := ExecCmd("hostname")
	if err != nil {
		log.Printf("failed to get host name: %v", err)
		return "", err
	}

	hostName = strings.Replace(hostName, "\n", "", -1)

	return hostName, nil
}
