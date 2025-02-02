package dbwritefunctions

import (
	"context"
	"strings"
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
func CreateUser(queries *sqlc_generated.Queries, email, password string) (user sqlc_generated.User, err error, isDupEmail bool) {
	params := sqlc_generated.CreateUserParams{
		ID:             uuid.New(),
		IsAdmin:        false,
		Email:          email,
		HashedPassword: password,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	/*
		in a loop, check to see if the attempted creation of the user returns any errors
		in the event that an error is returned, then we need to check if the error is a unique constaint violation
		if it is a unique violation there are two possible cases:
			1 - the uuid has been duplicated. while this is rare, it should still be a case that is handled by creating a new uuid
			2 - the email has been duplicated. in which case, we should just be returning the error and alerting the user
	*/
	for {
		createdUser, checkErr := queries.CreateUser(context.Background(), params)
		if checkErr != nil {
			isUniqueViolation, pqErr, _ := errorutils.CheckIsUniqueConstraintPqError(checkErr)
			//first check to see if the violation is due to a violation of a unique constraint
			switch isUniqueViolation {
			case true:
				if strings.Contains(pqErr.Message, "unique_user_id") {
					params.ID = uuid.New()
					continue
				}
				return user, checkErr, true
			case false:
				return user, checkErr, false
				//if the error returned is not due to a unique constraint, then it is due to internal error, return 500
			}
		}
		user = createdUser
		return user, nil, false
	}
	// if there is a problem with craating the user in the database, then 500 and early return
}
