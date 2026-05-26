package handlers

import (
	"MovieTrackerBack/internal/application"
	"MovieTrackerBack/internal/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
	chatService *application.ChatService
}

func NewChatHandler(chatService *application.ChatService) *ChatHandler {
	return &ChatHandler{chatService: chatService}
}

// POST /chat
func (h *ChatHandler) Chat(c *gin.Context) {
	var msg domain.ChatMessage
	if err := c.ShouldBindJSON(&msg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reply, err := h.chatService.SaveMessage(msg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"reply": reply})
}
