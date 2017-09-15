package newsBundle

import (
	"github.com/happeens/imgblog-api/app"
	"github.com/happeens/imgblog-api/model"
)

func init() {
	var newsCtrl = newsController{}

	news := app.Router.Group("/news")
	{
		news.GET("", newsCtrl.index)
		news.GET("/latest", newsCtrl.latest)
		news.POST("", app.RequireRole(model.AdminRole), newsCtrl.create)
		news.PUT("/:id", app.RequireRole(model.AdminRole), newsCtrl.update)
		news.DELETE("/:id", app.RequireRole(model.AdminRole), newsCtrl.destroy)
	}

	app.Log.Info("newsbundle registered")
}
