package qlog

import (
	"bytes"
	"fmt"
	"github.com/karantin2020/fasttemplate"
	"github.com/karantin2020/qlog/buffer"
	"io"
	"os"
	"strings"
	"sync"
	"unicode"
	"unsafe"
)

type TemplateOptions struct {
	ErrHandle     io.Writer
	OutHandle     io.Writer
	ErrLevel      Level
	OutLevel      Level
	FieldsName    string
	MessageName   string
	ErrorName     string
	TimestampName string

	upperTags []bool
}

var (
	bPool *sync.Pool
)

type iBuffer struct {
	bb [][]byte
	pl PairList
	fb []byte
}

func init() {
	bPool = &sync.Pool{
		New: func() interface{} {
			return &iBuffer{
				bb: make([][]byte, 0, 20),
				pl: make(PairList, 0, 10),
				fb: make([]byte, 0, 512),
			}
		},
	}
}

func newBuffer() *iBuffer {
	buf := bPool.Get().(*iBuffer)
	buf.reset()
	return buf
}

func (b *iBuffer) reset() {
	for i := range b.bb {
		b.bb[i] = nil
	}
	b.pl = b.pl[:0]
	b.fb = b.fb[:0]
}

func (b *iBuffer) free() {
	bPool.Put(b)
}

type Pair struct {
	Key   []byte
	Value []byte
}

type PairList []Pair

func (pl PairList) String() string {
	bb := bufferPool.Get()
	bb.Write(Str2Bytes("{"))
	for i, _ := range pl {
		bb.Write(Str2Bytes("\""))
		bb.Write(pl[i].Key)
		bb.Write(Str2Bytes("\""))
		bb.Write(Str2Bytes(":"))
		bb.Write(pl[i].Value)
		if i < len(pl)-1 {
			bb.Write(Str2Bytes(","))
		}
	}
	bb.Write(Str2Bytes("}"))
	s := bb.String()
	bb.Free()
	return s
}

func (pl PairList) iString(fb []byte) []byte {
	// write "{"
	fb = append(fb, 123)
	for i, _ := range pl {
		// write "\""
		fb = append(fb, 34)
		fb = append(fb, pl[i].Key...)
		// write "\""
		fb = append(fb, 34)
		// write ":"
		fb = append(fb, 58)
		fb = append(fb, pl[i].Value...)
		if i < len(pl)-1 {
			// write ","
			fb = append(fb, 44)
		}
	}
	// write "}"
	fb = append(fb, 125)
	return fb
}

type cacheBuffer struct {
	tbuf [][]byte
	fbuf PairList
}

type cacheList struct {
	loggers []*Logger
	buffer  []*cacheBuffer
}

// Assume that every logger and it's parent notepad have immutable fields lists.
// Their fields lists can be changed only by creating new instances of
// logger and it's parent notepad
func (cl *cacheList) fromCache(l *Logger) (*cacheBuffer, bool) {
	for i, ll := range cl.loggers {
		if ll == l {
			return cl.buffer[i], true
		}
	}
	return nil, false
}

func (cl *cacheList) toCache(l *Logger, ib *iBuffer) {
	cl.loggers = append(cl.loggers, l)
	cl.buffer = append(cl.buffer, &cacheBuffer{})
	cl.buffer[len(cl.buffer)-1].copyFrom(ib)
}

func (cb *cacheBuffer) copyTo(ib *iBuffer) {
	copy(ib.bb, cb.tbuf)
	ib.pl = append(ib.pl, cb.fbuf...)
}

func (cb *cacheBuffer) copyFrom(ib *iBuffer) {
	cb.tbuf = append(cb.tbuf, ib.bb...)
	cb.fbuf = append(cb.fbuf, ib.pl...)
}

