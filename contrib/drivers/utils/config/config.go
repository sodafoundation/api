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

package config

import (
	"io/ioutil"

	log "github.com/golang/glog"
	"gopkg.in/yaml.v2"
)

type PoolProperties struct {
	DiskType        string `yaml:"diskType,omitempty"`
	AZ              string `yaml:"AZ,omitempty"`
	AccessProtocol  string `yaml:"accessProtocol,omitempty"`
	ThinProvisioned bool   `yaml:"thinProvisioned,omitempty"`
	Compressed      bool   `yaml:"compressed,omitempty"`
	// Besides those basic pool properties above, vendors can configure some
	// advanced features (IOPS, throughout, latency, etc) themselves, all these
	// properties can be exposed to controller scheduler and filtered by
	// selector in a extensible way.
	Advanced map[string]interface{} `yaml:"advanced,omitempty"`
}

func Parse(conf interface{}, p string) (interface{}, error) {
	confYaml, err := ioutil.ReadFile(p)
	if err != nil {
		log.Fatalf("Read config yaml file (%s) failed, reason:(%v)", p, err)
		return nil, err
	}
	if err = yaml.Unmarshal(confYaml, conf); err != nil {
		log.Fatalf("Parse error: %v", err)
		return nil, err
	}
	return conf, nil
}

func BuildDefaultPoolParam(proper PoolProperties) map[string]interface{} {
	param := make(map[string]interface{})
	if proper.DiskType != "" {
		param["diskType"] = proper.DiskType
	}
	if proper.AccessProtocol != "" {
		param["accessProtocol"] = proper.AccessProtocol
	}
	if proper.ThinProvisioned {
		param["thinProvisioned"] = proper.ThinProvisioned
	}
	if proper.Compressed {
		param["compressed"] = proper.Compressed
	}
	if len(proper.Advanced) != 0 {
		param["advanced"] = proper.Advanced
	}

	return param
}
