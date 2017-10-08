package model

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

const DraftC = "drafts"

type Draft struct {
	ID         bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Author     UserPartial   `json:"author"`
	Title      LocalString   `json:"title"`
	Intro      LocalString   `json:"intro"`
	TitleImage string        `bson:"titleImage" json:"titleImage"`
	Sections   []PostSection `json:"sections"`

	Category string   `json:"category"`
	Tags     []string `json:"tags"`

	Created time.Time  `json:"created"`
	Updated *time.Time `json:"updated"`
	Deleted *time.Time `json:"deleted"`
}
