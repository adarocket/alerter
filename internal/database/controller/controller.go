package controller

import (
	"database/sql"
	"github.com/adarocket/alerter/internal/database/db"
)

var alertNodeController AlertNode
var alertController Alert

// InitializeControllerInstances - init alert and alertNode instances
func InitializeControllerInstances(database *sql.DB) {
	alertController.db = db.NewAlertInstance(database)
	alertNodeController.db = db.NewAlertNodeInstance(database)
}

func GetAlertControllerInstance() Alert {
	return alertController
}

func GetAlertNodeControllerInstance() AlertNode {
	return alertNodeController
}
