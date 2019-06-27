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

/*
This module defines some essential configuration infos for all storage drivers.

*/

package config

import (
	"io/ioutil"

	log "github.com/golang/glog"
	"github.com/opensds/opensds/pkg/model"
	"gopkg.in/yaml.v2"
)

type PoolProperties struct {
	// The storage type of the storage pool.
	// One of: "block", "file" or "object".
	StorageType string `yaml:"storageType,omitempty"`

	// The locality that pool belongs to.
	AvailabilityZone string `yaml:"availabilityZone,omitempty"`

	// Map of keys and StoragePoolExtraSpec object that represents the properties
	// of the pool, such as supported capabilities.
	// +optional
	Extras model.StoragePoolExtraSpec `yaml:"extras,omitempty"`

	// The volumes belong to the pool can be attached more than once.
	MultiAttach bool `yaml:"multiAttach,omitempty"`
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
