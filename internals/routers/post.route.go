package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ntisrangga142/chat/internals/handlers"
	"github.com/ntisrangga142/chat/internals/middlewares"
	"github.com/ntisrangga142/chat/internals/repositories"
	"github.com/redis/go-redis/v9"
)

func InitPost(ctx *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	repo := repositories.NewPostRepository(db)
	handler := handlers.NewPostHandler(repo, rdb)

	post := ctx.Group("/post")

	post.GET("", middlewares.Authentication, handler.GetFollowingPosts)
	post.GET("/:id", middlewares.Authentication, handler.GetPostDetail)
	post.POST("", middlewares.Authentication, handler.CreatePost)

	post.POST("/:id/like", handler.LikePost)
	post.DELETE("/:id/like", handler.UnlikePost)

	post.POST("/comment", handler.CreateComment)
	post.GET("/:id/comment", handler.GetAllCommentsByPost)
}
