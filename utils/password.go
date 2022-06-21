package utils

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword returns a BCrpt hash of a hashed password
func HashPassword(password string) (string, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		// return an empty hash string and wrap with an error message
		return "", fmt.Errorf("failed to hash the password, %w", err)
	}

	return string(hashPassword), nil
}

// CheckPassword checks if the provided password is correct or not
func CheckPassword(password, hashPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))

}
