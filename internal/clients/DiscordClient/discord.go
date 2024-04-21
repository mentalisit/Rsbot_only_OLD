package DiscordClient

import (
	"Rsbot_only/internal/config"
	"Rsbot_only/internal/models"
	"Rsbot_only/internal/storage"
	"Rsbot_only/pkg/clientDiscord"
	"github.com/bwmarrin/discordgo"
	"github.com/mentalisit/logger"
)

type Discord struct {
	ChanRsMessage chan models.InMessage
	s             *discordgo.Session
	log           *logger.Logger
	storage       *storage.Storage
	corpConfigRS  map[string]models.CorporationConfig
}

func NewDiscord(log *logger.Logger, st *storage.Storage, cfg *config.ConfigBot) *Discord {
	ds, err := clientDiscord.NewDiscord(log, cfg)
	if err != nil {
		log.ErrorErr(err)
	}

	DS := &Discord{
		s:             ds,
		log:           log,
		storage:       st,
		ChanRsMessage: make(chan models.InMessage, 10),
		corpConfigRS:  st.CorpConfigRS,
	}
	go ds.AddHandler(DS.messageHandler)
	go ds.AddHandler(DS.messageReactionAdd)
	go ds.AddHandler(DS.slash)
	go DS.ready()

	return DS
}
