# Gaziri
[![Build Status](https://travis-ci.org/hypersolid/gaziri.svg?branch=master)](https://travis-ci.org/hypersolid/gaziri)
[![Go Report Card](https://goreportcard.com/badge/github.com/hypersolid/gaziri)](https://goreportcard.com/report/github.com/hypersolid/gaziri)
## About
It's a tiny DRY lib for goroutine pooling and throttling.

1. Provide a worker function that takes `interface{}` as input and returns `interface{}` as output.
1. Set your limits on number of running tasks per second `maxWorkersPerSecond` and task in general `maxWorkers`.
1. Shoot your tasks int `Input` channel and listen to `Output`.

## Quickstart
```golang
pool := gaziri.NewPool(
  func(value interface{}) interface{} {
    // works goes here
    time.Sleep(2 * time.Second)
    return value.(int) * 2
  },
  100, //  max workers
  30,  //  max rps
)

go func() {
  for i := 0; i < jobs; i++ {
    pool.Input <- i // enqueue tasks
  }
}()

for i := 0; i < jobs; i++ {
  <-pool.Output // get results
}
```
