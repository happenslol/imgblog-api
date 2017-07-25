package postBundle

import (
	"errors"

	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"

	"github.com/happeens/imgblog-api/app"
	"github.com/happeens/imgblog-api/model"
)

type voteController struct{}

//TODO get this to work
func (voteController) showUserVotes(c *gin.Context) {
	var result []model.Vote
	err := app.DB().C(model.VoteC).Find(
		bson.M{"user._id": bson.ObjectIdHex(c.Param("id"))},
	).All(&result)

	if err != nil {
		app.DbError(c, err)
		return
	}

	app.Ok(c, result)
}

//TODO get this to work
func (voteController) showPostVotes(c *gin.Context) {
	var result []model.Vote
	err := app.DB().C(model.VoteC).Find(
		bson.M{
			"parentType": model.PostVote,
			"parentId":   bson.ObjectIdHex(c.Param("id")),
		},
	).All(&result)

	if err != nil {
		app.DbError(c, err)
		return
	}

	app.Ok(c, result)
}

func (voteController) showCommentVotes(c *gin.Context) {
	var result []model.Vote
	err := app.DB().C(model.VoteC).Find(
		bson.M{
			"parentType": model.CommentVote,
			"parentId":   bson.ObjectIdHex(c.Param("id")),
		},
	).All(&result)

	if err != nil {
		app.DbError(c, err)
		return
	}

	app.Ok(c, result)
}

type createVoteRequest struct {
	VoteType string `json:"voteType" binding:"required,eq=up|eq=down"`
}

func (voteController) createPostVote(c *gin.Context) {
	var json createVoteRequest
	err := c.BindJSON(&json)
	if err != nil {
		app.BadRequest(c, errors.New("Vote type must be 'up' or 'down'"))
		return
	}

	userName, _ := c.Get("user")
	user := model.User{}
	err = app.DB().C(model.UserC).Find(bson.M{"name": userName}).One(&user)
	if err != nil {
		app.DbError(c, err)
		return
	}

	upvoteAmount := 0
	downvoteAmount := 0

	if json.VoteType == "up" {
		upvoteAmount = 1
	}

	if json.VoteType == "down" {
		downvoteAmount = 1
	}

	parentID := bson.ObjectIdHex(c.Param("id"))

	insert := model.Vote{
		ID:         bson.NewObjectId(),
		VoteType:   json.VoteType,
		ParentType: model.PostVote,
		ParentID:   parentID,
		UserID:     user.ID,
	}

	//TODO handle double votes
	err = app.DB().C(model.VoteC).Insert(&insert)
	if err != nil {
		app.DbError(c, err)
		return
	}

	err = app.DB().C(model.PostC).Update(
		bson.M{"_id": parentID},
		bson.M{"$inc": bson.M{
			"upvotes":   upvoteAmount,
			"downvotes": downvoteAmount,
		}},
	)
	if err != nil {
		app.DbError(c, err)
		return
	}

	app.Created(c, insert.ID)
}

func (voteController) createCommentVote(c *gin.Context) {
	var json createVoteRequest
	err := c.BindJSON(&json)
	if err != nil {
		app.BadRequest(c, errors.New("Vote type must be 'up' or 'down'"))
		return
	}

	userName, _ := c.Get("user")
	user := model.User{}
	err = app.DB().C(model.UserC).Find(bson.M{"name": userName}).One(&user)
	if err != nil {
		app.DbError(c, err)
		return
	}

	upvoteAmount := 0
	downvoteAmount := 0

	if json.VoteType == "up" {
		upvoteAmount = 1
	}

	if json.VoteType == "down" {
		downvoteAmount = 1
	}

	parentID := bson.ObjectIdHex(c.Param("id"))

	insert := model.Vote{
		ID:         bson.NewObjectId(),
		VoteType:   json.VoteType,
		ParentType: model.CommentVote,
		ParentID:   parentID,
		UserID:     user.ID,
	}

	//TODO handle double votes
	err = app.DB().C(model.VoteC).Insert(&insert)
	if err != nil {
		app.DbError(c, err)
		return
	}

	err = app.DB().C(model.PostC).Update(
		bson.M{"comments._id": parentID},
		bson.M{"$inc": bson.M{
			"comments.$.upvotes":   upvoteAmount,
			"comments.$.downvotes": downvoteAmount,
		}},
	)
	if err != nil {
		app.DbError(c, err)
		return
	}

	app.Created(c, insert.ID)
}
