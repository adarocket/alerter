package database

import (
	"github.com/adarocket/alerter/database/structs"
	"log"
)

const getDataFromAlertsTable = `
	SELECT id, name, checked_field, type_checker
	FROM alerts
`

func (p sqlite) GetDataFromAlerts() ([]structs.AlertsTable, error) {
	rows, err := p.dbConn.Query(getDataFromAlertsTable)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var alerts []structs.AlertsTable

	for rows.Next() {
		alert := structs.AlertsTable{}
		err := rows.Scan(&alert.ID, &alert.Name, &alert.CheckedField, &alert.TypeChecker)
		if err != nil {
			log.Println(err)
			continue
		}

		alerts = append(alerts, alert)
	}

	return alerts, nil
}

const getDataFromAlertNodeTable = `
	SELECT alert_id, normal_from, normal_to, critical_from, critical_to, frequncy 
	FROM alert_node
	WHERE alert_id = $1
`

func (p sqlite) GetDataFromAlertNode(alertId int64) (structs.AlertNodeTable, error) {
	rows, err := p.dbConn.Query(getDataFromAlertNodeTable, alertId)
	if err != nil {
		log.Println(err)
		return structs.AlertNodeTable{}, err
	}
	defer rows.Close()

	var alertNode structs.AlertNodeTable

	for rows.Next() {
		err = rows.Scan(&alertNode.AlertID, &alertNode.NormalFrom,
			&alertNode.NormalTo, &alertNode.CriticalFrom, &alertNode.CriticalTo, &alertNode.Frequency)
		if err != nil {
			log.Println(err)
			return structs.AlertNodeTable{}, err
		}
	}

	return alertNode, nil
}

const getDataFromAlert = `
	SELECT id, name, checked_field, type_checker
	FROM alerts
	WHERE id = $1
`

func (p sqlite) GetDataFromAlert(id int64) (structs.AlertsTable, error) {
	rows, err := p.dbConn.Query(getDataFromAlert, id)
	if err != nil {
		log.Println(err)
		return structs.AlertsTable{}, err
	}
	defer rows.Close()

	var alert structs.AlertsTable

	for rows.Next() {
		err = rows.Scan(&alert.ID, &alert.Name, &alert.CheckedField, &alert.TypeChecker)
		if err != nil {
			log.Println(err)
			return structs.AlertsTable{}, err
		}
	}

	return alert, nil
}

const setDataInAlertsTable = `
	update alerts
	set (name, checked_field, type_checker)
		= ($1,$2,$3)
	where id = $4
`

func (p sqlite) SetDataInAlertsTable(table structs.AlertsTable) error {
	_, err := p.dbConn.Exec(setDataInAlertsTable,
		table.Name, table.CheckedField,
		table.TypeChecker, table.ID)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

const setDataInAlertNode = `
	update alert_node
	set (normal_from, normal_to, critical_from, critical_to, frequncy)
		= ($1,$2,$3,$4,$5)
	WHERE alert_id = $6
`

func (p sqlite) SetDataInAlertNode(table structs.AlertNodeTable) error {
	_, err := p.dbConn.Exec(setDataInAlertNode,
		table.NormalFrom, table.NormalTo, table.CriticalFrom,
		table.CriticalTo, table.Frequency, table.AlertID)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

const fillAlertsTable = `
	insert OR IGNORE into alerts 
	(id, name, checked_field, type_checker) 
	VALUES (0, 'sizeCache check', 'SizeCache', 'change_up'),
	       (1, 'cpuState check', 'CpuState', 'interval'),
	       (2, 'blocks check', 'Blocks', 'more')
`

const fillAlertsNodeTable = `
	insert OR IGNORE into alert_node
	(alert_id, normal_from, normal_to, critical_from, critical_to, frequncy)
	VALUES (0, 0.0, 10.0, 10.0, 20.0, 'normal'),
	       (1, 5.0, 16.0, 16.0, 20.0, 'normal'),
	       (2, 0.0, 10.0, 10.0, 20.0, 'normal')
`

func (p sqlite) FillTables() error {
	_, err := p.dbConn.Exec(fillAlertsTable)
	if err != nil {
		log.Fatal(err)
	}

	_, err = p.dbConn.Exec(fillAlertsNodeTable)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
