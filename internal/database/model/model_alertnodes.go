package model

import (
	"database/sql"
	"errors"
	"log"
)

type ModelAlertNode interface {
	GetAlertNodeByIdAndNodeUuid(alertId int64, nodeUuid string) (AlertNode, error)
	GetAlertNodesByID(alertId int64) ([]AlertNode, error)
	CreateAlertNode(alertNode AlertNode) error
	DeleteAlertNode(alertNodeID int64, nodeUuid string) error
	UpdateAlertNode(table AlertNode) error
	GetAlertsByNodeUuid(nodeUuid string) ([]AlertNodeAndAlert, error)
}

type AlertNodeAndAlert struct {
	Name         string
	CheckedField string
	TypeChecker  string
	AlertID      int64
	NormalFrom   float64
	NormalTo     float64
	CriticalFrom float64
	CriticalTo   float64
	Frequency    string
	NodeUuid     string
}

type alertNodeDB struct {
	dbConn *sql.DB
}

func NewAlertNodeInstance(dbConn *sql.DB) ModelAlertNode {
	return &alertNodeDB{dbConn: dbConn}
}

const getAlertsByNodeUuid = `
	SELECT alert_id, normal_from, normal_to, critical_from, 
	       critical_to, frequncy, node_uuid, name, checked_field, type_checker
	FROM alert_node
	INNER JOIN alerts a on alert_node.alert_id = a.id
	AND  node_uuid = $2
`

func (p *alertNodeDB) GetAlertsByNodeUuid(nodeUuid string) ([]AlertNodeAndAlert, error) {
	rows, err := p.dbConn.Query(getAlertsByNodeUuid, nodeUuid)
	if err != nil {
		log.Println(err)
		return []AlertNodeAndAlert{}, err
	}
	defer rows.Close()

	var nodes []AlertNodeAndAlert
	for rows.Next() {
		var node AlertNodeAndAlert
		err = rows.Scan(&node.AlertID, &node.NormalFrom, &node.NormalTo, &node.CriticalFrom,
			&node.CriticalTo, &node.Frequency, &node.NodeUuid, &node.Name,
			&node.CheckedField, &node.TypeChecker)
		if err != nil {
			log.Println(err)
			return []AlertNodeAndAlert{}, err
		}

		nodes = append(nodes, node)
	}

	return nodes, nil
}

const updateAlertNode = `
	UPDATE alert_node
	SET node_uuid = $1,
	    alert_id = $2,
	    critical_from = $3,
	    critical_to = $4,
	    normal_from = $5,
	    normal_to = $6,
	    frequncy = $7
	WHERE
    	node_uuid = $1 AND alert_id = $2
`

func (p *alertNodeDB) UpdateAlertNode(alertNode AlertNode) error {
	_, err := p.dbConn.Exec(updateAlertNode,
		alertNode.NodeUuid, alertNode.AlertID, alertNode.CriticalFrom, alertNode.CriticalTo,
		alertNode.NormalFrom, alertNode.NormalTo, alertNode.Frequency)

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

func (p *alertNodeDB) GetAlertNodesByID(alertId int64) ([]AlertNode, error) {
	rows, err := p.dbConn.Query(getAlertNodes, alertId)
	if err != nil {
		log.Println(err)
		return []AlertNode{}, err
	}
	defer rows.Close()

	alertNodes := make([]AlertNode, 0, 10)

	for rows.Next() {
		var alertNode AlertNode
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

func (p *alertNodeDB) GetAlertNodeByIdAndNodeUuid(alertId int64, nodeUuid string) (AlertNode, error) {
	rows, err := p.dbConn.Query(getAlertNode, alertId, nodeUuid)
	if err != nil {
		log.Println(err)
		return AlertNode{}, err
	}
	defer rows.Close()

	var alertNode AlertNode

	for rows.Next() {
		err = rows.Scan(&alertNode.AlertID, &alertNode.NormalFrom,
			&alertNode.NormalTo, &alertNode.CriticalFrom,
			&alertNode.CriticalTo, &alertNode.Frequency, &alertNode.NodeUuid)
		if err != nil {
			log.Println(err)
			return AlertNode{}, err
		}

		return alertNode, nil
	}

	return alertNode, errors.New("nothing found")
}

const createAlertNode = `
	INSERT INTO alert_node
	(normal_from, normal_to, critical_from, critical_to, frequncy, alert_id, node_uuid)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
`

func (p *alertNodeDB) CreateAlertNode(alertNode AlertNode) error {
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

func (p *alertNodeDB) DeleteAlertNode(alertNodeID int64, nodeUuid string) error {
	_, err := p.dbConn.Exec(deleteAlertNode, alertNodeID, nodeUuid)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
