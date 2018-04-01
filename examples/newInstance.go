package main

import (
	"github.com/karantin2020/qlog"
)

func example_new() {
	nlog := qlog.New("newInst", qlog.InfoLevel).
		SetOutput(qlog.ColorTemplate)
	// fmt.Printf("%+v\n", nlog)
	nlog.INFO.Msgf("failed to fetch %s", "URL")
	// time.Sleep(time.Millisecond * 200)
	nlog.INFO.Msg("failed to fetch 'URL'")

	newlog := nlog.WithFields(
		qlog.F{Key: "service", Value: "new"},
		qlog.F{Key: "source", Value: "after"},
	)
	// fmt.Printf("%+v\n", newlog)
	newlog.INFO.Msgf("failed to fetch %s", "URL")
	newlog.INFO.Msg("failed to fetch 'URL'")

	newlog2 := nlog.WithFields(
		qlog.F{Key: "kara", Value: 123},
	)
	newlog2.Name = []byte("nInst2")
	newlog2.INFO.Msgf("test sublogger %d", 123)
	newlog2.WARN.Msg("End of sublogger test")
}
