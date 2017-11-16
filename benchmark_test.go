package qlog

import (
	"errors"
	"io/ioutil"
	"testing"
	// "time"
)

var (
	errExample  = errors.New("fail")
	fakeMessage = "Test logging, but use a somewhat realistic message length."
)

func BenchmarkLogNoOutput(b *testing.B) {
	log := New(InfoLevel).
		Timestamp()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.Info(fakeMessage)
		}
	})
}

func BenchmarkLogEmpty(b *testing.B) {
	log := New(InfoLevel).
		Timestamp().
		SetOutput(Template("${time}\t${LEVEL}\t${message}\t${fields}\n", func(topts *TemplateOptions) error {
			topts.ErrHandle = ioutil.Discard
			topts.OutHandle = ioutil.Discard
			return nil
		}))
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.Info("")
		}
	})
}

func BenchmarkLogDisabled(b *testing.B) {
	log := New(InfoLevel).
		Timestamp().
		SetOutput(Template("${time}\t${LEVEL}\t${message}\t${fields}\n", func(topts *TemplateOptions) error {
			topts.ErrHandle = ioutil.Discard
			topts.OutHandle = ioutil.Discard
			return nil
		}))
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.Debug(fakeMessage)
		}
	})
}

func BenchmarkInfo(b *testing.B) {
	log := New(InfoLevel).
		Timestamp().
		SetOutput(Template("${time}\t${LEVEL}\t${message}\t${fields}\n", func(topts *TemplateOptions) error {
			topts.ErrHandle = ioutil.Discard
			topts.OutHandle = ioutil.Discard
			return nil
		}))
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.Info(fakeMessage)
		}
	})
}

func BenchmarkFields(b *testing.B) {
	log := New(InfoLevel).
		Timestamp().
		SetOutput(Template("${time}\t${LEVEL}\t${message}\t${fields}\n", func(topts *TemplateOptions) error {
			topts.ErrHandle = ioutil.Discard
			topts.OutHandle = ioutil.Discard
			return nil
		}))
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.INFO.Fields(
				F{"service", "new"},
				F{"source", "after"},
			).Msg(fakeMessage)
		}
	})
}

func BenchmarkWithFields(b *testing.B) {
	log := New(InfoLevel).
		Timestamp().
		WithFields(
			F{"service", "new"},
			F{"source", "after"},
		).
		SetOutput(Template("${time}\t${LEVEL}\t${message}\t${fields}\n", func(topts *TemplateOptions) error {
			topts.ErrHandle = ioutil.Discard
			topts.OutHandle = ioutil.Discard
			return nil
		}))
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.Info(fakeMessage)
		}
	})
}
