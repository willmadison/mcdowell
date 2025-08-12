package mcdowell_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/nlopes/slack"
	"github.com/stretchr/testify/assert"
	"github.com/willmadison/mcdowell"
)

type captured struct {
	Path        string
	ContentType string
	Body        []byte
	JSON        map[string]any
	Form        url.Values
}

// startFakeSlack returns a test server that records requests and responds OK
func startFakeSlack(t *testing.T) (*httptest.Server, *captured) {
	t.Helper()
	var cap captured

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cap.Path = r.URL.Path
		cap.ContentType = r.Header.Get("Content-Type")
		b, _ := io.ReadAll(r.Body)
		cap.Body = b

		// slack-go may send JSON or urlencoded; handle both
		if strings.Contains(cap.ContentType, "application/json") {
			_ = json.Unmarshal(b, &cap.JSON)
		} else {
			// Either form or "payload=JSON"
			if v, err := url.ParseQuery(string(b)); err == nil {
				cap.Form = v
				if p := v.Get("payload"); p != "" {
					_ = json.Unmarshal([]byte(p), &cap.JSON)
				}
			}
		}

		w.Header().Set("Content-Type", "application/json")

		// minimal OK reply for chat.postMessage
		w.Write([]byte(`{"ok":true,"channel":"C123","ts":"123.456","message":{}}`))
	}))

	return srv, &cap
}

func TestBotHandlesTeamJoinEvents(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	srv, captured := startFakeSlack(t)
	t.Cleanup(srv.Close)

	client := slack.New("dummyToken", slack.OptionAPIURL(srv.URL+"/"))

	m, err := mcdowell.NewBot(ctx, client, mcdowell.WithTesting())
	assert.Nil(t, err)

	testUser := slack.User{
		ID:   "dummyID",
		Name: "Test User",
	}

	e := &slack.TeamJoinEvent{
		User: testUser,
	}

	err = m.OnTeamJoined(e)
	assert.Nil(t, err)

	assert.Equal(t, testUser.ID, captured.Form.Get("channel"))

	expected := `Yo Test User!

I’d like to welcome you to the Atlanta Black Tech Family. Our mission is to improve the quality, quantity, and connections for people of African descent within the overall Metro Atlanta tech ecosystem.

Please click on “Channels” to browse all of our sub-communities, and join the ones that are most relevant to you. Enjoy your time, and help us build the communities by inviting others in your network.`

	assert.Equal(t, expected, captured.Form.Get("text"))

	actual_as_user, err := strconv.ParseBool(captured.Form.Get("as_user"))
	assert.Nil(t, err)
	assert.True(t, actual_as_user)

	actual_link_names, err := strconv.Atoi(captured.Form.Get("link_names"))
	assert.Nil(t, err)
	assert.Equal(t, 1, actual_link_names)
}

func TestShowMeTheMoney(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	srv, captured := startFakeSlack(t)
	t.Cleanup(srv.Close)

	client := slack.New("dummyToken", slack.OptionAPIURL(srv.URL+"/"))

	m, err := mcdowell.NewBot(ctx, client, mcdowell.WithTesting())
	assert.Nil(t, err)

	message := slack.Msg{
		Channel: "#general",
		User:    "willmadison",
		Text:    "They gonna have to show me the money!",
	}

	e := &slack.MessageEvent{
		Msg: message,
	}

	err = m.OnNewMessage(e)
	assert.Nil(t, err)

	fmt.Printf("captured: %+v\n", *captured)

	assert.Equal(t, e.Channel, captured.Form.Get("channel"))

	actual_unfurl_links, err := strconv.ParseBool(captured.Form.Get("unfurl_links"))
	assert.Nil(t, err)
	assert.True(t, actual_unfurl_links)

	raw_actual_attachments := captured.Form.Get("attachments")

	var actual_attachments []slack.Attachment

	err = json.Unmarshal([]byte(raw_actual_attachments), &actual_attachments)
	assert.Nil(t, err)

	assert.True(t, len(actual_attachments) > 0)

	expectedURL := "https://novembrepleut.files.wordpress.com/2011/06/zamundamoney_100.png"

	assert.Equal(t, expectedURL, actual_attachments[0].ImageURL)
	assert.Equal(t, `The boy has got his own money!`, actual_attachments[0].Text)
}

