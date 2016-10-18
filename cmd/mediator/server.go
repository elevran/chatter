package main

import (
	"encoding/json"
	"net/http"

	"github.com/elevran/chatter/pkg/gameon"
	"github.com/gorilla/websocket"
)

var (
	SupportedVersions = []int{1}
)

type Server struct {
	client   *Client
	sessions *SessionManager
}

func newServer() *Server {
	s := &Server{
		client:   newClient(),
		sessions: NewSessionManager(),
	}

	return s
}

func (s *Server) handleHTTP(w http.ResponseWriter, r *http.Request) {
	var upgrader websocket.Upgrader
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	s.handleWebsocket(conn)
}

func (s *Server) handleWebsocket(conn *websocket.Conn) {
	session := s.sessions.NewSession(conn)

	s.ack(session)
	go s.handleMessages(session)

	select {
	case <-session.Closed():
		conn.Close()
	}
}

func (s *Server) handleMessages(session *Session) {
	for {
		_, bytes, err := session.Conn.ReadMessage()
		if err != nil {
			session.Close()
			return
		}

		msg, err := parseMessage(bytes)
		if err != nil {
			session.Close()
			return
		}

		var payload interface{}
		switch msg.Direction {
		case "roomHello":
			payload = &gameon.Hello{}
		case "roomGoodbye":
			payload = &gameon.Goodbye{}
		case "room":
			payload = &gameon.Command{}
		default:
			session.Close()
			return
		}

		err = json.Unmarshal(msg.Payload, payload)
		if err != nil {
			session.Close()
			return
		}

		switch payload := payload.(type) {
		case *gameon.Hello:
			s.handleHello(payload, session)
		case *gameon.Goodbye:
			s.handleGoodbye(payload, session)
		case *gameon.Command:
			s.handleCommand(payload, session)
		}
	}
}

func (s *Server) ack(session *Session) {
	ack := gameon.Ack{
		Version: SupportedVersions,
	}
	ackBytes, err := json.Marshal(ack)
	if err != nil {
		session.Close()
		return
	}

	msg := &gameon.Message{
		Direction: "ack",
		Payload:   ackBytes,
	}

	s.sendMessage(msg, session)
}

func (s *Server) handleHello(hello *gameon.Hello, session *Session) {
	session.SetUserID(hello.UserID)

	resp, err := s.client.doHello(hello)
	if err != nil {
		return
	}

	s.handleResponse(resp)
}

func (s *Server) handleGoodbye(goodbye *gameon.Goodbye, session *Session) {
	resp, err := s.client.doGoodbye(goodbye)
	if err != nil {
		return
	}

	s.handleResponse(resp)

	session.Close()
}

func (s *Server) handleCommand(command *gameon.Command, session *Session) {
	resp, err := s.client.doCommand(command)
	if err != nil {
		return
	}

	s.handleResponse(resp)
}

func (s *Server) handleResponse(resp *gameon.MessageCollection) {
	for _, msg := range resp.Messages {
		if msg.Recipient == "*" {
			s.broadcastMessage(&msg, s.sessions.GetUserSessions())
		} else {
			session := s.sessions.GetUserSession(msg.Recipient)
			if session != nil {
				s.sendMessage(&msg, session)
			}
		}
	}
}

func (s *Server) sendMessage(msg *gameon.Message, session *Session) {
	bytes, err := formatMessage(msg)
	if err != nil {
		return
	}

	err = session.Conn.WriteMessage(websocket.TextMessage, bytes)
	if err != nil {
		session.Close()
	}
}

func (s *Server) broadcastMessage(msg *gameon.Message, sessions []*Session) {
	bytes, err := formatMessage(msg)
	if err != nil {
		return
	}

	for _, session := range sessions {
		err := session.Conn.WriteMessage(websocket.TextMessage, bytes)
		if err != nil {
			session.Close()
		}
	}

}
