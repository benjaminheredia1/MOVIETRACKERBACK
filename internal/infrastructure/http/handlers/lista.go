package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

type ListaHandler struct {
}

func NewListaHandler() *ListaHandler {
	return &ListaHandler{}
}

func (h *ListaHandler) GetAll(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented"})
}

func (h *ListaHandler) GetByID(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented"})
}

func (h *ListaHandler) Add(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented"})
}

func (h *ListaHandler) Delete(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented"})
}
