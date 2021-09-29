package sqllite

import (
	"github.com/adarocket/alerter/database"
	"log"
)

const updateAlertNode = `
	INSERT OR REPLACE INTO 
	alert_node (alert_id, normal_from, normal_to, critical_from, critical_to, frequncy)
	VALUES ($1,$2,$3,$4,$5,$6)
`

func (p sqlite) UpdateAlertNode(alertNode database.AlertNode) error {
	_, err := p.dbConn.Exec(updateAlertNode,
		alertNode.AlertID, alertNode.NormalFrom, alertNode.NormalTo,
		alertNode.CriticalFrom, alertNode.CriticalTo, alertNode.Frequency)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

const getAlertNodes = `
	SELECT alert_id, normal_from, normal_to, critical_from, critical_to, frequncy 
	FROM alert_node
	WHERE alert_id = $1
`

func (p sqlite) GetNodeAlertByID(alertId int64) (database.AlertNode, error) {
	rows, err := p.dbConn.Query(getAlertNodes, alertId)
	if err != nil {
		log.Println(err)
		return database.AlertNode{}, err
	}
	defer rows.Close()

	var alertNode database.AlertNode

	for rows.Next() {
		err = rows.Scan(&alertNode.AlertID, &alertNode.NormalFrom,
			&alertNode.NormalTo, &alertNode.CriticalFrom, &alertNode.CriticalTo, &alertNode.Frequency)
		if err != nil {
			log.Println(err)
			return database.AlertNode{}, err
		}
	}

	return alertNode, nil
}

const createAlertNode = `
	INSERT INTO alert_node
	(normal_from, normal_to, critical_from, critical_to, frequncy, alert_id)
	VALUES ($1, $2, $3, $4, $5, $6)
`

func (p sqlite) CreateAlertNode(alertNode database.AlertNode) error {
	_, err := p.dbConn.Exec(createAlertNode,
		alertNode.NormalFrom, alertNode.NormalTo, alertNode.CriticalFrom,
		alertNode.CriticalTo, alertNode.Frequency, alertNode.AlertID)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

const deleteAlertNode = `
	DELETE FROM alert_node
	WHERE alert_id = $1
`

func (p sqlite) DeleteAlertNode(alertNodeID int64) error {
	_, err := p.dbConn.Exec(deleteAlertNode, alertNodeID)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
