package postBundle

import (
	"github.com/happeens/imgblog-api/app"
	"github.com/happeens/imgblog-api/model"
)

func init() {
	var postCtrl = postController{}

	posts := app.Router.Group("/posts")
	{
		posts.GET("", postCtrl.Index)
		posts.POST("", app.RequireRole(model.AdminRole), postCtrl.Create)
		posts.DELETE("/:id", app.RequireRole(model.Admin), postCtrl.Destroy)
	}
}
