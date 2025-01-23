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
