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

package dorado

import (
	"strconv"
	"time"

	log "github.com/golang/glog"
	. "github.com/opensds/opensds/contrib/drivers/utils/config"
	"github.com/opensds/opensds/pkg/model"
	"github.com/opensds/opensds/pkg/utils/config"
	uuid "github.com/satori/go.uuid"
	"gopkg.in/yaml.v2"
)

/*
Naming Map:
metrics             --> OceanStor
iops                --> Throughput(IOPS)(IO/s)
bandwidth           --> Bandwidth(MB/s) / Block Bandwidth(MB/s)
latency             --> Average I/O Latency(us)
service_time        --> Service Time(Excluding Queue Time)(ms)
cache_hit_ratio     --> % Hit
cpu_usage           --> CPU Usage(%)
*/
// Todo: Move this Yaml config to a file
// Todo: Add resources for "volume", "disk" and "port".
var data = `
resources:
  - resource: pool
    metrics:
      - iops
      - bandwidth
      - latency
      - service_time
      - utilization_prcnt
    units:
      - tps
      - mbs
      - microsecond
      - ms
      - prcnt
  - resource: controller
    metrics:
      - iops
      - bandwidth
      - latency
      - service_time
      - cache_hit_ratio
      - cpu_usage
    units:
      - tps
      - mbs
      - microsecond
      - ms
      - prcnt
      - prcnt`

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
func getMetricToUnitMap(resourceType string) map[string]string {
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
		if resources.Resource == resourceType {
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
			for _, m := range resources.Metrics {
				supportedMetrics = append(supportedMetrics, m)
			}
			break
		}
	}

	return supportedMetrics, nil
}

func (d *MetricDriver) CollectPerformanceMetrics(resId string, metricList []string) (map[string]float64, error) {
	name2id := map[string]string{
		KMetricIOPS:               PerfIOPS,
		KMetricBandwidth:          PerfBandwidth,
		KMetricLatency:            PerfLatency,
		KMetricServiceTime:        PerfServiceTime,
		KMetricUtilizationPercent: PerfUtilizationPercent,
		KMetricCacheHitRatio:      PerfCacheHitRatio,
		KMetricCpuUsage:           PerfCpuUsage,
	}

	var idList = make([]string, len(metricList))
	for i, name := range metricList {
		idList[i] = name2id[name]
	}

	perfMap, err := d.client.GetPerformance(resId, idList)
	if err != nil {
		log.Errorf("get performance data failed: %s", err)
		return nil, err
	}

	var metricMap = make(map[string]float64)
	for _, name := range metricList {
		v, _ := strconv.ParseFloat(perfMap[name2id[name]], 64)
		metricMap[name] = v
	}
	return metricMap, nil
}

func (d *MetricDriver) CollectControllerMetrics() ([]*model.MetricSpec, error) {
	// get Metrics to unit map
	metricToUnitMap := getMetricToUnitMap(MetricResourceTypeController)
	//validate metric support list
	supportedMetrics, err := d.GetMetricList(MetricResourceTypeController)
	if supportedMetrics == nil {
		log.Infof("no metrics found in the supported metric list")
	}
	controllers, err := d.client.ListControllers()
	if err != nil {
		log.Errorf("get controller failed: %s", err)
		return nil, err
	}

	var tempMetricArray []*model.MetricSpec
	for _, controller := range controllers {
		// TODO: the controller id need to be optional
		resId := ObjectTypeController + ":" + controller.Id
		metricMap, err := d.CollectPerformanceMetrics(resId, supportedMetrics)
		if err != nil {
			log.Errorf("get performance data failed: %s", err)
			return nil, err
		}

		for _, element := range supportedMetrics {
			metricValue := &model.Metric{
				Timestamp: getCurrentUnixTimestamp(),
				Value:     metricMap[element],
			}
			metricValues := make([]*model.Metric, 0)
			metricValues = append(metricValues, metricValue)
			metric := &model.MetricSpec{
				InstanceID:   resId,
				InstanceName: resId,
				Job:          "HuaweiOceanStor",
				Labels:       map[string]string{"controller": resId},
				Component:    MetricResourceTypeController,
				Name:         element,
				Unit:         metricToUnitMap[element],
				AggrType:     "",
				MetricValues: metricValues,
			}
			tempMetricArray = append(tempMetricArray, metric)
		}
	}
	return tempMetricArray, nil
}

func (d *MetricDriver) CollectPoolMetrics() ([]*model.MetricSpec, error) {
	// get Metrics to unit map
	metricToUnitMap := getMetricToUnitMap(MetricResourceTypePool)
	//validate metric support list
	supportedMetrics, err := d.GetMetricList(MetricResourceTypePool)
	if supportedMetrics == nil {
		log.Infof("no metrics found in the supported metric list")
	}

	poolAll, err := d.client.ListStoragePools()
	if err != nil {
		log.Errorf("get controller failed: %s", err)
		return nil, err
	}
	// Filter unsupported pools
	var pools []StoragePool
	for _, p := range poolAll {
		if _, ok := d.conf.Pool[p.Name]; !ok {
			continue
		}
		pools = append(pools, p)
	}

	var tempMetricArray []*model.MetricSpec
	for _, pool := range pools {
		// TODO: the controller id need to be optional
		resId := ObjectTypePool + ":" + pool.Id
		metricMap, err := d.CollectPerformanceMetrics(resId, supportedMetrics)
		if err != nil {
			log.Errorf("get performance data failed: %s", err)
			return nil, err
		}
		poolId := uuid.NewV5(uuid.NamespaceOID, pool.Name).String()
		for _, element := range supportedMetrics {
			metricValue := &model.Metric{
				Timestamp: getCurrentUnixTimestamp(),
				Value:     metricMap[element],
			}
			metricValues := make([]*model.Metric, 0)
			metricValues = append(metricValues, metricValue)
			metric := &model.MetricSpec{
				InstanceID:   poolId,
				InstanceName: pool.Name,
				Job:          "HuaweiOceanStor",
				Labels:       map[string]string{"pool": poolId},
				Component:    MetricResourceTypePool,
				Name:         element,
				Unit:         metricToUnitMap[element],
				AggrType:     "",
				MetricValues: metricValues,
			}
			tempMetricArray = append(tempMetricArray, metric)
		}
	}
	return tempMetricArray, nil
}

//	CollectMetrics: Driver entry point to collect metrics. This will be invoked by the dock
//	[]*model.MetricSpec	-> the array of metrics to be returned
func (d *MetricDriver) CollectMetrics() ([]*model.MetricSpec, error) {
	var metricFunList = []func() ([]*model.MetricSpec, error){
		d.CollectControllerMetrics, d.CollectPoolMetrics,
	}

	var metricAll []*model.MetricSpec
	for _, f := range metricFunList {
		metric, err := f()
		if err != nil {
			log.Errorf("get metric failed: %v", err)
			return nil, err
		}
		metricAll = append(metricAll, metric...)
	}

	return metricAll, nil
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
