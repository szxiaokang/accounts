/**
 * @project Accounts
 * @filename logs.go
 * @author kangyun@outlook.com
 * @copyright Copyright (C) kangyun@outlook.com
 * @datetime 2023/2/19 10:22
 * @version 1.0
 * @description
 * 日志相关函数
 */

package base

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

const (
	LogLevelTrace = 0
	LogLevelDebug = 1
	LogLevelInfo  = 2
	LogLevelError = 3
	LogLevelFatal = 4
	LogLevelCount = 5
)

var (
	strLevels = []string{"trace", "debug", "info", "error", "fatal"}

	defaultLogger = &MyLogger{}

	inited            bool
	needOutputConsole bool
)

type MyLogger struct {
	inited bool

	logLevel int
	loggers  [LogLevelCount]*logWrapper

	fileWriter
}

type RawLogger struct {
	fileWriter
	logger *log.Logger
}

type fileWriter struct {
	logRoot   string
	fileName  string
	isDataLog bool //是否是数据日志

	f                   *os.File
	mutex               sync.RWMutex
	lastUpdateHourInDay int64
}

func NewFileWriter(path, name string, isDataLog bool) *fileWriter {
	return &fileWriter{logRoot: path, fileName: name, isDataLog: isDataLog}
}

func (self *fileWriter) Write(p []byte) (int, error) {
	now := time.Now()
	hoursInDay := now.Unix() / 3600
	self.mutex.Lock()
	defer self.mutex.Unlock()
	if self.f == nil || hoursInDay != self.lastUpdateHourInDay { // 15:04:05
		var filename string
		if self.isDataLog {
			//数据日志文件格式单独处理
			prefix := 1
			if now.Minute() > 30 {
				prefix = 2
			}
			filename = fmt.Sprintf("%s.log.%s0%d", self.fileName, now.Format("2006010215"), prefix)
		} else {
			filename = fmt.Sprintf("%s%s.log", self.fileName, now.Format("2006010215"))
		}
		f, err := os.OpenFile(fmt.Sprintf("%s%s", self.logRoot, filename), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600|0644)

		if err != nil {
			return 0, err
		}
		if self.f != nil {
			self.f.Close() //nolint:gosec
		}
		self.f = f
		self.lastUpdateHourInDay = hoursInDay
	}
	n, err := self.f.Write(p)
	if err != nil {
		self.f = nil
	}
	return n, err
}

type logWrapper struct {
	logger *log.Logger
}

func (self *logWrapper) initLogger(w io.Writer, prefix string) {
	self.logger = log.New(w, prefix, log.Lshortfile|log.LstdFlags)
}

func (self *logWrapper) log(text string) {
	if !inited {
		fmt.Println("!!!->", text)
		return
	}

	// self.logger.Println(text)
	err := self.logger.Output(4, text)
	if err != nil {
		fmt.Println("logger error|", err)
	}
	if needOutputConsole {
		fmt.Println(self.logger.Prefix(), time.Now().Format("15:04:05"), text)
	}
}

func Fatal(v ...interface{}) {
	defaultLogger.Fatal(v...)
}

func Error(v ...interface{}) {
	defaultLogger.Error(v...)
}

func Errorf(format string, v ...interface{}) {
	defaultLogger.Errorf(format, v...)
}

func Info(v ...interface{}) {
	defaultLogger.Info(v...)
}

func Infof(format string, v ...interface{}) {
	defaultLogger.Infof(format, v...)
}

func Debug(v ...interface{}) {
	defaultLogger.Debug(v...)
}
func Debugf(format string, v ...interface{}) {
	defaultLogger.Debugf(format, v...)
}

func Trace(v ...interface{}) {
	defaultLogger.Trace(v...)
}
func Tracef(format string, v ...interface{}) {
	defaultLogger.Tracef(format, v...)
}

func Println(v ...interface{}) {
	fmt.Println(v...)
}

func getLevelByStr(strLevel string) int {
	for i, v := range strLevels {
		if v == strLevel {
			return i
		}
	}
	return LogLevelDebug
}

func getLevelStr(level int) string {
	if level >= 0 && level < len(strLevels) {
		return strLevels[level]
	}
	return strLevels[0]
}

