// Copyright 2019 The OpenSDS Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package targets

import (
	"bytes"
	"errors"
	"io"
	"os"
	"os/exec"
	"strings"

	log "github.com/golang/glog"
	"github.com/opensds/opensds/pkg/utils"
)

const (
	opensdsNvmeofPrefix = "opensds-Nvmeof"
	NvmetDir            = "/sys/kernel/config/nvmet"
)

type NvmeofTarget interface {
	AddNvmeofSubsystem(volId, tgtNqn, path, initiator string) (string, error)
	RemoveNvmeofSubsystem(volId, nqn string) error
	GetNvmeofSubsystem(nqn string) (string, error)
	CreateNvmeofTarget(volId, tgtIqn, path, initiator, transtype string) error
	GetNvmeofTarget(nqn, transtype string) (bool, error)
	RemoveNvmeofTarget(volId, nqn, transtype string) error
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
	t.execCmd("modprobe", "nvmet-tcp")
	t.execCmd("modprobe", "nvmet-fc")
}

func (t *NvmeoftgtTarget) getTgtConfPath(volId string) string {
	return NvmetDir + "/" + opensdsNvmeofPrefix + volId
}

func (t *NvmeoftgtTarget) convertTranstype(transtype string) string {
	var portid string
	switch transtype {
	case "fc":
		portid = "3"
	case "rdma":
		portid = "2"
	default:
		portid = "1"
		log.Infof("default nvmeof transtype : tcp")
	}
	return portid
}

func (t *NvmeoftgtTarget) AddNvmeofSubsystem(volId, tgtNqn, path, initiator string) (string, error) {
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
			t.RemoveNvmeofSubsystem(volId, tgtNqn)
			return "", err
		}
	} else {
		// allow specific initiators to connect to this target
		var initiatorInfo = initiator
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
			t.RemoveNvmeofSubsystem(volId, tgtNqn)
			return "", err
		}
	}

	// get volume namespaceid
	namespaceid := t.Getnamespaceid(volId)
	if namespaceid == "" {
		t.RemoveNvmeofSubsystem(volId, tgtNqn)
		return "", errors.New("null namesapce")
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
		t.RemoveNvmeofSubsystem(volId, tgtNqn)
		return "", err
	}

	enablepath := namespace + "/enable"
	err = t.WriteWithIo(enablepath, "1")
	if err != nil {
		log.Errorf("Fail to set device path")
		t.RemoveNvmeofSubsystem(volId, tgtNqn)
		return "", err
	}
	log.Infof("new added subsys : %s", sysdir)
	return sysdir, nil
}

func (t *NvmeoftgtTarget) GetNvmeofSubsystem(nqn string) (string, error) {
	subsysdir := NvmetDir + "/subsystems/" + nqn
	if _, err := os.Stat(subsysdir); err == nil {
		return subsysdir, nil

	} else if os.IsNotExist(err) {
		return "", nil

	} else {
		log.Errorf("can not get nvmeof subsystem")
		return "", err
	}

}

