package inform

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/adarocket/alerter/cache"
	"github.com/adarocket/alerter/checker"
	"github.com/adarocket/alerter/client"
	"github.com/adarocket/alerter/database"
	"github.com/adarocket/alerter/database/structs"
	pb "github.com/adarocket/alerter/proto"
	"github.com/tidwall/gjson"
	"log"
)

const msgTemplate = "current value: %g, normal value from %g to %g"

func CheckNodes(newNode interface{}, notifyClient *client.NotifierClient) error {
	cacheInstance := cache.GetCacheInstance()
	oldNode := cacheInstance.GetOldNodeByType(newNode)

	newNodeJSON, err := json.Marshal(&newNode)
	if err != nil {
		log.Println(err)
		return err
	}

	alerts, err := database.Sqllite.GetDataFromAlerts()
	if err != nil {
		log.Println(err)
		return err
	}

	for _, alert := range alerts {
		value := gjson.Get(string(newNodeJSON), alert.CheckedField)
		if !value.Exists() {
			log.Println("val not exist")
			continue
		}

		var msg *pb.Request
		switch alert.TypeChecker {
		case checker.IntervalT.String():
			msg, err = IntervalCreateMsg(alert, value)
			if err != nil {
				log.Println(err)
				continue
			}
		case checker.ChangeUpT.String():
			if oldNode == nil {
				continue
			}

			msg, err = ChangeUpCreateMsg(oldNode, alert, value)
			if err != nil {
				log.Println(err)
				continue
			}
		case checker.ChangeDownT.String():
		case checker.DateT.String():
		default:
			log.Println("undefiend checker type")
			continue
		}

		if msg.Value != "" {
			msg.TypeMessage = fmt.Sprintf("Node %s info, field %s",
				alert.Name, alert.CheckedField)
			if err := notifyClient.SendMessage(msg); err != nil {
				log.Println(err)
				continue
			}
		}
	}

	return nil
}

func IntervalCreateMsg(alert structs.AlertsTable, value gjson.Result) (*pb.Request, error) {
	alertNode, err := database.Sqllite.GetDataFromAlertNode(alert.ID)
	if err != nil {
		log.Println(err)
		return &pb.Request{}, err
	}

	msg := pb.Request{}
	diffVal, _ := checker.Checker(alertNode.NormalFrom, alertNode.NormalTo, value.String(), alert.TypeChecker)

	if diffVal > alertNode.NormalTo && diffVal < alertNode.NormalFrom {
		msg.Value = fmt.Sprintf(msgTemplate,
			diffVal,
			alertNode.NormalFrom,
			alertNode.NormalTo)

		msg.Frequency = alertNode.Frequency
	}

	return &msg, nil
}

func ChangeUpCreateMsg(oldNode interface{}, alert structs.AlertsTable, value gjson.Result) (*pb.Request, error) {
	alertNode, err := database.Sqllite.GetDataFromAlertNode(alert.ID)
	if err != nil {
		log.Println(err)
		return &pb.Request{}, err
	}

	oldNodeJSON, _ := json.Marshal(&oldNode)
	oldValue := gjson.Get(string(oldNodeJSON), alert.CheckedField)
	if !value.Exists() {
		log.Println("val not exist")
		return &pb.Request{}, errors.New("val not exist")
	}

	msg := pb.Request{}
	diffVal, _ := checker.Checker(oldValue.Value(), value.String(), nil, alert.TypeChecker)

	if diffVal > alertNode.NormalTo && diffVal < alertNode.NormalFrom {
		msg.Value = fmt.Sprintf(msgTemplate,
			diffVal,
			alertNode.NormalFrom,
			alertNode.NormalTo)

		msg.Frequency = alertNode.Frequency
	}

	return &msg, nil
}
