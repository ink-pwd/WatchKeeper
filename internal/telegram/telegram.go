package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ink-pwd/WatchKeeper/internal/scheduler"
	"github.com/ink-pwd/WatchKeeper/internal/utils"
)

type TelegramBot struct {
	bot     *tgbotapi.BotAPI
	sched   *scheduler.Scheduler
	timeOut int
}

func CreateTelegramBot(bot *tgbotapi.BotAPI, sched *scheduler.Scheduler, timeOut int) *TelegramBot {
	return &TelegramBot{
		bot:     bot,
		sched:   sched,
		timeOut: timeOut,
	}
}

func (t *TelegramBot) ListenServ() {
	var (
		u       tgbotapi.UpdateConfig
		updates tgbotapi.UpdatesChannel
		upd     tgbotapi.Update
		chatID  int64
		link    string
		err     error
	)
	u = tgbotapi.NewUpdate(0)
	u.Timeout = t.timeOut

	updates = t.bot.GetUpdatesChan(u)

	for upd = range updates {
		if upd.Message != nil {
			chatID = upd.Message.Chat.ID
			switch upd.Message.Command() {
			case "add":
				link, err = utils.GetHostName(upd.Message.CommandArguments())
				if err != nil {
					t.SendMessage(chatID, "Enter the correct link type {https://google.com}")
					continue
				}
				/*
					Если запрос введен корректно - добавляем его в планировщик задач.
				*/
				err = t.sched.AddTask(chatID, link)
				if err != nil {
					t.SendMessage(chatID, "Error adding link, try again later")
					continue
				}
				t.SendMessage(chatID, "Success add link")
			default:
				t.SendMessage(chatID, "Use /add {link} for add website to the monitoring list")
			}
		}
	}
}

func (t *TelegramBot) SendMessage(chatID int64, text string) {
	var (
		msg tgbotapi.MessageConfig
	)
	msg = tgbotapi.NewMessage(chatID, text)
	t.bot.Send(msg)
}
