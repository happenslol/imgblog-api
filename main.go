package main

import (
	"github.com/happeens/imgblog-api/app"

	_ "github.com/happeens/imgblog-api/bundle/imageBundle"
	_ "github.com/happeens/imgblog-api/bundle/newsBundle"
	_ "github.com/happeens/imgblog-api/bundle/postBundle"
	_ "github.com/happeens/imgblog-api/bundle/userBundle"
)

func main() {
	port := ":" + app.Env("PORT", "8000")
	app.Router.Run(port)
}
