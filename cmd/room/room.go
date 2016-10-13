//
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type room struct {
}

func newChatRoom() *room {
	return &room{}
}

func (r *room) hello(resp http.ResponseWriter, req *http.Request) {
	defer closeRequestBody(req)

	if req.Method != "POST" {
		resp.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var hello Hello
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&hello)

	if err != nil || hello.UserID == "" {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	location := Location{
		Type:        "location",
		Name:        "chat room",
		Description: "a darkly lit room, there are people here, some are walking around, some are standing in groups",
	}
	welcome := Event{
		Type:     "event",
		Bookmark: time.Now().UTC().String(),
		Content: map[string]string{
			hello.UserID: "welcome",
			"*":          fmt.Sprintf("%s(%s) just enetered the room", hello.UserID, hello.Username),
		},
	}
	array := []interface{}{location, welcome}
	b, err := json.Marshal(&array)
	raw := json.RawMessage(b)
	body, err := json.Marshal(&Response{
		Destination: "player",
		Recipient:   hello.UserID,
		Payload:     &raw,
	})
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(http.StatusOK)
	resp.Write(body)
}

func (r *room) goodbye(resp http.ResponseWriter, req *http.Request) {
	defer closeRequestBody(req)

	if req.Method != "POST" {
		resp.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var bye Goodbye
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&bye)
	if err != nil || bye.UserID == "" {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	body, err := json.Marshal(&Response{
		Destination: "player",
		Recipient:   bye.UserID,
		// Payload is array of JSON raw messages (location, event announcing user departure and sending
		// "farewell and thanks for all the fish" to user)
	})

	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(http.StatusOK)
	resp.Write(body)
}

func (r *room) message(resp http.ResponseWriter, req *http.Request) {
	defer closeRequestBody(req)

	if req.Method != "POST" {
		resp.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var message PlayerMessage
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&message)
	if err != nil || message.UserID == "" || message.Content == "" {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	if message.Content[0] == '/' { // room command
		r.handleCommand(message, resp)

	} else {
		r.handleChat(message, resp)
	}
}

func (r *room) handleCommand(message PlayerMessage, resp http.ResponseWriter) {
	words := strings.Fields(message.Content)
	command := strings.ToLower(words[0])

	var reply string

	switch command {
	case "/examine":
		reply = "Shouldn't you be mingling?"
	case "/go":
		reply = "Sorry to see you go..."
	case "/inventory":
		reply = "There is nothing here"
	case "/look":
		reply = "It's just a room"
	default:
		reply = fmt.Sprintf("Don't know how to %s", command[1:len(command)])
	}

	r.sendEvent(message.UserID, reply, resp)
}

func (r *room) handleChat(message PlayerMessage, resp http.ResponseWriter) {
	r.sendEvent("*", message.Content, resp)
}

func (r *room) sendEvent(recipient, message string, resp http.ResponseWriter) {
	event := Event{
		Type:     "event",
		Bookmark: time.Now().UTC().String(),
		Content: map[string]string{
			recipient: message,
		},
	}

	body, err := json.Marshal(&event)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}
	raw := json.RawMessage(body)

	body, err = json.Marshal(&Response{
		Destination: "player",
		Recipient:   recipient,
		Payload:     &raw,
	})

	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(http.StatusOK)
	resp.Write(body)
}

func closeRequestBody(r *http.Request) {
	io.Copy(ioutil.Discard, r.Body)
	r.Body.Close()
}
