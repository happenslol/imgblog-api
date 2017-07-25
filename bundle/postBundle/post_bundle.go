package postBundle

import (
	"github.com/happeens/imgblog-api/app"
	"github.com/happeens/imgblog-api/model"
)

func init() {
	var postCtrl = postController{}
	var voteCtrl = voteController{}

	app.Router.GET("/search", postCtrl.Search)

	posts := app.Router.Group("/posts")
	{
		posts.GET("", postCtrl.Index)
		posts.GET("/:slug", postCtrl.Show)
		posts.POST("", app.RequireRole(model.AdminRole), postCtrl.Create)
		posts.DELETE("/:id", app.RequireRole(model.AdminRole), postCtrl.Destroy)

		posts.POST("/:id/comments", app.RequireAuth(), postCtrl.CreateComment)
	}

	votes := app.Router.Group("/votes")
	{
		votes.GET("/user/:id", voteCtrl.showUserVotes)
		votes.GET("/post/:id", voteCtrl.showPostVotes)
		votes.GET("/comment/:id", voteCtrl.showCommentVotes)

		votes.POST("/post/:id", app.RequireAuth(), voteCtrl.createPostVote)
		votes.POST("/comment/:id", app.RequireAuth(), voteCtrl.createCommentVote)
	}
}
