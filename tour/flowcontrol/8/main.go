package main

import (
	"fmt"
	"math"
)

func Sqrt(x float64) float64 {
	const maxn = 10

	z, prevz := 1.0, 1.0
	for n := 0; n < maxn; n++ {
		// Move z so that z*z gets closer to x.
		z -= (z*z - x) / (2 * z)

		// Early exit if z has stopped changing.
		if z == prevz {
			break
		}

		prevz = z
	}

	return z
}

func main() {
	fmt.Println("Sqrt:     ", Sqrt(1))
	fmt.Println("math.Sqrt:", math.Sqrt(1))
}
