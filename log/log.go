package log

import (
	"github.com/karantin2020/qlog"
)

var (
	DEBUG    *qlog.Logger
	INFO     *qlog.Logger
	WARN     *qlog.Logger
	ERROR    *qlog.Logger
	CRITICAL *qlog.Logger
	PANIC    *qlog.Logger
	FATAL    *qlog.Logger

	LOG *qlog.Logger

	defaultNotepad *qlog.Notepad
)

func init() {
	defaultNotepad = qlog.New("", qlog.InfoLevel).
		SetOutput(qlog.Template("${time}\t${LEVEL}\t${message}\t${fields}\n"))
	reloadDefaultNotepad()
}

func reloadDefaultNotepad() {
	DEBUG = defaultNotepad.DEBUG
	INFO = defaultNotepad.INFO
	WARN = defaultNotepad.WARN
	ERROR = defaultNotepad.ERROR
	CRITICAL = defaultNotepad.CRITICAL
	PANIC = defaultNotepad.PANIC
	FATAL = defaultNotepad.FATAL

	LOG = defaultNotepad.LOG
}

// log.Print("hello world")
// Output: {"level":"debug","time":2017-11-04T15:09:54+00:00,"message":"hello world"}
//
// log.Info("hello world")
// Output: {"level":"info","time":2017-11-04T15:09:54+00:00,"message":"hello world"}
//
// log.INFO.Fields(qlog.F{"service", service}).Msg("hello world")
// Output: {"level":"info","time":2017-11-04T15:09:54+00:00,"message":"hello world", "service":"myservice"}
//
// log.
//     Err(err).
//     Fields(qlog.F{"service", service}).
//     Fatalf("Cannot start %s", service)
// Output: {"level":"fatal","time":2017-11-04T15:09:54+00:00,"message":"Cannot start myservice","error":"some error","service":"myservice"}
// Exit 1
//
// New Log instance to manage different outputs
// logger := qlog.New(qlog.InfoLevel).
// 	SetOutput(Template("${time}\t${level}\t${source}\t${message}\t${fields}"), ...).
// 	WithFields(
// 		qlog.F{"module", "main"},
// 	)
// logger.Fields(qlog.F{"foo", "bar"}).Info("hello world")
// Output: 2017-11-04T15:09:54+00:00    INFO    main.go:15    hello world     {"module":"main","foo":"bar"}
//
// Formatter default is JSON (use github.com/json-iterator/go). May be Template, Text, nil (to use internal format for Writer) etc
//
// Sub-loggers let you chain loggers with additional context
// flog := log.WithFields(
// 	qlog.F{"module", "foo"}
// )
// flog.Info("hello world")
// Output: {"level":"info","time":2017-11-04T15:09:54+00:00,"message":"hello world","module":"foo"}

func SetLevel(lvl uint8) {
	defaultNotepad.SetLevel(lvl)
	defaultNotepad.Level = qlog.InitLevel(lvl)
	reloadDefaultNotepad()
}
func WithFields(flds ...qlog.F) *qlog.Notepad {
	return defaultNotepad.WithFields(flds...)
}

func Debug(msg string) {
	defaultNotepad.Debug(msg)
}
func Debugf(format string, a ...interface{}) {
	defaultNotepad.Debugf(format, a...)
}
func Info(msg string) {
	defaultNotepad.Info(msg)
}
func Infof(format string, a ...interface{}) {
	defaultNotepad.Infof(format, a...)
}
func Warn(msg string) {
	defaultNotepad.Warn(msg)
}
func Warnf(format string, a ...interface{}) {
	defaultNotepad.Warnf(format, a...)
}
func Error(msg string) {
	defaultNotepad.Error(msg)
}
func Errorf(format string, a ...interface{}) {
	defaultNotepad.Errorf(format, a...)
}
func Critical(msg string) {
	defaultNotepad.Critical(msg)
}
func Criticalf(format string, a ...interface{}) {
	defaultNotepad.Criticalf(format, a...)
}
func Panic(msg string) {
	defaultNotepad.Panic(msg)
}
func Panicf(format string, a ...interface{}) {
	defaultNotepad.Panicf(format, a...)
}
func Fatal(msg string) {
	defaultNotepad.Fatal(msg)
}
func Fatalf(format string, a ...interface{}) {
	defaultNotepad.Fatalf(format, a...)
}
func Log(msg string) {
	defaultNotepad.Log(msg)
}
func Logf(format string, a ...interface{}) {
	defaultNotepad.Logf(format, a...)
}
