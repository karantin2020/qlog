package qlog

import (
	"encoding/json"
	"fmt"
	// "github.com/karantin2020/qlog/buffer"
	"errors"
	"time"
)

type Logger struct {
	// Logs are written to that destination
	Output []Output
	// Pointer to notepad parent instance
	Notepad *Notepad
	// Hooks for the logger instance. These allow firing events based on logging
	// levels and log entries. For example, to send errors to an error tracking
	// service, log to StatsD or dump the core on fatal errors.
	Hooks []Hook
	// The logging level the logger should log at. This is typically (and defaults
	// to) `logrus.Info`, which allows Info(), Warn(), Error() and Fatal() to be
	// logged.
	Level Level
	// Logger context
	Context []Field
	// Enable flag
	Enable bool
}

// Notepad is where you leave a note!
type Notepad struct {
	DEBUG    *Logger
	INFO     *Logger
	WARN     *Logger
	ERROR    *Logger
	CRITICAL *Logger
	PANIC    *Logger
	FATAL    *Logger

	LOG *Logger
	// Level is min log level to use
	Level Level
	// Notepad name field
	Name []byte
	// All log entries pass through the formatter before logged to Out. The
	// included formatters are `TextFormatter` and `JSONFormatter` for which
	// TextFormatter is the default. In development (when a TTY is attached) it
	// logs with colors, but to a file it wouldn't. You can easily implement your
	// own that implements the `Formatter` interface, see the `README` or included
	// formatters for examples.
	Formatter []Formatter
	// Notepad context
	Context []Field
	// Loggers is pointer to inner loggers
	Loggers [7]**Logger
	// Options set notebook configs
	Options LogConfig
}

type LogConfig struct {
	// TimestampFieldName is the field name used for the timestamp field.
	TimestampFieldName string // "time"

	// LevelFieldName is the field name used for the level field.
	LevelFieldName string // "level"

	// MessageFieldName is the field name used for the message field.
	MessageFieldName string // "message"

	// ErrorFieldName is the field name used for error fields.
	ErrorFieldName string // "error"

	// CallerFieldName is the field name used for caller fields.
	CallerFieldName string // "caller"

	// FieldsName is the field name used for fields.
	FieldsName string // "fields"

	// ErrorFunc generates error field from message if no error was passed to Entry
	ErrorFunc func(string) error

	// TimeFieldFormat defines the time format of the Time field type.
	// Possible values are:
	// -  predefined time pkg layouts for use in Time.Format and time.Parse
	// -  "Unix" (formats as strconv Int with Time.Unix())
	// -  "UnixNano" (formats as strconv Int with Time.UnixNano())
	// -  "UnixMilli" (formats as strconv Int with Time.UnixNano()/1000000)
	// -  "UnixMicro" (formats as strconv Int with Time.UnixNano()/1000)
	TimeFieldFormat string // time.RFC3339

	// TimestampFunc defines the function called to generate a timestamp.
	TimestampFunc func() time.Time // time.Now

	// DurationFieldUnit defines the unit for time.Duration type fields added
	// using the Dur method.
	DurationFieldUnit time.Duration // time.Millisecond

	// DurationFieldInteger renders Dur fields as integer instead of float if
	// set to true.
	DurationFieldInteger bool // true

	// InterfaceMarshaler is used to marshal arbitrary data fields
	InterfaceMarshaler func(v interface{}) ([]byte, error)
}

func (l *Logger) AddHook(h Hook) {
	// fmt.Printf("Took hook: %#v\n", h)
	// fmt.Printf("Took logger: %#v\n", l)
	// fmt.Printf("Took logger hooks: %#v\n", l.Hooks)
	l.Hooks = append(l.Hooks, h)
}

