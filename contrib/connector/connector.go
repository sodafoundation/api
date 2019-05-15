// Copyright (c) 2018 Huawei Technologies Co., Ltd. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package connector

import (
	"fmt"
	"log"

	"github.com/opensds/opensds/contrib/drivers/utils/config"
)

const (
	FcDriver = config.FCProtocol
	PortName = "port_name"
	NodeName = "node_name"
	Wwpn     = "wwpn"
	Wwnn     = "wwnn"

	IscsiDriver = config.ISCSIProtocol
	Iqn         = "iqn"

	RbdDriver = config.RBDProtocol

	NvmeofDriver = config.NVMEOFProtocol
	Nqn          = "nqn"
)

// Connector implementation
type Connector interface {
	Attach(map[string]interface{}) (string, error)
	Detach(map[string]interface{}) error
	GetInitiatorInfo() (string, error)
}

var cnts = map[string]Connector{}

// NewConnector implementation
func NewConnector(cType string) Connector {
	if cnt, exist := cnts[cType]; exist {
		return cnt
	}

	log.Printf("%s is not registered to connector", cType)
	return nil
}

// RegisterConnector implementation
func RegisterConnector(cType string, cnt Connector) error {
	if _, exist := cnts[cType]; exist {
		return fmt.Errorf("Connector %s already exist", cType)
	}

	cnts[cType] = cnt
	return nil
}

// UnregisterConnector implementation
func UnregisterConnector(cType string) {
	if _, exist := cnts[cType]; !exist {
		return
	}

	delete(cnts, cType)
	return
}
