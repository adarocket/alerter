package nodesinfo

import (
	"encoding/json"
	"fmt"
	"github.com/adarocket/alerter/internal/cache"
	"github.com/adarocket/alerter/internal/database"
	checker2 "github.com/adarocket/alerter/internal/nodesinfo/checker"
	pb "github.com/adarocket/proto/proto-gen/notifier"
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

	alerts, err := database.Db.GetAlerts()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var messages []*pb.SendNotifier
	for _, alert := range alerts {
		value := gjson.Get(string(newNodeJSON), alert.CheckedField)
		if !value.Exists() {
			log.Println("val not exist")
			continue
		}

		alertNode, err := database.Db.GetNodeAlertByID(alert.ID)
		if err != nil {
			log.Println(err)
			continue
		}

		var msg pb.SendNotifier
		var diffVal float64
		switch alert.TypeChecker {
		case checker2.IntervalT.String():
			diffVal, err = checker2.Checker(alertNode.NormalFrom,
				alertNode.NormalTo, value.String(), alert.TypeChecker)
			if err != nil {
				log.Println(err)
				continue
			}
		case checker2.ChangeUpT.String():
			if oldNode == nil {
				continue
			}

			oldNodeJSON, _ := json.Marshal(&oldNode)
			oldValue := gjson.Get(string(oldNodeJSON), alert.CheckedField)
			if !value.Exists() {
				log.Println("val not exist")
				continue
			}

			diffVal, err = checker2.Checker(oldValue.Value(), value.String(), nil, alert.TypeChecker)
			if err != nil {
				log.Println(err)
				continue
			}
		case checker2.ChangeDownT.String():
		case checker2.DateT.String():
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

			messages = append(messages, &msg)
		}
	}

	return messages, nil
}
