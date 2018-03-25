package qlog

import (
	"bytes"
	"fmt"

	"github.com/karantin2020/fasttemplate"
	// "github.com/karantin2020/qlog/buffer"
	"io"
	"os"
	"strings"
	"sync"
	"unicode"
	"unsafe"
)

type TemplateOptions struct {
	ErrHandle       io.Writer
	OutHandle       io.Writer
	ErrLevel        uint8
	OutLevel        uint8
	LogName         string
	TimestampName   string
	LevelName       string
	MessageName     string
	ErrorName       string
	FieldsName      string
	FieldsStyle     string
	FieldsSeparator byte

	upperTags map[string]bool
}

var (
	bPool *sync.Pool
)

var (
	DefaultTemplate = Template("[${name}] ${time}\t${LEVEL}\t${message}\t${fields}\n")
	ColorTemplate   = Template("[${name}] \x1b[36m${time}\x1b[0m\t\x1b[33m${LEVEL}\x1b[0m\t\x1b[32m${message}\x1b[0m\t${fields}\n")
)

type iBuffer struct {
	fb  []byte
	fbt [256]byte
}

func init() {
	bPool = &sync.Pool{
		New: func() interface{} {
			return new(iBuffer)
		},
	}
}

func newBuffer() *iBuffer {
	buf := bPool.Get().(*iBuffer)
	buf.reset()
	return buf
}

func (b *iBuffer) reset() {
	b.fb = b.fbt[:0]
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
	bb := bufferPool.Get().(*bytes.Buffer)
	bb.Reset()
	bb.Write([]byte{'{'})
	for i, _ := range pl {
		bb.Write([]byte{'"'})
		bb.Write(pl[i].Key)
		bb.Write([]byte{'"', ':'})
		bb.Write(pl[i].Value)
		if i < len(pl)-1 {
			bb.Write([]byte{','})
		}
	}
	bb.Write([]byte{'}'})
	s := bb.String()
	bufferPool.Put(bb)
	return s
}

func (pl PairList) JSONByte(fb []byte, separator byte) []byte {
	fb = append(fb, '{')
	for i, _ := range pl {
		fb = append(fb, '"')
		fb = append(fb, pl[i].Key...)
		fb = append(fb, '"')
		fb = append(fb, separator)
		fb = append(fb, pl[i].Value...)
		if i < len(pl)-1 {
			fb = append(fb, ',')
		}
	}
	fb = append(fb, '}')
	return fb
}

func Template(template string, opts ...func(*TemplateOptions) error) func(np *Notepad) {
	options := defaultTemplateOptions()
	for _, fn := range opts {
		fn(options)
	}
	t, err := fasttemplate.NewTemplate(template, "${", "}")
	emptyByteSlice := []byte{}

	// Assumption that all tags with starting upper case letter
	// have all upper case letters
	for k := range t.Tags {
		for _, r := range t.Tags[k] {
			if unicode.IsUpper(r) {
				t.Tags[k] = strings.ToLower(t.Tags[k])
				options.upperTags[t.Tags[k]] = true
			}
			break
		}
	}
	if err != nil {
		panic("unexpected error when parsing template: " + err.Error())
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
		tmpOut := func(wio io.Writer, topts *TemplateOptions) Output {
			return func(e *Entry) {
				buf := newBuffer()
				// buf.bb = buf.bb[:len(t.Tags)]
				// fmt.Printf("\nStart\n%#v\n\n", buf.pl)
				// fmt.Printf("%#v\n", e.Time)
				// fmt.Printf("%s\n", e.bufferTime)
				buf.fb = GetEntryFields(e, buf.fb, topts.FieldsSeparator)
				_, err := t.ExecuteFunc(wio, func(w io.Writer, tag string) (int, error) {
					upper := options.upperTags[tag]
					var outBytes []byte
					switch tag {
					case topts.LogName:
						outBytes = e.Logger.Notepad.Name
					case topts.TimestampName:
						outBytes = e.bufferTime
					case topts.LevelName:
						if upper {
							outBytes = e.Logger.Level.CapitalBytes()
							upper = false
						} else {
							outBytes = e.Logger.Level.ToBytes()
						}
					case topts.MessageName:
						outBytes = e.Message
					case topts.ErrorName:
						outBytes = Str2Bytes(e.ErrorFld.Error())
					case topts.FieldsName:
						outBytes = buf.fb
					default:
						outBytes = emptyByteSlice
					}
					if upper {
						return w.Write(bytes.ToUpper(outBytes))
					}
					return w.Write(outBytes)
				})
				buf.free()
				if err != nil {
					panic(fmt.Sprintf("qlog template logging error: %s", err))
				}
			}
		}
		for tlv, logger := range np.Loggers {
			level := uint8(tlv)
			switch {
			case level >= options.ErrLevel:
				(*logger).Output = append((*logger).Output, tmpOut(options.ErrHandle, options))

			case level >= options.OutLevel:
				(*logger).Output = append((*logger).Output, tmpOut(options.OutHandle, options))
			}
		}
	}
}

func defaultTemplateOptions() *TemplateOptions {
	return &TemplateOptions{
		ErrHandle:       os.Stderr,
		OutHandle:       os.Stdout,
		ErrLevel:        ErrorLevel,
		OutLevel:        InfoLevel,
		LogName:         "name",
		TimestampName:   "time",
		LevelName:       "level",
		MessageName:     "message",
		FieldsName:      "fields",
		ErrorName:       "error",
		FieldsStyle:     "json",
		FieldsSeparator: ':',
		upperTags:       make(map[string]bool),
	}
}

func GetEntryFields(e *Entry, buf []byte, sep byte) []byte {

	// fmt.Printf("%#v\n", e.Logger.Notepad.Context)
	// fmt.Printf("%#v\n", e.Logger.Context)
	// fmt.Printf("%#v\n", e.Data)
	buf = append(buf, '{')
	buf = addFld(e.Logger.Notepad.Context, buf, sep)
	buf = addFld(e.Logger.Context, buf, sep)
	buf = addFld(e.Data, buf, sep)
	buf = append(buf, '}')
	return buf
}

func addFld(data []Field, buf []byte, sep byte) []byte {
	for i, _ := range data {
		buf = append(buf, '"')
		buf = append(buf, Str2Bytes(data[i].Key)...)
		buf = append(buf, '"')
		buf = append(buf, sep)
		buf = append(buf, data[i].Buffer.Bytes()...)
		if i < len(data)-1 {
			buf = append(buf, ',')
		}
	}
	return buf
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
