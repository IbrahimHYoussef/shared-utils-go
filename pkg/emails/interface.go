// Package emails contains email sender interfaces, SMTP implementation, and
// template helpers.
package emails

import "strings"

// IEmailSender describes a type that can send an email message.
type IEmailSender interface {
	// SendEmail sends an email to one recipient with the given subject and body.
	SendEmail(to string, subject string, body string) error
}

// ReplacePlaceholders replaces {{key}} placeholders in template with values
// from placeholders.
func ReplacePlaceholders(template string, placeholders map[string]string) string {
	result := template
	for key, value := range placeholders {
		placeholder := "{{" + key + "}}"
		result = strings.Replace(result, placeholder, value, -1)
	}
	return result
}
