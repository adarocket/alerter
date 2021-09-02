package structs

type Database interface {
	GetDataFromAlerts() ([]AlertsTable, error)
	GetDataFromAlertNode(alertId int64) ([]AlertNodeTable, error)
	FillTables() error
	CreateTables() error
}
