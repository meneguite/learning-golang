package main

import (
	"fmt"
	"math"
)

func fibonacci() func() int {
	x := 0
	return func() int {
		s := math.Sqrt(5)
		p1 := 1 / s
		p2 := math.Pow((1+s)/2, float64(x))
		p3 := math.Pow((1-s)/2, float64(x))
		r := p1 * (p2 - p3)
		x++
		return int(r)
	}
}

func main() {
	f := fibonacci()
	for i := 0; i < 10; i++ {
		fmt.Println(f())
	}
}
