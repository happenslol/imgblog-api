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

const postsPageSize = 10

type postController struct{}

func (postController) Index(c *gin.Context) {
	queryObject := bson.M{
		"deleted": nil,
	}

	if cat := c.DefaultQuery("cat", ""); cat != "" {
		queryObject["category"] = cat
	}

	query := app.DB().C(model.PostC).Find(queryObject)

	pageQuery := c.DefaultQuery("page", "")
	pageSizeQuery := c.DefaultQuery("pageSize", string(postsPageSize))

	if pageQuery != "" {
		page, err := strconv.Atoi(pageQuery)
		if err != nil {
			app.BadRequest(c, err)
			return
		}

		size, err := strconv.Atoi(pageSizeQuery)
		if err != nil {
			app.BadRequest(c, err)
			return
		}

		query.Skip(page * size).Limit(size)
	}

	if sortBy := c.DefaultQuery("sort", ""); sortBy != "" {
		query.Sort(sortBy)
	} else {
		query.Sort("created")
	}

	var result []model.Post
	if err := query.All(&result); err != nil {
		app.DbError(c, err)
		return
	}

	app.Ok(c, result)
}

func (postController) Show(c *gin.Context) {
	var result model.Post
	if err := app.DB().C(model.PostC).Find(
		bson.M{"slug": c.Param("slug")},
	).One(&result); err != nil {
		app.DbError(c, err)
		return
	}

	app.Ok(c, result)
}

//TODO better search/tag search
func (postController) Search(c *gin.Context) {
	var result []model.Post
	query := c.DefaultQuery("query", "")
	if query == "" {
		app.BadRequest(c, errors.New("empty search query not allowed"))
		return
	}

	if err := app.DB().C(model.PostC).Find(
		bson.M{"$or": []bson.M{
			bson.M{"title.en": bson.RegEx{Pattern: query, Options: "i"}},
			bson.M{"title.de": bson.RegEx{Pattern: query, Options: "i"}},
		}},
	).All(&result); err != nil {
		app.DbError(c, err)
		return
	}

	app.Ok(c, result)
}

type createRequest struct {
	Title      model.LocalString   `json:"title" binding:"required"`
	Intro      model.LocalString   `json:"intro" binding:"required"`
	Sections   []model.PostSection `json:"sections" binding:"required"`
	TitleImage string              `json:"titleImage" binding:"required"`
	Category   string              `json:"category" binding:"required"`
	Tags       []string            `json:"tags"`
}

func (postController) Create(c *gin.Context) {
	var req createRequest
	if err := c.BindJSON(&req); err != nil {
		app.BadRequest(c, err)
		return
	}

	user := model.User{}
	userName, _ := c.Get("user")

	if err := app.DB().C(model.UserC).Find(
		bson.M{"name": userName},
	).One(&user); err != nil {
		app.DbError(c, err)
		return
	}

	slugParts := strings.Split(strings.ToLower(req.Title["en"]), " ")
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
		bson.M{"slug": bson.RegEx{
			Pattern: slugBuffer.String(), Options: ""},
		},
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
		Title:      req.Title,
		Intro:      req.Intro,
		Slug:       slugBuffer.String(),
		TitleImage: req.TitleImage,
		Sections:   req.Sections,
		Comments:   []model.Comment{},

		Category: req.Category,
		Tags:     req.Tags,

		Upvotes: 0,

		Created: time.Now(),
		Updated: nil,
		Deleted: nil,
	}

	if err := app.DB().C(model.PostC).Insert(
		&insert,
	); err != nil {
		app.DbError(c, err)
		return
	}

	app.Created(c, insert)
}

func nameInArray(name string, array []string) bool {
	for _, item := range array {
		if name == item {
			return true
		}
	}

	return false
}

func (postController) Update(c *gin.Context) {
	var req createRequest
	if err := c.BindJSON(&req); err != nil {
		app.BadRequest(c, err)
		return
	}

	postId := bson.ObjectIdHex(c.Param("id"))

	if err := app.DB().C(model.PostC).Update(
		bson.M{"_id": postId},
		bson.M{"$set": req},
	); err != nil {
		app.DbError(c, err)
		return
	}

	var updated model.Post
	if err := app.DB().C(model.PostC).FindId(
		postId,
	).One(&updated); err != nil {
		app.DbError(c, err)
		return
	}

	app.Ok(c, updated)
}

func (postController) Destroy(c *gin.Context) {
	err := app.DB().C(model.PostC).Update(
		bson.M{"_id": c.Param("id")},
		bson.M{"$set": bson.M{"deleted": time.Now()}},
	)

	if err != nil {
		app.DbError(c, err)
		return
	}

	app.Deleted(c)
}
