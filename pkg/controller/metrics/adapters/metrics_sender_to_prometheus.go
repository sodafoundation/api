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
package adapters

import (
	"os"
	"strconv"
	"time"

	log "github.com/golang/glog"
	"github.com/opensds/opensds/pkg/model"
	. "github.com/opensds/opensds/pkg/utils/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
)

var nodeExporterFolder = CONF.OsdsLet.NodeExporterWatchFolder

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

				if CONF.OsdsLet.PrometheusPushMechanism == "NodeExporter" {
					// do the actual sending work here, by writing to the file of the node_exporter of prometheus
					writeToFile(work)
					log.Info("GetMetricsSenderToPrometheus processed metrics write to node exporter")
				} else if CONF.OsdsLet.PrometheusPushMechanism == "PushGateway" {
					// alternatively, we could also push the metrics to the push gateway of prometheus
					sendToPushGateway(work)
					log.Info("GetMetricsSenderToPrometheus processed metrics send to push gateway")
				}

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

	// form the label string
	labelStr := "{"
	for labelName, labelValue := range metrics.Labels {
		labelStr = labelStr + labelName + "=" + `"` + labelValue + `"`
		labelStr = labelStr + ","
	}
	// replace the last , with } to complete the set of labels
	labelStr = labelStr[:len(labelStr)-1] + "}"

	metricName := getMetricName(metrics)

	finalString += metricName + labelStr + " " + strconv.FormatFloat(metrics.MetricValues[0].Value, 'f', 2, 64) + "\n"

	// make a new file with current timestamp
	timeStamp := strconv.FormatInt(time.Now().Unix(), 10)
	// form the temp file name
	tempFName := nodeExporterFolder + metricName + ".prom.temp"
	// form the actual file name
	fName := nodeExporterFolder + metricName + ".prom"

	// write to the temp file
	f, err := os.Create(tempFName)
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
	log.Infoln(finalString)
	err = f.Close()
	if err != nil {
		log.Error(err)
		return
	}
	// this is done so that the exporter never sees an incomplete file
	renameErr := os.Rename(tempFName, fName)
	if renameErr != nil {
		log.Errorf("error %s renaming metrics file %s to %s", renameErr.Error(), tempFName, fName)
	}
}

func getMetricName(metrics *model.MetricSpec) string {
	// form the full metric name
	metricName := metrics.Job + "_" + metrics.Component + "_" + metrics.Name + "_" + metrics.Unit
	// is this an aggregated metric ? add the aggregation type
	if metrics.AggrType != "" {
		metricName = metricName + "_" + metrics.AggrType
	}
	return metricName
}

func sendToPushGateway(metrics *model.MetricSpec) {

	metricName := getMetricName(metrics)

	gauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: metricName,
		Help: "",
	})
	gauge.SetToCurrentTime()
	gauge.Set(metrics.MetricValues[0].Value)

	pusher := push.New(CONF.OsdsLet.PushGatewayUrl, "push_gateway").
		Collector(gauge)
	for lKey, lValue := range metrics.Labels {
		pusher = pusher.Grouping(lKey, lValue)
	}
	// add the metric name here, to differentiate between various metrics
	pusher = pusher.Grouping("metricname", metricName)

	if err := pusher.Push(); err != nil {
		log.Errorf("error when pushing gauge for metric name=%s;timestamp=%v:value=%v to Pushgateway:%s", metricName, metrics.MetricValues[0].Timestamp, metrics.MetricValues[0].Value, err)

	}
	log.Infof("completed push gauge for metric name=%s;timestamp=%v:value=%v to Pushgateway", metricName, metrics.MetricValues[0].Timestamp, metrics.MetricValues[0].Value)

}
