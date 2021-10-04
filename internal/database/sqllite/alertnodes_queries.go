package sqllite

import (
	"github.com/adarocket/alerter/internal/database"
	"log"
)

const updateAlertNode = `
	INSERT OR REPLACE INTO 
	alert_node (alert_id, normal_from, normal_to, critical_from, critical_to, frequncy, node_uuid)
	VALUES ($1,$2,$3,$4,$5,$6,$7)
`

func (p sqlite) UpdateAlertNode(alertNode database.AlertNode) error {
	_, err := p.dbConn.Exec(updateAlertNode,
		alertNode.AlertID, alertNode.NormalFrom, alertNode.NormalTo,
		alertNode.CriticalFrom, alertNode.CriticalTo, alertNode.Frequency, alertNode.NodeUuid)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

const getAlertNodes = `
	SELECT alert_id, normal_from, normal_to, critical_from, critical_to, frequncy, node_uuid 
	FROM alert_node
	WHERE alert_id = $1
`

func (p sqlite) GetAlertNodesByID(alertId int64) ([]database.AlertNode, error) {
	rows, err := p.dbConn.Query(getAlertNodes, alertId)
	if err != nil {
		log.Println(err)
		return []database.AlertNode{}, err
	}
	defer rows.Close()

	alertNodes := make([]database.AlertNode, 0, 10)

	for rows.Next() {
		var alertNode database.AlertNode
		err = rows.Scan(&alertNode.AlertID, &alertNode.NormalFrom,
			&alertNode.NormalTo, &alertNode.CriticalFrom,
			&alertNode.CriticalTo, &alertNode.Frequency, &alertNode.NodeUuid)
		if err != nil {
			log.Println(err)
			return alertNodes, err
		}

		alertNodes = append(alertNodes, alertNode)
	}

	return alertNodes, nil
}

const getAlertNode = `
	SELECT alert_id, normal_from, normal_to, critical_from, critical_to, frequncy, node_uuid 
	FROM alert_node
	WHERE alert_id = $1 AND node_uuid = $2
`

func (p sqlite) GetAlertNodeByIdAndNodeUuid(alertId int64, nodeUuid string) (database.AlertNode, error) {
	rows, err := p.dbConn.Query(getAlertNode, alertId, nodeUuid)
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
	(normal_from, normal_to, critical_from, critical_to, frequncy, alert_id, node_uuid)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
`

func (p sqlite) CreateAlertNode(alertNode database.AlertNode) error {
	_, err := p.dbConn.Exec(createAlertNode,
		alertNode.NormalFrom, alertNode.NormalTo, alertNode.CriticalFrom,
		alertNode.CriticalTo, alertNode.Frequency, alertNode.AlertID, alertNode.NodeUuid)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

const deleteAlertNode = `
	DELETE FROM alert_node
	WHERE alert_id = $1 AND node_uuid = $2
`

func (p sqlite) DeleteAlertNode(alertNodeID int64, nodeUuid string) error {
	_, err := p.dbConn.Exec(deleteAlertNode, alertNodeID, nodeUuid)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
