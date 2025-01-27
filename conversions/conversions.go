package conversions

import (
	"strconv"
	"time"
)

func GetDateFromDateString(input string) (returnedDateTime time.Time, err error) {
	dateTime, err := time.Parse(time.DateOnly, input)
	return dateTime, err
}

func GetUint64FromString(input string) (largeInt int64, err error) {
	largeInt, err = strconv.ParseInt(input, 10, 64)
	return largeInt, err
}

func GetDateOnlyTimeStamp(input time.Time) (dateOnlyOutput time.Time) {
	y, m, d := input.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}

func GetDateOnlyTimeStampFromDateString(input string) (returnedTimeStamp time.Time, err error) {
	dateTime, err := GetDateFromDateString(input)
	if err != nil {
		return time.Time{}, err
	}
	return GetDateOnlyTimeStamp(dateTime), nil
}
