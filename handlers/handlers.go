package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"encoding/json"

	"github.com/google/uuid"
	"github.com/sohWenMing/finance_server/internal/auth"
	dbwritefunctions "github.com/sohWenMing/finance_server/internal/database/db_write_functions"
	"github.com/sohWenMing/finance_server/internal/database/sqlc_generated"
	tokenmapping "github.com/sohWenMing/finance_server/mapping"

	errorutils "github.com/sohWenMing/finance_server/error_utils"
	usermapping "github.com/sohWenMing/finance_server/mapping/user_mapping"
)

func CreateUserHandler(queries *sqlc_generated.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bodyJson, shouldReturn := processUserBodyToJSON(r, w)
		if shouldReturn {
			return
		}
		// after decoding, then try to encrypt the password and return hashed password using BCrypt

		hashedPassword, err := auth.GenerateHashedPassword(bodyJson.Password)
		if err != nil {
			writeInternalError(w)
			return
		}

		// if hashing fails, then write 500 and return

		params := sqlc_generated.CreateUserParams{
			ID:             uuid.New(),
			IsAdmin:        false,
			Email:          bodyJson.Email,
			HashedPassword: hashedPassword,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		var user sqlc_generated.User

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
					w.WriteHeader(http.StatusConflict)
					w.Header().Set("content-type", "text/plain")
					w.Write([]byte(fmt.Sprintf("email %s is already being used", params.Email)))
					return
				case false:
					writeInternalError(w)
					return
					//if the error returned is not due to a unique constraint, then it is due to internal error, return 500
				}
			}
			user = createdUser
			break

		}
		// if there is a problem with craating the user in the database, then 500 and early return
		createdUserResJson := usermapping.CreatedUserResponse{
			IsSuccess: true,
			UserId:    user.ID.String(),
		}

		bodyString, err := json.Marshal(createdUserResJson)
		if err != nil {
			writeInternalError(w)
			return
		}

		w.WriteHeader(200)
		w.Header().Set("content-type", "application/json")
		w.Write([]byte(bodyString))

	}

}

func LoginUserHandler(getDuration func() time.Duration, queries *sqlc_generated.Queries, secret []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		bodyJson, shouldReturn := processUserBodyToJSON(r, w)
		if shouldReturn {
			return
		}
		user, dbErr := queries.GetUserByEmail(context.Background(), bodyJson.Email)
		if dbErr != nil {
			writeInternalError(w)
			return
		}
		compareHashErr := auth.CheckHashedPassword(bodyJson.Password, user.HashedPassword)
		if compareHashErr != nil {
			writeUnauthorizedError(w, "email and password do not match")
			return
		}

		generateJWTAndRefreshAndSendResponse(user, getDuration, secret, w, queries)
	}

}

