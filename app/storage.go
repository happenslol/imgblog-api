package app

import (
	"os"
)

var StoragePath string

func initStorage() {
	StoragePath = Env("STORAGE", "")

	if _, err := os.Stat(StoragePath); os.IsNotExist(err) {
		panic("storage path not found: " + StoragePath)
	}

	Router.Static("/static", StoragePath)
}
