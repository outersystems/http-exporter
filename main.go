package main

import (
	"flag"
	"fmt"
	"httpexporter"
	"log"
	"net/http"
)

var (
	githash   string // $(git rev-list -1 HEAD | cut -c -7)
	goVersion string // $(go version)
	buildDate string // $(date)
)

func main() {
	// flags
	port := flag.String("port", ":8080", "Listening port")
	target := flag.String("target", "http://127.0.0.1:8081", "Target")
	metrics := flag.String("metrics", ":9696", "Port on which the metrics will be available")

	flag.Parse()

	fmt.Printf("Listening on: %s\n", *port)
	fmt.Printf("Redirecting to: %s\n", *target)
	fmt.Printf("Metrics on: %s\n", *metrics)

	// proxy
	prx := httpexporter.New(*target, *metrics)

	// server
	http.Handle("/", prx.Handler())
	log.Printf("Handled and serving\n")
	log.Fatal(http.ListenAndServe(*port, nil))
}
