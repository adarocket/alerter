package inform

import (
	"encoding/json"
	"errors"
	"github.com/adarocket/alerter/cache"
	"github.com/adarocket/alerter/checker"
	"github.com/adarocket/alerter/database"
	"github.com/adarocket/alerter/database/structs"
	"github.com/tidwall/gjson"
	"log"
)

func CheckNodes(newNode interface{}) (string, error) {
	cacheInstance := cache.GetCacheInstance()
	oldNode := cacheInstance.GetOldNodeByType(newNode)

	newNodeJSON, err := json.Marshal(&newNode)
	if err != nil {
		log.Println(err)
		return "", err
	}

	alerts, err := database.Sqllite.GetDataFromAlerts()
	if err != nil {
		log.Println(err)
		return "", err
	}

	for _, alert := range alerts {
		value := gjson.Get(string(newNodeJSON), alert.CheckedField)
		if !value.Exists() {
			log.Println("val not exist")
			continue
		}

		switch alert.TypeChecker {
		case checker.IntervalT.String():
			_, _ = IntervalCreateMsg(alert, value)
		case checker.ChangeUpT.String():
			if oldNode == nil {
				continue
			}
			_, _ = ChangeUpCreateMsg(oldNode, alert, value)
		case checker.ChangeDownT.String():
		case checker.DateT.String():
		default:
			log.Println("undefiend checker type")
		}
	}

	return "", nil
}

func IntervalCreateMsg(alert structs.AlertsTable, value gjson.Result) (string, error) {
	alertNode, err := database.Sqllite.GetDataFromAlertNode(alert.ID)
	if err != nil {
		log.Println(err)
		return "", err
	}
	_, _ = checker.Checker(alertNode.NormalFrom, alertNode.NormalTo, value.String(), alert.TypeChecker)
	return "", nil
}

func ChangeUpCreateMsg(oldNode interface{}, alert structs.AlertsTable, value gjson.Result) (string, error) {
	if oldNode != nil {
		oldNodeJSON, _ := json.Marshal(&oldNode)
		oldValue := gjson.Get(string(oldNodeJSON), alert.CheckedField)
		if !value.Exists() {
			log.Println("")
			return "", errors.New("")
		}
		_, _ = checker.Checker(oldValue.Value(), value.String(), nil, alert.TypeChecker)
	}

	return "", errors.New("old not exist")
}
