package crypto

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// GenerateOTP generates a numeric one-time password with exactly length digits.
//
// The generated value is produced with crypto/rand and left-padded with zeros if
// needed. GenerateOTP returns an error when length is less than 1.
func GenerateOTP(length int) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("OTP length must be greater than 0")
	}

	// Calculate max value: 10^length
	maxValue := big.NewInt(10)
	maxValue.Exp(maxValue, big.NewInt(int64(length)), nil)

	// Generate a random number between 0 and (10^length - 1)
	n, err := rand.Int(rand.Reader, maxValue)
	if err != nil {
		return "", fmt.Errorf("failed to generate OTP: %w", err)
	}

	// Format with leading zeros to ensure specified length
	formatStr := fmt.Sprintf("%%0%dd", length)
	otp := fmt.Sprintf(formatStr, n.Int64())
	return otp, nil
}
