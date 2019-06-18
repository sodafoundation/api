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

package scms

import (
	"os/exec"
	"strconv"
	"errors"
)

const CMS_ADM = "/opt/cmsagent/cmsadm"

const (
	CMS_CREATE = "--create"
	CMS_DELETE = "--remove"
	CMS_START = "--up"
	CMS_STOP = "--down"
	CMS_QUERY = "--query"
)

type CmsAdm struct {
}


func NewCmsAdm () *CmsAdm {
	return &CmsAdm{}
}

func (c *CmsAdm) CreateTask(t *CmsTask, arg ...string) ([]byte, error) {
	var argv = []string {CMS_CREATE}

	var option string
	option = "-b " + strconv.FormatInt(t.bandwidth, 10)
	argv = append(argv, option)

	if t.cdpFlag {
		option = "-j"
		argv = append(argv, option)
	}

	for svolId := range t.taskVolumes {
		tvolId := t.taskVolumes[svolId]
		svol := t.volumeList[svolId]
		tvol := t.volumeList[tvolId]

		option = ("-D " + svol.VolumeName + "," + tvol.VolumeName)
		argv = append(argv, option)
	}

	return cmdExec(CMS_ADM, argv)
}

func (c *CmsAdm) DeleteTask(arg ...string) ([]byte, error) {
	var argv = []string {CMS_DELETE}
	return cmdExec(CMS_ADM, argv)
}

func (c *CmsAdm) Up() ([]byte, error) {
	var argv = []string {CMS_START}
	return cmdExec(CMS_ADM, argv)
}

func (c *CmsAdm) Down() ([]byte, error) {
	var argv = []string{CMS_STOP}
	return cmdExec(CMS_ADM, argv)
}

func (c *CmsAdm) Query() ([]byte, error) {
	var argv = []string{CMS_QUERY}
	return cmdExec(CMS_ADM, argv)
}

func cmdExec(cmd string, argv[]string) ([]byte, error) {
	out, err := exec.Command(cmd, argv[0:]...).Output()
	if err != nil {
		err = errors.New(string(out))
	}
	return out,err
}
