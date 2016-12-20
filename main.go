package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/rakyll/hey/requester"
)

const addr = "127.0.0.1:3000"

func main() {
	log.SetFlags(log.Lshortfile)

	batch := flag.Bool("batch", false, "if true, will not wait after printing metrics")
	flag.Parse()

	log.Println("serving on", addr)
	go func() {
		log.Fatalln(serve(addr))
	}()

	// wait for the server to start
	time.Sleep(1 * time.Second)

	req, err := http.NewRequest(http.MethodGet, "http://"+addr, nil)
	if err != nil {
		log.Fatalln(err)
	}

	work := &requester.Work{
		Request: req,
		N:       300,
		C:       2,
	}

	log.Println("working…")
	work.Run()

	log.Println("metrics")
	err = printMetrics(addr, os.Stdout)
	if err != nil {
		log.Fatalln(err)
	}

	// wait for interrupt to end
	if !*batch {
		os.Stdout.Write([]byte("\n"))
		log.Println("press ^C to exit")
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
	}

}

func work(w io.Writer) {
	sleepTime := time.Duration(rand.Intn(500)) * time.Millisecond
	time.Sleep(sleepTime)
	fmt.Fprintln(w, sleepTime)
}
