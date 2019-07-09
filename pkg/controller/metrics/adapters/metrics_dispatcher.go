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
	log "github.com/golang/glog"
	"github.com/opensds/opensds/pkg/model"
)

// A buffered channel that we can send work requests on.
var MetricsQueue = make(chan *model.MetricSpec, 100)

func SendMetricToRegisteredSenders(metrics *model.MetricSpec) {

	// Push the work onto the queue.
	MetricsQueue <- metrics
	log.Info("Send metrics request queued")

	return
}

func StartDispatcher() {

	listMetricSenders := make([]MetricsSenderIntf, 0)

	// initialize Prometheus sender
	senderStructProm := PrometheusMetricsSender{}
	promMetricSender := senderStructProm.GetMetricsSender()

	senderStructKafka := KafkaMetricsSender{}
	kafkaMetricsSender := senderStructKafka.GetMetricsSender()

	// add to list
	listMetricSenders = append(listMetricSenders, promMetricSender, kafkaMetricsSender)

	// start all senders
	for _, metricSender := range listMetricSenders {
		metricSender.Start()
	}

	// start wait loop of dispatcher
	go func() {
		for {
			select {
			case work := <-MetricsQueue:
				log.Info("Received send metrics request")
				go func() {
					for _, metricsSender := range listMetricSenders {
						metricsSender.AssignMetricsToSend(work)
					}

				}()
			}
		}
	}()
}
