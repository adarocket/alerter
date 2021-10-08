package controller

import (
	"github.com/adarocket/alerter/internal/database"
	"log"
)

type Alert struct {
	db database.ModelAlert
}

func (a *Alert) GetAlertByID(alertID int64) (database.Alerts, error) {
	objs, err := a.db.GetAlertByID(alertID)
	if err != nil {
		log.Println(err)
		return objs, err
	}

	return objs, nil
}

func (a *Alert) GetAlerts() ([]database.Alerts, error) {
	objs, err := a.db.GetAlerts()
	if err != nil {
		log.Println(err)
		return objs, err
	}

	return objs, nil
}

func (a *Alert) CreateAlert(alert database.Alerts) error {
	err := a.db.CreateAlert(alert)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (a *Alert) UpdateAlert(alert database.Alerts) error {
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
