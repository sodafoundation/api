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

func TestSendMetricToRegisteredSenders(t *testing.T) {
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
			SendMetricToRegisteredSenders(tt.args.metrics)
		})
	}
}

func TestStartDispatcher(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "test1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			StartDispatcher()
		})
	}
}
