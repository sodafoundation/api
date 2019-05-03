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
package main

import (
	"fmt"
	log "github.com/golang/glog"
	"github.com/opensds/opensds/contrib/drivers/lvm"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"sync"
)

//struct for lvm  collector that contains pointers
//to prometheus descriptors for each metric we expose.
type lvmCollector struct {
	mu sync.Mutex
	//volume metrics
	VolumeIOPS            *prometheus.Desc
	VolumeReadThroughput  *prometheus.Desc
	VolumeWriteThroughput *prometheus.Desc
	VolumeResponseTime    *prometheus.Desc
	VolumeServiceTime     *prometheus.Desc
	VolumeUtilization     *prometheus.Desc
	//Disk metrics
	DiskIOPS            *prometheus.Desc
	DiskReadThroughput  *prometheus.Desc
	DiskWriteThroughput *prometheus.Desc
	DiskResponseTime    *prometheus.Desc
	DiskServiceTime     *prometheus.Desc
	DiskUtilization     *prometheus.Desc
}

/* rr */
//constructor for lvm collector that
//initializes every descriptor and returns a pointer to the collector
func newLvmCollector() *lvmCollector {
	var labelKeys = []string{"device"}

	return &lvmCollector{
		VolumeIOPS: prometheus.NewDesc("OpensSDS_Volume_IOPS_tps",
			"Shows IOPS",
			labelKeys, nil,
		),
		VolumeReadThroughput: prometheus.NewDesc("OpensSDS_Volume_ReadThroughput_KBs",
			"Shows ReadThroughput",
			labelKeys, nil,
		),
		VolumeWriteThroughput: prometheus.NewDesc("OpensSDS_Volume_WriteThroughput_KBs",
			"Shows ReadThroughput",
			labelKeys, nil,
		),
		VolumeResponseTime: prometheus.NewDesc("OpensSDS_Volume_ResponseTime_ms",
			"Shows ReadThroughput",
			labelKeys, nil,
		),
		VolumeServiceTime: prometheus.NewDesc("OpensSDS_Volume_ServiceTime_ms",
			"Shows ServiceTime",
			labelKeys, nil,
		),
		VolumeUtilization: prometheus.NewDesc("OpensSDS_Volume_Utilization_prcnt",
			"Shows Utilization in percentage",
			labelKeys, nil,
		),
		DiskIOPS: prometheus.NewDesc("OpensSDS_Volume_IOPS_tps",
			"Shows IOPS",
			labelKeys, nil,
		),
		DiskReadThroughput: prometheus.NewDesc("OpensSDS_Disk_ReadThroughput_KBs",
			"Shows ReadThroughput",
			labelKeys, nil,
		),
		DiskWriteThroughput: prometheus.NewDesc("OpensSDS_Disk_WriteThroughput_KBs",
			"Shows ReadThroughput",
			labelKeys, nil,
		),
		DiskResponseTime: prometheus.NewDesc("OpensSDS_Disk_ResponseTime_ms",
			"Shows ReadThroughput",
			labelKeys, nil,
		),
		DiskServiceTime: prometheus.NewDesc("OpensSDS_Disk_ServiceTime_ms",
			"Shows ServiceTime",
			labelKeys, nil,
		),
		DiskUtilization: prometheus.NewDesc("OpensSDS_Disk_Utilization_prcnt",
			"Shows Utilization in percentage",
			labelKeys, nil,
		),
	}

}

//Describe function.
//It essentially writes all descriptors to the prometheus desc channel.
func (c *lvmCollector) Describe(ch chan<- *prometheus.Desc) {

	//Update this section with the each metric you create for a given collector
	ch <- c.VolumeIOPS
	ch <- c.VolumeReadThroughput
	ch <- c.VolumeWriteThroughput
	ch <- c.VolumeResponseTime
	ch <- c.VolumeServiceTime
	ch <- c.VolumeUtilization
	ch <- c.DiskIOPS
	ch <- c.DiskReadThroughput
	ch <- c.DiskWriteThroughput
	ch <- c.DiskResponseTime
	ch <- c.DiskServiceTime
	ch <- c.DiskUtilization
}

type Config struct {
	Type    string   `type`
	Devices []string `devices`
}

type Configs struct {
	Cfgs []Config `resources`
}

