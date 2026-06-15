package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func RunMigrations(database *sql.DB) error {

	_, err := database.Exec(`
		CREATE TABLE IF NOT EXISTS migrations (
			id SERIAL PRIMARY KEY,
			name TEXT UNIQUE NOT NULL,
			applied_at TIMESTAMP DEFAULT NOW()
		)
	`)
	if err != nil {
		return err
	}

	migrationPath := "internal/db/migrations"

	files, err := os.ReadDir(migrationPath)
	if err != nil {
		return err
	}

	for _, file := range files {

		if !strings.HasSuffix(file.Name(), ".up.sql") {
			continue
		}

		var exists bool

		err = database.QueryRow(
			"SELECT EXISTS (SELECT 1 FROM migrations WHERE name=$1)",
			file.Name(),
		).Scan(&exists)

		if err != nil {
			return err
		}

		if exists {
			fmt.Println("Skipping:", file.Name())
			continue
		}

		filePath := filepath.Join(migrationPath, file.Name())

		sqlBytes, err := os.ReadFile(filePath)
		if err != nil {
			return err
		}

		fmt.Println("Running migration:", file.Name())

		_, err = database.Exec(string(sqlBytes))
		if err != nil {
			return err
		}

		_, err = database.Exec(
			"INSERT INTO migrations(name) VALUES($1)",
			file.Name(),
		)
		if err != nil {
			return err
		}

		fmt.Println("Applied:", file.Name())
	}

	return nil
}