func (t *NvmeoftgtTarget) CreateNvmeofTarget(volId, tgtNqn, path, initiator, transtype string) error {

	if tgtexisted, err := t.GetNvmeofTarget(tgtNqn, transtype); tgtexisted == true && err == nil {
		log.Infof("Nvmeof target %s with transtype %s has existed", tgtNqn, transtype)
		return nil
	} else if err != nil {
		log.Errorf("can not get nvmeof target %s with transport type %s", tgtNqn, transtype)
		return err
	}

	var subexisted string
	subexisted, err := t.GetNvmeofSubsystem(tgtNqn)
	if err != nil {
		log.Errorf("can not get nvmeof subsystem %s ", tgtNqn)
		return err
	} else if subexisted == "" {
		log.Infof("add new nqn subsystem %s", tgtNqn)
		subexisted, err = t.AddNvmeofSubsystem(volId, tgtNqn, path, initiator)
		log.Infof("new subdir: %s", subexisted)
	} else {
		log.Infof("%s subsystem has existed", tgtNqn)
	}

	subexisted = NvmetDir + "/subsystems/" + tgtNqn
	log.Infof("new subdir: %s", subexisted)
	//	subexisted, err = t.GetNvmeofSubsystem(tgtNqn)
	//	log.Infof("new subdir: %s ", subexisted)
	//	if  subexisted == "" {
	//	log.Infof("still no subsystem after add new subsystem")
	//		//t.RemoveNvmeofSubsystem(volId, tgtNqn)
	//		return errors.New("still can not get subsystem after add new one")
	//	}
	//
	//create port
	portid := t.convertTranstype(transtype)
	portspath := NvmetDir + "/ports/" + portid
	if exist, _ := utils.PathExists(portspath); !exist {
		//log.Errorf(portspath)
		os.MkdirAll(portspath, 0755)
	}

	// get target ip
	// here the ip should be the ip interface of the specific nic
	// for example, if transport type is rdma, then the rdma ip should be used.
	// here just set the generic ip address since tcp is the default choice.
	ippath := portspath + "/addr_traddr"
	ip, err := t.execCmd("hostname", "-I")
	if err != nil {
		log.Errorf("fail to get target ipv4 address")
		t.RemoveNvmeofTarget(volId, tgtNqn, transtype)
		return err
	}

	ip = strings.Split(ip, " ")[0]
	err = t.WriteWithIo(ippath, ip)
	if err != nil {
		log.Errorf("Fail to set target ip")
		t.RemoveNvmeofTarget(volId, tgtNqn, transtype)
		return err
	}

	trtypepath := portspath + "/addr_trtype"
	err = t.WriteWithIo(trtypepath, transtype)
	if err != nil {
		log.Errorf("Fail to set transport type")
		t.RemoveNvmeofTarget(volId, tgtNqn, transtype)
		return err
	}

	trsvcidpath := portspath + "/addr_trsvcid"
	err = t.WriteWithIo(trsvcidpath, "4420")
	if err != nil {
		log.Errorf("Fail to set ip port")
		t.RemoveNvmeofTarget(volId, tgtNqn, transtype)
		return err
	}

	adrfampath := portspath + "/addr_adrfam"
	err = t.WriteWithIo(adrfampath, "ipv4")
	if err != nil {
		log.Errorf("Fail to set ip family")
		t.RemoveNvmeofTarget(volId, tgtNqn, transtype)
		return err
	}

	// create a soft link
	portssub := portspath + "/subsystems/" + tgtNqn
	_, err = t.execCmd("ln", "-s", subexisted, portssub)
	if err != nil {
		log.Errorf("Fail to create link")
		t.RemoveNvmeofTarget(volId, tgtNqn, transtype)
		return err
	}

	// check
	info, err := t.execCmd("dmesg", `|grep "enabling port"`)
	if err != nil || info == "" {
		log.Errorf("nvme target is not listening on the port")
		t.RemoveNvmeofTarget(volId, tgtNqn, transtype)
		return err
	}
	log.Info("create nvme target")
	return nil
}

func (t *NvmeoftgtTarget) GetNvmeofTarget(nqn, transtype string) (bool, error) {
	portid := t.convertTranstype(transtype)

	targetlinkpath := NvmetDir + "/ports/" + portid + "/subsystems/" + nqn
	if _, err := os.Lstat(targetlinkpath); err == nil {
		return true, nil

	} else if os.IsNotExist(err) {
		return false, nil

	} else {
		log.Errorf("can not get nvmeof target")
		return false, err
	}

}

