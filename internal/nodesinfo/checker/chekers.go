package checker

import (
	"errors"
	"math"
	"strconv"
	"time"
)

// Interval returns -val if val < a1 and +val if val > a2, return 0 if ok
func Interval(bottomLine, topLine, val float64) float64 {
	if val > topLine {
		return math.Abs(val) - math.Abs(topLine)
	}

	if val < bottomLine {
		return math.Abs(bottomLine) - math.Abs(val)
	}

	return val
}

func ChangeUp(oldVal, newVal float64) float64 {
	return newVal - oldVal
}

func ChangeDown(oldVal, newVal float64) float64 {
	return oldVal - newVal
}

func DateCheck(dat time.Time) float64 {
	return float64(dat.Unix() - time.Now().Unix())
}

func Checker(a1, a2, a3 interface{}, checkerType string) (float64, error) {
	var a1float64 float64
	var a2float64 float64
	var a3float64 float64

	if a1 != nil {
		if a1Str, isTrt := a1.(string); isTrt {
			a1float64, _ = strconv.ParseFloat(a1Str, 64)
		}
		if a1Float, isTrt := a1.(float64); isTrt {
			a1float64 = a1Float
		}
	}
	if a2 != nil {
		if a1Str, isTrt := a2.(string); isTrt {
			a2float64, _ = strconv.ParseFloat(a1Str, 64)
		}
		if a1Float, isTrt := a2.(float64); isTrt {
			a2float64 = a1Float
		}
	}
	if a3 != nil {
		if a1Str, isTrt := a3.(string); isTrt {
			a3float64, _ = strconv.ParseFloat(a1Str, 64)
		}
		if a1Float, isTrt := a3.(float64); isTrt {
			a3float64 = a1Float
		}
	}

	switch checkerType {
	case IntervalT.String():
		return Interval(a1float64, a2float64, a3float64), nil
	case ChangeUpT.String():
		return ChangeUp(a1float64, a2float64), nil
	case ChangeDownT.String():
		return ChangeDown(a1float64, a2float64), nil
	case DateT.String():
		var dat time.Time
		if a, isTr := a1.(time.Time); isTr {
			dat = a
		} else {
			return 0, errors.New("invalid date type")
		}

		return DateCheck(dat), nil
	default:
		return 0, errors.New("checker type not support")
	}
}
