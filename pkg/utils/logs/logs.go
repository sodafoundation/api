// Copyright 2017 The OpenSDS Authors.
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

package logs

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/golang/glog"
	"github.com/opensds/opensds/pkg/utils"
)

const DefaultLogDir = "/var/log/opensds"

// flushDaemon periodically flushes the log file buffers.
func flushDaemon(period time.Duration) {
	for range time.NewTicker(period).C {
		glog.Flush()
	}
}

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

// flush log when be interrupted.
func handleInterrupt() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGQUIT)
	go func() {
		sig := <-sigs
		fmt.Println(sig)
		FlushLogs()
		os.Exit(-1)
	}()
}

func InitLogs(LogFlushFrequency time.Duration) {
	log.SetOutput(GlogWriter{})
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	logDir := flag.CommandLine.Lookup("log_dir").Value.String()
	if exist, _ := utils.PathExists(logDir); !exist {
		os.MkdirAll(logDir, 0755)
	}
	glog.Infof("[Info] LogFlushFrequency: %v", LogFlushFrequency)
	go flushDaemon(LogFlushFrequency)
	handleInterrupt()
}

func FlushLogs() {
	glog.Flush()
}

//timeNow: Current system time, accurate to microseconds, for example:
//timeNow := time.Now().UTC().Format("2006-01-02 15:04:05.000000")
//LogPrefix: multi-cloud、hotpot、nbp
//os.Getpid: current process id
//logLevel: Log level,such as INFO、WARNING、ERROR、FATAL.
//filename: go file name.
//line: Line Number.
//funcname: function name
//logMsg: log detail Messages
//errorCode: if logLevel is ERROR ,else errorCode=""
func GetLogContent(logPrefix string, pid string, logLevel string, fileName string, line string, funcName string, logMsg string, errorCode string) string {
	timeNow := time.Now().UTC().Format("2006-01-02 15:04:05.000000")
	var logContent string
	if logLevel != "ERROR" {
		logContent = strings.Join([]string{"[", logLevel, "]", timeNow, logPrefix, "[ pid:", pid, "]", "[", fileName, ":", line, funcName, "]", logMsg}, " ")
	} else {
		logContent = strings.Join([]string{"[", logLevel, "]", "[ ErrorCode:", errorCode, "]", timeNow, logPrefix, "[ pid:", pid, "]", "[", fileName, ":", line, funcName, "]", logMsg}, " ")
	}
	return logContent
}
