package controller

import (
	"github.com/adarocket/alerter/internal/database/db"
	"log"
)

type Alert struct {
	db db.ModelAlert
}

func (a *Alert) GetAlertByID(alertID int64) (db.Alerts, error) {
	objs, err := a.db.GetAlertByID(alertID)
	if err != nil {
		log.Println(err)
		return objs, err
	}

	return objs, nil
}

func (a *Alert) GetAlerts() ([]db.Alerts, error) {
	objs, err := a.db.GetAlerts()
	if err != nil {
		log.Println(err)
		return objs, err
	}

	return objs, nil
}

func (a *Alert) CreateAlert(alert db.Alerts) error {
	err := a.db.CreateAlert(alert)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (a *Alert) UpdateAlert(alert db.Alerts) error {
	err := a.db.UpdateAlert(alert)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (a *Alert) DeleteAlert(alertID int64) error {
	err := a.db.DeleteAlert(alertID)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
