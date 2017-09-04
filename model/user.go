package model

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const UserC = "users"

const (
	AdminRole = "admin"
	UserRole  = "user"
)

type User struct {
	ID       bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Name     string        `json:"name"`
	Password string        `json:"-"`
	Email    string        `json:"email"`
	Role     string        `json:"role"`
}

func (u User) ToPartial() UserPartial {
	return UserPartial{
		ID:   u.ID,
		Name: u.Name,
	}
}

type UserPartial struct {
	ID   bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Name string        `json:"name"`
}

func EnsureUserIndices(db *mgo.Database) error {
	nameIndex := mgo.Index{
		Key:        []string{"name"},
		Unique:     true,
		DropDups:   false,
		Background: true,
	}

	emailIndex := mgo.Index{
		Key:        []string{"email"},
		Unique:     true,
		DropDups:   false,
		Background: true,
	}

	err := db.C(UserC).EnsureIndex(nameIndex)
	if err != nil {
		return err
	}

	err = db.C(UserC).EnsureIndex(emailIndex)
	if err != nil {
		return err
	}

	return nil
}
