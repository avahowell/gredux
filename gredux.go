// gredux implements a state atom, similar to the
// javascript library "redux".  The purpose is to greatly simplify writing
// complex, concurrent event-driven systems.
// LICENSE: MIT
package gredux

import (
	"sync"
)

// state defines the immutable state of the Atom.
type State map[string]interface{}

// Reducer defines a func that receives a
// new Action when dispatched into the Atom.
type Reducer func(State, Action)

// Atom implements the state atom type, otherwise known as a 'store' in flux.
type Atom struct {
	mu      sync.RWMutex
	reducer Reducer
	state   State
}

// New creates a new gredux Atom
func New(initialState map[string]interface{}) *Atom {
	at := Atom{
		state: initialState,
	}
	return &at
}

// ActionID defines a unique identifier for an
// Action, switched over in a reducer func.
type ActionID string

// Action defines an action that can be dispatched in to the Atom.
type Action struct {
	ID   ActionID
	data interface{}
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
