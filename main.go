package main

import (
	"embed"
	"github.com/adarocket/alerter/internal/database"
	"github.com/adarocket/alerter/internal/database/sqllite"
	"github.com/adarocket/alerter/internal/web"
	"log"
)

//go:embed data/*.html
var webUI embed.FS

func main() {
	log.SetFlags(log.Lshortfile)

	sqllite.InitDatabase("sqlDB.db")
	sqllite.Sqllite.CreateTables()
	database.Db = sqllite.Sqllite

	web.WebUI = webUI
	web.StartServer(":8080")
}
