package postBundle

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
	"time"

	"github.com/happeens/imgblog-api/app"
	"github.com/happeens/imgblog-api/model"
)

type draftController struct{}

func (draftController) Index(c *gin.Context) {
	var result []model.Draft
	if err := app.DB().C(model.DraftC).Find(
		bson.M{"deleted": nil},
	).All(&result); err != nil {
		app.DbError(c, err)
		return
	}

	app.Ok(c, result)
}

func (draftController) Create(c *gin.Context) {
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

	toSave := model.Draft{
		ID:         bson.NewObjectId(),
		Author:     user.ToPartial(),
		Title:      req.Title,
		Intro:      req.Intro,
		TitleImage: req.TitleImage,
		Sections:   req.Sections,

		Category: req.Category,
		Tags:     req.Tags,

		Created: time.Now(),
	}

	if err := app.DB().C(model.DraftC).Insert(
		&toSave,
	); err != nil {
		app.DbError(c, err)
		return
	}

	app.Created(c, toSave)
}

func (draftController) Update(c *gin.Context) {
	var req createRequest
	if err := c.BindJSON(&req); err != nil {
		app.BadRequest(c, err)
		return
	}

	id := bson.ObjectIdHex(c.Param("id"))
	toSave := model.ToMap(req)
	toSave["updated"] = time.Now()

	if err := app.DB().C(model.DraftC).Update(
		bson.M{"_id": id},
		bson.M{"$set": toSave},
	); err != nil {
		app.DbError(c, err)
		return
	}

	var updated model.Draft
	if err := app.DB().C(model.DraftC).FindId(
		id,
	).One(&updated); err != nil {
		app.DbError(c, err)
		return
	}

	app.Ok(c, updated)
}

func (draftController) Destroy(c *gin.Context) {
	if err := app.DB().C(model.DraftC).Update(
		bson.M{"_id": bson.ObjectIdHex(c.Param("id"))},
		bson.M{"$set": bson.M{"deleted": time.Now()}},
	); err != nil {
		app.DbError(c, err)
		return
	}

	app.Deleted(c)
}

func (draftController) Publish(c *gin.Context) {
	id := bson.ObjectIdHex(c.Param("id"))

	var draft model.Draft
	if err := app.DB().C(model.DraftC).FindId(id).One(
		&draft,
	); err != nil {
		app.DbError(c, err)
		return
	}

	toSave := model.Post{
		ID:         bson.NewObjectId(),
		Author:     draft.Author,
		Title:      draft.Title,
		Intro:      draft.Intro,
		Slug:       createSlug(draft.Title["en"]),
		TitleImage: draft.TitleImage,
		Sections:   draft.Sections,
		Comments:   []model.Comment{},

		Category: draft.Category,
		Tags:     draft.Tags,

		Created: time.Now(),
	}

	if err := app.DB().C(model.PostC).Insert(
		&toSave,
	); err != nil {
		app.DbError(c, err)
		return
	}

	if err := app.DB().C(model.DraftC).Update(
		bson.M{"_id": id},
		bson.M{"$set": bson.M{"deleted": time.Now()}},
	); err != nil {
		app.DbError(c, err)
		return
	}

	app.Created(c, toSave)
}
