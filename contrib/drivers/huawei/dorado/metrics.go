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
package dorado

import (
	"strconv"
	"time"

	log "github.com/golang/glog"
	. "github.com/opensds/opensds/contrib/drivers/utils/config"
	"github.com/opensds/opensds/pkg/model"
	"github.com/opensds/opensds/pkg/utils/config"
	"gopkg.in/yaml.v2"
)

// Todo: Move this Yaml config to a file

var data = `
resources:
  - resource: volume
    metrics:
      - iops
      - read_throughput
      - write_throughput
      - response_time
      - service_time
      - utilization_prcnt
    units:
      - tps
      - kbs
      - kbs
      - ms
      - ms
      - '%'
  - resource: pool
    metrics:
      - iops
      - read_throughput
      - write_throughput
      - response_time
      - service_time
      - utilization_prcnt
    units:
      - tps
      - kbs
      - kbs
      - ms
      - ms
      - '%'
  - resource: controller
    metrics:
      - iops
      - read_throughput
      - write_throughput
      - response_time
      - service_time
      - utilization_prcnt
    units:
      - tps
      - kbs
      - kbs
      - ms
      - ms
      - '%'`

type Config struct {
	Resource string
	Metrics  []string
	Units    []string
}

type Configs struct {
	Cfgs []Config `yaml:"resources"`
}
type MetricDriver struct {
	conf   *DoradoConfig
	client *DoradoClient
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

	err := yaml.Unmarshal(source, &configs)
	if err != nil {
		log.Fatalf("unmarshal error: %v", err)
	}
	metricToUnitMap := make(map[string]string)
	for _, resources := range configs.Cfgs {
		switch resources.Resource {
		//ToDo: Other Cases needs to be added
		case "volume", "pool", "controller":
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
	err = yaml.Unmarshal(source, &configs)
	if err != nil {
		log.Fatalf("unmarshal error: %v", err)
	}

	for _, resources := range configs.Cfgs {
		if resources.Resource == resourceType {
			switch resourceType {
			case "volume", "pool", "controller":
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
	controllers, err := d.client.ListControllers()
	if err != nil {
		log.Errorf("get controller failed: %s", err)
		return nil, err
	}

	var tempMetricArray []*model.MetricSpec
	for _, controller := range controllers {
		// TODO: the controller id need to be optional
		uuid := ObjectTypeController + ":" + controller.Id
		dataIdList := []string{PerfIOPS, PerfReadThroughput, PerfWriteThroughput,
			PerfResponseTime, PerfServiceTime, PerfUtilizationPrcnt}
		metricMap, err := d.client.GetPerformance(uuid, dataIdList)
		if err != nil {
			log.Errorf("get performance data failed: %s", err)
			return nil, err
		}

		name2id := map[string]string{
			"iops":              PerfIOPS,
			"read_throughput":   PerfReadThroughput,
			"write_throughput":  PerfWriteThroughput,
			"response_time":     PerfResponseTime,
			"service_time":      PerfServiceTime,
			"utilization_prcnt": PerfUtilizationPrcnt,
		}
		for _, element := range supportedMetrics {
			val, _ := strconv.ParseFloat(metricMap[name2id[element]], 64)
			metricValue := &model.Metric{
				Timestamp: getCurrentUnixTimestamp(),
				Value:     val,
			}
			metricValues := make([]*model.Metric, 0)
			metricValues = append(metricValues, metricValue)
			metric := &model.MetricSpec{
				InstanceID:   uuid,
				InstanceName: uuid,
				Job:          "HuaweiOceanStor",
				Labels:       map[string]string{},
				Component:    "controller",
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

func (d *MetricDriver) Setup() (err error) {
	// Read huawei dorado config file
	path := config.CONF.OsdsDock.Backends.HuaweiDorado.ConfigPath
	if "" == path {
		path = defaultConfPath
	}

	conf := &DoradoConfig{}
	Parse(conf, path)

	d.conf = conf
	d.client, err = NewClient(&d.conf.AuthOptions)
	if err != nil {
		log.Errorf("Get new client failed, %v", err)
		return err
	}
	return nil
}

func (*MetricDriver) Teardown() error { return nil }
