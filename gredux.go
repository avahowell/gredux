// gredux implements a state atom, similar to the
// javascript library "redux".  The purpose is to greatly simplify writing
// complex, concurrent event-driven systems.
// LICENSE: MIT
package gredux

// state defines the immutable state of the Atom.
type state map[string]interface{}

// Reducer defines a func that receives a 
// new Action when dispatched into the Atom.
type Reducer func(state, Action)

// Atom implements the state atom type, otherwise known as a 'store' in flux.
type Atom struct {
	reducers []Reducer
	st state
}

// New creates a new gredux Atom
func New() (*Atom) {
	at := Atom{
		st: make(state),
	}
	return &at
}

// ActionID defines a unique identifier for an 
// Action, switched over in a reducer func.
type ActionID string

// Action defines an action that can be dispatched in to the Atom.
type Action struct {
	ID ActionID
	data interface{}
}

// AddReducer adds a reducer to the atom
func (at *Atom) AddReducer(r Reducer) {
	at.reducers = append(at.reducers, r)
}

// Dispatch dispatches an Action into the Atom.
func (at *Atom) Dispatch(action Action) {
	for _, reducer := range at.reducers {
		reducer(at.st, action)
	}
}
