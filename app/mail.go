package app

import (
	"bytes"
	"html/template"

	"github.com/mailgun/mailgun-go"
)

var sender string
var domain string
var mg mailgun.Mailgun
var t *template.Template

const (
	WelcomeMail    = "welcome.template.html"
	NewsletterMail = "newsletter.template.html"
	NewPostMail    = "new-post.template.html"
)

type MailContent map[string]interface{}

func initMail() {
	domain = Env("MAILGUN_DOMAIN", "")
	sender = Env("MAILGUN_SENDER", "") + "@" + domain

	pubKey := Env("MAILGUN_KEY_PUBLIC", "")
	secKey := Env("MAILGUN_KEY_SECRET", "")
	mg = mailgun.NewMailgun(domain, secKey, pubKey)

	t = template.Must(template.ParseGlob("templates/*.template.html"))
	Log.Infof("Templates loaded%v", t.DefinedTemplates())
}

func SendMail(
	content MailContent,
	template, subject string,
	recipients ...string,
) error {
	b := new(bytes.Buffer)
	if err := t.ExecuteTemplate(b, template, content); err != nil {
		return err
	}

	msg := mailgun.NewMessage(
		sender,
		subject,
		b.String(),
	)

	for _, r := range recipients {
		msg.AddRecipient(r)
	}

	// res, id, err := mg.Send(msg)
	// Log.Debugf("res: %v, id: %v, err: %v", res, id, err)

	return nil
}
