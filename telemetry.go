package main

import (
	"io"
	"net/http"
	"time"

	"github.com/arussellsaw/telemetry"
	"pocketgophers.com/10-to-instrument/reporters"
)

var (
	tel          = telemetry.New("", 2*time.Minute)
	requests     = telemetry.NewCounter(tel, "requests", 2*time.Minute)
	responseTime = telemetry.NewAverage(tel, "responseTime", 2*time.Minute)
)

func serve(addr string) error {
	http.Handle("/metrics", reporters.TelemetryHandler{Tel: tel})
	http.HandleFunc("/", handler)
	return http.ListenAndServe(addr, nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	work(w)

	responseTime.Add(tel, time.Since(start).Seconds())
	requests.Add(tel, 1)
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
