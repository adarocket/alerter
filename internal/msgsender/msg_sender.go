package msgsender

import (
	"github.com/adarocket/alerter/internal/client"
	"github.com/adarocket/proto/proto-gen/notifier"
	"log"
)

const MaxFrequency = 1

type MsgSender struct {
	notifyClient *client.NotifierClient
	stack        map[KeyMsg]int64
}

type KeyMsg struct {
	NodeUuid  string
	NodeField string
}

type BodyMsg struct {
	Notify    *notifier.SendNotifier
	Frequency int64
}

func CreateMsgSender(notifyClient *client.NotifierClient) MsgSender {
	newMap := make(map[KeyMsg]int64)
	return MsgSender{
		notifyClient: notifyClient,
		stack:        newMap,
	}
}

func (s *MsgSender) updateNotifierInStack(notifier BodyMsg, keyMsgSender KeyMsg) error {
	if _, exist := s.stack[keyMsgSender]; exist && notifier.Frequency == MaxFrequency {
		delete(s.stack, keyMsgSender)
		if err := s.notifyClient.SendMessage(notifier.Notify); err != nil {
			log.Println(err)
			return err
		}
	} else if notifier.Frequency == MaxFrequency {
		s.notifyClient.SendMessage(notifier.Notify)
	} else if exist {
		val := s.stack[keyMsgSender]
		if (val - 1) <= 0 {
			if err := s.notifyClient.SendMessage(notifier.Notify); err != nil {
				log.Println(err)
				return err
			}

			s.stack[keyMsgSender] = notifier.Frequency

			return nil
		}
		s.stack[keyMsgSender] = val - 1
	} else {
		if err := s.notifyClient.SendMessage(notifier.Notify); err != nil {
			log.Println(err)
			return err
		}

		s.stack[keyMsgSender] = notifier.Frequency
	}

	return nil
}

// AddNotifiersToStack - add new notifiers to stack or delete old notifier if new map doest have old one
func (s *MsgSender) AddNotifiersToStack(messages map[KeyMsg]BodyMsg) {
	for key, _ := range s.stack {
		if _, exist := messages[key]; !exist {
			s.notifyClient.SendMessage(&notifier.SendNotifier{
				TextMessage: "Node " + key.NodeUuid + " field " + key.NodeField + " Field now stable",
			})
			delete(s.stack, key)
			continue
		}
	}

	for key, sendNotifier := range messages {
		if err := s.updateNotifierInStack(sendNotifier, key); err != nil {
			log.Println(err)
		}
	}
}
