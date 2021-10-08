package msgsender

import (
	"errors"
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

var tickSend = [...]int64{
	20,
	10,
	1,
}

func (tf TypeFrequency) String() string {
	return frequencyTypes[tf-1]
}

func (tf TypeFrequency) GetTick() int64 {
	return tickSend[tf-1]
}

func GetTickFrequency(frequency string) (int64, error) {
	for i, frequencyType := range frequencyTypes {
		if frequency == frequencyType {
			return tickSend[i], nil
		}
	}

	return 0, errors.New("invalid type")
}
