package nodesinfo

import (
	"encoding/json"
	"fmt"
	"github.com/adarocket/alerter/internal/cache"
	"github.com/adarocket/alerter/internal/database"
	"github.com/adarocket/alerter/internal/nodesinfo/checker"
	"github.com/adarocket/alerter/internal/nodesinfo/msgsender"
	pb "github.com/adarocket/proto/proto-gen/notifier"
	"github.com/tidwall/gjson"
	"log"
)

const msgTemplate = "current value: %g, normal value from %g to %g"

type MsgNodeField struct {
	NodeUuid  string
	NodeField string
	*pb.SendNotifier
}

func CheckFieldsOfNode(newNode interface{}, key cache.KeyCache) ([]MsgNodeField, error) {
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

	var messages []MsgNodeField
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

		msg := MsgNodeField{SendNotifier: &pb.SendNotifier{}}
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

		if diffVal > alertNode.CriticalTo || diffVal < alertNode.CriticalFrom {
			msg.Frequency = msgsender.Normal.String()
		} else if diffVal > alertNode.NormalTo || diffVal < alertNode.NormalFrom {
			msg.Frequency = msgsender.Normal.String()
		}

		msg.Value = fmt.Sprintf(msgTemplate, diffVal, alertNode.NormalFrom, alertNode.NormalTo)
		msg.NodeField = alert.CheckedField
		msg.NodeUuid = key.Key

		messages = append(messages, msg)
	}

	return messages, nil
}
