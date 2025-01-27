package conversions

import (
	"testing"
	"time"

	testhelpers "github.com/sohWenMing/finance_server/test_helpers"
)

func TestGetDateOnlyTimeStamp(t *testing.T) {
	currTimeUTC := time.Now().UTC()
	dateStamp := GetDateOnlyTimeStamp(currTimeUTC)
	testhelpers.AssertStringVals(t, currTimeUTC.Format(time.DateOnly), dateStamp.Format(time.DateOnly))
	//assertion is that the dateonly portion of both is current
}
