package validator

import (
	"regexp"
	"strings"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

var (
	hasLetterRegex  = regexp.MustCompile(`[A-Za-z]`)
	hasDigitRegex   = regexp.MustCompile(`\d`)
	hasSpecialRegex = regexp.MustCompile(`[@$!%*#?&]`)
)

// IsValidEmail checks if an email address is valid
func IsValidEmail(email string) bool {
	return emailRegex.MatchString(strings.TrimSpace(email))
}

// IsValidPassword checks if a password meets minimum requirements:
// - At least one letter
// - At least one digit
// - At least one special character (@$!%*#?&)
func IsValidPassword(password string) bool {
	return hasLetterRegex.MatchString(password) &&
		hasDigitRegex.MatchString(password) &&
		hasSpecialRegex.MatchString(password)
}

func IsValidOtp(otp string) bool {
	if len(otp) == 0 {
		return false
	}
	for _, ch := range otp {
		if ch < '0' || ch > '9' {
			return false
		}
	}
	return true
}

func IsValidUuid(uuid string) bool {
	uuidRegex := regexp.MustCompile(`^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$`)
	return uuidRegex.MatchString(uuid)
}
