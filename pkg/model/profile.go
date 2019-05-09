// Copyright (c) 2017 Huawei Technologies Co., Ltd. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

/*
This module implements the common data structure.

*/

package model

import (
	"encoding/json"
	"reflect"
	"strings"

	"github.com/golang/glog"
)

// An OpenSDS profile is identified by a unique name and ID. With additional
// profile properties, each profile contains a set of tags of storage
// capabilities which are desirable features for a class of applications.
type ProfileSpec struct {
	*BaseModel

	// The name of the profile.
	Name string `json:"name,omitempty"`

	// The description of the profile.
	// +optional
	Description string `json:"description,omitempty"`

	// The storage type of the profile.
	// One of: block, file or object.
	StorageType string `json:"storageType,omitempty"`

	// ProvisioningProperties represents some suggested properties for performing
	// provisioning policies.
	// +optional
	ProvisioningProperties ProvisioningPropertiesSpec `json:"provisioningProperties,omitempty"`

	// ReplicationProperties represents some suggested properties for performing
	// replicaiton policies.
	// +optional
	ReplicationProperties ReplicationPropertiesSpec `json:"replicationProperties,omitempty"`

	// SnapshotProperties represents some suggested properties for performing
	// snapshot policies.
	// +optional
	SnapshotProperties SnapshotPropertiesSpec `json:"snapshotProperties,omitempty"`

	// DataProtectionProperties represents some suggested properties for
	// performing data protection policies.
	// +optional
	DataProtectionProperties DataProtectionPropertiesSpec `json:"dataProtectionProperties,omitempty"`

	// CustomProperties is a map of keys and JSON object that represents the
	// customized properties of profile, such as requested capabilities
	// including diskType, latency, deduplicaiton, compression and so forth.
	// +optional
	CustomProperties CustomPropertiesSpec `json:"customProperties,omitempty"`
}

func NewProfileFromJson(s string) *ProfileSpec {
	p := &ProfileSpec{}
	err := json.Unmarshal([]byte(s), p)
	if err != nil {
		glog.Errorf("Unmarshal json to ProfileSpec failed, %v", err)
	}
	return p
}

func (p *ProfileSpec) ToJson() string {
	b, err := json.Marshal(p)
	if err != nil {
		glog.Errorf("ProfileSpec convert to json failed, %v", err)
	}
	return string(b)
}

type ProvisioningPropertiesSpec struct {
	// DataStorage represents some suggested data storage capabilities.
	DataStorage DataStorageLoS `json:"dataStorage,omitempty"`
	// IOConnectivity represents some suggested IO connectivity capabilities.
	IOConnectivity IOConnectivityLoS `json:"ioConnectivity,omitempty"`
}

func (pps ProvisioningPropertiesSpec) IsEmpty() bool {
	r := reflect.DeepEqual(ProvisioningPropertiesSpec{}, pps)
	return r
}

type ReplicationPropertiesSpec struct {
	// DataProtection represents some suggested data protection capabilities.
	DataProtection DataProtectionLoS `json:"dataProtection,omitempty"`
	// ReplicaInfos represents some suggested data replication information.
	ReplicaInfos struct {
		// The enumeration literal specifies whether the target elements will be
		// updated synchronously or asynchronously. The possible values of this
		// property could be:
		// * Active: Active-Active (i.e. bidirectional) synchronous updates.
		// * Adaptive: Allows implementation to switch between synchronous
		//   and asynchronous modes.
		// * Asynchronous: Asynchronous updates.
		// * Synchronous: Synchronous updates.
		ReplicaUpdateMode string `json:"replicaUpdateMode,omitempty"`
		// ConsistencyEnabled indicates that the source and target shall be
		// consistent. The default value is false.
		ConsistencyEnabled bool `json:"consistencyEnabled,omitempty"`
		// ReplicationPeriod shall be an ISO 8601 duration that defines the
		// duration of performing replication operation. For example,
		// "P3Y6M4DT12H30M5S" represents a duration of "3 years, 6 months,
		// 4 days, 12 hours, 30 minutes and 5 seconds".
		ReplicationPeriod string `json:"replicationPeriod,omitempty"`
		// ReplicationBandwidth specifies the maximum bandwidth for performing
		// replication operation.
		// +units:Mb/s
		ReplicationBandwidth int64 `json:"replicationBandwidth,omitempty"`
	} `json:"replicaInfos,omitempty"`
}

func (rps ReplicationPropertiesSpec) IsEmpty() bool {
	if (ReplicationPropertiesSpec{}) == rps {
		return true
	}
	return false
}

type SnapshotPropertiesSpec struct {
	// The property defines how to take snapshots.
	Schedule struct {
		// This vaule is represented as a string in ISO 8601 datetime format.
		// ISO 8601 sets out an internationally agreed way to represent dates:
		// yyyy-mm-ddThh:mm:ss.ffffff. For example, "3:53pm, September 15, 2008"
		// is represented as "2008-09-15T15:53:00".
		Datetime string `json:"datetime,omitempty"`
		// The value specifies the duration of executing a operation, which
		// contains three options including Daily, Weekly and Monthly.
		Occurrence string `json:"occurrence,omitempty"`
	} `json:"schedule,omitempty"`
	Retention struct {
		// The value specifies the total number of snapshots for retention.
		// +optional
		Number int64 `json:"number,omitempty"`
		// The value specifies the duration of snapshots for retention.
		// +optional
		// +units:day
		Duration int64 `json:"duration,omitempty"`
	} `json:"retention,omitempty"`
	Topology struct {
		Bucket string `json:"bucket,omitempty"` // This is virtual bucket managed by multi-cloud
	} `json:"topology,omitempty"`
}

func (sps SnapshotPropertiesSpec) IsEmpty() bool {
	if (SnapshotPropertiesSpec{}) == sps {
		return true
	}
	return false
}

type DataProtectionPropertiesSpec struct {
	// DataProtection represents some suggested data protection capabilities.
	DataProtection DataProtectionLoS `json:"dataProtection,omitempty"`
	// ConsistencyEnabled indicates that the source and target shall be
	// consistent. The default value is false.
	ConsistencyEnabled bool `json:"consistencyEnabled,omitempty"`
}

func (dps DataProtectionPropertiesSpec) IsEmpty() bool {
	if (DataProtectionPropertiesSpec{}) == dps {
		return true
	}
	return false
}

// CustomPropertiesSpec is a dictionary object that contains unique keys and
// JSON objects.
type CustomPropertiesSpec map[string]interface{}

func (cps CustomPropertiesSpec) IsEmpty() bool {
	if nil == cps {
		return true
	}
	return false
}

func (cps CustomPropertiesSpec) Encode() []byte {
	parmBody, _ := json.Marshal(&cps)
	return parmBody
}

func (cps CustomPropertiesSpec) GetCapabilitiesProperties() map[string]interface{} {
	caps := make(map[string]interface{})
	if cps.IsEmpty() {
		return caps
	}
	for k, v := range cps {
		words := strings.Split(k, ":")
		if len(words) > 1 && words[0] == "capabilities" {
			caps[words[1]] = v
		}
	}
	return caps
}
