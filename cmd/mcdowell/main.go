package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"context"

	"github.com/gorilla/mux"
	"github.com/nlopes/slack"
	"github.com/willmadison/mcdowell"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	botToken := os.Getenv("ABT_SLACK_BOT_TOKEN")
	devMode := os.Getenv("ABT_SLACK_BOT_DEV_MODE") == "true"

	options := []func(*mcdowell.Bot){}

	if devMode {
		options = append(options, mcdowell.WithDebug())
	}

	if botToken == "" {
		log.Fatalln("slack bot token is required for proper operation!")
	}

	client := slack.New(botToken)
	rtm := client.NewRTM()
	go rtm.ManageConnection()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bot, err := mcdowell.NewBot(ctx, client, options...)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		log.Println("listening for incoming events from slack...")

		for msg := range rtm.IncomingEvents {
			switch message := msg.Data.(type) {
			case *slack.MessageEvent:
				go bot.OnNewMessage(message)
			case *slack.TeamJoinEvent:
				go bot.OnTeamJoined(message)
			}
		}
	}()

	go func() {
		r := mux.NewRouter()

		r.HandleFunc("/health", func(w http.ResponseWriter, request *http.Request) {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, `{
				"botVersion": "Tip"
			}`)
		}).Name("healthCheck").Methods("GET")

		s := http.Server{
			Addr:         ":8088",
			Handler:      r,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
		}

		log.Println("serving healthCheck request(s) on", s.Addr)
		log.Fatal(s.ListenAndServe())
	}()

	log.Println("McDowell's is now open for business!!!")
	select {}
}
