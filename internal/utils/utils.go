package utils

import (
	appError "github.com/sborsh1kmusora/auth/internal/errors"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

func VerifyPassword(hashedPassword string, candidatePassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(candidatePassword))
	if err != nil {
		return appError.ErrInvalidCredentials
	}
	return nil
}
