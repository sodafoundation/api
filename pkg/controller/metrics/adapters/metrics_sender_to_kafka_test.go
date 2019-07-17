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

func TestKafkaMetricsSender_GetMetricsSender(t *testing.T) {
	type fields struct {
		Queue    chan *model.MetricSpec
		QuitChan chan bool
	}
	expectedMetricSender := &KafkaMetricsSender{
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
			p := &KafkaMetricsSender{
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

func TestKafkaMetricsSender_Start(t *testing.T) {
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
			p := &KafkaMetricsSender{
				Queue:    tt.fields.Queue,
				QuitChan: tt.fields.QuitChan,
			}
			p.Start()
			p.AssignMetricsToSend(SamplemetricsSpec[0])
			p.Stop()

		})
	}
}

func TestKafkaMetricsSender_Stop(t *testing.T) {
	type fields struct {
		Queue    chan *model.MetricSpec
		QuitChan chan bool
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{name: "test1", fields: fields{Queue: make(chan *model.MetricSpec),
			QuitChan: make(chan bool)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &KafkaMetricsSender{
				Queue:    tt.fields.Queue,
				QuitChan: tt.fields.QuitChan,
			}
			p.Stop()
		})
	}
}
