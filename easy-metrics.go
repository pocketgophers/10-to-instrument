package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/admobi/easy-metrics"
)

var (
	requests        = metrics.NewCounter("requests")
	avgResponseTime = metrics.NewGauge("avgResponseTime")
	responseTime    = &timing{}
)

func serve(addr string) error {
	r, err := metrics.NewTrackRegistry("stats", 100, time.Second, false)
	if err != nil {
		log.Fatalln(err)
	}

	err = r.AddMetrics(requests, avgResponseTime)
	if err != nil {
		log.Fatalln(err)
	}

	http.HandleFunc("/", handler)
	return http.ListenAndServe(addr, nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	begin := time.Now()

	work(w)

	responseTime.Observe(time.Since(begin))
	requests.Inc()
}

func printMetrics(addr string, w io.Writer) error {
	fmt.Fprintf(w, "Please look at http://%s/easy-metrics?show=stats\n", addr)
	return nil
}

type timing struct {
	count int64
	sum   time.Duration
}

func (t *timing) Observe(d time.Duration) {
	t.count++
	t.sum += d

	avgResponseTime.Set(t.sum.Seconds() / float64(t.count))
}

func (t timing) String() string {
	avg := time.Duration(t.sum.Nanoseconds() / t.count)
	return fmt.Sprintf("\"%v\"", avg)
}
