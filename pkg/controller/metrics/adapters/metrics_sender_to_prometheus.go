// Copyright (c) 2019 The OpenSDS Authors All Rights Reserved.
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
package adapters

import (
	"os"
	"strconv"
	"time"

	log "github.com/golang/glog"
	"github.com/opensds/opensds/pkg/model"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
)

var nodeExporterFolder = "/root/prom_nodeexporter_folder/"

type PrometheusMetricsSender struct {
	Queue    chan *model.MetricSpec
	QuitChan chan bool
}

func (p *PrometheusMetricsSender) GetMetricsSender() MetricsSenderIntf {
	sender := PrometheusMetricsSender{}
	sender.Queue = make(chan *model.MetricSpec)
	sender.QuitChan = make(chan bool)
	return &sender
}

func (p *PrometheusMetricsSender) Start() {
	go func() {
		for {
			select {
			case work := <-p.Queue:
				// Receive a work request.
				log.Infof("GetMetricsSenderToPrometheus received metrics for instance %s\n and metrics %f\n", work.InstanceID, work.MetricValues[0].Value)

				// do the actual sending work here, by writing to the file of the node_exporter of prometheus
				writeToFile(work)

				// alternatively, we could also push the metrics to the push gateway of prometheus
				sendToPushGateway(work)

				log.Info("GetMetricsSenderToPrometheus processed metrics")

			case <-p.QuitChan:
				return
			}

		}
	}()
}
func (p *PrometheusMetricsSender) Stop() {
	go func() {
		p.QuitChan <- true
	}()
}

func (p *PrometheusMetricsSender) AssignMetricsToSend(request *model.MetricSpec) {
	p.Queue <- request
}

func writeToFile(metrics *model.MetricSpec) {

	// get the string ready to be written
	var finalString = ""

	finalString += metrics.Name + " " + strconv.FormatFloat(metrics.MetricValues[0].Value, 'f', 2, 64) + "\n"

	// make a new file with current timestamp
	timeStamp := strconv.FormatInt(time.Now().Unix(), 10)
	f, err := os.Create(nodeExporterFolder + metrics.InstanceID + ".prom")
	if err != nil {
		log.Error(err)
		return
	}
	_, err = f.WriteString(finalString)
	if err != nil {
		log.Error(err)
		f.Close()
		return
	}
	log.Infof("metrics written successfully at time %s", timeStamp)
	err = f.Close()
	if err != nil {
		log.Error(err)
		return
	}
}

func sendToPushGateway(metrics *model.MetricSpec) {

	completionTime := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: metrics.Name,
		Help: "",
	})
	completionTime.SetToCurrentTime()
	completionTime.Set(metrics.MetricValues[0].Value)

	if err := push.New("http://localhost:9091", "push_gateway").
		Collector(completionTime).
		Grouping("l1", "v1").
		Push(); err != nil {
		log.Errorf("Could not push completion time to Pushgateway:%s", err)
	}
	log.Info("Completed push completion time to Pushgateway")
}
