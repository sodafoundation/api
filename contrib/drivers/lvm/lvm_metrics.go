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
	"strconv"
	"strings"
	"time"

	log "github.com/golang/glog"
	"github.com/opensds/opensds/pkg/model"
	"gopkg.in/yaml.v2"
)

// Supported metrics
var data = `
resources:
  - resource: volume
    metrics:
      - iops
      - read_throughput
      - write_throughput
      - response_time
      - service_time
      - utilization
    units:
      - tps
      - kbs
      - kbs
      - ms
      - ms
      - prcnt
  - resource: disk
    metrics:
      - iops
      - read_throughput
      - write_throughput
      - response_time
      - service_time
      - utilization
    units:
      - tps
      - kbs
      - kbs
      - ms
      - ms
      - prcnt`

type Config struct {
	Resource string
	Metrics  []string
	Units    []string
}

type Configs struct {
	Cfgs []Config `resources`
}
type MetricDriver struct {
	cli *MetricCli
}

func getCurrentUnixTimestamp() int64 {
	now := time.Now()
	secs := now.Unix()
	return secs
}
func getMetricToUnitMap() map[string]string {

	//construct metrics to value map
	var configs Configs
	//Read supported metric list from yaml config
	//Todo: Move this to read from file
	source := []byte(data)

	error := yaml.Unmarshal(source, &configs)
	if error != nil {
		log.Fatalf("unmarshal error: %v", error)
	}
	metricToUnitMap := make(map[string]string)
	for _, resources := range configs.Cfgs {
		switch resources.Resource {
		//ToDo: Other Cases needs to be added
		case "volume", "disk":
			for index, metricName := range resources.Metrics {

				metricToUnitMap[metricName] = resources.Units[index]

			}
		}
	}
	return metricToUnitMap
}

// 	getMetricList:- is  to get the list of supported metrics for given resource type
//	supportedMetrics -> list of supported metrics
func (d *MetricDriver) GetMetricList(resourceType string) (supportedMetrics []string, err error) {
	var configs Configs

	//Read supported metric list from yaml config
	source := []byte(data)
	error := yaml.Unmarshal(source, &configs)
	if error != nil {
		log.Fatalf("unmarshal error: %v", error)
	}

	for _, resources := range configs.Cfgs {
		if resources.Resource == resourceType {
			switch resourceType {
			case "volume", "disk":
				for _, m := range resources.Metrics {
					supportedMetrics = append(supportedMetrics, m)

				}
			}
		}
	}

	return supportedMetrics, nil
}

//	CollectMetrics: Driver entry point to collect metrics. This will be invoked by the dock
//	[]*model.MetricSpec	-> the array of metrics to be returned
func (d *MetricDriver) CollectMetrics() ([]*model.MetricSpec, error) {

	// get Metrics to unit map
	metricToUnitMap := getMetricToUnitMap()
	//validate metric support list
	supportedMetrics, err := d.GetMetricList("volume")
	if supportedMetrics == nil {
		log.Infof("no metrics found in the  supported metric list")
	}
	// discover lvm volumes
	volumeList, vgList, err := d.cli.DiscoverVolumes()
	if err != nil {
		log.Errorf("discover volume function returned error, err: %v", err)
	}
	// discover lvm physical volumes
	DiskList, err := d.cli.DiscoverDisks()
	if err != nil {
		log.Errorf("discover disk returned error, err: %v", err)
	}
	metricMap, labelMap, err := d.cli.CollectMetrics(supportedMetrics)
	if err != nil {
		log.Errorf("collect metrics returned error, err: %v", err)
	}
	var tempMetricArray []*model.MetricSpec
	// fill volume metrics
	for i, volume := range volumeList {
		convrtedVolID := convert(volume, vgList[i])
		aMetricMap := metricMap[convrtedVolID]
		aLabelMap := labelMap[convrtedVolID]
		for _, element := range supportedMetrics {
			val, _ := strconv.ParseFloat(aMetricMap[element], 64)
			metricValue := &model.Metric{
				Timestamp: getCurrentUnixTimestamp(),
				Value:     val,
			}
			metricValues := make([]*model.Metric, 0)
			metricValues = append(metricValues, metricValue)
			metric := &model.MetricSpec{
				InstanceID:   volume,
				InstanceName: aMetricMap["InstanceName"],
				Job:          "lvm",
				Labels:       aLabelMap,
				Component:    "volume",
				Name:         element,
				Unit:         metricToUnitMap[element],
				AggrType:     "",
				MetricValues: metricValues,
			}
			tempMetricArray = append(tempMetricArray, metric)
		}
	}
	// fill disk  metrics
	for _, disk := range DiskList {
		convrtedVolID := formatDiskName(disk)
		aMetricMap := metricMap[convrtedVolID]
		aLabelMap := labelMap[convrtedVolID]
		for _, element := range supportedMetrics {
			val, _ := strconv.ParseFloat(aMetricMap[element], 64)
			metricValue := &model.Metric{
				Timestamp: getCurrentUnixTimestamp(),
				Value:     val,
			}
			metricValues := make([]*model.Metric, 0)
			metricValues = append(metricValues, metricValue)
			metric := &model.MetricSpec{
				InstanceID:   disk,
				InstanceName: aMetricMap["InstanceName"],
				Job:          "lvm",
				Labels:       aLabelMap,
				Component:    "disk",
				Name:         element,
				Unit:         metricToUnitMap[element],
				AggrType:     "",
				MetricValues: metricValues,
			}
			tempMetricArray = append(tempMetricArray, metric)
		}
	}
	metricArray := tempMetricArray
	return metricArray, err
}
func convert(instanceID string, vg string) string {
	// systat utilities (sar/iostat) returnes  volume with -- instead of -, so we need to modify volume name to map lvs output
	instanceID = strings.Replace(instanceID, "-", "--", -1)
	vg = strings.Replace(vg, "-", "--", -1)
	//add opensds--volumes--default-- to the start of volume
	instanceID = vg + "-" + instanceID
	return instanceID
}
func formatDiskName(instanceID string) string {
	// systat(sar/iostat) returns only disk name. We need to add /dev/ to match with pvs output
	instanceID = strings.Replace(instanceID, "/dev/", "", -1)
	return instanceID
}
func (d *MetricDriver) Setup() error {

	cli, err := NewMetricCli()
	if err != nil {
		return err
	}
	d.cli = cli
	return nil
}

func (*MetricDriver) Teardown() error { return nil }
