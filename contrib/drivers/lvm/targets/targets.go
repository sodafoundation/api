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

const (
	iscsiTgtPrefix = "iqn.2017-10.io.opensds:"
)

// Target is an interface for exposing some operations of different targets,
// currently support iscsiTarget.
type Target interface {
	CreateExport(volId, path, hostIp, initiator string) (map[string]interface{}, error)

	RemoveExport(volId string) error
}

// NewTarget method creates a new iscsi target.
func NewTarget(bip string, tgtConfDir string) Target {
	return &iscsiTarget{
		ISCSITarget: NewISCSITarget(bip, tgtConfDir),
	}
}

type iscsiTarget struct {
	ISCSITarget
}

func (t *iscsiTarget) CreateExport(volId, path, hostIp, initiator string) (map[string]interface{}, error) {
	tgtIqn := iscsiTgtPrefix + volId
	if err := t.CreateISCSITarget(volId, tgtIqn, path, hostIp, initiator, []string{}); err != nil {
		return nil, err
	}
	lunId := t.GetLun(path)
	return map[string]interface{}{
		"targetDiscovered": true,
		"targetIQN":        tgtIqn,
		"targetPortal":     t.ISCSITarget.(*tgtTarget).BindIp + ":3260",
		"discard":          false,
		"targetLun":        lunId,
	}, nil
}

func (t *iscsiTarget) RemoveExport(volId string) error {
	tgtIqn := iscsiTgtPrefix + volId
	return t.RemoveISCSITarget(volId, tgtIqn)
}