// New func returns new instance of notepad
func New(name string, lvl uint8, opts ...func(*LogConfig) error) *Notepad {
	chkLevel(lvl)
	n := &Notepad{}
	n.Name = []byte(name)
	n.Loggers = [7]**Logger{&n.DEBUG, &n.INFO, &n.WARN, &n.ERROR, &n.CRITICAL, &n.PANIC, &n.FATAL}
	n.Level = InitLevel(lvl)
	n.Formatter = make([]Formatter, 0, 3)
	n.Context = make([]Field, 0, 7)
	n.Options = LogConfig{
		TimestampFieldName:   "time",
		LevelFieldName:       "level",
		MessageFieldName:     "message",
		ErrorFieldName:       "error",
		CallerFieldName:      "caller",
		FieldsName:           "fields",
		ErrorFunc:            errors.New,
		TimeFieldFormat:      "2006-01-02T15:04:05.000Z0700", // or time.RFC3339
		TimestampFunc:        time.Now,
		DurationFieldUnit:    time.Millisecond,
		DurationFieldInteger: true,
		InterfaceMarshaler:   json.Marshal,
	}
	for _, fn := range opts {
		fn(&n.Options)
	}
	n.LOG = n.DEBUG
	n.init()
	return n
}

// init creates the Loggers for each level depending on the notepad levels.
func (n *Notepad) init() {
	for t, logger := range n.Loggers {
		level := uint8(t)
		if level >= n.Level.n {
			*logger = NewLogger(n, level)
			(*logger).Output = make([]Output, 0, 3)
		}
	}
}

func (n *Notepad) copy() *Notepad {
	newnp := *n
	newnp.Formatter = make([]Formatter, len(n.Formatter), cap(n.Formatter))
	copy(newnp.Formatter, n.Formatter)
	newnp.Context = make([]Field, 0, 7)
	newnp.Context = append(newnp.Context, n.Context...)
	// newnp.CtxBuffer = make([]*buffer.Buffer, len(n.CtxBuffer), cap(n.CtxBuffer))
	// copy(newnp.CtxBuffer, n.CtxBuffer)
	newnp.Loggers = [7]**Logger{&newnp.DEBUG,
		&newnp.INFO, &newnp.WARN, &newnp.ERROR,
		&newnp.CRITICAL, &newnp.PANIC, &newnp.FATAL}
	for tl, logger := range n.Loggers {
		if *logger != nil {
			*newnp.Loggers[tl] = (*logger).copy()
			(*newnp.Loggers[tl]).Notepad = &newnp
		}
	}

	return &newnp
}

func (n *Notepad) free() {
	for j := range n.Loggers {
		(*n.Loggers[j]).free()
	}
	// for i := range n.CtxBuffer {
	// 	n.CtxBuffer[i].Free()
	// }
}

func TimeFormat(format string) func(*LogConfig) error {
	return func(lc *LogConfig) error {
		lc.TimeFieldFormat = format
		return nil
	}
}

func (np *Notepad) SetTimeFormat(format string) *Notepad {
	np.Options.TimeFieldFormat = format
	return np
}

func (np *Notepad) SetLevel(lvl uint8) {
	chkLevel(lvl)
	np.Level = InitLevel(lvl)
	np.init()
}

func (np *Notepad) AddHook(lvl uint8, h Hook) {
	chkLevel(lvl)
	(*np.Loggers[int(lvl)]).AddHook(h)
}

func (np *Notepad) AddHooks(lvl uint8, hs ...Hook) *Notepad {
	for t, logger := range np.Loggers {
		if *logger != nil && uint8(t) >= lvl {
			for _, h := range hs {
				np.AddHook(np.Level.n, h)
			}
		}
	}
	return np
}

func (np *Notepad) SetOutput(fns ...func(*Notepad)) *Notepad {
	for _, fn := range fns {
		fn(np)
	}
	return np
}

func (np *Notepad) AddField(f F) {
	AddField(f, &np.Context, &np.Options)
}

func (np *Notepad) WithFields(flds ...F) *Notepad {
	newnp := np.copy()
	for _, f := range flds {
		newnp.AddField(f)
	}
	return newnp
}

