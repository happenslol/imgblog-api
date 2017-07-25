package postBundle

import (
	"github.com/happeens/imgblog-api/app"
	"github.com/happeens/imgblog-api/model"
)

func init() {
	var postCtrl = postController{}
	var voteCtrl = voteController{}

	app.Router.GET("/search", postCtrl.search)

	posts := app.Router.Group("/posts")
	{
		posts.GET("", postCtrl.index)
		posts.GET("/:slug", postCtrl.show)
		posts.POST("", app.RequireRole(model.AdminRole), postCtrl.create)
		posts.PUT("/:id", app.RequireRole(model.AdminRole), postCtrl.update)
		posts.DELETE("/:id", app.RequireRole(model.AdminRole), postCtrl.destroy)

		posts.POST("/:id/comments", app.RequireAuth(), postCtrl.createComment)
		posts.PUT("/:id/comments/:commentId", app.RequireAuth(), postCtrl.updateComment)
		posts.DELETE("/:id/comments/:commentId", app.RequireAuth(), postCtrl.destroyComment)
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
