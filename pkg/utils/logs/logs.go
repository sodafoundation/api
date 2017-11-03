// Copyright (c) 2017 OpenSDS Authors.
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

package logs

import (
	"flag"
	"log"
	"os"

	"github.com/golang/glog"
)

const DefaultLogDir = "/var/log/opensds"

func init() {
	//Set OpenSDS default log directory.
	flag.CommandLine.VisitAll(func(flag *flag.Flag) {
		if flag.Name == "log_dir" {
			flag.DefValue = DefaultLogDir
			flag.Value.Set(DefaultLogDir)
		}
	})
}

type GlogWriter struct{}

func (writer GlogWriter) Write(data []byte) (n int, err error) {
	glog.Info(string(data))
	return len(data), nil
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func InitLogs() {
	log.SetOutput(GlogWriter{})
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	logDir := flag.CommandLine.Lookup("log_dir").Value.String()
	if exist, _ := PathExists(logDir); !exist {
		os.MkdirAll(logDir, 0755)
	}
}

func FlushLogs() {
	glog.Flush()
}
