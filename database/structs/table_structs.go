package structs

type AlertsTable struct {
	ID           int64
	Name         string
	CheckedField string
	TypeChecker  string
}

type AlertNodeTable struct {
	AlertID      int64
	NormalFrom   float64
	NormalTo     float64
	CriticalFrom float64
	CriticalTo   float64
	Frequency    string
}