func TestLetMeHoldSomething(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	srv, captured := startFakeSlack(t)
	t.Cleanup(srv.Close)

	client := slack.New("dummyToken", slack.OptionAPIURL(srv.URL+"/"))

	m, err := mcdowell.NewBot(ctx, client, mcdowell.WithTesting())
	assert.Nil(t, err)

	message := slack.Msg{
		Channel: "#general",
		User:    "willmadison",
		Text:    "Hey bruh let me hold something...",
	}

	e := &slack.MessageEvent{
		Msg: message,
	}

	err = m.OnNewMessage(e)
	assert.Nil(t, err)

	assert.Equal(t, e.Channel, captured.Form.Get("channel"))

	actual_as_user, err := strconv.ParseBool(captured.Form.Get("as_user"))
	assert.Nil(t, err)
	assert.True(t, actual_as_user)

	actual_unfurl_links, err := strconv.ParseBool(captured.Form.Get("unfurl_links"))
	assert.Nil(t, err)
	assert.True(t, actual_unfurl_links)

	raw_actual_attachments := captured.Form.Get("attachments")

	var actual_attachments []slack.Attachment

	err = json.Unmarshal([]byte(raw_actual_attachments), &actual_attachments)
	assert.Nil(t, err)

	assert.True(t, len(actual_attachments) > 0)

	expectedURL := "https://novembrepleut.files.wordpress.com/2011/06/zamundamoney_100.png"

	assert.Equal(t, expectedURL, actual_attachments[0].ImageURL)
	assert.Equal(t, "I got you!", actual_attachments[0].Text)
}

func TestSoulGlo(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	srv, captured := startFakeSlack(t)
	t.Cleanup(srv.Close)

	client := slack.New("dummyToken", slack.OptionAPIURL(srv.URL+"/"))

	m, err := mcdowell.NewBot(ctx, client, mcdowell.WithTesting())
	assert.Nil(t, err)

	message := slack.Msg{
		Channel: "#general",
		User:    "willmadison",
		Text:    "let your soul glow",
	}

	e := &slack.MessageEvent{
		Msg: message,
	}

	err = m.OnNewMessage(e)
	assert.Nil(t, err)

	assert.Equal(t, e.Channel, captured.Form.Get("channel"))

	actual_as_user, err := strconv.ParseBool(captured.Form.Get("as_user"))
	assert.Nil(t, err)
	assert.True(t, actual_as_user)

	actual_unfurl_links, err := strconv.ParseBool(captured.Form.Get("unfurl_links"))
	assert.Nil(t, err)
	assert.True(t, actual_unfurl_links)

	raw_actual_attachments := captured.Form.Get("attachments")

	var actual_attachments []slack.Attachment

	err = json.Unmarshal([]byte(raw_actual_attachments), &actual_attachments)
	assert.Nil(t, err)

	assert.True(t, len(actual_attachments) > 0)

	expectedURL := "https://media.giphy.com/media/3Gz3vy81HkDa8/giphy.gif"

	assert.Equal(t, expectedURL, actual_attachments[0].ImageURL)
	assert.Empty(t, actual_attachments[0].Text)
}

func TestQueenToBe(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	srv, captured := startFakeSlack(t)
	t.Cleanup(srv.Close)

	client := slack.New("dummyToken", slack.OptionAPIURL(srv.URL+"/"))

	m, err := mcdowell.NewBot(ctx, client, mcdowell.WithTesting())
	assert.Nil(t, err)

	message := slack.Msg{
		Channel: "#general",
		User:    "willmadison",
		Text:    "I'm just looking for a queen",
	}

	e := &slack.MessageEvent{
		Msg: message,
	}

	err = m.OnNewMessage(e)
	assert.Nil(t, err)

	assert.Equal(t, e.Channel, captured.Form.Get("channel"))

	actual_as_user, err := strconv.ParseBool(captured.Form.Get("as_user"))
	assert.Nil(t, err)
	assert.True(t, actual_as_user)

	actual_unfurl_links, err := strconv.ParseBool(captured.Form.Get("unfurl_links"))
	assert.Nil(t, err)
	assert.True(t, actual_unfurl_links)

	raw_actual_attachments := captured.Form.Get("attachments")

	var actual_attachments []slack.Attachment

	err = json.Unmarshal([]byte(raw_actual_attachments), &actual_attachments)
	assert.Nil(t, err)

	assert.True(t, len(actual_attachments) > 0)

	expectedURL := "https://img.memesuper.com/bc7ab2796bdb983d5434fc842efcee0b_coming-to-america-aha-meme-coming-to-america_500-263.gif"

	assert.Equal(t, expectedURL, actual_attachments[0].ImageURL)
}

func TestIgnoresBotMessages(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	srv, captured := startFakeSlack(t)
	t.Cleanup(srv.Close)

	client := slack.New("dummyToken", slack.OptionAPIURL(srv.URL+"/"))

	m, err := mcdowell.NewBot(ctx, client, mcdowell.WithTesting())
	assert.Nil(t, err)

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
	assert.Nil(t, err)

	raw_actual_attachments := captured.Form.Get("attachments")

	var actual_attachments []slack.Attachment

	err = json.Unmarshal([]byte(raw_actual_attachments), &actual_attachments)
	assert.NotNil(t, err)

	assert.True(t, len(actual_attachments) == 0)
}
