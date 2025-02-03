package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"encoding/json"

	"github.com/sohWenMing/finance_server/internal/auth"
	dbwritefunctions "github.com/sohWenMing/finance_server/internal/database/db_write_functions"
	"github.com/sohWenMing/finance_server/internal/database/sqlc_generated"
	tokenmapping "github.com/sohWenMing/finance_server/mapping"

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
		user, err, isDupEmail := dbwritefunctions.CreateUser(queries, bodyJson.Email, hashedPassword)
		if err != nil {
			switch isDupEmail {
			case true:
				w.WriteHeader(http.StatusBadRequest)
				w.Header().Set("content-type", "text/plain")
				w.Write([]byte(fmt.Sprintf("email %s is already being used", bodyJson.Email)))
				return
			case false:
				writeInternalError(w)
				return
			}
		}

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
	bearerString := strings.TrimPrefix(authHeader, "Bearer ")
	if bearerString == "" {
		return tokenString, errors.New("token not found in Authorization header in request")
	}
	return bearerString, nil
}
