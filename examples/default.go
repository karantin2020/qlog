package main

import (
	"github.com/karantin2020/qlog"
	"github.com/karantin2020/qlog/log"
)

func example_default() {
	log.Debug("failed to fetch URL")
	log.Info("failed to fetch URL")
	log.Error("failed to fetch URL")
	log.INFO.Fields(
		qlog.F{"service", "new"},
		qlog.F{"source", "after"},
	).Msgf("failed to fetch %s", "URL")
	log.Warn("failed to fetch URL")
}
