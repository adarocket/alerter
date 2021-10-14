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
	Notify    *notifier.SendNotifier
	Frequency string
	tickSend  int64
}

func CreateMsgSender(notifyClient *client.NotifierClient) MsgSender {
	newMap := make(map[KeyMsgSender]ValueMsgSender)
	return MsgSender{
		notifyClient: notifyClient,
		stack:        newMap,
	}
}

func (s *MsgSender) updateNotifierInStack(notifier ValueMsgSender, keyMsgSender KeyMsgSender) error {
	if _, exist := s.stack[keyMsgSender]; exist && notifier.Frequency == Max.String() {
		delete(s.stack, keyMsgSender)
		if err := s.notifyClient.SendMessage(notifier.Notify); err != nil {
			log.Println(err)
			return err
		}
	} else if notifier.Frequency == Max.String() {
		s.notifyClient.SendMessage(notifier.Notify)
	} else if exist {
		val := s.stack[keyMsgSender]
		if (val.tickSend - 1) <= 0 {
			delete(s.stack, keyMsgSender)
			if err := s.notifyClient.SendMessage(notifier.Notify); err != nil {
				log.Println(err)
				return err
			}
			val.tickSend--
			val.Notify = notifier.Notify
			s.stack[keyMsgSender] = val
		}
	} else {
		valueMsgSender := ValueMsgSender{}
		if err := s.notifyClient.SendMessage(notifier.Notify); err != nil {
			log.Println(err)
			return err
		}

		tickDur, err := GetTickFrequency(notifier.Frequency)
		if err != nil {
			log.Println(err)
			return err
		}

		valueMsgSender.tickSend = tickDur
		valueMsgSender.Notify = notifier.Notify
		s.stack[keyMsgSender] = valueMsgSender
	}

	return nil
}

func (s *MsgSender) AddNotifiersToStack(messages map[KeyMsgSender]ValueMsgSender) {
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
