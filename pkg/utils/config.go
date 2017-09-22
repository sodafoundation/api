// Copyright (c) 2016 OpenSDS Authors.
//
//    Licensed under the Apache License, Version 2.0 (the "License"); you may
//    not use this file except in compliance with the License. You may obtain
//    a copy of the License at
//
//         http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//    WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//    License for the specific language governing permissions and limitations
//    under the License.
package utils

import (
	gflag "flag"
	"github.com/go-ini/ini"
	"log"
	"reflect"
	"strconv"
	"strings"
)

const (
	ConfKeyName = iota
	ConfDefaultValue
)

type OsdsLet struct {
	ApiEndpoint string `conf:"api_endpoint,localhost:50040"`
	Graceful    bool   `conf:"graceful,true"`
	SocketOrder string `conf:"socket_order"`
}

type OsdsDock struct {
	ApiEndpoint string `conf:"api_endpoint,localhost:50050"`
}

type Database struct {
	Credential string `conf:"credential,username:password@tcp(ip:port)/dbname"`
	Driver     string `conf:"driver,etcd"`
	Endpoint   string `conf:"endpoint,localhost:2379,localhost:2380"`
}

type Default struct {
}

type Config struct {
	Default  `conf:"default"`
	OsdsLet  `conf:"osdslet"`
	OsdsDock `conf:"osdsdock"`
	Database `conf:"database"`
	Flag     FlagSet
}

func setSectionValue(section string, v reflect.Value, cfg *ini.File) {
	for i := 0; i < v.Type().NumField(); i++ {

		field := v.Field(i)
		tag := v.Type().Field(i).Tag.Get("conf")
		tags := strings.SplitN(tag, ",", 2)
		if !field.CanSet() {
			continue
		}

		var strVal = ""
		if cfg != nil {
			key, _ := cfg.Section(section).GetKey(tags[ConfKeyName])
			strVal = key.Value()
		} else if len(tags) > 1 {
			strVal = tags[ConfDefaultValue]
		} else {
			continue
		}

		switch field.Kind() {
		case reflect.Bool:
			val, _ := strconv.ParseBool(strVal)
			field.SetBool(val)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			val, _ := strconv.ParseInt(strVal, 10, 64)
			field.SetInt(val)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			val, _ := strconv.ParseUint(strVal, 10, 64)
			field.SetUint(val)
		case reflect.Float32, reflect.Float64:
			val, _ := strconv.ParseFloat(strVal, 64)
			field.SetFloat(val)
		case reflect.String:
			field.SetString(strVal)
		default:
		}
	}
}

func initConf(confFile string, conf interface{}) {
	cfg, err := ini.Load(confFile)
	if err != nil && confFile != "" {
		log.Println("[Info] Read configuration failed, use default value")
	}
	t := reflect.TypeOf(conf)
	v := reflect.ValueOf(conf)
	for i := 0; i < t.Elem().NumField(); i++ {
		field := v.Elem().Field(i)
		section := t.Elem().Field(i).Tag.Get("conf")
		setSectionValue(section, field, cfg)
	}
}

//New a Config and init default value.
func newConfig() *Config {
	var conf *Config = new(Config)
	initConf("", conf)
	return conf
}

func (c *Config) Load(confFile string) {
	gflag.StringVar(&confFile, "config-file", confFile, "The configuration file of OpenSDS")
	c.Flag.Parse()
	initConf(confFile, CONF)
	c.Flag.AssignValue()
}

var CONF *Config = newConfig()
