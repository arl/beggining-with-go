package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/arl/shorten"
)

func main() {
	addrFlag := ""
	flag.StringVar(&addrFlag, "addr", ":8080", "server listen address")
	flag.Parse()

	s := newServer(shorten.NewMemoryStore())

	if err := http.ListenAndServe(addrFlag, s); err != nil {
		log.Fatalf("ListenAndServe: %v", err)
	}
}
