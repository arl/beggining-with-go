package main

import (
	"fmt"
)

type ErrNegativeSqrt float64

func (e ErrNegativeSqrt) Error() string {
	return fmt.Sprintf("cannot Sqrt negative number: %f", e)
}

func Sqrt(x float64) (float64, error) {
	if x < 0 {
		return 0, ErrNegativeSqrt(x)
	}

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

	return z, nil
}

func main() {
	fmt.Println(Sqrt(2))
	fmt.Println(Sqrt(-2))
}
