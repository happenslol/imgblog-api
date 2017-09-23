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
	var catQuery interface{}
	cat := c.DefaultQuery("cat", "home")
	if cat == "home" {
		catQuery = bson.M{"$exists": false}
	} else {
		catQuery = cat
	}

	var result model.News
	err := app.DB().C(model.NewsC).Find(bson.M{
		"deleted":  nil,
		"category": catQuery,
	}).Sort("created").One(&result)
	if err != nil {
		app.DbError(c, err)
		return
	}

	app.Ok(c, result)
}

type sendRequest struct {
	TitleImage string              `json:"titleImage"`
	Sections   []newsletterSection `json:"sections"`
}

type newsletterSection struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

func (newsController) Send(c *gin.Context) {
	var json sendRequest
	if err := c.BindJSON(&json); err != nil {
		app.BadRequest(c, err)
		return
	}

	var recipients []model.User
	if err := app.DB().C(model.UserC).Find(
		bson.M{
			"email": bson.M{"$ne": nil},
			"mailSettings.receiveNewsletters": bson.M{"$eq": true},
		},
	).All(&recipients); err != nil {
		app.DbError(c, err)
		return
	}

	app.Ok(c, gin.H{"sent": len(recipients)})
}

type createRequest struct {
	Content  model.LocalString `json:"content"`
	Image    string            `json:"image"`
	Category string            `json:"category"`
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
		ID:       bson.NewObjectId(),
		Author:   user.ToPartial(),
		Content:  json.Content,
		Image:    json.Image,
		Category: json.Category,

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
	Content  model.LocalString `json:"content"`
	Image    string            `json:"image"`
	Category string            `json:"category"`
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
		"content":  json.Content,
		"image":    json.Image,
		"category": json.Category,
		"updated":  time.Now(),
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
