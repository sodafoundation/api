// Copyright (c) 2017 Huawei Technologies Co., Ltd. All Rights Reserved.
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

package config

import (
	"reflect"
	"testing"
)

type TestStruct struct {
	Bool    bool    `conf:"bool,true"`
	Int     int     `conf:"int,-456789"`
	Int8    int8    `conf:"int8,-111"`
	Int16   int     `conf:"int16,-4567"`
	Int32   int     `conf:"int32,-456789"`
	Int64   int64   `conf:"int64,-456789"`
	Uint    uint    `conf:"uint,456789"`
	Uint8   uint8   `conf:"uint8,111"`
	Uint16  uint16  `conf:"uint16,4567"`
	Uint32  uint32  `conf:"uint32,456789"`
	Uint64  uint64  `conf:"uint64,456789"`
	Float32 float32 `conf:"float32,0.456789"`
	Float64 float64 `conf:"float64,0.456789"`
	String  string  `conf:"string,DefaultValue"`
}

type TestSliceStruct struct {
	SliceBool           []bool    `conf:"slice_bool,False,True,False"`
	SliceString         []string  `conf:"slice_string,slice,string,test"`
	SliceInt            []int     `conf:"slice_int,4,-5,6"`
	SliceInt8           []int8    `conf:"slice_int8,4,-5,6"`
	SliceInt16          []int16   `conf:"slice_int16,4,-5,6"`
	SliceInt32          []int32   `conf:"slice_int32,4,-5,6"`
	SliceInt64          []int64   `conf:"slice_int64,4,-5,6"`
	SliceUint           []uint    `conf:"slice_uint,4,5,6"`
	SliceUint8          []uint8   `conf:"slice_uint8,4,5,6"`
	SliceUint16         []uint16  `conf:"slice_uint16,4,5,6"`
	SliceUint32         []uint32  `conf:"slice_uint32,4,5,6"`
	SliceUint64         []uint64  `conf:"slice_uint64,4,5,6"`
	SliceFloat32        []float32 `conf:"slice_float32,4,-0.5,0.6"`
	SliceFloat64        []float64 `conf:"slice_float64,4,-0.5,0.6"`
	SliceNotExistString []string  `conf:"slice_not_exist_string,not,exist,string"`
}
type TestEmbedSectionStruct struct {
	TestStruct      `conf:"test_struct"`
	TestSliceStruct `conf:"test_slice_struct"`
}
type TestConfig struct {
	TestStruct      `conf:"test_struct"`
	TestSliceStruct `conf:"test_slice_struct"`
	TestEmbedSectionStruct
}

