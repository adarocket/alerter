package nodesinfo

import (
	"encoding/json"
	"fmt"
	"github.com/adarocket/alerter/internal/cache"
	"github.com/adarocket/alerter/internal/controller"
	"github.com/adarocket/alerter/internal/nodesinfo/checker"
	"github.com/adarocket/alerter/internal/nodesinfo/msgsender"
	pb "github.com/adarocket/proto/proto-gen/notifier"
	"github.com/tidwall/gjson"
	"log"
	"time"
)

const msgTemplateNormal = "current value: %s, normal value from %g to %g"
const msgTemplateCritical = "current value: %s, critical value from %g to %g"
const msgTemplateType = "Node %s, field %s"

type MsgNodeField struct {
	NodeUuid  string
	NodeField string
	*pb.SendNotifier
}

func CheckFieldsOfNode(newNode interface{}, key cache.KeyCache) (map[msgsender.KeyMsgSender]*pb.SendNotifier, error) {
	cacheInstance := cache.GetCacheInstance()
	oldNode := cacheInstance.GetOldNodeByType(newNode, key)

	newNodeJSON, err := json.Marshal(&newNode)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	db := controller.GetControllerInstance().AlertNode
	alerts, err := db.GetCrossAlertNodeAndAlert(key.Key)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	messages := map[msgsender.KeyMsgSender]*pb.SendNotifier{}
	for _, alert := range alerts {
		value := gjson.Get(string(newNodeJSON), alert.CheckedField)
		if !value.Exists() {
			log.Println("val not exist")
			continue
		}

		var oldValue interface{}
		if oldNode != nil {
			oldNodeJSON, _ := json.Marshal(&oldNode)
			oldVal := gjson.Get(string(oldNodeJSON), alert.CheckedField)
			if !value.Exists() {
				log.Println("val not exist")
				continue
			}
			oldValue = oldVal.String()
		}

		msg := pb.SendNotifier{}
		var diffVal float64

		switch alert.TypeChecker {
		case checker.IntervalT.String():
			diffVal, err = checker.Checker(alert.NormalFrom,
				alert.NormalTo, value.String(), alert.TypeChecker)
			if err != nil {
				log.Println(err)
				continue
			}
		case checker.ChangeUpT.String():
			if oldNode == nil {
				continue
			}
			diffVal, err = checker.Checker(oldValue, value.String(), nil, alert.TypeChecker)
			if err != nil {
				log.Println(err)
				continue
			}
		case checker.ChangeDownT.String():
		case checker.DateT.String():
			tm, err := time.Parse("2006-01-02 15:04:05 -0700 MST", value.String())
			if err != nil {
				log.Println(err)
				continue
			}
			diffVal, err = checker.Checker(tm, nil, nil, alert.TypeChecker)
			if err != nil {
				log.Println(err)
				continue
			}
		default:
			log.Println("undefined checker type")
			continue
		}

		if diffVal > alert.CriticalTo || diffVal < alert.CriticalFrom {
			msg.Frequency = msgsender.Normal.String()
			msg.Value = fmt.Sprintf(msgTemplateCritical, value.String(), alert.NormalFrom, alert.NormalTo)
		} else if diffVal > alert.NormalTo || diffVal < alert.NormalFrom {
			msg.Value = fmt.Sprintf(msgTemplateNormal, value.String(), alert.NormalFrom, alert.NormalTo)
			msg.Frequency = alert.Frequency
		} else {
			continue
		}

		msg.TypeMessage = fmt.Sprintf(msgTemplateType, key.Key, alert.CheckedField)
		messagesKey := msgsender.KeyMsgSender{
			NodeUuid:  key.Key,
			NodeField: alert.CheckedField,
		}

		messages[messagesKey] = &msg
	}

	return messages, nil
}
