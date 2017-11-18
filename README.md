# qLog

Blazing fast, structured, leveled logging in Go.
Like zap or zerolog but more flexible

## Installation

`go get -u github.com/karantin2020/qlog`

## Quick Start

Examples are available in the folder ./examples.

To use default logger.

```go
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
	log.Critical("failed to fetch URL")
	log.Panic("failed to fetch URL")
}
```

Output is:
```sh
2017-11-19T01:12:20.569+0500    INFO    failed to fetch URL     {}
2017-11-19T01:12:20.569+0500    ERROR   failed to fetch URL     {"error":"failed to fetch URL"}
2017-11-19T01:12:20.569+0500    INFO    failed to fetch URL     {"service":"new","source":"after"}
2017-11-19T01:12:20.569+0500    WARN    failed to fetch URL     {}
2017-11-19T01:12:20.569+0500    CRITICAL        failed to fetch URL     {"error":"failed to fetch URL"}
2017-11-19T01:12:20.569+0500    PANIC   failed to fetch URL     {"error":"failed to fetch URL"}
panic: failed to fetch URL

goroutine 1 [running]:
github.com/karantin2020/qlog.(*Entry).errMsg(0xc467284400, 0x4fe87d, 0x13, 0xc467280001)
        $GOPATH/github.com/karantin2020/qlog/entry.go:125 +0x1ab
github.com/karantin2020/qlog.(*Entry).Panic(0xc467284400, 0x4fe87d, 0x13)
        $GOPATH/github.com/karantin2020/qlog/entry.go:156 +0x65
github.com/karantin2020/qlog.(*Notepad).Panic(0xc467280000, 0x4fe87d, 0x13)
        $GOPATH/github.com/karantin2020/qlog/qlog.go:298 +0x56
github.com/karantin2020/qlog/log.Panic(0x4fe87d, 0x13)
        $GOPATH/github.com/karantin2020/qlog/log/log.go:97 +0x41
main.example_default()
        $GOPATH/github.com/karantin2020/qlog/examples/default.go:22 +0x1c8
main.main()
        $GOPATH/github.com/karantin2020/qlog/examples/main.go:4 +0x20
```


If you need to change some configs then use customized version of logger.

```go
import (
	"github.com/karantin2020/qlog"
)

func example_new() {
	nlog := qlog.New(qlog.InfoLevel).
		Timestamp().
		SetOutput(qlog.Template("${time}\t${level}\t${message}\t${fields}\n"))
	nlog.INFO.Msgf("failed to fetch %s", "URL")
	nlog.INFO.Msg("failed to fetch 'URL'")

	newlog := nlog.WithFields(
		qlog.F{"service", "new"},
		qlog.F{"source", "after"},
	)
	newlog.INFO.Msgf("failed to fetch %s", "URL")
	newlog.INFO.Msg("failed to fetch 'URL'")
}
```

Output is:

```sh
2017-11-19T01:15:42.822+0500    info    failed to fetch URL     {}
2017-11-19T01:15:42.822+0500    info    failed to fetch 'URL'   {}
2017-11-19T01:15:42.822+0500    info    failed to fetch URL     {"service":"new","source":"after"}
2017-11-19T01:15:42.822+0500    info    failed to fetch 'URL'   {"service":"new","source":"after"}
```

See documentation in code.

## Performance

For now only text output is implemented. It's performance is equal to uber/zap and zerolog.

## Development Status: Stable

All APIs of qlog and log packages are stable before version 2.

## Contributing

Contributing is welcome.

<hr>

Released under the [MIT License](LICENSE.txt).
