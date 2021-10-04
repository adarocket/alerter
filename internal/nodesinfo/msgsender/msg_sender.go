package msgsender

import (
	"fmt"
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

func (s *MsgSender) updateNotifierInStack(notifier *notifier.SendNotifier, keyMsgSender KeyMsgSender) {
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
			s.notifyClient.SendMessage(notifier)

			tmDur, err := GetTimeFrequency(notifier.GetFrequency())
			if err != nil {
				log.Println(err)
				return
			}

			for {
				time.Sleep(tmDur)
				if val, existAfter := s.stack[keyMsgSender]; existAfter {
					fmt.Println("Send msg")
					s.notifyClient.SendMessage(val)
				} else {
					fmt.Println("delete this")
					return
				}
			}
		}()
	}
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

	for sender, sendNotifier := range messages {
		s.updateNotifierInStack(sendNotifier, sender)
	}
}
