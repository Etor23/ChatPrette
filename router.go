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
	conversationRepo := repos.NewConversationRepo(db)

	// ===== Handlers =====
	userHandler := handlers.NewUserHandler(userRepo)
	conversationHandler := handlers.NewConversationHandler(conversationRepo)
	authHandler := handlers.NewAuthHandler(fbClient.Auth, userRepo)

	// ===== Middleware =====
	authMiddleware := middleware.AuthMiddleware(fbClient.Auth)

	// ===== Routes =====
	api := r.Group("/api")
	{
		// Auth (pública)
		api.POST("/auth/login", authHandler.Login)
		api.GET("/auth/me", authMiddleware, authHandler.Me)

		// Users
		api.GET("/users", userHandler.GetAllUsers)
		api.GET("/users/get/:_id", userHandler.GetUserById)
		api.POST("/users", userHandler.CreateUser)
		api.PUT("/users/:_id", userHandler.UpdateUser)
		api.DELETE("/users/:_id", userHandler.DeleteUser)
		
		//Conversations
		api.GET("/conversations", conversationHandler.GetUserConversations)
		api.GET("/conversations/get/:_id", conversationHandler.GetConversationById)
		api.POST("/conversations/dm", conversationHandler.CreateDm)
		api.POST("/conversations/group", conversationHandler.CreateGroup)
		api.PATCH("/conversations/:_id/name", conversationHandler.UpdateGroupName)
		api.PATCH("/conversations/:_id/members/add", conversationHandler.AddMember)
		api.PATCH("/conversations/:_id/members/remove", conversationHandler.RemoveMember)
		api.DELETE("/conversations/:_id", conversationHandler.DeleteConversation)
	}

	// Health check
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	return r
}
