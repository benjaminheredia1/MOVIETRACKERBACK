package http

import (
	"MovieTrackerBack/internal/infrastructure/http/handlers"

	"github.com/gin-gonic/gin"
)

func NewRouter(
	itemsHandler *handlers.ItemsHandler,
	listaHandler *handlers.ListaHandler,
) *gin.Engine {
	r := gin.Default()
	api := r.Group("/api")

	// Items
	api.GET("/items", itemsHandler.GetAll)
	api.GET("/items/:id", itemsHandler.GetByID)
	api.POST("/items", itemsHandler.Add)
	api.PATCH("/items/:id/watched", itemsHandler.MarkAsWatched)
	api.DELETE("/items/:id", itemsHandler.Delete)

	// Listas
	api.GET("/lists", listaHandler.GetAll)
	api.GET("/lists/:id", listaHandler.GetByID)
	api.POST("/lists", listaHandler.Add)
	api.DELETE("/lists/:id", listaHandler.Delete)

	// Busqueda TMDB
	api.GET("/search", itemsHandler.Search)

	// Chat
	api.POST("/chat", itemsHandler.Chat)

	return r
}
