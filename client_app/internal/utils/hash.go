package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

const (
	// nkhaali defaute li howa 10 
	passwordHashCost = bcrypt.DefaultCost
)

func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), passwordHashCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedBytes), nil
}
// compare the password from db annd the entred one
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
