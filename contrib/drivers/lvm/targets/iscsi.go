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
	"os"
	"os/exec"
	"strconv"
	"strings"

	log "github.com/golang/glog"
	"github.com/opensds/opensds/pkg/utils"
)

const (
	opensdsPrefix = "opensds-"
	tgtAdminCmd   = "tgt-admin"
)

type ISCSITarget interface {
	CreateISCSITarget(volId, tgtIqn, path, hostIp, initiator string, chapAuth []string) error
	GetISCSITarget(iqn string) int
	RemoveISCSITarget(volId, iqn string) error
	GetLun(path string) int
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

func (t *tgtTarget) GetLun(path string) int {
	out, err := t.execCmd(tgtAdminCmd, "--show")
	if err != nil {
		log.Errorf("Fail to exec '%s' to display iscsi target:%v", tgtAdminCmd, err)
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

func (t *tgtTarget) getTgtConfPath(volId string) string {
	return t.TgtConfDir + "/" + opensdsPrefix + volId + ".conf"
}

func (t *tgtTarget) CreateISCSITarget(volId, tgtIqn, path, hostIp, initiator string, chapAuth []string) error {

	if exist, _ := utils.PathExists(t.TgtConfDir); !exist {
		os.MkdirAll(t.TgtConfDir, 0755)
	}

	var charStr string
	if len(chapAuth) != 0 {
		charStr = fmt.Sprintf("incominguser %s %s", chapAuth[0], chapAuth[1])
	}

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
	var initiatorAddr = "initiator-address" + hostIp
	var initiatorName string
	if initiator != "ALL" {
		initiatorName = "initiator-name " + initiator
	}

	confStr := fmt.Sprintf(tgtConfFormatter, tgtIqn, path, "iscsi", charStr, initiatorAddr, initiatorName, "on")
	f, err := os.Create(t.getTgtConfPath(volId))
	if err != nil {
		return err
	}
	defer f.Close()
	f.WriteString(confStr)
	f.Sync()

	if info, err := t.execCmd(tgtAdminCmd, "--update", tgtIqn); err != nil {
		log.Errorf("Fail to exec '%s' to create iscsi target, %s,%v", tgtAdminCmd, string(info), err)
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
	out, err := t.execCmd(tgtAdminCmd, "--show")
	if err != nil {
		log.Errorf("Fail to exec '%s' to display iscsi target:%v", tgtAdminCmd, err)
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
	tgtConfPath := t.getTgtConfPath(volId)
	if exist, _ := utils.PathExists(tgtConfPath); !exist {
		log.Warningf("Volume path %s does not exist, nothing to remove.", tgtConfPath)
		return nil
	}

	if info, err := t.execCmd(tgtAdminCmd, "--force", "--delete", iqn); err != nil {
		log.Errorf("Fail to exec '%s' to forcely remove iscsi target, %s, %v",
			tgtAdminCmd, string(info), err)
		return err
	}

	os.Remove(tgtConfPath)
	return nil
}

func (*tgtTarget) execCmd(name string, cmd ...string) (string, error) {
	ret, err := exec.Command(name, cmd...).Output()
	log.V(8).Infoln("Command:", cmd, strings.Join(cmd, " "))
	log.V(8).Infof("result:%s", string(ret))
	if err != nil {
		log.V(8).Info("error info:", err)
	}
	return string(ret), err
}
