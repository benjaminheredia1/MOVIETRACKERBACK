package main

import (
	"fmt"
	"os"

	"MovieTrackerBack/internal/infrastructure/http"
	"MovieTrackerBack/internal/infrastructure/postgres"
)

func main() {
	// TODO: Inicializar servicios y handlers reales
	// itemsHandler := handlers.NewItemsHandler(...)
	// listaHandler := handlers.NewListaHandler(...)

	router := http.NewRouter(nil, nil) // Pasando nil temporalmente

	db, err := postgres.NewConnection(os.Getenv("host"), os.Getenv("port"), os.Getenv("user"), os.Getenv("password"), os.Getenv("dbname"))
	if err != nil {
		panic(err)
	}
	defer func() {
		fmt.Println("cerrando conexion a la base de datos")
		db.Close()
	}()

	router.Run(os.Getenv("app_port"))


}
