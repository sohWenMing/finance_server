package testhelpers

import "testing"

func AssertNoError(t *testing.T, err error) {
	if err != nil {
		t.Errorf("didn't expect error, got %v\n", err)
	}
}

func AssertHasError(t *testing.T, err error) {
	if err == nil {
		t.Error("expected error, didn't get one\n")
	}
}

func AssertIntVals(t *testing.T, got, want int) {
	if got != want {
		t.Errorf("got: %d\nwant: %d", got, want)
	}
}

func AssertStringVals(t *testing.T, got, want string) {
	if got != want {
		t.Errorf("got: %s\nwant: %s", got, want)
	}
}

func AssertBool(t *testing.T, got, want bool) {

	if got != want {
		t.Errorf("got: %t\nwant: %t", got, want)
	}
}