func (t *NvmeoftgtTarget) RemoveNvmeofSubsystem(volId, nqn string) error {
	log.Info("removing subsystem", nqn)
	tgtConfPath := NvmetDir + "/subsystems/" + nqn
	if exist, _ := utils.PathExists(tgtConfPath); !exist {
		log.Warningf("Volume path %s does not exist, nothing to remove.", tgtConfPath)
		return nil
	}

	// remove namespace， whether it succeed or not, the removement should be executed.
	ns := t.Getnamespaceid(volId)
	if ns == "" {
		log.Infof("can not find volume %s's namespace", volId)
		// return errors.New("null namespace")
	}
	naspPath := NvmetDir + "/subsystems/" + nqn + "/namespaces/" + ns
	info, err := t.execCmd("rmdir", naspPath)
	if err != nil {
		log.Infof("can not rm nasp")
		// return err
	}

	// remove namespaces ; if it allows all initiators ,then this dir should be empty
	// if it allow specific hosts ,then here remove all the hosts
	cmd := "rm -f " + NvmetDir + "/subsystems/" + nqn + "/allowed_hosts/" + "*"
	info, err = t.execBash(cmd)
	if err != nil {
		log.Infof("can not rm allowed hosts")
		log.Infof(info)
		// return err
	}

	// remove subsystem
	syspath := NvmetDir + "/subsystems/" + nqn
	info, err = t.execCmd("rmdir", syspath)
	if err != nil {
		log.Infof("can not rm subsys")
		return err
	}
	return nil
}

func (t *NvmeoftgtTarget) RemoveNvmeofPort(nqn, transtype string) error {
	log.Infof("removing nvmeof port %s", transtype)
	portid := t.convertTranstype(transtype)

	portpath := NvmetDir + "/ports/" + portid + "/subsystems/" + nqn

	//  port's link has to be removed first or the subsystem cannot be removed
	tgtConfPath := NvmetDir + "/subsystems/" + nqn
	if exist, _ := utils.PathExists(tgtConfPath); !exist {
		log.Warningf("Volume path %s does not exist, nothing to remove.", tgtConfPath)
		return nil
	}

	info, err := t.execCmd("rm", "-f", portpath)
	if err != nil {
		log.Errorf("can not rm nvme port transtype: %s, nqn: %s", transtype, nqn)
		log.Errorf(info)
		return err
	}
	return nil
}

func (t *NvmeoftgtTarget) RemoveNvmeofTarget(volId, nqn, transtype string) error {
	log.Infof("removing target %s", nqn)
	if tgtexisted, err := t.GetNvmeofTarget(nqn, transtype); err != nil {
		log.Errorf("can not get nvmeof target %s with type %s", nqn, transtype)
		return err
	} else if tgtexisted == false {
		log.Infof("nvmeof target %s with type %s does not exist", nqn, transtype)
	} else {
		err = t.RemoveNvmeofPort(nqn, transtype)
		if err != nil {
			return err
		}
	}

	if subexisted, err := t.GetNvmeofSubsystem(nqn); err != nil {
		log.Errorf("can not get nvmeof subsystem %s ", nqn)
		return err
	} else if subexisted == "" {
		log.Errorf("subsystem %s does not exist", nqn)
		return nil
	} else {
		err = t.RemoveNvmeofSubsystem(volId, nqn)
		if err != nil {
			log.Errorf("can not remove nvme subsystem %s", nqn)
			return err
		}
	}
	return nil
}

func (*NvmeoftgtTarget) execCmd(name string, cmd ...string) (string, error) {
	ret, err := exec.Command(name, cmd...).Output()
	if err != nil {
		log.Errorf("error info: %v", err)
	}
	return string(ret), err
}

func (*NvmeoftgtTarget) execBash(name string) (string, error) {
	ret, err := exec.Command("/bin/sh", "-c", name).Output()
	if err != nil {
		log.Errorf("error info in sh %v ", err)
	}
	return string(ret), err
}

func (*NvmeoftgtTarget) WriteWithIo(name, content string) error {
	fileObj, err := os.OpenFile(name, os.O_RDWR, 0644)
	if err != nil {
		log.Errorf("Failed to open the file %v", err)
		return err
	}
	if _, err := io.WriteString(fileObj, content); err == nil {
		log.Infof("Successful appending to the file with os.OpenFile and io.WriteString.%s", content)
		return nil
	}
	return err
}

func (t *NvmeoftgtTarget) Getnamespaceid(volId string) string {
	var buffer bytes.Buffer
	for _, rune := range volId {
		// nvme target namespace dir should not be like 00 or 0 ,
		// so only digits range from 1 to 9 are accepted
		if rune >= '1' && rune <= '9' {
			buffer.WriteRune(rune)
		}
	}
	return buffer.String()[0:2]
}
