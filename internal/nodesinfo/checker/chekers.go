package checker

import (
	"encoding/json"
	"errors"
	"github.com/hashicorp/go-version"
	"io"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"
)

// rebuild it
func IntervalCalculate(bottomLine, topLine float64, val string) (float64, error) {
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

// ChangeUpCalculate - returns 0 if newVal - oldVal < 0, in other cases diff
func ChangeUpCalculate(oldVal, newVal string) (float64, error) {
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

func ChangeDownCalculate(oldVal, newVal string) (float64, error) {
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

// DateDiffCalculate - returns (val - current data) in days
func DateDiffCalculate(val time.Time) (float64, error) {
	return val.Sub(time.Now()).Hours(), nil
}

// EqualCheck - returns abs(currValF - equalValF)
func EqualCheck(currVal, equalVal string) (float64, error) {
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

func IsCardanoVersionOutdated(currVersion string) (bool, error) {
	res, err := http.Get("https://api.github.com/repos/adarocket/informer/releases/latest")
	if err != nil {
		log.Println(err)
		return false, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return false, err
	}
	res.Body.Close()

	if res.StatusCode > 299 {
		log.Printf("Response failed with status code: %d and\nbody: %s\n",
			res.StatusCode, body)
		return false, err
	}

	var info infoAboutLibGitHub
	err = json.Unmarshal(body, &info)
	if err != nil {
		log.Println(err)
		return false, err
	}

	v1, err := version.NewVersion(currVersion)
	v2, err := version.NewVersion(info.TagName)
	if v1.LessThan(v2) {
		return true, nil
	}

	return false, nil
}

type infoAboutLibGitHub struct {
	TagName string `json:"tag_name"`
}

func GetIntBool(val string) (float64, error) {
	currVal, err := strconv.ParseBool(val)
	if err != nil {
		log.Println(err)
		return -1, err
	}

	if currVal {
		return 1, nil
	}

	return 0, nil
}
