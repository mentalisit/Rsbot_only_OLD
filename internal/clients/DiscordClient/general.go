package DiscordClient

import (
	"Rsbot_only/internal/clients/DiscordClient/transmitter"
	"Rsbot_only/internal/models"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"time"
)

//lang ok

func (d *Discord) AddButtonsQueue(level string) []discordgo.MessageComponent {
	// Создание кнопки
	buttonOk := discordgo.Button{
		Style:    discordgo.SecondaryButton,
		Label:    level + "+",
		CustomID: level + "+",
		Emoji: &discordgo.ComponentEmoji{
			Name: emOK,
		},
	}
	buttonCancel := discordgo.Button{
		Style:    discordgo.SecondaryButton,
		Label:    level + "-",
		CustomID: level + "-",
		Emoji: &discordgo.ComponentEmoji{
			Name: emCancel,
		},
	}
	buttonRsStart := discordgo.Button{
		Style:    discordgo.SecondaryButton,
		Label:    level + "++",
		CustomID: level + "++",
		Emoji: &discordgo.ComponentEmoji{
			Name: emRsStart,
		},
	}
	buttonPl30 := discordgo.Button{
		Style:    discordgo.SecondaryButton,
		Label:    "+30",
		CustomID: level + "+++",
		Emoji: &discordgo.ComponentEmoji{
			Name: emPl30,
		},
	}

	// Создание компонентов с кнопкой
	return []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				buttonOk,
				buttonCancel,
				buttonRsStart,
				buttonPl30,
			},
		},
	}
}

func (d *Discord) DeleteMessage(chatid, mesid string) {
	_ = d.s.ChannelMessageDelete(chatid, mesid)
}
func (d *Discord) DeleteMesageSecond(chatid, mesid string, second int) {
	if second > 60 {
		d.storage.TimeDeleteMessage.TimerInsert(models.Timer{
			Dsmesid:  mesid,
			Dschatid: chatid,
			Timed:    second,
		})
	} else {
		go func() {
			time.Sleep(time.Duration(second) * time.Second)
			err := d.s.ChannelMessageDelete(chatid, mesid)
			if err != nil {
				fmt.Println("Ошибка удаления сообщения дискорда ", chatid, mesid, second)
			}
		}()
	}
}
func (d *Discord) EditComplexButton(dsmesid, dschatid string, Embeds discordgo.MessageEmbed, component []discordgo.MessageComponent) error {
	_, err := d.s.ChannelMessageEditComplex(&discordgo.MessageEdit{
		Content:    &mesContentNil,
		Embed:      &Embeds,
		ID:         dsmesid,
		Channel:    dschatid,
		Components: &component,
	})
	if err != nil {
		return err
	}
	return nil
}
func (d *Discord) Subscribe(nameid, argRoles, guildid string) int {
	g, err := d.s.State.Guild(guildid)
	if err != nil {
		d.log.ErrorErr(err)
		g, err = d.s.Guild(guildid)
		if err != nil {
			d.log.ErrorErr(err)
		}
	}

	exist, role := d.roleExists(g, argRoles)

	if !exist { //если нет роли
		role = d.createRole(argRoles, guildid)
	}

	member, err := d.s.GuildMember(guildid, nameid)
	if err != nil {
		d.log.ErrorErr(err)
	}
	var subscribe int = 0
	if exist {
		for _, _role := range member.Roles {
			if _role == role.ID {
				subscribe = 1
			}
		}
	}

	err = d.s.GuildMemberRoleAdd(guildid, nameid, role.ID)
	if err != nil {
		d.log.ErrorErr(err)
		subscribe = 2
	}

	return subscribe
}
func (d *Discord) Unsubscribe(nameid, argRoles, guildid string) int {
	var unsubscribe int = 0
	g, err := d.s.State.Guild(guildid)
	if err != nil {
		d.log.ErrorErr(err)
		g, err = d.s.Guild(guildid)
		if err != nil {
			d.log.ErrorErr(err)
		}
	}

	exist, role := d.roleExists(g, argRoles)
	if !exist { //если нет роли
		unsubscribe = 1
	}

	member, err := d.s.GuildMember(guildid, nameid)
	if err != nil {
		d.log.ErrorErr(err)
	}
	if exist {
		for _, _role := range member.Roles {
			if _role == role.ID {
				unsubscribe = 2
			}
		}
	}
	if unsubscribe == 2 {
		err = d.s.GuildMemberRoleRemove(guildid, nameid, role.ID)
		if err != nil {
			d.log.ErrorErr(err)
			unsubscribe = 3
		}
	}

	return unsubscribe
}
func (d *Discord) EditWebhook(text, username, chatID, mID string, guildID, avatarURL string) {
	if text == "" {
		return
	}

	web := transmitter.New(d.s, guildID, "KzBot", true, d.log)
	params := &discordgo.WebhookParams{
		Content:   text,
		Username:  username,
		AvatarURL: avatarURL,
	}
	err := web.Edit(chatID, mID, params)
	if err != nil {
		return
	}
}
func (d *Discord) EmbedDS(mapa map[string]string, numkz int, count int, dark bool) discordgo.MessageEmbed {
	textcount := ""
	if count == 1 {
		textcount = fmt.Sprintf("\n1️⃣ %s \n\n",
			mapa["name1"])
	} else if count == 2 {
		textcount = fmt.Sprintf("\n1️⃣ %s \n2️⃣ %s \n\n",
			mapa["name1"], mapa["name2"])
	} else if count == 3 {
		textcount = fmt.Sprintf("\n1️⃣ %s \n2️⃣ %s \n3️⃣ %s \n\n",
			mapa["name1"], mapa["name2"], mapa["name3"])
	} else {
		textcount = fmt.Sprintf("\n1️⃣ %s \n2️⃣ %s \n3️⃣ %s \n4️⃣ %s \n",
			mapa["name1"], mapa["name2"], mapa["name3"], mapa["name4"])
	}
	title := d.storage.Words.GetWords(mapa["lang"], "ocheredKz")
	if dark {
		title = d.storage.Words.GetWords(mapa["lang"], "ocheredTKz")
	}
	return discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{},
		Color:  16711680,
		Description: fmt.Sprintf("👇 %s <:rs:918545444425072671> %s (%d) ",
			d.storage.Words.GetWords(mapa["lang"], "jelaushieNa"), mapa["lvlkz"], numkz) +
			textcount,

		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Name: fmt.Sprintf(" %s %s\n%s %s\n%s %s",
					emOK, d.storage.Words.GetWords(mapa["lang"], "DlyaDobavleniya"),
					emCancel, d.storage.Words.GetWords(mapa["lang"], "DlyaVihodaIz"),
					emRsStart, d.storage.Words.GetWords(mapa["lang"], "prinuditelniStart")),
				Value:  d.storage.Words.GetWords(mapa["lang"], "DannieObnovleni") + ": ",
				Inline: true,
			}},
		Timestamp: time.Now().Format(time.RFC3339), // ТЕКУЩЕЕ ВРЕМЯ ДИСКОРДА
		Title:     title,
	}
}
func (d *Discord) ChannelTyping(ChannelID string) {
	err := d.s.ChannelTyping(ChannelID)
	if err != nil {
		d.log.ErrorErr(err)
		return
	}
}
