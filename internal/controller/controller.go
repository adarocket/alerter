package controller

import (
	"database/sql"
	"github.com/adarocket/alerter/internal/database"
)

var alertNodeController AlertNode
var alertController Alert

// InitializeControllerInstances - init alert and alertNode instances
func InitializeControllerInstances(db *sql.DB) {
	alertController.db = database.NewAlertInstance(db)
	alertNodeController.db = database.NewAlertNodeInstance(db)
}

func GetAlertControllerInstance() Alert {
	return alertController
}

func GetAlertNodeControllerInstance() AlertNode {
	return alertNodeController
}
