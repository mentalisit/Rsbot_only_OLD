package transmitter

import (
	"errors"
	"fmt"
	"github.com/mentalisit/logger"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

// A Transmitter represents a message manager for a single guild.
type Transmitter struct {
	session    *discordgo.Session
	guild      string
	title      string
	autoCreate bool

	// channelWebhooks maps from a channel ID to a webhook instance
	channelWebhooks map[string]*discordgo.Webhook

	mutex sync.RWMutex

	log *logger.Logger
}

// ErrWebhookNotFound is returned when a valid webhook for this channel/message combination does not exist
var ErrWebhookNotFound = errors.New("webhook for this channel and message does not exist")

func New(session *discordgo.Session, guild string, title string, autoCreate bool, log *logger.Logger) *Transmitter {
	return &Transmitter{
		session:    session,
		guild:      guild,
		title:      title,
		autoCreate: autoCreate,

		channelWebhooks: make(map[string]*discordgo.Webhook),

		log: log,
	}
}

// Send transmits a message to the given channel with the provided webhook data, and waits until Discord responds with message data.
func (t *Transmitter) Send(channelID string, params *discordgo.WebhookParams) (*discordgo.Message, error) {
	wh, err := t.getOrCreateWebhook(channelID)
	if err != nil {
		t.log.ErrorErr(err)
		return nil, err
	}

	msg, err := t.session.WebhookExecute(wh.ID, wh.Token, true, params)
	if err != nil {
		return nil, fmt.Errorf("execute failed: %w", err)
	}

	return msg, nil
}

// Edit will edit a message in a channel, if possible.
func (t *Transmitter) Edit(channelID string, messageID string, params *discordgo.WebhookParams) error {
	wh := t.getWebhook(channelID)

	if wh == nil {
		return ErrWebhookNotFound
	}

	uri := discordgo.EndpointWebhookToken(wh.ID, wh.Token) + "/messages/" + messageID
	_, err := t.session.RequestWithBucketID("PATCH", uri, params, discordgo.EndpointWebhookToken("", ""))
	if err != nil {
		return err
	}

	return nil
}

func (t *Transmitter) createWebhook(channel string) (*discordgo.Webhook, error) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	wh, err := t.session.WebhookCreate(channel, t.title+time.Now().Format(" 3:04:05PM"), "")
	if err != nil {
		t.log.ErrorErr(err)
		return nil, err
	}
	t.channelWebhooks[channel] = wh
	return wh, nil
}

func (t *Transmitter) getWebhook(channel string) *discordgo.Webhook {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	//t.log.Println("gethook 178")
	webhooks, err := t.session.ChannelWebhooks(channel)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	var webhook *discordgo.Webhook
	for _, i := range webhooks {
		//fmt.Printf("%s\n", i.User.Username)
		if i.User.Bot && i.User.Username == t.session.State.User.Username {
			webhook = i
			return webhook
		}
	}
	if webhook == nil {
		webhookCreate, err1 := t.session.WebhookCreate(channel, t.title, "")
		if err1 != nil {
			fmt.Println(err1)
			return nil
		}
		fmt.Println("webhookCreateNILL", webhookCreate.Token, webhookCreate.ID)
		return webhookCreate
	}

	if len(webhooks) == 0 {
		webhookCreate, err1 := t.session.WebhookCreate(channel, t.title, "")
		if err1 != nil {
			fmt.Println(err1)
			return nil
		}
		fmt.Println("webhookCreate", webhookCreate.Token, webhookCreate.ID)
		return webhookCreate
	}

	return nil
}

func (t *Transmitter) getOrCreateWebhook(channelID string) (*discordgo.Webhook, error) {
	// If we have a webhook for this channel, immediately return it
	wh := t.getWebhook(channelID)
	if wh != nil {
		return wh, nil
	}
	//t.log.Println(209)
	// Early exit if we don't want to automatically create one
	if !t.autoCreate {
		return nil, ErrWebhookNotFound
	}
	//t.log.Println(214)
	t.log.Info("Creating a webhook for " + channelID)
	wh, err := t.createWebhook(channelID)
	if err != nil {
		t.log.ErrorErr(err)
		return nil, fmt.Errorf("could not create webhook: %w", err)
	}

	return wh, nil
}
