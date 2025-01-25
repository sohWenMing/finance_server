package auth

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

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
	checkErr := bcrypt.CompareHashAndPassword(hashedPasswordInBytes, passwordInBytes)
	fmt.Printf("checkErr: %v\n", checkErr)
	return bcrypt.CompareHashAndPassword(hashedPasswordInBytes, passwordInBytes)
}
