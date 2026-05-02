package handlers

import (
	"chat-back/internal/dto"
	"chat-back/internal/repos"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	repo *repos.UserRepo
}

func NewUserHandler(repo *repos.UserRepo) *UserHandler {
	return &UserHandler{
		repo: repo,
	}
}

// GET /api/users/
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.repo.FindAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudieron obtener los usuarios"})
		return
	}

	var response []dto.UserResponse
	for _, user := range users {
		response = append(response, dto.UserResponse{
			ID:        user.ID.Hex(),
			Email:     user.Email,
			Username:  user.Username,
			AvatarURL: user.AvatarURL,
		})
	}

	c.JSON(http.StatusOK, response)
}

// GET /api/users/:_id
func (h *UserHandler) GetUserById(c *gin.Context) {
	id := c.Param("_id")
	user, err := h.repo.FindById(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		return
	}

	response := dto.UserResponse{
		ID:        user.ID.Hex(),
		Email:     user.Email,
		Username:  user.Username,
		AvatarURL: user.AvatarURL,
	}

	c.JSON(http.StatusOK, response)
}


// GET /api/users/search?q=...&limit=...
func (h *UserHandler) SearchUsers(c *gin.Context) {
	query := strings.TrimSpace(c.Query("q"))
	if query == "" {
		c.JSON(http.StatusOK, []dto.UserResponse{})
		return
	}

	limit := 10
	if rawLimit := c.Query("limit"); rawLimit != "" {
		parsed, err := strconv.Atoi(rawLimit)
		if err != nil || parsed <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "limit debe ser un entero positivo"})
			return
		}
		if parsed > 20 {
			parsed = 20
		}
		limit = parsed
	}

	users, err := h.repo.SearchByUsername(c.Request.Context(), query, int64(limit))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudieron buscar usuarios"})
		return
	}

	response := make([]dto.UserResponse, 0, len(users))
	for _, user := range users {
		response = append(response, dto.UserResponse{
			ID:        user.ID.Hex(),
			Email:     user.Email,
			Username:  user.Username,
			AvatarURL: user.AvatarURL,
		})
	}

	c.JSON(http.StatusOK, response)
}
// CreateUser y GetUserByEmail ELIMINADOS:
// - CreateUser era un agujero de seguridad (sin auth, sin firebase_uid)
// - GetUserByEmail no estaba registrado en ninguna ruta
