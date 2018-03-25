package qlog

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

type JsonOptions struct {
	ErrHandle     io.Writer
	OutHandle     io.Writer
	ErrLevel      uint8
	OutLevel      uint8
	LogName       string
	TimestampName string
	LevelName     string
	MessageName   string
	ErrorName     string
	FieldsName    string
}

func Json(opts ...func(*JsonOptions) error) func(np *Notepad) {
	options := defaultJsonOptions()
	for _, fn := range opts {
		fn(options)
	}
	return func(np *Notepad) {
		if options.OutLevel > _maxLevel || options.OutLevel < _minLevel || options.OutLevel < np.Level.n {
			panic("OutLevel is out of range")
		}
		if options.ErrLevel > _maxLevel || options.ErrLevel < _minLevel || options.ErrLevel < np.Level.n {
			panic("ErrLevel is out of range")
		}
		if options.OutLevel > options.ErrLevel {
			panic("OutLevel is higher than errLevel")
		}
		jsonOut := func(wio io.Writer, topts *JsonOptions) Output {
			return func(e *Entry) {
				bb := bufferPool.Get().(*bytes.Buffer)
				bb.Reset()
				bb.Write([]byte{'{'})
				for i, _ := range e.Data {
					bb.Write([]byte{'"'})
					bb.Write(Str2Bytes(e.Data[i].Key))
					bb.Write([]byte{'"', ':'})
					bb.Write(e.Data[i].Buffer.Bytes())
					if i < len(e.Data)-1 {
						bb.Write([]byte{','})
					}
				}
				bb.Write([]byte{'}'})
				if _, err := wio.Write(bb.Bytes()); err != nil {
					panic(fmt.Sprintf("qlog json logging error: %s", err))
				}
				bufferPool.Put(bb)
			}
		}
		for tlv, logger := range np.Loggers {
			level := uint8(tlv)
			switch {
			case level >= options.ErrLevel:
				(*logger).Output = append((*logger).Output, jsonOut(options.ErrHandle, options))

			case level >= options.OutLevel:
				(*logger).Output = append((*logger).Output, jsonOut(options.OutHandle, options))
			}
		}
	}
}

func defaultJsonOptions() *JsonOptions {
	return &JsonOptions{
		ErrHandle:     os.Stderr,
		OutHandle:     os.Stdout,
		ErrLevel:      ErrorLevel,
		OutLevel:      InfoLevel,
		LogName:       "name",
		TimestampName: "time",
		LevelName:     "level",
		MessageName:   "message",
		FieldsName:    "fields",
		ErrorName:     "error",
	}
}
