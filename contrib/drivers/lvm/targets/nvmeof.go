// Copyright (c) 2019 Intel Corporation, Ltd. All Rights Reserved.
//
//    Licensed under the Apache License, Version 2.0 (the "License"); you may
//    not use this file except in compliance with the License. You may obtain
//    a copy of the License at
//
//         http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//    WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//    License for the specific language governing permissions and limitations
//    under the License.

package targets

import (
	"bytes"
	"errors"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"

	log "github.com/golang/glog"
	"github.com/opensds/opensds/pkg/utils"
)

const (
	opensdsNvmeofPrefix = "opensds-Nvmeof"
	NvmetDir            = "/sys/kernel/config/nvmet"
)

type NvmeofTarget interface {
	CreateNvmeofTarget(volId, tgtIqn, path, hostIp, initiator string, chapAuth []string) error
	GetNvmeofTarget(iqn string) int
	RemoveNvmeofTarget(volId, iqn string) error
}

func NewNvmeofTarget(bip, tgtConfDir string) NvmeofTarget {
	return &NvmeoftgtTarget{
		TgtConfDir: tgtConfDir,
		BindIp:     bip,
	}
}

type NvmeoftgtTarget struct {
	BindIp     string
	TgtConfDir string
}

func (t *NvmeoftgtTarget) init() {
	t.execCmd("modprobe", "nvmet")
	t.execCmd("modprobe", "nvmet-rdma")
}

func (t *NvmeoftgtTarget) getTgtConfPath(volId string) string {
	return NvmetDir + "/" + opensdsNvmeofPrefix + volId
}

func (t *NvmeoftgtTarget) CreateNvmeofTarget(volId, tgtNqn, path, hostIp, initiator string, chapAuth []string) error {

	if exist, _ := utils.PathExists(NvmetDir); !exist {
		os.MkdirAll(NvmetDir, 0755)
	}
	sysdir := NvmetDir + "/subsystems/" + tgtNqn
	if exist, _ := utils.PathExists(sysdir); !exist {
		os.MkdirAll(sysdir, 0755)
	}

	var err error
	if initiator == "ALL" {
		// echo 1 > attr_allow_any_host
		attrfile := sysdir + "/attr_allow_any_host"
		content := "1"
		err = t.WriteWithIo(attrfile, content)
		if err != nil {
			log.Errorf("can not set attr_allow_any_host ")
			return err
		}
	} else {
		// allow specific initiators to connect to this target
		var initiatorInfo = "initiator:" + hostIp + ":" + initiator
		hostpath := NvmetDir + "/hosts"
		if exist, _ := utils.PathExists(hostpath); !exist {
			os.MkdirAll(hostpath, 0755)
		}

		hostDir := hostpath + "/" + initiatorInfo
		if exist, _ := utils.PathExists(hostDir); !exist {
			os.MkdirAll(hostDir, 0755)
		}

		// create symbolic link of host
		hostsys := sysdir + "/allowed_hosts/"
		_, err = t.execCmd("ln", "-s", hostDir, hostsys)
		if err != nil {
			log.Errorf("Fail to create host link: " + initiatorInfo)
			return err
		}
	}

	// get volume namespaceid
	namespaceid := t.Getnamespaceid(volId)
	if namespaceid == "" {
		return errors.New("null namesapce")
	}
	namespace := sysdir + "/namespaces/" + namespaceid
	if exist, _ := utils.PathExists(namespace); !exist {
		os.MkdirAll(namespace, 0755)
	}

	// volid as device path
	devpath := namespace + "/device_path"
	err = t.WriteWithIo(devpath, path)
	if err != nil {
		log.Errorf("Fail to set device path")
		return err
	}

	enablepath := namespace + "/enable"
	err = t.WriteWithIo(enablepath, "1")
	if err != nil {
		log.Errorf("Fail to set device path")
		return err
	}

	//create port
	portid := 1
	portspath := NvmetDir + "/ports/" + strconv.Itoa(portid)
	if exist, _ := utils.PathExists(portspath); !exist {
		//log.Errorf(portspath)
		os.MkdirAll(portspath, 0755)
	}

	// get target ip
	ippath := portspath + "/addr_traddr"
	ip, err := t.execCmd("hostname", "-I")
	if err != nil {
		log.Errorf("fail to get target ipv4 address")
		return err
	}
	// Set nvmeof parameters, rightnow only supports rdma
	// if built on virtual machine the return string ip may contain several ip addresses
	ip = strings.Split(ip, " ")[0]
	err = t.WriteWithIo(ippath, ip)
	if err != nil {
		log.Errorf("Fail to set target ip")
		return err
	}

	trtypepath := portspath + "/addr_trtype"
	err = t.WriteWithIo(trtypepath, "rdma")
	if err != nil {
		log.Errorf("Fail to set rdma type")
		return err
	}

	trsvcidpath := portspath + "/addr_trsvcid"
	err = t.WriteWithIo(trsvcidpath, "4420")
	if err != nil {
		log.Errorf("Fail to set ip port")
		return err
	}

	adrfampath := portspath + "/addr_adrfam"
	err = t.WriteWithIo(adrfampath, "ipv4")
	if err != nil {
		log.Errorf("Fail to set ip family")
		return err
	}

	// create a soft link
	portssub := portspath + "/subsystems/" + tgtNqn
	_, err = t.execCmd("ln", "-s", sysdir, portssub)
	if err != nil {
		log.Errorf("Fail to create link")
		return err
	}

	// check
	info, err := t.execCmd("dmesg", `|grep "enabling port"`)
	if err != nil || info == "" {
		log.Errorf("NVMe target is not listening on the port")
		return err
	}
	log.Info("create nvme target")
	return nil
}

