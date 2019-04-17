package diverters

import (
	"fmt"
	"github.com/opensds/opensds/pkg/model"
	"github.com/prometheus/client_golang/prometheus"
	"os"
	"strconv"
	"time"
	"github.com/prometheus/client_golang/prometheus/push"
)

var nodeExporterFolder = "/root/prom_nodeexporter_folder/"

type PrometheusMetricsSender struct {
	Queue chan model.MetricSpec
	QuitChan chan bool
}

func (p *PrometheusMetricsSender) GetMetricsSender() MetricsSenderIntf{
	sender := PrometheusMetricsSender{}
	sender.Queue = make(chan model.MetricSpec)
	sender.QuitChan = make(chan bool)
	return &sender
}

func (p *PrometheusMetricsSender) Start() {
	go func() {
		for {
			select {
			case work := <-p.Queue:
				// Receive a work request.
				fmt.Printf("GetMetricsSenderToPrometheus received metrics for instance %s\n and metrics %s\n", work.InstanceID, work.Value)

				// do the actual sending work here, by writing to the file of the node_exporter of prometheus
				writeToFile(work)

				// alternatively, we could also push the metrics to the push gateway of prometheus
				sendToPushGateway(work)

				fmt.Println("GetMetricsSenderToPrometheus processed metrics")

			case <-p.QuitChan:
				// We have been asked to stop.
				fmt.Println("GetMetricsSenderToPrometheus stopping\n")
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

func (p *PrometheusMetricsSender) AssignMetricsToSend(request model.MetricSpec){
	p.Queue <- request
}

func writeToFile(metrics model.MetricSpec){

	// get the string ready to be written
	var finalString = ""
	//for _,metric := range metrics.MetricValues{
		finalString += metrics.Name + " " + strconv.FormatFloat(metrics.Value,'f', 2,64) + "\n"
	//}
	// make a new file with current timestamp
	timeStamp := strconv.FormatInt(time.Now().Unix(), 10)
	f, err := os.Create(nodeExporterFolder + metrics.InstanceID + ".prom")
	if err != nil {
		fmt.Println(err)
		return
	}
		l, err := f.WriteString(finalString)
		if err != nil {
			fmt.Println(err)
			f.Close()
			return
		}
		fmt.Println(l, "metrics written successfully at time " + timeStamp)
		err = f.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
}

func sendToPushGateway(metrics model.MetricSpec){

	completionTime := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: metrics.Name,
		Help: "",
	})
	completionTime.SetToCurrentTime()
	completionTime.Set(metrics.Value)

	if err := push.New("http://localhost:9091", "push_gateway").
		Collector(completionTime).
		Grouping("l1", "v1").
		Push(); err != nil {
		fmt.Println("Could not push completion time to Pushgateway:", err)
	}
	fmt.Println("Completed push completion time to Pushgateway:")
}