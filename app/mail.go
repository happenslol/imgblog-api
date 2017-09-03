package app

import (
	"github.com/mailgun/mailgun-go"
)

var sender string
var domain string
var mg mailgun.Mailgun

func initMail() {
	domain = Env("MAILGUN_DOMAIN", "")
	sender = Env("MAILGUN_SENDER", "") + "@" + domain

	pubKey := Env("MAILGUN_KEY_PUBLIC", "")
	secKey := Env("MAILGUN_KEY_SECRET", "")
	mg = mailgun.NewMailgun(domain, secKey, pubKey)
}

func SendMail(recipient, content, subject string) error {
	msg := mailgun.NewMessage(
		sender,
		subject,
		content,
		recipient,
	)

	res, id, err := mg.Send(msg)
	Log.Debugf("res: %v, id: %v, err: %v", res, id, err)

	return err
}
