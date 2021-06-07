package main

import (
	"strings"

	"golang.org/x/tour/wc"
)

func WordCount(s string) map[string]int {
	m := make(map[string]int)

	words := strings.Fields(s)
	for _, word := range words {
		m[word]++
	}

	return m
}

func main() {
	wc.Test(WordCount)
}
