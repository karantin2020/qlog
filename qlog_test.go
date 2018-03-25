package qlog_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/karantin2020/qlog"
)

func TestNew(t *testing.T) {
	type args struct {
		lvl uint8
	}
	tests := []struct {
		name string
		args args
		want *qlog.Notepad
	}{
		{
			"Higher lvl new Notepad",
			args{100},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				p := recover()
				if p == nil {
					return
				}
				if tt.want != nil {
					t.Fatalf("Want New Notepad, got panic: %#v", p)
				}
			}()
			if got := qlog.New("", tt.args.lvl); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNotepad_AddHook(t *testing.T) {
	np := qlog.New("", qlog.InfoLevel)
	hook := qlog.Hook(func(e *qlog.Entry) {})
	tests := []struct {
		name      string
		fn        func() error
		wantError bool
	}{
		{
			"New Hook add",
			func() error {
				np.AddHook(qlog.WarnLevel, hook)
				// fmt.Printf("Took info hooks: %#v with len %d\n", np.INFO.Hooks, len(np.INFO.Hooks))
				if len(np.INFO.Hooks) != 0 {
					return errors.New("Wrong hook at Info level")
				}
				return nil
			},
			true,
		},
		{
			"New Hook add with error",
			func() error {
				// fmt.Printf("Took panic hooks: %#v with len %d\n", np.PANIC.Hooks, len(np.PANIC.Hooks))
				if len(np.PANIC.Hooks) != 0 {
					return errors.New("Wrong hook at Panic level")
				}
				return nil
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.fn()
			// fmt.Printf("Add hook error: %#v\n", err)
			if err != nil && !tt.wantError {
				t.Errorf("Got error %v", err)
			}
		})
	}
}
