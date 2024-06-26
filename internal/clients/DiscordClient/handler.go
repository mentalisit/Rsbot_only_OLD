package DiscordClient

import (
	"Rsbot_only/internal/models"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"time"
)

func (d *Discord) messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	if m.Message.WebhookID != "" {
		return
	}

	d.logicMix(m)

}
func (d *Discord) messageReactionAdd(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	message, err := s.ChannelMessage(r.ChannelID, r.MessageID)
	if err != nil {
		channel, err1 := s.Channel(r.ChannelID)
		if err1 != nil {
			d.log.Error(err1.Error())
			return
		}
		user, err2 := s.User(r.UserID)
		if err2 != nil {
			d.log.Error(err2.Error())
			return
		}
		d.log.Info(fmt.Sprintln(channel.Name, r.Emoji.Name, user.Username, err.Error()))
		return
	}

	if message.Author.ID == s.State.User.ID {
		d.readReactionQueue(r, message)
	}
}

func (d *Discord) slash(s *discordgo.Session, i *discordgo.InteractionCreate) {

	switch i.Type {

	case discordgo.InteractionApplicationCommand:
		{
			switch i.ApplicationCommandData().Name {
			case "module":
				// Обработка вашей слеш-команды
				d.handleModuleCommand(i)
			case "weapon":
				d.handleWeaponCommand(i)
			}
			commandHandlers := d.addSlashHandler()
			if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
				h(s, i)
			}
		}
	case discordgo.InteractionMessageComponent:
		d.handleButtonPressed(i)

	default:
		fmt.Printf("slash %+v\n", i.Type)
	}

}

func (d *Discord) ready() {
	commands := d.addSlashCommand()
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := d.s.ApplicationCommandCreate(d.s.State.User.ID, "", v)
		if err != nil {

			d.log.ErrorErr(err)
		}
		registeredCommands[i] = cmd
	}
}

