package config

import (
	gflag "flag"
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

//Create a Config and init default value.
func GetDefaultConfig() *Config {
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

var CONF *Config = new(Config)
