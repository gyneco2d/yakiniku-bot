package main

import (
	"encoding/json"
	"github.com/nlopes/slack"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type Config struct {
	BotUserToken string `json:"botUserToken"`
}

var (
	botId	string
	botName	string
)

func run(api *slack.Client) int {
	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for {
		select {
		case msg := <-rtm.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *slack.ConnectedEvent:
				botId = ev.Info.User.ID
				botName = ev.Info.User.Name

			case *slack.HelloEvent:
				log.Print("Hello, Event")

			case *slack.MessageEvent:
				text := ev.Text
				if ev.Type == "message" && strings.HasPrefix(text, "<@"+botId+">") {
					log.Printf("Message: %v\n", ev)
					rtm.SendMessage(rtm.NewOutgoingMessage("Hello world", ev.Channel))
				}

			case *slack.InvalidAuthEvent:
				log.Print("Invalid credentials")
				return 1

			}
		}
	}
}

func main() {
	jsonString, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatal(err)
	}
	c := new(Config)
	err = json.Unmarshal(jsonString, c)
	if err != nil {
		log.Println(err)
		os.Exit(2)
	}

	api := slack.New(c.BotUserToken)
	os.Exit(run(api))
}
