package sqllite

import (
	"github.com/adarocket/alerter/database"
	"log"
)

const getAlerts = `
	SELECT id, name, checked_field, type_checker
	FROM alerts
`

func (p sqlite) GetAlerts() ([]database.Alerts, error) {
	rows, err := p.dbConn.Query(getAlerts)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var alerts []database.Alerts

	for rows.Next() {
		alert := database.Alerts{}
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

func (p sqlite) GetAlertByID(id int64) (database.Alerts, error) {
	rows, err := p.dbConn.Query(getAlertByID, id)
	if err != nil {
		log.Println(err)
		return database.Alerts{}, err
	}
	defer rows.Close()

	var alert database.Alerts

	for rows.Next() {
		err = rows.Scan(&alert.ID, &alert.Name, &alert.CheckedField, &alert.TypeChecker)
		if err != nil {
			log.Println(err)
			return database.Alerts{}, err
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

func (p sqlite) UpdateAlert(alert database.Alerts) error {
	_, err := p.dbConn.Exec(updateAlert,
		alert.Name, alert.CheckedField,
		alert.TypeChecker, alert.ID)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

const createAlert = `
	INSERT INTO alerts
	(name, checked_field, type_checker, id)
	VALUES ($1,$2, $3, $4)
`

func (p sqlite) CreateAlert(alert database.Alerts) error {
	_, err := p.dbConn.Exec(createAlert,
		alert.Name, alert.CheckedField,
		alert.TypeChecker, alert.ID)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

const deleteAlert = `
	DELETE FROM alerts
	WHERE id = $1
`

func (p sqlite) DeleteAlert(alertID int64) error {
	_, err := p.dbConn.Exec(deleteAlert, alertID)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
