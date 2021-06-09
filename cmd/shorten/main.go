package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/arl/shorten"
)

func main() {
	addrFlag := ""
	flag.StringVar(&addrFlag, "addr", ":8080", "server listen address")
	flag.Parse()

	s := newServer(shorten.NewMemoryStore())

	go func() {
		if err := http.ListenAndServe(addrFlag, s); err != nil {
			log.Fatalf("ListenAndServe: %v", err)
		}
	}()
	log.Println("listening on", addrFlag)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch
	log.Println("quitting...")

	stats := s.statistics()
	fmt.Println("statistics:")
	fmt.Println("shortened:", stats.Shortened)
	fmt.Println("redirects:", stats.Redirected)
}
