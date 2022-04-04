package controller

import (
	"github.com/adarocket/alerter/internal/database/model"
	"log"
)

type AlertNode struct {
	db model.ModelAlertNode
}

func (c *AlertNode) GetAlertNodesByID(alertNodeID int64) ([]model.AlertNode, error) {
	objs, err := c.db.GetAlertNodesByID(alertNodeID)
	if err != nil {
		log.Println(err)
		return objs, err
	}

	return objs, nil
}

func (c *AlertNode) GetAlertsByNodeUuid(nodeUuid string) ([]model.AlertNodeAndAlert, error) {
	objs, err := c.db.GetAlertsByNodeUuid(nodeUuid)
	if err != nil {
		log.Println(err)
		return objs, err
	}

	return objs, nil
}

func (c *AlertNode) GetAlertNodeByIdAndNodeUuid(alertId int64, nodeUuid string) (model.AlertNode, error) {
	objs, err := c.db.GetAlertNodeByIdAndNodeUuid(alertId, nodeUuid)
	if err != nil {
		log.Println(err)
		return objs, err
	}

	return objs, nil
}

func (c *AlertNode) CreateAlertNode(alertNode model.AlertNode) error {
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

func (c *AlertNode) UpdateAlertNode(alertNode model.AlertNode) error {
	err := c.db.UpdateAlertNode(alertNode)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
