package main

import (
	"embed"
	"github.com/adarocket/alerter/internal/database"
	"github.com/adarocket/alerter/internal/database/sqllite"
	"github.com/adarocket/alerter/internal/nodesinfo"
	"github.com/adarocket/alerter/internal/web"
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
	go func() {
		nodesinfo.StartTracking()
	}()
	web.StartServer(":8080")
}
