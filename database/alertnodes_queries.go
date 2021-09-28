package database

import (
	"github.com/adarocket/alerter/database/structs"
	"log"
)

const updateAlertNode = `
	UPDATE alert_node
	SET (normal_from, normal_to, critical_from, critical_to, frequncy)
		= ($1,$2,$3,$4,$5)
	WHERE alert_id = $6
`

func (p sqlite) UpdateAlertNode(alertNode structs.AlertNodeTable) error {
	_, err := p.dbConn.Exec(updateAlertNode,
		alertNode.NormalFrom, alertNode.NormalTo, alertNode.CriticalFrom,
		alertNode.CriticalTo, alertNode.Frequency, alertNode.AlertID)

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

func (p sqlite) GetNodeAlertByID(alertId int64) (structs.AlertNodeTable, error) {
	rows, err := p.dbConn.Query(getAlertNodes, alertId)
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
