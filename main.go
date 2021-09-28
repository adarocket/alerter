package main

import (
	"embed"
	"github.com/adarocket/alerter/database"
	"github.com/adarocket/alerter/database/sqllite"
	"github.com/adarocket/alerter/web"
	"log"
)

//go:embed data/*.html
var webUI embed.FS

func main() {
	log.SetFlags(log.Lshortfile)

	sqllite.InitDatabase("sqlDB.db")
	sqllite.Sqllite.CreateTables()
	sqllite.Sqllite.FillTables()
	database.Db = sqllite.Sqllite

	web.WebUI = webUI
	web.StartServer(":5400")
}
