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
	// bootstrap arrays
	bbt [20][]byte
	plt [10]Pair
	fbt [128]byte
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
	b.bb = b.bbt[:0]
	b.pl = b.plt[:0]
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
	bb.Write(Str2Bytes("{"))
	for i, _ := range pl {
		bb.Write([]byte{'"'})
		bb.Write(pl[i].Key)
		bb.Write([]byte{'"'})
		bb.Write([]byte{':'})
		bb.Write(pl[i].Value)
		if i < len(pl)-1 {
			bb.Write([]byte{','})
		}
	}
	bb.Write(Str2Bytes("}"))
	s := bb.String()
	bufferPool.Put(bb)
	return s
}

func (pl PairList) iString(fb []byte) []byte {
	fb = append(fb, Str2Bytes("{")...)
	for i, _ := range pl {
		fb = append(fb, Str2Bytes("\"")...)
		fb = append(fb, pl[i].Key...)
		fb = append(fb, Str2Bytes("\"")...)
		fb = append(fb, Str2Bytes(":")...)
		fb = append(fb, pl[i].Value...)
		if i < len(pl)-1 {
			fb = append(fb, Str2Bytes(",")...)
		}
	}
	fb = append(fb, Str2Bytes("}")...)
	return fb
}

func Template(template string, opts ...func(*TemplateOptions) error) func(np *Notepad) {
	options := newTemplateOptions()
	for _, fn := range opts {
		fn(options)
	}
	t, err := fasttemplate.NewTemplate(template, "${", "}")
	options.upperTags = make([]bool, len(t.Tags))
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
	if err != nil {
		panic("unexpected error when parsing template: " + err.Error())
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
				GetEntryFields(e, t.Tags, options.upperTags, &buf.bb, &buf.pl)
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

func GetEntryFields(e *Entry, tags []string, upperTags []bool, tfields *[][]byte, fields *PairList) {

	addFld := func(data []Field) {
		for _, v := range data {
			if i, ok := stringInSlice(v.Key, tags); ok {
				(*tfields)[i] = v.Buffer.Bytes()[1:]
				(*tfields)[i] = (*tfields)[i][0 : len((*tfields)[i])-1]
				if upperTags[i] {
					(*tfields)[i] = bytes.ToUpper((*tfields)[i])
				}
			} else {
				(*fields) = append((*fields), Pair{Str2Bytes(v.Key), v.Buffer.Bytes()})
			}
		}
	}
	addFld(e.Logger.Notepad.Context)
	addFld(e.Logger.Context)
	addFld(e.Data)
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
