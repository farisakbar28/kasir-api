package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func InitDB(connectionString string) (*sql.DB, error) {
	// Buka koneksi ke database
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	// Test koneksi — pastikan database bisa dijangkau
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	// Atur connection pool
	db.SetMaxOpenConns(25) // maksimal 25 koneksi sekaligus
	db.SetMaxIdleConns(5)  // maksimal 5 koneksi idle

	log.Println("Database connected successfully")
	return db, nil
}