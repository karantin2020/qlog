package main

import (
	// "fmt"
	"time"

	"github.com/karantin2020/qlog"
)

func example_new() {
	nlog := qlog.New(qlog.InfoLevel).
		SetOutput(qlog.ColorTemplate)
	// fmt.Printf("%+v\n", nlog)
	nlog.INFO.Msgf("failed to fetch %s", "URL")
	time.Sleep(time.Millisecond * 200)
	nlog.INFO.Msg("failed to fetch 'URL'")

	newlog := nlog.WithFields(
		qlog.F{Key: "service", Value: "new"},
		qlog.F{Key: "source", Value: "after"},
	)
	// fmt.Printf("%+v\n", newlog)
	newlog.INFO.Msgf("failed to fetch %s", "URL")
	newlog.INFO.Msg("failed to fetch 'URL'")
}
