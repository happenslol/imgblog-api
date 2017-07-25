package model

import (
	"gopkg.in/mgo.v2/bson"
)

const VoteC = "votes"

const (
	UpVoteType   = "up"
	DownVoteType = "down"
)

const (
	PostVote    = "post"
	CommentVote = "comment"
)

type Vote struct {
	ID         bson.ObjectId `bson:"_id,omitempty" json:"id"`
	VoteType   string        `bson:"voteType" json:"voteType"`
	ParentType string        `bson:"parentType" json:"parentType"`
	ParentID   bson.ObjectId `bson:"parentId" json:"parentId"`
	UserID     bson.ObjectId `bson:"userId" json:"userId"`
}
