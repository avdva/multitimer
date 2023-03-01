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
It is safe to use a Timer object concurrently. Timer may drop messages, if the reader does not read fast enough, so it's important to choose a correct chan capacity.

## API

```go
create:
New[T any]() // New returns a timer with capacity set to 1.
NewWithCapacity[T any](cap int) // NewWithCapacity returns a timer for given capacity.

```

## Contact

[Aleksandr Demakin](mailto:alexander.demakin@gmail.com)

## License

Source code is available under the [Apache License Version 2.0](/LICENSE).