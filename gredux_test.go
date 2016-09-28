package gredux

import (
	"testing"
	"time"
	"math/rand"
)

func TestDispatch(t *testing.T) {
	atom := New(make(State))
	atom.AddReducer(func(state State, action Action) {
		if action.ID == "test" {
			state["testSuccess"] = true
		}
	})
	atom.Dispatch(Action{"test", nil})
	if _, ok := atom.GetState()["testSuccess"]; !ok {
		t.Fatal("expected state atom to have key added by reducer")
	}
}

func TestDispatchIncrementDecrement(t *testing.T) {
	initialState := make(State)
	initialState["count"] = 0
	atom := New(initialState)
	atom.AddReducer(func(state State, action Action) {
		if action.ID == "increment" {
			state["count"] = state["count"].(int) + action.data.(int)
		}
		if action.ID == "decrement" {
			state["count"] = state["count"].(int) - action.data.(int)
		}
	})
	atom.Dispatch(Action{"increment", 5})
	val, ok := atom.GetState()["count"]
	if !ok {
		t.Fatal("state didnt have count")
	}
	if val != 5 {
		t.Fatal("count was not incremented")
	}
	atom.Dispatch(Action{"increment", 3})
	val, _ = atom.GetState()["count"]
	if val != 8 {
		t.Fatal("count was not incremented")
	}
	atom.Dispatch(Action{"decrement", 2})
	val, _ = atom.GetState()["count"]
	if val != 6 {
		t.Fatal("count was not decremented")
	}
}

func TestConcurrentDispatch(t *testing.T) {
	atom := New(make(State))
	atom.AddReducer(func(state State, action Action) {
		state["test"] = true
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
	atom.AddReducer(func(state State, action Action) {
		if action.ID == "increment" {
			state["count"] = state["count"].(int) + 1
		}
	})

	for i := 0; i < b.N; i++ {
		atom.Dispatch(Action{"increment", nil})
	}
}
