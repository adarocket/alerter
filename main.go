package main

import (
	"embed"
	"fmt"
	"github.com/Knetic/govaluate"
	"reflect"
)

//go:embed tmpl/*.html
var webUI embed.FS

func main() {
	expression, _ := govaluate.NewEvaluableExpression("10 ** 0.5")

	parameters := make(map[string]interface{}, 8)
	parameters["foo"] = -1

	result, _ := expression.Evaluate(parameters)
	fmt.Println(reflect.TypeOf(result), result)

	/*log.SetFlags(log.Lshortfile)
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
		nodesinfo.StartTracking(conf.TimeoutCheck, conf.NotifierAddr)
	}()
	web.StartServer(conf.WebServerAddr)*/
}
