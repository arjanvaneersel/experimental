package main

import (
	"fmt"
	"os"
	"github.com/nlopes/slack"
	"time"
)

func main() {
	pb := NewPhoneBot(false)
	if err := pb.Connect(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	pb.InvalidMentionCommand = pb.ErrorAsDM
	pb.InvalidDirectCommand = pb.ErrorAsDM

	pb.AddCommand(BotGlobalCommand, "hello", "Will greet you back", pb.HelloCommand)
	pb.AddCommand(BotDirectCommand, "hi", "Will greet you back", pb.HelloCommand)
	pb.AddCommand(BotMentionCommand, "kufte", "Talks about food. Takes 1 extra argument", pb.KufteCommand)

	pb.AddCommand(BotMentionCommand, "channel", "", pb.ChannelCommand)
	pb.AddCommand(BotDirectCommand, "channel", "", pb.ChannelCommand)

	pb.AddCommand(BotMentionCommand, "call", "", pb.CallCommand)
	pb.AddCommand(BotDirectCommand, "call", "", pb.CallCommand)

	pb.ListenerHandleFunc(func(pb *phoneBot){
		for msg := range pb.rtm.IncomingEvents {
			fmt.Print("Event Received: ")
			switch ev := msg.Data.(type) {
			case *slack.HelloEvent:
			// Ignore hello

			case *slack.ConnectedEvent:
				//fmt.Println("Infos:", ev.Info)
				//fmt.Println("Connection counter:", ev.ConnectionCount)
				pb.ProcessConnectedEvent(ev)

			//rtm.SendMessage(rtm.NewOutgoingMessage("Hello world", "#phonebot-development"))

			case *slack.MessageEvent:
				fmt.Printf("Message: %v\n", ev)
				pb.ProcessMessageEvent(ev)

			case *slack.PresenceChangeEvent:
				fmt.Printf("Presence Change: %v\n", ev)

			case *slack.LatencyReport:
				fmt.Printf("Current latency: %v\n", ev.Value)

			case *slack.RTMError:
				fmt.Printf("Error: %s\n", ev.Error())

			case *slack.InvalidAuthEvent:
				fmt.Printf("Invalid credentials")
				return

			default:
			// We don't care about other events
			}
		}
	})

	time.AfterFunc(30 * time.Second, func(){
		mockCall(pb)
	})

	pb.Listen()
}