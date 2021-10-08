package controller

import (
	"database/sql"
	"github.com/adarocket/alerter/internal/database"
)

var AlertNodeController AlertNode
var AlertController Alert

func InitializeControllerInstances(db *sql.DB) {
	AlertController.db = database.NewAlertInstance(db)
	AlertNodeController.db = database.NewAlertNodeInstance(db)
}

func GetAlertControllerInstance() Alert {
	return AlertController
}

func GetAlertNodeControllerInstance() AlertNode {
	return AlertNodeController
}
