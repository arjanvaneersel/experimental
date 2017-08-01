package main

import (
	"github.com/nlopes/slack"
	"fmt"
	"strconv"
)

func (pb *phoneBot) ErrorAsDM(ev *slack.MessageEvent, msg string) {
	pb.rtm.PostMessage(ev.User, msg, slack.PostMessageParameters{AsUser: true, LinkNames: 1})
}

func (pb *phoneBot) ErrorAsMsg(ev *slack.MessageEvent, msg string) {
	pb.rtm.SendMessage(pb.rtm.NewOutgoingMessage(msg, ev.Channel))
}

func (pb *phoneBot) CallCommand(ev *slack.MessageEvent, args []string) {
	errMsg := "Command error. Please type *help* for more information."
	if len(args) < 2 {
		// Not enough of arguments given
		pb.rtm.PostMessage(ev.User,errMsg, slack.PostMessageParameters{AsUser: true, LinkNames: 1})
		return
	}

	pargs := pb.ProcessArgs(args)
	switch(pargs[0]) {
		case "accept":
			if len(args) < 2 {
				pb.rtm.PostMessage(ev.User,errMsg, slack.PostMessageParameters{AsUser: true, LinkNames: 1})
				break
			}
			id, err := strconv.Atoi(pargs[1])
			if err != nil {
				pb.rtm.PostMessage(ev.User,fmt.Sprintf("%s Is an invalid ID", pargs[1]), slack.PostMessageParameters{AsUser: true, LinkNames: 1})
				break
			}
			acceptCall(pb, id , ev.Msg.User)

		default:
			// Invalid command
			pb.rtm.PostMessage(ev.User, errMsg, slack.PostMessageParameters{AsUser: true, LinkNames: 1})
	}
}

func (pb *phoneBot) ChannelCommand(ev *slack.MessageEvent, args []string) {
	errMsg := "[channel] Command error. Please type *help* for more information."
	if len(args) < 2 {
		pb.rtm.PostMessage(ev.User,errMsg, slack.PostMessageParameters{AsUser: true, LinkNames: 1})
		return
	}

	pargs := pb.ProcessArgs(args)
	switch(pargs[0]) {
	case "set":
		if len(args) >= 2 {
			if pargs[1] == "broadcast" {
				if pb.channelID != "" {
					pb.rtm.PostMessage(ev.User, "A broadcasting channel has already been set. Use *channel unset broadcast* to unset.", slack.PostMessageParameters{AsUser: true, LinkNames: 1})
					break
				}
				if pb.IsDirectMessage(ev) {
					if len(args) < 3 {
						pb.rtm.PostMessage(ev.User, "When using channel set broadcast as a direct message, you need to provide the channel ID.", slack.PostMessageParameters{AsUser: true, LinkNames: 1})
						break
					}
					 _, err := pb.rtm.GetChannelInfo(args[2])
					//channels, err := pb.rtm.GetChannels(true)
					if err != nil {
						pb.rtm.PostMessage(ev.User, err.Error(), slack.PostMessageParameters{AsUser: true, LinkNames: 1})
						break
					}
					pb.channelID = args[2]
				} else {
					pb.channelID = ev.Channel
				}
				pb.rtm.SendMessage(pb.rtm.NewOutgoingMessage(":mega: This channel has been set as broadcasting channel.", pb.channelID))
			}
		} else {
			// Insufficient arguments
			pb.rtm.PostMessage(ev.User, errMsg, slack.PostMessageParameters{AsUser: true, LinkNames: 1})
		}
	case "unset":
		if len(args) >= 2 {
			if pargs[1] == "broadcast" {
				ch := pb.channelID
				pb.channelID = ""
				pb.rtm.SendMessage(pb.rtm.NewOutgoingMessage(":mega: This channel has been unset as broadcasting channel.", ch))
			}
		} else {
			// Insufficient arguments
			pb.rtm.PostMessage(ev.User, errMsg, slack.PostMessageParameters{AsUser: true, LinkNames: 1})
		}
	case "id":
		chInfo, err := pb.rtm.GetChannelInfo(ev.Channel)
		if err != nil {
			pb.rtm.PostMessage(ev.User, err.Error(), slack.PostMessageParameters{AsUser: true, LinkNames: 1})
		}
		pb.rtm.PostMessage(ev.User, fmt.Sprintf("The ID of channel *#%s* is *%s*", chInfo.Name, ev.Channel), slack.PostMessageParameters{AsUser: true, LinkNames: 1})
	default:
		// Invalid command
		pb.rtm.PostMessage(ev.User, errMsg, slack.PostMessageParameters{AsUser: true, LinkNames: 1})
	}
}


func (pb *phoneBot) HelloCommand(ev *slack.MessageEvent, args []string) {
	sender, _ := pb.SenderName(&ev.Msg)
	pb.rtm.SendMessage(pb.rtm.NewOutgoingMessage(fmt.Sprintf("Oh hello there, %s", sender), ev.Channel))
}

func (pb *phoneBot) KufteCommand(ev *slack.MessageEvent, args []string) {
	sender, _ := pb.SenderName(&ev.Msg)
	if args == nil {
		pb.rtm.SendMessage(pb.rtm.NewOutgoingMessage(fmt.Sprintf("@%s I like kufte and kepeke too!", sender), ev.Channel))
	} else {
		pb.rtm.SendMessage(pb.rtm.NewOutgoingMessage(fmt.Sprintf("@%s I like kufte and kepeke too! But what is a %s?", sender, args[0]), ev.Channel))
	}
}


