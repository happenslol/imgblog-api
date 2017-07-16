package app

import (
	"gopkg.in/mgo.v2"
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
}

func DB() *mgo.Database {
	return db
}
