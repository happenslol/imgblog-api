package model

import (
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/happeens/imgblog-api/app"
)

const PostC = "posts"

type Post struct {
	ID         bson.ObjectId   `bson:"_id,omitempty" json:"id"`
	Author     UserPartial     `json:"author"`
	Title      app.LocalString `json:"title"`
	Slug       string          `json:"slug"`
	TitleImage string          `json:"titleImage"`
	Content    app.LocalString `json:"content"`
	Images     []string        `json:"images"`
	Comments   []Comment       `json:"comments"`
	Created    time.Time       `json:"created"`
	Updated    *time.Time      `json:"updated"`
	Deleted    *time.Time      `json:"deleted"`
}

type Comment struct {
	ID      bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Author  UserPartial   `json:"author"`
	Content string        `json:"content"`
	Created time.Time     `json:"created"`
	Updated *time.Time    `json:"updated"`
	Deleted *time.Time    `json:"deleted"`
}
