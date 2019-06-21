// Copyright 2018 The OpenSDS Authors.
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

package multicloud

import (
	"fmt"
	"reflect"
	"testing"
)

const (
	testConfFile = "./testdata/multi-cloud.yaml"
)

func TestLoadConf(t *testing.T) {
	m := &MultiCloud{}
	conf, err := m.loadConf(testConfFile)
	if err != nil {
		t.Errorf("load conf file failed")
	}
	expect := &MultiCloudConf{
		Endpoint:      "http://127.0.0.1:8088",
		UploadTimeout: DefaultUploadTimeout,
		AuthOptions: AuthOptions{
			Strategy:        "keystone",
			AuthUrl:         "http://127.0.0.1/identity",
			DomainName:      "Default",
			UserName:        "admin",
			Password:        "opensds@123",
			TenantName:      "admin",
			PwdEncrypter:    "aes",
			EnableEncrypted: false,
		},
	}
	fmt.Printf("%+v", conf)
	if !reflect.DeepEqual(expect, conf) {
		t.Errorf("load conf file error")
	}
}
