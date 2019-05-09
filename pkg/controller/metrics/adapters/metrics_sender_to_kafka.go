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
	"context"
	"strconv"

	log "github.com/golang/glog"
	"github.com/opensds/opensds/pkg/model"
	"github.com/segmentio/kafka-go"
)

type KafkaMetricsSender struct {
	Queue    chan *model.MetricSpec
	QuitChan chan bool
}

func (p *KafkaMetricsSender) GetMetricsSender() MetricsSenderIntf {
	sender := KafkaMetricsSender{}
	sender.Queue = make(chan *model.MetricSpec)
	sender.QuitChan = make(chan bool)
	return &sender
}

func (p *KafkaMetricsSender) Start() {
	go func() {
		for {
			select {
			case work := <-p.Queue:
				// Receive a work request.
				log.Infof("GetMetricsSenderToKafka received metrics for instance %s\n and metrics %f\n", work.InstanceID, work.MetricValues[0].Value)

				// do the actual sending work here
				// make a writer that produces to topic-A, using the least-bytes distribution
				w := kafka.NewWriter(kafka.WriterConfig{
					Brokers:  []string{"localhost:9092"},
					Topic:    "test",
					Balancer: &kafka.LeastBytes{},
				})

				// get the string ready to be written
				var finalString = ""
				//for _,metric := range work.MetricValues{
				finalString += work.Name + " " + strconv.FormatFloat(work.MetricValues[0].Value, 'f', 2, 64) + "\n"

				w.WriteMessages(context.Background(),
					kafka.Message{
						Key:   []byte("Key-A"),
						Value: []byte(finalString),
					})

				w.Close()

				log.Info("GetMetricsSenderToKafka processed metrics")

			case <-p.QuitChan:
				return
			}

		}
	}()
}
func (p *KafkaMetricsSender) Stop() {
	go func() {
		p.QuitChan <- true
	}()
}

func (p *KafkaMetricsSender) AssignMetricsToSend(request *model.MetricSpec) {
	p.Queue <- request
}
