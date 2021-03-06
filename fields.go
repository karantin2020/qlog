package qlog

import (
	"bytes"
	"fmt"
	"math"
	"strconv"
	"time"
	"unicode/utf8"
)

type F struct {
	Key   string
	Value interface{}
}

func (e *Entry) AddField(f F) {
	AddField(f, &e.Data, &e.Logger.Notepad.Options)
}

func (l *Logger) Fields(flds ...F) *Entry {
	e := l.NewEntry()
	for i, _ := range flds {
		e.AddField(flds[i])
	}
	return e
}

func AddField(f F, data *[]Field, opts *LogConfig) {
	var (
		dst   *bytes.Buffer
		found bool
	)
	for i, fld := range *data {
		if fld.Key == f.Key {
			dst = &(*data)[i].Buffer
			found = true
			break
		}
	}
	if !found {
		*data = append(*data, Field{Key: f.Key, Value: f.Value})
	}
	dst = &(*data)[len(*data)-1].Buffer
	dst.Reset()
	buf := bytesPool.Get().(*[]byte)
	*buf = (*buf)[:0]
	switch val := f.Value.(type) {
	case string:
		dst.Write(AppendString(*buf, val))
	case []byte:
		dst.Write(AppendBytes(*buf, val))
	case error:
		dst.Write(AppendError(*buf, val))
	case []error:
		dst.Write(AppendErrors(*buf, val))
	case bool:
		dst.Write(AppendBool(*buf, val))
	case int:
		dst.Write(AppendInt(*buf, val))
	case int8:
		dst.Write(AppendInt8(*buf, val))
	case int16:
		dst.Write(AppendInt16(*buf, val))
	case int32:
		dst.Write(AppendInt32(*buf, val))
	case int64:
		dst.Write(AppendInt64(*buf, val))
	case uint:
		dst.Write(AppendUint(*buf, val))
	case uint8:
		dst.Write(AppendUint8(*buf, val))
	case uint16:
		dst.Write(AppendUint16(*buf, val))
	case uint32:
		dst.Write(AppendUint32(*buf, val))
	case uint64:
		dst.Write(AppendUint64(*buf, val))
	case float32:
		dst.Write(AppendFloat32(*buf, val))
	case float64:
		dst.Write(AppendFloat64(*buf, val))
	case time.Time:
		dst.Write(AppendTime(*buf, val, opts.TimeFieldFormat))
	case time.Duration:
		dst.Write(AppendDuration(*buf, val,
			opts.DurationFieldUnit, opts.DurationFieldInteger))
	case []string:
		dst.Write(AppendStrings(*buf, val))
	case []bool:
		dst.Write(AppendBools(*buf, val))
	case []int:
		dst.Write(AppendInts(*buf, val))
	case []int8:
		dst.Write(AppendInts8(*buf, val))
	case []int16:
		dst.Write(AppendInts16(*buf, val))
	case []int32:
		dst.Write(AppendInts32(*buf, val))
	case []int64:
		dst.Write(AppendInts64(*buf, val))
	case []uint:
		dst.Write(AppendUints(*buf, val))
	// case []uint8:
	// 	dst.Write(AppendUints8(*buf, val))
	case []uint16:
		dst.Write(AppendUints16(*buf, val))
	case []uint32:
		dst.Write(AppendUints32(*buf, val))
	case []uint64:
		dst.Write(AppendUints64(*buf, val))
	case []float32:
		dst.Write(AppendFloats32(*buf, val))
	case []float64:
		dst.Write(AppendFloats64(*buf, val))
	case []time.Time:
		dst.Write(AppendTimes(*buf, val, opts.TimeFieldFormat))
	case []time.Duration:
		dst.Write(AppendDurations(*buf, val,
			opts.DurationFieldUnit, opts.DurationFieldInteger))
	case nil:
		dst.Write(append(*buf, "null"...))
	default:
		dst.Write(AppendInterface(*buf, val, opts.InterfaceMarshaler))
	}
	// if !found {
	// 	*buffer = append(*buffer, dst)
	// }
	bytesPool.Put(buf)
}

