package http

import (
    "github.com/gin-gonic/gin"
    "MovieTrackerBack/internal/infrastructure/http/handlers"
)

func NewRouter(
    itemsHandler *handlers.ItemsHandler,
    listaHandler *handlers.ListaHandler,
) *gin.Engine {
    r := gin.Default()

    // Items
    r.GET("/items", itemsHandler.GetAll)
    r.GET("/items/:id", itemsHandler.GetByID)
    r.POST("/items", itemsHandler.Add)
    r.PATCH("/items/:id/watched", itemsHandler.MarkAsWatched)
    r.DELETE("/items/:id", itemsHandler.Delete)

    // Listas
    r.GET("/lists", listaHandler.GetAll)
    r.GET("/lists/:id", listaHandler.GetByID)
    r.POST("/lists", listaHandler.Add)
    r.DELETE("/lists/:id", listaHandler.Delete)

    // Busqueda TMDB
    r.GET("/search", itemsHandler.Search)

    // Chat
    r.POST("/chat", itemsHandler.Chat)

    return r
}