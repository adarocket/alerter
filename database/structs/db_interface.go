package structs

type Database interface {
	GetAlerts() ([]AlertsTable, error)
	GetNodeAlertByID(alertId int64) (AlertNodeTable, error)
	GetAlertByID(id int64) (AlertsTable, error)
	CreateAlertNode(alertNode AlertNodeTable) error
	DeleteAlertNode(alertNodeID int64) error
	CreateAlert(alert AlertsTable) error
	DeleteAlert(alertID int64) error
	UpdateAlert(table AlertsTable) error
	UpdateAlertNode(table AlertNodeTable) error
	FillTables() error
	CreateTables() error
}
