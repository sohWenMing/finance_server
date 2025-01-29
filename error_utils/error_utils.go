package errorutils

import (
	"errors"

	"github.com/lib/pq"
)

func CheckPqError(err error) (pqErr *pq.Error, rawErr error) {
	rawErr = err
	for err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			return pqErr, rawErr
		}
		rawErr = errors.Unwrap(err)
	}
	return nil, err
}

// function checks if the error passed in is a postgres related error, will keep unwrapping until the is nil

func CheckIsUniqueConstraintPqError(err error) (isUniqueViolation bool, pqErr *pq.Error, rawErr error) {
	pqErr, rawErr = CheckPqError(err)
	if pqErr == nil {
		return false, nil, rawErr
	}
	if pqErr.Code == "23505" {
		return true, pqErr, rawErr
	}
	return false, pqErr, rawErr

}
