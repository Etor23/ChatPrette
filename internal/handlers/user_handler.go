package handlers

import (
	"chat-back/internal/dto"
	"chat-back/internal/models"
	"chat-back/internal/repos"
	"net/http"
	"time"

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

func (h *UserHandler) CreateUser(c *gin.Context) {
	var body dto.CreateUserRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := models.User{
		//ID:        body.ID,
		Email:     body.Email,
		Username:  body.Username,
		AvatarURL: body.AvatarURL,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := h.repo.Create(c.Request.Context(), &user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo crear el usuario"})
		return
	}
	response := dto.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		AvatarURL: user.AvatarURL,
	}
	c.JSON(http.StatusCreated, response)
}

func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.repo.FindAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudieron obtener los usuarios"})
		return
	}

	var response []dto.UserResponse
	for _, user := range users {
		response = append(response, dto.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			Username:  user.Username,
			AvatarURL: user.AvatarURL,
		})
	}

	c.JSON(http.StatusOK, response)
}

func (h *UserHandler) GetUserById(c *gin.Context) {
	id := c.Param("_id")

	user, err := h.repo.FindById(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		return
	}

	response := dto.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		AvatarURL: user.AvatarURL,
	}
	c.JSON(http.StatusOK, response)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("_id")
	var body dto.UpdateUserRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := models.User{
		Email:     body.Email,
		Username:  body.Username,
		AvatarURL: body.AvatarURL,
		UpdatedAt: time.Now(),
	}

	err := h.repo.Update(c.Request.Context(), id, &user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		return
	}

	// Obtener el usuario actualizado para retornarlo
	updatedUser, err := h.repo.FindById(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener usuario actualizado"})
		return
	}

	response := dto.UserResponse{
		ID:        updatedUser.ID,
		Email:     updatedUser.Email,
		Username:  updatedUser.Username,
		AvatarURL: updatedUser.AvatarURL,
	}
	c.JSON(http.StatusOK, response)
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("_id")

	err := h.repo.Delete(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Usuario eliminado exitosamente"})
}