func (d *Discord) removeCommand(guildid string) {
	registeredCommands, err := d.s.ApplicationCommands(d.s.State.User.ID, guildid)
	if err != nil {
		d.log.Fatal(err.Error())
	}

	for _, v := range registeredCommands {
		err = d.s.ApplicationCommandDelete(d.s.State.User.ID, guildid, v.ID)
		if err != nil {
			d.log.ErrorErr(err)
		}
	}
	fmt.Println("удалены")
}
func (d *Discord) addSlashCommand() []*discordgo.ApplicationCommand {
	return []*discordgo.ApplicationCommand{
		{
			Name:        "help",
			Description: "Общая справка",
		},
		{
			Name:        "helpqueue",
			Description: "Очередь КЗ",
		},
		{
			Name:        "helpnotification",
			Description: "Уведомления",
		},
		{
			Name:        "helpevent",
			Description: "Ивент КЗ",
		},
		{
			Name:        "helptop",
			Description: "ТОП лист",
		},
		{
			Name:        "helpicon",
			Description: "Работа с иконками",
		},
		{
			Name:        "module",
			Description: "Выберите нужный модуль и уровень / Select the desired module and level",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "module",
					Description: "Выберите модуль / Select module",
					Required:    true,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "Ингибитор КЗ / RSE",
							Value: "RSE",
						},
						{
							Name:  "Генезис / Genesis",
							Value: "GENESIS",
						},
						{
							Name:  "Обогатить / Enrich",
							Value: "ENRICH",
						},
						// Добавьте другие модули по мере необходимости
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "level",
					Description: "Выберите уровень / Select level",
					Required:    true,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "Уровень / Level 0",
							Value: 0,
						},
						{
							Name:  "Уровень / Level 1",
							Value: 1,
						}, {
							Name:  "Уровень / Level 2",
							Value: 2,
						}, {
							Name:  "Уровень / Level 3",
							Value: 3,
						}, {
							Name:  "Уровень / Level 4",
							Value: 4,
						}, {
							Name:  "Уровень / Level 5",
							Value: 5,
						}, {
							Name:  "Уровень / Level 6",
							Value: 6,
						}, {
							Name:  "Уровень / Level 7",
							Value: 7,
						}, {
							Name:  "Уровень / Level 8",
							Value: 8,
						}, {
							Name:  "Уровень / Level 9",
							Value: 9,
						}, {
							Name:  "Уровень / Level 10",
							Value: 10,
						}, {
							Name:  "Уровень / Level 11",
							Value: 11,
						}, {
							Name:  "Уровень / Level 12",
							Value: 12,
						}, {
							Name:  "Уровень / Level 13",
							Value: 13,
						}, {
							Name:  "Уровень / Level 14",
							Value: 14,
						}, {
							Name:  "Уровень / Level 15",
							Value: 15,
						},
						// Добавьте другие уровни по мере необходимости
					},
				},
			},
		},
		{
			Name:        "weapon",
			Description: "Выберите основное оружие / Select your main weapon",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "weapon",
					Description: "Выберите оружие / Select weapon",
					Required:    true,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "Артобстрел / Barrage",
							Value: "barrage",
						},
						{
							Name:  "Лазер / Laser",
							Value: "laser",
						},
						{
							Name:  "Цепной луч / Chain ray",
							Value: "chainray",
						},
						{
							Name:  "Батарея / Battery",
							Value: "battery",
						},
						{
							Name:  "Залповая батарея / Mass battery",
							Value: "massbattery",
						},
						{
							Name:  "Пусковая установка / Dart launcher",
							Value: "dartlauncher",
						},
						{
							Name:  "Ракетная установка / Rocket launcher",
							Value: "rocketlauncher",
						},
						{
							Name:  "Удалить оружие / Remove weapon",
							Value: "Remove",
						},
						// Добавьте другие модули по мере необходимости
					},
				},
			},
		},
	}
}
func (d *Discord) addSlashHandler() map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"help": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{

					Content: models.Help,
				},
			})
			go func() {
				time.Sleep(1 * time.Minute)
				s.InteractionResponseDelete(i.Interaction)
			}()
		},
		"helpqueue": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: models.HelpQueue,
				},
			})
			go func() {
				time.Sleep(1 * time.Minute)
				s.InteractionResponseDelete(i.Interaction)
			}()
		},
		"helpnotification": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Уведомления:\n" +
						"	Подписаться на уведомления о начале очереди: +[4-11]\n" +
						"+10 -подписаться на уведомления о начале очереди на КЗ 10ур.\n\n" +
						"	Подписаться на уведомление, если в очереди 3 человека: ++[4-11]\n" +
						"++10 -подписаться на уведомления о наличии 3х человек в очереди на КЗ 10ур.\n\n" +
						"	Отключить уведомления о начале сбора: -[5-11]\n" +
						"-9 -отключить уведомления о начале сборе на КЗ 9ур.\n\n" +
						"	Отключить уведомления 3/4 в очереди: --[5-11]\n" +
						"--9 -отключить уведомления о наличии 3х человек в очереди на КЗ 9ур.",
				},
			})
			go func() {
				time.Sleep(1 * time.Minute)
				s.InteractionResponseDelete(i.Interaction)
			}()
		},
		"helpevent": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: models.HelpEvent,
				},
			})
			go func() {
				time.Sleep(1 * time.Minute)
				s.InteractionResponseDelete(i.Interaction)
			}()
		},
		"helptop": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: models.HelpTop,
				},
			})
			go func() {
				time.Sleep(1 * time.Minute)
				s.InteractionResponseDelete(i.Interaction)
			}()
		},
		"helpicon": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: models.HelpIcon,
				},
			})
			go func() {
				time.Sleep(1 * time.Minute)
				s.InteractionResponseDelete(i.Interaction)
			}()
		},
	}

	return commandHandlers
}
