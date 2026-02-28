package main

import (
	"log"

	"chat-back/internal/config"
	"chat-back/internal/db"
)

func main() {

	cfg := config.Load()

	mongoConn, err := db.NewMongo(cfg.MongoURI, cfg.MongoDB)
	if err != nil {
		log.Fatal(err)
	}

	r := SetupRouter(mongoConn.Database)

	r.Run(":" + cfg.Port)
}