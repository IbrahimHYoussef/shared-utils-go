package crypto

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

// DefaultSecretTokenBytes is the number of random bytes GenerateSecretToken
// uses when numBytes is less than 1.
const DefaultSecretTokenBytes = 32

// GenerateSecretToken returns numBytes of crypto/rand randomness encoded as
// unpadded base64url, safe to put in URLs (links in emails, API keys, session
// tokens). When numBytes is less than 1, DefaultSecretTokenBytes is used.
func GenerateSecretToken(numBytes int) (string, error) {
	if numBytes < 1 {
		numBytes = DefaultSecretTokenBytes
	}
	buf := make([]byte, numBytes)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("failed to generate secret token: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}

// HashToken returns the lowercase hex SHA-256 of token, for storing
// high-entropy secrets (refresh tokens, link tokens, OTP codes) without
// keeping them in plaintext. Use HashPassword for passwords.
func HashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

// CheckTokenHash reports whether token hashes to hash, comparing in constant
// time.
func CheckTokenHash(token, hash string) bool {
	return subtle.ConstantTimeCompare([]byte(HashToken(token)), []byte(hash)) == 1
}
