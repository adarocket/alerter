package structs

// FIXME какой от этого смысл?
type Database interface {
	GetAlerts() ([]AlertsTable, error)
	GetNodeAlertByID(alertId int64) (AlertNodeTable, error)
	GetAlertByID(id int64) (AlertsTable, error)
	UpdateAlert(table AlertsTable) error
	UpdateAlertNode(table AlertNodeTable) error
	FillTables() error
	CreateTables() error
}
