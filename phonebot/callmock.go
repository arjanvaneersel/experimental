package main

import (
	"math/rand"
	"log"
	"fmt"
	"github.com/nlopes/slack"
	"time"
)

type connections map[string]*call

func (cn connections) Connected(o string) bool {
	_, ok := cn[o]
	if !ok {
		return false
	}

	return cn[o] != nil
}

var operators = make(connections)

type call struct {
	ID int
	Name string
	Subject string
	Chan string
}

var calls []*call = []*call{
	&call{ID:0, Name:"Snoepie Piklov", Subject:"Buy bones"},
	&call{ID:1, Name:"Djoslin Piklova", Subject:"Play with the hippopotamus"},
	&call{ID:2, Name:"Bobcho de Brave", Subject:"Eat some granules"},
	&call{ID:3,Name:"Petko Magareto", Subject:"Go to the grass field"},
}

func mockCall(pb *phoneBot) {
	rand.Seed(time.Now().UTC().UnixNano())
	log.Println("*************************** Time to mock a call")

	if err := pb.PublishCall(calls[rand.Intn(len(calls))]); err != nil {
		log.Printf("*************************** %s", err.Error())
	}
}

func acceptCall(pb *phoneBot, id int, o string) {
	call := calls[id]
	if call.Chan != "" {
		pb.rtm.PostMessage(o, fmt.Sprintf("Call %d has already been accepted by an operator", call.ID), slack.PostMessageParameters{AsUser: true, LinkNames: 1})
		return
	}
	call.Chan = o
	operators[o] = call
	u, err := pb.rtm.GetUserInfo(o)
	if err != nil {
		pb.rtm.PostMessage(o, err.Error(), slack.PostMessageParameters{AsUser: true, LinkNames: 1})
		return
	}
	pb.rtm.SendMessage(pb.rtm.NewOutgoingMessage(fmt.Sprintf("Call %d of %s has been accepted by %s.", call.ID, call.Name, u.Name), pb.channelID))
	pb.rtm.PostMessage(o, fmt.Sprintf("You have accepted call %d of %s.\nThe subject is: %s", call.ID, call.Name, call.Subject), slack.PostMessageParameters{AsUser: true, LinkNames: 1})
}

func hangupCall(pb *phoneBot, o string) {
	call, ok := operators[o]
	if !ok {
		// Non existing operator
	}
	operators[o] = nil
	pb.rtm.PostMessage(call.Chan, fmt.Sprintf("Call %d of %s has been terminated.", call.ID, call.Name), slack.PostMessageParameters{AsUser: true, LinkNames: 1})
	call.Chan = ""
}
