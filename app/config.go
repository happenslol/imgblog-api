package app

import (
	"os"

	"github.com/joho/godotenv"
)

func initConf() {
	err := godotenv.Load()
	if err != nil {
		Log.Criticalf("error loading env: %v", err)
	}
}

func Env(key string, def string) string {
	val := os.Getenv(key)
	if val == "" {
		return def
	}

	return val
}