func AppendBool(dst []byte, val bool) []byte {
	return strconv.AppendBool(dst, val)
}

func AppendBools(dst []byte, vals []bool) []byte {
	if len(vals) == 0 {
		return append(dst, '[', ']')
	}
	dst = append(dst, '[')
	dst = strconv.AppendBool(dst, vals[0])
	if len(vals) > 1 {
		for _, val := range vals[1:] {
			dst = strconv.AppendBool(append(dst, ','), val)
		}
	}
	dst = append(dst, ']')
	return dst
}

func AppendInt(dst []byte, val int) []byte {
	return strconv.AppendInt(dst, int64(val), 10)
}

func AppendInts(dst []byte, vals []int) []byte {
	if len(vals) == 0 {
		return append(dst, '[', ']')
	}
	dst = append(dst, '[')
	dst = strconv.AppendInt(dst, int64(vals[0]), 10)
	if len(vals) > 1 {
		for _, val := range vals[1:] {
			dst = strconv.AppendInt(append(dst, ','), int64(val), 10)
		}
	}
	dst = append(dst, ']')
	return dst
}

func AppendInt8(dst []byte, val int8) []byte {
	return strconv.AppendInt(dst, int64(val), 10)
}

func AppendInts8(dst []byte, vals []int8) []byte {
	if len(vals) == 0 {
		return append(dst, '[', ']')
	}
	dst = append(dst, '[')
	dst = strconv.AppendInt(dst, int64(vals[0]), 10)
	if len(vals) > 1 {
		for _, val := range vals[1:] {
			dst = strconv.AppendInt(append(dst, ','), int64(val), 10)
		}
	}
	dst = append(dst, ']')
	return dst
}

func AppendInt16(dst []byte, val int16) []byte {
	return strconv.AppendInt(dst, int64(val), 10)
}

func AppendInts16(dst []byte, vals []int16) []byte {
	if len(vals) == 0 {
		return append(dst, '[', ']')
	}
	dst = append(dst, '[')
	dst = strconv.AppendInt(dst, int64(vals[0]), 10)
	if len(vals) > 1 {
		for _, val := range vals[1:] {
			dst = strconv.AppendInt(append(dst, ','), int64(val), 10)
		}
	}
	dst = append(dst, ']')
	return dst
}

func AppendInt32(dst []byte, val int32) []byte {
	return strconv.AppendInt(dst, int64(val), 10)
}

func AppendInts32(dst []byte, vals []int32) []byte {
	if len(vals) == 0 {
		return append(dst, '[', ']')
	}
	dst = append(dst, '[')
	dst = strconv.AppendInt(dst, int64(vals[0]), 10)
	if len(vals) > 1 {
		for _, val := range vals[1:] {
			dst = strconv.AppendInt(append(dst, ','), int64(val), 10)
		}
	}
	dst = append(dst, ']')
	return dst
}

func AppendInt64(dst []byte, val int64) []byte {
	return strconv.AppendInt(dst, val, 10)
}

func AppendInts64(dst []byte, vals []int64) []byte {
	if len(vals) == 0 {
		return append(dst, '[', ']')
	}
	dst = append(dst, '[')
	dst = strconv.AppendInt(dst, vals[0], 10)
	if len(vals) > 1 {
		for _, val := range vals[1:] {
			dst = strconv.AppendInt(append(dst, ','), val, 10)
		}
	}
	dst = append(dst, ']')
	return dst
}

func AppendUint(dst []byte, val uint) []byte {
	return strconv.AppendUint(dst, uint64(val), 10)
}

func AppendUints(dst []byte, vals []uint) []byte {
	if len(vals) == 0 {
		return append(dst, '[', ']')
	}
	dst = append(dst, '[')
	dst = strconv.AppendUint(dst, uint64(vals[0]), 10)
	if len(vals) > 1 {
		for _, val := range vals[1:] {
			dst = strconv.AppendUint(append(dst, ','), uint64(val), 10)
		}
	}
	dst = append(dst, ']')
	return dst
}

func AppendUint8(dst []byte, val uint8) []byte {
	return strconv.AppendUint(dst, uint64(val), 10)
}