func TestFunctionAllType(t *testing.T) {
	conf := &TestConfig{}
	initConf("testdata/opensds.conf", conf)
	if conf.TestStruct.Bool != true {
		t.Error("Test TestStuct Bool error")
	}
	if conf.TestStruct.Int != -123456 {
		t.Error("Test TestStuct Int error")
	}
	if conf.TestStruct.Int8 != -123 {
		t.Error("Test TestStuct Int8 error")
	}
	if conf.TestStruct.Int16 != -1234 {
		t.Error("Test TestStuct Int16 error")
	}
	if conf.TestStruct.Int32 != -123456 {
		t.Error("Test TestStuct Int32 error")
	}
	if conf.TestStruct.Int64 != -123456 {
		t.Error("Test TestStuct Int64 error")
	}
	if conf.TestStruct.Uint != 123456 {
		t.Error("Test TestStuct Uint error")
	}
	if conf.TestStruct.Uint8 != 123 {
		t.Error("Test TestStuct Uint8 error")
	}
	if conf.TestStruct.Uint16 != 12345 {
		t.Error("Test TestStuct Uint16 error")
	}
	if conf.TestStruct.Uint32 != 123456 {
		t.Error("Test TestStuct Uint32 error")
	}
	if conf.TestStruct.Uint64 != 123456 {
		t.Error("Test TestStuct Uint64 error")
	}
	if conf.TestStruct.Float32 != 0.123456 {
		t.Error("Test TestStuct Float32 error")
	}
	if conf.TestStruct.Float64 != 0.123456 {
		t.Error("Test TestStuct Float64 error")
	}
	if conf.TestStruct.String != "HelloWorld" {
		t.Error("Test TestStuct String error")
	}
	if !reflect.DeepEqual(conf.TestSliceStruct.SliceString, []string{"slice", "string", "test"}) {
		t.Error("Test TestSliceStruct String error")
	}
	if !reflect.DeepEqual(conf.TestSliceStruct.SliceBool, []bool{false, true, false}) {
		t.Error("Test TestSliceStruct bool error")
	}
	if !reflect.DeepEqual(conf.TestSliceStruct.SliceInt, []int{1, -2, 3}) {
		t.Error("Test TestSliceStruct int error")
	}
	if !reflect.DeepEqual(conf.TestSliceStruct.SliceInt8, []int8{1, -2, 3}) {
		t.Error("Test TestSliceStruct int8 error")
	}
	if !reflect.DeepEqual(conf.TestSliceStruct.SliceInt16, []int16{1, -2, 3}) {
		t.Error("Test TestSliceStruct int16 error")
	}
	if !reflect.DeepEqual(conf.TestSliceStruct.SliceInt32, []int32{1, -2, 3}) {
		t.Error("Test TestSliceStruct int32 error")
	}
	if !reflect.DeepEqual(conf.TestSliceStruct.SliceInt64, []int64{1, -2, 3}) {
		t.Error("Test TestSliceStruct int64 error")
	}
	if !reflect.DeepEqual(conf.TestSliceStruct.SliceUint, []uint{1, 2, 3}) {
		t.Error("Test TestSliceStruct uint error")
	}
	if !reflect.DeepEqual(conf.TestSliceStruct.SliceUint8, []uint8{1, 2, 3}) {
		t.Error("Test TestSliceStruct uint8 error")
	}
	if !reflect.DeepEqual(conf.TestSliceStruct.SliceUint16, []uint16{1, 2, 3}) {
		t.Error("Test TestSliceStruct uint16 error")
	}
	if !reflect.DeepEqual(conf.TestSliceStruct.SliceUint32, []uint32{1, 2, 3}) {
		t.Error("Test TestSliceStruct uint32 error")
	}
	if !reflect.DeepEqual(conf.TestSliceStruct.SliceUint64, []uint64{1, 2, 3}) {
		t.Error("Test TestSliceStruct uint64 error")
	}
	if !reflect.DeepEqual(conf.TestSliceStruct.SliceFloat32, []float32{1, -0.2, 0.3}) {
		t.Error("Test TestSliceStruct float32 error")
	}
	if !reflect.DeepEqual(conf.TestSliceStruct.SliceFloat64, []float64{1, -0.2, 0.3}) {
		t.Error("Test TestSliceStruct float64 error")
	}
	if !reflect.DeepEqual(conf.TestSliceStruct.SliceNotExistString, []string{"not", "exist", "string"}) {
		t.Error("Test TestSliceStruct String error")
	}
	if conf.TestEmbedSectionStruct.TestStruct.String != "HelloWorld" {
		t.Error("Test TestEmbedSectionStruct.TestStuct String error")
	}
	if !reflect.DeepEqual(conf.TestEmbedSectionStruct.TestSliceStruct.SliceString, []string{"slice", "string", "test"}) {
		t.Error("Test TestEmbedSectionStruct.TestSliceStruct String error")
	}
}

