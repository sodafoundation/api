// Copyright 2017 The OpenSDS Authors.
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

package config

import (
	"reflect"
	"testing"

	"github.com/opensds/opensds/pkg/model"
)

type Config struct {
	ConfigFile string                    `yaml:"configFile,omitempty"`
	Pool       map[string]PoolProperties `yaml:"pool,flow"`
}

func TestParse(t *testing.T) {
	var conf, path = &Config{}, "testdata/config.yaml"
	var expectedConfig = &Config{
		ConfigFile: "/etc/ceph/ceph.conf",
		Pool: map[string]PoolProperties{
			"rbd": {
				StorageType:      "block",
				AvailabilityZone: "default",
				MultiAttach:      true,
				Extras: model.StoragePoolExtraSpec{
					DataStorage: model.DataStorageLoS{
						ProvisioningPolicy: "Thin",
						IsSpaceEfficient:   true,
					},
					IOConnectivity: model.IOConnectivityLoS{
						AccessProtocol: "rbd",
						MaxIOPS:        8000000,
						MaxBWS:         700,
					},
					Advanced: map[string]interface{}{
						"diskType": "SSD",
						"latency":  "3ms",
					},
				},
			},
		},
	}

	result, err := Parse(conf, path)
	if err != nil {
		t.Errorf("Failed to parse path %s to Config struct: %v\n", path, err)
	}
	if !reflect.DeepEqual(result, expectedConfig) {
		t.Errorf("Expected %+v, got %+v\n", expectedConfig, result)
	}
}
