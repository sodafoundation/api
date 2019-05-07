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
	"strconv"
	"time"

	log "github.com/golang/glog"
	"github.com/opensds/opensds/pkg/model"
	"gopkg.in/yaml.v2"
)

// Todo: Move this Yaml config to a file

var data = `
resources:
  - resource: volume
    metrics:
      - IOPS
      - ReadThroughput
      - WriteThroughput
      - ResponseTime
      - ServiceTime
      - UtilizationPercentage
    units:
      - tps
      - KB/s
      - KB/s
      - ms
      - ms
      - '%'
  - resource: pool
    metrics:
      - ReadRequests
      - WriteRequests
      - ReponseTime
    units:
      - tps
      - KB/s
      - KB/s
      - ms
      - ms
      - '%'
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

	//construct metrics to value map
	var configs Configs
	//Read supported metric list from yaml config
	//Todo: Move this to read from file
	source := []byte(data)

	error := yaml.Unmarshal(source, &configs)
	if error != nil {
		log.Fatalf("Unmarshal error: %v", error)
	}
	metricToUnitMap := make(map[string]string)
	for _, resources := range configs.Cfgs {
		switch resources.Resource {
		//ToDo: Other Cases needs to be added
		case "volume":
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

	//Read supported metric list from yaml config
	//Todo: Move this to read from file
	source := []byte(data)
	error := yaml.Unmarshal(source, &configs)
	if error != nil {
		log.Fatalf("Unmarshal error: %v", error)
	}

	for _, resources := range configs.Cfgs {
		switch resources.Resource {
		//ToDo: Other Cases needs to be added
		case "volume":
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
	supportedMetrics, err := d.ValidateMetricsSupportList(metricsList, "volume")
	if supportedMetrics == nil {
		log.Infof("No metrics found in the  supported metric list")
	}
	metricMap, err := d.cli.CollectMetrics(supportedMetrics, instanceID)

	var tempMetricArray []*model.MetricSpec
	for _, element := range metricsList {
		val, _ := strconv.ParseFloat(metricMap[element], 64)
		//Todo: See if association  is required here, resource discovery could fill this information
		associatorMap := make(map[string]string)
		associatorMap["device"] = metricMap["InstanceName"]
		metricValue := &model.Metric{
			Timestamp: getCurrentUnixTimestamp(),
			Value:     val,
		}
		metricValues := make([]*model.Metric, 0)
		metricValues = append(metricValues, metricValue)

		metric := &model.MetricSpec{
			InstanceID:   instanceID,
			InstanceName: metricMap["InstanceName"],
			Job:          "OpenSDS",
			Labels:       associatorMap,
			//Todo Take Componet from Post call, as of now it is only for volume
			Component: "Volume",
			Name:      element,
			//Todo : Fill units according to metric type
			Unit: metricToUnitMap[element],
			//Todo : Get this information dynamically ( hard coded now , as all are direct values
			AggrType:     "",
			MetricValues: metricValues,
		}
		tempMetricArray = append(tempMetricArray, metric)
	}
	metricArray := tempMetricArray
	return metricArray, err
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
