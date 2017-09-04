package app

import (
	"gopkg.in/mgo.v2"

	"github.com/happeens/imgblog-api/model"
)

var db *mgo.Database

func initDb() {
	// get database config
	host := Env("DB_HOST", "localhost")
	port := Env("DB_PORT", "27017")
	user := Env("DB_USERNAME", "")
	pass := Env("DB_PASSWORD", "")
	name := Env("DB_DATABASE", "")

	// construct database dial string
	dialString := "mongodb://" +
		user + ":" +
		pass + "@" +
		host + ":" +
		port + "/" +
		name

	con, err := mgo.Dial(dialString)
	if err != nil {
		panic(err)
	}

	db = con.DB(name)

	// Ensure all indices
	if err = model.EnsureUserIndices(db); err != nil {
		panic(err)
	}

	if err = model.EnsureVoteIndices(db); err != nil {
		panic(err)
	}

	if err = model.EnsurePostIndices(db); err != nil {
		panic(err)
	}
}

func DB() *mgo.Database {
	return db
}
