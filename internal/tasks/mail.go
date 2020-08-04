package tasks

import (
	"strconv"

	"clevergo.tech/osenv"
	"github.com/go-mail/mail"
)

func SendMail(to []string, subject, html string) error {
	msg := mail.NewMessage()
	msg.SetHeader("From", osenv.MustGet("MAILER_SENDER"))
	msg.SetHeader("To", to...)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", html)

	port, err := strconv.Atoi(osenv.MustGet("MAILER_PORT"))
	if err != nil {
		return err
	}
	d := mail.NewDialer(osenv.MustGet("MAILER_HOST"), port, osenv.MustGet("MAILER_USERNAME"), osenv.MustGet("MAILER_PASSWORD"))
	d.StartTLSPolicy = mail.MandatoryStartTLS

	return d.DialAndSend(msg)
}
