// Copyright (c) 2017 Huawei Technologies Co., Ltd. All Rights Reserved.
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
	"github.com/opensds/opensds/pkg/utils"
	"os"
)

type ISCSITarget interface {
	CreateISCSITarget(volId, tgtIqn, path, initiator string, chapAuth []string) error
	GetISCSITarget(iqn string) int
	RemoveISCSITarget(volId, iqn string) error

	AddLun(tid, lun int, path string) error
	GetLun(path string) int
	RemoveLun(tid, lun int) error

	BindInitiatorName(tid, initiator string) error
	UnbindInitiatorName(tid, initiator string) error

	BindInitiatorAddress(tid, initiator string) error
	UnbindInitiatorAddress(tid, initiator string) error
}

func NewISCSITarget(bip, tgtConfDir string) ISCSITarget {
	return &tgtTarget{
		TgtConfDir: tgtConfDir,
		BindIp:     bip,
	}
}

type tgtTarget struct {
	BindIp     string
	TgtConfDir string
}

func (t *tgtTarget) AddLun(tid, lun int, path string) error {
	var cmd = []string{
		"--lld", "iscsi",
		"--op", "new",
		"--mode", "logicalunit",
		"--tid", fmt.Sprint(tid),
		"--lun", fmt.Sprint(lun),
		"--backing-store", path,
	}
	if _, err := t.execCmd("tgtadm", cmd...); err != nil {
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
	out, err := t.execCmd("tgtadm", cmd...)
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

func (t *tgtTarget) RemoveLun(tid, lun int) error {
	var cmd = []string{
		"--lld", "iscsi",
		"--op", "delete",
		"--mode", "logicalunit",
		"--tid", fmt.Sprint(tid),
		"--lun", fmt.Sprint(lun),
	}
	if _, err := t.execCmd("tgtadm", cmd...); err != nil {
		log.Error("Fail to exec 'tgtadm' to remove lun from iscsi target:", err)
		return err
	}

	return nil
}

func (t *tgtTarget) CreateISCSITarget(volId, tgtIqn, path, initiator string, chapAuth []string) error {

	if exist, _ := utils.PathExists(t.TgtConfDir); !exist {
		os.MkdirAll(t.TgtConfDir, 0755)
	}

	var charStr string
	if len(chapAuth) != 0 {
		charStr = fmt.Sprintf("incominguser %s %s", chapAuth[0], chapAuth[1])
	}

	tgtConfPath := t.TgtConfDir + "/" + volId + ".conf"

	var tgtConfFormatter = `
<target %s>
	backing-store %s
	driver %s
	%s
	%s
	%s
	write-cache %s
</target>
`
	var initiatorAddr string
	var initiatorName string
	if initiator != "ALL" {
		initiatorAddr = "initiator-address ALL"
		initiatorName = "initiator-name " + initiator
	}

	confStr := fmt.Sprintf(tgtConfFormatter, tgtIqn, path, "iscsi", charStr, initiatorAddr, initiatorName, "on")
	f, err := os.Create(tgtConfPath)
	if err != nil {
		return err
	}
	defer f.Close()
	f.WriteString(confStr)
	f.Sync()

	if info, err := t.execCmd("tgt-admin", "--update", tgtIqn); err != nil {
		log.Error("Fail to exec 'tgtadm' to create iscsi target, ", info, err)
		return err
	}

	if t.GetISCSITarget(tgtIqn) == -1 {
		log.Errorf("Failed to create iscsi target for Volume "+
			"ID: %s. It could be caused by problem "+
			"with concurrency. "+
			"Also please ensure your tgtd config "+
			"file contains 'include %s/*'",
			volId, t.TgtConfDir)
		return fmt.Errorf("failed to create volume(%s) attachment", volId)
	}
	return nil
}

func (t *tgtTarget) GetISCSITarget(iqn string) int {
	var cmd = []string{
		"--lld", "iscsi",
		"--op", "show",
		"--mode", "target",
	}
	out, err := t.execCmd("tgtadm", cmd...)
	if err != nil {
		log.Error("Fail to exec 'tgtadm' to display iscsi target:", err)
		return -1
	}

	var tid = -1
	for _, line := range strings.Split(out, "\n") {
		if strings.Contains(line, iqn) {
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

func (t *tgtTarget) RemoveISCSITarget(volId, iqn string) error {
	tgtConfPath := t.TgtConfDir + "/" + volId + ".conf"
	if exist, _ := utils.PathExists(tgtConfPath); !exist {
		log.Warningf("Volume path %s does not exist, nothing to remove.", tgtConfPath)
		return nil
	}

	if info, err := t.execCmd("tgt-admin", "--force", "--delete", iqn); err != nil {
		log.Error("Fail to exec 'tgtadm' to forcely remove iscsi target", info, err)
		return err
	}

	os.Remove(tgtConfPath)
	return nil
}
func (t *tgtTarget) isBindInitiatorName(tid, initiator string) bool {
	var cmd = []string{
		"--lld", "iscsi",
		"--op", "show",
		"--mode", "target",
	}
	out, err := t.execCmd("tgtadm", cmd...)
	if err != nil {
		log.Error("Fail to exec 'tgtadm' to display iscsi target:", err)
		return false
	}

	indent := 0
	tgtPrefix := fmt.Sprintf("Target %d", tid)
	var lines = strings.Split(out, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, tgtPrefix) {
			indent++
			continue
		}
		if indent == 1 && strings.HasPrefix(line, "    ACL information:") {
			indent++
			continue
		}
		if indent == 2 {
			// ACL bind ip indent 8 spaces
			if !strings.HasPrefix(line, "        ") {
				return false
			}
			if strings.Contains(line, initiator) {
				return true
			}
		}
	}
	return false
}

func (t *tgtTarget) BindInitiatorName(tid, initiator string) error {
	if t.isBindInitiatorName(tid, initiator) {
		log.Infof("Specified initiator %s has already been binded to target %d", initiator, tid)
		return nil
	}
	var cmd = []string{
		"--lld", "iscsi",
		"--op", "bind",
		"--mode", "target",
		"--tid", fmt.Sprint(tid),
		"--initiator-name", initiator,
	}
	if _, err := t.execCmd("tgtadm", cmd...); err != nil {
		log.Error("Fail to exec 'tgtadm' to bind iscsi target:", err)
		return err
	}

	return nil
}

func (t *tgtTarget) UnbindInitiatorName(tid, initiator string) error {
	var cmd = []string{
		"--lld", "iscsi",
		"--op", "unbind",
		"--mode", "target",
		"--tid", fmt.Sprint(tid),
		"--initiator-name", initiator,
	}
	if _, err := t.execCmd("tgtadm", cmd...); err != nil {
		log.Error("Fail to exec 'tgtadm' to unbind iscsi target:", err)
		return err
	}

	return nil
}

func (t *tgtTarget) BindInitiatorAddress(tid, initiator string) error {
	var cmd = []string{
		"--lld", "iscsi",
		"--op", "bind",
		"--mode", "target",
		"--tid", fmt.Sprint(tid),
		"--initiator-address", initiator,
	}
	if _, err := t.execCmd("tgtadm", cmd...); err != nil {
		log.Error("Fail to exec 'tgtadm' to bind iscsi target:", err)
		return err
	}

	return nil
}

func (t *tgtTarget) UnbindInitiatorAddress(tid, initiator string) error {
	var cmd = []string{
		"--lld", "iscsi",
		"--op", "unbind",
		"--mode", "target",
		"--tid", fmt.Sprint(tid),
		"--initiator-address", initiator,
	}
	if _, err := t.execCmd("tgtadm", cmd...); err != nil {
		log.Error("Fail to exec 'tgtadm' to unbind iscsi target:", err)
		return err
	}

	return nil
}

func (*tgtTarget) execCmd(name string, cmd ...string) (string, error) {
	ret, err := exec.Command(name, cmd...).Output()
	log.V(8).Infoln("Command: ", "tgtadm ", strings.Join(cmd, " "))
	log.V(8).Infof("result:%s", string(ret))
	if err != nil {
		log.V(8).Info("error info:", err)
	}
	return string(ret), err
}
