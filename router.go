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
	conversationRepo := repos.NewConversationRepo(db)

	// ===== Handlers =====
	userHandler := handlers.NewUserHandler(userRepo)
	conversationHandler := handlers.NewConversationHandler(conversationRepo)

	// ===== Routes =====
	api := r.Group("/api")
	{
		//Users
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