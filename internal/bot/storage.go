package bot

import "Rsbot_only/internal/models"

func (b *Bot) CheckCorpNameConfig(corpName string) (bool, models.CorporationConfig) {
	for _, config := range b.configCorp {
		if config.CorpName == corpName {
			return true, config
		}
	}
	return false, models.CorporationConfig{}
}
