// Copyright (c) 2016 Huawei Technologies Co., Ltd. All Rights Reserved.
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
This module defines an standard table of storage plugin. The default storage
plugin is Cinder plugin. If you want to use other storage plugin, just modify
Init() method.

*/

package storagePlugins

import (
	"crypto/tls"
	"net/http"
	"net/url"

	"adapter/storagePlugins/cinder"
	"adapter/storagePlugins/coprhd"
)

type StoragePlugin interface {
	//Any initialization the volume driver does while starting.
	Setup()
	//Any operation the volume driver does while stoping.
	Unset()

	CreateVolume(name string, size int) (string, error)

	GetVolume(volID string) (string, error)

	GetAllVolumes(allowDetails bool) (string, error)

	UpdateVolume(volID string, name string) (string, error)

	DeleteVolume(volID string) (string, error)

	Mount(host string, volID string)

	Unmount(host string, volID string)
}

func Init(resourceType string) StoragePlugin {
	switch resourceType {
	case "cinder":
		var plugin StoragePlugin = &cinder.CinderPlugin{
			"http://162.3.140.36:35357/v2.0",
			"admin",
			"huawei",
			"admin",
		}
		return plugin
	case "coprhd":
		var plugin StoragePlugin = &coprhd.Driver{
			"https://coprhd.emc.com",
			url.UserPassword("admin", "password"),
			&http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
				},
			},
		}
		return plugin
	default:
		return nil
	}
}
