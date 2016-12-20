package main

import (
	"io"
	"net/http"
	"time"

	"github.com/square/inspect/metrics"
)

var (
	requests     = metrics.NewCounter()
	responseTime = metrics.NewStatsTimer(time.Nanosecond, 300)
)

func serve(addr string) error {
	m := metrics.NewMetricContext("webapp")
	m.Register(requests, "requests")
	m.Register(responseTime, "responseTime")
	http.HandleFunc("/metrics", m.HttpJsonHandler)

	http.HandleFunc("/", handler)
	return http.ListenAndServe(addr, nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	t := responseTime.Start()

	work(w)

	responseTime.Stop(t)
	requests.Add(1)
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