func AppendUints8(dst []byte, vals []uint8) []byte {
	if len(vals) == 0 {
		return append(dst, '[', ']')
	}
	dst = append(dst, '[')
	dst = strconv.AppendUint(dst, uint64(vals[0]), 10)
	if len(vals) > 1 {
		for _, val := range vals[1:] {
			dst = strconv.AppendUint(append(dst, ','), uint64(val), 10)
		}
	}
	dst = append(dst, ']')
	return dst
}

func AppendUint16(dst []byte, val uint16) []byte {
	return strconv.AppendUint(dst, uint64(val), 10)
}

func AppendUints16(dst []byte, vals []uint16) []byte {
	if len(vals) == 0 {
		return append(dst, '[', ']')
	}
	dst = append(dst, '[')
	dst = strconv.AppendUint(dst, uint64(vals[0]), 10)
	if len(vals) > 1 {
		for _, val := range vals[1:] {
			dst = strconv.AppendUint(append(dst, ','), uint64(val), 10)
		}
	}
	dst = append(dst, ']')
	return dst
}

func AppendUint32(dst []byte, val uint32) []byte {
	return strconv.AppendUint(dst, uint64(val), 10)
}

func AppendUints32(dst []byte, vals []uint32) []byte {
	if len(vals) == 0 {
		return append(dst, '[', ']')
	}
	dst = append(dst, '[')
	dst = strconv.AppendUint(dst, uint64(vals[0]), 10)
	if len(vals) > 1 {
		for _, val := range vals[1:] {
			dst = strconv.AppendUint(append(dst, ','), uint64(val), 10)
		}
	}
	dst = append(dst, ']')
	return dst
}

func AppendUint64(dst []byte, val uint64) []byte {
	return strconv.AppendUint(dst, uint64(val), 10)
}

func AppendUints64(dst []byte, vals []uint64) []byte {
	if len(vals) == 0 {
		return append(dst, '[', ']')
	}
	dst = append(dst, '[')
	dst = strconv.AppendUint(dst, vals[0], 10)
	if len(vals) > 1 {
		for _, val := range vals[1:] {
			dst = strconv.AppendUint(append(dst, ','), val, 10)
		}
	}
	dst = append(dst, ']')
	return dst
}

func AppendFloat(dst []byte, val float64, bitSize int) []byte {
	// JSON does not permit NaN or Infinity. A typical JSON encoder would fail
	// with an error, but a logging library wants the data to get thru so we
	// make a tradeoff and store those types as string.
	switch {
	case math.IsNaN(val):
		return append(dst, `"NaN"`...)
	case math.IsInf(val, 1):
		return append(dst, `"+Inf"`...)
	case math.IsInf(val, -1):
		return append(dst, `"-Inf"`...)
	}
	return strconv.AppendFloat(dst, val, 'f', -1, bitSize)
}

func AppendFloat32(dst []byte, val float32) []byte {
	return AppendFloat(dst, float64(val), 32)
}

func AppendFloats32(dst []byte, vals []float32) []byte {
	if len(vals) == 0 {
		return append(dst, '[', ']')
	}
	dst = append(dst, '[')
	dst = AppendFloat(dst, float64(vals[0]), 32)
	if len(vals) > 1 {
		for _, val := range vals[1:] {
			dst = AppendFloat(append(dst, ','), float64(val), 32)
		}
	}
	dst = append(dst, ']')
	return dst
}

func AppendFloat64(dst []byte, val float64) []byte {
	return AppendFloat(dst, val, 64)
}

func AppendFloats64(dst []byte, vals []float64) []byte {
	if len(vals) == 0 {
		return append(dst, '[', ']')
	}
	dst = append(dst, '[')
	dst = AppendFloat(dst, vals[0], 32)
	if len(vals) > 1 {
		for _, val := range vals[1:] {
			dst = AppendFloat(append(dst, ','), val, 64)
		}
	}
	dst = append(dst, ']')
	return dst
}

func AppendInterface(dst []byte, i interface{}, fn func(v interface{}) ([]byte, error)) []byte {
	marshaled, err := fn(i)
	if err != nil {
		return AppendString(dst, fmt.Sprintf("marshaling error: %v", err))
	}
	return append(dst, marshaled...)
}

