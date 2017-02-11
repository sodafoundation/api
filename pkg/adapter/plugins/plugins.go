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

package plugins

import (
	"crypto/tls"
	"errors"
	"net/http"
	"net/url"

	"github.com/opensds/opensds/pkg/adapter/plugins/cinder"
	"github.com/opensds/opensds/pkg/adapter/plugins/coprhd"
	"github.com/opensds/opensds/pkg/adapter/plugins/manila"
)

type VolumePlugin interface {
	//Any initialization the volume driver does while starting.
	Setup()
	//Any operation the volume driver does while stoping.
	Unset()

	CreateVolume(name string, size int) (string, error)

	GetVolume(volID string) (string, error)

	GetAllVolumes(allowDetails bool) (string, error)

	UpdateVolume(volID string, name string) (string, error)

	DeleteVolume(volID string) (string, error)

	MountVolume(volID, host, mountpoint string) (string, error)

	UnmountVolume(volID string, attachement string) (string, error)
}

type SharePlugin interface {
	//Any initialization the file share driver does while starting.
	Setup()
	//Any operation the file share driver does while stoping.
	Unset()

	CreateShare(name string, shrType string, shrProto string, size int) (string, error)

	GetShare(shrID string) (string, error)

	GetAllShares(allowDetails bool) (string, error)

	UpdateShare(shrID string, name string) (string, error)

	DeleteShare(shrID string) (string, error)
}

func InitVP(resourceType string) (VolumePlugin, error) {
	switch resourceType {
	case "cinder":
		var plugin VolumePlugin = &cinder.CinderPlugin{
			Host:        "http://162.3.140.36:35357/v2.0",
			Username:    "admin",
			Password:    "huawei",
			ProjectName: "admin",
		}
		return plugin, nil
	case "coprhd":
		var plugin VolumePlugin = &coprhd.Driver{
			"https://coprhd.emc.com",
			url.UserPassword("admin", "password"),
			&http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
				},
			},
		}
		return plugin, nil
	default:
		err := errors.New("Can't find this resource type in backend storage.")
		return nil, err
	}
}

func InitSP(resourceType string) (SharePlugin, error) {
	switch resourceType {
	case "manila":
		var plugin SharePlugin = &manila.ManilaPlugin{
			Host:        "http://162.3.140.36:35357/v2.0",
			Username:    "admin",
			Password:    "huawei",
			ProjectName: "admin",
		}
		return plugin, nil
	default:
		err := errors.New("Can't find this resource type in backend storage.")
		return nil, err
	}
}
