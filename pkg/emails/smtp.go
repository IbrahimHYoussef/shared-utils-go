package emails

import (
	"net/smtp"
)

// EmailSender sends HTML email through an SMTP server.
type EmailSender struct {
	// Host is the SMTP server host name.
	Host        string
	port        string
	fromAddress string
	address     string
	password    string
}

// New returns an SMTP EmailSender configured with server and authentication
// settings.
func New(host, port, fromAddress, address, password string) *EmailSender {
	return &EmailSender{
		Host:        host,
		port:        port,
		fromAddress: fromAddress,
		address:     address,
		password:    password,
	}
}

// SendEmail sends an HTML email to one recipient through the configured SMTP
// server.
func (e EmailSender) SendEmail(to string, subject string, body string) error {
	toList := []string{to}
	auth := smtp.PlainAuth("", e.address, e.password, e.Host)
	adderess := e.Host + ":" + e.port

	// Add MIME headers for HTML email support
	message := []byte("MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
		"Subject: " + subject + "\r\n" +
		"From: " + e.fromAddress + "\r\n" +
		"To: " + to + "\r\n" +
		"\r\n" +
		body + "\r\n")

	err := smtp.SendMail(adderess, auth, e.fromAddress, toList, message)
	if err != nil {
		return err
	}
	return nil
}
