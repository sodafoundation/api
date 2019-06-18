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
package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"

	log "github.com/golang/glog"
	"github.com/opensds/opensds/contrib/drivers/lvm"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const DefaultPort = "8080"

// struct for lvm  collector that contains pointers
// to prometheus descriptors for each metric we expose.
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

// constructor for lvm collector that
// initializes every descriptor and returns a pointer to the collector
func newLvmCollector() *lvmCollector {
	var labelKeys = []string{"device"}

	return &lvmCollector{
		VolumeIOPS: prometheus.NewDesc("lvm_volume_iops_tps",
			"Shows IOPS",
			labelKeys, nil,
		),
		VolumeReadThroughput: prometheus.NewDesc("lvm_volume_read_throughput_kbs",
			"Shows ReadThroughput",
			labelKeys, nil,
		),
		VolumeWriteThroughput: prometheus.NewDesc("lvm_volume_write_throughput_kbs",
			"Shows ReadThroughput",
			labelKeys, nil,
		),
		VolumeResponseTime: prometheus.NewDesc("lvm_volume_response_time_ms",
			"Shows ReadThroughput",
			labelKeys, nil,
		),
		VolumeServiceTime: prometheus.NewDesc("lvm_volume_service_time_ms",
			"Shows ServiceTime",
			labelKeys, nil,
		),
		VolumeUtilization: prometheus.NewDesc("lvm_volume_utilization_prcnt",
			"Shows Utilization in percentage",
			labelKeys, nil,
		),
		DiskIOPS: prometheus.NewDesc("lvm_disk_iops_tps",
			"Shows IOPS",
			labelKeys, nil,
		),
		DiskReadThroughput: prometheus.NewDesc("lvm_disk_read_throughput_kbs",
			"Shows Disk ReadThroughput",
			labelKeys, nil,
		),
		DiskWriteThroughput: prometheus.NewDesc("lvm_disk_write_throughput_kbs",
			"Shows Write Throughput",
			labelKeys, nil,
		),
		DiskResponseTime: prometheus.NewDesc("lvm_disk_response_time_ms",
			"Shows Disk Response Time",
			labelKeys, nil,
		),
		DiskServiceTime: prometheus.NewDesc("lvm_disk_service_time_ms",
			"Shows ServiceTime",
			labelKeys, nil,
		),
		DiskUtilization: prometheus.NewDesc("lvm_disk_utilization_prcnt",
			"Shows Utilization in percentage",
			labelKeys, nil,
		),
	}

}

// Describe function.
// It essentially writes all descriptors to the prometheus desc channel.
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
	Cfgs []*Config `resources`
}

// Collect implements required collect function for all promehteus collectors
func (c *lvmCollector) Collect(ch chan<- prometheus.Metric) {

	c.mu.Lock()
	defer c.mu.Unlock()

	metricDriver := lvm.MetricDriver{}
	metricDriver.Setup()

	metricArray, _ := metricDriver.CollectMetrics()
	for _, metric := range metricArray {
		lableVals := []string{metric.InstanceName}
		switch metric.Component {
		case "volume":
			switch metric.Name {
			case "iops":
				ch <- prometheus.MustNewConstMetric(c.VolumeIOPS, prometheus.GaugeValue, metric.MetricValues[0].Value, lableVals...)
			case "read_throughput":
				ch <- prometheus.MustNewConstMetric(c.VolumeReadThroughput, prometheus.GaugeValue, metric.MetricValues[0].Value, lableVals...)
			case "write_throughput":
				ch <- prometheus.MustNewConstMetric(c.VolumeWriteThroughput, prometheus.GaugeValue, metric.MetricValues[0].Value, lableVals...)
			case "response_time":
				ch <- prometheus.MustNewConstMetric(c.VolumeResponseTime, prometheus.GaugeValue, metric.MetricValues[0].Value, lableVals...)
			case "service_time":
				ch <- prometheus.MustNewConstMetric(c.VolumeServiceTime, prometheus.GaugeValue, metric.MetricValues[0].Value, lableVals...)
			case "utilization_prcnt":
				ch <- prometheus.MustNewConstMetric(c.VolumeUtilization, prometheus.GaugeValue, metric.MetricValues[0].Value, lableVals...)

			}
		case "disk":
			switch metric.Name {
			case "iops":
				ch <- prometheus.MustNewConstMetric(c.DiskIOPS, prometheus.GaugeValue, metric.MetricValues[0].Value, lableVals...)
			case "read_throughput":
				ch <- prometheus.MustNewConstMetric(c.DiskReadThroughput, prometheus.GaugeValue, metric.MetricValues[0].Value, lableVals...)
			case "write_throughput":
				ch <- prometheus.MustNewConstMetric(c.DiskWriteThroughput, prometheus.GaugeValue, metric.MetricValues[0].Value, lableVals...)
			case "response_time":
				ch <- prometheus.MustNewConstMetric(c.DiskResponseTime, prometheus.GaugeValue, metric.MetricValues[0].Value, lableVals...)
			case "service_time":
				ch <- prometheus.MustNewConstMetric(c.DiskServiceTime, prometheus.GaugeValue, metric.MetricValues[0].Value, lableVals...)
			case "utilization_prcnt":
				ch <- prometheus.MustNewConstMetric(c.DiskUtilization, prometheus.GaugeValue, metric.MetricValues[0].Value, lableVals...)

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

	var portNo string
	if len(os.Args) > 1 {
		portNo = validateCliArg(os.Args[1])
	} else {
		portNo = DefaultPort
	}

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
