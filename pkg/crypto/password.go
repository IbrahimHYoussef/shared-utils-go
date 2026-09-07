// Package crypto contains helpers for password hashing, UUID generation, and
// one-time password generation.
package crypto

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword generates a bcrypt hash of password using bcrypt.DefaultCost.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash reports whether password matches a bcrypt hash.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
