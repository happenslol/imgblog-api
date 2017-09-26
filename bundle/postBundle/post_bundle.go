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
		posts.PUT("/:id", app.RequireRole(model.AdminRole), postCtrl.Update)
		posts.DELETE("/:id", app.RequireRole(model.AdminRole), postCtrl.Destroy)

		posts.POST("/:id/comments", app.RequireAuth(), postCtrl.CreateComment)
		posts.PUT("/:id/comments/:commentId", app.RequireAuth(), postCtrl.UpdateComment)
		posts.DELETE("/:id/comments/:commentId", app.RequireAuth(), postCtrl.DestroyComment)
	}

	votes := app.Router.Group("/votes")
	{
		votes.GET("/user/:id", voteCtrl.ShowUserVotes)
		votes.GET("/post/:id", voteCtrl.ShowPostVotes)
		votes.GET("/comment/:id", voteCtrl.ShowCommentVotes)

		votes.POST("/post/:id", app.RequireAuth(), voteCtrl.CreatePostVote)
		votes.POST("/comment/:id", app.RequireAuth(), voteCtrl.CreateCommentVote)

		votes.DELETE("/post/:id", app.RequireAuth(), voteCtrl.DeletePostVote)
		votes.DELETE("/comment/:id", app.RequireAuth(), voteCtrl.DeleteCommentVote)
	}

	app.Log.Info("postbundle registered")
}
