package main

import (
	"fmt"
	"os"

	"MovieTrackerBack/internal/application"
	"MovieTrackerBack/internal/infrastructure/http"
	"MovieTrackerBack/internal/infrastructure/http/handlers"
	"MovieTrackerBack/internal/infrastructure/openai"
	"MovieTrackerBack/internal/infrastructure/postgres"
	"MovieTrackerBack/internal/infrastructure/redis"
	"MovieTrackerBack/internal/infrastructure/tmdb"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	host := os.Getenv("host")
	port := os.Getenv("port")
	user := os.Getenv("user")
	password := os.Getenv("password")
	dbname := os.Getenv("dbname")

	postgres.RunMigrations(host, port, user, password, dbname)

	db, err := postgres.NewConnection(host, port, user, password, dbname)
	if err != nil {
		panic(err)
	}
	defer func() {
		fmt.Println("cerrando conexion a la base de datos")
		db.Close()
	}()

	//conexion a redis
	redisClient := redis.NewConnection(os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"), os.Getenv("REDIS_PASSWORD"))
	redisCache := redis.NewCacheRepository(redisClient)

	//Cliente de openai que mas adelante se inyectara
	model := os.Getenv("OPENAI_MODEL")
	if model == "" {
		model = "gpt-4"
	}

	// Inyeccion de dependencias para TMDB
	tmdbClient := tmdb.NewTMDBClient(os.Getenv("TMDB_API_KEY"), os.Getenv("TMDB_BASE_URL"), redisCache)

	//Inyeccion de dependencias para los items
	// Se agrego la inyeccion de TMDBClient al servicio de items para que pueda realizar búsquedas en TMDB
	itemRepo := postgres.NewItemRepository(db)
	itemService := application.NewItemsService(itemRepo, tmdbClient)
	itemsHandler := handlers.NewItemsHandler(itemService)

	//Inyeccion de dependencias para las listas
	listaRepo := postgres.NewListaRepository(db)
	listaService := application.NewListaService(listaRepo)
	listaHandler := handlers.NewListaHandler(listaService)

	//Inyeccion de dependencias para el chat
	openaiClient := openai.NewOpenAIClient(redisCache, os.Getenv("OPENAI_API_KEY"), model, itemRepo, listaRepo)
	chatService := application.NewChatService(openaiClient)
	chatHandler := handlers.NewChatHandler(chatService)

	router := http.NewRouter(itemsHandler, listaHandler, chatHandler)

	router.Run(":" + os.Getenv("app_port"))

}
