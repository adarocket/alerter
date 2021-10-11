package main

import (
	"embed"
	"github.com/adarocket/alerter/internal/controller"
	"github.com/adarocket/alerter/internal/database"
	"github.com/adarocket/alerter/internal/nodesinfo"
	"github.com/adarocket/alerter/internal/web"
	"log"
)

//go:embed data/*.html
var webUI embed.FS

func main() {
	log.SetFlags(log.Lshortfile)

	dbConn := database.InitDatabase("sqlDB.db")
	controller.InitializeControllerInstances(dbConn)
	//database.CreateTables()
	//database.FillTables()

	web.WebUI = webUI
	go func() {
		nodesinfo.StartTracking()
	}()
	web.StartServer(":8080")
}
