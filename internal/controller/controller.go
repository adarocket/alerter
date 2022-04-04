package controller

import (
	"database/sql"
	"github.com/adarocket/alerter/internal/database/model"
)

var alertNodeController AlertNode
var alertController Alert

// InitializeControllerInstances - init alert and alertNode instances
func InitializeControllerInstances(database *sql.DB) {
	alertController.db = model.NewAlertInstance(database)
	alertNodeController.db = model.NewAlertNodeInstance(database)
}

func GetAlertControllerInstance() Alert {
	return alertController
}

func GetAlertNodeControllerInstance() AlertNode {
	return alertNodeController
}
