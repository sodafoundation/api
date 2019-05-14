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
	"fmt"
	"reflect"
	"testing"

	"github.com/opensds/opensds/pkg/model"
	"github.com/opensds/opensds/pkg/utils/exec"
)

var metricMap map[string]float64 = map[string]float64{"IOPS": 3.16, "ReadThroughput": 4.17, "WriteThroughput": 134.74, "ResponseTime": 2.67, "ServiceTime": 4.00, "UtilizationPercentage": 1.26}
var metricToUnitMap map[string]string = map[string]string{"IOPS": "tps", "ReadThroughput": "KB/s", "WriteThroughput": "KB/s", "ResponseTime": "ms", "ServiceTime": "ms", "UtilizationPercentage": "%"}
var respMap map[string]*MetricFakeResp = map[string]*MetricFakeResp{
	"sar": {`05:26:43  IST       DEV       tps     rkB/s     wkB/s   areq-sz    aqu-sz     await     svctm     %util
			05:26:44      loop0      0.00      0.00      0.00      0.00      0.00      0.00      0.00      0.00
			05:26:44      loop1      0.00      0.00      0.00      0.00      0.00      0.00      0.00      0.00
			05:26:44      loop2      0.00      0.00      0.00      0.00      0.00      0.00      0.00      0.00
			05:26:44      loop3      0.00      0.00      0.00      0.00      0.00      0.00      0.00      0.00
			05:26:44      loop4      0.00      0.00      0.00      0.00      0.00      0.00      0.00      0.00
			05:26:44      loop5      0.00      0.00      0.00      0.00      0.00      0.00      0.00      0.00
			05:26:44      loop6      0.00      0.00      0.00      0.00      0.00      0.00      0.00      0.00
			05:26:44      loop7      0.00      0.00      0.00      0.00      0.00      0.00      0.00      0.00
			05:26:44      opensds--volumes--default-volume--b902e771--8e02--4099--b601--a6b3881f8      3.16      4.17    134.74     42.67      0.01      2.67      4.00      1.26
			05:26:44      loop8      0.00      0.00      0.00      0.00      0.00      0.00      0.00      0.00
			05:26:44      loop9      0.00      0.00      0.00      0.00      0.00      0.00      0.00      0.00
			05:26:44      loop10      0.00      0.00      0.00      0.00      0.00      0.00      0.00      0.00`, nil},
}
var expctdMetricList []string = []string{"IOPS", "ReadThroughput", "WriteThroughput", "ResponseTime", "ServiceTime", "UtilizationPercentage"}

func TestMetricDriverSetup(t *testing.T) {
	var d = &MetricDriver{}
	if err := d.Setup(); err != nil {
		t.Errorf("setup lvm metric  driver failed: %+v\n", err)
	}
}

type MetricFakeExecuter struct {
	RespMap map[string]*MetricFakeResp
}

type MetricFakeResp struct {
	out string
	err error
}

func (f *MetricFakeExecuter) Run(name string, args ...string) (string, error) {
	var cmd = name
	if name == "env" {
		cmd = args[1]
	}
	v, ok := f.RespMap[cmd]
	if !ok {
		return "", fmt.Errorf("can't find specified op: %s", args[1])
	}
	return v.out, v.err
}

func NewMetricFakeExecuter(respMap map[string]*MetricFakeResp) exec.Executer {
	return &MetricFakeExecuter{RespMap: respMap}
}

func TestValidateMetricsSupportList(t *testing.T) {
	var md = &MetricDriver{}
	md.Setup()
	returnedMetricList, err := md.ValidateMetricsSupportList(expctdMetricList, "volume")
	if err != nil {
		t.Error("failed to validate metric list:", err)
	}
	if !reflect.DeepEqual(expctdMetricList, returnedMetricList) {
		t.Errorf("expected %+v, got %+v\n", expctdMetricList, returnedMetricList)
	}
}

func TestCollectMetrics(t *testing.T) {
	var md = &MetricDriver{}
	md.Setup()
	md.cli.RootExecuter = NewMetricFakeExecuter(respMap)
	md.cli.BaseExecuter = NewMetricFakeExecuter(respMap)
	var tempMetricArray []*model.MetricSpec
	for _, element := range expctdMetricList {
		val := metricMap[element]
		expctdLabels := make(map[string]string)
		expctdLabels["device"] = "opensds--volumes--default-volume--b902e771--8e02--4099--b601--a6b3881f8"
		expctdmetricValue := &model.Metric{
			Timestamp: 123456,
			Value:     val,
		}
		expctdMetricValues := make([]*model.Metric, 0)
		expctdMetricValues = append(expctdMetricValues, expctdmetricValue)
		metric := &model.MetricSpec{
			InstanceID:   "b902e771-8e02-4099-b601-a6b3881f8",
			InstanceName: "opensds--volumes--default-volume--b902e771--8e02--4099--b601--a6b3881f8",
			Job:          "OpenSDS",
			Labels:       expctdLabels,
			Component:    "Volume",
			Name:         element,
			Unit:         metricToUnitMap[element],
			AggrType:     "",
			MetricValues: expctdMetricValues,
		}
		tempMetricArray = append(tempMetricArray, metric)
	}
	expectedMetrics := tempMetricArray
	retunMetrics, err := md.CollectMetrics(expctdMetricList, "b902e771-8e02-4099-b601-a6b3881f8")
	if err != nil {
		t.Error("failed to collect stats:", err)
	}
	// we can't use deep equal on metric spec objects as the timesatmp calulation is time.Now() in driver
	// validate equivalence of each metricspec fields against expected except timestamp
	var b bool = true
	for i, m := range expectedMetrics {
		b = b && reflect.DeepEqual(m.InstanceName, retunMetrics[i].InstanceName)
		b = b && reflect.DeepEqual(m.InstanceID, retunMetrics[i].InstanceID)
		b = b && reflect.DeepEqual(m.Job, retunMetrics[i].Job)
		b = b && reflect.DeepEqual(m.Labels, retunMetrics[i].Labels)
		b = b && reflect.DeepEqual(m.Component, retunMetrics[i].Component)
		b = b && reflect.DeepEqual(m.Unit, retunMetrics[i].Unit)
		b = b && reflect.DeepEqual(m.AggrType, retunMetrics[i].AggrType)
		for j, v := range m.MetricValues {
			b = b && reflect.DeepEqual(v.Value, retunMetrics[i].MetricValues[j].Value)
		}
	}
	if !b {
		t.Errorf("expected metric spec")
		printMetricSpec(expectedMetrics)
		t.Errorf("returned metric spec")
		printMetricSpec(retunMetrics)
	}
}

func printMetricSpec(m []*model.MetricSpec) {
	for _, p := range m {
		fmt.Errorf("%+v\n", p)
		for _, v := range p.MetricValues {
			fmt.Errorf("%+v\n", v)
		}
	}

}
