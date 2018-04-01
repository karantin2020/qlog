package main

import (
	"time"

	"github.com/karantin2020/qlog"
)

func json_new() {
	nlog := qlog.New("JSON", qlog.InfoLevel).
		SetOutput(qlog.Json())
	// fmt.Printf("%+v\n", nlog)
	nlog.ERROR.Msgf("failed to fetch %s", "URL")
	time.Sleep(time.Millisecond * 200)
	nlog.INFO.Msg("failed to fetch 'URL'")

	newlog := nlog.WithFields(
		qlog.F{Key: "srv", Value: "new"},
		qlog.F{Key: "src", Value: "after"},
	)
	// fmt.Printf("%+v\n", newlog)
	newlog.INFO.Msgf("failed to fetch %s", "URL")
	newlog.INFO.Msg("failed to fetch 'URL'")
}
