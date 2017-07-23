package app

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var Router *gin.Engine

func init() {
	initLogger()

	Log.Debug("Loading environment vars...")
	initConf()

	Log.Debug("Connecting to database...")
	initDb()

	Log.Debug("Setting up auth...")
	initAuth()

	if Env("ENV", "dev") == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	Router = gin.Default()
	Router.Use(cors.Default())
}
