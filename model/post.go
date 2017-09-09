package model

import (
	"encoding/json"
	"fmt"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const PostC = "posts"

type Post struct {
	ID         bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Author     UserPartial   `json:"author"`
	Title      LocalString   `json:"title"`
	Slug       string        `json:"slug"`
	TitleImage string        `bson:"titleImage" json:"titleImage"`
	Sections   []PostSection `json:"sections"`
	Images     []string      `json:"images"`
	Comments   []Comment     `json:"comments"`

	Upvotes   int `json:"upvotes"`
	Downvotes int `json:"downvotes"`

	Created time.Time  `json:"created"`
	Updated *time.Time `json:"updated"`
	Deleted *time.Time `json:"deleted"`
}

func (p *Post) SetBSON(raw bson.Raw) error {
	decoded := new(struct {
		ID         bson.ObjectId            `bson:"_id,omitempty" json:"id"`
		Author     UserPartial              `json:"author"`
		Title      LocalString              `json:"title"`
		Slug       string                   `json:"slug"`
		TitleImage string                   `bson:"titleImage" json:"titleImage"`
		Sections   []map[string]interface{} `json:"sections"`
		Images     []string                 `json:"images"`
		Comments   []Comment                `json:"comments"`

		Upvotes   int `json:"upvotes"`
		Downvotes int `json:"downvotes"`

		Created time.Time  `json:"created"`
		Updated *time.Time `json:"updated"`
		Deleted *time.Time `json:"deleted"`
	})

	err := raw.Unmarshal(decoded)
	if err != nil {
		return err
	}

	p.ID = decoded.ID
	p.Author = decoded.Author
	p.Title = decoded.Title
	p.Slug = decoded.Slug
	p.TitleImage = decoded.TitleImage
	p.Images = decoded.Images
	p.Comments = decoded.Comments
	p.Upvotes = decoded.Upvotes
	p.Downvotes = decoded.Downvotes
	p.Created = decoded.Created
	p.Updated = decoded.Updated
	p.Deleted = decoded.Deleted

	sections := make([]PostSection, len(decoded.Sections))
	for i, s := range decoded.Sections {
		if s["type"] == "text" {
			content := s["content"].(map[string]interface{})
			langs := make(PostText, len(content))
			for lang, langContent := range content {
				langString := langContent.(string)
				langs[lang] = langString
			}
			sections[i] = &langs
		} else if s["type"] == "image" {
			content := s["content"].(string)
			contentString := PostImage(content)
			sections[i] = &contentString
		}
	}

	p.Sections = sections

	fmt.Printf("decoded: %v\n", decoded)

	return nil
}

type PostSection interface {
	Type() interface{}
	Content() interface{}
}

type PostText LocalString

func (PostText) Type() interface{} {
	return "text"
}

func (p PostText) Content() interface{} {
	return p
}

func (p *PostText) MarshalJSON() (b []byte, e error) {
	return json.Marshal(map[string]interface{}{
		"type":    p.Type(),
		"content": p.Content(),
	})
}

func (p *PostText) GetBSON() (interface{}, error) {
	return struct {
		Type    string      `json:"type" bson:"type"`
		Content interface{} `json:"content" bson:"content"`
	}{
		Type:    p.Type().(string),
		Content: p.Content().(PostText),
	}, nil
}

func (p *PostText) SetBSON(raw bson.Raw) error {
	decoded := new(struct {
		Type    string      `json:"type" bson:"type"`
		Content interface{} `json:"content" bson:"content"`
	})

	bsonErr := raw.Unmarshal(decoded)

	if bsonErr != nil {
		return bsonErr
	}

	pt := decoded.Content.(PostText)
	p = &pt
	return nil
}

type PostImage string

func (PostImage) Type() interface{} {
	return "image"
}

func (p PostImage) Content() interface{} {
	return p
}

func (p *PostImage) MarshalJSON() (b []byte, e error) {
	return json.Marshal(map[string]interface{}{
		"type":    p.Type(),
		"content": p.Content(),
	})
}

func (p *PostImage) GetBSON() (interface{}, error) {
	return struct {
		Type    string      `json:"type" bson:"type"`
		Content interface{} `json:"content" bson:"content"`
	}{
		Type:    p.Type().(string),
		Content: p.Content().(PostImage),
	}, nil
}

func (p *PostImage) SetBSON(raw bson.Raw) error {
	decoded := new(struct {
		Type    string `json:"type" bson:"type"`
		Content string `json:"content" bson:"content"`
	})

	bsonErr := raw.Unmarshal(decoded)

	if bsonErr != nil {
		return bsonErr
	}

	pi := PostImage(decoded.Content)
	p = &pi

	return nil
}

type Comment struct {
	ID       bson.ObjectId  `bson:"_id,omitempty" json:"id"`
	ParentID *bson.ObjectId `bson:",omitempty" json:"parentId"`
	Author   UserPartial    `json:"author"`
	Content  string         `json:"content"`

	Upvotes   int `json:"upvotes"`
	Downvotes int `json:"downvotes"`

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
