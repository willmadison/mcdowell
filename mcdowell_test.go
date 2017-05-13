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
