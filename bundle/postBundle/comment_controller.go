package postBundle

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kennygrant/sanitize"
	"gopkg.in/mgo.v2/bson"

	"github.com/happeens/imgblog-api/app"
	"github.com/happeens/imgblog-api/model"
)

type createCommentRequest struct {
	ParentID string `json:"parentId"`
	Content  string `json:"content" binding:"required"`
}

func (postController) CreateComment(c *gin.Context) {
	var req createCommentRequest
	err := c.BindJSON(&req)
	if err != nil {
		app.BadRequest(c, err)
		return
	}

	userName, _ := c.Get("user")
	user := model.User{}
	err = app.DB().C(model.UserC).Find(bson.M{"name": userName}).One(&user)
	if err != nil {
		app.DbError(c, err)
		return
	}

	newComment := model.Comment{
		ID:       bson.NewObjectId(),
		Author:   user.ToPartial(),
		ParentID: nil,
		Content:  sanitize.HTML(req.Content),

		Upvotes: 0,

		Created: time.Now(),
		Updated: nil,
		Deleted: nil,
	}

	if req.ParentID != "" {
		parentID := bson.ObjectIdHex(req.ParentID)
		var count int
		count, err = app.DB().C(model.PostC).FindId(parentID).Count()
		if count != 1 {
			//TODO tell user that parent was not found
			app.BadRequest(c, err)
			return
		}

		newComment.ParentID = &parentID
	}

	postID := bson.ObjectIdHex(c.Param("id"))
	err = app.DB().C(model.PostC).Update(
		bson.M{"_id": postID},
		bson.M{"$push": bson.M{"comments": newComment}},
	)

	if err != nil {
		app.DbError(c, err)
		return
	}

	app.Created(c, newComment.ID)
}

func (postController) UpdateComment(c *gin.Context) {
	var req createCommentRequest
	err := c.BindJSON(&req)
	if err != nil {
		app.BadRequest(c, err)
		return
	}

	user := model.User{}
	userName, _ := c.Get("user")
	err = app.DB().C(model.UserC).Find(bson.M{"name": userName}).One(&user)
	if err != nil {
		app.DbError(c, err)
		return
	}

	var postExists int
	postExists, err = app.DB().C(model.PostC).Find(
		bson.M{
			"_id":               bson.ObjectIdHex(c.Param("id")),
			"comments._id":      bson.ObjectIdHex(c.Param("commentId")),
			"comments.user._id": user.ID,
		},
	).Count()

	if err != nil {
		app.DbError(c, err)
		return
	}

	if postExists == 0 {
		//TODO better error reporting
		app.DbError(c, errors.New("comment not found"))
		return
	}

	err = app.DB().C(model.PostC).Update(
		bson.M{
			"_id":               bson.ObjectIdHex(c.Param("id")),
			"comments._id":      bson.ObjectIdHex(c.Param("commentId")),
			"comments.user._id": user.ID,
		},
		bson.M{"$set": bson.M{
			"comments.$.content": req.Content,
			"comments.$.updated": time.Now(),
		}},
	)

	if err != nil {
		app.DbError(c, err)
		return
	}

	app.Ok(c, gin.H{"updated": c.Param("id")})
}

func (postController) DestroyComment(c *gin.Context) {
	user := model.User{}
	userName, _ := c.Get("user")
	err := app.DB().C(model.UserC).Find(bson.M{"name": userName}).One(&user)
	if err != nil {
		app.DbError(c, err)
		return
	}

	var postExists int
	postExists, err = app.DB().C(model.PostC).Find(
		bson.M{
			"_id":               bson.ObjectIdHex(c.Param("id")),
			"comments._id":      bson.ObjectIdHex(c.Param("commentId")),
			"comments.user._id": user.ID,
		},
	).Count()

	if err != nil {
		app.DbError(c, err)
		return
	}

	if postExists == 0 {
		//TODO better error reporting
		app.DbError(c, errors.New("comment not found"))
		return
	}

	err = app.DB().C(model.PostC).Update(
		bson.M{
			"_id":               bson.ObjectIdHex(c.Param("id")),
			"comments._id":      bson.ObjectIdHex(c.Param("commentId")),
			"comments.user._id": user.ID,
		},
		bson.M{"$set": bson.M{"comments.$.deleted": time.Now()}},
	)

	if err != nil {
		app.DbError(c, err)
		return
	}

	app.Ok(c, gin.H{"deleted": 1})
}
