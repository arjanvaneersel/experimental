package main

import (
	"log"
	"os"

	"github.com/nlopes/slack"
	"errors"
	"strings"
	"fmt"
)

type commandType uint

const (
	BotGlobalCommand commandType = 0
	BotMentionCommand commandType = 1
	BotDirectCommand commandType = 2
)

type botCommand struct {
	Description string
	Do func(*slack.MessageEvent, []string)
}

type botCommands map[string]*botCommand

type phoneBot struct {
	id string
	rtm *slack.RTM
	channelID string
	ListenHandler func(pb *phoneBot)
	GlobalCommands botCommands
	InvalidGlobalCommand func(*slack.MessageEvent, string)
	MentionCommands botCommands
	InvalidMentionCommand func(*slack.MessageEvent, string)
	DirectCommands botCommands
	InvalidDirectCommand func(*slack.MessageEvent, string)
	CaseSensitive bool

}

func (b *phoneBot) Connect() error {
	token := os.Getenv("PB_SLACK_TOKEN")
	if token == "" {
		return errors.New("PB_SLACK_TOKEN is not set in environment")
	}

	api := slack.New(token)
	logger := log.New(os.Stdout, "phone-bot: ", log.Lshortfile|log.LstdFlags)
	slack.SetLogger(logger)
	api.SetDebug(true)

	b.rtm = api.NewRTM()
	go b.rtm.ManageConnection()
	return nil
}

func (b *phoneBot) IsConnected() bool {
	return b.id != "" && b.rtm != nil
}

func (b *phoneBot) ListenerHandleFunc(f func(pb *phoneBot)) {
	b.ListenHandler = f
}

func (b *phoneBot) Listen() {
	b.ListenHandler(b)
}

func (b *phoneBot) ProcessArgs(a []string) (r []string) {
	r = make([]string, len(a))
	copy(r, a)
	for i, v := range a {
		if !b.CaseSensitive {
			r[i] = strings.ToLower(v)
		}
	}
	return
}

func (b *phoneBot) SentByMe(msg *slack.Msg) bool {
	return msg.User == b.id
}

func (b *phoneBot) SenderName(msg *slack.Msg) (string, error) {
	u, err := b.rtm.GetUserInfo(msg.User)
	if err != nil {
		return "", err
	}
	return u.Name, nil
}

func (b *phoneBot) IsDirectMessage(ev *slack.MessageEvent) bool {
	return strings.HasPrefix(ev.Channel, "D")
}

func (b *phoneBot) MentionsMe(msg *slack.Msg) bool {
	return strings.HasPrefix(msg.Text, "<@" + b.id + ">")
}

func (b *phoneBot) ProcessConnectedEvent(ev *slack.ConnectedEvent) {
	b.id = ev.Info.User.ID
}

func (b *phoneBot) ProcessMessageEvent(ev *slack.MessageEvent) {
	if !b.SentByMe(&ev.Msg) {
		parts := strings.Fields(ev.Msg.Text)
		if b.IsDirectMessage(ev) {
			if operators.Connected(ev.Msg.User) {
				if ev.Msg.Text == "//END" {
					hangupCall(b, ev.Msg.User)
				}
			} else {
				var c string = parts[0]
				if !b.CaseSensitive {
					c = strings.ToLower(parts[0])
				}
				cmd, exists := b.DirectCommands[c]
				if !exists {
					// Command doesn't esist
					if b.InvalidDirectCommand != nil {
						b.InvalidDirectCommand(ev, "I don't understand that. Please type *help* for more information about my commands.")
					}
				} else {
					if len(parts) > 1 {
						cmd.Do(ev, parts[1:])
					} else {
						cmd.Do(ev, nil)
					}
				}
			}
		} else if b.MentionsMe(&ev.Msg) {
			var c string = parts[1]
			if !b.CaseSensitive {
				c = strings.ToLower(parts[1])
			}
			cmd, exists := b.MentionCommands[c]
			if !exists {
				// Command doesn't esist
				if b.InvalidMentionCommand != nil {
					b.InvalidMentionCommand(ev, "I don't understand that. Please type *help* for more information about my commands.")
				}
			} else {
				if len(parts) > 2 {
					cmd.Do(ev, parts[2:])
				} else {
					cmd.Do(ev, nil)
				}
			}
		} else {
			var c string = parts[0]
			if !b.CaseSensitive {
				c = strings.ToLower(parts[0])
			}
			cmd, exists := b.GlobalCommands[c]
			if !exists {
				// Command doesn't esist
				if b.InvalidGlobalCommand != nil {
					sender, _ := b.SenderName(&ev.Msg)
					b.InvalidGlobalCommand(ev, fmt.Sprintf("%s: I don't understand that. Type *help* for more information about my commands", sender))
				}
			} else {
				if len(parts) > 1 {
					cmd.Do(ev, parts[1:])
				} else {
					cmd.Do(ev, nil)
				}
			}
		}
	}
}

