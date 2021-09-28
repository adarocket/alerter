package main

import (
	"embed"
	"github.com/adarocket/alerter/database"
	"github.com/adarocket/alerter/web"
	"log"
)

//go:embed data/*.html
var webUI embed.FS

func main() {
	log.SetFlags(log.Lshortfile)

	database.InitDatabase()
	database.Sqllite.CreateTables()
	database.Sqllite.FillTables()

	web.WebUI = webUI
	web.StartServer(":8080")
}
