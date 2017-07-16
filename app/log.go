package app

import (
	"github.com/op/go-logging"
	"os"
)

var Log = logging.MustGetLogger("app")
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)

func initLogger() {
	logBackend := logging.NewLogBackend(os.Stdout, "", 0)
	backendFormatted := logging.NewBackendFormatter(logBackend, format)
	logging.SetBackend(backendFormatted)
}