func TestFunctionDefaultValue(t *testing.T) {
	conf := &TestConfig{}
	initConf("NotExistFile", conf)

	if conf.TestStruct.Bool != true {
		t.Error("Test TestStuct Bool error")
	}
	if conf.TestStruct.Int != -456789 {
		t.Error("Test TestStuct Int error")
	}
	if conf.TestStruct.Int8 != -111 {
		t.Error("Test TestStuct Int8 error")
	}
	if conf.TestStruct.Int16 != -4567 {
		t.Error("Test TestStuct Int16 error")
	}
	if conf.TestStruct.Int32 != -456789 {
		t.Error("Test TestStuct Int32 error")
	}
	if conf.TestStruct.Int64 != -456789 {
		t.Error("Test TestStuct Int64 error")
	}
	if conf.TestStruct.Uint != 456789 {
		t.Error("Test TestStuct Uint error")
	}
	if conf.TestStruct.Uint8 != 111 {
		t.Error("Test TestStuct Uint8 error")
	}
	if conf.TestStruct.Uint16 != 4567 {
		t.Error("Test TestStuct Uint16 error")
	}
	if conf.TestStruct.Uint32 != 456789 {
		t.Error("Test TestStuct Uint32 error")
	}
	if conf.TestStruct.Uint64 != 456789 {
		t.Error("Test TestStuct Uint64 error")
	}
	if conf.TestStruct.Float32 != 0.456789 {
		t.Error("Test TestStuct Float32 error")
	}
	if conf.TestStruct.Float64 != 0.456789 {
		t.Error("Test TestStuct Float64 error")
	}
	if conf.TestStruct.String != "DefaultValue" {
		t.Error("Test TestStuct String error")
	}

	if !reflect.DeepEqual(conf.TestSliceStruct.SliceString, []string{"slice", "string", "test"}) {
		t.Error("Test TestSliceStruct String error")
	}
	if !reflect.DeepEqual(conf.TestSliceStruct.SliceBool, []bool{false, true, false}) {
		t.Error("Test TestSliceStruct bool error")
	}
	if !reflect.DeepEqual(conf.TestSliceStruct.SliceInt, []int{4, -5, 6}) {
		t.Error("Test TestSliceStruct bool error")
	}
	if !reflect.DeepEqual(conf.TestSliceStruct.SliceInt8, []int8{4, -5, 6}) {
		t.Error("Test TestSliceStruct bool error")
	}
	if !reflect.DeepEqual(conf.TestSliceStruct.SliceInt16, []int16{4, -5, 6}) {
		t.Error("Test TestSliceStruct bool error")
	}
	if !reflect.DeepEqual(conf.TestSliceStruct.SliceInt32, []int32{4, -5, 6}) {
		t.Error("Test TestSliceStruct bool error")
	}
	if !reflect.DeepEqual(conf.TestSliceStruct.SliceInt64, []int64{4, -5, 6}) {
		t.Error("Test TestSliceStruct bool error")
	}
	if !reflect.DeepEqual(conf.TestSliceStruct.SliceUint, []uint{4, 5, 6}) {
		t.Error("Test TestSliceStruct bool error")
	}
	if !reflect.DeepEqual(conf.TestSliceStruct.SliceUint8, []uint8{4, 5, 6}) {
		t.Error("Test TestSliceStruct bool error")
	}
	if !reflect.DeepEqual(conf.TestSliceStruct.SliceUint16, []uint16{4, 5, 6}) {
		t.Error("Test TestSliceStruct bool error")
	}
	if !reflect.DeepEqual(conf.TestSliceStruct.SliceUint32, []uint32{4, 5, 6}) {
		t.Error("Test TestSliceStruct bool error")
	}
	if !reflect.DeepEqual(conf.TestSliceStruct.SliceUint64, []uint64{4, 5, 6}) {
		t.Error("Test TestSliceStruct bool error")
	}
	if !reflect.DeepEqual(conf.TestSliceStruct.SliceFloat32, []float32{4, -0.5, 0.6}) {
		t.Error("Test TestSliceStruct bool error")
	}
	if !reflect.DeepEqual(conf.TestSliceStruct.SliceFloat64, []float64{4, -0.5, 0.6}) {
		t.Error("Test TestSliceStruct bool error")
	}
}

