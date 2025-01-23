package database

import (
	"testing"

	testhelpers "github.com/sohWenMing/finance_server/test_helpers"
)

func TestDBConnection(t *testing.T) {
	_, err := ConnectToDB("../../../.env")
	testhelpers.AssertNoError(t, err)
}
