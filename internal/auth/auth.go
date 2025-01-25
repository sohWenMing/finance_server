package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var signingMethod jwt.SigningMethod = jwt.SigningMethodHS256

// functionality for generating BCrypt and validating BCrypt password
func GenerateHashedPassword(password string) (hashedPassword string, err error) {
	if len(password) == 0 {
		return "", errors.New("password cannot be empty string")
	}
	inputInBytes := []byte(password)
	hashBytes, err := bcrypt.GenerateFromPassword(inputInBytes, bcrypt.DefaultCost)
	return string(hashBytes), err
}

func CheckHashedPassword(password, hashedPassword string) (err error) {
	passwordInBytes := []byte(password)
	hashedPasswordInBytes := []byte(hashedPassword)
	return bcrypt.CompareHashAndPassword(hashedPasswordInBytes, passwordInBytes)
}

//functionality for generation and validation of jwt token

type CustomClaims struct {
	IsAdmin bool   `json:"is_admin"`
	UserId  string `json:"user_id"`
	jwt.RegisteredClaims
}

type ReturnedClaims struct {
	IsAdmin bool   `json:"is_admin"`
	UserId  string `json:"user_id"`
}

func GenerateJWTToken(userId string, isAdmin bool, validity time.Duration, signingKey []byte) (signedString string, err error) {
	claims := CustomClaims{
		isAdmin,
		userId,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(validity)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(signingMethod, claims)
	signedString, err = token.SignedString(signingKey)
	return signedString, err
}

func ValidateAndReturnClaims(secret []byte, tokenString string) (returnedClaims ReturnedClaims, err error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("signing method does not match requirement")
		}
		return secret, nil
	})
	if err != nil {
		return returnedClaims, err
	}
	claims, ok := token.Claims.(*CustomClaims)
	if ok && token.Valid {
		returnedClaims.IsAdmin = claims.IsAdmin
		returnedClaims.UserId = claims.UserId
		return returnedClaims, nil
	}
	return returnedClaims, errors.New("falsed to parse or invalid token")
}
