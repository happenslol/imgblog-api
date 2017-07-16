package app

import (
	"github.com/joho/godotenv"
)

var env map[string]string

func initConf() {
	var err error
	env, err = godotenv.Read()
	if err != nil {
		Log.Criticalf("error loading env: %v", err)
	}
}

func Env(key string, def string) string {
	val := env[key]
	if val == "" {
		return def
	}

	return val
}
