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
	DiskType              string   `yaml:"diskType,omitempty"`
	AZ                    string   `yaml:"AZ,omitempty"`
	RecoveryTimeObjective int      `yaml:"recoveryTimeObjective,omitempty"`
	ProvisioningPolicy    []string `yaml:"provisioningPolicy,omitempty"`
	AccessProtocol        string   `yaml:"accessProtocol,omitempty"`
	MaxIOPS               int      `yaml:"maxIOPS,omitempty"`
	Compress              bool     `yaml:"compress,omitempty"`
	Dedupe                bool     `yaml:"dedupe,omitempty"`
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
	if proper.RecoveryTimeObjective != 0 {
		param["recoveryTimeObjective"] = proper.RecoveryTimeObjective
	}
	if len(proper.ProvisioningPolicy) != 0 {
		param["provisioningPolicy"] = proper.ProvisioningPolicy
	}
	if proper.AccessProtocol != "" {
		param["accessProtocol"] = proper.AccessProtocol
	}
	if proper.MaxIOPS != 0 {
		param["maxIOPS"] = proper.MaxIOPS
	}
	if proper.Compress {
		param["compress"] = proper.Compress
	}
	if proper.Dedupe {
		param["dedupe"] = proper.Dedupe
	}

	return param
}
