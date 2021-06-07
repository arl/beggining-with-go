package main

import "fmt"

// fibonacci is a function that returns
// a function that returns an int.
func fibonacci() func() int {
	// First two numbers of the Fibonacci sequence.
	a, b := 0, 1
	return func() int {
		// Save current number.
		ret := a

		// Compute next numbers (b -> a and b -> a+b)
		a, b = b, a+b

		return ret
	}
}

func main() {
	f := fibonacci()
	for i := 0; i < 10; i++ {
		fmt.Println(f())
	}
}
