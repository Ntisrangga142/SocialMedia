package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ntisrangga142/chat/internals/handlers"
	"github.com/ntisrangga142/chat/internals/repositories"
	"github.com/redis/go-redis/v9"
)

func InitAuth(ctx *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	repo := repositories.NewAuthRepo(db, rdb)
	handler := handlers.NewAuthHandler(repo, rdb)

	auth := ctx.Group("/auth")

	// Login
	auth.POST("", handler.Login)

	// Register
	auth.POST("/register", handler.Register)
}
