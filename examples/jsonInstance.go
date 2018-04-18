package main

import (
	"github.com/karantin2020/qlog"
)

func json_new() {
	fakeMessage := "failed to fetch 'URL'"
	nlog := qlog.New("JSON", qlog.InfoLevel, qlog.TimeFormat("UnixMicro")).
		SetOutput(qlog.Json())
	nlog.ERROR.Msgf("failed to fetch %s", "URL")
	nlog.INFO.Msg(fakeMessage)
	nlog.Info(fakeMessage)

	newlog := nlog.WithFields(
		qlog.F{Key: "srv", Value: "new"},
		qlog.F{Key: "src", Value: "after"},
	).SetTimeFormat("Unix")
	newlog.Infof("failed to fetch %s", "URL")
	newlog.Info(fakeMessage)
	nlog.Panic(fakeMessage)
}