const hex = "0123456789abcdef"

func AppendStrings(dst []byte, vals []string) []byte {
	if len(vals) == 0 {
		return append(dst, '[', ']')
	}
	dst = append(dst, '[')
	dst = AppendString(dst, vals[0])
	if len(vals) > 1 {
		for _, val := range vals[1:] {
			dst = AppendString(append(dst, ','), val)
		}
	}
	dst = append(dst, ']')
	return dst
}

// AppendString encodes the input string to json and appends
// the encoded string to the input byte slice.
//
// The operation loops though each byte in the string looking
// for characters that need json or utf8 encoding. If the string
// does not need encoding, then the string is appended in it's
// entirety to the byte slice.
// If we encounter a byte that does need encoding, switch up
// the operation and perform a byte-by-byte read-encode-append.
func AppendString(dst []byte, s string) []byte {
	// Start with a double quote.
	dst = append(dst, '"')
	dst = AppendStringNoQuotes(dst, s)
	// End with a double quote
	return append(dst, '"')
}

func AppendStringNoQuotes(dst []byte, s string) []byte {
	// Loop through each character in the string.
	for i := 0; i < len(s); i++ {
		// Check if the character needs encoding. Control characters, slashes,
		// and the double quote need json encoding. Bytes above the ascii
		// boundary needs utf8 encoding.
		if s[i] < 0x20 || s[i] > 0x7e || s[i] == '\\' || s[i] == '"' {
			// We encountered a character that needs to be encoded. Switch
			// to complex version of the algorithm.
			dst = appendStringComplex(dst, s, i)
			return append(dst, '"')
		}
	}
	// The string has no need for encoding an therefore is directly
	// appended to the byte slice.
	return append(dst, s...)
}

// appendStringComplex is used by appendString to take over an in
// progress JSON string encoding that encountered a character that needs
// to be encoded.
func appendStringComplex(dst []byte, s string, i int) []byte {
	start := 0
	for i < len(s) {
		b := s[i]
		if b >= utf8.RuneSelf {
			r, size := utf8.DecodeRuneInString(s[i:])
			if r == utf8.RuneError && size == 1 {
				// In case of error, first append previous simple characters to
				// the byte slice if any and append a remplacement character code
				// in place of the invalid sequence.
				if start < i {
					dst = append(dst, s[start:i]...)
				}
				dst = append(dst, `\ufffd`...)
				i += size
				start = i
				continue
			}
			i += size
			continue
		}
		if b >= 0x20 && b <= 0x7e && b != '\\' && b != '"' {
			i++
			continue
		}
		// We encountered a character that needs to be encoded.
		// Let's append the previous simple characters to the byte slice
		// and switch our operation to read and encode the remainder
		// characters byte-by-byte.
		if start < i {
			dst = append(dst, s[start:i]...)
		}
		switch b {
		case '"', '\\':
			dst = append(dst, '\\', b)
		case '\b':
			dst = append(dst, '\\', 'b')
		case '\f':
			dst = append(dst, '\\', 'f')
		case '\n':
			dst = append(dst, '\\', 'n')
		case '\r':
			dst = append(dst, '\\', 'r')
		case '\t':
			dst = append(dst, '\\', 't')
		default:
			dst = append(dst, '\\', 'u', '0', '0', hex[b>>4], hex[b&0xF])
		}
		i++
		start = i
	}
	if start < len(s) {
		dst = append(dst, s[start:]...)
	}
	return dst
}

// AppendBytes is a mirror of appendString with []byte arg
func AppendBytes(dst, s []byte) []byte {
	dst = append(dst, '"')
	for i := 0; i < len(s); i++ {
		if s[i] < 0x20 || s[i] > 0x7e || s[i] == '\\' || s[i] == '"' {
			dst = appendBytesComplex(dst, s, i)
			return append(dst, '"')
		}
	}
	dst = append(dst, s...)
	return append(dst, '"')
}

