// Copyright (c) 2019 The OpenSDS Authors All Rights Reserved.
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

	// the instance on which the metrics are to be collected
	InstanceId string `json:"instanceId,omitempty"`

	// the list of metrics to be collected
	Metrics []string `json:"metrics,omitempty"`
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
	// Following are the labels associated with Metric, same as Prometheus labels

	//Example: {device="dm-0",instance="121.244.95.60:12419",job="prometheus"}

	// Instance ID -\> volumeID/NodeID

	InstanceID string `json:"instanceID,omitempty"`

	// instance name -\> volume name / node name etc.

	InstanceName string `json:"instanceName,omitempty"`

	// job -\> Prometheus/openSDS

	Job string `json:"job,omitempty"`

	/*associator - Some metric would need specific fields to relate components.

	  Use case could be to query volumes of a particular pool. Attaching the related

	  components as labels would help us to form promQl query efficiently.
	  Example: node_disk_read_bytes_total{instance="121.244.95.60"}
	  Above query will respond with all disks associated with node 121.244.95.60
	  Since associated components vary, we will keep a map in metric struct to denote
	  the associated component type as key and component name as value
	  Example: associator[pool]=pool1 */

	Labels map[string]string `json:"labels,omitempty"`

	// Following fields can be used to form a unique metric name

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
	Value     float64 `json:"value,omitempty"`
}