//Collect implements required collect function for all promehteus collectors
func (c *lvmCollector) Collect(ch chan<- prometheus.Metric) {

	c.mu.Lock()
	defer c.mu.Unlock()

	//Implement logic here to determine proper metric value to return to prometheus
	//for each descriptor
	metricList := []string{"IOPS", "ReadThroughput", "WriteThroughput", "ResponseTime", "ServiceTime", "UtilizationPercentage"}
	filename := "resources.yaml"
	source, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal("file "+filename+"can't read", err)
	}
	var config Configs
	err1 := yaml.Unmarshal(source, &config)
	if err1 != nil {
		log.Fatalf("error: %v", err)
	}

	//volumeList := []string{"sda", "loop0"}
	metricDriver := lvm.MetricDriver{}
	metricDriver.Setup()
	for _, resource := range config.Cfgs {
		switch resource.Type {
		case "volume":
			for _, volume := range resource.Devices {
				metricArray, _ := metricDriver.CollectMetrics(metricList, volume)
				fmt.Println(metricArray)
				for _, metric := range metricArray {
					lableVals := []string{metric.InstanceName}
					//unitLabel := "Unit:"+metric.Unit
					switch metric.Name {
					case "IOPS":
						ch <- prometheus.MustNewConstMetric(c.VolumeIOPS, prometheus.GaugeValue, metric.MetricValues[0].Value, lableVals...)
					case "ReadThroughput":
						ch <- prometheus.MustNewConstMetric(c.VolumeReadThroughput, prometheus.GaugeValue, metric.MetricValues[0].Value, lableVals...)
					case "WriteThroughput":
						ch <- prometheus.MustNewConstMetric(c.VolumeWriteThroughput, prometheus.GaugeValue, metric.MetricValues[0].Value, lableVals...)
					case "ResponseTime":
						ch <- prometheus.MustNewConstMetric(c.VolumeResponseTime, prometheus.GaugeValue, metric.MetricValues[0].Value, lableVals...)
					case "ServiceTime":
						ch <- prometheus.MustNewConstMetric(c.VolumeServiceTime, prometheus.GaugeValue, metric.MetricValues[0].Value, lableVals...)

					case "UtilizationPercentage":
						ch <- prometheus.MustNewConstMetric(c.VolumeUtilization, prometheus.GaugeValue, metric.MetricValues[0].Value, lableVals...)

					}
				}
			}

		case "disk":
			for _, volume := range resource.Devices {
				metricArray, _ := metricDriver.CollectMetrics(metricList, volume)
				fmt.Println(metricArray)
				for _, metric := range metricArray {
					lableVals := []string{metric.Labels["device"]}
					//unitLabel := "Unit:"+metric.Unit
					switch metric.Name {
					case "IOPS":
						ch <- prometheus.MustNewConstMetric(c.DiskIOPS, prometheus.GaugeValue, metric.MetricValues[0].Value, lableVals...)
					case "ReadThroughput":
						ch <- prometheus.MustNewConstMetric(c.DiskReadThroughput, prometheus.GaugeValue, metric.MetricValues[0].Value, lableVals...)
					case "WriteThroughput":
						ch <- prometheus.MustNewConstMetric(c.DiskWriteThroughput, prometheus.GaugeValue, metric.MetricValues[0].Value, lableVals...)
					case "ResponseTime":
						ch <- prometheus.MustNewConstMetric(c.DiskResponseTime, prometheus.GaugeValue, metric.MetricValues[0].Value, lableVals...)
					case "ServiceTime":
						ch <- prometheus.MustNewConstMetric(c.DiskServiceTime, prometheus.GaugeValue, metric.MetricValues[0].Value, lableVals...)
					case "UtilizationPercentage":
						ch <- prometheus.MustNewConstMetric(c.DiskUtilization, prometheus.GaugeValue, metric.MetricValues[0].Value, lableVals...)

					}
				}
			}
		}

	}

}
func validateCliArg(arg1 string) string {
	num, err := strconv.Atoi(arg1)
	if (err != nil) || (num > 65535) {

		fmt.Println("please enter a valid port number")
		os.Exit(1)
	}
	return arg1
}

// main function for lvm exporter
// lvm exporter is a independent process which user can start if required
func main() {

	portNo := validateCliArg(os.Args[1])

	//Create a new instance of the lvmcollector and
	//register it with the prometheus client.
	lvm := newLvmCollector()
	prometheus.MustRegister(lvm)

	//This section will start the HTTP server and expose
	//any metrics on the /metrics endpoint.
	http.Handle("/metrics", promhttp.Handler())
	log.Info("lvm exporter begining to serve on port :" + portNo)
	log.Fatal(http.ListenAndServe(":"+portNo, nil))
}
