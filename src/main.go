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
				log.Printf("Message: %v\n", ev)
				if ev.Type == "message" && strings.HasPrefix(ev.Text, "<@"+botId+">") {
					commands(ev, api, rtm)
				}

			case *slack.InvalidAuthEvent:
				log.Print("Invalid credentials")
				return 1

			}
		}
	}
}

func commands(ev *slack.MessageEvent, api *slack.Client, rtm *slack.RTM) {
	user := ev.User
	text := ev.Text
	channel := ev.Channel

	content := strings.TrimLeft(text, "<@"+botId+"> ")
	if strings.HasPrefix(content, "hello") {
		rtm.SendMessage(rtm.NewOutgoingMessage("Hello, " + user, ev.Channel))
	} else if strings.HasPrefix(content, "list") {
		channelInfo, err := api.GetChannelInfo(channel)
		if err != nil {
			rtm.SendMessage(rtm.NewOutgoingMessage("channel not found", channel))
			return
		}

		var names []string
		for _, v := range channelInfo.Members {
			userProfile, err := api.GetUserProfile(v, true)
			if err != nil {
				log.Print(err)
				return
			}
			names = append(names, userProfile.DisplayName)
		}

		for _, name := range names {
			rtm.SendMessage(rtm.NewOutgoingMessage(name, channel))
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
