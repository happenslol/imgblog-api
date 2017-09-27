package model

import (
	"encoding/json"
	"errors"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const PostC = "posts"

type Post struct {
	ID         bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Author     UserPartial   `json:"author"`
	Title      LocalString   `json:"title"`
	Intro      LocalString   `json:"intro"`
	Slug       string        `json:"slug"`
	TitleImage string        `bson:"titleImage" json:"titleImage"`
	Sections   []PostSection `json:"sections"`
	Comments   []Comment     `json:"comments"`

	Category string   `json:"category"`
	Tags     []string `json:"tags"`

	Upvotes int `json:"upvotes"`

	Created time.Time  `json:"created"`
	Updated *time.Time `json:"updated"`
	Deleted *time.Time `json:"deleted"`
}

type PostSection struct {
	Content SectionContent
}

type SectionContent interface {
	Type() string
	Content() interface{}
}

func (p *PostSection) MarshalJSON() (b []byte, e error) {
	return json.Marshal(map[string]interface{}{
		"type":    p.Content.Type(),
		"content": p.Content.Content(),
	})
}

func (p *PostSection) UnmarshalJSON(b []byte) error {
	decoded := new(struct {
		Type    string      `json:"type" bson:"type"`
		Content interface{} `json:"content" bson:"content"`
	})

	if err := json.Unmarshal(b, &decoded); err != nil {
		return err
	}

	switch decoded.Type {
	case "text":
		content := decoded.Content.(map[string]interface{})
		langs := make(PostText, len(content))
		for lang, c := range content {
			s := c.(string)
			langs[lang] = s
		}
		p.Content = langs
	case "image":
		content := decoded.Content.(string)
		p.Content = PostImage(content)
	default:
		return errors.New("unrecognized section type")
	}

	return nil
}

func (p PostSection) GetBSON() (interface{}, error) {
	return struct {
		Type    string      `json:"type" bson:"type"`
		Content interface{} `json:"content" bson:"content"`
	}{
		Type:    p.Content.Type(),
		Content: p.Content.Content(),
	}, nil
}

func (p *PostSection) SetBSON(raw bson.Raw) error {
	decoded := new(struct {
		Type    string      `json:"type" bson:"type"`
		Content interface{} `json:"content" bson:"content"`
	})

	if err := raw.Unmarshal(decoded); err != nil {
		return err
	}

	switch decoded.Type {
	case "text":
		content := decoded.Content.(bson.M)
		langs := make(PostText, len(content))
		for lang, c := range content {
			s := c.(string)
			langs[lang] = s
		}
		p.Content = langs
	case "image":
		content := decoded.Content.(string)
		p.Content = PostImage(content)
	default:
		return errors.New("unrecognized section type")
	}

	return nil
}

type PostText LocalString

func (PostText) Type() string {
	return "text"
}

func (p PostText) Content() interface{} {
	return p
}

type PostImage string

func (PostImage) Type() string {
	return "image"
}

func (p PostImage) Content() interface{} {
	return p
}

type Comment struct {
	ID       bson.ObjectId  `bson:"_id,omitempty" json:"id"`
	ParentID *bson.ObjectId `bson:"parentId,omitempty" json:"parentId,omitempty"`
	Author   UserPartial    `json:"author"`
	Content  string         `json:"content"`

	Upvotes int `json:"upvotes"`

	Created time.Time  `json:"created"`
	Updated *time.Time `json:"updated"`
	Deleted *time.Time `json:"deleted"`
}

func EnsurePostIndices(db *mgo.Database) error {
	slugIndex := mgo.Index{
		Key:        []string{"slug"},
		Unique:     true,
		DropDups:   false,
		Background: true,
	}

	return db.C(PostC).EnsureIndex(slugIndex)
}