func (np *Notepad) Timestamp() *Notepad {
	np.AddField(F{Key: np.Options.TimestampFieldName, Value: np.Options.TimestampFunc()})
	return np
}

func NewLogger(np *Notepad, level uint8) *Logger {
	lgr := &Logger{
		Notepad: np,
		Hooks:   make([]Hook, 0, 3),
		Level:   InitLevel(level),
		// Context: make([]Field, 0, 10),
		Enable: true,
	}
	lgr.Context = make([]Field, 0, 7)
	return lgr
}

func (l *Logger) copy() *Logger {
	newl := *l
	newl.Context = make([]Field, 0, 7)
	newl.Context = append(newl.Context, l.Context...)
	// newl.CtxBuffer = make([]*buffer.Buffer, len(l.CtxBuffer), cap(l.CtxBuffer))
	// copy(newl.CtxBuffer, l.CtxBuffer)
	return &newl
}

func (l *Logger) free() {
	// for i := range l.CtxBuffer {
	// 	l.CtxBuffer[i].Free()
	// }
}

func (l *Logger) AddField(f F) {
	AddField(f, &l.Context, &l.Notepad.Options)
}

func (l *Logger) WithFields(flds ...F) *Logger {
	for _, f := range flds {
		l.AddField(f)
	}
	return l
}

func (np *Notepad) Debug(msg string) {
	if np.DEBUG == nil {
		return
	}
	np.DEBUG.NewEntry().Debug(msg)
}

func (np *Notepad) Debugf(format string, a ...interface{}) {
	np.Debug(fmt.Sprintf(format, a...))
}

func (np *Notepad) Info(msg string) {
	if np.INFO == nil {
		return
	}
	np.INFO.NewEntry().Info(msg)
}

func (np *Notepad) Infof(format string, a ...interface{}) {
	np.Info(fmt.Sprintf(format, a...))
}

func (np *Notepad) Warn(msg string) {
	if np.WARN == nil {
		return
	}
	np.WARN.NewEntry().Warn(msg)
}

func (np *Notepad) Warnf(format string, a ...interface{}) {
	np.Warn(fmt.Sprintf(format, a...))
}

func (np *Notepad) Error(msg string) {
	if np.ERROR == nil {
		return
	}
	np.ERROR.NewEntry().Error(msg)
}

func (np *Notepad) Errorf(format string, a ...interface{}) {
	np.Error(fmt.Sprintf(format, a...))
}

func (np *Notepad) Critical(msg string) {
	if np.CRITICAL == nil {
		return
	}
	np.CRITICAL.NewEntry().Critical(msg)
}

func (np *Notepad) Criticalf(format string, a ...interface{}) {
	np.Critical(fmt.Sprintf(format, a...))
}

func (np *Notepad) Panic(msg string) {
	if np.PANIC == nil {
		return
	}
	np.PANIC.NewEntry().Panic(msg)
}

func (np *Notepad) Panicf(format string, a ...interface{}) {
	np.Panic(fmt.Sprintf(format, a...))
}

func (np *Notepad) Fatal(msg string) {
	if np.FATAL == nil {
		return
	}
	np.FATAL.NewEntry().Fatal(msg)
}

func (np *Notepad) Fatalf(format string, a ...interface{}) {
	np.Fatal(fmt.Sprintf(format, a...))
}

func (np *Notepad) Log(msg string) {
	if np.LOG == nil {
		return
	}
	np.LOG.NewEntry().Log(msg)
}

func (np *Notepad) Logf(format string, a ...interface{}) {
	np.Log(fmt.Sprintf(format, a...))
}

func (l *Logger) Msg(msg string) {
	e := l.NewEntry()
	// e.Level = l.Level
	e.Message = append(e.Message, Str2Bytes(msg)...)
	e.Process()
}

func (l *Logger) Msgf(format string, a ...interface{}) {
	l.Msg(fmt.Sprintf(format, a...))
}

func chkLevel(lvl uint8) {
	if lvl > _maxLevel || lvl < _minLevel {
		panic("Logging level is out of range")
	}
}
