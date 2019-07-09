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

import "github.com/opensds/opensds/contrib/drivers/utils/config"

const (
	iscsiTgtPrefix  = "iqn.2017-10.io.opensds:"
	nvmeofTgtPrefix = "nqn.2019-01.com.opensds:nvme:"
)

// Target is an interface for exposing some operations of different targets,
// currently support iscsiTarget.
type Target interface {
	CreateExport(volId, path, hostIp, initiator string, chapAuth []string) (map[string]interface{}, error)

	RemoveExport(volId, hostIp string) error
}

// NewTarget method creates a new target based on its type.
func NewTarget(bip string, tgtConfDir string, access string) Target {
	switch access {
	case config.ISCSIProtocol:
		return &iscsiTarget{
			ISCSITarget: NewISCSITarget(bip, tgtConfDir),
		}
	case config.NVMEOFProtocol:
		return &nvmeofTarget{
			NvmeofTarget: NewNvmeofTarget(bip, tgtConfDir),
		}
	default:
		return nil
	}
}

type iscsiTarget struct {
	ISCSITarget
}

func (t *iscsiTarget) CreateExport(volId, path, hostIp, initiator string, chapAuth []string) (map[string]interface{}, error) {
	tgtIqn := iscsiTgtPrefix + volId
	if err := t.CreateISCSITarget(volId, tgtIqn, path, hostIp, initiator, chapAuth); err != nil {
		return nil, err
	}
	lunId := t.GetLun(path)
	conn := map[string]interface{}{
		"targetDiscovered": true,
		"targetIQN":        []string{tgtIqn},
		"targetPortal":     []string{t.ISCSITarget.(*tgtTarget).BindIp + ":3260"},
		"discard":          false,
		"targetLun":        lunId,
	}
	if len(chapAuth) == 2 {
		conn["authMethod"] = "chap"
		conn["authUserName"] = chapAuth[0]
		conn["authPassword"] = chapAuth[1]
	}
	return conn, nil
}

func (t *iscsiTarget) RemoveExport(volId, hostIp string) error {
	tgtIqn := iscsiTgtPrefix + volId
	return t.RemoveISCSITarget(volId, tgtIqn, hostIp)
}

type nvmeofTarget struct {
	NvmeofTarget
}

func (t *nvmeofTarget) CreateExport(volId, path, hostIp, initiator string, chapAuth []string) (map[string]interface{}, error) {
	tgtNqn := nvmeofTgtPrefix + volId
	// So far nvmeof transtport type is defaultly set as tcp because of its widely use, but it can also be rdma/fc.
	// The difference of transport type leads to different performance of volume attachment latency and iops.
	// This choice of transport type depends on 3 following factors:
	// 1. initiator's latency/iops requiremnet
	// 2. initiator's availiable nic(whether the inititator can use rdma/fc/tcpip)
	// 3. target server's availiable nic(whether the target server can use rdma/fc/tcpip)
	// According to the opensds architecture, it is a more approprite way for the opensds controller
	//to take the decision in the future.
	var transtype string
	transtype = "tcp"
	if err := t.CreateNvmeofTarget(volId, tgtNqn, path, initiator, transtype); err != nil {
		return nil, err
	}
	conn := map[string]interface{}{
		"targetDiscovered": true,
		"targetNQN":        tgtNqn,
		"targetIP":         t.NvmeofTarget.(*NvmeoftgtTarget).BindIp,
		"targetPort":       "4420",
		"hostNqn":          initiator,
		"discard":          false,
		"transporType":     transtype,
	}

	return conn, nil
}

func (t *nvmeofTarget) RemoveExport(volId, hostIp string) error {
	tgtNqn := nvmeofTgtPrefix + volId
	// So far nvmeof transtport type is defaultly set as tcp because of its widely use, but it can also be rdma/fc.
	// The difference of transport type leads to different performance of volume attachment latency and iops.
	// This choice of transport type depends on 3 following factors:
	// 1. initiator's latency/iops requiremnet
	// 2. initiator's availiable nic(whether the inititator can use rdma/fc/tcpip)
	// 3. target server's availiable nic(whether the target server can use rdma/fc/tcpip)
	// According to the opensds architecture, it is a more approprite way for the opensds controller
	//to take the decision in the future.
	var transtype string
	transtype = "tcp"
	return t.RemoveNvmeofTarget(volId, tgtNqn, transtype)
}
