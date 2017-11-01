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
	"os/exec"
	"strings"

	log "github.com/golang/glog"
)

type ISCSITarget interface {
	CreateISCSITarget(*ISCSITargetOpts) (string, error)

	RemoveISCSITarget(*ISCSITargetOpts) error
}

type ISCSITargetOpts struct {
	Name       string
	Tid        string
	Lun        string
	Path       string
	VolumeId   string
	VolumeName string
	Parameters map[string]interface{}
}

func NewISCSITarget() ISCSITarget {
	return &tgtTarget{}
}

type tgtTarget struct{}

func (*tgtTarget) CreateISCSITarget(opt *ISCSITargetOpts) (string, error) {
	return "", nil
}

func (*tgtTarget) RemoveISCSITarget(opt *ISCSITargetOpts) error {
	return nil
}

func (*tgtTarget) getISCSITarget(volumeId string) (string, error) {
	return "0", nil
}

func (*tgtTarget) getTarget(iqn string) (string, error) {
	out, err := execCmd("targetcli iscsi/ ls")
	if err != nil {
		log.Error("Fail to exec 'targetcli iscsi/ ls':", err)
		return "", err
	}

	for _, line := range strings.Split(out, "\n") {
		if strings.Contains(line, iqn) {
			return iqn, nil
		}
	}

	return "", nil

}

func execCmd(cmd string) (string, error) {
	ret, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		log.Error(err.Error())
		return "", err
	}
	return string(ret), nil
}
