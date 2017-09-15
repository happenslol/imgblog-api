package newsBundle

import (
	"github.com/happeens/imgblog-api/app"
	"github.com/happeens/imgblog-api/model"
)

func init() {
	var newsCtrl = newsController{}

	news := app.Router.Group("/news")
	{
		news.GET("", newsCtrl.Index)
		news.GET("/latest", newsCtrl.Latest)
		news.POST("", app.RequireRole(model.AdminRole), newsCtrl.Create)
		news.PUT("/:id", app.RequireRole(model.AdminRole), newsCtrl.Update)
		news.DELETE("/:id", app.RequireRole(model.AdminRole), newsCtrl.Destroy)
	}

	app.Log.Info("newsbundle registered")
}
