// Copyright 2017 The OpenSDS Authors.
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
	"errors"
	"fmt"
	"io/ioutil"
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
	RemoveISCSITarget(volId, iqn, hostIp string) error
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

type configMap map[string][]string

func (t *tgtTarget) CreateISCSITarget(volId, tgtIqn, path, hostIp, initiator string, chapAuth []string) error {
	if hostIp == "" {
		return errors.New("create ISCSI target failed, host ip cannot be empty")
	}

	if exist, _ := utils.PathExists(t.TgtConfDir); !exist {
		os.MkdirAll(t.TgtConfDir, 0755)
	}

	config := make(configMap)

	configFile := t.getTgtConfPath(volId)

	if IsExist(configFile) {
		data, err := ioutil.ReadFile(configFile)
		if err != nil {
			return err
		}
		config.parse(string(data))
	}

	var charStr string
	if len(chapAuth) != 0 {
		charStr = fmt.Sprintf("%s %s", chapAuth[0], chapAuth[1])
		config.updateConfigmap("incominguser", charStr)
	}

	config.updateConfigmap("initiator-address", hostIp)
	config.updateConfigmap("driver", "iscsi")
	config.updateConfigmap("backing-store", path)
	config.updateConfigmap("write-cache", "on")

	err := config.writeConfig(configFile, tgtIqn)
	if err != nil {
		log.Errorf("failed to update config file %s %v", t.getTgtConfPath(volId), err)
		return err
	}

	if info, err := t.execCmd(tgtAdminCmd, "--force", "--update", tgtIqn); err != nil {
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

func (t *tgtTarget) RemoveISCSITarget(volId, iqn, hostIp string) error {
	if hostIp == "" {
		return errors.New("remove ISCSI target failed, host ip cannot be empty")
	}

	tgtConfPath := t.getTgtConfPath(volId)
	if exist, _ := utils.PathExists(tgtConfPath); !exist {
		log.Warningf("Volume path %s does not exist, nothing to remove.", tgtConfPath)
		return nil
	}

	config := make(configMap)

	data, err := ioutil.ReadFile(tgtConfPath)
	if err != nil {
		return err
	}

	config.parse(string(data))

	ips := config["initiator-address"]
	for i, v := range ips {
		if v == hostIp {
			ips = append(ips[:i], ips[i+1:]...)
			break
		}
	}
	config["initiator-address"] = ips
	if len(ips) == 0 {
		if info, err := t.execCmd(tgtAdminCmd, "--force", "--delete", iqn); err != nil {
			log.Errorf("Fail to exec '%s' to forcely remove iscsi target, %s, %v",
				tgtAdminCmd, string(info), err)
			return err
		}

		os.Remove(tgtConfPath)
	} else {
		err := config.writeConfig(t.getTgtConfPath(volId), iqn)
		if err != nil {
			log.Errorf("failed to update config file %s %v", t.getTgtConfPath(volId), err)
			return err
		}

		if info, err := t.execCmd(tgtAdminCmd, "--force", "--update", iqn); err != nil {
			log.Errorf("Fail to exec '%s' to create iscsi target, %s,%v", tgtAdminCmd, string(info), err)
			return err
		}

		if t.GetISCSITarget(iqn) == -1 {
			log.Errorf("Failed to create iscsi target for Volume "+
				"ID: %s. It could be caused by problem "+
				"with concurrency. "+
				"Also please ensure your tgtd config "+
				"file contains 'include %s/*'",
				volId, t.TgtConfDir)
			return fmt.Errorf("failed to create volume(%s) attachment", volId)
		}
	}

	return nil
}

func (*tgtTarget) execCmd(name string, cmd ...string) (string, error) {
	ret, err := exec.Command(name, cmd...).Output()
	log.Infoln("Command:", cmd, strings.Join(cmd, " "))
	log.V(8).Infof("result:%s", string(ret))
	if err != nil {
		log.Error("error info:", err)
	}
	return string(ret), err
}

func (m *configMap) parse(data string) {
	var lines = strings.Split(data, "\n")

	for _, line := range lines {
		for _, key := range []string{"backing-store", "driver", "initiator-address", "write-cache"} {
			if strings.Contains(line, key) {
				s := strings.TrimSpace(line)
				if (*m)[key] == nil {
					(*m)[key] = []string{strings.Split(s, " ")[1]}
				} else {
					(*m)[key] = append((*m)[key], strings.Split(s, " ")[1])
				}
			}
		}
	}
}

func (m *configMap) updateConfigmap(key, value string) {
	v := (*m)[key]
	if v == nil {
		(*m)[key] = []string{value}
	} else {
		if !utils.Contains(v, value) {
			v = append(v, value)
			(*m)[key] = v
		}
	}
}

func (m configMap) writeConfig(file, tgtIqn string) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()

	f.WriteString(fmt.Sprintf("<target %s>\n", tgtIqn))
	for k, v := range m {
		for _, vl := range v {
			f.WriteString(fmt.Sprintf("        %s %s\n", k, vl))
		}
	}
	f.WriteString("</target>")
	f.Sync()
	return nil
}

func IsExist(f string) bool {
	_, err := os.Stat(f)
	return err == nil || os.IsExist(err)
}
