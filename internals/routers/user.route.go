package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ntisrangga142/chat/internals/handlers"
	"github.com/ntisrangga142/chat/internals/middlewares"
	"github.com/ntisrangga142/chat/internals/repositories"
	"github.com/redis/go-redis/v9"
)

func InitUser(ctx *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	repo := repositories.NewUserRepository(db)
	handler := handlers.NewUserHandler(repo, rdb)

	user := ctx.Group("/user")
	user.Use(middlewares.Authentication)

	// Get Profile
	user.GET("", handler.GetProfile)
	// Update Profile
	user.PATCH("", handler.UpdateProfile)

	// Follow
	user.POST(":id", handler.Follow)
	// Unfollow
	user.DELETE(":id", handler.Unfollow)
	// Get Followers
	user.GET("/follower", handler.GetFollowers)
	// Get Following
	user.GET("/following", handler.GetFollowing)
}
