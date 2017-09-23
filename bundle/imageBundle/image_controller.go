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

	sizes := strings.Split(c.DefaultPostForm("sizes", "full"), ",")
	uuid := uuid.NewV4()
	pathOriginal := fmt.Sprintf("%v/%v-o.%v", app.StoragePath, uuid, parts[1])

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

	for _, size := range sizes {
		err = saveSize(c, uuid.String(), size, imgOriginal)

		if err != nil {
			app.ServerError(c, err)
			return
		}
	}

	app.Ok(c, gin.H{"name": uuid})
}

func saveSize(c *gin.Context, uuid, size string, original image.Image) error {
	var w, h, wr, hr int
	if size == "news" {
		w, h, wr, hr = 150, 200, 3, 4
	} else if size == "thumb" {
		w, h, wr, hr = 800, 200, 16, 4
	} else if size == "header" {
		w, h, wr, hr = 1920, 720, 16, 6
	} else if size == "full" {
		w, h, wr, hr = 1920, 1080, 16, 9
	} else {
		return errors.New("invalid size given")
	}

	bounds := original.Bounds()

	x := getDefaultInt(c, fmt.Sprintf("%v-x", size), bounds.Max.X/2)
	y := getDefaultInt(c, fmt.Sprintf("%v-y", size), bounds.Max.Y/2)

	cropped, err := cutter.Crop(original, cutter.Config{
		Width:   wr,
		Height:  hr,
		Anchor:  image.Point{x, y},
		Mode:    cutter.Centered,
		Options: cutter.Ratio,
	})

	if err != nil {
		return err
	}

	resized := resize.Resize(uint(w), uint(h), cropped, resize.Lanczos3)

	path := fmt.Sprintf("%v/%v-%v.jpg", app.StoragePath, uuid, size)
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	jpeg.Encode(out, resized, nil)
	return nil
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
