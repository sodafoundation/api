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
	"fmt"
)

type CmsVolume struct {
	VolumeId string
	VolumeName string
}

type CmsTask struct {
	bandwidth int64
	cdpFlag bool
	taskVolumes map[string]string
	volumeList map[string]CmsVolume
}

func NewCmsTask(bandwidth int64, cdpFlag bool) *CmsTask {
	return &CmsTask{
		bandwidth: bandwidth,
		cdpFlag: cdpFlag,
		taskVolumes: make(map[string]string),
		volumeList: make(map[string]CmsVolume),
    }
}

func checkVolume(c *CmsTask, volumeId string) bool {
	_, find := c.volumeList[volumeId]

	return find
}

func (t *CmsTask) AddVolume(source CmsVolume, target CmsVolume) error {
	if findSource := checkVolume(t, source.VolumeId); findSource {
		return fmt.Errorf("source volume[%s] already exists", source.VolumeId)
	}

	if findTarget := checkVolume(t, target.VolumeId); findTarget {
		return fmt.Errorf("target volume[%s] already exists", target.VolumeId)
	}

	t.taskVolumes[source.VolumeId] = target.VolumeId
	t.volumeList[source.VolumeId] = source
	t.volumeList[target.VolumeId] = target

	return nil
}
