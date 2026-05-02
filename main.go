package main

import (
	"log"
	"os"
	"time"

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

	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		secretKey = "your-secret-key-change-in-production"
	}
	jwtManager := auth.NewJWTManager(secretKey, 24*time.Hour)

	// Router recibe el manager JWT compartido
	r := SetupRouter(mongoConn.Database, jwtManager)

	log.Printf(" Servidor corriendo en http://localhost:%s\n", cfg.Port)
	r.Run(":" + cfg.Port)
}
