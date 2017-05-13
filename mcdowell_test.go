package mcdowell_test

import (
	"context"
	"testing"

	"strings"

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

	prefix := "Welcome " + testUser.Name + "!"
	if !strings.HasPrefix(client.lastMessage, prefix) {
		t.Error("got:", client.lastMessage, "which didn't start with the expected preamble:", prefix)
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

	expected := "The boy has got his own money! https://novembrepleut.files.wordpress.com/2011/06/zamundamoney_100.png"
	if client.lastMessage != expected {
		t.Error("got:", client.lastMessage, "wanted:", expected)
	}

	if !client.params.AsUser {
		t.Error("expected message to be posted as the Cleo McDowell!")
	}

	if !client.params.UnfurlLinks {
		t.Error("expected message to be posted as the Cleo McDowell!")
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

	expected := "I got you! https://novembrepleut.files.wordpress.com/2011/06/zamundamoney_100.png"
	if client.lastMessage != expected {
		t.Error("got:", client.lastMessage, "wanted:", expected)
	}

	if !client.params.AsUser {
		t.Error("expected message to be posted as the Cleo McDowell!")
	}

	if !client.params.UnfurlLinks {
		t.Error("expected links to be unfurled!")
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
	if client.lastMessage != expected {
		t.Error("got:", client.lastMessage, "wanted:", expected)
	}

	if !client.params.AsUser {
		t.Error("expected message to be posted as the Cleo McDowell!")
	}

	if !client.params.UnfurlLinks {
		t.Error("expected links to be unfurled!")
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

	if client.lastChannel != "" {
		t.Error("sent message to", client.lastChannel, "expected no interaction.")
	}

	if client.lastMessage != "" {
		t.Error("got:", client.lastMessage, "expected no interaction.")
	}
}
