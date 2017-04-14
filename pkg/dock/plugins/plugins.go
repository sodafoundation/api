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
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/opensds/opensds/pkg/dock/plugins/cinder"
	"github.com/opensds/opensds/pkg/dock/plugins/coprhd"
	"github.com/opensds/opensds/pkg/dock/plugins/manila"
)

type VolumePlugin interface {
	//Any initialization the volume driver does while starting.
	Setup()
	//Any operation the volume driver does while stoping.
	Unset()

	CreateVolume(name string, volType string, size int32) (string, error)

	GetVolume(volID string) (string, error)

	GetAllVolumes(allowDetails bool) (string, error)

	DeleteVolume(volID string) (string, error)

	AttachVolume(volID string) (string, error)

	DetachVolume(device string) (string, error)
}

type SharePlugin interface {
	//Any initialization the file share driver does while starting.
	Setup()
	//Any operation the file share driver does while stoping.
	Unset()

	CreateShare(name string, shrType string, shrProto string, size int32) (string, error)

	GetShare(shrID string) (string, error)

	GetAllShares(allowDetails bool) (string, error)

	DeleteShare(shrID string) (string, error)

	AttachShare(shrID string) (string, error)

	DetachShare(device string) (string, error)
}

type cinderConfig struct {
	Host        string   `json:"host"`
	Methods     []string `json:"methods"`
	Username    string   `json:"username"`
	Password    string   `json:"password"`
	ProjectId   string   `json:"projectId"`
	ProjectName string   `json:"projectName"`
}

type manilaConfig struct {
	Host        string   `json:"host"`
	Methods     []string `json:"methods"`
	Username    string   `json:"username"`
	Password    string   `json:"password"`
	ProjectId   string   `json:"projectId"`
	ProjectName string   `json:"projectName"`
}

type coprHDConfig struct {
	Host     string `json:"host"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type pluginsConfig struct {
	Cinder cinderConfig `json:"cinder"`
	Manila manilaConfig `json:"manila"`
	CoprHD coprHDConfig `json:"coprhd"`
}

func InitVP(resourceType string) (VolumePlugin, error) {
	config := readBackendConfigFile()

	switch resourceType {
	case "cinder":
		return &cinder.CinderPlugin{
			Host:        config.Cinder.Host,
			Methods:     config.Cinder.Methods,
			Username:    config.Cinder.Username,
			Password:    config.Cinder.Password,
			ProjectId:   config.Cinder.ProjectId,
			ProjectName: config.Cinder.ProjectName,
		}, nil
	case "coprhd":
		return &coprhd.Driver{
			Url:   config.CoprHD.Host,
			Creds: url.UserPassword(config.CoprHD.Username, config.CoprHD.Password),
			HttpClient: &http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
				},
			},
		}, nil
	default:
		err := errors.New("Can't find this resource type in backend storage.")
		return nil, err
	}
}

func InitSP(resourceType string) (SharePlugin, error) {
	config := readBackendConfigFile()

	switch resourceType {
	case "manila":
		return &manila.ManilaPlugin{
			Host:        config.Manila.Host,
			Methods:     config.Manila.Methods,
			Username:    config.Manila.Username,
			Password:    config.Manila.Password,
			ProjectId:   config.Manila.ProjectId,
			ProjectName: config.Manila.ProjectName,
		}, nil
	default:
		err := errors.New("Can't find this resource type in backend storage.")
		return nil, err
	}
}

// readBackendConfigFile provides access to credentials in backend resource plugins.
func readBackendConfigFile() *pluginsConfig {
	var config *pluginsConfig

	userJSON, err := ioutil.ReadFile("/etc/opensds/config.json")
	if err != nil {
		log.Println("ReadFile json failed:", err)
	}
	if err = json.Unmarshal(userJSON, config); err != nil {
		log.Println("Unmarshal json failed:", err)
	}
	return config
}