func TestOpensdsConfig(t *testing.T) {
	CONF.Load("testdata/opensds.conf")
	if CONF.OsdsLet.ApiEndpoint != "localhost:50040" {
		t.Error("Test OsdsLet.ApiEndpoint error")
	}
	if CONF.OsdsLet.Graceful != true {
		t.Error("Test OsdsLet.Graceful error")
	}
	if CONF.OsdsLet.SocketOrder != "inc" {
		t.Error("Test OsdsLet.SocketOrder error")
	}
	if CONF.OsdsDock.ApiEndpoint != "localhost:50050" {
		t.Error("Test OsdsDock.ApiEndpoint error")
	}
	if CONF.OsdsDock.EnabledBackends[0] != "ceph" {
		t.Error("OsdsDock.EnabledBackends[0] error")
	}
	if CONF.OsdsDock.EnabledBackends[1] != "cinder" {
		t.Error("Test OsdsDock.EnabledBackends[1] error")
	}
	if CONF.OsdsDock.EnabledBackends[2] != "sample" {
		t.Error("Test OsdsDock.EnabledBackends[2] error")
	}
	if CONF.Database.Credential != "opensds:password@127.0.0.1:3306/dbname" {
		t.Error("Test Database.Credential error")
	}
	if CONF.Database.Endpoint != "localhost:2379,localhost:2380" {
		t.Error("Test Database.Endpoint error")
	}
	if CONF.Database.Driver != "etcd" {
		t.Error("Test Database.Driver error")
	}
	if CONF.Backends.Ceph.Name != "ceph" {
		t.Error("Test Ceph.Backends.Name error")
	}
	if CONF.Ceph.Name != "ceph" {
		t.Error("Test Ceph.Name error")
	}
	if CONF.Ceph.Description != "Ceph Test" {
		t.Error("Test Ceph.Description error")
	}
	if CONF.Ceph.DriverName != "ceph" {
		t.Error("Test Ceph.DriverName error")
	}
	if CONF.Ceph.ConfigPath != "/etc/opensds/driver/ceph.yaml" {
		t.Error("Test Ceph.ConfigPath error")
	}
	if CONF.Cinder.Name != "cinder" {
		t.Error("Test Cinder.Name error")
	}
	if CONF.Cinder.Description != "Cinder Test" {
		t.Error("Test Cinder.Description error")
	}
	if CONF.Cinder.DriverName != "cinder" {
		t.Error("Test Cinder.DriverName error")
	}
	if CONF.Cinder.ConfigPath != "/etc/opensds/driver/cinder.yaml" {
		t.Error("Test Cinder.ConfigPath error")
	}
	if CONF.Sample.Name != "sample" {
		t.Error("Test Sample.Name error")
	}
	if CONF.Sample.Description != "Sample Test" {
		t.Error("Test Sample.Description error")
	}
	if CONF.Sample.DriverName != "sample" {
		t.Error("Test Sample.DriverName error")
	}
	if CONF.Sample.ConfigPath != "/etc/opensds/driver/sample.yaml" {
		t.Error("Test Sample.ConfigPath error")
	}
	if CONF.LVM.Name != "lvm" {
		t.Error("Test LVM.Name error")
	}
	if CONF.LVM.Description != "LVM Test" {
		t.Error("Test Sample.Description error")
	}
	if CONF.LVM.DriverName != "lvm" {
		t.Error("Test LVM.DriverName error")
	}
	if CONF.LVM.ConfigPath != "/etc/opensds/driver/lvm.yaml" {
		t.Error("Test LVM.ConfigPath error")
	}
	bm := GetBackendsMap()
	if bm["ceph"].Name != "ceph" {
		t.Error("Test bm[\"ceph\"].Name error")
	}
	if bm["cinder"].Name != "cinder" {
		t.Error("Test bm[\"cinder\"].Name error")
	}
	if bm["sample"].Name != "sample" {
		t.Error("Test bm[\"sample\"].Name error")
	}
	if _, ok := bm["lvm"]; !ok {
		t.Error("Test bm[\"lvm\"].Name error")
	}
}
