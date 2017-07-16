package model

import (
	"gopkg.in/mgo.v2/bson"
)

const PostC = "posts"

type Post struct {
	ID      bson.ObjectId `bson:"_id,omitempty"`
	Title   string
	Content string
}
