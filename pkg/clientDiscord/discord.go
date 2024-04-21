package clientDiscord

import (
	"Rsbot_only/internal/config"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/mentalisit/logger"
)

func NewDiscord(log *logger.Logger, cfg *config.ConfigBot) (*discordgo.Session, error) {
	DSBot, err := discordgo.New("Bot " + cfg.Token.TokenDiscord)
	if err != nil {
		log.Panic("Ошибка запуска дискорда" + err.Error())
		return nil, err
	}

	err = DSBot.Open()
	if err != nil {
		log.Panic("Ошибка открытия ДС " + err.Error())
	}
	fmt.Println("Бот Дискорд загружен ")
	return DSBot, nil
}
