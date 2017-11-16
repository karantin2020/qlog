package qlog

import (
	"encoding/json"
	"fmt"
	"github.com/karantin2020/qlog/buffer"
	"github.com/pkg/errors"
	"time"
)

type Logger struct {
	// Logs are written to that destination
	Output []Output
	// Pointer to notepad parent instance
	Notepad *Notepad

	// // Logs are written to that destination
	// Out io.Writer

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
	// Context buffer
	CtxBuffer []*buffer.Buffer
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
	// All log entries pass through the formatter before logged to Out. The
	// included formatters are `TextFormatter` and `JSONFormatter` for which
	// TextFormatter is the default. In development (when a TTY is attached) it
	// logs with colors, but to a file it wouldn't. You can easily implement your
	// own that implements the `Formatter` interface, see the `README` or included
	// formatters for examples.
	Formatter []Formatter
	// Notepad context
	Context []Field
	// Context buffer
	CtxBuffer []*buffer.Buffer
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

	// ErrorFunc generates error field from message if no error was passed to Entry
	ErrorFunc func(string) error

	// TimeFieldFormat defines the time format of the Time field type.
	// If set to an empty string, the time is formatted as an UNIX timestamp
	// as integer.
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
	l.Hooks = append(l.Hooks, h)
}

// New func returns new instance of notepad
func New(lvl Level /*outLevel, errLevel Level, outHandle, errHandle io.Writer*/) *Notepad {
	if lvl > _maxLevel || lvl < _minLevel {
		panic("Logging level is out of range")
	}
	n := &Notepad{}
	n.Loggers = [7]**Logger{&n.DEBUG, &n.INFO, &n.WARN, &n.ERROR, &n.CRITICAL, &n.PANIC, &n.FATAL}
	n.Level = lvl
	n.Formatter = make([]Formatter, 0, 3)
	n.Context = make([]Field, 0, 10)
	n.Options = LogConfig{
		TimestampFieldName:   "time",
		LevelFieldName:       "level",
		MessageFieldName:     "message",
		ErrorFieldName:       "error",
		CallerFieldName:      "caller",
		ErrorFunc:            errors.New,
		TimeFieldFormat:      "2006-01-02T15:04:05.000Z0700", // or time.RFC3339
		TimestampFunc:        time.Now,
		DurationFieldUnit:    time.Millisecond,
		DurationFieldInteger: true,
		InterfaceMarshaler:   json.Marshal,
	}
	n.LOG = n.DEBUG
	n.init()
	return n
}

// init creates the Loggers for each level depending on the notepad levels.
func (n *Notepad) init() {
	for t, logger := range n.Loggers {
		level := Level(t)
		if level >= n.Level {
			*logger = NewLogger(n, level)
			(*logger).AddField(F{n.Options.LevelFieldName, (*logger).Level.String()})
		}
	}
}

func (n *Notepad) copy() *Notepad {
	newnp := *n
	newnp.Formatter = make([]Formatter, len(n.Formatter), cap(n.Formatter))
	copy(newnp.Formatter, n.Formatter)
	newnp.Context = make([]Field, len(n.Context), cap(n.Context))
	copy(newnp.Context, n.Context)
	newnp.CtxBuffer = make([]*buffer.Buffer, len(n.CtxBuffer), cap(n.CtxBuffer))
	copy(newnp.CtxBuffer, n.CtxBuffer)
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

func (np *Notepad) AddHook(lvl Level, h Hook) {
	for t, logger := range np.Loggers {
		if lvl >= Level(t) {
			(*logger).AddHook(h)
		}
	}
}

func (np *Notepad) Hooks(hs ...Hook) *Notepad {
	for _, h := range hs {
		np.AddHook(np.Level, h)
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
	AddField(f, &np.Context, &np.CtxBuffer, &np.Options)
}

func (np *Notepad) WithFields(flds ...F) *Notepad {
	newnp := np.copy()
	for _, f := range flds {
		newnp.AddField(f)
	}
	return newnp
}

func (np *Notepad) Timestamp() *Notepad {
	np.AddField(F{np.Options.TimestampFieldName, np.Options.TimestampFunc()})
	return np
}

func NewLogger(pn *Notepad, level Level) *Logger {
	return &Logger{
		Notepad: pn,
		Hooks:   make([]Hook, 0, 3),
		Level:   level,
		Context: make([]Field, 0, 10),
		Enable:  true,
	}
}

func (l *Logger) copy() *Logger {
	newl := *l
	newl.Context = make([]Field, len(l.Context), cap(l.Context))
	copy(newl.Context, l.Context)
	newl.CtxBuffer = make([]*buffer.Buffer, len(l.CtxBuffer), cap(l.CtxBuffer))
	copy(newl.CtxBuffer, l.CtxBuffer)
	return &newl
}

func (l *Logger) AddField(f F) {
	AddField(f, &l.Context, &l.CtxBuffer, &l.Notepad.Options)
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
	e := np.DEBUG.NewEntry()
	e.Debug(msg)
}

func (np *Notepad) Info(msg string) {
	if np.INFO == nil {
		return
	}
	e := np.INFO.NewEntry()
	e.Info(msg)
}

func (np *Notepad) Warn(msg string) {
	if np.WARN == nil {
		return
	}
	e := np.WARN.NewEntry()
	e.Warn(msg)
}

func (np *Notepad) Error(msg string) {
	if np.ERROR == nil {
		return
	}
	e := np.ERROR.NewEntry()
	e.Error(msg)
}

func (np *Notepad) Critical(msg string) {
	if np.CRITICAL == nil {
		return
	}
	e := np.CRITICAL.NewEntry()
	e.Critical(msg)
}

func (np *Notepad) Panic(msg string) {
	if np.PANIC == nil {
		return
	}
	e := np.PANIC.NewEntry()
	e.Panic(msg)
}

func (np *Notepad) Fatal(msg string) {
	if np.FATAL == nil {
		return
	}
	e := np.FATAL.NewEntry()
	e.Fatal(msg)
}

func (np *Notepad) Log(msg string) {
	if np.LOG == nil {
		return
	}
	e := np.LOG.NewEntry()
	e.Log(msg)
}

func (l *Logger) Msg(msg string) {
	e := l.NewEntry()
	e.Level = l.Level
	e.Message = msg
	e.Process()
}

func (l *Logger) Msgf(format string, a ...interface{}) {
	l.Msg(fmt.Sprintf(format, a...))
}
