package controller

import (
	"github.com/adarocket/alerter/internal/database"
	"log"
)

type alert struct {
	db database.ModelAlert
}

func (a *alert) GetAlertByID(alertID int64) (database.Alerts, error) {
	objs, err := a.db.GetAlertByID(alertID)
	if err != nil {
		log.Println(err)
		return objs, err
	}

	return objs, nil
}

func (a *alert) GetAlerts() ([]database.Alerts, error) {
	objs, err := a.db.GetAlerts()
	if err != nil {
		log.Println(err)
		return objs, err
	}

	return objs, nil
}

func (a *alert) CreateAlert(alert database.Alerts) error {
	err := a.db.CreateAlert(alert)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (a *alert) UpdateAlert(alert database.Alerts) error {
	err := a.db.UpdateAlert(alert)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (a *alert) DeleteAlert(alertID int64) error {
	err := a.db.DeleteAlert(alertID)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
