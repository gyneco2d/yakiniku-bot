package main

import (
	"encoding/json"
	"github.com/nlopes/slack"
	"io/ioutil"
	"log"
	"math/rand"
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
		list := getMembersList(api, rtm, channel)
		str := strings.Join(list, "\n")
		rtm.SendMessage(rtm.NewOutgoingMessage(str, channel))
	} else if strings.HasPrefix(content, "kuji") {
		list := getMembersList(api, rtm, channel)
		shuffle(list)
		rtm.SendMessage(rtm.NewOutgoingMessage(list[0], channel))
	}
}

func shuffle(data []string) {
	n := len(data)
	for i := n-1; i >= 0; i-- {
		j := rand.Intn(i + 1)
		data[i], data[j] = data[j], data[i]
	}
}

func getMembersList(api *slack.Client, rtm *slack.RTM, channel string) []string {
	channelInfo, err := api.GetChannelInfo(channel)
	if err != nil {
		log.Print(err)
		rtm.SendMessage(rtm.NewOutgoingMessage("channel not found", channel))
		return []string{}
	}

	var names []string
	for _, v := range channelInfo.Members {
		userInfo, err := api.GetUserInfo(v)
		if err != nil {
			log.Print(err)
			return []string{}
		}

		if userInfo.IsBot {
			continue
		}

		if userInfo.Profile.DisplayName == "" {
			names = append(names, userInfo.Profile.RealName)
		} else {
			names = append(names, userInfo.Profile.DisplayName)
		}
	}
	return names
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
