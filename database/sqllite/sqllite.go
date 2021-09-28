package sqllite

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var Sqllite sqlite

// InitDatabase ...
func InitDatabase(sqlLitePathDB string) {
	db, err := sql.Open("sqlite3", sqlLitePathDB)
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
