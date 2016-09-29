# gredux
[![Go Report Card](https://goreportcard.com/badge/github.com/johnathanhowell/gredux)](https://goreportcard.com/report/github.com/johnathanhowell/gredux)
[![GoDoc](https://godoc.org/github.com/johnathanhowell/gredux?status.svg)](https://godoc.org/github.com/johnathanhowell/gredux)

gredux is a golang implementation of a [redux](https://github.com/reactjs/redux)-esque state container. The aim is to provide a structure for writing applications which have consistent, predictable behaviour.

## Example Usage

```go

import (
	"github.com/johnathanhowell/gredux"
)

// Create an initial state for the state Atom
initialState := make(gredux.State)
initialState["count"] = 0

// Instantiate a new store around this state
store := gredux.New(initialState)

// Create a reducer which increments "count" when it receives an "increment" 
// action, and decrements when it receives a "decrement" action.
store.Reducer(func(state gredux.State, action gredux.Action) gredux.State {
	if action.ID == "increment" {
		state["count"] = state["count"].(int) + action.data.(int)
	}
	if action.ID == "decrement" {
		state["count"] = state["count"].(int) - action.data.(int)
	}
	return state
})

store.Dispatch(Action{"increment", 5})
store.Dispatch(Action{"decrement", 2})

fmt.Println(store.GetState()["count"].(int)) // prints 3

// Register a func to be called after each state update
store.AfterUpdate(func(state State) {
	fmt.Println(state)
})
store.Dispatch(Action{"decrement", 2})
```

Note that mutating `initialState`, `state` (the state passed to the reducer), or the map returned by `GetState()` will not mutate the store's internal state. `initalState` is copied into the store's state in `New()`, `state` is passed a copy of the state map, and `GetState()` returns a copy of the state map. This does incur a performance penalty, however it provides immutablity assuming callers never access the unexported `store.state` field directly.