func (b *phoneBot) AddCommand(t commandType, cmd string, desc string, f func(ev *slack.MessageEvent, args []string)) error {
	if b.GlobalCommands == nil {
		b.GlobalCommands = make(botCommands)
	}
	if b.MentionCommands == nil {
		b.MentionCommands = make(botCommands)
	}
	if b.DirectCommands == nil {
		b.DirectCommands = make(botCommands)
	}

	switch(t) {
	case BotGlobalCommand:
		b.GlobalCommands[cmd] = &botCommand{desc, f}
	case BotMentionCommand:
		b.MentionCommands[cmd] = &botCommand{desc, f}
	case BotDirectCommand:
		b.DirectCommands[cmd] = &botCommand{desc, f}
	default:
		return errors.New("Invalid command type")
	}
	return nil
}

func (b *phoneBot) HelpText() string {
	// Global commands
	var g string
	for k, v := range b.GlobalCommands {
		if g == "" {
			g = "*Global commands*\n"
		}
		g = fmt.Sprintf("%s```[%s]\n%s```\n", g, k, v.Description)
	}

	// Mention commands
	var m string
	for k, v := range b.MentionCommands {
		if m == "" {
			m = "*Mention me commands*\n"
		}
		m = fmt.Sprintf("%s```[%s]\n%s```\n", m, k, v.Description)
	}

	// Direct commands
	var d string
	for k, v := range b.DirectCommands {
		if d == "" {
			d = "*Direct message commands*\n"
		}
		d = fmt.Sprintf("%s```[%s]\n%s```\n", d, k, v.Description)
	}

	var cs string
	if b.CaseSensitive {
		cs = "\n\r:exclamation: Commands are case sensitive"
	} else {
		cs = "\n\r:exclamation: Commands are not case sensitive"
	}

	return fmt.Sprintf("Phonebot version 0.1.\n\r*Listing of all commands:*\n\r%s\n%s\n%s\n%s", g, m, d, cs)
}

func (b *phoneBot) HelpMentionCommand(ev *slack.MessageEvent, args []string) {
	b.rtm.PostMessage(ev.User, b.HelpText(), slack.PostMessageParameters{AsUser: true, LinkNames: 1})
}

func (b *phoneBot) HelpDirectCommand(ev *slack.MessageEvent, args []string) {
	b.rtm.SendMessage(b.rtm.NewOutgoingMessage(b.HelpText(), ev.Channel))
}

func (b *phoneBot) PublishCall(c *call) error {
	if b.channelID == "" {
		return errors.New("Can't publish call. No broadcasting channel has been set.")
	}

	ch, err  := b.rtm.GetChannelInfo(b.channelID)
	if err != nil {
		return err
	}
	b.rtm.SendMessage(b.rtm.NewOutgoingMessage(fmt.Sprintf(":phone: Call from *%s* about *%s*. Use ID *%d* to accept the call", c.Name, c.Subject, c.ID), ch.ID))
	return nil
}

func NewPhoneBot(cs bool) *phoneBot {
	pb := &phoneBot{CaseSensitive: cs}
	pb.AddCommand(BotDirectCommand, "help", "Prints the bot's commands", pb.HelpDirectCommand)
	pb.AddCommand(BotMentionCommand, "help", "Prints the bot's commands", pb.HelpMentionCommand)
	return pb
}
