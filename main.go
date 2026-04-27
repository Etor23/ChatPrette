package main

import (
	"log"

	"chat-back/internal/auth"
	"chat-back/internal/config"
	"chat-back/internal/db"
)

func main() {

	cfg := config.Load()

	// MongoDB
	mongoConn, err := db.NewMongo(cfg.MongoURI, cfg.MongoDB)
	if err != nil {
		log.Fatal(" Error conectando a MongoDB:", err)
	}

	// Firebase                                      ← NUEVO
	firebaseAuth, err := auth.NewFirebaseProvider(cfg.FirebaseCredentials)
	if err != nil {
		log.Fatal(" Error inicializando Firebase:", err)
	}

	// Router (ahora recibe Firebase también)        ← MODIFICADO
	r := SetupRouter(mongoConn.Database, firebaseAuth)

	log.Printf(" Servidor corriendo en http://localhost:%s\n", cfg.Port)
	r.Run(":" + cfg.Port)
}
