package dbwritefunctions

import (
	"context"
	"time"

	"github.com/google/uuid"
	errorutils "github.com/sohWenMing/finance_server/error_utils"
	"github.com/sohWenMing/finance_server/internal/auth"
	"github.com/sohWenMing/finance_server/internal/database/sqlc_generated"
)

func CreateRefreshToken(queries *sqlc_generated.Queries, validFor time.Duration, userId uuid.UUID) (refreshToken sqlc_generated.RefreshToken, err error) {
	for {
		params, err := MapRefreshTokenParams(validFor, userId)
		if err != nil {
			return refreshToken, err
		}
		//if err returned from MapRefreshToken, is internal error. early return
		refreshToken, err, isUniqueViolation := SaveRefreshTokenToDB(queries, params)
		switch err != nil {
		case true:
			if isUniqueViolation {
				continue
			}
			return refreshToken, err
		case false:
			return refreshToken, nil
		}
	}
}
func SaveRefreshTokenToDB(queries *sqlc_generated.Queries,
	params sqlc_generated.CreateRefreshTokenParams) (refreshToken sqlc_generated.RefreshToken, err error, isUniqueViolation bool) {
	/*
		There are two edge cases that need to be considered
		1. The uuid that is generated violates the unique constraint for the id of the refresh token table
		2. the refresh token string that is generated violates the unique constrant for the refreshToken of the refresh token table
	*/

	refreshToken, err = queries.CreateRefreshToken(context.Background(), params)
	// if error is returned, function needs to return whether or not the error was due to unique violation so program will know it needs to remap and reattempt
	if err != nil {
		isUniqueViolation, _, rawErr := errorutils.CheckIsUniqueConstraintPqError(err)
		return refreshToken, rawErr, isUniqueViolation
	}
	return refreshToken, nil, false

}

func MapRefreshTokenParams(validFor time.Duration, userId uuid.UUID) (params sqlc_generated.CreateRefreshTokenParams, err error) {
	refreshTokenString, _, generateErr := auth.GenerateRefreshToken()
	//if err happens during the generation, constitutes internal error. early return
	if generateErr != nil {
		return params, generateErr
	}
	params = sqlc_generated.CreateRefreshTokenParams{
		ID:        uuid.New(),
		UserID:    userId,
		Token:     refreshTokenString,
		ExpiresOn: time.Now().Add(validFor),
		CreatedOn: time.Now(),
		UpdatedOn: time.Now(),
	}
	return params, nil
}
