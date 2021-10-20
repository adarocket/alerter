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
	"time"
)

const msgTemplateType = "Node %s, field %s, cheker type: %s"

func CheckFieldsOfNode(newNode interface{}, key cache.KeyCache,
	alerts []database.AlertNodeAndAlert) (map[msgsender.KeyMsgSender]msgsender.ValueMsgSender, error) {

	cacheInstance := cache.GetCacheInstance()
	oldNode := cacheInstance.GetOldNodeByType(newNode, key)

	newNodeJSON, err := json.Marshal(&newNode)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	messages := map[msgsender.KeyMsgSender]msgsender.ValueMsgSender{}
	for _, alert := range alerts {
		value := gjson.Get(string(newNodeJSON), alert.CheckedField)
		if !value.Exists() {
			log.Println("val not exist")
			continue
		}

		var oldValue string
		if oldNode != nil {
			oldNodeJSON, _ := json.Marshal(&oldNode)
			oldVal := gjson.Get(string(oldNodeJSON), alert.CheckedField)
			if !value.Exists() {
				log.Println("val not exist")
				continue
			}
			oldValue = oldVal.String()
		}

		msg := msgsender.ValueMsgSender{
			Notify: &pb.SendNotifier{},
		}
		var diffVal float64

		switch alert.TypeChecker {
		case checker.IntervalT.String():
			diffVal, err = checker.IntervalTest(alert.NormalFrom,
				alert.NormalTo, value.String())
			if err != nil {
				log.Println(err)
				continue
			}
		case checker.ChangeUpT.String():
			if oldValue == "" {
				continue
			}
			diffVal, err = checker.ChangeUpTest(oldValue, value.String())
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
			diffVal, err = checker.DateCheckTest(tm)
			if err != nil {
				log.Println(err)
				continue
			}
		case checker.EqualT.String():
			if oldValue == "" {
				continue
			}
			diffVal, err = checker.EqualCheckTest(value.String(), oldValue)
			if err != nil {
				log.Println(err)
				continue
			}
		case checker.MoreT.String():
			diffVal, _, _, err = checker.ParseToFloat64(value.String(), nil, nil)
			if err != nil {
				log.Println(err)
				continue
			}
		default:
			log.Println("undefined checker type")
			continue
		}

		if diffVal > alert.CriticalTo || diffVal < alert.CriticalFrom {
			msg.Frequency = msgsender.Max.String()
			msg.Notify.From = fmt.Sprintf("%f", alert.CriticalFrom)
			msg.Notify.To = fmt.Sprintf("%f", alert.CriticalTo)
		} else if diffVal > alert.NormalTo || diffVal < alert.NormalFrom {
			msg.Frequency = alert.Frequency
			msg.Notify.From = fmt.Sprintf("%f", alert.NormalFrom)
			msg.Notify.To = fmt.Sprintf("%f", alert.NormalTo)
		} else {
			continue
		}

		msg.Notify.CurrentVal = value.String()
		msg.Notify.TextMessage = fmt.Sprintf(msgTemplateType, key.Key, alert.CheckedField, alert.TypeChecker)

		messagesKey := msgsender.KeyMsgSender{
			NodeUuid:  key.Key,
			NodeField: alert.CheckedField,
		}

		messages[messagesKey] = msg
	}

	return messages, nil
}
