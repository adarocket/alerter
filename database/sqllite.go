package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var Sqllite sqlite
var SqlLitePathDB = "sqlDB.db"

// InitDatabase ...
func InitDatabase() {
	db, err := sql.Open("sqlite3", SqlLitePathDB)
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	Sqllite.dbConn = db
}

type sqlite struct {
	dbConn *sql.DB
}
