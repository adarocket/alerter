package telegram

import (
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"log"
)

type TelegramMessenger struct {
}

const telegramToken = "1808950021:AAENv6bK4qf9KzpUXGqiK9NRH4uy_VtQX5k"

func (t *TelegramMessenger) StartBot(msgStatus chan string) {
	bot, err := tgbotapi.NewBotAPI(telegramToken)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)
	quit := make(chan int)
	funcIsRunning := false

	for update := range updates {
		if update.Message == nil {
			continue
		}
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		var reply string
		sendStatus := func() {
			for {
				funcIsRunning = true
				select {
				case m := <-msgStatus:
					reply := m
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
					if _, err := bot.Send(msg); err != nil {
						log.Println(err)
						break
					}
				case <-quit:
					funcIsRunning = false
					return
				}
			}
		}

		switch update.Message.Command() {
		case "stop":
			reply = "уже остановлен"
			if funcIsRunning {
				reply = "остановлен"
				quit <- 0
			}
		case "start":
			reply = "уже запущен"
			if !funcIsRunning {
				reply = "Привет. Подключаюсь к alerter"
				go sendStatus()
			}
		default:
			reply = "неизвестная команда"
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		bot.Send(msg)
	}
}
