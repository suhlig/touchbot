package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/nlopes/slack"
)

func atMentioned(message string, userID string) bool {
	prefix := fmt.Sprintf("<@%s>", userID)
	return strings.Contains(message, prefix)
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

			case *slack.ConnectedEvent:
				fmt.Printf("Connected as %s <@%s>\n", ev.Info.User.Name, ev.Info.User.ID)

			case *slack.MessageEvent:
				userInfo, err := rtm.GetUserInfo(ev.User)

				if err != nil {
					fmt.Printf("Error: %s\n", err)
				} else {
					fmt.Printf("Message from %v in #%v: %v\n", userInfo.Name, ev.Channel, ev.Text)

					if strings.HasPrefix(ev.Channel, "D") {
						fmt.Printf("DM from %s\n", userInfo.Name)
						rtm.SendMessage(rtm.NewOutgoingMessage("What's up?", ev.Channel))
					}

					info := rtm.GetInfo()

					if ev.User != info.User.ID && atMentioned(ev.Text, info.User.ID) {
						msg := fmt.Sprintf("What's up, <@%s>?", ev.User)
						rtm.SendMessage(rtm.NewOutgoingMessage(msg, ev.Channel))
					}
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
