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
	"fmt"
	"log"
	"net"
	"os/exec"
	"strings"
)

// ExecCmd Log and convert the result of exec.Command
func ExecCmd(name string, arg ...string) (string, error) {
	log.Printf("Command: %s %s:\n", name, strings.Join(arg, " "))
	info, err := exec.Command(name, arg...).CombinedOutput()
	return string(info), err
}

// GetFSType returns the File System Type of device
func GetFSType(device string) (string, error) {
	log.Printf("GetFSType: %s\n", device)

	var fsType string
	blkidCmd := "blkid"
	out, err := ExecCmd("blkid", device)
	if err != nil {
		log.Printf("failed to GetFSType: %v cmd: %s output: %s\n",
			err, blkidCmd, string(out))
		return fsType, nil
	}

	for _, v := range strings.Split(string(out), " ") {
		if strings.Contains(v, "TYPE=") {
			fsType = strings.Split(v, "=")[1]
			fsType = strings.Replace(fsType, "\"", "", -1)
			fsType = strings.Replace(fsType, "\n", "", -1)
			fsType = strings.Replace(fsType, "\r", "", -1)
			return fsType, nil
		}
	}

	return fsType, nil
}

// Format device by File System Type
func Format(device string, fsType string) error {
	log.Printf("Format device: %s fstype: %s\n", device, fsType)

	mkfsCmd := fmt.Sprintf("mkfs.%s", fsType)

	_, err := exec.LookPath(mkfsCmd)
	if err != nil {
		if err == exec.ErrNotFound {
			return fmt.Errorf("%q executable not found in $PATH", mkfsCmd)
		}
		return err
	}

	mkfsArgs := []string{}
	mkfsArgs = append(mkfsArgs, device)
	if fsType == "ext4" || fsType == "ext3" {
		mkfsArgs = []string{"-F", device}
	}

	out, err := ExecCmd(mkfsCmd, mkfsArgs...)
	if err != nil {
		return fmt.Errorf("formatting disk failed: %v cmd: '%s %s' output: %q",
			err, mkfsCmd, strings.Join(mkfsArgs, " "), string(out))
	}

	return nil
}

// Mount device into mount point
func Mount(device, mountpoint, fsType string, mountFlags []string) error {
	log.Printf("Mount device: %s mountpoint: %s, fsType: %s, mountFlags:　%v\n", device, mountpoint, fsType, mountFlags)

	_, err := ExecCmd("mkdir", "-p", mountpoint)
	if err != nil {
		log.Printf("failed to mkdir: %v\n", err)
		return err
	}

	mountArgs := []string{}

	mountArgs = append(mountArgs, "-t", fsType)

	if len(mountFlags) > 0 {
		mountArgs = append(mountArgs, "-o", strings.Join(mountFlags, ","))
	}

	mountArgs = append(mountArgs, device)
	mountArgs = append(mountArgs, mountpoint)

	_, err = exec.Command("mount", mountArgs...).CombinedOutput()
	if err != nil {
		log.Printf("failed to mount: %v\n", err)
		return err
	}
	return nil
}

// Umount from mountpoint
func Umount(mountpoint string) error {
	log.Printf("Umount mountpoint: %s\n", mountpoint)

	_, err := ExecCmd("umount", mountpoint)
	if err != nil {
		log.Printf("failed to Umount: %v\n", err)
		return err
	}
	return nil
}

// GetHostIP return Host IP
func GetHostIP() string {
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

// GetHostName ...
func GetHostName() (string, error) {
	hostName, err := ExecCmd("hostname")
	if err != nil {
		log.Printf("failed to get host name: %v\n", err)
		return "", err
	}
	hostName = strings.Replace(hostName, "\n", "", -1)
	return hostName, nil
}

// IsMounted ...
func IsMounted(target string) (bool, error) {
	findmntCmd := "findmnt"
	_, err := exec.LookPath(findmntCmd)
	if err != nil {
		if err == exec.ErrNotFound {
			msg := fmt.Sprintf("%s executable not found in $PATH, err: %v\n", findmntCmd, err)
			log.Printf(msg)
			return false, fmt.Errorf(msg)
		}
		log.Printf("failed to check IsMounted %v\n", err)
		return false, err
	}

	findmntArgs := []string{"--target", target}

	log.Printf("findmnt args is %s\n", findmntArgs)

	out, err := ExecCmd(findmntCmd, findmntArgs...)
	if err != nil {
		// findmnt exits with non zero exit status if it couldn't find anything
		if strings.TrimSpace(string(out)) == "" {
			return false, nil
		}

		errIsMounted := fmt.Errorf("checking mounted failed: %v cmd: %s output: %s",
			err, findmntCmd, string(out))

		log.Printf("checking mounted failed: %v\n", errIsMounted)
		return false, errIsMounted
	}

	log.Printf("checking mounted result is %s\n", strings.TrimSpace(string(out)))
	if strings.TrimSpace(string(out)) == "" {
		return false, nil
	}

	line := strings.Split(string(out), "\n")

	if strings.Split(line[1], " ")[0] != target {
		return false, nil
	}

	return true, nil
}
