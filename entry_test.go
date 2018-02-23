package qlog_test

import (
	"testing"
	"time"

	"github.com/karantin2020/qlog"
	"github.com/stretchr/testify/assert"
)

func TestLogger_NewEntry(t *testing.T) {
	lgr := qlog.New(qlog.InfoLevel)
	e1 := lgr.INFO.NewEntry()
	time.Sleep(time.Millisecond * 20)
	e2 := lgr.INFO.NewEntry()
	e3 := lgr.WARN.NewEntry()
	tests := []struct {
		name      string
		fn        func() error
		wantError bool
	}{
		{
			"New Entry pass",
			func() error {
				e := lgr.INFO.NewEntry()
				assert.NotNil(t, e)
				return nil
			},
			false,
		},
		{
			"Check entry level",
			func() error {
				assert.Equal(t, e1.Logger.Level, qlog.InitLevel(qlog.InfoLevel))
				return nil
			},
			false,
		},
		{
			"Check entry time",
			func() error {
				assert.Equal(t, e3.Time.After(e2.Time), true)
				return nil
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.fn(); err != nil && !tt.wantError {
				t.Errorf("Got error %v", err)
			}
		})
	}
	e1.Free()
	e2.Free()
	e3.Free()
}