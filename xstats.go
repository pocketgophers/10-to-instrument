package main

import (
	"io"
	"net/http"
	"time"

	"github.com/rs/xstats"
	"github.com/rs/xstats/prometheus"
)

var sender = prometheus.NewHandler()
var stats = xstats.New(sender)

func serve(addr string) error {
	http.Handle("/metrics", sender)

	http.HandleFunc("/", handler)
	return http.ListenAndServe(addr, nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	begin := time.Now()

	work(w)

	stats.Histogram("responseTime", time.Since(begin).Seconds())
	stats.Count("requests", 1)
}

func printMetrics(addr string, w io.Writer) error {
	resp, err := http.Get("http://" + addr + "/metrics")
	if err != nil {
		return err
	}

	_, err = io.Copy(w, resp.Body)
	if err != nil {
		return err
	}

	return resp.Body.Close()
}
