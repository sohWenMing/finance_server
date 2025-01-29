package auth

import (
	"encoding/hex"
	"fmt"
	"strings"
	"testing"
	"time"

	testhelpers "github.com/sohWenMing/finance_server/test_helpers"
)

var testJWTSecret []byte = []byte("this is the test jwt secret")

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

func TestGenerateJWTToken(t *testing.T) {
	signedString, shouldReturn := generateTestToken(t, 10*time.Second)
	if shouldReturn {
		return
	}
	if signedString == "" {
		t.Errorf("signedString returned should not be empty string")
		return
	}
}

func TestValidateAndReturnJWTToken(t *testing.T) {
	type testStruct struct {
		name          string
		validity      time.Duration
		isExpectError bool
	}
	tests := []testStruct{
		{
			"test should pass, token should be valid at time of checking",
			20 * time.Second,
			false,
		},
		{
			"test should fail, token should not be valid at time of checking",
			0 * time.Second,
			true,
		},
	}
	for _, test := range tests {
		signedString, shouldReturn := generateTestToken(t, test.validity)
		if shouldReturn {
			return
		}
		claims, err := ValidateAndReturnClaims(testJWTSecret, signedString)

		switch test.isExpectError {
		case false:
			testhelpers.AssertNoError(t, err)
			if err != nil {
				return
			}
			testhelpers.AssertBool(t, claims.IsAdmin, false)
			testhelpers.AssertStringVals(t, claims.UserId, "user_id")
		case true:
			testhelpers.AssertHasError(t, err)
		}
	}
}

func TestGenerateRefreshToken(t *testing.T) {
	refreshToken, writtenByteSlice, err := GenerateRefreshToken()
	testhelpers.AssertNoError(t, err)
	decodedBytes, err := hex.DecodeString(refreshToken)
	testhelpers.AssertNoError(t, err)
	testhelpers.AssertIntVals(t, len(writtenByteSlice), len(decodedBytes))
	if len(writtenByteSlice) != len(decodedBytes) {
		return
	}
	for i := 0; i < len(writtenByteSlice); i++ {
		if writtenByteSlice[i] != decodedBytes[i] {
			t.Error("bytes decoded and written differ")
			return
		}
	}
	testhelpers.AssertIntVals(t, len(refreshToken), 32)
	fmt.Printf("refreshToken generated: %s\n", refreshToken)
}

func generateTestToken(t *testing.T, validity time.Duration) (string, bool) {
	signedString, err := GenerateJWTToken("user_id", false, validity, testJWTSecret)
	testhelpers.AssertNoError(t, err)
	if err != nil {
		return "", true
	}
	return signedString, false
}

func changeHashedPassword(hashedPassword string) (changedHashPassword string) {
	startChar := 1
	lastCharIndex := len(hashedPassword) - 1
	for fmt.Sprintf("%d", startChar) == hashedPassword[lastCharIndex:] {
		startChar += 1
	}
	changedHashPassword = hashedPassword[:lastCharIndex] + fmt.Sprintf("%d", startChar)
	return changedHashPassword
}
