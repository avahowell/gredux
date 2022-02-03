package gredux

import (
	"sync"
)

type (
	// State is the state of the gredux Store.
	State interface{}

	// Reducer is the func which receives actions dispatched
	// using Store.Dispatch() and updates the internal state.
	Reducer func(State, Action) State

	// Action defines a dispatchable data type that triggers updates in the Store.
	Action struct {
		ID   string
		Data interface{}
	}

	// Store defines an immutable store of state.
	// The current state of the Store can be received by calling State()
	// but the state can only be changed by a Reducer as the result of a Dispatch'd Action.
	Store struct {
		mu      sync.RWMutex
		reducer Reducer
		state   State
		update  func(State)
		hooks   map[string][]func(State)
	}
)

// New instantiates a new gredux Store.
// initialState should be the struct used to define the Store's state.
func New(initialState State) *Store {
	st := Store{
		reducer: func(s State, a Action) State {
			return s
		},
		state: initialState,
		hooks: map[string][]func(State){},
	}
	return &st
}

// Reducer sets the store's reducer function to the function `r`.
func (st *Store) Reducer(r Reducer) {
	st.reducer = r
}

// AfterUpdate sets Store's update func. `update` is called after each
// dispatch with a copy of the new state.
func (st *Store) AfterUpdate(update func(State)) {
	st.update = update
}

// AddHook adds a hook which is invoked for specific actions that are specified as second argument
// The hook must not dispatch another actions otherwise deadlock will happen!
func (st *Store) AddHook(hook func(State), actions []string)  {
	for _, action := range actions {
		if actionHooks, exists := st.hooks[action]; exists{
			st.hooks[action] = append(actionHooks, hook)
		} else {
			st.hooks[action] = []func(State){hook}
		}
	}
}

// getState returns a copy of Store's current state map.
func (st *Store) getState() State {
	return st.state
}

// State returns a copy of the current state.
func (st *Store) State() State {
	st.mu.RLock()
	defer st.mu.RUnlock()
	return st.getState()
}

// Dispatch dispatches an Action into the Store.
func (st *Store) Dispatch(action Action) {
	st.mu.Lock()
	defer st.mu.Unlock()

	st.state = st.reducer(st.getState(), action)
	if st.update != nil {
		st.update(st.getState())
	}

	actionHooks := st.hooks[action.ID]
	for _, hook := range actionHooks {
		hook(st.getState())
	}
}
