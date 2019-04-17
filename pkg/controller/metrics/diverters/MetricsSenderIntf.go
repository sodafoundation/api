package diverters

import "github.com/opensds/opensds/pkg/model"

type MetricsSenderIntf interface {

	//GetMetricsSender(queue chan MetricRequest, quitChan chan bool) MetricsSenderIntf
	GetMetricsSender() MetricsSenderIntf
	AssignMetricsToSend(request model.MetricSpec)
	Start()
	Stop()
}

