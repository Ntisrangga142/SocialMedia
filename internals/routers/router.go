package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	docs "github.com/ntisrangga142/chat/docs"
	"github.com/ntisrangga142/chat/internals/middlewares"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRouter(db *pgxpool.Pool, rdb *redis.Client) *gin.Engine {
	router := gin.Default()

	docs.SwaggerInfo.Title = "Social Media API"
	docs.SwaggerInfo.Description = "This is a sample API for social media"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	router.GET("/chat/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	middlewares.InitRedis(rdb)

	router.Static("/avatar", "./public/profile")
	router.Static("/img", "./public/post")

	InitAuth(router, db, rdb)
	InitUser(router, db, rdb)
	InitPost(router, db, rdb)
	InitNotif(router, db, rdb)

	return router
}
