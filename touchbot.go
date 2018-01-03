package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/nlopes/slack"
)

func atMentioned(message string, userID string) bool  {
  prefix := fmt.Sprintf("<@%s> ", userID)
  return strings.HasPrefix(message, prefix)
}

func main() {
	token := os.Getenv("SLACK_TOKEN")

  if token == "" {
    os.Stderr.WriteString("Error: Missing environment variable SLACK_TOKEN.")
    os.Exit(1)
  }

	api := slack.New(token)
	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for {
		select {
		case msg := <-rtm.IncomingEvents:
			switch ev := msg.Data.(type) {

			case *slack.MessageEvent:
        fmt.Printf("Message from %v: %v\n", ev.User, ev.Text)

        info := rtm.GetInfo()

        if ev.User != info.User.ID && atMentioned(ev.Text, info.User.ID) {
          msg := fmt.Sprintf("Hi %s, what's up buddy?", ev.User)
					rtm.SendMessage(rtm.NewOutgoingMessage(msg, ev.Channel))
				}

			case *slack.RTMError:
				fmt.Printf("Error: %s\n", ev.Error())

			case *slack.InvalidAuthEvent:
        os.Stderr.WriteString("Error: Invalid credentials. Check the value of SLACK_TOKEN.")
        os.Exit(1)

			default:
        // Ignore other events..
  			// fmt.Printf("Unexpected: %v\n", msg.Data)
			}
		}
	}
}
