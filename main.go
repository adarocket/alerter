package main

import (
	"github.com/adarocket/alerter/database"
	"github.com/adarocket/alerter/nodesinfo"
	"log"
)

func main() {
	log.SetFlags(log.Lshortfile)
	database.InitDatabase()
	/*database.Sqllite.CreateTables()
	database.Sqllite.FillTables()
	database.Sqllite.GetDataFromAlerts()
	database.Sqllite.GetDataFromAlertNode(0)*/

	nodesinfo.StartTracking()
}
