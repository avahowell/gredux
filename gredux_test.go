package gredux

import (
	"math/rand"
	"testing"
	"time"
)

func TestDispatch(t *testing.T) {
	atom := New(make(State))
	atom.Dispatch(Action{"test", nil})
	atom.Reducer(func(state State, action Action) State {
		if action.ID == "test" {
			state["testSuccess"] = true
		}
		return state
	})
	atom.Dispatch(Action{"test", nil})
	if _, ok := atom.State()["testSuccess"]; !ok {
		t.Fatal("expected state atom to have key added by reducer")
	}
}

func TestDispatchUpdate(t *testing.T) {
	atom := New(make(State))
	atom.Reducer(func(state State, action Action) State {
		if action.ID == "test" {
			state["testSuccess"] = true
		}
		return state
	})
	done := make(chan struct{})
	atom.AfterUpdate(func(state State) {
		defer close(done)
		if state["testSuccess"] != true {
			t.Fatal()
		}
	})
	atom.Dispatch(Action{"test", nil})
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("OnUpdate func was not called after dispatch after 1 second")
	}
}

func TestDispatchIncrementDecrement(t *testing.T) {
	initialState := make(State)
	initialState["count"] = 0
	atom := New(initialState)
	atom.Reducer(func(state State, action Action) State {
		if action.ID == "increment" {
			state["count"] = state["count"].(int) + action.data.(int)
		}
		if action.ID == "decrement" {
			state["count"] = state["count"].(int) - action.data.(int)
		}
		return state
	})
	atom.Dispatch(Action{"increment", 5})
	val, ok := atom.State()["count"]
	if !ok {
		t.Fatal("state didnt have count")
	}
	if val != 5 {
		t.Fatal("count was not incremented")
	}
	atom.Dispatch(Action{"increment", 3})
	val, _ = atom.State()["count"]
	if val != 8 {
		t.Fatal("count was not incremented")
	}
	atom.Dispatch(Action{"decrement", 2})
	val, _ = atom.State()["count"]
	if val != 6 {
		t.Fatal("count was not decremented")
	}
}

func TestConcurrentDispatch(t *testing.T) {
	atom := New(make(State))
	atom.Reducer(func(state State, action Action) State {
		state["test"] = true
		return state
	})
	for i := 0; i < 10; i++ {
		go func() {
			time.Sleep(time.Second * time.Duration(rand.Int()))
			atom.Dispatch(Action{"test", nil})
		}()
	}
}

func BenchmarkDispatch(b *testing.B) {
	initialState := make(State)
	initialState["count"] = 0
	atom := New(initialState)
	atom.Reducer(func(state State, action Action) State {
		if action.ID == "increment" {
			state["count"] = state["count"].(int) + 1
		}
		return state
	})

	for i := 0; i < b.N; i++ {
		atom.Dispatch(Action{"increment", nil})
	}
}
