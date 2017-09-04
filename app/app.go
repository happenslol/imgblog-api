package app

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var Router *gin.Engine

func init() {
	initLogger()

	Log.Info("Loading environment vars")
	initConf()

	Log.Info("Connecting to database")
	initDb()

	Log.Info("Setting up auth")
	initAuth()

	Log.Info("Initializing mailgun api")
	initMail()

	Log.Info("Initializing captcha")
	initCaptcha()

	Log.Info("Checking storage")
	initStorage()

	if Env("ENV", "dev") == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	Router = gin.Default()
	Router.Use(cors.Default())
}
