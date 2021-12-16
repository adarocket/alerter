package nodesinfo

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Knetic/govaluate"
	"github.com/adarocket/alerter/internal/database/db"
	"github.com/adarocket/alerter/internal/msgsender"
	"github.com/adarocket/alerter/internal/nodesinfo/checker"
	pb "github.com/adarocket/proto/proto-gen/notifier"
	"github.com/tidwall/gjson"
	"log"
	"time"
)

const msgTemplateType = "Node %s, field %s, cheker type: %s"

// CheckFieldsOfNode - checking fields of node and create notifiers messages
func CheckFieldsOfNode(newNode interface{}, oldNode interface{},
	alerts []db.AlertNodeAndAlert) (map[msgsender.KeyMsg]msgsender.BodyMsg, error) {

	newNodeJSON, err := json.Marshal(&newNode)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	messages := map[msgsender.KeyMsg]msgsender.BodyMsg{}
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

		diffVal, err := calculateDiffVal(alert, oldValue, value)
		if err != nil {
			log.Println(err)
			continue
		}

		msg, err := createMsg(alert, diffVal, value)
		if err != nil {
			log.Println(err)
			continue
		}

		messagesKey := msgsender.KeyMsg{
			NodeUuid:  alert.NodeUuid,
			NodeField: alert.CheckedField,
		}

		messages[messagesKey] = msg
	}

	return messages, nil
}

// example : formula = diffVal ** 0.5
func calculateFrequency(diffValue float64, formula string) (float64, error) {
	expression, err := govaluate.NewEvaluableExpression(formula)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	parameters := make(map[string]interface{}, 8)
	parameters["diffVal"] = diffValue

	result, err := expression.Evaluate(parameters)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	frequency, _, _, err := checker.ParseToFloat64(result, nil, nil)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	//fmt.Println(reflect.TypeOf(result), result)

	return frequency, nil
}

func calculateDiffVal(alert db.AlertNodeAndAlert, oldValue string, value gjson.Result) (float64, error) {
	var diffVal float64
	var err error

	switch alert.TypeChecker {
	case checker.IntervalT.String():
		diffVal, err = checker.IntervalTest(alert.NormalFrom,
			alert.NormalTo, value.String())
		if err != nil {
			log.Println(err)
			return 0, err
		}
	case checker.ChangeUpT.String():
		if oldValue == "" {
			return 0, errors.New("old val does not exist")
		}
		diffVal, err = checker.ChangeUpTest(oldValue, value.String())
		if err != nil {
			log.Println(err)
			return 0, err
		}
	case checker.ChangeDownT.String():
	case checker.DateT.String():
		tm, err := time.Parse("2006-01-02 15:04:05 -0700 MST", value.String())
		if err != nil {
			log.Println(err)
			return 0, err
		}
		diffVal, err = checker.DateCheckTest(tm)
		if err != nil {
			log.Println(err)
			return 0, err
		}
	case checker.EqualT.String():
		if oldValue == "" {
			return 0, errors.New("old val does not exist")
		}
		diffVal, err = checker.EqualCheckTest(value.String(), oldValue)
		if err != nil {
			log.Println(err)
			return 0, err
		}
	case checker.MoreT.String():
		diffVal, _, _, err = checker.ParseToFloat64(value.String(), nil, nil)
		if err != nil {
			log.Println(err)
			return 0, err
		}
	default:
		log.Println("undefined checker type")
		return 0, errors.New("undefined checker type")
	}

	return diffVal, nil
}

func createMsg(alert db.AlertNodeAndAlert,
	diffVal float64, value gjson.Result) (msgsender.BodyMsg, error) {
	msg := msgsender.BodyMsg{
		Notify: &pb.SendNotifier{},
	}

	if diffVal > alert.CriticalTo || diffVal < alert.CriticalFrom {
		msg.Frequency = msgsender.MaxFrequency
		msg.Notify.From = fmt.Sprintf("%f", alert.CriticalFrom)
		msg.Notify.To = fmt.Sprintf("%f", alert.CriticalTo)
	} else if diffVal > alert.NormalTo || diffVal < alert.NormalFrom {
		frequency, err := calculateFrequency(diffVal, alert.Frequency)
		if err != nil {
			log.Println(err)
			return msgsender.BodyMsg{}, err
		}
		msg.Frequency = int64(frequency)
		msg.Notify.From = fmt.Sprintf("%f", alert.NormalFrom)
		msg.Notify.To = fmt.Sprintf("%f", alert.NormalTo)
	} else {
		return msgsender.BodyMsg{},
			errors.New("something goes wrong while creating notify")
	}

	msg.Notify.CurrentVal = value.String()
	msg.Notify.TextMessage = fmt.Sprintf(msgTemplateType, alert.NodeUuid, alert.CheckedField, alert.TypeChecker)

	return msg, nil
}
