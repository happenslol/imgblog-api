package postBundle

import (
	"github.com/happeens/imgblog-api/app"
	"github.com/happeens/imgblog-api/model"
)

func init() {
	var postCtrl = postController{}

	app.Router.GET("/search", postCtrl.Search)

	posts := app.Router.Group("/posts")
	{
		posts.GET("", postCtrl.Index)
		posts.GET("/:slug", postCtrl.Show)
		posts.POST("", app.RequireRole(model.AdminRole), postCtrl.Create)
		posts.DELETE("/:id", app.RequireRole(model.AdminRole), postCtrl.Destroy)

		posts.POST("/:id/comments", app.RequireAuth(), postCtrl.CreateComment)
	}
}
