package qlog

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
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
				writeFieldFirst(bb, Str2Bytes(topts.LogName), e.Logger.Notepad.Name)
				buf := newBuffer()
				buf.fb = strconv.AppendInt(buf.fb, e.Time.UnixNano(), 10)
				writeField(bb, Str2Bytes(topts.TimestampName), buf.fb)
				buf.free()
				writeField(bb, Str2Bytes(topts.LevelName), e.Logger.Level.ToBytes())
				writeField(bb, Str2Bytes(topts.MessageName), e.Message)
				if e.ErrorFld != nil {
					writeField(bb, Str2Bytes(topts.ErrorName), Str2Bytes(e.ErrorFld.Error()))
				}

				writeData(bb, e.Logger.Notepad.Context)
				writeData(bb, e.Logger.Context)
				writeData(bb, e.Data)

				bb.Write([]byte{'}', '\n'})
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
		LogName:       "n",
		TimestampName: "t",
		LevelName:     "l",
		MessageName:   "m",
		FieldsName:    "f",
		ErrorName:     "e",
	}
}

func writeFieldFirst(w io.Writer, name, content []byte) {
	w.Write([]byte{'"'})
	w.Write(name)
	w.Write([]byte{'"', ':'})
	w.Write([]byte{'"'})
	w.Write(content)
	w.Write([]byte{'"'})
}

func writeField(w io.Writer, name, content []byte) {
	w.Write([]byte{','})
	writeFieldFirst(w, name, content)
}

func writeData(w io.Writer, data []Field) {
	for i, _ := range data {
		w.Write([]byte{',', '"'})
		w.Write(Str2Bytes(data[i].Key))
		w.Write([]byte{'"', ':'})
		w.Write(data[i].Buffer.Bytes())
	}
}
