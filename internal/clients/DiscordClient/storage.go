package DiscordClient

import (
	"Rsbot_only/internal/models"
)

// CheckChannelConfigDS RsConfig
func (d *Discord) CheckChannelConfigDS(chatid string) (channelGood bool, config models.CorporationConfig) {
	for _, corpporationConfig := range d.corpConfigRS {
		if corpporationConfig.DsChannel == chatid {
			return true, corpporationConfig
		}
	}
	return false, models.CorporationConfig{}
}

// AddDsCorpConfig add RsConfig
func (d *Discord) AddDsCorpConfig(chatName, chatid, guildid, lang string) {
	c := models.CorporationConfig{
		CorpName:  chatName,
		DsChannel: chatid,
		Country:   lang,
		Guildid:   guildid,
	}
	c.MesidDsHelp = d.HelpChannelUpdate(c)
	d.storage.ConfigRs.InsertConfigRs(c)
	d.corpConfigRS[c.CorpName] = c
	d.log.Info(chatName + " Добавлена в конфиг корпораций ")
}
