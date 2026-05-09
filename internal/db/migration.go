package db

import (
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// Migrate up the db based on the url and files in the directory ./migrations
func RunMigrations(dbURL string) error {
	// 1. Create a new migration object
	m, err := migrate.New("file://migrations", dbURL)
	if err != nil {
		return fmt.Errorf("migration init failed: %w", err)
	}
	defer m.Close()

	// 2. Lounch teh command to migrate up all files in migrations .up.sql in order
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration failed: %w", err)
	}
	log.Println("Migrations applied succesfully")
	return nil
}
