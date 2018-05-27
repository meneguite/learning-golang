package main

import (
	"math/rand"
	"sync"
	"time"
)

func executeFunction(wg *sync.WaitGroup, x int) {
	rt := rand.Int31n(1000)
	time.Sleep(time.Duration(rt) * time.Millisecond)

	println("Goroutine: ", x)
	wg.Done()
}

func main() {
	println("Starting")
	rand.Seed(time.Now().Unix())

	var wg sync.WaitGroup
	wg.Add(10)

	for i := 0; i < 10; i++ {
		go executeFunction(&wg, i)
	}

	wg.Wait()

	println("Finished")
}
