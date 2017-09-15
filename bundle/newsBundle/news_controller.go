package newsBundle

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
	"time"

	"github.com/happeens/imgblog-api/app"
	"github.com/happeens/imgblog-api/model"
)

type newsController struct{}

func (newsController) Index(c *gin.Context) {
	var result []model.News
	if err := app.DB().C(model.NewsC).Find(
		bson.M{"deleted": nil},
	).All(&result); err != nil {
		app.DbError(c, err)
		return
	}

	app.Ok(c, result)
}

func (newsController) Latest(c *gin.Context) {
	var result model.News
	err := app.DB().C(model.NewsC).Find(
		bson.M{"deleted": nil},
	).Sort("created").One(&result)
	if err != nil {
		app.DbError(c, err)
		return
	}

	app.Ok(c, result)
}

type createRequest struct {
	Content model.LocalString `json:"content"`
	Image   string            `json:"image"`
}

func (newsController) Create(c *gin.Context) {
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

	insert := model.News{
		ID:      bson.NewObjectId(),
		Author:  user.ToPartial(),
		Content: json.Content,
		Image:   json.Image,

		Created: time.Now(),
		Updated: nil,
		Deleted: nil,
	}

	err = app.DB().C(model.NewsC).Insert(&insert)
	if err != nil {
		app.DbError(c, err)
		return
	}

	app.Ok(c, insert)
}

type updateRequest struct {
	Content model.LocalString `json:"content"`
	Image   string            `json:"image"`
}

func (newsController) Update(c *gin.Context) {
	var json updateRequest
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

	var news model.News
	err = app.DB().C(model.NewsC).FindId(
		bson.ObjectIdHex(c.Param("id")),
	).One(&news)
	if err != nil {
		//TODO not found error
		app.BadRequest(c, err)
		return
	}

	if news.Author.ID != user.ID {
		app.Unauthorized(c)
		return
	}

	update := bson.M{
		"content": json.Content,
		"image":   json.Image,
		"updated": time.Now(),
	}

	err = app.DB().C(model.NewsC).Update(
		bson.M{"_id": bson.ObjectIdHex(c.Param("id"))},
		bson.M{"$set": update},
	)

	if err != nil {
		app.DbError(c, err)
		return
	}

	var updated model.News
	err = app.DB().C(model.NewsC).FindId(news.ID).One(&updated)

	if err != nil {
		app.DbError(c, err)
		return
	}

	app.Ok(c, updated)
}

func (newsController) Destroy(c *gin.Context) {
	if err := app.DB().C(model.NewsC).Update(
		bson.M{"_id": bson.ObjectIdHex(c.Param("id"))},
		bson.M{"$set": bson.M{"deleted": time.Now()}},
	); err != nil {
		app.DbError(c, err)
		return
	}

	//TODO deleted return
	app.Ok(c, gin.H{"deleted": c.Param("id")})
}