func Start(root string, l string) error {
	if inited {
		return nil
	}
	inited = true

	if len(root) > 0 {
		err := os.MkdirAll(root, 0750)
		if err != nil {
			return err
		}
	}
	defaultLogger.Init(root, l, "")

	// loggers[logFileTypeError].fileName = "error"
	needOutputConsole = runtime.GOOS == "windows"

	Info("put log to:", root, " log level:", l)
	return nil
}

// return logRoot, logLevel
func GetDesc() (string, string) {
	return defaultLogger.logRoot, getLevelStr(defaultLogger.logLevel)
}

func SetLevel(l string) {
	defaultLogger.SetLevel(l)
}

// ////////////////////
func (self *MyLogger) Init(logRoot string, logLevel, filename string) {
	self.SetLevel(logLevel)
	self.logRoot = logRoot
	self.fileName = filename
	if self.inited {
		return
	}
	self.inited = true

	for k := range self.loggers {
		v := &logWrapper{}
		v.initLogger(self, strings.ToUpper("\n["+getLevelStr(k)+"]"))
		self.loggers[k] = v
	}
}

func concatArgs(v ...interface{}) string {
	o := fmt.Sprintln(v...)
	lo := len(o)
	if lo > 0 {
		return o[:lo-1]
	}
	return o
}
func (self *MyLogger) Fatal(v ...interface{}) {
	if self.logLevel <= LogLevelFatal {
		self.loggers[LogLevelFatal].log(concatArgs(v...))
		panic("log fatal")
	}
}

func (self *MyLogger) Error(v ...interface{}) {
	if self.logLevel <= LogLevelError {
		self.loggers[LogLevelError].log(concatArgs(v...))
	}
}

func (self *MyLogger) Errorf(format string, v ...interface{}) {
	if self.logLevel <= LogLevelError {
		self.loggers[LogLevelError].log(fmt.Sprintf(format, v...))
	}
}

func (self *MyLogger) Info(v ...interface{}) {
	if self.logLevel <= LogLevelInfo {
		self.loggers[LogLevelInfo].log(concatArgs(v...))
	}
}

func (self *MyLogger) Infof(format string, v ...interface{}) {
	if self.logLevel <= LogLevelInfo {
		self.loggers[LogLevelInfo].log(fmt.Sprintf(format, v...))
	}
}

func (self *MyLogger) Debug(v ...interface{}) {
	if self.logLevel <= LogLevelDebug {
		self.loggers[LogLevelDebug].log(concatArgs(v...))
	}
}
func (self *MyLogger) Debugf(format string, v ...interface{}) {
	if self.logLevel <= LogLevelDebug {
		self.loggers[LogLevelDebug].log(fmt.Sprintf(format, v...))
	}
}

func (self *MyLogger) Trace(v ...interface{}) {
	if self.logLevel <= LogLevelTrace {
		self.loggers[LogLevelTrace].log(concatArgs(v...))
	}
}
func (self *MyLogger) Tracef(format string, v ...interface{}) {
	if self.logLevel <= LogLevelTrace {
		self.loggers[LogLevelTrace].log(fmt.Sprintf(format, v...))
	}
}

func (self *MyLogger) SetLevel(l string) {
	self.logLevel = getLevelByStr(l)
}

type LogInst struct {
	*MyLogger
	inited            bool
	needOutputConsole bool
}

func (self *RawLogger) Init(logRoot string, filename string, prefix string, flag int) {
	self.logRoot = logRoot
	self.fileName = filename
	self.logger = log.New(self, prefix, flag)
}

func (self *RawLogger) Log(v string) {
	err := self.logger.Output(2, v)
	if err != nil {
		fmt.Println("unknown log err|", err)
	}
}

func (self *LogInst) Start(root string, l string) error {
	if self.inited {
		return nil
	}
	self.inited = true

	self.MyLogger = &MyLogger{}

	if len(root) > 0 {
		err := os.MkdirAll(root, 0755) //nolint:gosec
		if err != nil {
			return err
		}
	}
	self.MyLogger.Init(root, l, "")

	self.needOutputConsole = runtime.GOOS == "windows"

	Info("put log to:", root, " log level:", l)
	return nil
}
