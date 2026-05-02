package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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
		log.Fatal("Error conectando a MongoDB:", err)
	}

	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		secretKey = "your-secret-key-change-in-production"
	}
	jwtManager := auth.NewJWTManager(secretKey, 24*time.Hour)

	r := SetupRouter(mongoConn.Database, jwtManager)

	// Crear servidor HTTP manualmente para poder hacer shutdown graceful
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	// Lanzar servidor en goroutine
	go func() {
		log.Printf("Servidor corriendo en http://localhost:%s\n", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error en servidor: %v\n", err)
		}
	}()

	// Esperar señal de cierre (Ctrl+C, kill, etc.)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Apagando servidor...")

	// Dar 5 segundos para terminar requests en curso
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Cerrar servidor HTTP
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Shutdown forzado:", err)
	}

	// Desconectar MongoDB
	if err := mongoConn.Client.Disconnect(ctx); err != nil {
		log.Fatal("Error desconectando MongoDB:", err)
	}

	log.Println("Servidor apagado correctamente")
}
