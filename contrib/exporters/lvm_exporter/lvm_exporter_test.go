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
package main

import (
	"reflect"
	"sync"
	"testing"

	"github.com/opensds/opensds/contrib/drivers"
	"github.com/opensds/opensds/pkg/model"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	SampleVolumeIOPS *prometheus.Desc = prometheus.NewDesc("lvm_volume_iops_tps",
		"Shows IOPS",
		[]string{"device"}, nil,
	)
	SampleVolumeReadThroughput *prometheus.Desc = prometheus.NewDesc("lvm_volume_read_throughput_kbs",
		"Shows ReadThroughput",
		[]string{"device"}, nil,
	)
	SampleVolumeWriteThroughput *prometheus.Desc = prometheus.NewDesc("lvm_volume_write_throughput_kbs",
		"Shows ReadThroughput",
		[]string{"device"}, nil,
	)
	SampleVolumeResponseTime *prometheus.Desc = prometheus.NewDesc("lvm_volume_response_time_ms",
		"Shows ReadThroughput",
		[]string{"device"}, nil,
	)
	SampleVolumeServiceTime *prometheus.Desc = prometheus.NewDesc("lvm_volume_service_time_ms",
		"Shows ServiceTime",
		[]string{"device"}, nil,
	)
	SampleVolumeUtilization *prometheus.Desc = prometheus.NewDesc("lvm_volume_utilization_prcnt",
		"Shows Utilization in percentage",
		[]string{"device"}, nil,
	)
	SampleDiskIOPS *prometheus.Desc = prometheus.NewDesc("lvm_disk_iops_tps",
		"Shows IOPS",
		[]string{"device"}, nil,
	)
	SampleDiskReadThroughput *prometheus.Desc = prometheus.NewDesc("lvm_disk_read_throughput_kbs",
		"Shows Disk ReadThroughput",
		[]string{"device"}, nil,
	)
	SampleDiskWriteThroughput *prometheus.Desc = prometheus.NewDesc("lvm_disk_write_throughput_kbs",
		"Shows Write Throughput",
		[]string{"device"}, nil,
	)
	SampleDiskResponseTime *prometheus.Desc = prometheus.NewDesc("lvm_disk_response_time_ms",
		"Shows Disk Response Time",
		[]string{"device"}, nil,
	)
	SampleDiskServiceTime *prometheus.Desc = prometheus.NewDesc("lvm_disk_service_time_ms",
		"Shows ServiceTime",
		[]string{"device"}, nil,
	)
	SampleDiskUtilization *prometheus.Desc = prometheus.NewDesc("lvm_disk_utilization_prcnt",
		"Shows Utilization in percentage",
		[]string{"device"}, nil,
	)
	SamplemetricsSpec1 = []*model.MetricSpec{
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
			Name:         "read_throughput",
			Unit:         "kbs",
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
			Name:         "write_throughput",
			Unit:         "kbs",
			AggrType:     "",
			MetricValues: []*model.Metric{
				&model.Metric{
					Timestamp: 1561465759,
					Value:     32.14,
				},
			},
		},
	}
	SamplemetricsSpec2 = []*model.MetricSpec{
		{InstanceID: "volume-38aa800e-dc4b-4e01-9a0f-926f58ee2f14",
			InstanceName: "opensds--volumes--default-volume--38aa800e--dc4b--4e01--9a0f--926f58ee2f14",
			Job:          "lvm",
			Labels:       map[string]string{"device": "opensds--volumes--default-volume--38aa800e--dc4b--4e01--9a0f--926f58ee2f14"},
			Component:    "disk",
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
			Component:    "disk",
			Name:         "read_throughput",
			Unit:         "kbs",
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
			Component:    "disk",
			Name:         "write_throughput",
			Unit:         "kbs",
			AggrType:     "",
			MetricValues: []*model.Metric{
				&model.Metric{
					Timestamp: 1561465759,
					Value:     32.14,
				},
			},
		},
	}
	execCount = 0
)