func CheckRefreshTokenAndGetNewJWTHandler(getDuration func() time.Duration,
	queries *sqlc_generated.Queries, secret []byte) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		var refreshTokenBody tokenmapping.RefreshTokenJSON
		decoder := json.NewDecoder(r.Body)
		jsonDecodeErr := decoder.Decode(&refreshTokenBody)
		if jsonDecodeErr != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("content-type", "text/plain")
			w.Write([]byte("request body could not be parsed"))
			return
		}

		refreshToken, err := queries.GetRefreshTokenInfoByToken(context.Background(), refreshTokenBody.RefreshToken)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Header().Set("content-type", "text/plain")
			w.Write([]byte("refresh token is invalid"))
			return
			//if err is returned can take it that the token could not be found, so just write unauthorized and return
		}
		if time.Now().After(refreshToken.ExpiresOn) {
			w.WriteHeader(http.StatusUnauthorized)
			w.Header().Set("content-type", "text/plain")
			w.Write([]byte("refresh token has expired"))
			return
			//refresh token expired, return
		}
		user, err := queries.GetUserById(context.Background(), refreshToken.UserID)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Header().Set("content-type", "text/plain")
			w.Write([]byte("refresh token is invalid"))
			return
		}
		generateJWTAndRefreshAndSendResponse(user, getDuration, secret, w, queries)

	}
}
func generateJWTAndRefreshAndSendResponse(user sqlc_generated.User, getDuration func() time.Duration, secret []byte, w http.ResponseWriter, queries *sqlc_generated.Queries) {
	tokenString, err := auth.GenerateJWTToken(user.ID.String(), user.IsAdmin, getDuration(), secret)
	if err != nil {
		writeInternalError(w)
		return
	}

	refreshToken, err := dbwritefunctions.CreateRefreshToken(queries, getDuration(), user.ID)
	if err != nil {
		writeInternalError(w)
		return
	}

	responseJson := usermapping.LoginResponse{
		IsSuccess:    true,
		AccessToken:  tokenString,
		RefreshToken: refreshToken.Token,
	}
	resBytes, marshalErr := json.Marshal(responseJson)
	if marshalErr != nil {
		writeInternalError(w)
	}
	w.WriteHeader(200)
	w.Header().Set("content-type", "application/json")
	w.Write([]byte(resBytes))
}

func PingHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Header().Set("content-type", "text/plain")
	w.Write([]byte("OK"))
}

func FileServerMiddleWare(fileServerHandler http.Handler) http.Handler {
	return fsMiddleWareGenerator(fileServerHandler)
}

func CheckAccessMiddleware(next func(w http.ResponseWriter, r *http.Request), secret []byte, isAdminRequired bool) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, err := checkJWTTokenValidAndReturnClaims(r, secret)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Header().Set("content-type", "text/plain")
			w.Write([]byte(err.Error()))
			return
		}

		if isAdminRequired {
			if !claims.IsAdmin {

				w.WriteHeader(http.StatusUnauthorized)
				w.Header().Set("content-type", "text/plain")
				w.Write([]byte("user does not have access"))
				return
			}
		}

		next(w, r)
	})

}

func fsMiddleWareGenerator(next http.Handler) http.Handler {
	/*
		http.HandlerFunc takes in a function with a signature of func(http.ResponseWriter http.Request) which
		has a ServeHTTP method, satisfying the http.Handler type. Allows for passing of other handlers into the
		function, which allows the function to act as a middleware layer
	*/
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		next.ServeHTTP(w, r)
	})
}
func processUserBodyToJSON(r *http.Request, w http.ResponseWriter) (usermapping.UserJSON, bool) {
	var bodyJson usermapping.UserJSON
	decoder := json.NewDecoder(r.Body)
	jsonDecodeErr := decoder.Decode(&bodyJson)
	if jsonDecodeErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		return usermapping.UserJSON{}, true
	}
	return bodyJson, false
}

func writeInternalError(w http.ResponseWriter) {

	w.WriteHeader(500)
	w.Header().Set("content-type", "text/plain")
	w.Write([]byte("500: Internal Error"))
}

func writeUnauthorizedError(w http.ResponseWriter, response string) {
	w.WriteHeader(http.StatusUnauthorized)
	w.Header().Set("content-type", "text/plain")
	w.Write([]byte(response))
}
func checkJWTTokenValidAndReturnClaims(r *http.Request, secret []byte) (claims auth.ReturnedClaims, err error) {
	tokenString, err := getJWTBearerToken(r)
	if err != nil {
		return claims, err
	}
	returnedClaims, err := auth.ValidateAndReturnClaims(secret, tokenString)
	if err != nil {
		return claims, err
	}
	return returnedClaims, nil

}
func getJWTBearerToken(r *http.Request) (tokenString string, err error) {
	authHeader := r.Header.Get("Authorization")
	bearerString := strings.TrimPrefix(authHeader, "Bearer: ")
	if bearerString == "" {
		return tokenString, errors.New("token not found in Authorization header in request")
	}
	return bearerString, nil
}
