package gredux

import (
	"sync"
)

type (
	// State is the state of the gredux Atom.
	State map[string]interface{}

	// Reducer is the func which receives actions dispatched
	// using Atom.Dispatch() and updates the internal state.
	Reducer func(State, Action) State

	// Action defines a dispatchable data type that triggers updates in the Atom.
	Action struct {
		ID   string
		data interface{}
	}

	// Atom defines an immutable store of state.
	// The current state of the Atom can be received by calling GetState()
	// but the state can only be changed by a Reducer as the result of a Dispatch'd Action.
	Atom struct {
		mu      sync.RWMutex
		reducer Reducer
		state   State
		update  func(State)
	}
)

// New instantiates a new gredux Atom. initialState should be an initialized State map.
func New(initialState State) *Atom {
	at := Atom{
		state: make(State),
		reducer: func(s State, a Action) State {
			return s
		},
	}
	for k, v := range initialState {
		at.state[k] = v
	}
	return &at
}

// Reducer sets the atom's reducer function to the function `r`.
func (at *Atom) Reducer(r Reducer) {
	at.reducer = r
}

// AfterUpdate sets Atom's update func. `update` is called after each
// dispatch with a copy of the new state.
func (at *Atom) AfterUpdate(update func(State)) {
	at.update = update
}

// getState returns a copy of Atom's current state map.
func (at *Atom) getState() State {
	currentState := make(State)
	for k, v := range at.state {
		currentState[k] = v
	}
	return currentState
}

// State returns a copy of the current state.
func (at *Atom) State() State {
	at.mu.RLock()
	defer at.mu.RUnlock()
	return at.getState()
}

// Dispatch dispatches an Action into the Atom.
func (at *Atom) Dispatch(action Action) {
	at.mu.Lock()
	defer at.mu.Unlock()
	at.state = at.reducer(at.getState(), action)
	if at.update != nil {
		at.update(at.getState())
	}
}
