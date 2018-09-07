package main

import (
  "fmt"
  "os"
  "io/ioutil"
  "log"
  "encoding/json"
	"github.com/nlopes/slack"
)

type Config struct {
  BotUserToken string `json:"botUserToken"`
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

  fmt.Printf("botUserToken：%#v", c.BotUserToken)

	api := slack.New(c.BotUserToken)

	rtm := api.NewRTM()

	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			if ev.Text == "ping" {
				rtm.SendMessage(rtm.NewOutgoingMessage("pong", ev.Channel))
			}
		}
	}
}
