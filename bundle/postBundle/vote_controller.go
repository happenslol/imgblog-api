package postBundle

import (
	"errors"

	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"

	"github.com/happeens/imgblog-api/app"
	"github.com/happeens/imgblog-api/model"
)

type voteController struct{}

func (voteController) ShowUserVotes(c *gin.Context) {
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

func (voteController) ShowPostVotes(c *gin.Context) {
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

func (voteController) ShowCommentVotes(c *gin.Context) {
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

func (voteController) CreatePostVote(c *gin.Context) {
	userName, _ := c.Get("user")
	user := model.User{}
	if err := app.DB().C(model.UserC).Find(
		bson.M{"name": userName},
	).One(&user); err != nil {
		app.DbError(c, err)
		return
	}

	parentID := bson.ObjectIdHex(c.Param("id"))

	if existingCount, _ := app.DB().C(model.VoteC).Find(
		bson.M{
			"parentType": model.PostVote,
			"parentId":   parentID,
			"userId":     user.ID,
		},
	).Count(); existingCount > 0 {
		app.BadRequest(c, errors.New("you can only vote once!"))
		return
	}

	insert := model.Vote{
		ID:         bson.NewObjectId(),
		ParentType: model.PostVote,
		ParentID:   parentID,
		UserID:     user.ID,
	}

	if err := app.DB().C(model.VoteC).Insert(
		&insert,
	); err != nil {
		app.DbError(c, err)
		return
	}

	if err := app.DB().C(model.PostC).Update(
		bson.M{"_id": parentID},
		bson.M{"$inc": bson.M{
			"upvotes": 1,
		}},
	); err != nil {
		app.DbError(c, err)
		return
	}

	app.Ok(c, insert)
}

func (voteController) DeletePostVote(c *gin.Context) {
	userName, _ := c.Get("user")
	user := model.User{}
	if err := app.DB().C(model.UserC).Find(
		bson.M{"name": userName},
	).One(&user); err != nil {
		app.DbError(c, err)
		return
	}

	parentID := bson.ObjectIdHex(c.Param("id"))
	if err := app.DB().C(model.VoteC).Remove(
		bson.M{
			"parentType": model.PostVote,
			"parentId":   parentID,
			"userId":     user.ID,
		},
	); err != nil {
		app.DbError(c, err)
		return
	}

	if err := app.DB().C(model.PostC).Update(
		bson.M{"_id": parentID},
		bson.M{"$inc": bson.M{
			"upvotes": -1,
		}},
	); err != nil {
		app.DbError(c, err)
		return
	}

	app.Ok(c, gin.H{"deleted": 1})
}

func (voteController) CreateCommentVote(c *gin.Context) {
	userName, _ := c.Get("user")
	user := model.User{}
	if err := app.DB().C(model.UserC).Find(
		bson.M{"name": userName},
	).One(&user); err != nil {
		app.DbError(c, err)
		return
	}

	parentID := bson.ObjectIdHex(c.Param("id"))

	if existingCount, _ := app.DB().C(model.VoteC).Find(
		bson.M{
			"parentType": model.CommentVote,
			"parentId":   parentID,
			"userId":     user.ID,
		},
	).Count(); existingCount > 0 {
		app.BadRequest(c, errors.New("you can only vote once!"))
		return
	}

	insert := model.Vote{
		ID:         bson.NewObjectId(),
		ParentType: model.CommentVote,
		ParentID:   parentID,
		UserID:     user.ID,
	}

	if err := app.DB().C(model.VoteC).Insert(
		&insert,
	); err != nil {
		app.DbError(c, err)
		return
	}

	if err := app.DB().C(model.PostC).Update(
		bson.M{"comments._id": parentID},
		bson.M{"$inc": bson.M{
			"comments.$.upvotes": 1,
		}},
	); err != nil {
		app.DbError(c, err)
		return
	}

	app.Ok(c, insert)
}

func (voteController) DeleteCommentVote(c *gin.Context) {
	userName, _ := c.Get("user")
	user := model.User{}
	if err := app.DB().C(model.UserC).Find(
		bson.M{"name": userName},
	).One(&user); err != nil {
		app.DbError(c, err)
		return
	}

	parentID := bson.ObjectIdHex(c.Param("id"))
	if err := app.DB().C(model.VoteC).Remove(
		bson.M{
			"parentType": model.CommentVote,
			"parentId":   parentID,
			"userId":     user.ID,
		},
	); err != nil {
		app.DbError(c, err)
		return
	}

	if err := app.DB().C(model.PostC).Update(
		bson.M{"comments._id": parentID},
		bson.M{"$inc": bson.M{
			"comments.$.upvotes": -1,
		}},
	); err != nil {
		app.DbError(c, err)
		return
	}

	app.Ok(c, gin.H{"deleted": 1})
}
