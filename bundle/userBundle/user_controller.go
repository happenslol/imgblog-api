package userBundle

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/happeens/xkcdpass"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"

	"github.com/happeens/imgblog-api/app"
	"github.com/happeens/imgblog-api/model"
)

type userController struct{}

type authenticateRequest struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (userController) Authenticate(c *gin.Context) {
	var json authenticateRequest
	err := c.BindJSON(&json)
	if err != nil {
		app.BadRequest(c, err)
		return
	}

	user := model.User{}

	err = app.DB().C(model.UserC).Find(bson.M{"name": json.Name}).One(&user)
	if err != nil {
		app.DbError(c, err)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(json.Password))
	if err != nil {
		app.Unauthorized(c)
		return
	}

	token := app.CreateToken(user.Name, user.Role)

	app.Ok(c, gin.H{"token": token})
}

type registerRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Captcha  string `json:"captcha" binding:"required"`
	Password string `json:"password"`

	MailSettings model.MailSettings `json:"mailSettings"`
}

func (userController) Register(c *gin.Context) {
	var json registerRequest
	err := c.BindJSON(&json)
	if err != nil {
		app.BadRequest(c, err)
		return
	}

	// Check if captcha is valid
	if !app.ConfirmCaptcha(c.ClientIP(), json.Captcha) {
		app.BadRequest(c, errors.New("Captcha could not be confirmed"))
		return
	}

	// Create password if needed
	var pass string
	if json.Password == "" {
		pass = xkcdpass.GenerateDefaultChecked()
	} else {
		pass = json.Password
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)

	// Insert user into database
	insert := model.User{
		ID:       bson.NewObjectId(),
		Name:     json.Name,
		Password: string(hash[:]),
		Email:    json.Email,
		Role:     model.UserRole,

		MailSettings: json.MailSettings,
	}

	if err = app.DB().C(model.UserC).Insert(&insert); err != nil {
		app.DbError(c, err)
		return
	}

	if err := app.SendMail(
		app.MailContent{
			"some": "content goes here",
		},
		app.WelcomeMail,
		fmt.Sprintf("welcome, %v!", json.Name),
		json.Email,
	); err != nil {
		app.ServerError(c, err)
		return
	}

	app.Created(c, "0")
}

func (userController) Me(c *gin.Context) {
	user, _ := c.Get("user")

	var result model.User
	err := app.DB().C(model.UserC).Find(
		bson.M{"name": user},
	).One(&result)

	if err != nil {
		app.DbError(c, err)
		return
	}

	app.Ok(c, result)
}

func (userController) Index(c *gin.Context) {
	var result []model.User
	if err := app.DB().C(model.UserC).Find(nil).All(&result); err != nil {
		app.DbError(c, err)
		return
	}

	app.Ok(c, result)
}

func (userController) Show(c *gin.Context) {
	result := model.User{}
	id := bson.ObjectIdHex(c.Param("id"))

	if err := app.DB().C(model.UserC).FindId(id).One(&result); err != nil {
		app.DbError(c, err)
		return
	}

	app.Ok(c, result)
}

type createRequest struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Role     string `json:"role"`
}

func (userController) Create(c *gin.Context) {
	var json createRequest
	err := c.BindJSON(&json)
	if err != nil {
		app.BadRequest(c, err)
		return
	}

	role := json.Role
	if role != model.AdminRole && role != model.UserRole {
		app.BadRequest(c, err)
		return
	}

	hash, err := bcrypt.GenerateFromPassword(
		[]byte(json.Password),
		bcrypt.DefaultCost,
	)

	insert := model.User{
		ID:       bson.NewObjectId(),
		Name:     json.Name,
		Password: string(hash[:]),
		Email:    json.Email,
		Role:     role,

		MailSettings: model.MailSettings{
			ReceivePostNotifications: false,
			ReceiveNewsletters:       true,
		},
	}

	err = app.DB().C(model.UserC).Insert(&insert)
	if err != nil {
		app.DbError(c, err)
		return
	}

	app.Created(c, insert.ID)
}

type updateSettingsRequest struct {
	MailSettings model.MailSettings `json:"mailSettings"`
}

func (userController) UpdateSettings(c *gin.Context) {
	var json updateSettingsRequest
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

	update := bson.M{"mailSettings": json.MailSettings}

	err = app.DB().C(model.UserC).Update(
		bson.M{"_id": user.ID},
		bson.M{"$set": update},
	)

	if err != nil {
		app.DbError(c, err)
		return
	}

	app.Ok(c, gin.H{
		"mailSettings": json.MailSettings,
	})
}

func (userController) Destroy(c *gin.Context) {
	id := bson.ObjectIdHex(c.Param("id"))
	err := app.DB().C(model.UserC).RemoveId(id)
	if err != nil {
		app.DbError(c, err)
		return
	}

	app.Ok(c, gin.H{"deleted": 1})
}
