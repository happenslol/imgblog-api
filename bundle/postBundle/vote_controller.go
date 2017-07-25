package postBundle

import (
	"errors"

	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"

	"github.com/happeens/imgblog-api/app"
	"github.com/happeens/imgblog-api/model"
)

type voteController struct{}

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
	var response gin.H

	existingCount, _ := app.DB().C(model.VoteC).Find(
		bson.M{
			"parentType": model.PostVote,
			"parentId":   parentID,
			"userId":     user.ID,
		},
	).Count()

	// If there already is a vote
	if existingCount > 0 {
		existingVote := model.Vote{}
		err = app.DB().C(model.VoteC).Find(
			bson.M{
				"parentType": model.PostVote,
				"parentId":   parentID,
				"userId":     user.ID,
			},
		).One(&existingVote)

		if err != nil {
			app.DbError(c, err)
			return
		}

		// Find it, check the vote type and adjust the values accordingly
		if existingVote.VoteType == json.VoteType {
			// Nothing happens at all in this case and we can instantly return
			app.Ok(c, gin.H{"response": "nothing updated, vote already exists"})
			return
		}

		// Otherwise, flip the amounts
		if upvoteAmount == 0 {
			upvoteAmount = -1
		}

		if downvoteAmount == 0 {
			downvoteAmount = -1
		}

		// Update Database entry if they were different
		err = app.DB().C(model.VoteC).Update(
			bson.M{"_id": existingVote.ID},
			bson.M{"$set": bson.M{"voteType": json.VoteType}},
		)

		if err != nil {
			app.DbError(c, err)
			return
		}

		response = gin.H{"updated": existingVote.ID}
	} else {
		// Insert new vote otherwise
		insert := model.Vote{
			ID:         bson.NewObjectId(),
			VoteType:   json.VoteType,
			ParentType: model.PostVote,
			ParentID:   parentID,
			UserID:     user.ID,
		}

		err = app.DB().C(model.VoteC).Insert(&insert)
		if err != nil {
			app.DbError(c, err)
			return
		}

		response = gin.H{"created": insert.ID}
	}

	// Update post
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

	//TODO find a better solution for response codes
	app.Ok(c, response)
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
	var response gin.H

	existingCount, _ := app.DB().C(model.VoteC).Find(
		bson.M{
			"parentType": model.CommentVote,
			"parentId":   parentID,
			"userId":     user.ID,
		},
	).Count()

	// If there already is a vote
	if existingCount > 0 {
		existingVote := model.Vote{}
		err = app.DB().C(model.VoteC).Find(
			bson.M{
				"parentType": model.CommentVote,
				"parentId":   parentID,
				"userId":     user.ID,
			},
		).One(&existingVote)

		if err != nil {
			app.DbError(c, err)
			return
		}

		// Find it, check the vote type and adjust the values accordingly
		if existingVote.VoteType == json.VoteType {
			// Nothing happens at all in this case and we can instantly return
			app.Ok(c, gin.H{"response": "nothing updated, vote already exists"})
			return
		}

		// Otherwise, flip the amounts
		if upvoteAmount == 0 {
			upvoteAmount = -1
		}

		if downvoteAmount == 0 {
			downvoteAmount = -1
		}

		// Update Database entry if they were different
		err = app.DB().C(model.VoteC).Update(
			bson.M{"_id": existingVote.ID},
			bson.M{"$set": bson.M{"voteType": json.VoteType}},
		)

		if err != nil {
			app.DbError(c, err)
			return
		}

		response = gin.H{"updated": existingVote.ID}
	} else {
		// Insert new vote otherwise
		insert := model.Vote{
			ID:         bson.NewObjectId(),
			VoteType:   json.VoteType,
			ParentType: model.CommentVote,
			ParentID:   parentID,
			UserID:     user.ID,
		}

		err = app.DB().C(model.VoteC).Insert(&insert)
		if err != nil {
			app.DbError(c, err)
			return
		}

		response = gin.H{"created": insert.ID}
	}

	// Update comment
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

	//TODO find a better solution for response codes
	app.Ok(c, response)
}
