package emails

import (
	"net/smtp"
)

type EmailSender struct {
	Host        string
	port        string
	fromAddress string
	address     string
	password    string
}

func New(host, port, fromAddress, address, password string) *EmailSender {
	return &EmailSender{
		Host:        host,
		port:        port,
		fromAddress: fromAddress,
		address:     address,
		password:    password,
	}
}

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
