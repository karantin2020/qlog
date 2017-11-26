package main

import (
	"github.com/karantin2020/qlog"
	"github.com/karantin2020/qlog/log"
)

func example_default() {
	log.Debug("Debug: failed to fetch URL")
	log.Info("Info: failed to fetch URL")
	log.Error("Error: failed to fetch URL")
	log.INFO.Fields(
		qlog.F{"service", "new"},
		qlog.F{"source", "after"},
	).Msgf("Info: failed to fetch %s", "URL")
	log.Warn("Warn: failed to fetch URL")
	log.Critical("Critical: failed to fetch URL")
	// log.Panic("Panic: failed to fetch URL")
}
