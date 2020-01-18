// Command sloth is an HTTP tarpit.
package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func garbage() []byte {
	var (
		num = rand.Intn(1024 * 1024)
		buf = make([]byte, num)
	)

	rand.Seed(time.Now().UnixNano())
	rand.Read(buf)
	return buf
}

func usage() {
	fmt.Printf(
		"sloth: A stupid HTTP tarpit.\n\n" +
			"usage:\n" +
			"  sloth -host {host} -port {port} -wait {N}\n" +
			"arguments:\n" +
			"  flag    description         default\n" +
			"  ----    -----------         -------\n" +
			"  -host   host                localhost\n" +
			"  -port   port                8080\n" +
			"  -wait   wait interval (ms)  8000\n",
	)
}

func tarpit(t *time.Ticker, w http.ResponseWriter, r *http.Request) {
	var buf []byte

	log.Printf("%s\t%s\t%s%s", r.RemoteAddr, r.Method, r.Host, r.URL.Path)

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Connection", "keep-alive")

	for {
		select {
		case <-t.C:
			buf = garbage()
			buf = append(buf, []byte{'\r', '\n'}...)
			w.Write(buf)
			w.(http.Flusher).Flush()
		}
	}
}

func main() {
	flag.Usage = usage
	var (
		host = flag.String("host", "localhost", "host")
		port = flag.String("port", "8080", "port")
		wait = flag.Uint("wait", 8000, "wait interval")

		tick *time.Ticker
	)

	flag.Parse()

	tick = time.NewTicker(time.Duration(*wait) * time.Millisecond)

	http.HandleFunc("/", func(wtr http.ResponseWriter, rdr *http.Request) {
		tarpit(tick, wtr, rdr)
	})

	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", *host, *port), nil))
}
