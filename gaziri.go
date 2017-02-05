package gaziri

import (
	"sync"
	"time"
)

// Pool is a pool of goroutines which processes tasks from Input channel and writes them to Output
type Pool struct {
	maxWorkers int

	Input  chan interface{}
	Output chan interface{}

	waitGroup      *sync.WaitGroup
	spawnedWorkers chan struct{}

	workerFunc func(interface{}) interface{}

	ticker *time.Ticker
}

// NewPool creates pool which spawns workerFunc only
// at maxWorkersPerSecond rate
// and uses maxWorkers when at maximum load
func NewPool(
	workerFunc func(interface{}) interface{},
	maxWorkers int,
	maxWorkersPerSecond int,
) *Pool {
	p := new(Pool)
	p.maxWorkers = maxWorkers
	p.Input = make(chan interface{}, maxWorkers)
	p.Output = make(chan interface{}, maxWorkers)
	p.spawnedWorkers = make(chan struct{}, maxWorkers)
	p.workerFunc = workerFunc
	p.waitGroup = new(sync.WaitGroup)
	p.ticker = time.NewTicker(time.Second / time.Nanosecond / time.Duration(maxWorkersPerSecond))

	go poolManager(p)

	return p
}

// WorkersCount returns number of spawned workers at this moment
func (pool *Pool) WorkersCount() int {
	return len(pool.spawnedWorkers)
}

func poolManager(pool *Pool) {
	cleanup := func() {
		<-pool.spawnedWorkers
		pool.waitGroup.Done()
	}
	for range pool.ticker.C {
		task := <-pool.Input
		pool.spawnedWorkers <- struct{}{}
		pool.waitGroup.Add(1)
		go func() {
			defer cleanup()
			pool.Output <- pool.workerFunc(task)
		}()
	}
	pool.waitGroup.Wait()
}
