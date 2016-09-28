package gredux

import (
	"testing"
	"time"
	"math/rand"
)

// TestDispatch tests that actions are correctly dispatched to each reducer
func TestDispatch(t *testing.T) {
	atom := New()
	testAction := Action{"TestAction", "just a test"}
	testReducer := func(state State, action Action) {
		if action.ID == testAction.ID {
			state["testSuccess"] = true
		}
	}
	atom.AddReducer(testReducer)
	atom.Dispatch(testAction)
	if _, ok := atom.GetState()["testSuccess"]; !ok {
		t.Fatal("expected state atom to have key added by reducer")
	}
}

// TestConcurrentDispatch verifies that concurrently dispatching actions works as expected.
func TestConcurrentDispatch(t *testing.T) {
	atom := New()
	testReducer := func(state State, action Action) {
		state["test"] = true
	}
	atom.AddReducer(testReducer)
	for i := 0; i < 10; i++ {
		go func() {
			time.Sleep(time.Second * time.Duration(rand.Int()))
			atom.Dispatch(Action{"test", nil})
		}()
	}
}