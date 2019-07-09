// Copyright 2019 The OpenSDS Authors.
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

type CollectMetricSpec struct {
	*BaseModel

	// the storage type(driver type) on which the metrics are to be collected
	DriverType string `json:"driverType,omitempty"`
}

type GetMetricSpec struct {
	*BaseModel

	// the instance on which the metrics are to be collected
	InstanceId string `json:"instanceId,omitempty"`

	// the name of the metric to retrieve
	MetricName string `json:"metricName,omitempty"`

	StartTime string `json:"startTime,omitempty"`

	EndTime string `json:"endTime,omitempty"`
}

type MetricSpec struct {
	/* Following are the fields used to form name and labels associated with a Metric, same as Prometheus guage name and labels
	Example: node_disk_read_bytes_total{device="dm-0",instance="121.244.95.60:12419",job="prometheus"}
	guage name can be formed by appending Job_Component_Name_Unit_AggrType */

	// Instance ID -\> volumeID/NodeID
	InstanceID string `json:"instanceID,omitempty"`

	// instance name -\> volume name / node name etc.
	InstanceName string `json:"instanceName,omitempty"`

	// job -\> Prometheus/openSDS
	Job string `json:"job,omitempty"`

	/*Labels - There can be multiple componets/properties  associated with a metric , these are catured using this map
	  Example: Labels[pool]="pool1";Labels[device]="dm-0" */
	Labels map[string]string `json:"labels,omitempty"`

	// component -\> disk/logicalVolume/VG etc
	Component string `json:"component,omitempty"`

	// name -\> metric name -\> readRequests/WriteRequests/Latency etc
	Name string `json:"name,omitempty"`

	// unit -\> seconds/bytes/MBs etc
	Unit string `json:"unit,omitempty"`

	// Can be used to determine Total/Avg etc
	AggrType string `json:"aggrType,omitempty"`

	/*If isAggregated ='True' then type of aggregation can be set in this field
	  ie:- if collector is aggregating some metrics and producing a new metric of
	  higher level constructs, then this field can be set as 'Total' to indicate it is
	  aggregated/derived from other metrics.*/

	MetricValues []*Metric `json:"metricValues,omitempty"`
}

type Metric struct {
	Timestamp int64   `json:"timestamp,omitempty"`
	Value     float64 `json:"value"`
}

type NoParam struct{}

type UrlSpec struct {
	Name string `json:"name,omitempty"`
	Url  string `json:"url,omitempty"`
	Desc string `json:"description,omitempty"`
}

type UrlDesc struct {
	Url  string `json:"url,omitempty"`
	Desc string `json:"desc,omitempty"`
}
