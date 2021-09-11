package main

import (
	"embed"
	"github.com/adarocket/alerter/database"
	"github.com/adarocket/alerter/nodesinfo"
	"github.com/adarocket/alerter/web"
	"log"
)

//go:embed data/*.html
var webUI embed.FS

func main() {
	log.SetFlags(log.Lshortfile)

	database.InitDatabase()
	web.WebUI = webUI

	go func() {
		web.StartServer()
	}()

	nodesinfo.StartTracking()
}
