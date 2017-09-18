package utils

import (
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

type TestConfig struct {
	TestStruct `conf:"test_struct"`
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
}

func TestOpensdsConfig(t *testing.T) {
	CONF.Load("testdata/opensds.conf")
	if CONF.OsdsLet.ApiEndpoint != "localhost:50040" {
		t.Error("Test OsdsLet.ApiEndpoint error")
	}
	if CONF.OsdsLet.Graceful != true {
		t.Error("Test OsdsLet.Graceful error")
	}
	if CONF.OsdsLet.LogFile != "/var/log/opensds/osdslet.log" {
		t.Error("Test OsdsLet.LogFile error")
	}
	if CONF.OsdsLet.SocketOrder != "inc" {
		t.Error("Test OsdsLet.SocketOrder error")
	}
	if CONF.OsdsDock.ApiEndpoint != "localhost:50050" {
		t.Error("Test OsdsLet.ApiEndpoint error")
	}
	if CONF.OsdsDock.LogFile != "/var/log/opensds/osdsdock.log" {
		t.Error("Test OsdsLet.ApiEndpoint error")
	}
	if CONF.Database.Credential != "opensds:password@127.0.0.1:3306/dbname" {
		t.Error("Test OsdsLet.ApiEndpoint error")
	}
	if CONF.Database.Endpoint != "localhost:2379,localhost:2380" {
		t.Error("Test OsdsLet.ApiEndpoint error")
	}
	if CONF.Database.Driver != "etcd" {
		t.Error("Test OsdsLet.ApiEndpoint error")
	}
}
