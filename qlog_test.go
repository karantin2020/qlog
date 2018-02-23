package qlog_test

import (
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
			if got := qlog.New(tt.args.lvl); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
