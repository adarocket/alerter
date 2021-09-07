package main

import (
	"embed"
	"github.com/adarocket/alerter/web"
	"log"
)

//go:embed data/*.html
var webUI embed.FS

func main() {
	/*log.SetFlags(log.Lshortfile)
	database.InitDatabase()
	database.Sqllite.CreateTables()
	database.Sqllite.FillTables()
	database.Sqllite.GetDataFromAlerts()
	database.Sqllite.GetDataFromAlertNode(0)

	nodesinfo.StartTracking()*/
	log.SetFlags(log.Lshortfile)
	web.WebUI = webUI
	web.StartServer()
}
