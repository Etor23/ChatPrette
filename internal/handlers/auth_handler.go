package handlers

import (
	"net/http"
	"time"

	"chat-back/internal/auth"
	"chat-back/internal/dto"
	"chat-back/internal/models"
	"chat-back/internal/repos"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthHandler struct {
	repo *repos.UserRepo
}

func NewAuthHandler(repo *repos.UserRepo) *AuthHandler {
	return &AuthHandler{repo: repo}
}

// POST /api/auth/register
// El usuario YA se registró en Firebase (frontend).
// Aquí creamos su perfil en nuestra base de datos.
func (h *AuthHandler) Register(c *gin.Context) {
	// 1. Parsear body
	var body dto.RegisterRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. Datos que vienen del middleware (Firebase ya verificó el token)
	firebaseUID := auth.GetFirebaseUID(c)
	email := auth.GetEmail(c)

	// 3. ¿Ya se registró antes?
	exists, err := h.repo.ExistsByFirebaseUID(c.Request.Context(), firebaseUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al verificar usuario"})
		return
	}
	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "Este usuario ya está registrado"})
		return
	}

	// 4. ¿Username ya tomado?
	taken, err := h.repo.ExistsByUsername(c.Request.Context(), body.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al verificar username"})
		return
	}
	if taken {
		c.JSON(http.StatusConflict, gin.H{"error": "Ese username ya está en uso"})
		return
	}

	// 5. ¿Email ya existe?
	emailTaken, err := h.repo.ExistsByEmail(c.Request.Context(), email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al verificar email"})
		return
	}
	if emailTaken {
		c.JSON(http.StatusConflict, gin.H{"error": "Ese email ya está registrado"})
		return
	}

	// 6. Crear usuario
	user := &models.User{
		ID:          primitive.NewObjectID(),
		FirebaseUID: firebaseUID,
		Email:       email,
		Username:    body.Username,
		AvatarURL:   body.AvatarURL,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := h.repo.Create(c.Request.Context(), user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo crear el usuario"})
		return
	}

	// 7. Responder
	c.JSON(http.StatusCreated, dto.LoginResponse{
		User: dto.UserResponse{
			ID:        user.ID.Hex(),
			Email:     user.Email,
			Username:  user.Username,
			AvatarURL: user.AvatarURL,
		},
		IsNew: true,
	})
}

// POST /api/auth/login
// El usuario inicia sesión con Firebase (frontend).
// Aquí verificamos que existe en nuestra BD y retornamos su perfil.
func (h *AuthHandler) Login(c *gin.Context) {
	// 1. Datos del middleware
	firebaseUID := auth.GetFirebaseUID(c)

	// 2. Buscar en nuestra BD
	user, err := h.repo.FindByFirebaseUID(c.Request.Context(), firebaseUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al buscar usuario"})
		return
	}

	// 3. Si no existe, necesita registrarse primero
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "not_registered",
			"message": "Usuario no registrado. Regístrate primero.",
		})
		return
	}

	// 4. Retornar perfil
	c.JSON(http.StatusOK, dto.LoginResponse{
		User: dto.UserResponse{
			ID:        user.ID.Hex(),
			Email:     user.Email,
			Username:  user.Username,
			AvatarURL: user.AvatarURL,
		},
		IsNew: false,
	})
}

// GET /api/auth/me
// Retorna el perfil del usuario autenticado.
func (h *AuthHandler) GetMe(c *gin.Context) {
	firebaseUID := auth.GetFirebaseUID(c)

	user, err := h.repo.FindByFirebaseUID(c.Request.Context(), firebaseUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al buscar usuario"})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		return
	}

	c.JSON(http.StatusOK, dto.UserResponse{
		ID:        user.ID.Hex(),
		Email:     user.Email,
		Username:  user.Username,
		AvatarURL: user.AvatarURL,
	})
}
