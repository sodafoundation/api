// Copyright 2019 The OpenSDS Authors.
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
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/opensds/opensds/pkg/utils/constants"
)

const (
	OpensdsEndpoint = "OPENSDS_ENDPOINT"
)

var (
	cacert string
)

// Client is a struct for exposing some operations of opensds resources.
type Client struct {
	*ProfileMgr
	*DockMgr
	*PoolMgr
	*VolumeMgr
	*VersionMgr
	*ReplicationMgr
	*FileShareMgr

	cfg *Config
}

// Config is a struct that defines some options for calling the Client.
type Config struct {
	Endpoint    string
	AuthOptions AuthOptions
}

// NewClient method creates a new Client.
func NewClient(c *Config) (*Client, error) {
	// If endpoint field not specified,use the default value localhost.
	if c.Endpoint == "" {
		c.Endpoint = constants.DefaultOpensdsEndpoint
		log.Printf("WARNING: OpenSDS Endpoint is not specified, use default(%s)\n", c.Endpoint)
	}

	// If https is enabled, CA cert file should be provided.
	u, _ := url.Parse(c.Endpoint)
	if u.Scheme == "https" {
		cacert = constants.OpensdsCaCertFile
		_, err := os.Stat(cacert)
		if err != nil {
			if os.IsNotExist(err) {
				return nil, fmt.Errorf("CA file(%s) doesn't exist", cacert)
			}
		}
	}

	var r Receiver
	var err error
	switch c.AuthOptions.(type) {
	case *NoAuthOptions:
		r = NewReceiver()
	case *KeystoneAuthOptions:
		r, err = NewKeystoneReceiver(c.AuthOptions.(*KeystoneAuthOptions))
		if err != nil {
			return nil, fmt.Errorf("keystone authentication failed")
		}
	default:
		log.Println("WARNING: Not support auth options, use default(noauth).")
		r = NewReceiver()
		c.AuthOptions = NewNoauthOptions(constants.DefaultTenantId)
	}

	t := c.AuthOptions.GetTenantId()
	return &Client{
		cfg:            c,
		ProfileMgr:     NewProfileMgr(r, c.Endpoint, t),
		DockMgr:        NewDockMgr(r, c.Endpoint, t),
		PoolMgr:        NewPoolMgr(r, c.Endpoint, t),
		VolumeMgr:      NewVolumeMgr(r, c.Endpoint, t),
		VersionMgr:     NewVersionMgr(r, c.Endpoint, t),
		ReplicationMgr: NewReplicationMgr(r, c.Endpoint, t),
		FileShareMgr:   NewFileShareMgr(r, c.Endpoint, t),
	}, nil
}

// Reset method is defined to clean Client struct.
func (c *Client) Reset() *Client {
	return &Client{}
}

func processListParam(args []interface{}) (string, error) {
	var urlParam []string
	var output string

	// If args is empty, just return the output immeadiately.
	if len(args) == 0 {
		return "", nil
	}
	// Add some limits for input args parameter.
	if len(args) != 1 {
		return "", errors.New("args should only support one parameter")
	}
	filter, ok := args[0].(map[string]string)
	if !ok {
		return "", errors.New("args element type should be map[string]string")
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
		output = strings.Join(urlParam, "&")
	}

	return output, nil
}
