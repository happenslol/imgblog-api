package userBundle

import (
	"github.com/happeens/imgblog-api/app"
	"github.com/happeens/imgblog-api/model"
)

func init() {
	var userCtrl = userController{}

	app.Router.POST("/authenticate", userCtrl.Authenticate)
	app.Router.POST("/register", userCtrl.Create)
	app.Router.GET("/me", app.RequireAuth(), userCtrl.Me)

	users := app.Router.Group("/users")
	{
		users.GET("", app.RequireRole(model.AdminRole), userCtrl.Index)
		users.GET("/:id", app.RequireRole(model.AdminRole), userCtrl.Show)
		users.DELETE("/:id", app.RequireRole(model.AdminRole), userCtrl.Destroy)
	}
}
