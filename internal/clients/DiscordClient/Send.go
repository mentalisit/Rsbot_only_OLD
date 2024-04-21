package DiscordClient

import (
	"Rsbot_only/internal/clients/DiscordClient/transmitter"
	"Rsbot_only/internal/models"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"time"
)

var mesContentNil string

func (d *Discord) SendEmbedText(chatid, title, text string) *discordgo.Message {
	Emb := &discordgo.MessageEmbed{
		Author:      &discordgo.MessageEmbedAuthor{},
		Color:       16711680,
		Description: text,
		Title:       title,
	}
	m, err := d.s.ChannelMessageSendEmbed(chatid, Emb)
	if err != nil {
		d.log.Error("chatid " + chatid + " " + err.Error())
		return &discordgo.Message{}
	}
	return m
}
func (d *Discord) SendChannelDelSecond(chatid, text string, second int) {
	if text != "" {
		message, err := d.s.ChannelMessageSend(chatid, text)
		if err != nil {
			d.log.ErrorErr(err)
			return
		}
		if second <= 60 {
			go func() {
				time.Sleep(time.Duration(second) * time.Second)
				_ = d.s.ChannelMessageDelete(chatid, message.ID)
			}()
		} else {
			d.storage.TimeDeleteMessage.TimerInsert(models.Timer{
				Dsmesid:  message.ID,
				Dschatid: chatid,
				Timed:    second,
			})
		}
	}
}
func (d *Discord) SendComplexContent(chatid, text string) (mesId string) { //отправка текста комплексного сообщения
	mesCompl, err := d.s.ChannelMessageSendComplex(chatid, &discordgo.MessageSend{
		Content: text})
	if err != nil {
		channel, _ := d.s.Channel(chatid)
		d.log.Info("Ошибка отправки комплексного сообщения text " + channel.Name)
		d.log.ErrorErr(err)
		mesCompl, err = d.s.ChannelMessageSendComplex(chatid, &discordgo.MessageSend{
			Content: text})
		if err == nil {
			return mesCompl.ID
		}
		return ""
	}
	return mesCompl.ID
}
func (d *Discord) SendComplex(chatid string, embeds discordgo.MessageEmbed, component []discordgo.MessageComponent) (mesId string) { //отправка текста комплексного сообщения
	mesCompl, err := d.s.ChannelMessageSendComplex(chatid, &discordgo.MessageSend{
		Content:    mesContentNil,
		Embed:      &embeds,
		Components: component,
	})
	if err != nil {
		channel, _ := d.s.Channel(chatid)
		d.log.Info("Ошибка отправки комплексного сообщения embed " + channel.Name)
		d.log.ErrorErr(err)
		mesCompl, err = d.s.ChannelMessageSendComplex(chatid, &discordgo.MessageSend{
			Content: mesContentNil,
			Embed:   &embeds,
		})
		if err == nil {
			return mesCompl.ID
		}
		return ""
	}
	return mesCompl.ID
}
func (d *Discord) Send(chatid, text string) (mesId string) { //отправка текста
	message, err := d.s.ChannelMessageSend(chatid, text)
	if err != nil {
		d.log.ErrorErr(err)
	}
	return message.ID
}

func (d *Discord) SendEmbedTime(chatid, text string) (mesId string) { //отправка текста с двумя реакциями
	message, err := d.s.ChannelMessageSendComplex(chatid, &discordgo.MessageSend{
		Content: text,
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Style:    discordgo.SecondaryButton,
						Label:    "+",
						CustomID: "+",
						Emoji: &discordgo.ComponentEmoji{
							Name: emPlus},
					},

					&discordgo.Button{
						Style:    discordgo.SecondaryButton,
						Label:    "-",
						CustomID: "-",
						Emoji: &discordgo.ComponentEmoji{
							Name: emMinus},
					}}},
		},
	})
	if err != nil {
		d.log.ErrorErr(err)
		return ""
	}
	return message.ID
}

func (d *Discord) SendWebhook(text, username, chatid, guildId, Avatar string) (mesId string) {
	if text == "" {
		return ""
	}
	web := transmitter.New(d.s, guildId, "KzBot", true, d.log)
	pp := discordgo.WebhookParams{
		Content:   text,
		Username:  username,
		AvatarURL: Avatar,
	}
	mes, err := web.Send(chatid, &pp)
	if err != nil {
		fmt.Println(err)
		m := d.Send(chatid, text)
		return m
	}
	return mes.ID
}
