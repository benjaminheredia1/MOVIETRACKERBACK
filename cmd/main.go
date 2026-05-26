package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"

	"MovieTrackerBack/internal/application"
	"MovieTrackerBack/internal/infrastructure/http"
	"MovieTrackerBack/internal/infrastructure/http/handlers"
	"MovieTrackerBack/internal/infrastructure/postgres"
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

	//Inyeccion de dependencias
	itemRepo := postgres.NewItemRepository(db)
	itemService := application.NewItemsListService(itemRepo)
	itemsHandler := handlers.NewItemsHandler(itemService)

	router := http.NewRouter(itemsHandler, nil)

	router.Run(":" + os.Getenv("app_port"))

}
