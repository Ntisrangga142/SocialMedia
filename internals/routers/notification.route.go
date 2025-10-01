package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ntisrangga142/chat/internals/handlers"
	"github.com/ntisrangga142/chat/internals/middlewares"
	"github.com/ntisrangga142/chat/internals/repositories"
	"github.com/redis/go-redis/v9"
)

func InitNotif(ctx *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	repo := repositories.NewNotificationRepository(db)
	handler := handlers.NewNotificationHandler(repo)

	notif := ctx.Group("/notif")
	notif.Use(middlewares.Authentication)

	notif.GET("", handler.GetUnreadNotifications)

}
