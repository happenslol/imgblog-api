package imageBundle

import (
	"errors"
	"fmt"
	"image"
	"os"
	"strconv"
	"strings"

	_ "image/jpeg"
	_ "image/png"

	"github.com/gin-gonic/gin"
	"github.com/oliamb/cutter"
	"github.com/satori/go.uuid"

	"github.com/happeens/imgblog-api/app"
)

const thumbWidth = 800
const thumbHeight = 200

const fullWidth = 1920
const fullHeight = 600

type imageController struct{}

func (imageController) Upload(c *gin.Context) {
	uuid := uuid.NewV4()

	cropTop := 0
	cropLeft := 0

	var err error
	if c.PostForm("top") != "" {
		cropTop, err = strconv.Atoi(c.PostForm("top"))
		if err != nil {
			app.BadRequest(c, errors.New("Invalid top position"))
			return
		}
	}

	if c.PostForm("left") != "" {
		cropLeft, err = strconv.Atoi(c.PostForm("left"))
		if err != nil {
			app.BadRequest(c, errors.New("Invalid left position"))
			return
		}
	}

	file, err := c.FormFile("file")
	if err != nil {
		app.BadRequest(c, err)
		return
	}

	parts := strings.Split(file.Filename, ".")

	if len(parts) < 2 {
		app.BadRequest(c, errors.New("invalid filename"))
		return
	}

	if parts[1] != "jpg" && parts[1] != "png" {
		app.BadRequest(c, errors.New("invalid filetype"))
		return
	}

	pathOriginal := fmt.Sprintf("%v/%v-o.%v", app.StoragePath, uuid, parts[1])
	pathThumb := fmt.Sprintf("%v/%v-thumb.%v", app.StoragePath, uuid, parts[1])
	pathFull := fmt.Sprintf("%v/%v-full.%v", app.StoragePath, uuid, parts[1])

	fmt.Printf("o: %v, thumb: %v, full: %v", pathOriginal, pathThumb, pathFull)

	// save original file to disk
	err = c.SaveUploadedFile(file, pathOriginal)
	if err != nil {
		app.ServerError(c, err)
		return
	}

	// load image
	fileOriginal, err := os.Open(pathOriginal)
	if err != nil {
		app.ServerError(c, err)
		return
	}

	imgOriginal, _, err := image.Decode(fileOriginal)
	if err != nil {
		app.ServerError(c, err)
		return
	}

	app.Ok(c, gin.H{"name": uuid})
}
