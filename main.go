package main

import (
	"embed"
)

//go:embed data/*.html
var webUI embed.FS

func main() {
	//in := controller.InitializeControllerInstance()

	/*log.SetFlags(log.Lshortfile)

	database.InitDatabase("sqlDB.db")
	database.CreateTables()
	//database.FillTables()

	web.WebUI = webUI
	go func() {
		nodesinfo.StartTracking()
	}()
	web.StartServer(":8080")*/
}
