package imageBundle

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"

	"github.com/happeens/imgblog-api/app"
)

type imageController struct{}

func (imageController) Upload(c *gin.Context) {
	uuid := uuid.NewV4()
	file, err := c.FormFile("file")
	if err != nil {
		app.BadRequest(c, err)
		return
	}

	filename := fmt.Sprintf("./static/%s.jpg", uuid)

	//TODO to jpg or png always
	err = c.SaveUploadedFile(file, filename)
	if err != nil {
		app.BadRequest(c, err)
		return
	}

	app.Ok(c, gin.H{"image": filename})
}
