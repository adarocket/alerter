package database

import (
	"database/sql"
	"log"
)

type ModelAlert interface {
	GetAlerts() ([]Alerts, error)
	GetAlertByID(id int64) (Alerts, error)
	CreateAlert(alert Alerts) error
	DeleteAlert(alertID int64) error
	UpdateAlert(table Alerts) error
}

type alertDB struct {
	dbConn *sql.DB
}

func NewAlertInstance(dbConn *sql.DB) ModelAlert {
	return &alertDB{dbConn: dbConn}
}

const getAlerts = `
	SELECT id, name, checked_field, type_checker
	FROM alerts
`

func (p *alertDB) GetAlerts() ([]Alerts, error) {
	rows, err := p.dbConn.Query(getAlerts)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var alerts []Alerts

	for rows.Next() {
		alert := Alerts{}
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

func (p alertDB) GetAlertByID(id int64) (Alerts, error) {
	rows, err := p.dbConn.Query(getAlertByID, id)
	if err != nil {
		log.Println(err)
		return Alerts{}, err
	}
	defer rows.Close()

	var alert Alerts

	for rows.Next() {
		err = rows.Scan(&alert.ID, &alert.Name, &alert.CheckedField, &alert.TypeChecker)
		if err != nil {
			log.Println(err)
			return Alerts{}, err
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

func (p *alertDB) UpdateAlert(alert Alerts) error {
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

func (p *alertDB) CreateAlert(alert Alerts) error {
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

func (p *alertDB) DeleteAlert(alertID int64) error {
	_, err := p.dbConn.Exec(deleteAlert, alertID)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
