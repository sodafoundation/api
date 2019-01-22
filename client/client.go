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

// Config is a struct that defines some options for calling the Client.
type Config struct {
	Endpoint    string
	AuthOptions AuthOptions
}

var config *Config

// Client is a struct for exposing some operations of opensds resources.
type Client struct {
	*ProfileMgr
	*DockMgr
	*PoolMgr
	*VolumeMgr
	*VersionMgr
	*ReplicationMgr
}

// NewClient method creates a new Client.
func NewClient(c *Config) *Client {
	// If endpoint field not specified,use the default value localhost.
	if c.Endpoint == "" {
		c.Endpoint = constants.DefaultOpensdsEndpoint
		log.Printf("Warnning: OpenSDS endpoint is not specified using the default value(%s)", c.Endpoint)
	}

	var receiver Receiver
	switch c.AuthOptions.(type) {
	case *NoAuthOptions:
		receiver = NewReceiver()
	case *KeystoneAuthOptions:
		receiver = NewKeystoneReciver()
		c.AuthOptions, _ = GetToken(c.AuthOptions.(*KeystoneAuthOptions))
	default:
		log.Printf("Warnning: Not support auth options, use default")
		c.AuthOptions = NewNoauthOptions(constants.DefaultTenantId)
		receiver = NewReceiver()
	}

	config = c

	return &Client{
		ProfileMgr:     NewProfileMgr(receiver),
		DockMgr:        NewDockMgr(receiver),
		PoolMgr:        NewPoolMgr(receiver),
		VolumeMgr:      NewVolumeMgr(receiver),
		VersionMgr:     NewVersionMgr(receiver),
		ReplicationMgr: NewReplicationMgr(receiver),
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
