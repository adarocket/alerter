package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

// InitDatabase ...
func InitDatabase(sqlLitePathDB string) *sql.DB {
	db, err := sql.Open("sqlite3", sqlLitePathDB)
	if err != nil {
		log.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	return db
}
