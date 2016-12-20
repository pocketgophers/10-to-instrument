This is the example code for https://pocketgophers.com/10-to-instrument

Run using `go run main.go metric.go`

Done this way so you don't have to install all metrics libraries just to play with one.

## Support
- `main.go` - harness for running metric test files; test file must implement `serve(addr string)` w
- `all.go` - runs all metrics in current directory

## Implementing a metric
- `serve(addr string) error` - starts the server
- `printMetrics(addr string, w io.Writer)` - writes metrics to w

## Metrics
- `codahale.go` - [github.com/codahale/metrics](https://github.com/codahale/metrics)
- `expvar.go` - [expvar](//golang.org/pkg/expvar)
- `health.go` - [github.com/gocraft/health](https://github.com/gocraft/health)
- `instruments.go` - [github.com/heroku/instruments](https://github.com/heroku/instruments)
- `monkit.go` - [gopkg.in/spacemonkeygo/monkit.v2](https://github.com/spacemonkeygo/monkit/)
- `prometheus.go` [github.com/prometheus/client_golang/prometheus](https://github.com/prometheus/client_golang/prometheus)
- `telemetry.go` - [github.com/arussellsaw/telemetry](https://github.com/arussellsaw/telemetry) FORKED reporters
- `xstats.go` - [github.com/rs/xstats](https://github.com/rs/xstats) NOT HAPPY WITH OUTPUT

## to implement
- `inspect.go` - [github.com/square/inspect/metrics](https://godoc.org/github.com/square/inspect/metrics) RESPONSETIME is not working
- `say.go` - [github.com/go-say/say](github.com/go-say/say)
