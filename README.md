# cow
Package cow implements a Call Out Wheel which schedules and fires callbacks in the background of your Go programs.

This library was inspired by a Uber Go library which was later made private.

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
