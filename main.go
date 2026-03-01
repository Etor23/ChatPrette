package main

import (
	"log"

	"chat-back/internal/config"
	"chat-back/internal/db"
	"chat-back/internal/firebase"
)

func main() {

	cfg := config.Load()

	mongoConn, err := db.NewMongo(cfg.MongoURI, cfg.MongoDB)
	if err != nil {
		log.Fatal(err)
	}

	fbClient, err := firebase.NewFirebaseClient(cfg.FirebaseCredentials)
	if err != nil {
		log.Fatal("Error inicializando Firebase:", err)
	}

	r := SetupRouter(mongoConn.Database, fbClient)

	r.Run(":" + cfg.Port)
}
