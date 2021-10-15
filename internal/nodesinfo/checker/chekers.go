package checker

import (
	"errors"
	"log"
	"math"
	"strconv"
	"time"
)

func IntervalTest(bottomLine, topLine, val interface{}) (float64, error) {
	if bottomLine == nil || topLine == nil || val == nil {
		return 0, errors.New("on of the args is nil")
	}

	_, _, valF, err := parseToFloat64(nil, nil, val)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	return valF, nil
}

func ChangeUpTest(oldVal, newVal interface{}) (float64, error) {
	if oldVal == nil || newVal == nil {
		return 0, errors.New("on of the args is nil")
	}

	oldValF, newValF, _, err := parseToFloat64(oldVal, newVal, nil)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	if newValF-oldValF < 0 {
		return 0, nil
	} else {
		return newValF - oldValF, nil
	}
}

func ChangeDownTest(oldVal, newVal interface{}) (float64, error) {
	if oldVal == nil || newVal == nil {
		return 0, errors.New("on of the args is nil")
	}

	oldValF, newValF, _, err := parseToFloat64(oldVal, newVal, nil)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	if newValF-oldValF > 0 {
		return 0, nil
	} else {
		return newValF - oldValF, nil
	}
}

func DateCheckTest(val interface{}) (float64, error) {
	if val == nil {
		return 0, errors.New("on of the args is nil")
	}

	var dat time.Time
	if a, isTr := val.(time.Time); isTr {
		dat = a
	} else {
		return 0, errors.New("invalid date type")
	}

	return float64(dat.Unix() - time.Now().Unix()), nil
}

func EqualCheckTest(currVal, equalVal interface{}) (float64, error) {
	if currVal == nil || equalVal == nil {
		return 0, errors.New("on of the args is nil")
	}

	currValF, equalValF, _, err := parseToFloat64(currVal, equalVal, nil)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	return math.Abs(currValF - equalValF), nil
}

func parseToFloat64(val1, val2, val3 interface{}) (val1F float64,
	val2F float64, val3F float64, err error) {

	if val1 != nil {
		if str, ok := val1.(string); ok {
			val1F, _ = strconv.ParseFloat(str, 64)
		} else if val1F, ok = val1.(float64); ok {
		} else {
			err = errors.New("invalid type")
			return
		}
	}
	if val2 != nil {
		if str, ok := val2.(string); ok {
			val2F, _ = strconv.ParseFloat(str, 64)
		} else if val2F, ok = val2.(float64); ok {
		} else {
			err = errors.New("invalid type")
			return
		}
	}
	if val3 != nil {
		if str, ok := val3.(string); ok {
			val3F, _ = strconv.ParseFloat(str, 64)
		} else if val3F, ok = val3.(float64); ok {
		} else {
			err = errors.New("invalid type")
			return
		}
	}

	return
}
