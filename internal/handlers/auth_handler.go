package handlers

// La lógica es: el frontend hace login con Firebase y te manda el id_token.
//  Tú lo verificas, extraes el email/uid y buscas o creas el usuario en MongoDB (patrón upsert).
import (
	"context"
	"net/http"
	"time"

	"chat-back/internal/dto"
	"chat-back/internal/models"
	"chat-back/internal/repos"

	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthHandler struct {
	authClient *auth.Client
	userRepo   *repos.UserRepo
}

func NewAuthHandler(authClient *auth.Client, userRepo *repos.UserRepo) *AuthHandler {
	return &AuthHandler{
		authClient: authClient,
		userRepo:   userRepo,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var body dto.LoginRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verificar token con Firebase
	token, err := h.authClient.VerifyIDToken(context.Background(), body.IDToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
		return
	}

	email, _ := token.Claims["email"].(string)
	name, _ := token.Claims["name"].(string)
	picture, _ := token.Claims["picture"].(string)

	// Buscar usuario en MongoDB
	user, err := h.userRepo.FindByEmail(c.Request.Context(), email)
	isNew := false

	if err == mongo.ErrNoDocuments {
		// Primera vez → crear usuario
		isNew = true
		user = &models.User{
			ID:        primitive.NewObjectID().Hex(),
			Email:     email,
			Username:  name,
			AvatarURL: picture,
			BirthDate: time.Time{}, // El usuario lo actualiza después
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := h.userRepo.Create(c.Request.Context(), user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo crear el usuario"})
			return
		}
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al buscar usuario"})
		return
	}

	c.JSON(http.StatusOK, dto.AuthResponse{
		User: dto.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			Username:  user.Username,
			AvatarURL: user.AvatarURL,
			BirthDate: user.BirthDate.Format("2006-01-02"),
		},
		IsNew: isNew,
	})
}

// Devuelve el perfil del usuario autenticado (requiere middleware)
func (h *AuthHandler) Me(c *gin.Context) {
	email, _ := c.Get("firebase_email")

	user, err := h.userRepo.FindByEmail(c.Request.Context(), email.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		return
	}

	c.JSON(http.StatusOK, dto.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		AvatarURL: user.AvatarURL,
		BirthDate: user.BirthDate.Format("2006-01-02"),
	})
}
