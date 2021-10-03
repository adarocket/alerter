package msgsender

import (
	"errors"
	"time"
)

type TypeFrequency int

const (
	Low TypeFrequency = 1 + iota
	Normal
	Max
)

var frequencyTypes = [...]string{
	"low",
	"normal",
	"max",
}

var frequencyTypesTime = [...]time.Duration{
	time.Minute * 10,
	time.Minute * 5,
	time.Minute,
}

func (tf TypeFrequency) String() string {
	return frequencyTypes[tf-1]
}

func (tf TypeFrequency) GetTime() time.Duration {
	return frequencyTypesTime[tf-1]
}

func GetTimeFrequency(frequency string) (time.Duration, error) {
	for i, frequencyType := range frequencyTypes {
		if frequency == frequencyType {
			return frequencyTypesTime[i], nil
		}
	}

	return 0, errors.New("invalid type")
}
