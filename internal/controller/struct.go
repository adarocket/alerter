package controller

import (
	"github.com/adarocket/alerter/internal/config"
	"github.com/adarocket/alerter/internal/database"
	"log"
)

var instance StructController

type StructController struct {
	AlertNode alertNode
	Alert     alert
}

func InitializeControllerInstance() StructController {
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatalln(err)
	}

	db := database.InitDatabase(conf.SqlLitePathDB)
	instance.Alert.db = database.NewAlertInstance(db)
	instance.AlertNode.db = database.NewAlertNodeInstance(db)

	return instance
}

func GetControllerInstance() StructController {
	return instance
}
