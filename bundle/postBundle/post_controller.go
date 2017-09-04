package postBundle

import (
	"bytes"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kennygrant/sanitize"
	"gopkg.in/mgo.v2/bson"

	"github.com/happeens/imgblog-api/app"
	"github.com/happeens/imgblog-api/model"
)

const postsPageSize = 2

type postController struct{}

func (postController) index(c *gin.Context) {
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

func (postController) show(c *gin.Context) {
	var result model.Post
	err := app.DB().C(model.PostC).Find(
		bson.M{"slug": c.Param("slug")},
	).One(&result)

	if err != nil {
		app.DbError(c, err)
	}

	app.Ok(c, result)
}

func (postController) search(c *gin.Context) {
	var results []model.Post
	query := c.DefaultQuery("query", "")
	if query == "" {
		app.BadRequest(c, errors.New("empty search query not allowed"))
		return
	}

	err := app.DB().C(model.PostC).Find(
		bson.M{"$or": []bson.M{
			bson.M{"title.en": bson.RegEx{Pattern: query, Options: "i"}},
			bson.M{"title.de": bson.RegEx{Pattern: query, Options: "i"}},
		}},
	).All(&results)

	if err != nil {
		app.DbError(c, err)
	}

	app.Ok(c, results)
}

type createRequest struct {
	Title      model.LocalString `json:"title" binding:"required"`
	Content    model.LocalString `json:"content" binding:"required"`
	TitleImage string            `json:"titleImage" binding:"required"`
	Images     []string          `json:"images" binding:"required"`
}

func (postController) create(c *gin.Context) {
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

	slugParts := strings.Split(strings.ToLower(json.Title["en"]), " ")
	slugLength := 4
	if slugLength > len(slugParts) {
		slugLength = len(slugParts)
	}

	var slugBuffer bytes.Buffer
	for i := 0; i < slugLength; i++ {
		slugBuffer.WriteString(sanitize.Name(slugParts[i]))
		if i < (slugLength - 1) {
			slugBuffer.WriteString("-")
		}
	}

	var slugLikePosts []model.Post
	app.DB().C(model.PostC).Find(
		bson.M{"slug": bson.RegEx{Pattern: slugBuffer.String(), Options: ""}},
	).All(&slugLikePosts)

	if len(slugLikePosts) > 0 {
		var slugLikes []string
		for _, post := range slugLikePosts {
			slugLikes = append(slugLikes, post.Slug)
		}

		slugIndex := 0
		slugBuffer.WriteString("-")
		slugBuffer.WriteString(strconv.Itoa(slugIndex))

		for nameInArray(slugBuffer.String(), slugLikes) {
			slugBuffer.Truncate(len(slugBuffer.String()) - 1)
			slugIndex++
			slugBuffer.WriteString(strconv.Itoa(slugIndex))
		}
	}

	insert := model.Post{
		ID:         bson.NewObjectId(),
		Author:     user.ToPartial(),
		Title:      json.Title,
		Slug:       slugBuffer.String(),
		TitleImage: json.TitleImage,
		Content:    json.Content,
		Images:     json.Images,
		Comments:   []model.Comment{},

		Upvotes:   0,
		Downvotes: 0,

		Created: time.Now(),
		Updated: nil,
		Deleted: nil,
	}

	err = app.DB().C(model.PostC).Insert(&insert)
	if err != nil {
		app.DbError(c, err)
		return
	}

	app.Created(c, insert.ID)
}

func nameInArray(name string, array []string) bool {
	for _, item := range array {
		if name == item {
			return true
		}
	}

	return false
}

func (postController) update(c *gin.Context) {
	var json createRequest
	update := bson.M{
		"title":      json.Title,
		"content":    json.Content,
		"titleImage": json.TitleImage,
		"images":     json.Images,
		"updated":    time.Now(),
	}

	err := app.DB().C(model.PostC).Update(
		bson.M{"_id": bson.ObjectIdHex(c.Param("id"))},
		update,
	)

	if err != nil {
		app.DbError(c, err)
		return
	}

	app.Ok(c, gin.H{"updated": c.Param("id")})
}

func (postController) destroy(c *gin.Context) {
	err := app.DB().C(model.PostC).Update(
		bson.M{"_id": c.Param("id")},
		bson.M{"$set": bson.M{"deleted": time.Now()}},
	)

	if err != nil {
		app.DbError(c, err)
		return
	}

	app.Ok(c, gin.H{"deleted": 1})
}

type createCommentRequest struct {
	ParentID string `json:"parentId"`
	Content  string `json:"content" binding:"required"`
}

func (postController) createComment(c *gin.Context) {
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
		ID:       bson.NewObjectId(),
		Author:   user.ToPartial(),
		ParentID: nil,
		Content:  sanitize.HTML(json.Content),

		Upvotes:   0,
		Downvotes: 0,

		Created: time.Now(),
		Updated: nil,
		Deleted: nil,
	}

	if json.ParentID != "" {
		parentID := bson.ObjectIdHex(json.ParentID)
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

func (postController) updateComment(c *gin.Context) {
	var json createCommentRequest
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
			"comments.$.content": json.Content,
			"comments.$.updated": time.Now(),
		}},
	)

	if err != nil {
		app.DbError(c, err)
		return
	}

	app.Ok(c, gin.H{"updated": c.Param("id")})
}

func (postController) destroyComment(c *gin.Context) {
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
