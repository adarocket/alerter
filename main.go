package main

import (
	"embed"
	"github.com/adarocket/alerter/internal/config"
	"github.com/adarocket/alerter/internal/controller"
	"github.com/adarocket/alerter/internal/database"
	"github.com/adarocket/alerter/internal/nodesinfo"
	"github.com/adarocket/alerter/internal/web"
	"log"
)

//go:embed tmpl/*.html
var webUI embed.FS

func main() {
	log.SetFlags(log.Lshortfile)
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	dbConn := database.InitDatabase(conf.SqlLitePathDB)
	controller.InitializeControllerInstances(dbConn)
	//database.CreateTables()
	//database.FillTables()

	web.WebUI = webUI
	go func() {
		nodesinfo.StartTracking()
	}()
	web.StartServer(conf.WebServerAddr)
}
