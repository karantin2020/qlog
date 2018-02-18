package qlog

import (
	"fmt"
)

// A Level is a logging priority. Higher levels are more important.
type Level struct {
	n uint8
	b []byte
}

const (
	// DebugLevel logs are typically voluminous, and are usually disabled in
	// production.
	DebugLevel uint8 = iota
	// InfoLevel is the default logging priority.
	InfoLevel
	// WarnLevel logs are more important than Info, but don't need individual
	// human review.
	WarnLevel
	// ErrorLevel logs are high-priority. If an application is running smoothly,
	// it shouldn't generate any error-level logs.
	ErrorLevel
	// CriticalLevel logs are particularly important errors. In development the
	// logger panics after writing the message.
	CriticalLevel
	// PanicLevel logs a message, then panics.
	PanicLevel
	// FatalLevel logs a message, then calls os.Exit(1).
	FatalLevel

	_minLevel = DebugLevel
	_maxLevel = FatalLevel
)

var (
	_debug    = []byte("debug")
	_info     = []byte("info")
	_warn     = []byte("warn")
	_error    = []byte("error")
	_critical = []byte("critical")
	_panic    = []byte("panic")
	_fatal    = []byte("fatal")
)

func InitLevel(lvl uint8) Level {
	l := Level{n: lvl}
	l.b = l.ToBytes()
	return l
}

// String returns a lower-case ASCII representation of the log level.
func (l Level) String() string {
	switch l.n {
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warn"
	case ErrorLevel:
		return "error"
	case CriticalLevel:
		return "critical"
	case PanicLevel:
		return "panic"
	case FatalLevel:
		return "fatal"
	default:
		return fmt.Sprintf("Level(%d)", l)
	}
}

func (l Level) AppendBytes(dst []byte) []byte {
	switch l.n {
	case DebugLevel:
		return AppendStringNoQuotes(dst, "debug")
	case InfoLevel:
		return AppendStringNoQuotes(dst, "info")
	case WarnLevel:
		return AppendStringNoQuotes(dst, "warn")
	case ErrorLevel:
		return AppendStringNoQuotes(dst, "error")
	case CriticalLevel:
		return AppendStringNoQuotes(dst, "critical")
	case PanicLevel:
		return AppendStringNoQuotes(dst, "panic")
	case FatalLevel:
		return AppendStringNoQuotes(dst, "fatal")
	default:
		return AppendStringNoQuotes(dst, fmt.Sprintf("Level(%d)", l))
	}
}

func (l Level) ToBytes() []byte {
	switch l.n {
	case DebugLevel:
		return _debug
	case InfoLevel:
		return _info
	case WarnLevel:
		return _warn
	case ErrorLevel:
		return _error
	case CriticalLevel:
		return _critical
	case PanicLevel:
		return _panic
	case FatalLevel:
		return _fatal
	default:
		return Str2Bytes(fmt.Sprintf("Level(%d)", l))
	}
}

// CapitalString returns an all-caps ASCII representation of the log level.
func (l Level) CapitalString() string {
	// Printing levels in all-caps is common enough that we should export this
	// functionality.
	switch l.n {
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	case CriticalLevel:
		return "CRITICAL"
	case PanicLevel:
		return "PANIC"
	case FatalLevel:
		return "FATAL"
	default:
		return fmt.Sprintf("LEVEL(%d)", l)
	}
}

// MarshalText marshals the Level to text. Note that the text representation
// drops the -Level suffix (see example).
func (l Level) MarshalText() ([]byte, error) {
	return []byte(l.String()), nil
}
