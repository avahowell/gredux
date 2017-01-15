package gredux

import (
	"math/rand"
	"testing"
	"time"
)

func TestDispatch(t *testing.T) {
	type testState struct {
		success bool
	}
	store := New(testState{false})
	store.Dispatch(Action{"test", nil})
	store.Reducer(func(state State, action Action) (State, interface{}) {
		switch action.ID {
		case "test":
			return testState{true}, nil
		default:
			return state, nil
		}
	})
	store.Dispatch(Action{"test", nil})
	if st := store.State().(testState); !st.success {
		t.Fatal("expected reducer to set success")
	}
}

func TestDispatchUpdate(t *testing.T) {
	type testState struct {
		success bool
	}
	store := New(testState{false})
	store.Reducer(func(state State, action Action) (State, interface{}) {
		switch action.ID {
		case "test":
			return testState{true}, nil
		default:
			return state, nil
		}
	})
	done := make(chan struct{})
	store.AfterUpdate(func(state State) {
		defer close(done)
		if !state.(testState).success {
			t.Fatal()
		}
	})
	store.Dispatch(Action{"test", nil})
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("OnUpdate func was not called after dispatch after 1 second")
	}
}

func TestDispatchIncrementDecrement(t *testing.T) {
	type counterState struct {
		count int
	}
	store := New(counterState{0})
	store.Reducer(func(state State, action Action) (State, interface{}) {
		switch action.ID {
		case "increment":
			return counterState{state.(counterState).count + action.Data.(int)}, nil
		case "decrement":
			return counterState{state.(counterState).count - action.Data.(int)}, nil
		default:
			return state, nil
		}
	})
	store.Dispatch(Action{"increment", 5})
	if val := store.State().(counterState).count; val != 5 {
		t.Fatal("increment did not increment correctly")
	}
	store.Dispatch(Action{"increment", 3})
	if val := store.State().(counterState).count; val != 8 {
		t.Fatal("increment did not increment correctly")
	}
	store.Dispatch(Action{"decrement", 2})
	if val := store.State().(counterState).count; val != 6 {
		t.Fatal("decrement did not decrement correctly")
	}
}

func TestConcurrentDispatch(t *testing.T) {
	type testState struct {
		success bool
	}
	store := New(testState{false})
	store.Reducer(func(state State, action Action) (State, interface{}) {
		return testState{true}, nil
	})
	for i := 0; i < 10; i++ {
		go func() {
			time.Sleep(time.Second * time.Duration(rand.Int()))
			store.Dispatch(Action{"test", nil})
		}()
	}
}

// TestStoreImmutability verifies that mutating the state passed
// to AfterUpdate, Reducer, or returned by State() does not effect the internal state.
func TestStoreImmutability(t *testing.T) {
	type testState struct {
		success bool
		mutated bool
	}
	store := New(testState{false, false})
	store.Reducer(func(state State, action Action) (State, interface{}) {
		st := state.(testState)
		if st.mutated {
			t.Fatal("state was mutated")
		}
		st.mutated = true
		switch action.ID {
		case "test":
			return testState{true, false}, nil
		default:
			return state, nil
		}
	})
	i := 0
	done := make(chan struct{})
	store.AfterUpdate(func(state State) {
		i++
		if i == 2 {
			defer close(done)
		}
		st := state.(testState)
		if st.mutated {
			t.Fatal("state was mutated")
		}
		st.mutated = true
	})
	store.Dispatch(Action{"test", nil})
	store.Dispatch(Action{"test", nil})
	st := store.State().(testState)
	st.mutated = true
	select {
	case <-done:
		// success!
	case <-time.After(time.Second):
		t.Fatal("TestStoreImmutability timed out after one second")
	}
	if store.State().(testState).mutated {
		t.Fatal("store was mutated")
	}
}

func BenchmarkDispatch(b *testing.B) {
	type counterState struct {
		count int
	}
	store := New(counterState{0})
	store.Reducer(func(state State, action Action) (State, interface{}) {
		switch action.ID {
		case "increment":
			return counterState{state.(counterState).count + 1}, nil
		default:
			return state, nil
		}
	})

	for i := 0; i < b.N; i++ {
		store.Dispatch(Action{"increment", nil})
	}
}

func TestConcurrentSelectorSupport(t *testing.T) {
	type testState struct {
		success bool
	}
	store := New(testState{false})
	store.Reducer(func(state State, action Action) (State, interface{}) {
		switch action.ID {
		case "select":
			return state, state.(testState).success
		default:
			return testState{true}, nil
    }
	})
	for i := 0; i < 10; i++ {
		go func() {
			time.Sleep(time.Second * time.Duration(rand.Int()))
			ret := store.Dispatch(Action{"select", nil})
			if ret != false {
				t.Error("Selector failed")
			}
		}()
	}
}
