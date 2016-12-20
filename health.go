package main

import (
	"io"
	"net/http"
	"time"

	"github.com/gocraft/health"
)

var stream = health.NewStream()

func serve(addr string) error {
	// use 2 minute intervals because the interval does not start when
	// the program starts, but at the beginning of each minute
	// Since traffic generation takes about 35 seconds, the metrics
	// were being split into two intervals, which would be OK in real
	// life, but not for this example
	sink := health.NewJsonPollingSink(2*time.Minute, 5*time.Minute)
	stream.AddSink(sink)
	http.Handle("/health", sink)

	http.HandleFunc("/", handler)
	return http.ListenAndServe(addr, nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	job := stream.NewJob("handler")

	work(w)

	job.Complete(health.Success)
}

func printMetrics(addr string, w io.Writer) error {
	resp, err := http.Get("http://" + addr + "/health")
	if err != nil {
		return err
	}

	_, err = io.Copy(w, resp.Body)
	if err != nil {
		return err
	}

	return resp.Body.Close()
}
