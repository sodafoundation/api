// Copyright 2019 The OpenSDS Authors.
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

package fusionstorage

import (
	"reflect"
	"testing"

	"github.com/opensds/opensds/contrib/drivers/utils/config"
	. "github.com/opensds/opensds/contrib/drivers/utils/config"
	"github.com/opensds/opensds/pkg/model"
)

func TestSetup(t *testing.T) {
	path := "./testdata/fusionstorage.yaml"

	conf := &Config{}
	Parse(conf, path)

	expect := &Config{
		AuthOptions: AuthOptions{
			Username:        "admin",
			Password:        "IaaS@PORTAL-CLOUD9!",
			PwdEncrypter:    "aes",
			EnableEncrypted: false,
			Url:             "https://8.46.195.74:28443",
			FmIp:            "8.46.195.74",
			FsaIp: []string{
				"8.46.195.71",
				"8.46.195.72",
				"8.46.195.73",
			},
			Version: "6.3",
		},
		Pool: map[string]config.PoolProperties{
			"0": {
				StorageType:      "block",
				AvailabilityZone: "default",
				MultiAttach:      true,
				Extras: model.StoragePoolExtraSpec{
					DataStorage: model.DataStorageLoS{
						ProvisioningPolicy: "Thin",
						IsSpaceEfficient:   false,
					},
					IOConnectivity: model.IOConnectivityLoS{
						AccessProtocol: "iscsi",
						MaxIOPS:        7000000,
						MaxBWS:         600,
					},
					Advanced: map[string]interface{}{
						"diskType": "SSD",
						"latency":  "3ms",
					},
				},
			},
			"1": {
				StorageType:      "block",
				AvailabilityZone: "default",
				MultiAttach:      true,
				Extras: model.StoragePoolExtraSpec{
					DataStorage: model.DataStorageLoS{
						ProvisioningPolicy: "Thin",
						IsSpaceEfficient:   false,
					},
					IOConnectivity: model.IOConnectivityLoS{
						AccessProtocol: "iscsi",
						MaxIOPS:        3000000,
						MaxBWS:         300,
					},
					Advanced: map[string]interface{}{
						"diskType": "SSD",
						"latency":  "500ms",
					},
				},
			},
		},
	}
	if !reflect.DeepEqual(expect, conf) {
		t.Errorf("Test driver setup failed:\n expect:%v\n got:\t%v", expect, conf)
	}
}
