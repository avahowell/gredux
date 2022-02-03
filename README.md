# gredux
[![Go Report Card](https://goreportcard.com/badge/github.com/avahowell/gredux)](https://goreportcard.com/report/github.com/johnathanhowell/gredux)
[![Build Status](https://travis-ci.org/avahowell/gredux.svg?branch=master)](https://travis-ci.org/johnathanhowell/gredux)
[![codecov](https://codecov.io/gh/avahowell/gredux/branch/master/graph/badge.svg)](https://codecov.io/gh/johnathanhowell/gredux)
[![GoDoc](https://godoc.org/github.com/avahowell/gredux?status.svg)](https://godoc.org/github.com/johnathanhowell/gredux)

gredux is a golang implementation of a [redux](https://github.com/reactjs/redux)-esque state container. The aim is to provide a structure for writing applications which have consistent, predictable behaviour.

## Example Usage

```go

import (
	"github.com/avahowell/gredux"
)

// Create an initial state for the Store
type counterState struct {
	count int
}

// Instantiate a new store around this state
store := gredux.New(counterState{0})

// Create a reducer which increments "count" when it receives an "increment" 
// action, and decrements when it receives a "decrement" action.
store.Reducer(func(state gredux.State, action gredux.Action) gredux.State {
	switch action.ID {
	case "increment":
		return counterState{state.(counterState).count + action.Data.(int)}
	case "decrement":
		return counterState{state.(counterState).count - action.Data.(int)}
	default:
		return state
	}
})

store.Dispatch(Action{"increment", 5})
store.Dispatch(Action{"decrement", 2})

fmt.Println(store.State().(counterState).count) // prints 3

// Register a func to be called after each state update
store.AfterUpdate(func(state State) {
	fmt.Println(state.(counterState).count) // prints the count after every state update
})

// Register a hook for given action that will be invoked everytime that action is dispatched 
// This hook must not dispatch another action otherwise deadlock will happen!
store.AddHook(func(state State) {
	fmt.Println(state.(counterState).count) // prints the count after every state update
}, []string{"increment"})
store.Dispatch(Action{"decrement", 2})
```

## License
The MIT License (MIT)
