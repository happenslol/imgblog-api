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
	VoteType   string        `json:"voteType"`
	ParentType string        `json:"parentType"`
	ParentID   bson.ObjectId `json:"parentId"`
	User       UserPartial   `json:"user"`
}
