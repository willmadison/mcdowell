package mcdowell_test

import (
	"context"
	"testing"

	"github.com/nlopes/slack"
	"github.com/willmadison/mcdowell"
)

type mockSlack struct {
	lastChannel, lastMessage string
	params                   slack.PostMessageParameters
}

func (m *mockSlack) PostMessage(channel, message string, params slack.PostMessageParameters) (string, string, error) {
	m.lastChannel = channel
	m.lastMessage = message
	m.params = params
	return "", "", nil
}

func (m *mockSlack) GetUsers() ([]slack.User, error) {
	return []slack.User{}, nil
}

func TestBotHandlesTeamJoinEvents(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client := &mockSlack{}

	m, err := mcdowell.NewBot(ctx, client, mcdowell.WithTesting())
	if err != nil {
		t.Fatal("encountered an unexpected error creating a new McDowell instance:", err)
	}

	testUser := slack.User{
		ID:   "dummyID",
		Name: "Test User",
	}

	e := &slack.TeamJoinEvent{
		User: testUser,
	}

	err = m.OnTeamJoined(e)
	if err != nil {
		t.Fatal("encountered an unexpected error handling an OnTeamJoined event:", err)
	}

	if client.lastChannel != testUser.ID {
		t.Error("sent message to", client.lastChannel, "expected message to go to", testUser.ID)
	}

	expected := `Yo Test User!

I’d like to welcome you to the Atlanta Black Tech Family. Our mission is to improve the quality, quantity, and connections for people of African descent within the overall Metro Atlanta tech ecosystem.

Please click on “Channels” to browse all of our sub-communities, and join the ones that are most relevant to you. Enjoy your time, and help us build the communities by inviting others in your network.`
	if client.lastMessage != expected {
		t.Error("got:", client.lastMessage, "wanted:", expected)
	}

	if !client.params.AsUser {
		t.Error("expected message to be posted as the Cleo McDowell!")
	}

	if client.params.LinkNames == 0 {
		t.Error("got LinkNames =", client.params.LinkNames, "expected at least 1!")
	}
}

func TestShowMeTheMoney(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client := &mockSlack{}

	m, err := mcdowell.NewBot(ctx, client, mcdowell.WithTesting())
	if err != nil {
		t.Fatal("encountered an unexpected error creating a new McDowell instance:", err)
	}

	message := slack.Msg{
		Channel: "#general",
		User:    "willmadison",
		Text:    "They gonna have to show me the money!",
	}

	e := &slack.MessageEvent{
		Msg: message,
	}

	err = m.OnNewMessage(e)
	if err != nil {
		t.Fatal("encountered an unexpected error handling an OnTeamJoined event:", err)
	}

	if client.lastChannel != e.Channel {
		t.Error("sent message to", client.lastChannel, "expected message to go to", e.Channel)
	}

	expected := "https://novembrepleut.files.wordpress.com/2011/06/zamundamoney_100.png"
	if !client.params.UnfurlLinks {
		t.Error("expected links to be unfurled!")
	}

	if len(client.params.Attachments) == 0 {
		t.Fatal("expected there to be at least one attachment!")
	}

	if client.params.Attachments[0].ImageURL != expected {
		t.Error("incorrect Soul Glo reaction url!")
	}

	if client.params.Attachments[0].Text != "The boy has got his own money!" {
		t.Fatal("incorrect Soul Glo reaction url!")
	}
}

func TestLetMeHoldSomething(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client := &mockSlack{}

	m, err := mcdowell.NewBot(ctx, client, mcdowell.WithTesting())
	if err != nil {
		t.Fatal("encountered an unexpected error creating a new McDowell instance:", err)
	}

	message := slack.Msg{
		Channel: "#general",
		User:    "willmadison",
		Text:    "Hey bruh let me hold something...",
	}

	e := &slack.MessageEvent{
		Msg: message,
	}

	err = m.OnNewMessage(e)
	if err != nil {
		t.Fatal("encountered an unexpected error handling an OnTeamJoined event:", err)
	}

	if client.lastChannel != e.Channel {
		t.Error("sent message to", client.lastChannel, "expected message to go to", e.Channel)
	}

	if !client.params.AsUser {
		t.Error("expected message to be posted as the Cleo McDowell!")
	}

	if !client.params.UnfurlLinks {
		t.Error("expected links to be unfurled!")
	}

	if len(client.params.Attachments) == 0 {
		t.Fatal("expected there to be at least one attachment!")
	}

	expected := "https://novembrepleut.files.wordpress.com/2011/06/zamundamoney_100.png"
	if client.params.Attachments[0].ImageURL != expected {
		t.Fatal("incorrect Soul Glo reaction url!")
	}

	if client.params.Attachments[0].Text != "I got you!" {
		t.Fatal("incorrect Soul Glo reaction url!")
	}
}

func TestSoulGlo(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client := &mockSlack{}

	m, err := mcdowell.NewBot(ctx, client, mcdowell.WithTesting())
	if err != nil {
		t.Fatal("encountered an unexpected error creating a new McDowell instance:", err)
	}

	message := slack.Msg{
		Channel: "#general",
		User:    "willmadison",
		Text:    "let your soul glow",
	}

	e := &slack.MessageEvent{
		Msg: message,
	}

	err = m.OnNewMessage(e)
	if err != nil {
		t.Fatal("encountered an unexpected error handling an OnTeamJoined event:", err)
	}

	if client.lastChannel != e.Channel {
		t.Error("sent message to", client.lastChannel, "expected message to go to", e.Channel)
	}

	expected := "https://media.giphy.com/media/3Gz3vy81HkDa8/giphy.gif"

	if !client.params.AsUser {
		t.Error("expected message to be posted as the Cleo McDowell!")
	}

	if !client.params.UnfurlLinks {
		t.Error("expected links to be unfurled!")
	}

	if len(client.params.Attachments) == 0 {
		t.Fatal("expected there to be at least one attachment!")
	}

	if client.params.Attachments[0].ImageURL != expected {
		t.Fatal("incorrect Soul Glo reaction url! got:", client.params.Attachments[0].ImageURL)
	}
}

func TestIgnoresBotMessages(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client := &mockSlack{}

	m, err := mcdowell.NewBot(ctx, client, mcdowell.WithTesting())
	if err != nil {
		t.Fatal("encountered an unexpected error creating a new McDowell instance:", err)
	}

	message := slack.Msg{
		BotID:   "someId",
		Channel: "#general",
		User:    "mcdowell",
		Text:    "They gonna have to show me the money!",
		SubType: "bot_message",
	}

	e := &slack.MessageEvent{
		Msg: message,
	}

	err = m.OnNewMessage(e)
	if err != nil {
		t.Fatal("encountered an unexpected error handling an OnTeamJoined event:", err)
	}

	if len(client.params.Attachments) > 0 {
		t.Error("got:", len(client.params.Attachments), "attachments, expected no interaction.")
	}
}
