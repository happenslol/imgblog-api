package imageBundle

import (
	"errors"
	"fmt"
	"image"
	"os"
	"strconv"
	"strings"

	"image/jpeg"
	_ "image/png"

	"github.com/gin-gonic/gin"
	"github.com/nfnt/resize"
	"github.com/oliamb/cutter"
	"github.com/satori/go.uuid"

	"github.com/happeens/imgblog-api/app"
)

const thumbWidth = 800
const thumbHeight = 200

const fullWidth = 1920
const fullHeight = 720

type imageController struct{}

func (imageController) Upload(c *gin.Context) {
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

	uuid := uuid.NewV4()
	pathOriginal := fmt.Sprintf("%v/%v-o.%v", app.StoragePath, uuid, parts[1])
	pathThumb := fmt.Sprintf("%v/%v-thumb.jpg", app.StoragePath, uuid)
	pathFull := fmt.Sprintf("%v/%v-full.jpg", app.StoragePath, uuid)

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
	fileOriginal.Close()

	originalBounds := imgOriginal.Bounds()

	x := getDefaultInt(c, "x", originalBounds.Max.X/2)
	y := getDefaultInt(c, "y", originalBounds.Max.Y/2)

	croppedFull, err := cutter.Crop(imgOriginal, cutter.Config{
		Width:   16,
		Height:  6,
		Anchor:  image.Point{x, y},
		Mode:    cutter.Centered,
		Options: cutter.Ratio,
	})

	if err != nil {
		app.ServerError(c, err)
		return
	}

	thumbX := getDefaultInt(c, "thumbx", originalBounds.Max.X/2)
	thumbY := getDefaultInt(c, "thumby", originalBounds.Max.Y/2)

	croppedThumb, err := cutter.Crop(imgOriginal, cutter.Config{
		Width:   16,
		Height:  4,
		Anchor:  image.Point{thumbX, thumbY},
		Mode:    cutter.Centered,
		Options: cutter.Ratio,
	})

	if err != nil {
		app.ServerError(c, err)
		return
	}

	imgFull := resize.Resize(uint(fullWidth), uint(fullHeight), croppedFull, resize.Lanczos3)
	imgThumb := resize.Resize(uint(thumbWidth), uint(thumbHeight), croppedThumb, resize.Bilinear)

	fullOut, err := os.Create(pathFull)
	if err != nil {
		app.ServerError(c, err)
		return
	}
	defer fullOut.Close()

	thumbOut, err := os.Create(pathThumb)
	if err != nil {
		app.ServerError(c, err)
		return
	}
	defer thumbOut.Close()

	jpeg.Encode(fullOut, imgFull, nil)
	jpeg.Encode(thumbOut, imgThumb, nil)

	app.Ok(c, gin.H{"name": uuid})
}

func getDefaultInt(c *gin.Context, name string, def int) int {
	val := c.DefaultPostForm(name, strconv.Itoa(def))
	i, err := strconv.Atoi(val)
	if err != nil {
		return 0
	}

	if i < 0 {
		return 0
	}

	return i
}
