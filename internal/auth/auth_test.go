package auth

import (
	"fmt"
	"strings"
	"testing"

	testhelpers "github.com/sohWenMing/finance_server/test_helpers"
)

func TestGenerateHashedPassword(t *testing.T) {
	type testStruct struct {
		testName    string
		input       string
		isExpectErr bool
	}
	tests := []testStruct{
		{
			testName:    "test should pass with no error",
			input:       "test-password",
			isExpectErr: false,
		},
		{
			testName:    "test empty input should fail",
			input:       "",
			isExpectErr: true,
		},
		{
			testName:    "test password too long should fail",
			input:       strings.Repeat("a", 73),
			isExpectErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			_, err := GenerateHashedPassword(test.input)
			switch test.isExpectErr {
			case true:
				testhelpers.AssertHasError(t, err)
			case false:
				testhelpers.AssertNoError(t, err)
			}
		})
	}
}

func TestCheckHashedPassword(t *testing.T) {
	type testStruct struct {
		name          string
		password      string
		retrievedHash string
		isExpectErr   bool
		isAppend      bool
		isTrim        bool
		isChange      bool
	}
	tests := []testStruct{
		{
			name:          "test should pass",
			password:      "test-password",
			retrievedHash: "",
			isExpectErr:   false,
			isAppend:      false,
			isTrim:        false,
			isChange:      false,
		},
		{
			name:          "test should fail - append to hash",
			password:      "test-password",
			retrievedHash: "",
			isExpectErr:   false,
			isAppend:      true,
			isTrim:        false,
			isChange:      false,
		},
		{
			name:          "test should fail - subtract from hash",
			password:      "test-password",
			retrievedHash: "",
			isExpectErr:   true,
			isAppend:      false,
			isTrim:        true,
			isChange:      false,
		},
		{
			name:          "test should fail - change first char hash ",
			password:      "test-password",
			retrievedHash: "",
			isExpectErr:   true,
			isAppend:      false,
			isTrim:        false,
			isChange:      true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			hashedPassword, err := GenerateHashedPassword(test.password)
			test.retrievedHash = hashedPassword
			if test.isAppend {
				test.retrievedHash = hashedPassword + "this should absolutely fucking fail"
			}
			if test.isTrim {
				test.retrievedHash = hashedPassword[:5]
			}
			if test.isChange {
				test.retrievedHash = changeHashedPassword(hashedPassword)
			}
			testhelpers.AssertNoError(t, err)
			switch test.isExpectErr {
			case false:
				testhelpers.AssertNoError(t, CheckHashedPassword(test.password, test.retrievedHash))
			case true:
				testhelpers.AssertHasError(t, CheckHashedPassword(test.password, test.retrievedHash))
			}
		})
	}
}

func changeHashedPassword(hashedPassword string) (changedHashPassword string) {
	startChar := 1
	lastCharIndex := len(hashedPassword) - 1
	for fmt.Sprintf("%d", startChar) == hashedPassword[lastCharIndex:] {
		startChar += 1
	}
	changedHashPassword = hashedPassword[:lastCharIndex] + fmt.Sprintf("%d", startChar)
	fmt.Printf("hashedPassword: %s\n", hashedPassword)
	fmt.Printf("changedHashPassword: %s\n", changedHashPassword)
	return changedHashPassword
}
