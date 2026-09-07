// Package validator contains reusable validation interfaces and simple string
// validators.
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

// IsValidEmail reports whether email has a basic valid email address format.
func IsValidEmail(email string) bool {
	return emailRegex.MatchString(strings.TrimSpace(email))
}

// IsValidPassword reports whether password contains at least one letter, one
// digit, and one supported special character from @$!%*#?&.
func IsValidPassword(password string) bool {
	return hasLetterRegex.MatchString(password) &&
		hasDigitRegex.MatchString(password) &&
		hasSpecialRegex.MatchString(password)
}

// IsValidOtp reports whether otp is non-empty and contains only decimal digits.
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

// IsValidUuid reports whether uuid matches the canonical UUID string format.
func IsValidUuid(uuid string) bool {
	uuidRegex := regexp.MustCompile(`^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$`)
	return uuidRegex.MatchString(uuid)
}
