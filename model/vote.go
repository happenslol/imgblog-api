package model

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const VoteC = "votes"

const (
	PostVote    = "post"
	CommentVote = "comment"
)

type Vote struct {
	ID         bson.ObjectId `bson:"_id,omitempty" json:"id"`
	ParentType string        `bson:"parentType" json:"parentType"`
	ParentID   bson.ObjectId `bson:"parentId" json:"parentId"`
	UserID     bson.ObjectId `bson:"userId" json:"userId"`
}

func EnsureVoteIndices(db *mgo.Database) error {
	voteIndex := mgo.Index{
		Key:        []string{"parentId", "userId"},
		Unique:     true,
		DropDups:   true,
		Background: true,
	}

	err := db.C(VoteC).EnsureIndex(voteIndex)
	if err != nil {
		return err
	}

	return nil
}
