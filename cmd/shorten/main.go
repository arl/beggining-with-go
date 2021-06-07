package main

import (
	"flag"
	"fmt"
)

func main() {
	addrFlag := ""
	flag.StringVar(&addrFlag, "addr", ":8080", "server listen address")
	flag.Parse()

	fmt.Println(addrFlag)
}
