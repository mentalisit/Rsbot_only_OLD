package DiscordClient

import (
	"Rsbot_only/internal/models"
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strings"
	"time"
)

const (
	emOK      = "✅"
	emCancel  = "❎"
	emRsStart = "🚀"
	emPl30    = "⌛"
	emPlus    = "➕"
	emMinus   = "➖"
)

func (d *Discord) readReactionQueue(r *discordgo.MessageReactionAdd, message *discordgo.Message) {
	user, err := d.s.User(r.UserID)
	if err != nil {
		d.log.ErrorErr(err)
	}
	if user.ID != message.Author.ID {
		ok, config := d.CheckChannelConfigDS(r.ChannelID)
		if ok {
			in := models.InMessage{
				Tip:         "ds",
				Name:        user.Username,
				NameMention: user.Mention(),
				Ds: struct {
					Mesid   string
					Nameid  string
					Guildid string
					Avatar  string
				}{
					Mesid:   r.MessageID,
					Nameid:  user.ID,
					Guildid: config.Guildid,
					Avatar:  user.AvatarURL("128"),
				},

				Config: config,
				Option: models.Option{
					Reaction: true},
			}
			d.reactionUserRemove(r)

			if r.Emoji.Name == emPlus {
				in.Mtext = "+"
			} else if r.Emoji.Name == emMinus {
				in.Mtext = "-"
			} else if r.Emoji.Name == emOK || r.Emoji.Name == emCancel || r.Emoji.Name == emRsStart || r.Emoji.Name == emPl30 {
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				in.Lvlkz, err = d.storage.DbFunc.ReadMesIdDS(ctx, r.MessageID)
				if err == nil && in.Lvlkz != "" {
					if r.Emoji.Name == emOK {
						in.Timekz = "30"
						in.Mtext = in.Lvlkz + "+"
					} else if r.Emoji.Name == emCancel {
						in.Mtext = in.Lvlkz + "-"
					} else if r.Emoji.Name == emRsStart {
						in.Mtext = in.Lvlkz + "++"
					} else if r.Emoji.Name == emPl30 {
						in.Mtext = in.Lvlkz + "+++"
					}
				}
			}
			d.ChanRsMessage <- in
		}
	}
}

func (d *Discord) reactionUserRemove(r *discordgo.MessageReactionAdd) {
	err := d.s.MessageReactionRemove(r.ChannelID, r.MessageID, r.Emoji.Name, r.UserID)
	if err != nil {
		fmt.Println("Ошибка удаления эмоджи", err)
	}
}

func (d *Discord) logicMix(m *discordgo.MessageCreate) {
	if d.ifMentionBot(m) {
		return
	}
	if d.avatar(m) {
		return
	}
	d.AccesChatDS(m)

	//filter Rs
	ok, config := d.CheckChannelConfigDS(m.ChannelID)
	if ok {
		d.SendToRsFilter(m, config)
		return
	}
}

func (d *Discord) SendToRsFilter(m *discordgo.MessageCreate, config models.CorporationConfig) {
	if len(m.Attachments) > 0 {
		m.Content += m.Attachments[0].URL
	}
	if len(m.Message.Embeds) > 0 {
		m.Content += "\u200B"
	}
	in := models.InMessage{
		Mtext:       m.Content,
		Tip:         "ds",
		Name:        m.Author.Username,
		NameMention: m.Author.Mention(),
		Ds: struct {
			Mesid   string
			Nameid  string
			Guildid string
			Avatar  string
		}{
			Mesid:   m.ID,
			Nameid:  m.Author.ID,
			Guildid: m.GuildID,
			Avatar:  m.Author.AvatarURL("128"),
		},
		Config: config,
		Option: models.Option{InClient: true},
	}
	d.ChanRsMessage <- in

}
func (d *Discord) ifMentionBot(m *discordgo.MessageCreate) bool {
	after, found := strings.CutPrefix(m.Content, d.s.State.User.Mention())
	if found {
		if len(after) > 0 {
			split := strings.Split(after, " ")
			if split[0] == "help" || split[0] == "справка" || split[0] == "довідка" {
				//nujno sdelat obshuu spravku
				d.SendChannelDelSecond(m.ChannelID, "сорян в разработке", 10)
				return true
			}
		}

		d.DeleteMesageSecond(m.ChannelID, m.ID, 30)
		goodRs, _ := d.CheckChannelConfigDS(m.ChannelID)
		var text string
		if goodRs {
			text = fmt.Sprintf("%s че пингуешь? пиши Справка,или пиши создателю бота @Mentalisit#5159 ", m.Author.Mention())
		} else {
			text = fmt.Sprintf("%s че пингуешь? я же многофункциональный бот, Префикс доступен только после активации нужного режима \n Для получения справки пиши %s help",
				m.Author.Mention(), d.s.State.User.Mention())
		}
		d.SendChannelDelSecond(m.ChannelID, text, 30)
	}
	return found
}
