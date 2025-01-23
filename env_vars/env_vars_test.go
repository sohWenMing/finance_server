package envvars

import (
	"testing"

	testhelpers "github.com/sohWenMing/finance_server/test_helpers"
)

func TestLoadEnv(t *testing.T) {
	type testStruct struct {
		name          string
		path          string
		isErrExpected bool
	}

	tests := []testStruct{
		{
			"test should pass",
			"../.env",
			false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := LoadEnv(test.path)
			switch test.isErrExpected {
			case false:
				testhelpers.AssertNoError(t, err)
			case true:
				testhelpers.AssertHasError(t, err)
			}
		})
	}
}
