package handlers

import (
	"net/http"
	"os"
	"strings"
	"time"

	"chat-back/internal/auth"
	"chat-back/internal/dto"
	"chat-back/internal/models"
	"chat-back/internal/repos"
	"chat-back/internal/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	repo       *repos.UserRepo
	jwtManager *auth.JWTManager
}

func NewAuthHandler(repo *repos.UserRepo, jwtManager *auth.JWTManager) *AuthHandler {
	// Crear JWT Manager con secret key del environment o un valor por defecto
	if jwtManager == nil {
		secretKey := os.Getenv("JWT_SECRET_KEY")
		if secretKey == "" {
			secretKey = "your-secret-key-change-in-production"
		}
		jwtManager = auth.NewJWTManager(secretKey, 24*time.Hour)
	}

	return &AuthHandler{
		repo:       repo,
		jwtManager: jwtManager,
	}
}

// POST /api/auth/register
// Crea un usuario nuevo con email, password y username
// NO requiere autenticación previa
func (h *AuthHandler) Register(c *gin.Context) {
	var body dto.RegisterRequest

	// 1. Parsear y validar body
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()

	// 2. Verificar si el email ya existe en MongoDB
	exists, err := h.repo.ExistsByEmail(ctx, body.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al verificar email"})
		return
	}
	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "Email ya registrado"})
		return
	}

	// 3. Verificar si el username ya existe en MongoDB
	taken, err := h.repo.ExistsByUsername(ctx, body.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al verificar username"})
		return
	}
	if taken {
		c.JSON(http.StatusConflict, gin.H{"error": "Username ya en uso"})
		return
	}

	// 4. Hash de contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al procesar contraseña"})
		return
	}

	// Parse birthdate if provided
	var birthPtr *time.Time
	if body.Birthdate != "" {
		// try RFC3339 first, then date-only
		if t, err := time.Parse(time.RFC3339, body.Birthdate); err == nil {
			birthPtr = &t
		} else if t2, err2 := time.Parse("2006-01-02", body.Birthdate); err2 == nil {
			birthPtr = &t2
		}
	}

	user := &models.User{
		ID:           primitive.NewObjectID(),
		Email:        body.Email,
		Username:     body.Username,
		PasswordHash: string(hashedPassword),
		AvatarURL:    body.AvatarURL,
		Birthdate:    birthPtr,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := h.repo.Create(ctx, user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear perfil"})
		return
	}

	// 5. Generar JWT
	token, err := h.jwtManager.GenerateToken(user.ID.Hex(), user.Email, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al generar token"})
		return
	}

	// 7. Format birthdate and createdAt for response
	birthStr := utils.FormatDateOrEmpty(user.Birthdate)
	createdStr := utils.FormatTime(user.CreatedAt)

	// 8. Responder
	c.JSON(http.StatusCreated, dto.LoginResponse{
		Token: token,
		User: dto.UserResponse{
			ID:        user.ID.Hex(),
			Email:     user.Email,
			Username:  user.Username,
			AvatarURL: user.AvatarURL,
			Birthdate: birthStr,
			CreatedAt: createdStr,
		},
		IsNew: true,
	})
}

// POST /api/auth/login
// Autentica usuario con email y password
// NO requiere autenticación previa
func (h *AuthHandler) Login(c *gin.Context) {
	var body dto.LoginRequest

	// 1. Parsear y validar body
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()

	// 2. Buscar usuario en MongoDB por email
	user, err := h.repo.FindByEmail(ctx, body.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al buscar usuario"})
		return
	}

	// 3. Si no existe en MongoDB, necesita registrarse primero
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "not_registered",
			"message": "Usuario no registrado. Regístrate primero.",
		})
		return
	}

	// 4. Validar contraseña
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(body.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email o password inválido"})
		return
	}

	// 5. Generar JWT
	token, err := h.jwtManager.GenerateToken(user.ID.Hex(), user.Email, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al generar token"})
		return
	}

	// 6. Format birthdate and createdAt for response
	birthStr := utils.FormatDateOrEmpty(user.Birthdate)
	createdStr := utils.FormatTime(user.CreatedAt)

	// 7. Responder
	c.JSON(http.StatusOK, dto.LoginResponse{
		Token: token,
		User: dto.UserResponse{
			ID:        user.ID.Hex(),
			Email:     user.Email,
			Username:  user.Username,
			AvatarURL: user.AvatarURL,
			Birthdate: birthStr,
			CreatedAt: createdStr,
		},
		IsNew: false,
	})
}

