package app

import (
	"os"
)

var storagePath string

func initStorage() {
	storagePath = Env("STORAGE", "")

	if _, err := os.Stat(storagePath); os.IsNotExist(err) {
		panic("storage path not found: " + storagePath)
	}

	Router.Static("/static", storagePath)
}
