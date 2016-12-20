package main

import (
	"expvar"
	"fmt"
	"io"
	"net/http"
	"time"
)

var (
	requests     = expvar.NewInt("requests")
	responseTime = &timing{}
)

func serve(addr string) error {
	expvar.Publish("responseTime", responseTime)

	http.HandleFunc("/", handler)
	return http.ListenAndServe(addr, nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	begin := time.Now()

	work(w)

	responseTime.Observe(time.Since(begin))
	requests.Add(1)
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

type timing struct {
	count int64
	sum   time.Duration
}

func (t *timing) Observe(d time.Duration) {
	t.count++
	t.sum += d
}

func (t timing) String() string {
	avg := time.Duration(t.sum.Nanoseconds() / t.count)
	return fmt.Sprintf("\"%v\"", avg)
}
