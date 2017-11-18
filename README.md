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

See documentation in code.

## Performance

For now only text output is implemented. It's performance is equal to uber/zap and zerolog.

## Development Status: Stable

All APIs of qlog and log packages are stable before version 2.

## Contributing

Contributing is welcome.

<hr>

Released under the [MIT License](LICENSE.txt).
