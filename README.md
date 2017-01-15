# gredux
[![Go Report Card](https://goreportcard.com/badge/github.com/johnathanhowell/gredux)](https://goreportcard.com/report/github.com/johnathanhowell/gredux)
[![Build Status](https://travis-ci.org/johnathanhowell/gredux.svg?branch=master)](https://travis-ci.org/johnathanhowell/gredux)
[![codecov](https://codecov.io/gh/johnathanhowell/gredux/branch/master/graph/badge.svg)](https://codecov.io/gh/johnathanhowell/gredux)
[![GoDoc](https://godoc.org/github.com/johnathanhowell/gredux?status.svg)](https://godoc.org/github.com/johnathanhowell/gredux)

gredux is a golang implementation of a [redux](https://github.com/reactjs/redux)-esque state container. The aim is to provide a structure for writing applications which have consistent, predictable behaviour.

## Example Usage

```go

import (
	"github.com/johnathanhowell/gredux"
)

// Create an initial state for the Store
type counterState struct {
	count int
}

// Instantiate a new store around this state
store := gredux.New(counterState{0})

// Create a reducer which increments "count" when it receives an "increment"
// action, and decrements when it receives a "decrement" action.
// Note: the second return value is the value returned by Dispatch -- useful for selectors
store.Reducer(func(state gredux.State, action gredux.Action) (gredux.State, interface{}) {
	switch action.ID {
	case "increment":
		return counterState{state.(counterState).count + action.Data.(int)}, nil
	case "decrement":
		return counterState{state.(counterState).count - action.Data.(int)}, nil
  case "get"
		return state, state.(counterState).count
	default:
		return state, nil
	}
})

store.Dispatch(Action{"increment", 5})
store.Dispatch(Action{"decrement", 2})

fmt.Println(store.Dispatch(Action{"get", nil})) // prints 3

// Register a func to be called after each state update
store.AfterUpdate(func(state State) {
  fmt.Println(store.Dispatch(Action{"get", nil})) // prints the count after every state update
})
store.Dispatch(Action{"decrement", 2})
```

## License
The MIT License (MIT)
