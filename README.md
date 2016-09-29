# gredux

gredux is a golang implementation of a [redux](https://github.com/reactjs/redux)-esque state container. The aim is to provide a structure for writing applications which have consistent, predictable behaviour.

## Example Usage

```go

import (
	"github.com/johnathanhowell/gredux"
)

// Create an initial state for the state Atom
initialState := make(gredux.State)
initialState["count"] = 0

// Instantiate a new atom around this state
atom := gredux.New(initialState)

// Create a reducer which increments "count" when it receives an "increment" 
// action, and decrements when it receives a "decrement" action.
atom.Reducer(func(state gredux.State, action gredux.Action) gredux.State {
	if action.ID == "increment" {
		state["count"] = state["count"].(int) + action.data.(int)
	}
	if action.ID == "decrement" {
		state["count"] = state["count"].(int) - action.data.(int)
	}
	return state
})

atom.Dispatch(Action{"increment", 5})
atom.Dispatch(Action{"decrement", 2})
fmt.Println(atom.GetState()["count"].(int)) // prints 3
```

Note that mutating `initialState`, `state` (the state passed to the reducer), or the map returned by `GetState()` will not mutate the atom's internal state. `initalState` is copied into the atom's state in `New()`, `state` is passed a copy of the state map, and `GetState()` returns a copy of the state map. This does incur a performance penalty, however it provides immutablity assuming callers never access the unexported `atom.state` field directly.
