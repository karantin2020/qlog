# qLog ![Build Status](https://jarvi.ga/api/badges/karantin2020/qlog/status.svg)

Blazing fast, structured, leveled logging in Go.
Like zap or zerolog but more flexible

## Installation

`go get -u github.com/karantin2020/qlog`

Library supports go version => 1.9 because of type alias usage.

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

qlog.Template takes template string as the first argument. Substrings in ${...} are interpreted 
as field names or reserved words (fields, message...). If fields name (${...}) is in capital case
then output will be formatted in capital case too (see level name in examples).

See documentation in code.

## Performance

For now only text output is implemented. It's performance is equal to uber/zap and zerolog.

Benchmark results:

```bash
$ go test -benchmem  -bench=.
goos: linux
goarch: amd64
pkg: github.com/karantin2020/qlog
BenchmarkLogNoOutput-3   	20000000	        57.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkLogEmpty-3      	 5000000	       307 ns/op	      40 B/op	       2 allocs/op
BenchmarkLogDisabled-3   	1000000000	         1.80 ns/op	       0 B/op	       0 allocs/op
BenchmarkInfo-3          	 5000000	       282 ns/op	      40 B/op	       2 allocs/op
BenchmarkFields-3        	 3000000	       460 ns/op	      56 B/op	       4 allocs/op
BenchmarkWithFields-3    	 5000000	       251 ns/op	      40 B/op	       2 allocs/op
PASS
ok github.com/karantin2020/qlog 10.259s
```

## Development Status: Stable

All APIs of qlog and log packages are stable before version 2.

## Contributing

Contributing is welcome.

<hr>

Released under the [MIT License](LICENSE.txt).
