package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var dbConn *sql.DB

// InitDatabase ...
func InitDatabase(sqlLitePathDB string) {
	db, err := sql.Open("sqlite3", sqlLitePathDB)
	if err != nil {
		log.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	dbConn = db
}
