package msgsender

import (
	"github.com/adarocket/alerter/internal/client"
	"github.com/adarocket/proto/proto-gen/notifier"
	"log"
)

type MsgSender struct {
	notifyClient *client.NotifierClient
	stack        map[KeyMsgSender]ValueMsgSender
}

type KeyMsgSender struct {
	NodeUuid  string
	NodeField string
}

type ValueMsgSender struct {
	notify   *notifier.SendNotifier
	tickSend int64
}

func CreateMsgSender(notifyClient *client.NotifierClient) MsgSender {
	newMap := make(map[KeyMsgSender]ValueMsgSender)
	return MsgSender{
		notifyClient: notifyClient,
		stack:        newMap,
	}
}

func (s *MsgSender) updateNotifierInStack(notifier *notifier.SendNotifier, keyMsgSender KeyMsgSender) error {
	if _, exist := s.stack[keyMsgSender]; exist && notifier.Frequency == Max.String() {
		delete(s.stack, keyMsgSender)
		if err := s.notifyClient.SendMessage(notifier); err != nil {
			log.Println(err)
			return err
		}
	} else if notifier.Frequency == Max.String() {
		s.notifyClient.SendMessage(notifier)
	} else if exist {
		val := s.stack[keyMsgSender]
		if (val.tickSend - 1) <= 0 {
			delete(s.stack, keyMsgSender)
			if err := s.notifyClient.SendMessage(notifier); err != nil {
				log.Println(err)
				return err
			}
			val.tickSend--
			val.notify = notifier
			s.stack[keyMsgSender] = val
		}
	} else {
		valueMsgSender := ValueMsgSender{}
		if err := s.notifyClient.SendMessage(notifier); err != nil {
			log.Println(err)
			return err
		}

		tickDur, err := GetTickFrequency(notifier.GetFrequency())
		if err != nil {
			log.Println(err)
			return err
		}

		valueMsgSender.tickSend = tickDur
		valueMsgSender.notify = notifier
		s.stack[keyMsgSender] = valueMsgSender
	}

	return nil
}

func (s *MsgSender) AddNotifiersToStack(messages map[KeyMsgSender]*notifier.SendNotifier) {
	for key, _ := range s.stack {
		if _, exist := messages[key]; !exist {
			s.notifyClient.SendMessage(&notifier.SendNotifier{
				TypeMessage: "Node " + key.NodeUuid + " field " + key.NodeField,
				Value:       "Field now stable",
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
