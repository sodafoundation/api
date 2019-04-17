package diverters

import (
	"fmt"
	"github.com/opensds/opensds/pkg/model"
)

// A buffered channel that we can send work requests on.
var MetricsQueue = make(chan model.MetricSpec, 100)

func SendMetricToRegisteredSenders(metrics model.MetricSpec) {

	// Push the work onto the queue.
	MetricsQueue <- metrics
	fmt.Println("Send metrics request queued")

	return
}

func StartDispatcher() {

	listMetricSenders := make([]MetricsSenderIntf,0)

	// initialize Prometheus sender
	senderStructProm := PrometheusMetricsSender{}
	promMetricSender := senderStructProm.GetMetricsSender()

	senderStructKafka := KafkaMetricsSender{}
	kafkaMetricsSender := senderStructKafka.GetMetricsSender()

	// add to list
	listMetricSenders = append(listMetricSenders, promMetricSender, kafkaMetricsSender)

	// start all senders
	for _,metricSender := range listMetricSenders{
		metricSender.Start()
	}

	// start wait loop of dispatcher
	go func() {
		for {
			select {
			case work := <-MetricsQueue:
				fmt.Println("Received send metrics request")
				go func() {
					for i,metricsSender := range listMetricSenders{
						fmt.Println("Dispatching send metrics request to sender %d",i)
						metricsSender.AssignMetricsToSend(work)
					}

				}()
			}
		}
	}()
}
