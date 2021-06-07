package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func main() {
	addrFlag := ""
	flag.StringVar(&addrFlag, "addr", ":8080", "server listen address")
	flag.Parse()

	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello world!")
	})

	if err := http.ListenAndServe(addrFlag, nil); err != nil {
		log.Fatalf("ListenAndServe: %v", err)
	}
}
