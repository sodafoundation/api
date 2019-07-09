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
	"context"
	"encoding/json"

	log "github.com/golang/glog"
	"github.com/opensds/opensds/pkg/model"
	. "github.com/opensds/opensds/pkg/utils/config"
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
					Brokers:  []string{CONF.OsdsLet.KafkaEndpoint},
					Topic:    CONF.OsdsLet.KafkaTopic,
					Balancer: &kafka.LeastBytes{},
				})

				// send the whole struct as JSON onto the KAFKA topic
				byteArr, err := json.Marshal(*work)
				if err != nil {
					log.Errorf("error marshaling metrics to json, for kafka topic")
					return
				}
				finalString := string(byteArr)
				sendErr := w.WriteMessages(context.Background(),
					kafka.Message{
						Key:   []byte("metricspec"),
						Value: []byte(finalString),
					})
				if sendErr != nil {
					log.Errorf("error sending metrics to kafka topic")
				}
				_ = w.Close()

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
