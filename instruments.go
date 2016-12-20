package main

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/heroku/instruments"
)

var (
	requests     = instruments.NewCounter()
	responseTime = instruments.NewTimer(-1) // in ms
)

func serve(addr string) error {
	http.HandleFunc("/", handler)
	return http.ListenAndServe(addr, nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	begin := time.Now()

	work(w)

	responseTime.Since(begin)
	requests.Update(1)
}

func printMetrics(addr string, w io.Writer) error {
	fmt.Fprintf(w, "requests: %v\n", requests.Snapshot())

	fmt.Fprint(w, "response time distribution:\n")
	times := responseTime.Snapshot()
	for _, percent := range []float64{10, 25, 50, 75, 90, 95, 99} {
		q := instruments.Quantile(times, percent/100)
		d := time.Duration(q) * time.Millisecond

		fmt.Fprintf(w, "  %v%% in %v\n", percent, d)
	}
	return nil
}
