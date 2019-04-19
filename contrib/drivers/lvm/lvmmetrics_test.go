
package lvm

import (
	"fmt"
	"testing"
)



func TestMetricDriverSetup(t *testing.T) {
	var d = &MetricDriver{}


	if err := d.Setup(); err != nil {
		t.Errorf("Setup lvm metric  driver failed: %+v\n", err)
	}

}



func TestCollectMetrics(t *testing.T) {
	metricList:=[]string{"IOPS","ReadThroughput","WriteThroughput","ResponseTime"}
	var metricDriver = &MetricDriver{}
	metricDriver.Setup()
	metricArray,err:=metricDriver.CollectMetrics(metricList,"sda")
	if err != nil {
		t.Errorf("CollectMetrics call to lvm driver failed: %+v\n", err)
	}
	fmt.Println(metricArray)

}



