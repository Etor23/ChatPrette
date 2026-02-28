package main

import (
	"github.com/gin-gonic/gin"

	"chat-back/internal/handlers"
	"chat-back/internal/repos"

	"go.mongodb.org/mongo-driver/mongo"
)

func SetupRouter(db *mongo.Database) *gin.Engine {

	r := gin.Default()

	// ===== Repos =====
	userRepo := repos.NewUserRepo(db)

	// ===== Handlers =====
	userHandler := handlers.NewUserHandler(userRepo)

	// ===== Routes =====
	api := r.Group("/api")
	{
		api.GET("/users", userHandler.GetAllUsers)
		api.GET("/users/get/:_id", userHandler.GetUserById)
		api.POST("/users", userHandler.CreateUser)
	}

	// Health check
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	return r
}