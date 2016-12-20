package main

import (
	"io"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	requests = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "requests",
		Help: "number of requests",
	})
	responseTime = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name: "responseTime",
		Help: "response time",
	})
)

func serve(addr string) error {
	prometheus.MustRegister(requests)
	prometheus.MustRegister(responseTime)
	http.Handle("/metrics", prometheus.Handler())

	http.HandleFunc("/", handler)
	return http.ListenAndServe(addr, nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	work(w)

	responseTime.Observe(time.Since(start).Seconds())
	requests.Inc()
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
