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
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/go-ini/ini"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

type LogFormatter struct {
	TimestampFormat string
	LogFormat       string
}

const (
	debugLevel             = "debug"
	infoLevel              = "info"
	warnLevel              = "warn"
	errorLevel             = "error"
	path                   = "path"
	level                  = "level"
	format                 = "format"
	configFileName         = "/etc/opensds/opensds.conf"
	defaultLogPath         = "/var/log/opensds"
	defaultLogLevel        = "info"
	unknownHost            = "unknownhost"
	unknownUser            = "unknownuser"
	defaultLogFormat       = "[%time%] [%level%] [%filename%] [%funcName%():%lineNo%] [PID:%process%] %message%"
	defaultTimestampFormat = time.RFC3339
	logSection             = "log"
	tenMb				   = 10
	threeMonth			   = 100
)

func InitLogs() {
	path, level, format := readConfigurationFile()
	configureLogModule(path, level, format)
}

func configureLogModule(path, level, format string) {
	configureWriter(path, format)
	configureLevel(level)
}

func configureWriter(path, format string) error {
	logrus.SetFormatter(&LogFormatter{
		TimestampFormat: defaultTimestampFormat,
		LogFormat:       format + "\n",
	})
	logrus.SetOutput(&lumberjack.Logger{
		Filename: filepath.Join(path, logName()),
		MaxSize: tenMb,
		MaxAge: threeMonth,
		Compress: true,
	})
	return nil
}

func configureLevel(level string) {
	switch level {
	case debugLevel:
		logrus.SetLevel(logrus.DebugLevel)
	case infoLevel:
		logrus.SetLevel(logrus.InfoLevel)
	case warnLevel:
		logrus.SetLevel(logrus.WarnLevel)
	case errorLevel:
		logrus.SetLevel(logrus.ErrorLevel)
	}
	logrus.SetReportCaller(true)
}

func logName() (name string) {
	name = fmt.Sprintf("%s.%s.%s.log",
		filepath.Base(os.Args[0]),
		hostName(),
		userName())
	return name
}

func shortHostname(hostname string) string {
	if i := strings.Index(hostname, "."); i >= 0 {
		return hostname[:i]
	}
	return hostname
}

func hostName() string {
	host := unknownHost
	h, err := os.Hostname()
	if err == nil {
		host = shortHostname(h)
	}
	return host
}

func userName() string {
	userName := unknownUser
	current, err := user.Current()
	if err == nil {
		userName = current.Username
	}
	// Sanitize userName since it may contain filepath separators on Windows.
	userName = strings.Replace(userName, `\`, "_", -1)
	return userName
}

func readConfigurationFile() (cfgPath, cfgLevel, cfgFormat string) {
	cfgPath = defaultLogPath
	cfgLevel = defaultLogLevel
	cfgFormat = defaultLogFormat
	cfg, err := ini.Load(configFileName)
	if err != nil {
		log.Println("Failed to open config file")
		return cfgPath, cfgLevel, cfgFormat
	}
	if cfg.Section(logSection).HasKey(path) {
		cfgPath = cfg.Section(logSection).Key(path).String()
	}
	if cfg.Section(logSection).HasKey(level) {
		cfgLevel = strings.ToLower(cfg.Section(logSection).Key(level).String())
	}
	if cfg.Section(logSection).HasKey(format) {
		cfgFormat = cfg.Section(logSection).Key(format).String()
	}

	return cfgPath, cfgLevel, cfgFormat
}

func (f *LogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	output := f.LogFormat
	if output == "" {
		output = defaultLogFormat
	}

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = defaultTimestampFormat
	}

	output = strings.Replace(output, "%time%", entry.Time.Format(timestampFormat), 1)

	output = strings.Replace(output, "%message%", entry.Message, 1)

	level := strings.ToUpper(entry.Level.String())
	output = strings.Replace(output, "%level%", level, 1)

	output = strings.Replace(output, "%process%", strconv.Itoa(os.Getpid()), 1)

	output = strings.Replace(output, "%filename%", entry.Caller.File, 1)
	output = strings.Replace(output, "%lineNo%", strconv.Itoa(entry.Caller.Line), 1)
	output = strings.Replace(output, "%funcName%", entry.Caller.Function, 1)

	return []byte(output), nil
}
