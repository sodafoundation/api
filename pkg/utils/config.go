package utils

import (
	"github.com/Unknwon/goconfig"
	"log"
	"reflect"
	"strconv"
	"strings"
)

const (
	CONF_NAME=iota
	CONF_DEFULT_VLAUE
)

type OsdsLet struct {
	ApiEndpoint string `conf:"api_endpoint,localhost:50040"`
	Graceful    bool   `conf:"graceful,true"`
	LogFile     string `conf:"log_file,/var/log/opensds/osdslet.log"`
	SocketOrder string `conf:"socket_order"`
}

type OsdsDock struct {
	ApiEndpoint string `conf:"api_endpoint,localhost:50050"`
	LogFile     string `conf:"log_file,/var/log/opensds/osdsdock.log"`
}

type Database struct {
	Credential string `conf:"credential"`
	Driver     string `conf:"driver,etcd"`
	Endpoint   string `conf:"endpoint,localhost:2379,localhost:2380"`
}

type Default struct {
}

type Config struct {
	Default    `conf:"default"`
	OsdsLet    `conf:"osdslet"`
	OsdsDock   `conf:"osdsdock"`
	Database   `conf:"database"`
}

func setSectionValue(section string, v reflect.Value, cfg *goconfig.ConfigFile) {
	for i := 0; i < v.Type().NumField(); i++ {

		field := v.Field(i)
		tag := v.Type().Field(i).Tag.Get("conf")
		tags := strings.SplitN(tag, ",", 2)
		if !field.CanSet() {
			continue
		}

		var strVal = ""
		if len(tags) > 1 {
			strVal = tags[CONF_DEFULT_VLAUE]
		}
		if cfg != nil {
			strVal, _ = cfg.GetValue(section, tags[CONF_NAME])
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
	cfg, err := goconfig.LoadConfigFile(confFile)
	if err != nil {
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

var CONF * Config = new(Config)
func (conf *Config)Load(confFile string) {
		initConf(confFile, CONF)
}

