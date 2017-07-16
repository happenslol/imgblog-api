package postBundle

import (
	"gopkg.in/gin-gonic/gin.v1"
	"gopkg.in/mgo.v2/bson"

	"github.com/happeens/imgblog-api/app"
	"github.com/happeens/imgblog-api/model"
)

type postController struct{}

func (postController) Index(c *gin.Context) {
	var result []model.Post
	err := app.DB().C(model.PostC).Find(nil).All(&result)
	if err != nil {
		app.DbError(c, err)
		return
	}

	app.Ok(c, result)
}

type createRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

func (postController) Create(c *gin.Context) {
	var json createRequest
	err := c.BindJSON(&json)
	if err != nil {
		app.BadRequest(c, err)
		return
	}

	insert := model.Post{
		ID:      bson.NewObjectId(),
		Title:   json.Title,
		Content: json.Content,
	}

	err = app.DB().C(model.PostC).Insert(&insert)
	if err != nil {
		app.DbError(c, err)
		return
	}

	app.Created(c, insert.ID)
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
