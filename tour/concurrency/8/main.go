package main

import (
	"fmt"

	"golang.org/x/tour/tree"
)

// Walk walks the tree t sending all values
// from the tree to the channel ch.
func Walk(t *tree.Tree, ch chan int) {
	go func() {
		doWalk(t, ch)
		close(ch)
	}()
}

func doWalk(t *tree.Tree, ch chan int) {
	if t.Left != nil {
		doWalk(t.Left, ch)
	}
	ch <- t.Value
	if t.Right != nil {
		doWalk(t.Right, ch)
	}
}

// Same determines whether the trees
// t1 and t2 contain the same values.
func Same(t1, t2 *tree.Tree) bool {
	ch1 := make(chan int)
	ch2 := make(chan int)

	Walk(t1, ch1)
	Walk(t2, ch2)

	for {
		v1, ok1 := <-ch1
		v2, ok2 := <-ch2

		if ok1 != ok2 || v1 != v2 {
			return false
		}

		if !(ok1 || ok2) {
			break
		}
	}

	return true
}

func main() {
	// Walk
	fmt.Println("[Walk]")
	ch := make(chan int)
	Walk(tree.New(5), ch)
	for v := range ch {
		fmt.Println(v)
	}

	// Same
	fmt.Println("\n[Same]")
	fmt.Println("should be true: ", Same(tree.New(1000), tree.New(1000)))
	fmt.Println("should be false: ", Same(tree.New(1000), tree.New(999)))
}
