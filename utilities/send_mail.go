package utilities

import (
	"net/mail"
	"net/smtp"
	"os"
)

func SendEmail(to, subject, body string) error {
	from := os.Getenv("EMAIL")
	password := os.Getenv("SMTP_PASSWORD")
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Set up authentication information.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Compose the email
	msg := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body + "\r\n")

	// Send the email
	fromAddr := mail.Address{Address: from}
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, fromAddr.Address, []string{to}, msg)
	if err != nil {
		return err
	}
	return nil
}