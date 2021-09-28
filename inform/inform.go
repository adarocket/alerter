package inform

import (
	"encoding/json"
	"fmt"
	"github.com/adarocket/alerter/cache"
	"github.com/adarocket/alerter/checker"
	"github.com/adarocket/alerter/database"
	pb "github.com/adarocket/alerter/proto"
	"github.com/tidwall/gjson"
	"log"
)

const msgTemplate = "current value: %g, normal value from %g to %g"

func CheckFieldsOfNode(newNode interface{}, key cache.KeyCache) ([]*pb.SendNotifier, error) {
	cacheInstance := cache.GetCacheInstance()
	oldNode := cacheInstance.GetOldNodeByType(newNode, key)

	newNodeJSON, err := json.Marshal(&newNode)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	alerts, err := database.Sqllite.GetAlerts()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var msges []*pb.SendNotifier
	for _, alert := range alerts {
		value := gjson.Get(string(newNodeJSON), alert.CheckedField)
		if !value.Exists() {
			log.Println("val not exist")
			continue
		}

		alertNode, err := database.Sqllite.GetNodeAlertByID(alert.ID)
		if err != nil {
			log.Println(err)
			continue
		}

		var msg pb.SendNotifier
		var diffVal float64
		switch alert.TypeChecker {
		case checker.IntervalT.String():
			diffVal, err = checker.Checker(alertNode.NormalFrom,
				alertNode.NormalTo, value.String(), alert.TypeChecker)
			if err != nil {
				log.Println(err)
				continue
			}
		case checker.ChangeUpT.String():
			if oldNode == nil {
				continue
			}

			oldNodeJSON, _ := json.Marshal(&oldNode)
			oldValue := gjson.Get(string(oldNodeJSON), alert.CheckedField)
			if !value.Exists() {
				log.Println("val not exist")
				continue
			}

			diffVal, err = checker.Checker(oldValue.Value(), value.String(), nil, alert.TypeChecker)
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

		if diffVal > alertNode.NormalTo || diffVal < alertNode.NormalFrom {
			msg.Value = fmt.Sprintf(msgTemplate,
				diffVal,
				alertNode.NormalFrom,
				alertNode.NormalTo)
			msg.Frequency = alertNode.Frequency
			msg.TypeMessage = fmt.Sprintf("Node %s, uuid: %s info, field %s",
				key.TypeNode, key.Key, alert.CheckedField)

			msges = append(msges, &msg)
		}
	}

	return msges, nil
}
