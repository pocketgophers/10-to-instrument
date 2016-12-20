package main

import (
	"io"
	"net/http"
	"time"

	"github.com/codahale/metrics"
)

var (
	requests     = metrics.Counter("requests")
	responseTime = metrics.NewHistogram("responseTime", 0, time.Minute.Nanoseconds(), 4)
)

func serve(addr string) error {
	http.HandleFunc("/", handler)
	return http.ListenAndServe(addr, nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	begin := time.Now()

	work(w)

	responseTime.RecordValue(time.Since(begin).Nanoseconds())
	requests.Add()
}

func printMetrics(addr string, w io.Writer) error {
	resp, err := http.Get("http://" + addr + "/debug/vars")
	if err != nil {
		return err
	}

	_, err = io.Copy(w, resp.Body)
	if err != nil {
		return err
	}

	return resp.Body.Close()
}
