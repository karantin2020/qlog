package main

import (
	"fmt"
	"github.com/karantin2020/qlog"
)

func example_new() {
	nlog := qlog.New(qlog.InfoLevel).
		Timestamp().
		SetOutput(qlog.Template("${time}\t${LEVEL}\t${message}\t${fields}\n"))
	fmt.Printf("%+v\n", nlog)
	nlog.INFO.Msgf("failed to fetch %s", "URL")
	nlog.INFO.Msg("failed to fetch 'URL'")

	newlog := nlog.WithFields(
		qlog.F{"service", "new"},
		qlog.F{"source", "after"},
	)
	fmt.Printf("%+v\n", newlog)
	newlog.INFO.Msgf("failed to fetch %s", "URL")
	newlog.INFO.Msg("failed to fetch 'URL'")
}
