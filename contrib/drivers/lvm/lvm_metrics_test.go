// Copyright (c) 2019 The OpenSDS Authors.
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
package lvm

import (
	"testing"
)

func TestMetricDriverSetup(t *testing.T) {
	var d = &MetricDriver{}

	if err := d.Setup(); err != nil {
		t.Errorf("Setup lvm metric  driver failed: %+v\n", err)
	}

}

func TestCollectMetrics(t *testing.T) {

	metricList := []string{"IOPS", "ReadThroughput", "WriteThroughput", "ResponseTime", "ServiceTime", "UtilizationPercentage"}
	var metricDriver = &MetricDriver{}
	metricDriver.Setup()
	metricArray, err := metricDriver.CollectMetrics(metricList, "sda")
	if err != nil {
		t.Errorf("collectMetrics call to lvm driver failed: %+v\n", err)
	}
	for _, m := range metricArray {
		t.Log(*m)
	}

}
