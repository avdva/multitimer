# multitimer
The package provides a convenient way of managing several timers.

## Badges

![Build Status](https://github.com/avdva/multitimer/workflows/golangci-lint/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/avdva/multitimer)](https://goreportcard.com/report/github.com/avdva/multitimer)

## Installation

To start using this package, run:

```sh
$ go get github.com/avdva/multitimer
```

## Description

Timer sends payload to its chan after a delay. One can schedule several events, but only one real timer will be used.
It is safe to use a Timer object concurrently. Timer may drop messages, if the reader does not read fast enough, so it's important to choose a correct channel capacity.

## API

```go
create:
// New returns a timer with capacity set to 1.
New[T any]() 
// NewWithCapacity returns a timer for given capacity.
NewWithCapacity[T any](cap int)

use:
// Schedule schedules a timer to fire after the delay.
// The payload will be sent to C.
Schedule(delay time.Duration, payload T)
// ScheduleAt schedules a timer to fire at the specific moment.
// The payload will be sent to C.
ScheduleAt(when time.Time, payload T)
// Stop cancels all the timers.
Stop() 

```

## Examples

```go
timer := NewWithCapacity[int](10)
for i := 0; i < 10; i++ {
	timer.Schedule(100*time.Millisecond*time.Duration(i+1), i)
}
for i := 0; i < 10; i++ {
	fmt.Println(<-timer.C)
}
```

## Contact

[Aleksandr Demakin](mailto:alexander.demakin@gmail.com)

## License

Source code is available under the [Apache License Version 2.0](/LICENSE).