package app

import (
	"github.com/mailgun/mailgun-go"
)

var sender string
var domain string
var mg mailgun.Mailgun

const (
	WelcomeMail    = "templates/welcome.template.html"
	NewsletterMail = "templates/newsletter.template.html"
	NewPostMail    = "templates/new-post.template.html"
)

type MailContent map[string]interface{}

func initMail() {
	domain = Env("MAILGUN_DOMAIN", "")
	sender = Env("MAILGUN_SENDER", "") + "@" + domain

	pubKey := Env("MAILGUN_KEY_PUBLIC", "")
	secKey := Env("MAILGUN_KEY_SECRET", "")
	mg = mailgun.NewMailgun(domain, secKey, pubKey)
}

func SendMail(
	content MailContent,
	template string,
	subject string,
	recipients ...string,
) error {
	msg := mailgun.NewMessage(
		sender,
		subject,
		"some content",
	)

	for _, r := range recipients {
		msg.AddRecipient(r)
	}

	res, id, err := mg.Send(msg)
	Log.Debugf("res: %v, id: %v, err: %v", res, id, err)

	return err
}
