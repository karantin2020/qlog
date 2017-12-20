package qlog

import (
	"errors"
	"io/ioutil"
	"testing"
)

var (
	errExample  = errors.New("fail")
	errMsg      = "fail"
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

func BenchmarkError(b *testing.B) {
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
			log.Error(errMsg)
		}
	})
}

func BenchmarkInfoLower(b *testing.B) {
	log := New(InfoLevel).
		Timestamp().
		SetOutput(Template("${time}\t${level}\t${message}\t${fields}\n", func(topts *TemplateOptions) error {
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

func BenchmarkOneField(b *testing.B) {
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
				F{Key: "service", Value: "new"},
			).Msg(fakeMessage)
		}
	})
}

func BenchmarkTwoFields(b *testing.B) {
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
				F{Key: "service", Value: "new"},
				F{Key: "source", Value: "after"},
			).Msg(fakeMessage)
		}
	})
}

func BenchmarkOneFieldLower(b *testing.B) {
	log := New(InfoLevel).
		Timestamp().
		SetOutput(Template("${time}\t${level}\t${message}\t${fields}\n", func(topts *TemplateOptions) error {
			topts.ErrHandle = ioutil.Discard
			topts.OutHandle = ioutil.Discard
			return nil
		}))
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.INFO.Fields(
				F{Key: "service", Value: "new"},
			).Msg(fakeMessage)
		}
	})
}

func BenchmarkTwoFieldsLower(b *testing.B) {
	log := New(InfoLevel).
		Timestamp().
		SetOutput(Template("${time}\t${level}\t${message}\t${fields}\n", func(topts *TemplateOptions) error {
			topts.ErrHandle = ioutil.Discard
			topts.OutHandle = ioutil.Discard
			return nil
		}))
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.INFO.Fields(
				F{Key: "service", Value: "new"},
				F{Key: "source", Value: "after"},
			).Msg(fakeMessage)
		}
	})
}

func BenchmarkWithFields(b *testing.B) {
	log := New(InfoLevel).
		Timestamp().
		WithFields(
			F{Key: "service", Value: "new"},
			F{Key: "source", Value: "after"},
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
