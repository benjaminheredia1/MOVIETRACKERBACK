package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"MovieTrackerBack/internal/infrastructure/postgres"
)

func main() {
	// Cargar variables de entorno, ignorar error si no existe el archivo .env
	_ = godotenv.Load()

	host := os.Getenv("host")
	port := os.Getenv("port")
	user := os.Getenv("user")
	password := os.Getenv("password")
	dbname := os.Getenv("dbname")

	if host == "" {
		log.Println("Advertencia: Variables de entorno de base de datos no encontradas.")
	}

	db, err := postgres.NewConnection(host, port, user, password, dbname)
	if err != nil {
		log.Fatalf("Error conectando a la base de datos: %v", err)
	}
	defer db.Close()

	fmt.Println("Ejecutando seeders...")
	err = seedData(db)
	if err != nil {
		log.Fatalf("Error ejecutando seeders: %v", err)
	}
	fmt.Println("¡Seeders ejecutados correctamente!")
}

func seedData(db *sql.DB) error {
	// Limpiar datos existentes
	_, err := db.Exec(`TRUNCATE TABLE LIST_ITEMS, ITEMS, LISTA RESTART IDENTITY CASCADE`)
	if err != nil {
		return fmt.Errorf("error limpiando tablas: %w", err)
	}

	// 1. Crear listas
	listQuery := `INSERT INTO LISTA (name, description, created_at) VALUES ($1, $2, $3) RETURNING id`

	var favId, watchId int
	now := time.Now()
	
	err = db.QueryRow(listQuery, "Favoritos", "Mis películas y series favoritas", now).Scan(&favId)
	if err != nil {
		return fmt.Errorf("error insertando lista Favoritos: %w", err)
	}

	err = db.QueryRow(listQuery, "Watchlist", "Pendientes por ver", now).Scan(&watchId)
	if err != nil {
		return fmt.Errorf("error insertando lista Watchlist: %w", err)
	}

	// 2. Crear items llenando todas las columnas según la migración
	itemQuery := `
		INSERT INTO ITEMS (
			tmdb_id, adult, backdrop_path, name, original_name, overview, 
			poster_path, media_type, original_language, popularity, first_air_date, 
			softcore, genre_ids, origin_country, vote_average, vote_count, 
			list_id, status, comentary_user, calification_user, watched_at, added_at
		) 
		VALUES (
			$1, $2, $3, $4, $5, $6, 
			$7, $8, $9, $10, $11, 
			$12, $13, $14, $15, $16, 
			$17, $18, $19, $20, $21, $22
		) RETURNING id
	`

	var item1, item2 int

	// The Matrix
	err = db.QueryRow(itemQuery,
		603,                               // tmdb_id
		false,                             // adult
		"/l4QHerS0KWCcE6L009P5Rz6hXqf.jpg", // backdrop_path
		"The Matrix",                      // name
		"The Matrix",                      // original_name
		"Set in the 22nd century, The Matrix tells the story of a computer hacker who joins a group of underground insurgents fighting the vast and powerful computers who now rule the earth.", // overview
		"/f89U3ADr1oiB1s9GkdPOEpXUk5H.jpg", // poster_path
		"movie",                           // media_type
		"en",                              // original_language
		82.1234,                           // popularity
		"1999-03-30",                      // first_air_date
		false,                             // softcore
		"28,878",                          // genre_ids
		"US",                              // origin_country
		8.2,                               // vote_average
		24000,                             // vote_count
		favId,                             // list_id
		"watched",                         // status
		"Obra maestra de la ciencia ficción.", // comentary_user
		10.0,                              // calification_user
		now,                               // watched_at
		now,                               // added_at
	).Scan(&item1)
	if err != nil {
		return fmt.Errorf("error insertando item 1: %w", err)
	}

	// Inception
	err = db.QueryRow(itemQuery,
		27205,                             // tmdb_id
		false,                             // adult
		"/8ZTVqvKDQ8emSGUEMjsS4yHAwrp.jpg", // backdrop_path
		"Inception",                       // name
		"Inception",                       // original_name
		"Cobb, a skilled thief who commits corporate espionage by infiltrating the subconscious of his targets is offered a chance to regain his old life as payment for a task considered to be impossible: \"inception\", the implantation of another person's idea into a target's subconscious.", // overview
		"/oYuLEt3zVCKq57qu2F8dT7NIa6f.jpg", // poster_path
		"movie",                           // media_type
		"en",                              // original_language
		95.5678,                           // popularity
		"2010-07-15",                      // first_air_date
		false,                             // softcore
		"28,878,12",                       // genre_ids
		"US,GB",                           // origin_country
		8.4,                               // vote_average
		35000,                             // vote_count
		watchId,                           // list_id
		"pending",                         // status
		"Me la recomendaron muchísimo.",   // comentary_user
		0.0,                               // calification_user
		now,                               // watched_at
		now,                               // added_at
	).Scan(&item2)
	if err != nil {
		return fmt.Errorf("error insertando item 2: %w", err)
	}

	// 3. Relacionar items con listas (tabla pivot)
	listItemQuery := `INSERT INTO LIST_ITEMS (list_id, item_id, added_at) VALUES ($1, $2, $3)`

	_, err = db.Exec(listItemQuery, favId, item1, now)
	if err != nil {
		return fmt.Errorf("error insertando list_items 1: %w", err)
	}

	_, err = db.Exec(listItemQuery, watchId, item2, now)
	if err != nil {
		return fmt.Errorf("error insertando list_items 2: %w", err)
	}

	return nil
}
