package TelegramClient

import (
	"Rsbot_only/internal/models"
	tgbotapi "github.com/matterbridge/telegram-bot-api/v6"
	"strconv"
	"strings"
	"time"
)

func (t *Telegram) SendEmded(lvlkz string, chatid string, text string) int {
	a := strings.SplitN(chatid, "/", 2)
	chatId, err := strconv.ParseInt(a[0], 10, 64)
	if err != nil {
		t.log.ErrorErr(err)
	}
	ThreadID := 0
	if len(a) > 1 {
		ThreadID, err = strconv.Atoi(a[1])
		if err != nil {
			t.log.ErrorErr(err)
		}
	}
	var keyboardQueue = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(lvlkz+"+", lvlkz+"+"),
			tgbotapi.NewInlineKeyboardButtonData(lvlkz+"-", lvlkz+"-"),
			tgbotapi.NewInlineKeyboardButtonData(lvlkz+"++", lvlkz+"++"),
			tgbotapi.NewInlineKeyboardButtonData(lvlkz+"+30", lvlkz+"+++"),
		),
	)
	msg := tgbotapi.NewMessage(chatId, text)
	msg.MessageThreadID = ThreadID
	msg.ReplyMarkup = keyboardQueue
	message, _ := t.t.Send(msg)

	return message.MessageID

}
func (t *Telegram) SendEmbedTime(chatid string, text string) int {
	a := strings.SplitN(chatid, "/", 2)
	chatId, err := strconv.ParseInt(a[0], 10, 64)
	if err != nil {
		t.log.ErrorErr(err)
	}
	ThreadID := 0
	if len(a) > 1 {
		ThreadID, err = strconv.Atoi(a[1])
		if err != nil {
			t.log.ErrorErr(err)
		}
	}

	var keyboardQueue = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("+", "+"),
			tgbotapi.NewInlineKeyboardButtonData("-", "-"),
		),
	)
	msg := tgbotapi.NewMessage(chatId, text)
	msg.MessageThreadID = ThreadID
	msg.ReplyMarkup = keyboardQueue
	message, _ := t.t.Send(msg)

	return message.MessageID
}

// отправка сообщения в телегу
func (t *Telegram) SendChannel(chatid string, text string) int {
	a := strings.SplitN(chatid, "/", 2)
	chatId, err := strconv.ParseInt(a[0], 10, 64)
	if err != nil {
		t.log.ErrorErr(err)
	}
	ThreadID := 0
	if len(a) > 1 {
		ThreadID, err = strconv.Atoi(a[1])
		if err != nil {
			t.log.ErrorErr(err)
		}
	}
	m := tgbotapi.NewMessage(chatId, text)
	m.MessageThreadID = ThreadID
	tMessage, _ := t.t.Send(m)
	return tMessage.MessageID
}

func (t *Telegram) SendChannelDelSecond(chatid string, text string, second int) {
	a := strings.SplitN(chatid, "/", 2)
	chatId, err := strconv.ParseInt(a[0], 10, 64)
	if err != nil {
		t.log.ErrorErr(err)
	}
	ThreadID := 0
	if len(a) > 1 {
		ThreadID, err = strconv.Atoi(a[1])
		if err != nil {
			t.log.ErrorErr(err)
		}
	}
	m := tgbotapi.NewMessage(chatId, text)
	m.MessageThreadID = ThreadID
	tMessage, err1 := t.t.Send(m)
	if err1 != nil {
		t.log.Error(err1.Error())
	}
	if second <= 60 {
		go func() {
			time.Sleep(time.Duration(second) * time.Second)
			t.DelMessage(chatid, tMessage.MessageID)
		}()
	} else {
		t.storage.TimeDeleteMessage.TimerInsert(models.Timer{
			Tgmesid:  strconv.Itoa(tMessage.MessageID),
			Tgchatid: chatid,
			Timed:    second,
		})
	}
}

func (t *Telegram) ChatTyping(chatId string) {
	chatid, threadID := t.chat(chatId)
	typingConfig := tgbotapi.NewChatAction(chatid, tgbotapi.ChatTyping)
	typingConfig.MessageThreadID = threadID
	_, _ = t.t.Send(typingConfig)
}
