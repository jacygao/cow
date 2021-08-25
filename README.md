# cow 

![CI](https://github.com/jacygao/cow/actions/workflows/build.yml/badge.svg?branch=master)
[![codecov](https://codecov.io/gh/jacygao/cow/branch/master/graph/badge.svg)](https://codecov.io/gh/jacygao/cow)

Package cow implements a Call Out Wheel which schedules and fires callbacks in the background of your Go programs.

## Sample Usage

```

// Init call out wheel with default configuration
cli := cow.New()

// Start spins the wheel
cli.Start()

// Gracefully stop the wheel
defer cli.Stop()

// Schedule a callback in 5 seconds
cli.Schedule(5 * time.Seconds, []byte("your data"), func(data []byte){
    fmt.Printf("callback is triggerred with data %s", string(data))
})

```

## Configuration

```
// Init call out wheel with custom tick interval
cli := cow.New(cow.WithTickInterval(time.Second * 10))
```
