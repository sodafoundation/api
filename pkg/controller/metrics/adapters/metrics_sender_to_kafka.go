package adapters

import (
	"context"
	log "github.com/golang/glog"
	"github.com/opensds/opensds/pkg/model"
	"github.com/segmentio/kafka-go"
	"strconv"
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
				log.Infof("GetMetricsSenderToKafka received metrics for instance %s\n and metrics %s\n", work.InstanceID, work.Value)

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
				finalString += work.Name + " " + strconv.FormatFloat(work.Value, 'f', 2, 64) + "\n"

				w.WriteMessages(context.Background(),
					kafka.Message{
						Key:   []byte("Key-A"),
						Value: []byte(finalString),
					})

				w.Close()

				log.Info("GetMetricsSenderToKafka processed metrics")

			case <-p.QuitChan:
				// We have been asked to stop.
				//fmt.Println("GetMetricsSenderToKafka stopping\n")
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
