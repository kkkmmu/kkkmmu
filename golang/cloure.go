package main

import (
	"fmt"
)

func Producer() func(i int) int {
	j := 2

	return func(k int) int {
		return j * k
	}
}

func main() {
	t := Producer()
	fmt.Println(t(5))
	fmt.Println(t(6))
	nextInt := intSeq()
	// See the effect of the closure by calling `nextInt`
	// a few times.
	fmt.Println(nextInt())
	fmt.Println(nextInt())
	fmt.Println(nextInt())

	// To confirm that the state is unique to that
	// particular function, create and test a new one.
	newInts := intSeq()
	fmt.Println(newInts())

}

func intSeq() func() int {
	i := 0
	return func() int {
		i++
		return i
	}
}
