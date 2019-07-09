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
	"testing"

	"github.com/opensds/opensds/pkg/model"
)

var SamplemetricsSpec = []*model.MetricSpec{
	{InstanceID: "volume-38aa800e-dc4b-4e01-9a0f-926f58ee2f14",
		InstanceName: "opensds--volumes--default-volume--38aa800e--dc4b--4e01--9a0f--926f58ee2f14",
		Job:          "lvm",
		Labels:       map[string]string{"device": "opensds--volumes--default-volume--38aa800e--dc4b--4e01--9a0f--926f58ee2f14"},
		Component:    "volume",
		Name:         "iops",
		Unit:         "tps",
		AggrType:     "",
		MetricValues: []*model.Metric{
			&model.Metric{
				Timestamp: 1561465759,
				Value:     32.14,
			},
		},
	},
	{InstanceID: "volume-38aa800e-dc4b-4e01-9a0f-926f58ee2f14",
		InstanceName: "opensds--volumes--default-volume--38aa800e--dc4b--4e01--9a0f--926f58ee2f14",
		Job:          "lvm",
		Labels:       map[string]string{"device": "opensds--volumes--default-volume--38aa800e--dc4b--4e01--9a0f--926f58ee2f14"},
		Component:    "volume",
		Name:         "iops",
		Unit:         "tps",
		AggrType:     "total",
		MetricValues: []*model.Metric{
			&model.Metric{
				Timestamp: 1561465759,
				Value:     32.14,
			},
		},
	},
}

func TestPrometheusMetricsSender_GetMetricsSender(t *testing.T) {
	type fields struct {
		Queue    chan *model.MetricSpec
		QuitChan chan bool
	}
	expectedMetricSender := &PrometheusMetricsSender{
		Queue:    make(chan *model.MetricSpec),
		QuitChan: make(chan bool),
	}
	tests := []struct {
		name   string
		fields fields
		want   MetricsSenderIntf
	}{
		{name: "test1", fields: fields{Queue: make(chan *model.MetricSpec), QuitChan: make(chan bool)}, want: expectedMetricSender},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PrometheusMetricsSender{
				Queue:    tt.fields.Queue,
				QuitChan: tt.fields.QuitChan,
			}
			got := p.GetMetricsSender()
			if got == nil {
				t.Errorf("unexpected response for GetMetricsSender() = %v, ", got)

			}
		})
	}
}

func TestPrometheusMetricsSender_Start(t *testing.T) {
	type fields struct {
		Queue    chan *model.MetricSpec
		QuitChan chan bool
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{name: "test1", fields: fields{
			Queue:    make(chan *model.MetricSpec),
			QuitChan: make(chan bool)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PrometheusMetricsSender{
				Queue:    tt.fields.Queue,
				QuitChan: tt.fields.QuitChan,
			}
			p.Start()
			p.AssignMetricsToSend(SamplemetricsSpec[0])
			p.Stop()

			//t.SkipNow()
		})
	}
}

func Test_writeToFile(t *testing.T) {
	type args struct {
		metrics *model.MetricSpec
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "test1", args: args{metrics: SamplemetricsSpec[0]}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writeToFile(tt.args.metrics)
		})
	}
}

func Test_getMetricName(t *testing.T) {
	type args struct {
		metrics *model.MetricSpec
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "test1", args: args{metrics: SamplemetricsSpec[0]}, want: "lvm_volume_iops_tps"},
		{name: "test2", args: args{metrics: SamplemetricsSpec[1]}, want: "lvm_volume_iops_tps_total"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getMetricName(tt.args.metrics); got != tt.want {
				t.Errorf("getMetricName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sendToPushGateway(t *testing.T) {
	type args struct {
		metrics *model.MetricSpec
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "test1", args: args{metrics: SamplemetricsSpec[0]}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sendToPushGateway(tt.args.metrics)
		})
	}
}
