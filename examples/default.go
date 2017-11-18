package main

import (
	// "fmt"
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
	// fmt.Println(log.WARN)
	log.Warn("failed to fetch URL")
	// fmt.Println(log.CRITICAL)
	log.Critical("failed to fetch URL")
	// fmt.Println(log.FATAL)
	log.Panic("failed to fetch URL")
}