func (t *NvmeoftgtTarget) GetNvmeofTarget(nqn string) int {
	_, err := t.execCmd("cd", "/sys/kernel/config/nvmet/subsystems")
	if err != nil {
		log.Errorf("Fail to exec to enter nvme target dir:%v", err)
		return -1
	}
	_, err = t.execCmd("cd", nqn)
	if err != nil {
		log.Errorf("Fail to exec to display nvme target :%v", err)
		return -1
	}
	return 0
}

func (t *NvmeoftgtTarget) RemoveNvmeofTarget(volId, nqn string) error {
	log.Info("removing target",nqn)
	tgtConfPath := NvmetDir + "/subsystems/" + nqn
	if exist, _ := utils.PathExists(tgtConfPath); !exist {
		log.Warningf("Volume path %s does not exist, nothing to remove.", tgtConfPath)
		return nil
	}

	//  port's link has to be removed first or the subsystem cannot be removed
	portpath := NvmetDir + "/ports/1/subsystems/" + nqn
	info, err := t.execCmd("rm", "-f", portpath)
	if err != nil {
		log.Errorf("can not rm port")
		log.Errorf(info)
		return err
	}

	// remove namespace
	ns := t.Getnamespaceid(volId)
	if ns == "" {
		log.Errorf("can not find volume ", volId, "'s uuid")
		return errors.New("null namespace")
	}
	naspPath := NvmetDir + "/subsystems/" + nqn + "/namespaces/" + ns
	info, err = t.execCmd("rmdir", naspPath)
	if err != nil {
		log.Errorf("can not rm nasp")
		return err
	}

	// remove namespaces ; if it allows all initiators ,then this dir should be empty
	// if it allow specific hosts ,then here remove all the hosts
	cmd := "rm -f " + NvmetDir + "/subsystems/" + nqn + "/allowed_hosts/" + "*"
	info, err = t.execBash(cmd)
	if err != nil {
		log.Errorf("can not rm allowed hosts")
		log.Errorf(info)
		return err
	}

	// remove subsystem
	syspath := NvmetDir + "/subsystems/" + nqn
	info, err = t.execCmd("rmdir", syspath)
	if err != nil {
		log.Errorf("can not rm subsys")
		return err
	}
	return nil
}

func (*NvmeoftgtTarget) execCmd(name string, cmd ...string) (string, error) {
	ret, err := exec.Command(name, cmd...).Output()
	if err != nil {
		log.Errorf("error info:", err)
	}
	return string(ret), err
}

func (*NvmeoftgtTarget) execBash(name string) (string, error) {
	ret, err := exec.Command("/bin/sh", "-c", name).Output()
	if err != nil {
		log.Errorf("error info in sh ", err)
	}
	return string(ret), err
}

func (*NvmeoftgtTarget) WriteWithIo(name, content string) error {
	fileObj, err := os.OpenFile(name, os.O_RDWR, 0644)
	if err != nil {
		log.Errorf("Failed to open the file", err.Error())
		return err
	}
	if _, err := io.WriteString(fileObj, content); err == nil {
		log.Infof("Successful appending to the file with os.OpenFile and io.WriteString.", content)
		return nil
	}
	return err
}

func (t *NvmeoftgtTarget) Getnamespaceid(volId string) string {
	var buffer bytes.Buffer
	for _, rune := range volId {
		if rune >= '0' && rune <= '9' {
			buffer.WriteRune(rune)
		}
	}
	return buffer.String()[0:2]
}
