package emails

import "strings"

// interface for more mail senders in the future
type IEmailSender interface{
	SendEmail(to string, subject string, body string) error
}


func ReplacePlaceholders(template string, placeholders map[string]string) string {
	result := template
	for key, value := range placeholders {
		placeholder := "{{" + key + "}}"
		result = strings.Replace(result, placeholder, value, -1)
	}
	return result
}
