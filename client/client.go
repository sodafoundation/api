// Copyright (c) 2017 Huawei Technologies Co., Ltd. All Rights Reserved.
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

package client

import (
	"errors"
	"log"
	"strings"

	"github.com/opensds/opensds/pkg/utils/constants"
)

const (
	OpensdsEndpoint = "OPENSDS_ENDPOINT"
)

//// Client is a struct for exposing some operations of opensds resources.
//type Client struct {
//	cfg *Config
//}

// Config is a struct that defines some options for calling the Client.
type Client struct {
	*ProfileMgr
	*DockMgr
	*PoolMgr
	*VolumeMgr
	*VersionMgr
	*ReplicationMgr
}

// NewClient method creates a new Client.
func NewClient(endPoint string, authOptions AuthOptions) *Client {
	// If endpoint field not specified,use the default value localhost.
	if endPoint == "" {
		endPoint = constants.DefaultOpensdsEndpoint
		log.Printf("Warnning: OpenSDS Endpoint is not specified using the default value(%s)", endPoint)
	}

	var receiver Receiver
	switch authOptions.(type) {
	case *NoAuthOptions:
		receiver = NewReceiver(authOptions.(*NoAuthOptions))
	case *KeystoneAuthOptions:
		receiver = NewKeystoneReciver(authOptions.(*KeystoneAuthOptions))
	default:
		log.Printf("Warnning: Not support auth options, use default")

		authOptions = NewNoauthOptions(constants.DefaultTenantId)
		receiver = NewReceiver(authOptions.(*NoAuthOptions))
	}

	return &Client{
		ProfileMgr:     NewProfileMgr(receiver, endPoint),
		DockMgr:        NewDockMgr(receiver, endPoint),
		PoolMgr:        NewPoolMgr(receiver, endPoint),
		VolumeMgr:      NewVolumeMgr(receiver, endPoint),
		VersionMgr:     NewVersionMgr(receiver, endPoint),
		ReplicationMgr: NewReplicationMgr(receiver, endPoint),
	}
}

// Reset method is defined to clean Client struct.
func (c *Client) Reset() *Client {
	c = &Client{}
	return c
}

func processListParam(args []interface{}) (string, error) {
	var filter map[string]string
	var u string
	var urlParam []string

	if len(args) > 0 {
		if len(args) > 1 {
			return "", errors.New("only support one parameter that must be map[string]string")
		}
		filter = args[0].(map[string]string)
	}

	if filter != nil {
		for k, v := range filter {
			if v == "" {
				continue
			}
			urlParam = append(urlParam, k+"="+v)
		}
	}

	if len(urlParam) > 0 {
		u = strings.Join(urlParam, "&")
	}

	return u, nil
}
