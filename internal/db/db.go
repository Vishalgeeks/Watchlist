package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func Connect(connStr string) *sql.DB {

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	log.Println("Database connected successfully")

	// AUTO RUN MIGRATIONS
	if err := RunMigrations(db); err != nil {
		log.Fatal("Migration failed:", err)
	}

	log.Println("Migrations completed successfully")

	return db
}