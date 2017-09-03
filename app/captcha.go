package app

import (
	"github.com/dpapathanasiou/go-recaptcha"
)

func initCaptcha() {
	recaptcha.Init(Env("CAPTCHA", ""))
}

func ConfirmCaptcha(ip, res string) bool {
	if Env("ENV", "dev") == "dev" {
		return true
	}

	return recaptcha.Confirm(ip, res)
}