func Template(template string, opts ...func(*TemplateOptions) error) func(np *Notepad) {
	options := newTemplateOptions()
	for _, fn := range opts {
		fn(options)
	}
	t, err := fasttemplate.NewTemplate(template, "${", "}")
	options.upperTags = make([]bool, len(t.Tags))
	if err != nil {
		panic("unexpected error when parsing template: " + err.Error())
	}
	// Assumption that all tags with starting upper case letter
	// have all upper case letters
	for k := range t.Tags {
		for _, r := range t.Tags[k] {
			if unicode.IsUpper(r) {
				options.upperTags[k] = true
				t.Tags[k] = strings.ToLower(t.Tags[k])
			}
			break
		}
	}
	cache := &cacheList{
		loggers: []*Logger{},
		buffer:  []*cacheBuffer{},
	}
	return func(np *Notepad) {
		if options.OutLevel > _maxLevel || options.OutLevel < _minLevel || options.OutLevel < np.Level {
			panic("OutLevel is out of range")
		}
		if options.ErrLevel > _maxLevel || options.ErrLevel < _minLevel || options.ErrLevel < np.Level {
			panic("ErrLevel is out of range")
		}
		if options.OutLevel > options.ErrLevel {
			panic("OutLevel is higher than errLevel")
		}
		tmpOut := func(wio io.Writer) Output {
			return func(e *Entry) {
				buf := newBuffer()
				buf.bb = buf.bb[:len(t.Tags)]
				if c, ok := cache.fromCache(e.Logger); ok {
					c.copyTo(buf)
					GetEntryFields(e, t.Tags, options.upperTags, &buf.bb, &buf.pl, true)
				} else {
					GetEntryFields(e, t.Tags, options.upperTags, &buf.bb, &buf.pl, false)
					cache.toCache(e.Logger, buf)
					GetEntryFields(e, t.Tags, options.upperTags, &buf.bb, &buf.pl, true)
				}
				for i, t := range t.Tags {
					if t == options.FieldsName {
						buf.bb[i] = buf.pl.iString(buf.fb)
					}
					if t == options.MessageName {
						buf.bb[i] = Str2Bytes(e.Message)
					}
				}
				i := 0
				_, err := t.ExecuteFunc(wio, func(w io.Writer, tag string) (int, error) {
					i += 1
					return w.Write(buf.bb[i-1])
				})
				if err != nil {
					panic(fmt.Sprintf("unexpected error: %s", err))
				}
				buf.free()
			}
		}
		for tlv, logger := range np.Loggers {
			level := Level(tlv)
			switch {
			case level >= options.ErrLevel:
				(*logger).Output = append((*logger).Output, tmpOut(options.ErrHandle))

			case level >= options.OutLevel:
				(*logger).Output = append((*logger).Output, tmpOut(options.OutHandle))
			}
		}
	}
}

func newTemplateOptions() *TemplateOptions {
	return &TemplateOptions{
		ErrHandle:   os.Stderr,
		OutHandle:   os.Stdout,
		ErrLevel:    ErrorLevel,
		OutLevel:    InfoLevel,
		FieldsName:  "fields",
		MessageName: "message",
	}
}

func GetEntryFields(e *Entry, tags []string, upperTags []bool, tfields *[][]byte, fields *PairList, cached bool) {

	addFld := func(data []Field, encValues []*buffer.Buffer) {
		for j, v := range data {
			if i, ok := stringInSlice(v.Key, tags); ok {
				if upperTags[i] {
					(*tfields)[i] = bytes.ToUpper(encValues[j].Bytes()[1:])
					(*tfields)[i] = (*tfields)[i][0 : len((*tfields)[i])-1]
				} else {
					(*tfields)[i] = encValues[j].Bytes()[1:]
					(*tfields)[i] = (*tfields)[i][0 : len((*tfields)[i])-1]
				}
			} else {
				(*fields) = append((*fields), Pair{Str2Bytes(v.Key), encValues[j].Bytes()})
			}
		}
	}
	if !cached {
		addFld(e.Logger.Notepad.Context, e.Logger.Notepad.CtxBuffer)
		addFld(e.Logger.Context, e.Logger.CtxBuffer)
	} else {
		addFld(e.Data, e.Buffer)
	}
}

func stringInSlice(a string, list []string) (int, bool) {
	for i, b := range list {
		if b == a {
			return i, true
		}
	}
	return -1, false
}

func Str2Bytes(stringData string) []byte {
	temporaryData := (*[2]uintptr)(unsafe.Pointer(&stringData))
	bytesData := [3]uintptr{temporaryData[0], temporaryData[1], temporaryData[1]}
	return *(*[]byte)(unsafe.Pointer(&bytesData))
}

func Bytes2Str(bytesData []byte) string {
	return *(*string)(unsafe.Pointer(&bytesData))
}
