package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"chat-back/internal/config"
	"chat-back/internal/db"
)

func main() {

	// Cargar configuración
	cfg := config.Load()

	// Conectar Mongo
	mongoConn, err := db.NewMongo(cfg.MongoURI, cfg.MongoDB)
	if err != nil {
		log.Fatal("Mongo connection failed:", err)
	}

	log.Println("Mongo connected to:", mongoConn.Database.Name())

	// Crear router
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// Iniciar servidor
	r.Run(":" + cfg.Port)
}