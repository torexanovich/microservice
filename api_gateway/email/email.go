package email

import (
	"net/smtp"
)

func SendEmail(to []string, message []byte) error {
	from := "torexanovich.l@gmail.com"
	password := "lsroncdnplvxpeas"

	// smtp server configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Sending email.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	return err

}