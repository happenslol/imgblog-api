package model

import (
	"gopkg.in/mgo.v2/bson"
)

const UserC = "users"

const (
	AdminRole = "admin"
	UserRole  = "user"
)

type User struct {
	ID       bson.ObjectId `bson:"_id,omitempty"`
	Name     string
	Password string `json:"-"`
	Email    string
	Role     string
}
