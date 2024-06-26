package TelegramClient

import "Rsbot_only/internal/models"

func (t *Telegram) CheckChannelConfigTG(chatid string) (channelGood bool, config models.CorporationConfig) {
	for _, corpporationConfig := range t.corpConfigRS {
		if corpporationConfig.TgChannel == chatid {
			return true, corpporationConfig
		}
	}
	return false, models.CorporationConfig{}
}

// AddTgCorpConfig add RsConfig
func (t *Telegram) AddTgCorpConfig(chatName string, chatid, lang string) {
	c := models.CorporationConfig{
		CorpName:  chatName,
		Country:   lang,
		TgChannel: chatid,
	}
	t.storage.ConfigRs.InsertConfigRs(c)
	t.corpConfigRS[c.CorpName] = c
	t.log.Info(chatName + " Добавлена в конфиг корпораций ")
}
