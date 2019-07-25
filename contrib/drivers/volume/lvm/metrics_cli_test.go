// Copyright 2019 The OpenSDS Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package lvm

import (
	"reflect"
	"testing"

	"github.com/opensds/opensds/pkg/utils/exec"
)

func TestNewMetricCli(t *testing.T) {
	tests := []struct {
		name    string
		want    *MetricCli
		wantErr bool
	}{
		{name: "test1", want: &MetricCli{
			BaseExecuter: exec.NewBaseExecuter(),
			RootExecuter: exec.NewRootExecuter(),
		}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewMetricCli()
			if (err != nil) != tt.wantErr {
				t.Errorf("unexpected error output for NewMetricCli() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("unexpected cli output for  NewMetricCli() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetricCli_execute(t *testing.T) {
	type fields struct {
		BaseExecuter exec.Executer
		RootExecuter exec.Executer
	}
	type args struct {
		cmd []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{name: "test1", fields: fields{}, args: args{}, want: "env: ‘error_cmd’: No such file or directory\n", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := NewMetricCli()
			got, err := c.execute("env", "LC_ALL=C", "error_cmd")
			if (err != nil) != tt.wantErr {
				t.Errorf("difference in expected result of MetricCli.execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("difference in expected result of of MetricCli.execute() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isSarEnabled(t *testing.T) {
	type args struct {
		out string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "test1", args: args{out: "Please check if data collecting is enabled"}, want: false},
		{name: "test1", args: args{out: "Command 'sar' not found"}, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isSarEnabled(tt.args.out); got != tt.want {
				t.Errorf("difference in expected result of expected isSarEnabled() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetricCli_parseCommandOutput(t *testing.T) {

	type fields struct {
		BaseExecuter exec.Executer
		RootExecuter exec.Executer
	}
	type args struct {
		metricList []string

		out *MetricFakeResp
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{name: "test1", fields: fields{}, args: args{metricList: expctdMetricList, out: respMap["sar"]}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			returnMap := make(map[string]map[string]string)
			labelMap := make(map[string]map[string]string)
			metricMap := make(map[string]int)
			c := &MetricCli{
				BaseExecuter: tt.fields.BaseExecuter,
				RootExecuter: tt.fields.RootExecuter,
			}
			c.parseCommandOutput(tt.args.metricList, returnMap, labelMap, metricMap, tt.args.out.out)
		})
	}
}

func TestMetricCli_CollectMetrics(t *testing.T) {
	type fields struct {
		BaseExecuter exec.Executer
		RootExecuter exec.Executer
	}
	type args struct {
		metricList []string
	}
	var expctedLabelMap map[string]map[string]string = map[string]map[string]string{"loop2": map[string]string{"device": "loop2"}, "loop3": map[string]string{"device": "loop3"}, "loop4": map[string]string{"device": "loop4"}, "loop5": map[string]string{"device": "loop5"}, "loop6": map[string]string{"device": "loop6"}, "loop7": map[string]string{"device": "loop7"}, "loop9": map[string]string{"device": "loop9"}, "loop10": map[string]string{"device": "loop10"}, "opensds--volumes--default-volume-d96cc42b-b285-474e-aa98-c61e66df7461": map[string]string{"device": "opensds--volumes--default-volume-d96cc42b-b285-474e-aa98-c61e66df7461"}, "opensds--volumes--default-volume--b902e771--8e02--4099--b601--a6b3881f8": map[string]string{"device": "opensds--volumes--default-volume--b902e771--8e02--4099--b601--a6b3881f8"}}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[string]map[string]string
		wantErr bool
	}{
		{name: "test1", fields: fields{}, args: args{metricList: expctdMetricList}, want: expctedLabelMap, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli := &MetricCli{
				BaseExecuter: NewMetricFakeExecuter(respMap),
				RootExecuter: NewMetricFakeExecuter(respMap),
			}
			_, got1, err := cli.CollectMetrics(tt.args.metricList)
			if (err != nil) != tt.wantErr {
				t.Errorf("difference in expected result of MetricCli.CollectMetrics() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got1, tt.want) {
				t.Errorf("difference in expected result of MetricCli.CollectMetrics() got = %v, want %v", got1, tt.want)
			}

		})
	}
}

func TestMetricCli_DiscoverVolumes(t *testing.T) {
	type fields struct {
		BaseExecuter exec.Executer
		RootExecuter exec.Executer
	}
	tests := []struct {
		name    string
		fields  fields
		want    []string
		want1   []string
		wantErr bool
	}{
		{name: "test1", fields: fields{}, want: expctedVolList, want1: expectdVgs},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &MetricCli{
				BaseExecuter: NewMetricFakeExecuter(respMap),
				RootExecuter: NewMetricFakeExecuter(respMap),
			}
			got, got1, err := c.DiscoverVolumes()
			if (err != nil) != tt.wantErr {
				t.Errorf("difference in expected result of MetricCli.DiscoverVolumes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("difference in expected result of MetricCli.DiscoverVolumes() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("difference in expected result of MetricCli.DiscoverVolumes() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestMetricCli_DiscoverDisks(t *testing.T) {
	type fields struct {
		BaseExecuter exec.Executer
		RootExecuter exec.Executer
	}
	tests := []struct {
		name    string
		fields  fields
		want    []string
		wantErr bool
	}{
		{name: "test1", fields: fields{}, want: expctedDiskList, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &MetricCli{
				BaseExecuter: NewMetricFakeExecuter(respMap),
				RootExecuter: NewMetricFakeExecuter(respMap),
			}
			got, err := c.DiscoverDisks()
			if (err != nil) != tt.wantErr {
				t.Errorf("difference in expected result of MetricCli.DiscoverDisks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("difference in expected result of MetricCli.DiscoverDisks() = %v, want %v", got, tt.want)
			}
		})
	}
}