// GET /api/auth/me
// Retorna el perfil del usuario autenticado
// REQUIERE: Authorization header con Bearer <token>
func (h *AuthHandler) GetMe(c *gin.Context) {
	// El middleware ya verificó el token y dejó los datos en contexto
	userID := auth.GetUserID(c)

	user, err := h.repo.FindById(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al buscar usuario"})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		return
	}

	// Format birthdate and createdAt as ISO strings
	birthStr := utils.FormatDateOrEmpty(user.Birthdate)
	createdStr := utils.FormatTime(user.CreatedAt)

	c.JSON(http.StatusOK, dto.UserResponse{
		ID:        user.ID.Hex(),
		Email:     user.Email,
		Username:  user.Username,
		AvatarURL: user.AvatarURL,
		Birthdate: birthStr,
		CreatedAt: createdStr,
	})
}

// PUT /api/auth/me
// Actualiza el perfil del usuario (username, birthdate, avatar_url)
// REQUIERE: Authorization header con Bearer <token>
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	// El middleware ya verificó el token
	userID := auth.GetUserID(c)

	var body dto.UpdateProfileRequest

	// 1. Parsear y validar body
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()

	// 2. Si cambia username, verificar que no esté en uso (excepto el suyo)
	if body.Username != "" {
		currentUser, err := h.repo.FindById(ctx, userID)
		if err != nil || currentUser == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
			return
		}

		// Si el nuevo username es diferente al actual, verificar disponibilidad
		if body.Username != currentUser.Username {
			taken, err := h.repo.ExistsByUsername(ctx, body.Username)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al verificar username"})
				return
			}
			if taken {
				c.JSON(http.StatusConflict, gin.H{"error": "Username ya en uso"})
				return
			}
		}
	}

	// 3. Parsear birthdate si se proporciona
	var birthPtr *time.Time
	if body.Birthdate != "" {
		// try RFC3339 first, then date-only
		if t, err := time.Parse(time.RFC3339, body.Birthdate); err == nil {
			birthPtr = &t
		} else if t2, err2 := time.Parse("2006-01-02", body.Birthdate); err2 == nil {
			birthPtr = &t2
		}
	}

	// 4. Actualizar en la BD
	username := &body.Username
	if body.Username == "" {
		username = nil
	}

	avatarURL := &body.AvatarURL
	if body.AvatarURL == "" {
		avatarURL = nil
	}

	user, err := h.repo.Update(ctx, userID, username, birthPtr, avatarURL)
	if err != nil {
		if strings.Contains(err.Error(), "no encontrado") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar perfil"})
		}
		return
	}

	// 5. Formatear respuesta
	birthStr := utils.FormatDateOrEmpty(user.Birthdate)
	createdStr := utils.FormatTime(user.CreatedAt)

	c.JSON(http.StatusOK, dto.UserResponse{
		ID:        user.ID.Hex(),
		Email:     user.Email,
		Username:  user.Username,
		AvatarURL: user.AvatarURL,
		Birthdate: birthStr,
		CreatedAt: createdStr,
	})
}

// POST /api/auth/logout
// Invalida el token (opcional)
// REQUIERE: Authorization header con Bearer <token>
func (h *AuthHandler) Logout(c *gin.Context) {
	// En una implementación con blacklist, aquí guardarías el token en una lista negra
	// Por ahora, simplemente retornamos éxito - el frontend debe borrar el token

	c.JSON(http.StatusOK, gin.H{"message": "Sesión cerrada correctamente"})
}

// POST /api/auth/refresh
// Refresca el JWT (opcional)
// REQUIERE: Authorization header con Bearer <token> válido
func (h *AuthHandler) Refresh(c *gin.Context) {
	// El middleware ya verificó el token anterior
	userID := auth.GetUserID(c)

	user, err := h.repo.FindById(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al buscar usuario"})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		return
	}

	// Generar nuevo token
	newToken, err := h.jwtManager.GenerateToken(user.ID.Hex(), user.Email, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al generar token"})
		return
	}

	c.JSON(http.StatusOK, dto.RefreshTokenResponse{
		Token:     newToken,
		ExpiresIn: int64((24 * time.Hour).Seconds()), // 86400 segundos
	})
}
