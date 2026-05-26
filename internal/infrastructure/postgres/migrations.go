package postgres

import (
	"fmt"
	"log"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(host, port, user, password, dbname string) {
	databaseURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user, password, host, port, dbname,
	)

	m, err := migrate.New(
		"file://db/migrations",
		databaseURL,
	)
	if err != nil {
		log.Fatalf("error creando migrador: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("error corriendo migraciones: %v", err)
	}

	log.Println("migraciones aplicadas correctamente")
}
