package database

import (
	"github.com/adarocket/alerter/database/structs"
	"log"
)

const getAlerts = `
	SELECT id, name, checked_field, type_checker
	FROM alerts
`

func (p sqlite) GetAlerts() ([]structs.AlertsTable, error) {
	rows, err := p.dbConn.Query(getAlerts)
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

const getAlertByID = `
	SELECT id, name, checked_field, type_checker
	FROM alerts
	WHERE id = $1
`

func (p sqlite) GetAlertByID(id int64) (structs.AlertsTable, error) {
	rows, err := p.dbConn.Query(getAlertByID, id)
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

const updateAlert = `
	UPDATE alerts
	SET (name, checked_field, type_checker)
		= ($1,$2,$3)
	WHERE id = $4
`

func (p sqlite) UpdateAlert(alert structs.AlertsTable) error {
	_, err := p.dbConn.Exec(updateAlert,
		alert.Name, alert.CheckedField,
		alert.TypeChecker, alert.ID)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
