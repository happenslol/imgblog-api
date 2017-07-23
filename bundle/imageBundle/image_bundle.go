package imageBundle

import (
	"github.com/happeens/imgblog-api/app"
	"github.com/happeens/imgblog-api/model"
)

func init() {
	var imageCtrl = imageController{}

	images := app.Router.Group("/images")
	{
		images.POST("", app.RequireRole(model.AdminRole), imageCtrl.Upload)
	}
}
