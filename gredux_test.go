package gredux

import (
	"math/rand"
	"testing"
	"time"
)

func TestDispatch(t *testing.T) {
	store := New(make(State))
	store.Dispatch(Action{"test", nil})
	store.Reducer(func(state State, action Action) State {
		if action.ID == "test" {
			state["testSuccess"] = true
		}
		return state
	})
	store.Dispatch(Action{"test", nil})
	if _, ok := store.State()["testSuccess"]; !ok {
		t.Fatal("expected state store to have key added by reducer")
	}
}

func TestDispatchUpdate(t *testing.T) {
	store := New(make(State))
	store.Reducer(func(state State, action Action) State {
		if action.ID == "test" {
			state["testSuccess"] = true
		}
		return state
	})
	done := make(chan struct{})
	store.AfterUpdate(func(state State) {
		defer close(done)
		if state["testSuccess"] != true {
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
	initialState := make(State)
	initialState["count"] = 0
	store := New(initialState)
	store.Reducer(func(state State, action Action) State {
		if action.ID == "increment" {
			state["count"] = state["count"].(int) + action.Data.(int)
		}
		if action.ID == "decrement" {
			state["count"] = state["count"].(int) - action.Data.(int)
		}
		return state
	})
	store.Dispatch(Action{"increment", 5})
	val, ok := store.State()["count"]
	if !ok {
		t.Fatal("state didn't have count")
	}
	if val != 5 {
		t.Fatal("count was not incremented")
	}
	store.Dispatch(Action{"increment", 3})
	val, _ = store.State()["count"]
	if val != 8 {
		t.Fatal("count was not incremented")
	}
	store.Dispatch(Action{"decrement", 2})
	val, _ = store.State()["count"]
	if val != 6 {
		t.Fatal("count was not decremented")
	}
}

func TestConcurrentDispatch(t *testing.T) {
	store := New(make(State))
	store.Reducer(func(state State, action Action) State {
		state["test"] = true
		return state
	})
	for i := 0; i < 10; i++ {
		go func() {
			time.Sleep(time.Second * time.Duration(rand.Int()))
			store.Dispatch(Action{"test", nil})
		}()
	}
}

func BenchmarkDispatch(b *testing.B) {
	initialState := make(State)
	initialState["count"] = 0
	store := New(initialState)
	store.Reducer(func(state State, action Action) State {
		if action.ID == "increment" {
			state["count"] = state["count"].(int) + 1
		}
		return state
	})

	for i := 0; i < b.N; i++ {
		store.Dispatch(Action{"increment", nil})
	}
}
