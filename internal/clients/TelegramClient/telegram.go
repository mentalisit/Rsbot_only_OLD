package TelegramClient

import (
	"Rsbot_only/internal/config"
	"Rsbot_only/internal/models"
	"Rsbot_only/internal/storage"
	"Rsbot_only/pkg/clientTelegram"
	tgbotapi "github.com/matterbridge/telegram-bot-api/v6"
	"github.com/mentalisit/logger"
)

type Telegram struct {
	ChanRsMessage chan models.InMessage
	t             *tgbotapi.BotAPI
	log           *logger.Logger
	storage       *storage.Storage
	debug         bool
	corpConfigRS  map[string]models.CorporationConfig
}

func NewTelegram(log *logger.Logger, st *storage.Storage, cfg *config.ConfigBot) *Telegram {
	client, err := clientTelegram.NewTelegram(log, cfg)
	if err != nil {
		return nil
	}

	tg := &Telegram{
		ChanRsMessage: make(chan models.InMessage, 10),
		t:             client,
		log:           log,
		storage:       st,
		debug:         cfg.IsDebug,
		corpConfigRS:  st.CorpConfigRS,
	}
	go tg.update()

	return tg
}
func (t *Telegram) update() {
	ut := tgbotapi.NewUpdate(0)
	ut.Timeout = 60
	//получаем обновления от телеграм
	updates := t.t.GetUpdatesChan(ut)
	for update := range updates {
		if update.InlineQuery != nil {
			//t.handleInlineQuery(update.InlineQuery)
		} else if update.ChosenInlineResult != nil {
			//go t.handleChosenInlineResult(update.ChosenInlineResult)
		} else if update.CallbackQuery != nil {
			t.callback(update.CallbackQuery) //нажатия в чате
		} else if update.Message != nil {

			if update.Message.Chat.IsPrivate() { //если пишут боту в личку
				//t.ifPrivatMesage(update.Message)
			} else if update.Message.IsCommand() {
				t.updatesComand(update.Message) //если сообщение является командой
			} else { //остальные сообщения
				t.logicMix(update.Message, false)
			}
		} else if update.EditedMessage != nil {
			t.logicMix(update.EditedMessage, true)
		} else if update.MyChatMember != nil {
			t.myChatMember(update.MyChatMember)

		} else if update.ChatMember != nil {
			t.chatMember(update.ChatMember)

		} else {

		}
	}
}