func Test_newLvmCollector(t *testing.T) {
	tests := []struct {
		name string
		want *lvmCollector
	}{
		{
			name: "test1",
			want: &lvmCollector{
				mu:                    sync.Mutex{},
				VolumeIOPS:            SampleVolumeIOPS,
				VolumeReadThroughput:  SampleVolumeReadThroughput,
				VolumeWriteThroughput: SampleVolumeWriteThroughput,
				VolumeResponseTime:    SampleVolumeResponseTime,
				VolumeServiceTime:     SampleVolumeServiceTime,
				VolumeUtilization:     SampleVolumeUtilization,
				DiskIOPS:              SampleDiskIOPS,
				DiskReadThroughput:    SampleDiskReadThroughput,
				DiskWriteThroughput:   SampleDiskWriteThroughput,
				DiskResponseTime:      SampleDiskResponseTime,
				DiskServiceTime:       SampleDiskServiceTime,
				DiskUtilization:       SampleDiskUtilization,
				metricDriver:          drivers.InitMetricDriver("lvm"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newLvmCollector(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("unexpected result newLvmCollector() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_lvmCollector_Describe(t *testing.T) {
	type fields struct {
		mu                    sync.Mutex
		VolumeIOPS            *prometheus.Desc
		VolumeReadThroughput  *prometheus.Desc
		VolumeWriteThroughput *prometheus.Desc
		VolumeResponseTime    *prometheus.Desc
		VolumeServiceTime     *prometheus.Desc
		VolumeUtilization     *prometheus.Desc
		DiskIOPS              *prometheus.Desc
		DiskReadThroughput    *prometheus.Desc
		DiskWriteThroughput   *prometheus.Desc
		DiskResponseTime      *prometheus.Desc
		DiskServiceTime       *prometheus.Desc
		DiskUtilization       *prometheus.Desc
	}
	type args struct {
		ch chan *prometheus.Desc
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "test1",
			fields: fields{
				mu:                    sync.Mutex{},
				VolumeIOPS:            SampleVolumeIOPS,
				VolumeReadThroughput:  SampleVolumeReadThroughput,
				VolumeWriteThroughput: SampleVolumeWriteThroughput,
				VolumeResponseTime:    SampleVolumeResponseTime,
				VolumeServiceTime:     SampleVolumeServiceTime,
				VolumeUtilization:     SampleVolumeUtilization,
				DiskIOPS:              SampleDiskIOPS,
				DiskReadThroughput:    SampleDiskReadThroughput,
				DiskWriteThroughput:   SampleDiskWriteThroughput,
				DiskResponseTime:      SampleDiskResponseTime,
				DiskServiceTime:       SampleDiskServiceTime,
				DiskUtilization:       SampleDiskUtilization,
			},
			args: args{
				ch: make(chan *prometheus.Desc),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &lvmCollector{
				mu:                    tt.fields.mu,
				VolumeIOPS:            tt.fields.VolumeIOPS,
				VolumeReadThroughput:  tt.fields.VolumeReadThroughput,
				VolumeWriteThroughput: tt.fields.VolumeWriteThroughput,
				VolumeResponseTime:    tt.fields.VolumeResponseTime,
				VolumeServiceTime:     tt.fields.VolumeServiceTime,
				VolumeUtilization:     tt.fields.VolumeUtilization,
				DiskIOPS:              tt.fields.DiskIOPS,
				DiskReadThroughput:    tt.fields.DiskReadThroughput,
				DiskWriteThroughput:   tt.fields.DiskWriteThroughput,
				DiskResponseTime:      tt.fields.DiskResponseTime,
				DiskServiceTime:       tt.fields.DiskServiceTime,
				DiskUtilization:       tt.fields.DiskUtilization,
			}
			//fmt.Println(c)
			go Consume(tt.args.ch)
			c.Describe(tt.args.ch)

		})
	}
}
func Consume(ch <-chan *prometheus.Desc) {

	for {
		<-ch

	}
}
func Test_lvmCollector_Collect(t *testing.T) {
	type fields struct {
		mu                    sync.Mutex
		VolumeIOPS            *prometheus.Desc
		VolumeReadThroughput  *prometheus.Desc
		VolumeWriteThroughput *prometheus.Desc
		VolumeResponseTime    *prometheus.Desc
		VolumeServiceTime     *prometheus.Desc
		VolumeUtilization     *prometheus.Desc
		DiskIOPS              *prometheus.Desc
		DiskReadThroughput    *prometheus.Desc
		DiskWriteThroughput   *prometheus.Desc
		DiskResponseTime      *prometheus.Desc
		DiskServiceTime       *prometheus.Desc
		DiskUtilization       *prometheus.Desc
	}
	type args struct {
		ch chan prometheus.Metric
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "test1",
			fields: fields{
				mu:                    sync.Mutex{},
				VolumeIOPS:            SampleVolumeIOPS,
				VolumeReadThroughput:  SampleVolumeReadThroughput,
				VolumeWriteThroughput: SampleVolumeWriteThroughput,
				VolumeResponseTime:    SampleVolumeResponseTime,
				VolumeServiceTime:     SampleVolumeServiceTime,
				VolumeUtilization:     SampleVolumeUtilization,
				DiskIOPS:              SampleDiskIOPS,
				DiskReadThroughput:    SampleDiskReadThroughput,
				DiskWriteThroughput:   SampleDiskWriteThroughput,
				DiskResponseTime:      SampleDiskResponseTime,
				DiskServiceTime:       SampleDiskServiceTime,
				DiskUtilization:       SampleDiskUtilization,
			},
			args: args{
				ch: make(chan prometheus.Metric),
			},
		},
		{
			name: "test2",
			fields: fields{
				mu:                    sync.Mutex{},
				VolumeIOPS:            SampleVolumeIOPS,
				VolumeReadThroughput:  SampleVolumeReadThroughput,
				VolumeWriteThroughput: SampleVolumeWriteThroughput,
				VolumeResponseTime:    SampleVolumeResponseTime,
				VolumeServiceTime:     SampleVolumeServiceTime,
				VolumeUtilization:     SampleVolumeUtilization,
				DiskIOPS:              SampleDiskIOPS,
				DiskReadThroughput:    SampleDiskReadThroughput,
				DiskWriteThroughput:   SampleDiskWriteThroughput,
				DiskResponseTime:      SampleDiskResponseTime,
				DiskServiceTime:       SampleDiskServiceTime,
				DiskUtilization:       SampleDiskUtilization,
			},
			args: args{
				ch: make(chan prometheus.Metric),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &lvmCollector{
				mu:                    tt.fields.mu,
				VolumeIOPS:            tt.fields.VolumeIOPS,
				VolumeReadThroughput:  tt.fields.VolumeReadThroughput,
				VolumeWriteThroughput: tt.fields.VolumeWriteThroughput,
				VolumeResponseTime:    tt.fields.VolumeResponseTime,
				VolumeServiceTime:     tt.fields.VolumeServiceTime,
				VolumeUtilization:     tt.fields.VolumeUtilization,
				DiskIOPS:              tt.fields.DiskIOPS,
				DiskReadThroughput:    tt.fields.DiskReadThroughput,
				DiskWriteThroughput:   tt.fields.DiskWriteThroughput,
				DiskResponseTime:      tt.fields.DiskResponseTime,
				DiskServiceTime:       tt.fields.DiskServiceTime,
				DiskUtilization:       tt.fields.DiskUtilization,
				metricDriver:          NewFakeMetricDriver(),
			}
			go ConsumeMetrics(tt.args.ch)
			c.Collect(tt.args.ch)
		})
	}
}
func ConsumeMetrics(ch <-chan prometheus.Metric) {

	for {
		<-ch

	}
}

type fakeMetricDriver struct {
}

func (d *fakeMetricDriver) CollectMetrics() ([]*model.MetricSpec, error) {
	execCount++
	if execCount%2 == 0 {
		return SamplemetricsSpec1, nil
	} else {
		return SamplemetricsSpec2, nil
	}
}

func (d *fakeMetricDriver) GetMetricList(resourceType string) (supportedMetrics []string, err error) {
	return
}
func (d *fakeMetricDriver) Setup() error {

	return nil
}

func (*fakeMetricDriver) Teardown() error { return nil }

func NewFakeMetricDriver() drivers.MetricDriver {

	return &fakeMetricDriver{}
}

func Test_validateCliArg(t *testing.T) {
	type args struct {
		arg1 string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test1",
			args: args{
				arg1: "9601",
			},
			want: "9601",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validateCliArg(tt.args.arg1); got != tt.want {
				t.Errorf("validateCliArg() = %v, want %v", got, tt.want)
			}
		})
	}
}
