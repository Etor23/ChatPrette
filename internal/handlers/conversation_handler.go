package handlers

import (
	"chat-back/internal/dto"
	"chat-back/internal/models"
	"chat-back/internal/repos"
	"chat-back/internal/helpers"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ConversationHandler struct {
	repo *repos.ConversationRepo
}

func NewConversationHandler(repo *repos.ConversationRepo) *ConversationHandler {
	return &ConversationHandler{
		repo: repo,
	}
}

func (h *ConversationHandler) CreateDm(c *gin.Context) {
	var body dto.CreateDmRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	currentUserID := c.GetString("userID") // Obtener el ID del usuario autenticado
	conversation := models.Conversation{
		Type: "dm",
		Members: []string{currentUserID, body.OtherUserID},
	}

	err := h.repo.Create(c.Request.Context(), &conversation)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo crear el DM"})
	}
	response := dto.ConversationResponse{
		ID: conversation.ID.Hex(),
		Type: conversation.Type,
		Members: conversation.Members,
	}
	c.JSON(http.StatusCreated, response)
	
} 

func (h *ConversationHandler) CreateGroup(c *gin.Context) {
	var body dto.CreateGroupRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	currentUserID := c.GetString("userID") // Obtener el ID del usuario autenticado
	members := append([]string{currentUserID}, body.Members...) //Obtiene los miembros del grupo y agrega el usuario actual al inicio de la lista
	conversation := models.Conversation{
		Type: "group",
		Name: body.Name,
		Members: members,
	}

	err := h.repo.Create(c.Request.Context(), &conversation)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo crear el grupo"})
	}
	response := dto.ConversationResponse{
		ID: conversation.ID.Hex(),
		Type: conversation.Type,
		Name: conversation.Name,
		Members: conversation.Members,
	}
	c.JSON(http.StatusCreated, response)
	
} 

func (h *ConversationHandler) GetUserConversations(c *gin.Context) {
	currentUserID := c.GetString("userID") // Obtener el ID del usuario autenticado
	conversations, err := h.repo.FindByMember(c.Request.Context(), currentUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudieron obtener las conversaciones"})
		return
	}
	var response []dto.ConversationResponse
	for _, conv := range conversations {
		response = append(response, dto.ConversationResponse{
			ID: conv.ID.Hex(),
			Type: conv.Type,
			Name: conv.Name,
			Members: conv.Members,
			//TODO: Agregar LastMessageAt y LastMessagePreview una vez que este lo de messages
		})
	}
	c.JSON(http.StatusOK, response)
} 

func (h *ConversationHandler) GetConversationById(c *gin.Context) {
	id := c.Param("_id")
	conversation, err := h.repo.FindById(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo obtener la conversación"})
		return
	}
	response := dto.ConversationResponse{
		ID: conversation.ID.Hex(),
		Type: conversation.Type,
		Name: conversation.Name,
		Members: conversation.Members,
	}
	c.JSON(http.StatusOK, response)
} 

func (h *ConversationHandler) UpdateGroupName(c *gin.Context) {
	id := c.Param("_id")

	conversationID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de conversación inválido"})
		return
	}
	var body dto.UpdateGroupNameRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	currentUserID := c.GetString("userID") // Obtener el ID del usuario autenticado

	conversation, err := h.repo.FindById(c.Request.Context(), conversationID.Hex())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Conversación no encontrada"})
		return
	}

	if conversation.Type != "group" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Solo se pueden actualizar el nombre en grupos"})
		return
	}

	if !helpers.Contains(conversation.Members, currentUserID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "not a member of this group"})
		return
	}

	err = h.repo.UpdateName(c.Request.Context(), conversationID.Hex(), body.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo actualizar el nombre del grupo"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Nombre del grupo actualizado"})

}

func (h *ConversationHandler) AddMember(c *gin.Context) {
	id := c.Param("_id")

	conversationID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de conversación inválido"})
		return
	}
	var body dto.AddGroupMembersRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	currentUserID := c.GetString("userID") // Obtener el ID del usuario autenticado

	conversation, err := h.repo.FindById(c.Request.Context(), conversationID.Hex())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Conversación no encontrada"})
		return
	}

	if conversation.Type != "group" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Solo se pueden agregar miembros en grupos"})
		return
	}

	if !helpers.Contains(conversation.Members, currentUserID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "not a member of this group"})
		return
	}

	err = h.repo.AddMembers(c.Request.Context(), conversationID.Hex(), body.Members)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo agregar miembros al grupo"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Miembros agregados al grupo"})

}

func (h *ConversationHandler) RemoveMember(c *gin.Context) {
	id := c.Param("_id")

	conversationID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de conversación inválido"})
		return
	}
	var body dto.RemoveGroupMembersRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	currentUserID := c.GetString("userID") // Obtener el ID del usuario autenticado

	conversation, err := h.repo.FindById(c.Request.Context(), conversationID.Hex())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Conversación no encontrada"})
		return
	}

	if conversation.Type != "group" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Solo se pueden remover miembros en grupos"})
		return
	}

	if !helpers.Contains(conversation.Members, currentUserID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "not a member of this group"})
		return
	}

	err = h.repo.RemoveMembers(c.Request.Context(), conversationID.Hex(), body.Members)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo remover miembros del grupo"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Miembros removidos del grupo"})

}

func (h *ConversationHandler) DeleteConversation(c *gin.Context) {
	id := c.Param("_id")

	conversationID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de conversación inválido"})
		return
	}

	currentUserID := c.GetString("userID") // Obtener el ID del usuario autenticado

	conversation, err := h.repo.FindById(c.Request.Context(), conversationID.Hex())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Conversación no encontrada"})
		return
	}

	if !helpers.Contains(conversation.Members, currentUserID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "No permitido eliminar una conversación de la que no eres miembro"})
		return
	}

	err = h.repo.Delete(c.Request.Context(), conversationID.Hex())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo eliminar la conversación"})
		return
	}

	//TODO: Eliminar los mensajes de la conversación una vez que este lo de messages este implementado

	c.JSON(http.StatusOK, gin.H{"message": "Conversación eliminada correctamente"})

}