package database

var Db Database

type Database interface {
	GetAlerts() ([]Alerts, error)
	GetAlertNodeByIdAndNodeUuid(alertId int64, nodeUuid string) (AlertNode, error)
	GetAlertNodesByID(alertId int64) ([]AlertNode, error)
	GetAlertByID(id int64) (Alerts, error)
	CreateAlertNode(alertNode AlertNode) error
	DeleteAlertNode(alertNodeID int64) error
	CreateAlert(alert Alerts) error
	DeleteAlert(alertID int64) error
	UpdateAlert(table Alerts) error
	UpdateAlertNode(table AlertNode) error
	FillTables() error
	CreateTables() error
}
