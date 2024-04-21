package TelegramClient

import (
	"Rsbot_only/internal/models"
	"fmt"
	tgbotapi "github.com/matterbridge/telegram-bot-api/v6"
	"strconv"
)

func (t *Telegram) logicMix(m *tgbotapi.Message, edit bool) {
	t.accesChatTg(m) //это была начальная функция при добавлени бота в группу
	ThreadID := m.MessageThreadID
	if !m.IsTopicMessage && ThreadID != 0 {
		ThreadID = 0
	}
	ChatId := strconv.FormatInt(m.Chat.ID, 10) + fmt.Sprintf("/%d", ThreadID)

	// RsClient
	ok, config := t.CheckChannelConfigTG(ChatId)
	if ok {
		name := t.nameNick(m.From.UserName, m.From.FirstName, config.TgChannel)
		in := models.InMessage{
			Mtext:       m.Text,
			Tip:         "tg",
			Name:        name,
			NameMention: "@" + name,
			Tg: struct {
				Mesid int
			}{
				Mesid: m.MessageID,
			},
			Config: config,
			Option: models.Option{
				InClient: true,
			},
		}
		if in.Mtext == "" && config.Forward {
			t.DelMessageSecond(ChatId, strconv.Itoa(m.MessageID), 180)
		}

		t.ChanRsMessage <- in
	}

}
