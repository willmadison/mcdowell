package mcdowell

import (
	"context"

	"log"

	"strings"

	"github.com/nlopes/slack"
	"github.com/pkg/errors"
)

type (
	// Bot represents a single bot instance.
	Bot struct {
		id           string
		name         string
		client       SlackClient
		ctx          context.Context
		contributors map[string]string
		Debug        bool
		Testing      bool
	}

	// SlackClient represents the interface of methods we rely on from the Slack client.
	SlackClient interface {
		PostMessage(channel, message string, params slack.PostMessageParameters) (string, string, error)
		GetUsers() ([]slack.User, error)
	}
)

func (b *Bot) initialize() error {
	if b.Debug {
		log.Println("determining bot/contributor user IDs:")
	}

	users, err := b.client.GetUsers()
	if err != nil {
		return errors.WithStack(err)
	}

	b.contributors = map[string]string{}

	for _, user := range users {
		switch user.Name {
		case "willmadison", "xango":
			b.contributors[user.Name] = user.ID
		case b.name:
			if user.IsBot {
				b.id = user.ID
			}
		default:
			continue
		}
	}

	if b.Debug {
		log.Println("contributors:", b.contributors)
	}

	if b.id == "" && !b.Testing {
		return errors.New("could not find bot in the list of names, ensure the bot is called \"" + b.name + "\" ")
	}

	return nil
}

// OnTeamJoined handles the appropriate behavior for when new team members join our slack.
func (b *Bot) OnTeamJoined(event *slack.TeamJoinEvent) error {
	message := `Yo ` + event.User.Name + `!

I’d like to welcome you to the Atlanta Black Tech Family. Our mission is to improve the quality, quantity, and connections for people of African descent within the overall Metro Atlanta tech ecosystem.

Please click on “Channels” to browse all of our sub-communities, and join the ones that are most relevant to you. Enjoy your time, and help us build the communities by inviting others in your network.`

	params := slack.PostMessageParameters{AsUser: true, LinkNames: 1}
	_, _, err := b.client.PostMessage(event.User.ID, message, params)

	return err
}

var botEventTextToResponses = map[string]func(*Bot, *slack.MessageEvent) error{
	"show me the money":     heHasHisOwnMoney("The boy has got his own money!"),
	"let me hold something": heHasHisOwnMoney("I got you!"),
	"soul glo":              soulGlo,
}

func heHasHisOwnMoney(message string) func(*Bot, *slack.MessageEvent) error {
	return func(b *Bot, event *slack.MessageEvent) error {
		params := slack.PostMessageParameters{AsUser: true, UnfurlLinks: true}
		params.Attachments = []slack.Attachment{
			{
				Text:     message,
				ImageURL: "https://novembrepleut.files.wordpress.com/2011/06/zamundamoney_100.png",
			},
		}
		_, _, err := b.client.PostMessage(event.Channel, "", params)
		return err
	}
}

func soulGlo(b *Bot, event *slack.MessageEvent) error {
	params := slack.PostMessageParameters{AsUser: true, UnfurlLinks: true}
	params.Attachments = []slack.Attachment{
		{
			ImageURL: "https://media.giphy.com/media/3Gz3vy81HkDa8/giphy.gif",
		},
	}
	_, _, err := b.client.PostMessage(event.Channel, "", params)
	return err
}

// OnNewMessage handles the appropriate behavior for when new interesting
// messages happen in any channel the bot is listening in.
func (b *Bot) OnNewMessage(event *slack.MessageEvent) error {
	if event.BotID != "" || event.User == "" || event.SubType == "bot_message" {
		return nil
	}

	eventText := strings.Trim(strings.ToLower(event.Text), " \n\r")

	if b.Debug || b.Testing {
		log.Printf("event: %+v\n", *event)
		log.Println("got message:", eventText)
	}

	var err error
	for fragment, response := range botEventTextToResponses {
		if strings.Contains(eventText, fragment) {
			err = response(b, event)
		}
	}

	return err
}

// NewBot returns a new McDowell Bot instance ready to handle any events from Slack.
func NewBot(ctx context.Context, client SlackClient, options ...func(*Bot)) (*Bot, error) {
	b := &Bot{
		ctx:    ctx,
		client: client,
		name:   "mcdowell",
	}

	for _, option := range options {
		option(b)
	}

	err := b.initialize()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return b, nil
}

// WithDebug enables debug mode on the bot.
func WithDebug() func(*Bot) {
	return func(b *Bot) {
		b.Debug = true
	}
}

// WithTesting enables test mode on the bot.
func WithTesting() func(*Bot) {
	return func(b *Bot) {
		b.Testing = true
	}
}
