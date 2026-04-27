package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"chat-back/internal/auth"
	"chat-back/internal/handlers"
	"chat-back/internal/repos"

	"go.mongodb.org/mongo-driver/mongo"
)

func SetupRouter(db *mongo.Database, firebaseAuth *auth.FirebaseProvider) *gin.Engine {

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	// ===== Repos =====
	userRepo := repos.NewUserRepo(db)

	// ===== Handlers =====
	userHandler := handlers.NewUserHandler(userRepo)
	authHandler := handlers.NewAuthHandler(userRepo) // ← NUEVO

	// ===== Routes =====
	api := r.Group("/api")
	{
		// --- Auth (requieren token de Firebase) ---
		authRoutes := api.Group("/auth")
		authRoutes.Use(auth.Middleware(firebaseAuth)) //  Middleware protege estas rutas
		{
			authRoutes.POST("/register", authHandler.Register)
			authRoutes.POST("/login", authHandler.Login)
			authRoutes.GET("/me", authHandler.GetMe)
		}

		// --- Users (públicas por ahora, luego pueden protegerlas) ---
		api.GET("/users", userHandler.GetAllUsers)
		api.GET("/users/get/:_id", userHandler.GetUserById)
		api.POST("/users", userHandler.CreateUser)

		// Si quieren proteger las rutas de users también:
		// usersRoutes := api.Group("/users")
		// usersRoutes.Use(auth.Middleware(firebaseAuth))
		// {
		//     usersRoutes.GET("/", userHandler.GetAllUsers)
		//     usersRoutes.GET("/:_id", userHandler.GetUserById)
		// }
	}

	// Health check
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	return r
}
