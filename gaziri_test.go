package gaziri

import "testing"

const (
	capacity = 100
	jobs     = 1000
	rps      = 500
)

func TestPool(t *testing.T) {
	pool := NewPool(
		func(value interface{}) interface{} {
			return value.(int) * 2
		},
		capacity,
		rps,
	)

	go func() {
		for i := 0; i < jobs; i++ {
			pool.Input <- i
		}
	}()

	for i := 0; i < jobs; i++ {
		<-pool.Output
	}
}
