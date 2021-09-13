package structs

type Database interface {
	GetDataFromAlerts() ([]AlertsTable, error)
	GetDataFromAlertNode(alertId int64) (AlertNodeTable, error)
	GetDataFromAlert(id int64) (AlertsTable, error)
	UpdateDataInAlertsTable(table AlertsTable) error
	UpdateDataInAlertNode(table AlertNodeTable) error
	FillTables() error
	CreateTables() error
}
