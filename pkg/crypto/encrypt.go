package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
)

// EncryptionKeyBytes is the size of an AES-256 key.
const EncryptionKeyBytes = 32

// encryptedPrefix names the scheme, so a key or algorithm change can be told
// apart from old ciphertexts.
const encryptedPrefix = "v1:"

// ParseEncryptionKey decodes a base64 (standard or URL, padded or not) key of
// EncryptionKeyBytes bytes, as kept in configuration.
func ParseEncryptionKey(encoded string) ([]byte, error) {
	encoded = strings.TrimSpace(encoded)
	var key []byte
	var err error
	for _, encoding := range []*base64.Encoding{base64.StdEncoding, base64.RawStdEncoding, base64.URLEncoding, base64.RawURLEncoding} {
		if key, err = encoding.DecodeString(encoded); err == nil {
			break
		}
	}
	if err != nil {
		return nil, fmt.Errorf("encryption key is not base64: %w", err)
	}
	if len(key) != EncryptionKeyBytes {
		return nil, fmt.Errorf("encryption key must be %d bytes, got %d", EncryptionKeyBytes, len(key))
	}
	return key, nil
}

// Encrypt seals plaintext with AES-256-GCM under key and returns
// "v1:" + base64(nonce || ciphertext), for secrets kept at rest such as
// access tokens. Each call uses a fresh random nonce.
func Encrypt(key []byte, plaintext string) (string, error) {
	aead, err := newAEAD(key)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, aead.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}
	sealed := aead.Seal(nonce, nonce, []byte(plaintext), nil)
	return encryptedPrefix + base64.StdEncoding.EncodeToString(sealed), nil
}

// Decrypt opens what Encrypt returned. A ciphertext made with another key, or
// changed since, fails.
func Decrypt(key []byte, ciphertext string) (string, error) {
	if !strings.HasPrefix(ciphertext, encryptedPrefix) {
		return "", errors.New("ciphertext has no known prefix")
	}
	sealed, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(ciphertext, encryptedPrefix))
	if err != nil {
		return "", fmt.Errorf("ciphertext is not base64: %w", err)
	}
	aead, err := newAEAD(key)
	if err != nil {
		return "", err
	}
	if len(sealed) < aead.NonceSize() {
		return "", errors.New("ciphertext is too short")
	}
	plaintext, err := aead.Open(nil, sealed[:aead.NonceSize()], sealed[aead.NonceSize():], nil)
	if err != nil {
		return "", errors.New("ciphertext cannot be opened with this key")
	}
	return string(plaintext), nil
}

func newAEAD(key []byte) (cipher.AEAD, error) {
	if len(key) != EncryptionKeyBytes {
		return nil, fmt.Errorf("encryption key must be %d bytes, got %d", EncryptionKeyBytes, len(key))
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return cipher.NewGCM(block)
}
