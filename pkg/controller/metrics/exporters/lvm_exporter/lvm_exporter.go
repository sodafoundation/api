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
	"net/http"
	"os"
	"strconv"
	"sync"
)

//struct for lvm  collector that contains pointers
//to prometheus descriptors for each metric we expose.
type lvmCollector struct {
	mu         sync.Mutex
	//volume metrics
	IOPS            *prometheus.Desc
	ReadThroughput  *prometheus.Desc
	WriteThroughput *prometheus.Desc
	ResponseTime    *prometheus.Desc
	ServiceTime     *prometheus.Desc
	Utilization     *prometheus.Desc
}

/* rr */
//constructor for lvm collector that
//initializes every descriptor and returns a pointer to the collector
func newLvmCollector() *lvmCollector {
	var volumeLabel = []string{"volume"}
	return &lvmCollector{
		IOPS: prometheus.NewDesc("OpensSDS_Volume_IOPS_tps",
			"Shows IOPS",
			volumeLabel, nil,
		),
		ReadThroughput: prometheus.NewDesc("OpensSDS_Volume_ReadThroughput_KBs",
			"Shows ReadThroughput",
			volumeLabel, nil,
		),
		WriteThroughput: prometheus.NewDesc("OpensSDS_Volume_WriteThroughput_KBs",
			"Shows ReadThroughput",
			volumeLabel, nil,
		),
		ResponseTime: prometheus.NewDesc("OpensSDS_Volume_ResponseTime_ms",
			"Shows ReadThroughput",
			volumeLabel, nil,
		),
		ServiceTime: prometheus.NewDesc("OpensSDS_Volume_ServiceTime_ms",
			"Shows ServiceTime",
			volumeLabel, nil,
		),
		Utilization: prometheus.NewDesc("OpensSDS_Volume_Utilization_prcnt",
			"Shows Utilization in percentage",
			volumeLabel, nil,
		),
	}

}

//Describe function.
//It essentially writes all descriptors to the prometheus desc channel.
func (c *lvmCollector) Describe(ch chan<- *prometheus.Desc) {

	//Update this section with the each metric you create for a given collector
	ch <- c.IOPS
	ch <- c.ReadThroughput
	ch <- c.WriteThroughput
	ch <- c.ResponseTime
	ch <- c.ServiceTime
	ch <- c.Utilization
}

//Collect implements required collect function for all promehteus collectors
func (c *lvmCollector) Collect(ch chan<- prometheus.Metric) {

	c.mu.Lock()
	defer c.mu.Unlock()

	//Implement logic here to determine proper metric value to return to prometheus
	//for each descriptor
	metricList := []string{"IOPS", "ReadThroughput", "WriteThroughput", "ResponseTime", "ServiceTime", "UtilizationPercentage"}
	//Todo : Need to read list from a config file
	volumeList := []string{"sda", "loop0"}
	metricDriver := lvm.MetricDriver{}
	metricDriver.Setup()
	for _, volume := range volumeList {
		metricArray, _ := metricDriver.CollectMetrics(metricList, volume)
		fmt.Println(metricArray)
		for _, metric := range metricArray {
			instanceLabel := metric.InstanceID
			//unitLabel := "Unit:"+metric.Unit
			switch metric.Name {
			case "IOPS":
				ch <- prometheus.MustNewConstMetric(c.IOPS, prometheus.GaugeValue, metric.MetricValues[0].Value, instanceLabel)
			case "ReadThroughput":
				ch <- prometheus.MustNewConstMetric(c.ReadThroughput, prometheus.GaugeValue, metric.MetricValues[0].Value, instanceLabel)
			case "WriteThroughput":
				ch <- prometheus.MustNewConstMetric(c.WriteThroughput, prometheus.GaugeValue, metric.MetricValues[0].Value, instanceLabel)
			case "ResponseTime":
				ch <- prometheus.MustNewConstMetric(c.ResponseTime, prometheus.GaugeValue, metric.MetricValues[0].Value, instanceLabel)
			case "ServiceTime":
				ch <- prometheus.MustNewConstMetric(c.ServiceTime, prometheus.GaugeValue, metric.MetricValues[0].Value, instanceLabel)

			case "UtilizationPercentage":
				ch <- prometheus.MustNewConstMetric(c.Utilization, prometheus.GaugeValue, metric.MetricValues[0].Value, instanceLabel)

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
