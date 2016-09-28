package gredux

import (
	"sync"
)

type (
	// State is the state of the gredux Atom.
	State map[string]interface{}

	// Reducer is the func which receives actions dispatched
	// using Atom.Dispatch() and updates the internal state.
	Reducer func(State, Action)

	// ActionID defines a unique identifier for an Action.  
	// This ID is used to determine state updates in a Reducer.
	ActionID string

	// Action defines a dispatchable data type that triggers updates in the Atom.
	Action struct {
		ID   ActionID
		data interface{}
	}

	// Atom defines an immutable store of state.
	// The current state of the Atom can be received by calling GetState()
	// but the state can only be changed by a Reducer as the result of a Dispatch'd Action.
	Atom struct {
		mu sync.RWMutex
		reducer Reducer
		state State
	}
)

// New instantiates a new gredux Atom. initialState should be an initialized State map.
func New(initialState State) *Atom {
	at := Atom{
		state: initialState,
	}
	return &at
}

// Reducer sets the atom's reducer function to the function `r`.
func (at *Atom) Reducer(r Reducer) {
	at.reducer = r
}

// GetState returns a copy of the current state
func (at *Atom) GetState() State {
	at.mu.RLock()
	defer at.mu.RUnlock()
	currentState := make(State)
	for k, v := range at.state {
		currentState[k] = v
	}
	return currentState
}

// Dispatch dispatches an Action into the Atom.
func (at *Atom) Dispatch(action Action) {
	at.mu.Lock()
	defer at.mu.Unlock()
	at.reducer(at.state, action)
}
