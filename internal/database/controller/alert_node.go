package controller

import (
	"github.com/adarocket/alerter/internal/database/db"
	"log"
)

type AlertNode struct {
	db db.ModelAlertNode
}

func (c *AlertNode) GetAlertNodesByID(alertNodeID int64) ([]db.AlertNode, error) {
	objs, err := c.db.GetAlertNodesByID(alertNodeID)
	if err != nil {
		log.Println(err)
		return objs, err
	}

	return objs, nil
}

func (c *AlertNode) GetAlertsByNodeUuid(nodeUuid string) ([]db.AlertNodeAndAlert, error) {
	objs, err := c.db.GetAlertsByNodeUuid(nodeUuid)
	if err != nil {
		log.Println(err)
		return objs, err
	}

	return objs, nil
}

func (c *AlertNode) GetAlertNodeByIdAndNodeUuid(alertId int64, nodeUuid string) (db.AlertNode, error) {
	objs, err := c.db.GetAlertNodeByIdAndNodeUuid(alertId, nodeUuid)
	if err != nil {
		log.Println(err)
		return objs, err
	}

	return objs, nil
}

func (c *AlertNode) CreateAlertNode(alertNode db.AlertNode) error {
	err := c.db.CreateAlertNode(alertNode)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (c *AlertNode) DeleteAlertNode(alertNodeID int64, nodeUuid string) error {
	err := c.db.DeleteAlertNode(alertNodeID, nodeUuid)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (c *AlertNode) UpdateAlertNode(alertNode db.AlertNode) error {
	err := c.db.UpdateAlertNode(alertNode)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
