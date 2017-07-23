package postBundle

import (
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"

	"github.com/happeens/imgblog-api/app"
	"github.com/happeens/imgblog-api/model"
)

const postsPageSize = 2

type postController struct{}

func (postController) Index(c *gin.Context) {
	var result []model.Post

	pageString := c.DefaultQuery("page", "all")
	if pageString == "all" {
		err := app.DB().C(model.PostC).Find(nil).All(&result)
		if err != nil {
			app.DbError(c, err)
			return
		}

		app.Ok(c, result)
		return
	}

	page, err := strconv.Atoi(pageString)
	if err != nil {
		app.BadRequest(c, err)
		return
	}

	err = app.DB().C(model.PostC).Find(nil).Skip(
		page * postsPageSize,
	).Limit(postsPageSize).All(&result)

	if err != nil {
		app.DbError(c, err)
		return
	}

	app.Ok(c, result)
}

func (postController) Show(c *gin.Context) {
	var result model.Post
	err := app.DB().C(model.PostC).Find(bson.M{"slug": c.Param("slug")}).One(&result)
	if err != nil {
		app.DbError(c, err)
	}

	app.Ok(c, result)
}

type createRequest struct {
	Title      app.LocalString `json:"title" binding:"required"`
	Content    app.LocalString `json:"content" binding:"required"`
	TitleImage string
	Images     []string
}

func (postController) Create(c *gin.Context) {
	var json createRequest
	err := c.BindJSON(&json)
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

	slug := strings.Replace(strings.ToLower(json.Title["en"]), " ", "-", -1)

	insert := model.Post{
		ID:         bson.NewObjectId(),
		Author:     user.ToPartial(),
		Title:      json.Title,
		Slug:       slug,
		TitleImage: json.TitleImage,
		Content:    json.Content,
		Images:     json.Images,
		Comments:   []model.Comment{},
		Created:    time.Now(),
		Updated:    nil,
		Deleted:    nil,
	}

	err = app.DB().C(model.PostC).Insert(&insert)
	if err != nil {
		app.DbError(c, err)
		return
	}

	app.Created(c, insert.ID)
}

type createCommentRequest struct {
	Content string
}

func (postController) CreateComment(c *gin.Context) {
	var json createCommentRequest
	err := c.BindJSON(&json)
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
		ID:      bson.NewObjectId(),
		Author:  user.ToPartial(),
		Content: json.Content,
		Created: time.Now(),
		Updated: nil,
		Deleted: nil,
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

func (postController) Destroy(c *gin.Context) {
	id := bson.ObjectIdHex(c.Param("id"))
	err := app.DB().C(model.PostC).RemoveId(id)
	if err != nil {
		app.DbError(c, err)
		return
	}

	app.Ok(c, gin.H{"deleted": 1})
}
