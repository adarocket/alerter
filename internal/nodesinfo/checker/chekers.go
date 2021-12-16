package checker

import (
	"errors"
	"log"
	"math"
	"strconv"
	"time"
)

// rebuild it
func IntervalTest(bottomLine, topLine float64, val string) (float64, error) {
	if val == "" {
		return 0, errors.New("on of the args is nil")
	}

	valF, err := strconv.ParseFloat(val, 64)

	if err != nil {
		log.Println(err)
		return 0, err
	}

	return valF, nil
}

// ChangeUpTest - returns 0 if newVal - oldVal < 0, in other cases diff
func ChangeUpTest(oldVal, newVal string) (float64, error) {
	if oldVal == "" || newVal == "" {
		return 0, errors.New("on of the args is nil")
	}

	oldValF, err := strconv.ParseFloat(oldVal, 64)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	newValF, err := strconv.ParseFloat(newVal, 64)
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

func ChangeDownTest(oldVal, newVal string) (float64, error) {
	if oldVal == "" || newVal == "" {
		return 0, errors.New("on of the args is nil")
	}

	oldValF, err := strconv.ParseFloat(oldVal, 64)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	newValF, err := strconv.ParseFloat(newVal, 64)
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

// DateCheckTest - returns (val - current data) in days
func DateCheckTest(val time.Time) (float64, error) {
	return val.Sub(time.Now()).Hours(), nil
}

// EqualCheckTest - returns abs(currValF - equalValF)
func EqualCheckTest(currVal, equalVal string) (float64, error) {
	if currVal == "" || equalVal == "" {
		return 0, errors.New("on of the args is nil")
	}

	currValF, err := strconv.ParseFloat(currVal, 64)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	equalValF, err := strconv.ParseFloat(equalVal, 64)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	return math.Abs(currValF - equalValF), nil
}

/*func ParseToFloat64(val1, val2, val3 interface{}) (val1F float64,
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
}*/
