package postBundle

import (
	"bytes"
	"strconv"
	"strings"

	"github.com/kennygrant/sanitize"
	"gopkg.in/mgo.v2/bson"

	"github.com/happeens/imgblog-api/app"
	"github.com/happeens/imgblog-api/model"
)

func createSlug(s string) string {
	parts := strings.Split(strings.ToLower(s), " ")
	length := 4
	if length > len(parts) {
		length = len(parts)
	}

	var b bytes.Buffer
	for i := 0; i < length; i++ {
		b.WriteString(sanitize.Name(parts[i]))
		if i < (length - 1) {
			b.WriteString("-")
		}
	}

	var slugLikePosts []model.Post
	app.DB().C(model.PostC).Find(
		bson.M{"slug": bson.RegEx{
			Pattern: b.String(), Options: ""},
		},
	).All(&slugLikePosts)

	if len(slugLikePosts) > 0 {
		var slugLikes []string
		for _, post := range slugLikePosts {
			slugLikes = append(slugLikes, post.Slug)
		}

		slugIndex := 0
		b.WriteString("-")
		b.WriteString(strconv.Itoa(slugIndex))

		for nameInArray(b.String(), slugLikes) {
			b.Truncate(len(b.String()) - 1)
			slugIndex++
			b.WriteString(strconv.Itoa(slugIndex))
		}
	}

	return b.String()
}

func nameInArray(name string, array []string) bool {
	for _, item := range array {
		if name == item {
			return true
		}
	}

	return false
}
