package DiscordClient

import (
	"Rsbot_only/internal/models"
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strconv"
	"time"
)

// slash command module respond
func (d *Discord) handleModuleCommand(i *discordgo.InteractionCreate) {
	module := i.ApplicationCommandData().Options[0].StringValue()
	level := i.ApplicationCommandData().Options[1].IntValue()

	response := fmt.Sprintf("Выбран модуль: %s, уровень: %d", module, level)
	if level == 0 {
		response = fmt.Sprintf("Удален модуль: %s, уровень: %d", module, level)
	}
	// Отправка ответа
	err := d.s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: response,
		},
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	go func() {
		time.Sleep(20 * time.Second)
		err = d.s.InteractionResponseDelete(i.Interaction)
		if err != nil {
			return
		}
	}()
	d.updateModuleOrWeapon(i.Interaction.Member.User.Username, module, strconv.FormatInt(level, 10))
}

// slash command weapon respond
func (d *Discord) handleWeaponCommand(i *discordgo.InteractionCreate) {
	weapon := i.ApplicationCommandData().Options[0].StringValue()

	response := fmt.Sprintf("Установлено оружие: %s", weapon)

	// Отправка ответа
	err := d.s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: response,
		},
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	go func() {
		time.Sleep(20 * time.Second)
		err = d.s.InteractionResponseDelete(i.Interaction)
		if err != nil {
			return
		}
	}()
	d.updateModuleOrWeapon(i.Interaction.Member.User.Username, weapon, "")
}

func (d *Discord) updateModuleOrWeapon(username, module, level string) {
	rse := "<:rse:1199068829511335946> " + level
	genesis := "<:genesis:1199068748280242237> " + level
	enrich := "<:enrich:1199068793633251338> " + level
	if level == "0" {
		rse, genesis, enrich = "", "", ""
	}

	barrage := "<:barrage:1199084425393225782>"
	laser := "<:laser:1199084197571207339>"
	chainray := "<:chainray:1199073579577376888>"
	battery := "<:batteryw:1199072534562345021>"
	massbattery := "<:massbattery:1199072493760151593>"
	dartlauncher := "<:dartlauncher:1199072434674991145>"
	rocketlauncher := "<:rocketlauncher:1199071677548605562>"
	t := d.storage.Emoji.EmojiModuleReadUsers(context.Background(), username, "ds")
	if len(t.Name) == 0 {
		d.storage.Emoji.EmInsertEmpty(context.Background(), "ds", username)
	}
	switch module {
	case "RSE":
		d.storage.Emoji.ModuleUpdate(context.Background(), username, "ds", "1", rse)
	case "GENESIS":
		d.storage.Emoji.ModuleUpdate(context.Background(), username, "ds", "2", genesis)
	case "ENRICH":
		d.storage.Emoji.ModuleUpdate(context.Background(), username, "ds", "3", enrich)
	case "barrage":
		d.storage.Emoji.WeaponUpdate(context.Background(), username, "ds", barrage)
	case "laser":
		d.storage.Emoji.WeaponUpdate(context.Background(), username, "ds", laser)
	case "chainray":
		d.storage.Emoji.WeaponUpdate(context.Background(), username, "ds", chainray)
	case "battery":
		d.storage.Emoji.WeaponUpdate(context.Background(), username, "ds", battery)
	case "massbattery":
		d.storage.Emoji.WeaponUpdate(context.Background(), username, "ds", massbattery)
	case "dartlauncher":
		d.storage.Emoji.WeaponUpdate(context.Background(), username, "ds", dartlauncher)
	case "rocketlauncher":
		d.storage.Emoji.WeaponUpdate(context.Background(), username, "ds", rocketlauncher)
	case "Remove":
		d.storage.Emoji.WeaponUpdate(context.Background(), username, "ds", "")
	}
}

func (d *Discord) handleButtonPressed(i *discordgo.InteractionCreate) {
	ok, config := d.CheckChannelConfigDS(i.ChannelID)
	if ok {
		in := models.InMessage{
			Mtext:       i.MessageComponentData().CustomID,
			Tip:         "ds",
			Name:        i.Interaction.Member.User.Username,
			NameMention: i.Interaction.Member.User.Mention(),
			Ds: struct {
				Mesid   string
				Nameid  string
				Guildid string
				Avatar  string
			}{
				Mesid:   i.Interaction.Message.ID,
				Nameid:  i.Interaction.Member.User.ID,
				Guildid: i.Interaction.GuildID,
				Avatar:  i.Interaction.Member.User.AvatarURL("128"),
			},
			Config: config,
			Option: models.Option{Reaction: true},
		}
		d.ChanRsMessage <- in
		err := d.s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredMessageUpdate,
		})
		if err != nil {
			d.log.ErrorErr(err)
			return
		}
	}
}
