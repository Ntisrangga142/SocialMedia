package main

import (
	"context"
	"log"

	"github.com/joho/godotenv"
	"github.com/ntisrangga142/chat/internals/configs"
	"github.com/ntisrangga142/chat/internals/routers"
)

func main() {
	// Load ENV
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Init Database
	db, err := configs.InitDB()
	if err != nil {
		log.Printf("DB ERROR: %s", err.Error())
	}

	// Ping Database
	if err := configs.PingDB(db); err != nil {
		log.Printf("DB ERROR: %s", err.Error())
	}
	log.Println("Database connected")

	//Init Redis
	rdb := configs.InitRedis()
	if cmd := rdb.Ping(context.Background()); cmd.Err() != nil {
		log.Println("Ping to Redis failed\nCause: ", cmd.Err().Error())
		return
	}
	log.Println("Redis Connected")
	defer rdb.Close()

	router := routers.InitRouter(db, rdb)
	router.Run(":8080")
}
