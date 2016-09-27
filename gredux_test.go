package gredux

import (
	"testing"
)

// TestDispatch tests that actions are correctly dispatched to each reducer
func TestDispatch(t *testing.T) {
	atom := New()
	testAction := Action{"TestAction", "just a test"}
	testReducer := func(currentState state, action Action) {
		if action.ID == testAction.ID {
			currentState["testSuccess"] = true
		}
	}
	atom.AddReducer(testReducer)
	atom.Dispatch(testAction)
	if _, ok := atom.st["testSuccess"]; !ok {
		t.Fatal("expected state atom to have key added by reducer")
	}
}