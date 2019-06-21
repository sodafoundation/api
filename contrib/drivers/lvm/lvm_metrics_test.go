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

package lvm

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"

	"github.com/opensds/opensds/pkg/model"
	"github.com/opensds/opensds/pkg/utils/exec"
)

var metricMap map[string]float64 = map[string]float64{"iops": 3.16, "read_throughput": 4.17, "write_throughput": 134.74, "response_time": 2.67, "service_time": 4.00, "utilization": 1.26}
var metricToUnitMap map[string]string = map[string]string{"iops": "tps", "read_throughput": "kbs", "write_throughput": "kbs", "response_time": "ms", "service_time": "ms", "utilization": "prcnt"}
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
			05:26:44      opensds--volumes--default-volume-d96cc42b-b285-474e-aa98-c61e66df7461      6.26.00      8.27      268.74      84.67      0.02      4.67      8.00      2.46
			05:26:44      loop9      0.00      0.00      0.00      0.00      0.00      0.00      0.00      0.00
			05:26:44      loop10      8.77      12.22      32.56      78.01      0.00      0.01      0.01      12.21`, nil},
	"lvs": {`LV                                          VG                      Attr       LSize Pool Origin Data%  Meta%  Move Log Cpy%Sync Convert
volume-b902e771-8e02-4099-b601-a6b3881f8 opensds-volumes-default -wi-a----- 1.00g                                                    
volume-d96cc42b-b285-474e-aa98-c61e66df7461 opensds-volumes-default -wi-a----- 1.00g`, nil},
	"pvs": {`PV          VG                      Fmt  Attr PSize   PFree  
/dev/loop10 opensds-volumes-default lvm2 a--  <20.00g <18.00g`, nil},
}
var expectdVgs []string = []string{"opensds-volumes-default", "opensds-volumes-default"}
var expctdMetricList []string = []string{"iops", "read_throughput", "write_throughput", "response_time", "service_time", "utilization"}
var expctedVolList []string = []string{"volume-b902e771-8e02-4099-b601-a6b3881f8", "volume-d96cc42b-b285-474e-aa98-c61e66df7461"}
var expctedDiskList []string = []string{"/dev/loop10"}
var expctedLabelMap map[string]map[string]string = map[string]map[string]string{"loop10": map[string]string{"device": "loop10"}, "opensds--volumes--default-volume-d96cc42b-b285-474e-aa98-c61e66df7461": map[string]string{"device": "opensds--volumes--default-volume-d96cc42b-b285-474e-aa98-c61e66df7461"}, "opensds--volumes--default-volume--b902e771--8e02--4099--b601--a6b3881f8": map[string]string{"device": "opensds--volumes--default-volume--b902e771--8e02--4099--b601--a6b3881f8"}}
var expctedMetricMap map[string]map[string]string = map[string]map[string]string{"loop10": map[string]string{"InstanceName": "loop10", "iops": "8.77", "read_throughput": "12.22", "response_time": "0.01", "service_time": "0.01", "write_throughput": "32.56", "utilization": "12.21"},
	"opensds--volumes--default-volume-d96cc42b-b285-474e-aa98-c61e66df7461":   map[string]string{"InstanceName": "opensds--volumes--default-volume-d96cc42b-b285-474e-aa98-c61e66df7461", "iops": "6.26", "read_throughput": "8.27", "response_time": "4.67", "service_time": "8.00", "write_throughput": "268.74", "utilization": "2.46"},
	"opensds--volumes--default-volume--b902e771--8e02--4099--b601--a6b3881f8": map[string]string{"InstanceName": "opensds--volumes--default-volume--b902e771--8e02--4099--b601--a6b3881f8", "iops": "3.16", "read_throughput": "4.17", "response_time": "2.67", "service_time": "4.00", "write_throughput": "134.74", "utilization": "1.26"}}

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

func TestGetMetricList(t *testing.T) {
	var md = &MetricDriver{}
	md.Setup()
	returnedMetricList, err := md.GetMetricList("volume")
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
	for i, volume := range expctedVolList {
		convrtedVolID := convert(volume, expectdVgs[i])
		thismetricMAp := expctedMetricMap[convrtedVolID]
		thisLabelMap := expctedLabelMap[convrtedVolID]
		for _, element := range expctdMetricList {
			val, _ := strconv.ParseFloat(thismetricMAp[element], 64)

			expctdmetricValue := &model.Metric{
				Timestamp: 123456,
				Value:     val,
			}
			expctdMetricValues := make([]*model.Metric, 0)
			expctdMetricValues = append(expctdMetricValues, expctdmetricValue)
			metric := &model.MetricSpec{
				InstanceID:   volume,
				InstanceName: thismetricMAp["InstanceName"],
				Job:          "lvm",
				Labels:       thisLabelMap,
				Component:    "volume",
				Name:         element,
				Unit:         metricToUnitMap[element],
				AggrType:     "",
				MetricValues: expctdMetricValues,
			}
			tempMetricArray = append(tempMetricArray, metric)
		}
	}
	for _, disk := range expctedDiskList {
		convrtedVolID := formatDiskName(disk)
		thismetricMAp := expctedMetricMap[convrtedVolID]
		thisLabelMap := expctedLabelMap[convrtedVolID]
		fmt.Println(thismetricMAp)
		for _, element := range expctdMetricList {
			val, _ := strconv.ParseFloat(thismetricMAp[element], 64)
			expctdmetricValue := &model.Metric{
				Timestamp: 123456,
				Value:     val,
			}
			expctdMetricValues := make([]*model.Metric, 0)
			expctdMetricValues = append(expctdMetricValues, expctdmetricValue)
			metric := &model.MetricSpec{
				InstanceID:   disk,
				InstanceName: thismetricMAp["InstanceName"],
				Job:          "lvm",
				Labels:       thisLabelMap,
				Component:    "disk",
				Name:         element,
				Unit:         metricToUnitMap[element],
				AggrType:     "",
				MetricValues: expctdMetricValues,
			}
			tempMetricArray = append(tempMetricArray, metric)
		}
	}
	expectedMetrics := tempMetricArray
	retunMetrics, err := md.CollectMetrics()
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
		t.Log("expected metric spec")
		printMetricSpec(expectedMetrics)
		t.Log("returned metric spec")
		printMetricSpec(retunMetrics)
	}
}

func printMetricSpec(m []*model.MetricSpec) {
	for _, p := range m {
		fmt.Printf("%+v\n", p)
		for _, v := range p.MetricValues {
			fmt.Printf("%+v\n", v)
		}
	}

}
