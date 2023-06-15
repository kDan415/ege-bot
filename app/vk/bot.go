package vk

import (
	"context"
	"ege/app/ege"
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/SevereCloud/vksdk/v2/longpoll-bot"
	"log"
)

type Bot struct {
	lastResult  string
	chatID      int
	vk          *api.VK
	lp          *longpoll.LongPoll
	grabber     *ege.Grabber
	notificator *ege.Notificator
}

func NewBot(conf Config, grabber *ege.Grabber) (*Bot, error) {

	bot := &Bot{
		chatID:  conf.ChatID,
		grabber: grabber,
	}

	bot.vk = api.NewVK(conf.Token)

	groups, err := bot.vk.GroupsGetByID(nil)
	if err != nil {
		return nil, NewError(APIVKError, err)
	}

	bot.lp, err = longpoll.NewLongPoll(bot.vk, groups[0].ID)
	if err != nil {
		return nil, NewError(APIVKError, err)
	}

	bot.lp.MessageNew(bot.getLongPollHandler())

	bot.notificator = ege.NewNotificator(conf.Delay, grabber, func(s string) {
		bot.SendMsg(bot.chatID, s)
	})
	err = bot.notificator.Start()
	if err != nil {
		return nil, err
	}

	return bot, nil
}

func (bot *Bot) getLongPollHandler() func(_ context.Context, obj events.MessageNewObject) {
	return func(_ context.Context, obj events.MessageNewObject) {

		if obj.Message.PeerID != bot.chatID {
			bot.SendMsg(obj.Message.PeerID, "access is denied")
			return
		}

		switch obj.Message.Text {
		case "ping":
			bot.SendMsg(bot.chatID, "pong")
		case "stats":
			bot.SendMsg(bot.chatID, bot.notificator.GetStats().String())
		case "last":
			bot.SendMsg(bot.chatID, "Last notificator result:\n"+bot.grabber.GetLastResult().String())
		case "result":
			result, err := bot.grabber.GetResult()
			if err != nil {
				log.Print("Result for LongPoll: ", err.Error())
				bot.SendMsg(bot.chatID, err.Error())
				bot.SendMsg(bot.chatID, "Last result:\n"+bot.grabber.GetLastResult().String())
				return
			}
			bot.SendMsg(bot.chatID, result.String())
		case "start":
			if err := bot.notificator.Start(); err != nil {
				bot.SendMsg(bot.chatID, err.Error())
				return
			}
		case "stop":
			if err := bot.notificator.Stop(); err != nil {
				bot.SendMsg(bot.chatID, err.Error())
				return
			}
		}
		return
	}
}

func (bot *Bot) RunLongPool() error {
	if err := bot.lp.Run(); err != nil {
		return err
	}
	return nil
}

func (bot *Bot) StartNotificator() error {
	return bot.notificator.Start()
}

func (bot *Bot) StopNotificator() error {
	return bot.notificator.Stop()
}

func (bot *Bot) SendMsg(chatID int, msg string) {
	var err error
	b := params.NewMessagesSendBuilder()
	b.Message(msg)
	b.RandomID(0)
	b.PeerID(chatID)

	_, err = bot.vk.MessagesSend(b.Params)
	if err != nil {
		log.Print(NewError(APIVKError, err))
	}
}
