package main

import (
	"fmt"
	"time"

	"github.com/hypersolid/gaziri"
)

const (
	jobs = 100
)

func main() {
	pool := gaziri.NewPool(
		func(value interface{}) interface{} {
			time.Sleep(2 * time.Second)
			return value.(int) * 2
		},
		100, //  max workers
		30,  //  max rps
	)

	go func() {
		for i := 0; i < jobs; i++ {
			pool.Input <- i
		}
	}()

	for i := 0; i < jobs; i++ {
		<-pool.Output
		if i%10 == 0 {
			fmt.Println(pool.WorkersCount(), "workers")
		}
	}
}
