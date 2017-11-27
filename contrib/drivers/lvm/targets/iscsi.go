// Copyright (c) 2017 OpenSDS Authors.
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
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	log "github.com/golang/glog"
)

type ISCSITarget interface {
	CreateISCSITarget() error
	GetISCSITarget() int
	RemoveISCSITarget() error

	AddLun(lun int, path string) error
	GetLun(path string) int
	RemoveLun(lun int) error

	BindInitiator(initiator string) error
	UnbindInitiator(initiator string) error
}

func NewISCSITarget(tid int, name string, bip string) ISCSITarget {
	return &tgtTarget{
		Tid:    tid,
		TName:  name,
		BindIp: bip,
	}
}

type tgtTarget struct {
	Tid    int
	TName  string
	BindIp string
}

func (t *tgtTarget) AddLun(lun int, path string) error {
	var cmd = []string{
		"--lld", "iscsi",
		"--op", "new",
		"--mode", "logicalunit",
		"--tid", fmt.Sprint(t.Tid),
		"--lun", fmt.Sprint(lun),
		"--backing-store", path,
	}
	if _, err := t.execCmd(cmd); err != nil {
		log.Error("Fail to exec 'tgtadm' to add lun into iscsi target:", err)
		return err
	}

	return nil
}

func (t *tgtTarget) GetLun(path string) int {
	var cmd = []string{
		"--lld", "iscsi",
		"--op", "show",
		"--mode", "target",
	}
	out, err := t.execCmd(cmd)
	if err != nil {
		log.Error("Fail to exec 'tgtadm' to display iscsi target:", err)
		return -1
	}

	var lun = -1
	var lines = strings.Split(out, "\n")
	for num, line := range lines {
		if strings.Contains(line, path) {
			for i := 1; i < num; i++ {
				if strings.Contains(lines[num-i], "LUN") {
					lunString := strings.Fields(lines[num-i])[1]
					lun, err = strconv.Atoi(lunString)
					if err != nil {
						return -1
					}
					return lun
				}
			}
		}
	}
	log.Info("Got lun id:", lun)

	return -1
}

func (t *tgtTarget) RemoveLun(lun int) error {
	var cmd = []string{
		"--lld", "iscsi",
		"--op", "delete",
		"--mode", "logicalunit",
		"--tid", fmt.Sprint(t.Tid),
		"--lun", fmt.Sprint(lun),
	}
	if _, err := t.execCmd(cmd); err != nil {
		log.Error("Fail to exec 'tgtadm' to remove lun from iscsi target:", err)
		return err
	}

	return nil
}

func (t *tgtTarget) CreateISCSITarget() error {
	var cmd = []string{
		"--lld", "iscsi",
		"--op", "new",
		"--mode", "target",
		"--tid", fmt.Sprint(t.Tid),
		"-T", t.TName,
	}
	if _, err := t.execCmd(cmd); err != nil {
		log.Error("Fail to exec 'tgtadm' to create iscsi target:", err)
		return err
	}

	return nil
}

func (t *tgtTarget) GetISCSITarget() int {
	var cmd = []string{
		"--lld", "iscsi",
		"--op", "show",
		"--mode", "target",
	}
	out, err := t.execCmd(cmd)
	if err != nil {
		log.Error("Fail to exec 'tgtadm' to display iscsi target:", err)
		return -1
	}

	var tid = -1
	for _, line := range strings.Split(out, "\n") {
		if strings.Contains(line, t.TName) {
			tidString := strings.Fields(strings.Split(line, ":")[0])[1]
			tid, err = strconv.Atoi(tidString)
			if err != nil {
				return -1
			}
			break
		}
	}
	return tid
}

func (t *tgtTarget) RemoveISCSITarget() error {
	var cmd = []string{
		"--lld", "iscsi",
		"--op", "delete",
		"--force",
		"--mode", "target",
		"--tid", fmt.Sprint(t.Tid),
	}
	if _, err := t.execCmd(cmd); err != nil {
		log.Error("Fail to exec 'tgtadm' to forcely remove iscsi target:", err)
		return err
	}

	return nil
}

func (t *tgtTarget) BindInitiator(initiator string) error {
	var cmd = []string{
		"--lld", "iscsi",
		"--op", "bind",
		"--mode", "target",
		"--tid", fmt.Sprint(t.Tid),
		"-I", initiator,
	}
	if _, err := t.execCmd(cmd); err != nil {
		log.Error("Fail to exec 'tgtadm' to bind iscsi target:", err)
		return err
	}

	return nil
}

func (t *tgtTarget) UnbindInitiator(initiator string) error {
	var cmd = []string{
		"--lld", "iscsi",
		"--op", "unbind",
		"--mode", "target",
		"--tid", fmt.Sprint(t.Tid),
		"-I", initiator,
	}
	if _, err := t.execCmd(cmd); err != nil {
		log.Error("Fail to exec 'tgtadm' to unbind iscsi target:", err)
		return err
	}

	return nil
}

func (*tgtTarget) execCmd(cmd []string) (string, error) {
	ret, err := exec.Command("tgtadm", cmd...).Output()
	if err != nil {
		log.Error(err.Error())
		return "", err
	}
	return string(ret), nil
}

