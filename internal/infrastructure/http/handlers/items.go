package handlers

import (
    "net/http"
    "strconv"
    "github.com/gin-gonic/gin"
    "MovieTrackerBack/internal/application"
    "MovieTrackerBack/internal/domain"
)

type ItemsHandler struct {
    service *application.ItemsListService
}

func NewItemsHandler(service *application.ItemsListService) *ItemsHandler {
    return &ItemsHandler{service: service}
}

// GET /items
func (h *ItemsHandler) GetAll(c *gin.Context) {
    filters := domain.Filters{
        Status:    c.Query("status"),     // ?status=pending
        MediaType: c.Query("media_type"), // ?media_type=movie
        OrderBy:   c.Query("order_by"),   // ?order_by=added_at
        OrderDir:  c.Query("order_dir"),  // ?order_dir=DESC
    }

    items, err := h.service.GetAll(filters)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, items)
}

// GET /items/:id
func (h *ItemsHandler) GetByID(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "id invalido"})
        return
    }

    item, err := h.service.GetByID(id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "item no encontrado"})
        return
    }

    c.JSON(http.StatusOK, item)
}

// POST /items
func (h *ItemsHandler) Add(c *gin.Context) {
    var item domain.ITEM
    if err := c.ShouldBindJSON(&item); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := h.service.AddItem(item); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, item)
}

// PATCH /items/:id/watched
func (h *ItemsHandler) MarkAsWatched(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "id invalido"})
        return
    }

    var body struct {
        Rating     int    `json:"rating"`
        Commentary string `json:"commentary"`
    }
    if err := c.ShouldBindJSON(&body); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := h.service.MarkAsWatched(id, body.Rating, body.Commentary); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "marcado como visto"})
}

// DELETE /items/:id
func (h *ItemsHandler) Delete(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "id invalido"})
        return
    }

    if err := h.service.Delete(id); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "eliminado correctamente"})
}

// GET /search
func (h *ItemsHandler) Search(c *gin.Context) {
    // esto lo implementamos cuando hagamos el cliente TMDB
}

// POST /chat
func (h *ItemsHandler) Chat(c *gin.Context) {
    // esto lo implementamos cuando hagamos el agente
}