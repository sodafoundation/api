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

/*
This module implements cinder driver for OpenSDS. Cinder driver will pass
these operation requests about volume to gophercloud which is an OpenStack
Go SDK.

*/

package cinder

import (
	"github.com/bouk/monkey"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/opensds/opensds/pkg/utils/config"
	"testing"
)

func TestSetup(t *testing.T) {
	var opt gophercloud.AuthOptions
	defer monkey.UnpatchAll()
	monkey.Patch(openstack.AuthenticatedClient,
		func(options gophercloud.AuthOptions) (*gophercloud.ProviderClient, error) {
			opt = options
			return &gophercloud.ProviderClient{}, nil
		})
	monkey.Patch(openstack.NewBlockStorageV2,
		func(client *gophercloud.ProviderClient, eo gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error) {
			return &gophercloud.ServiceClient{}, nil
		})

	config.CONF.OsdsDock.CinderConfig = "testdata/cinder.yaml"
	d := Driver{}
	d.Setup()
	if opt.IdentityEndpoint != "http://192.168.56.104/identity" {
		t.Error("IdentityEndpoint error.")
	}
	if opt.DomainID != "Default" {
		t.Error("DomainID error.")
	}
	if opt.DomainName != "Default" {
		t.Error("DomainName error.")
	}
	if opt.Username != "admin" {
		t.Error("Username error.")
	}
	if opt.Password != "admin" {
		t.Error("Password error.")
	}
	if opt.TenantID != "04154b841eb644a3947506c54fa73c76" {
		t.Error("TenantID error.")
	}
	if opt.TenantName != "admin" {
		t.Error("TenantName error.")
	}
}
