// Copyright 2017 The OpenSDS Authors.
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
	"log"

	"github.com/opensds/opensds/pkg/utils/constants"
)

// Client is a struct for exposing some operations of opensds resources.
type Client struct {
	*ProfileMgr
	*DockMgr
	*PoolMgr
	*VolumeMgr

	cfg *Config
}

// Config is a struct that defines some options for calling the Client.
type Config struct {
	Endpoint string
}

// NewClient method creates a new Client.
func NewClient(c *Config) *Client {
	// If endpoint field not specified,use the default value localhost.
	if c.Endpoint == "" {
		c.Endpoint = constants.DefaultOpensdsEndpoint
		log.Printf("OpenSDS Endpoint is not specified using the default value(%s)", c.Endpoint)
	}

	return &Client{
		cfg:        c,
		ProfileMgr: NewProfileMgr(c.Endpoint),
		DockMgr:    NewDockMgr(c.Endpoint),
		PoolMgr:    NewPoolMgr(c.Endpoint),
		VolumeMgr:  NewVolumeMgr(c.Endpoint),
	}
}

// Reset method is defined to clean Client struct.
func (c *Client) Reset() *Client {
	c = &Client{}
	return c
}