// appendBytesComplex is a mirror of the appendStringComplex
// with []byte arg
func appendBytesComplex(dst, s []byte, i int) []byte {
	start := 0
	for i < len(s) {
		b := s[i]
		if b >= utf8.RuneSelf {
			r, size := utf8.DecodeRune(s[i:])
			if r == utf8.RuneError && size == 1 {
				if start < i {
					dst = append(dst, s[start:i]...)
				}
				dst = append(dst, `\ufffd`...)
				i += size
				start = i
				continue
			}
			i += size
			continue
		}
		if b >= 0x20 && b <= 0x7e && b != '\\' && b != '"' {
			i++
			continue
		}
		// We encountered a character that needs to be encoded.
		// Let's append the previous simple characters to the byte slice
		// and switch our operation to read and encode the remainder
		// characters byte-by-byte.
		if start < i {
			dst = append(dst, s[start:i]...)
		}
		switch b {
		case '"', '\\':
			dst = append(dst, '\\', b)
		case '\b':
			dst = append(dst, '\\', 'b')
		case '\f':
			dst = append(dst, '\\', 'f')
		case '\n':
			dst = append(dst, '\\', 'n')
		case '\r':
			dst = append(dst, '\\', 'r')
		case '\t':
			dst = append(dst, '\\', 't')
		default:
			dst = append(dst, '\\', 'u', '0', '0', hex[b>>4], hex[b&0xF])
		}
		i++
		start = i
	}
	if start < len(s) {
		dst = append(dst, s[start:]...)
	}
	return dst
}

func AppendError(dst []byte, err error) []byte {
	if err == nil {
		return append(dst, `null`...)
	}
	return AppendString(dst, err.Error())
}

func AppendErrors(dst []byte, errs []error) []byte {
	if len(errs) == 0 {
		return append(dst, '[', ']')
	}
	dst = append(dst, '[')
	if errs[0] != nil {
		dst = AppendString(dst, errs[0].Error())
	} else {
		dst = append(dst, "null"...)
	}
	if len(errs) > 1 {
		for _, err := range errs[1:] {
			if err == nil {
				dst = append(dst, ",null"...)
				continue
			}
			dst = AppendString(append(dst, ','), err.Error())
		}
	}
	dst = append(dst, ']')
	return dst
}

func AppendTime(dst []byte, t time.Time, format string) []byte {
	if format == "" {
		return AppendInt64(dst, t.Unix())
	}
	return append(t.AppendFormat(append(dst, '"'), format), '"')
}

func AppendTimes(dst []byte, vals []time.Time, format string) []byte {
	if format == "" {
		return appendUnixTimes(dst, vals)
	}
	if len(vals) == 0 {
		return append(dst, '[', ']')
	}
	dst = append(dst, '[')
	dst = append(vals[0].AppendFormat(append(dst, '"'), format), '"')
	if len(vals) > 1 {
		for _, t := range vals[1:] {
			dst = append(t.AppendFormat(append(dst, ',', '"'), format), '"')
		}
	}
	dst = append(dst, ']')
	return dst
}

func appendUnixTimes(dst []byte, vals []time.Time) []byte {
	if len(vals) == 0 {
		return append(dst, '[', ']')
	}
	dst = append(dst, '[')
	dst = strconv.AppendInt(dst, vals[0].Unix(), 10)
	if len(vals) > 1 {
		for _, t := range vals[1:] {
			dst = strconv.AppendInt(dst, t.Unix(), 10)
		}
	}
	dst = append(dst, ']')
	return dst
}

func AppendDuration(dst []byte, d time.Duration, unit time.Duration, useInt bool) []byte {
	if useInt {
		return strconv.AppendInt(dst, int64(d/unit), 10)
	}
	return AppendFloat64(dst, float64(d)/float64(unit))
}

func AppendDurations(dst []byte, vals []time.Duration, unit time.Duration, useInt bool) []byte {
	if len(vals) == 0 {
		return append(dst, '[', ']')
	}
	dst = append(dst, '[')
	dst = AppendDuration(dst, vals[0], unit, useInt)
	if len(vals) > 1 {
		for _, d := range vals[1:] {
			dst = AppendDuration(append(dst, ','), d, unit, useInt)
		}
	}
	dst = append(dst, ']')
	return dst
}
