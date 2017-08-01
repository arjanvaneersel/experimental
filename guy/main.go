package main

import (
	"os"
	"fmt"
	"log"
	"strings"
)

func main() {
	token := os.Getenv("GUY_TOKEN")
	if token == "" {
		fmt.Fprintf(os.Stderr, "SB_TOKEN is not set in environment\n")
		os.Exit(1)
	}
	ws, id, err := slackConnect(token)
	if err != nil {
		log.Fatal(err)
	}

	for {
		msg, err := getMessage(ws)
		if err != nil {
			log.Println(msg)
			log.Fatal(err)
		}

		if msg.Type == "message" && msg.User != id {
			//parts := strings.Fields(msg.Text)
			go func(msg Message) {
				t := analyzeMessage(msg.Text)
				if t != "" {
					msg.Text = t
					postMessage(ws, msg)
				}
			}(msg)
		}
	}
}

type Trigger struct {
	RespondTo []string
	Response string
}

func (t *Trigger) Contains(msg string) bool {
	for _, syn := range t.RespondTo {
		found := strings.Contains(msg, syn)
		if found {
			return true
		}
	}
	return false
}

var Triggers []Trigger = []Trigger{
	Trigger{
		RespondTo: []string{"populism", "populist"},
		Response:"There is only one good populism: My populism!",
	},
	Trigger{
		RespondTo: []string{"losing elections", "lost elections", "losing the elections", "lost the elections", "losing an election", "lost an election"},
		Response: "If we or our allies lost the elections then the elections must have been hacked by Russia!",
	},
	Trigger{
		RespondTo: []string{"russia", "russian"},
		Response: "We need to send our troops to the borders of those pilmeni tossers!",
	},
	Trigger{
		RespondTo: []string{"putin", "vladimir putin"},
		Response: "The only Putin I like, is your money *put in* my hand to waste it!",
	},
	Trigger{
		RespondTo: []string{"war", "warmonger", "warmongers", "warmongering"},
		Response: "So you think you know about war? Think again, infidel! Nobody is better at provoking Russia than me! I mean, we need a European army to protect us from Russian aggression!",
	},
	Trigger{
		RespondTo: []string{"europe"},
		Response: "Who needs a continent, if you can have a super state!",
	},
	Trigger{
		RespondTo: []string{"eu", "european union", "union"},
		Response: "People don't need this European Union. They need emperor Guy's XXXXXXL federal European Empire!",
	},
	Trigger{
		RespondTo: []string{"guy", "verhofstadt"},
		Response: "Don't make me shy, make me emperor!",
	},
}


func analyzeMessage(msg string) string {
	str := strings.ToLower(msg)
	for _, trigger := range Triggers {
		if trigger.Contains(str) {
			return trigger.Response
		}
	}
	return "I don't know what you're talking about, but I am sure that Russia is behind it."
}

