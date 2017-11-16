package qlog

import (
	"encoding/json"
	"fmt"
	"github.com/karantin2020/qlog/buffer"
	// "runtime"
	"sync"
	"time"
)

var (
	entryPool  *sync.Pool
	bufferPool *buffer.Pool
)

func init() {
	entryPool = &sync.Pool{
		New: func() interface{} {
			return &Entry{
				Data:   make([]Field, 0, 10),
				Buffer: make([]*buffer.Buffer, 0, 10),
			}
		},
	}
	bufferPool = buffer.NewPool()
}

type Entry struct {
	Logger  *Logger
	Data    []Field
	Time    time.Time
	Level   Level
	Message string
	ErrorF  error
	Buffer  []*buffer.Buffer
}

// A Field is a marshaling struct type used to add a key-value pair to a logger's
// context
type Field struct {
	Key   string
	Value interface{}
}

func (f Field) MarshalJSON() ([]byte, error) {
	b1, err := json.Marshal(f.Value)
	b2 := []byte(f.Key + ":")
	b2 = append(b2, b1...)
	return b2, err
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
	e.Data = e.Data[:0]
	e.Buffer = e.Buffer[:0]
}

func (e *Entry) Fields(fields ...Field) {
	for _, fld := range fields {
		e.AddField(fld)
	}
}

func (e *Entry) Err(err error) {
	e.AddField(F{e.Logger.Notepad.Options.ErrorFieldName, err})
}

func (e *Entry) Timestamp() {
	e.AddField(F{e.Logger.Notepad.Options.TimestampFieldName,
		e.Logger.Notepad.Options.TimestampFunc()})
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

func (e *Entry) errMsg(msg string) {
	e.Level = e.Logger.Level
	if e.ErrorF == nil {
		e.ErrorF = e.Logger.Notepad.Options.ErrorFunc(msg)
	} else {
		e.Message = msg
	}
	e.Process()
}

func (e *Entry) Error(msg string) {
	if e.Logger == nil || e.Logger != e.Logger.Notepad.ERROR {
		return
	}
	e.errMsg(msg)
}

func (e *Entry) Critical(msg string) {
	if e.Logger == nil || e.Logger != e.Logger.Notepad.CRITICAL {
		return
	}
	e.errMsg(msg)
}

func (e *Entry) Panic(msg string) {
	if e.Logger == nil || e.Logger != e.Logger.Notepad.PANIC {
		return
	}
	e.errMsg(msg)
}

func (e *Entry) Fatal(msg string) {
	if e.Logger == nil || e.Logger != e.Logger.Notepad.FATAL {
		return
	}
	e.errMsg(msg)
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
	} else {
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
	for _, buf := range e.Buffer {
		buf.Free()
	}
	entryPool.Put(e)
}
