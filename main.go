package main

import (
	"embed"
	"github.com/adarocket/alerter/internal"
	"github.com/adarocket/alerter/internal/config"
	"github.com/adarocket/alerter/internal/database/controller"
	"github.com/adarocket/alerter/internal/database/model"
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

	dbConn := model.InitDatabase(conf.SqlLitePathDB)

	controller.InitializeControllerInstances(dbConn)
	//database.CreateTables()
	//database.FillTables()

	internal.StartTracking(conf, dbConn)
	//web.WebUI = webUI
	//web.StartServer(conf.WebServerAddr)
}
