package main

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"golang.org/x/net/websocket"
	"sync/atomic"
	"github.com/pkg/errors"
)

type responseRtmStart struct {
	OK bool `json:"ok"`
	Error string `json:"error"`
	URL string `json:"url"`
	Self responseSelf `json:"self"`
}

type responseSelf struct {
	ID string `json:"id"`
}

type Message struct {
	ID      uint64 `json:"id"`
	Type    string `json:"type"`
	Channel *string `json:"channel"`
	User string `json:"user"`
	Text    string `json:"text"`
}

func slackStart(token string) (string, string, error){
	url := "https://slack.com/api/rtm.start?token=" + token
	resp, err := http.Get(url)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	var respObj responseRtmStart
	err = json.Unmarshal(body, &respObj)
	if err != nil {
		return "", "", err
	}

	if !respObj.OK {
		return "", "", errors.New(respObj.Error)
	}

	return respObj.URL, respObj.Self.ID, nil
}

func slackConnect(token string) (*websocket.Conn, string, error) {
	wsurl, id, err := slackStart(token)
	if err != nil {
		return nil, "", err
	}

	ws, err := websocket.Dial(wsurl, "", "https://api.slack.com")
	if err != nil {
		return nil, "", err
	}
	return ws, id, nil
}

func getMessage(ws *websocket.Conn) (m Message, err error) {
	err = websocket.JSON.Receive(ws, &m)
	return
}

var counter uint64

func postMessage(ws *websocket.Conn, msg Message) error {
	msg.ID = atomic.AddUint64(&counter, 1)
	return websocket.JSON.Send(ws, msg)
}