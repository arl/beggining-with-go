package main

import (
	"bytes"
	"io"
	"os"
	"strings"
)

var org = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
var rot13 = []byte("NOPQRSTUVWXYZABCDEFGHIJKLMnopqrstuvwxyzabcdefghijklm")

type rot13Reader struct {
	r io.Reader
}

func (r rot13Reader) Read(p []byte) (int, error) {
	n, err := r.r.Read(p)

	for i := 0; i < n; i++ {
		pos := bytes.IndexByte(org, p[i])
		if pos == -1 {
			// Ignore characters not defined in the rot13 _cypher_
			continue
		}

		p[i] = rot13[pos]
	}

	return n, err
}

func main() {
	s := strings.NewReader("Lbh penpxrq gur pbqr!")
	r := rot13Reader{s}
	io.Copy(os.Stdout, &r)
}
