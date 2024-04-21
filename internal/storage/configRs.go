package storage

import (
	"Rsbot_only/internal/models"
)

type ConfigRs interface {
	InsertConfigRs(c models.CorporationConfig)
	ReadConfigRs() []models.CorporationConfig
	DeleteConfigRs(c models.CorporationConfig)
	AutoHelpUpdateMesid(c models.CorporationConfig)
	AutoHelp() []models.CorporationConfig
}

func (s *Storage) DeleteConfigRs(c models.CorporationConfig) {
	s.ConfigRs.DeleteConfigRs(c)
	var a map[string]models.CorporationConfig
	a = make(map[string]models.CorporationConfig)
	b := s.ConfigRs.ReadConfigRs()
	for _, config := range b {
		a[config.CorpName] = config
	}
	s.CorpConfigRS = a
}
