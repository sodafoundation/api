package obs

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

type Level int

const (
	LEVEL_OFF   Level = 500
	LEVEL_ERROR Level = 400
	LEVEL_WARN  Level = 300
	LEVEL_INFO  Level = 200
	LEVEL_DEBUG Level = 100
)

var cacheCount = 50

var logLevelMap = map[Level]string{
	LEVEL_OFF:   "[OFF]: ",
	LEVEL_ERROR: "[ERROR]: ",
	LEVEL_WARN:  "[WARN]: ",
	LEVEL_INFO:  "[INFO]: ",
	LEVEL_DEBUG: "[DEBUG]: ",
}

type logConfType struct {
	level        Level
	logToConsole bool
	logFullPath  string
	maxLogSize   int64
	backups      int
}

func getDefaultLogConf() logConfType {
	return logConfType{
		level:        LEVEL_WARN,
		logToConsole: false,
		logFullPath:  "",
		maxLogSize:   1024 * 1024 * 30, //30MB
		backups:      10,
	}
}

var logConf logConfType

type loggerWrapper struct {
	fullPath string
	fd       *os.File
	queue    []string
	logger   *log.Logger
	index    int
	lock     *sync.RWMutex
}

func (lw *loggerWrapper) doInit() {
	lw.queue = make([]string, 0, cacheCount)
	lw.logger = log.New(lw.fd, "", 0)
	lw.lock = new(sync.RWMutex)
}

func (lw *loggerWrapper) rotate() {
	stat, err := lw.fd.Stat()
	if err == nil && stat.Size() >= logConf.maxLogSize {
		lw.fd.Sync()
		lw.fd.Close()

		if lw.index > logConf.backups {
			lw.index = 1
		}
		os.Rename(lw.fullPath, lw.fullPath+"."+IntToString(lw.index))
		lw.index += 1

		fd, err := os.OpenFile(lw.fullPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}
		lw.fd = fd
		lw.logger.SetOutput(lw.fd)
	}
}

func (lw *loggerWrapper) doFlush() {
	lw.rotate()
	for _, m := range lw.queue {
		lw.logger.Println(m)
	}
	lw.fd.Sync()
}

func (lw *loggerWrapper) doClose() {
	lw.doFlush()
	lw.fd.Close()
	lw.queue = nil
	lw.fd = nil
	lw.logger = nil
	lw.lock = nil
	lw.fullPath = ""
}

func (lw *loggerWrapper) printfWithLock(msg string) {
	lw.lock.Lock()
	defer lw.lock.Unlock()
	if len(lw.queue) >= cacheCount {
		lw.doFlush()
		lw.queue = make([]string, 0, cacheCount)
	} else {
		lw.queue = append(lw.queue, msg)
	}
}

func (lw *loggerWrapper) Printf(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	if len(lw.queue) >= cacheCount {
		lw.printfWithLock(msg)
	} else {
		lw.lock.RLock()
		if len(lw.queue) >= cacheCount {
			lw.lock.RUnlock()
			lw.printfWithLock(msg)
		} else {
			defer lw.lock.RUnlock()
			lw.queue = append(lw.queue, msg)
		}
	}
}

var consoleLogger *log.Logger
var fileLogger *loggerWrapper

var lock *sync.RWMutex = new(sync.RWMutex)

func isDebugLogEnabled() bool {
	return logConf.level <= LEVEL_DEBUG
}

func isErrorLogEnabled() bool {
	return logConf.level <= LEVEL_ERROR
}

func isWarnLogEnabled() bool {
	return logConf.level <= LEVEL_WARN
}

func isInfoLogEnabled() bool {
	return logConf.level <= LEVEL_INFO
}

func reset() {
	if fileLogger != nil {
		fileLogger.doClose()
		fileLogger = nil
	}
	consoleLogger = nil
	logConf = getDefaultLogConf()
}

func InitLog(logFullPath string, maxLogSize int64, backups int, level Level, logToConsole bool) error {
	return InitLogWithCacheCnt(logFullPath, maxLogSize, backups, level, logToConsole, 50)
}

func InitLogWithCacheCnt(logFullPath string, maxLogSize int64, backups int, level Level, logToConsole bool, cacheCnt int) error {
	lock.Lock()
	defer lock.Unlock()
	if cacheCnt <= 0 {
		cacheCnt = 50
	}
	cacheCount = cacheCnt
	reset()
	if fullPath := strings.TrimSpace(logFullPath); fullPath != "" {
		_fullPath, err := filepath.Abs(fullPath)
		if err != nil {
			return err
		}

		if !strings.HasSuffix(_fullPath, ".log") {
			_fullPath += ".log"
		}

		stat, err := os.Stat(_fullPath)
		if err == nil && stat.IsDir() {
			return errors.New(fmt.Sprintf("logFullPath:[%s] is a directory", _fullPath))
		} else if err := os.MkdirAll(filepath.Dir(_fullPath), os.ModePerm); err != nil {
			return err
		}

		fd, err := os.OpenFile(_fullPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return err
		}
		fileLogger = &loggerWrapper{fullPath: _fullPath, fd: fd, index: 1}

		if stat == nil {
			stat, err = os.Stat(_fullPath)
		}
		prefix := stat.Name() + "."
		walkFunc := func(path string, info os.FileInfo, err error) error {
			if name := info.Name(); strings.HasPrefix(name, prefix) {
				if i := StringToInt(name[len(prefix):], 0); i >= fileLogger.index {
					fileLogger.index = i + 1
				}
			}
			return nil
		}

		filepath.Walk(filepath.Dir(_fullPath), walkFunc)
		fileLogger.doInit()
	}
	if maxLogSize > 0 {
		logConf.maxLogSize = maxLogSize
	}
	if backups > 0 {
		logConf.backups = backups
	}
	logConf.level = level
	if logToConsole {
		consoleLogger = log.New(os.Stdout, "", log.LstdFlags)
	}
	return nil
}

func CloseLog() {
	if fileLogger != nil || consoleLogger != nil {
		lock.Lock()
		defer lock.Unlock()
		reset()
	}
}

func SyncLog() {
	if fileLogger != nil {
		lock.Lock()
		defer lock.Unlock()
		fileLogger.doFlush()
	}
}

func logEnabled() bool {
	return consoleLogger != nil || fileLogger != nil
}

func doLog(level Level, format string, v ...interface{}) {
	if logEnabled() && logConf.level <= level {
		lock.RLock()
		defer lock.RUnlock()
		msg := fmt.Sprintf(format, v...)
		if _, file, line, ok := runtime.Caller(1); ok {
			index := strings.LastIndex(file, "/")
			if index >= 0 {
				file = file[index+1:]
			}
			msg = fmt.Sprintf("%s:%d|%s", file, line, msg)
		}
		prefix := logLevelMap[level]
		if consoleLogger != nil {
			consoleLogger.Printf("%s%s", prefix, msg)
		}
		if fileLogger != nil {
			nowDate := FormatUtcNow("2006-01-02T15:04:05Z")
			fileLogger.Printf("%s %s%s", nowDate, prefix, msg)
		}
	}
}
