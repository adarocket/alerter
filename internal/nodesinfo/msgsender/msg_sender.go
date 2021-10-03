package msgsender

import (
	"github.com/adarocket/alerter/internal/client"
	"github.com/adarocket/proto/proto-gen/notifier"
	"log"
	"time"
)

type MsgSender struct {
	notifyClient *client.NotifierClient
	stack        map[KeyMsgSender]*notifier.SendNotifier
}

type KeyMsgSender struct {
	NodeUuid  string
	NodeField string
}

func CreateMsgSender(notifyClient *client.NotifierClient) MsgSender {
	newMap := make(map[KeyMsgSender]*notifier.SendNotifier)
	return MsgSender{
		notifyClient: notifyClient,
		stack:        newMap,
	}
}

func (s *MsgSender) AddNotifierToStack(notifier *notifier.SendNotifier, keyMsgSender KeyMsgSender) {
	if _, exist := s.stack[keyMsgSender]; exist && notifier.Frequency == Max.String() {
		delete(s.stack, keyMsgSender)
		s.notifyClient.SendMessage(notifier)
	} else if notifier.Frequency == Max.String() {
		s.notifyClient.SendMessage(notifier)
	} else if exist {
		s.stack[keyMsgSender] = notifier
	} else {
		go func() {
			s.stack[keyMsgSender] = notifier

			tmDur, err := GetTimeFrequency(notifier.GetFrequency())
			if err != nil {
				log.Println(err)
				return
			}

			time.Sleep(tmDur)
			if val, existAfter := s.stack[keyMsgSender]; existAfter {
				delete(s.stack, keyMsgSender)
				s.notifyClient.SendMessage(val)
			}
		}()
	}
}
