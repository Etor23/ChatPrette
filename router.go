package main

import (
	"github.com/gin-gonic/gin"

	"chat-back/internal/firebase"
	"chat-back/internal/handlers"
	"chat-back/internal/middleware"
	"chat-back/internal/repos"

	"go.mongodb.org/mongo-driver/mongo"
)

func SetupRouter(db *mongo.Database, fbClient *firebase.FirebaseClient) *gin.Engine {

	r := gin.Default()

	// ===== Repos =====
	userRepo := repos.NewUserRepo(db)

	// ===== Handlers =====
	userHandler := handlers.NewUserHandler(userRepo)
	authHandler := handlers.NewAuthHandler(fbClient.Auth, userRepo)

	// ===== Middleware =====
	authMiddleware := middleware.AuthMiddleware(fbClient.Auth)

	// ===== Routes =====
	api := r.Group("/api")
	{
		// Auth (pública)
		api.POST("/auth/login", authHandler.Login)
		api.GET("/auth/me", authMiddleware, authHandler.Me)

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
