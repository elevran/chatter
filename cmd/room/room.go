package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/elevran/chatter/pkg/gameon"
)

var exits = map[string]string{
	"N": "An old wooden door with a large arrow carved on its center",
	"S": "A heavy metal door with signs of rust",
	"W": "A gray, plain looking door",
	"E": "A door surrounded by a mysterious glow along it edges",
}

type room struct{}

func newChatRoom() *room {
	return &room{}
}

func (r *room) hello(resp http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		resp.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var hello gameon.Hello
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&hello)

	if err != nil || hello.UserID == "" {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	location := gameon.Message{
		Direction: "player",
		Recipient: hello.UserID,
		Payload: jsonMarshal(gameon.Location{
			Type:        "location",
			Name:        "Chatter",
			FullName:    "A chat room",
			Description: "a darkly lit room, there are people here, some are walking around, some are standing in groups",
			Exits:       exits,
			Commands:    map[string]string{},
			Inventory:   []string{},
		}),
	}

	welcome := gameon.Message{
		Direction: "player",
		Recipient: "*",
		Payload: jsonMarshal(gameon.Event{
			Type: "event",
			Content: map[string]string{
				hello.UserID: "Welcome!",
				"*":          fmt.Sprintf("%s has just entered the room", hello.Username),
			},
		}),
	}

	writeResponseMessages(resp, location, welcome)
}

func (r *room) goodbye(resp http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		resp.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var goodbye gameon.Goodbye
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&goodbye)

	if err != nil || goodbye.UserID == "" {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	farewell := gameon.Message{
		Direction: "player",
		Recipient: "*",
		Payload: jsonMarshal(gameon.Event{
			Type: "event",
			Content: map[string]string{
				goodbye.UserID: "Farewell!",
				"*":            fmt.Sprintf("%s has left the room", goodbye.Username),
			},
		}),
	}

	writeResponseMessages(resp, farewell)
}

func (r *room) room(resp http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		resp.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var command gameon.RoomCommand
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&command)
	if err != nil || command.UserID == "" || command.Content == "" {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	if strings.HasPrefix(command.Content, "/") {
		// slash command
		r.handleSlash(command, resp)
	} else {
		// chat command
		r.handleChat(command, resp)
	}
}

func (r *room) handleSlash(command gameon.RoomCommand, resp http.ResponseWriter) {
	words := strings.Fields(command.Content)
	commandName := strings.ToLower(words[0])

	replyWithEvent := func(content string) {
		event := gameon.Message{
			Direction: "player",
			Recipient: command.UserID,
			Payload: jsonMarshal(gameon.Event{
				Type: "event",
				Content: map[string]string{
					command.UserID: content,
				},
			}),
		}
		writeResponseMessages(resp, event)
	}

	switch commandName {
	case "/go":
		if len(words) < 2 {
			replyWithEvent("Go where?")
			break
		}

		exitID := strings.ToUpper(words[2])
		if _, ok := exits[exitID]; !ok {
			replyWithEvent("You probably don't wanna go there...")
			break
		}

		location := gameon.Message{
			Direction: "playerLocation",
			Recipient: command.UserID,
			Payload: jsonMarshal(gameon.PlayerLocation{
				Type:    "exit",
				Content: "You frantically run towards the exit",
				ExitID:  exitID,
			}),
		}
		writeResponseMessages(resp, location)

	case "/examine":
		replyWithEvent("Shouldn't you be mingling?")
	case "/inventory":
		replyWithEvent("There is nothing here")
	case "/look":
		replyWithEvent("It's just a room")
	default:
		replyWithEvent(fmt.Sprintf("Don't know how to %s", commandName[1:]))
	}
}

func (r *room) handleChat(command gameon.RoomCommand, resp http.ResponseWriter) {
	chat := gameon.Message{
		Direction: "player",
		Recipient: "*",
		Payload: jsonMarshal(gameon.Chat{
			Type:     "chat",
			Username: command.Username,
			Content:  command.Content,
		}),
	}

	writeResponseMessages(resp, chat)
}

func writeResponseMessages(resp http.ResponseWriter, messages ...gameon.Message) {
	bytes := jsonMarshal(gameon.MessageCollection{
		Messages: messages,
	})

	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(http.StatusOK)
	resp.Write(bytes)
}

func jsonMarshal(obj interface{}) []byte {
	bytes, _ := json.Marshal(obj)
	return bytes
}
