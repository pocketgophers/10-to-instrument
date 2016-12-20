package main

import (
	"io"
	"net/http"
)

func serve(addr string) error {
	http.HandleFunc("/", handler)
	return http.ListenAndServe(addr, nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	work(w)
}

func printMetrics(addr string, w io.Writer) error {
	w.Write([]byte("no metrics recorded\n"))
	return nil
}
