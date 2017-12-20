package qlog

import (
	// "encoding/json"
	"fmt"
	// "github.com/karantin2020/qlog/buffer"
	"bytes"
	"os"
	"sync"
	"time"
)

var (
	entryPool  *sync.Pool
	bufferPool *sync.Pool
	bytesPool  *sync.Pool
)

const (
	bufferFields = 7
)

func init() {
	entryPool = &sync.Pool{
		New: func() interface{} {
			return new(Entry)
		},
	}
	bytesPool = &sync.Pool{
		New: func() interface{} {
			b := make([]byte, 0, 64)
			return &b
		},
	}
	bufferPool = &sync.Pool{
		New: func() interface{} {
			return new(bytes.Buffer)
		},
	}
}

type Entry struct {
	Logger  *Logger
	Data    []Field
	Time    time.Time
	Level   Level
	Message string
	ErrorF  error
	buffer  [bufferFields]Field
}

// A Field is a marshaling struct type used to add a key-value pair to a logger's
// context
type Field struct {
	Key    string
	Value  interface{}
	Buffer bytes.Buffer
}

func (f Field) MarshalJSON() ([]byte, error) {
	b := bufferPool.Get().(*bytes.Buffer)
	b.Reset()
	defer bufferPool.Put(b)
	if _, err := b.Write(Str2Bytes(f.Key)); err != nil {
		return nil, err
	}
	if _, err := b.Write([]byte{':'}); err != nil {
		return nil, err
	}
	if _, err := b.Write(f.Buffer.Bytes()); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func (l *Logger) NewEntry() *Entry {
	entry, _ := entryPool.Get().(*Entry)
	entry.Reset()
	entry.Time = time.Now()
	entry.Logger = l
	entry.Level = l.Level
	return entry
}

func (e *Entry) Reset() {
	e.Logger = nil
	// e.Time = time.Time{}
	// e.Level = DebugLevel
	e.Message = ""
	for i := range e.buffer {
		e.buffer[i].Buffer.Reset()
	}
	e.Data = e.buffer[:0]
	// e.Buffer = e.Buffer[:0]
	e.ErrorF = nil
}

func (e *Entry) Fields(fields ...F) {
	for _, fld := range fields {
		e.AddField(fld)
	}
}

func (e *Entry) Err(err error) {
	e.AddField(F{Key: e.Logger.Notepad.Options.ErrorFieldName, Value: err})
}

func (e *Entry) Timestamp() {
	e.AddField(F{Key: e.Logger.Notepad.Options.TimestampFieldName,
		Value: e.Logger.Notepad.Options.TimestampFunc()})
}

// func (e *Entry) Caller(calldepth int) {
// 	_, file, line, ok := runtime.Caller(calldepth)
// 	if !ok {
// 		file = "???"
// 		line = 0
// 	}
// 	e.AddField(F{e.Logger.Notepad.Options.CallerFieldName,
// 		fmt.Sprintf("%s:%d", file, line)})
// }

func (e *Entry) Debug(msg string) {
	if e.Logger == nil || e.Logger != e.Logger.Notepad.DEBUG {
		return
	}
	e.Msg(msg)
}

func (e *Entry) Info(msg string) {
	if e.Logger == nil || e.Logger != e.Logger.Notepad.INFO {
		return
	}
	e.Msg(msg)
}

func (e *Entry) Warn(msg string) {
	if e.Logger == nil || e.Logger != e.Logger.Notepad.WARN {
		return
	}
	e.Msg(msg)
}

func (e *Entry) errMsg(msg string, panicErr, exitErr bool) {
	e.Level = e.Logger.Level
	// if e.ErrorF == nil {
	e.ErrorF = e.Logger.Notepad.Options.ErrorFunc(msg)
	e.AddField(F{Key: e.Logger.Notepad.Options.ErrorFieldName, Value: e.ErrorF})
	// } else {
	e.Message = e.ErrorF.Error()
	// }
	e.Process()
	if panicErr {
		panic(msg)
	} else if exitErr {
		os.Exit(1)
	}

}

func (e *Entry) Error(msg string) {
	if e.Logger == nil || e.Logger != e.Logger.Notepad.ERROR {
		return
	}
	e.errMsg(msg, false, false)
}

func (e *Entry) Critical(msg string) {
	if e.Logger == nil || e.Logger != e.Logger.Notepad.CRITICAL {
		return
	}
	// e.errMsg(msg)
	if e.Logger.Notepad.DEBUG != nil {
		e.errMsg(msg, true, false)
	} else {
		e.errMsg(msg, false, false)
	}
}

func (e *Entry) Panic(msg string) {
	if e.Logger == nil || e.Logger != e.Logger.Notepad.PANIC {
		return
	}
	// e.errMsg(msg)
	e.errMsg(msg, true, false)
}

func (e *Entry) Fatal(msg string) {
	if e.Logger == nil || e.Logger != e.Logger.Notepad.FATAL {
		return
	}
	e.errMsg(msg, false, true)
	// os.Exit(1)
}

func (e *Entry) Log(msg string) {
	if e.Logger == nil || e.Logger != e.Logger.Notepad.LOG {
		return
	}
	e.Msg(msg)
}

func (e *Entry) Msg(msg string) {
	if e.Logger == nil {
		return
	}
	e.Level = e.Logger.Level
	e.Message = msg
	e.Process()
}

func (e *Entry) Msgf(format string, a ...interface{}) {
	if e.Logger == nil {
		return
	}
	e.Msg(fmt.Sprintf(format, a...))
}

func (e *Entry) Process() {
	for _, frmt := range e.Logger.Notepad.Formatter {
		frmt(e)
	}
	for _, hook := range e.Logger.Hooks {
		hook(e)
	}
	if len(e.Logger.Output) == 1 {
		e.Logger.Output[0](e)
	} else if len(e.Logger.Output) > 1 {
		var wg sync.WaitGroup
		for _, out := range e.Logger.Output {
			wg.Add(1)
			go func(et *Entry, fn func(*Entry)) {
				// Decrement the counter when the goroutine completes.
				defer wg.Done()
				fn(et)
			}(e, out)
		}
		wg.Wait()
	}
	// fmt.Printf("%v\t%v\t%s\n", e.Data, e.ErrorF, e.Message)
	// fmt.Printf("%v\n", e.Logger.Context)
	// fmt.Printf("%v\n", e.Logger.Notepad.Context)
	// for _, buf := range e.Buffer {
	// 	buf.Free()
	// }
	entryPool.Put(e)
}
