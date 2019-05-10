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
package ceph

import (
	"fmt"
	log "github.com/golang/glog"
	"github.com/opensds/opensds/pkg/model"
	"gopkg.in/yaml.v2"
	"strconv"
	"time"
)

// TODO: Move this Yaml config to a file
var data = `
resources:
  - resource: pool
    metrics:
      - pool_used_bytes
      - pool_raw_used_bytes
      - pool_available_bytes
      - pool_objects_total
      - pool_dirty_objects_total
      - pool_read_total
      - pool_read_bytes_total
      - pool_write_total
      - pool_write_bytes_total
    units:
      - bytes
      - bytes
      - bytes
      - ""
      - ""
      - ""
      - bytes
      - ""
      - bytes
`

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

func metricInMetrics(metric string, metriclist []string) bool {
	for _, m := range metriclist {
		if m == metric {
			return true
		}
	}
	return false
}

func getCurrentUnixTimestamp() int64 {
	now := time.Now()
	secs := now.Unix()
	return secs
}

func getMetricToUnitMap() map[string]string {

	// Construct metrics to value map
	var configs Configs
	// Read supported metric list from yaml config
	// TODO: Move this to read from file
	source := []byte(data)

	error := yaml.Unmarshal(source, &configs)
	if error != nil {
		log.Fatalf("Unmarshal error: %v", error)
	}
	metricToUnitMap := make(map[string]string)
	for _, resources := range configs.Cfgs {
		switch resources.Resource {
		// TODO: Other Cases needs to be added
		case "pool":
			for index, metricName := range resources.Metrics {

				metricToUnitMap[metricName] = resources.Units[index]

			}
		}
	}
	return metricToUnitMap
}

// 	ValidateMetricsSupportList:- is  to check whether the posted metric list is in the uspport list of this driver
// 	metricList-> Posted metric list
//	supportedMetrics -> list of supported metrics
func (d *MetricDriver) ValidateMetricsSupportList(metricList []string, resourceType string) (supportedMetrics []string, err error) {
	var configs Configs

	// Read supported metric list from yaml config
	// TODO: Move this to read from file
	source := []byte(data)
	error := yaml.Unmarshal(source, &configs)
	if error != nil {
		log.Fatalf("Unmarshal error: %v", error)
	}

	for _, resources := range configs.Cfgs {
		switch resources.Resource {
		// TODO: Other Cases needs to be added
		case "pool":
			for _, metricName := range metricList {
				if metricInMetrics(metricName, resources.Metrics) {
					supportedMetrics = append(supportedMetrics, metricName)

				} else {
					log.Infof("metric:%s is not in the supported list", metricName)
				}
			}
		}
	}
	return supportedMetrics, nil
}

//	CollectMetrics: Driver entry point to collect metrics. This will be invoked by the dock
//	metricsList-> posted metric list
//	instanceID -> posted instanceID
//	metricArray	-> the array of metrics to be returned
func (d *MetricDriver) CollectMetrics(metricsList []string, instanceID string) ([]*model.MetricSpec, error) {

	// get Metrics to unit map
	metricToUnitMap := getMetricToUnitMap()
	//validate metric support list
	supportedMetrics, err := d.ValidateMetricsSupportList(metricsList, "pool")
	if supportedMetrics == nil {
		log.Infof("No metrics found in the  supported metric list")
	}
	metricMap, err := d.cli.CollectMetrics(supportedMetrics, instanceID)

	var tempMetricArray []*model.MetricSpec
	for label_val, _ := range metricMap {
		for _, element := range supportedMetrics {
			val, _ := strconv.ParseFloat(metricMap[label_val][element], 64)
			// TODO: See if association  is required here, resource discovery could fill this information
			associatorMap := make(map[string]string)
			associatorMap["cluster"] = "ceph"
			associatorMap["pool"] = label_val
			metricValue := &model.Metric{
				Timestamp: getCurrentUnixTimestamp(),
				Value:     val,
			}
			metricValues := make([]*model.Metric, 0)
			metricValues = append(metricValues, metricValue)

			metric := &model.MetricSpec{
				InstanceID:   instanceID,
				InstanceName: "",
				Job:          "OpenSDS",
				Labels:       associatorMap,
				// TODO Take Componet from Post call, as of now it is only for volume
				Component: "Pool",
				Name:      fmt.Sprintf("%s_%s", associatorMap["cluster"], element),
				// TODO : Fill units according to metric type
				Unit: metricToUnitMap[element],
				// TODO : Get this information dynamically ( hard coded now , as all are direct values
				AggrType:     "",
				MetricValues: metricValues,
			}
			tempMetricArray = append(tempMetricArray, metric)
		}
	}
	metricArray := tempMetricArray
	return metricArray, err
}

func (d *MetricDriver) Setup() error {

	return nil
}

func (*MetricDriver) Teardown() error { return nil }
