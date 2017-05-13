package mcdowell

import (
	"context"

	"log"

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
	message := `Welcome ` + event.User.Name + `!

We're so happy to have you as a part of the Atlanta Black Tech Family & Ecosystem.

Enjoy the community!`

	params := slack.PostMessageParameters{AsUser: true, LinkNames: 1}
	_, _, err := b.client.PostMessage(event.User.ID, message, params)

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
