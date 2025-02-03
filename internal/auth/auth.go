package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
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
		fmt.Printf("Error in validate and return claims: %v\n", err)
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

func GenerateRefreshToken() (refreshToken string, writtenByteSlice []byte, err error) {
	allocated := make([]byte, 16)
	//allocate memory for a slice of 16 bytes
	generateErr := generateRandBytes(allocated)
	if generateErr != nil {
		return refreshToken, []byte{}, generateErr
	}
	return hex.EncodeToString(allocated), allocated, nil

}

func generateRandBytes(byteSlice []byte) (err error) {
	_, readErr := rand.Read(byteSlice)
	if readErr != nil {
		return readErr
	}
	//function will write random bytes to the slice that is passed in, no need to return slice because slice is passed by reference
	return nil
}
