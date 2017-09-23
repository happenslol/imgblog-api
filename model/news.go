package model

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

const NewsC = "news"

type News struct {
	ID       bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Author   UserPartial   `json:"author"`
	Content  LocalString   `json:"content"`
	Image    string        `json:"image,omitempty" bson:",omitempty"`
	Category string        `json:"category,omitempty" bson:",omitempty"`

	Created time.Time  `json:"created"`
	Updated *time.Time `json:"updated"`
	Deleted *time.Time `json:"deleted"`
}
