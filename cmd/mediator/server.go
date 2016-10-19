package main

import (
	"encoding/json"
	"net/http"

	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/elevran/chatter/pkg/gameon"
	"github.com/gorilla/websocket"
)

var (
	SupportedVersions = []int{1}
)

type Server struct {
	client   *Client
	roomID   string
	sessions *SessionManager
}

func newServer(config *Config) *Server {
	s := &Server{
		client:   newClient(config),
		roomID:   config.RoomID,
		sessions: NewSessionManager(),
	}

	return s
}

func (s *Server) handleHTTP(w http.ResponseWriter, r *http.Request) {
	logrus.Debugf("Incoming HTTP request from %s", r.RemoteAddr)

	var upgrader websocket.Upgrader
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logrus.WithError(err).Errorf("Error upgrading HTTP to websocket connection")
		return
	}

	logrus.Debugf("Websocket connection established with %s", conn.RemoteAddr().String())
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
	// The loop runs forever, and is terminated only when an error occurs.
	// In such a case, attempt to close the session (in case not closed already).
	defer session.Close()

	for {
		_, bytes, err := session.Conn.ReadMessage()
		if err != nil {
			logrus.WithError(err).Errorf("Error reading websocket message")
			return
		}

		msg, err := parseMessage(bytes)
		if err != nil {
			logrus.WithError(err).Errorf("Error parsing websocket message")
			return
		}

		logrus.WithFields(messageToFields(msg)).Debugf("Websocket message received")

		// Validate the message recipient is our own room ID
		if s.roomID != "" && msg.Recipient != s.roomID {
			logrus.WithError(fmt.Errorf("recipient (%s) doesn't match expected room id (%s)", msg.Recipient, s.roomID)).
				Errorf("Invalid message received")
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
			logrus.WithError(fmt.Errorf("unrecognized message direction: %s", msg.Direction)).
				Errorf("Invalid message received")
			return
		}

		err = json.Unmarshal(msg.Payload, payload)
		if err != nil {
			logrus.WithError(err).Errorf("Error unmarshaling message payload")
			return
		}

		switch payload := payload.(type) {
		case *gameon.Hello:
			s.handleHello(payload, session)
		case *gameon.Goodbye:
			s.handleGoodbye(payload, session)
		case *gameon.Command:
			s.handleCommand(payload, session)
		default:
			logrus.WithError(fmt.Errorf("unrecognized payload type: %T", payload))
		}
	}
}

func (s *Server) ack(session *Session) {
	logrus.Debugf("Sending ack for websocket connection with remote address %s", session.Conn.RemoteAddr().String())

	ack := gameon.Ack{
		Version: SupportedVersions,
	}
	ackBytes, _ := json.Marshal(ack)

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
		logrus.WithError(err).Errorf("Error executing 'hello' with room service")
		return
	}

	s.handleResponse(resp)
}

func (s *Server) handleGoodbye(goodbye *gameon.Goodbye, session *Session) {
	defer session.Close()

	resp, err := s.client.doGoodbye(goodbye)
	if err != nil {
		logrus.WithError(err).Errorf("Error executing 'goodbye' with room service")
		return
	}

	s.handleResponse(resp)
}

func (s *Server) handleCommand(command *gameon.Command, session *Session) {
	resp, err := s.client.doCommand(command)
	if err != nil {
		logrus.WithError(err).Errorf("Error executing command with room service")
		return
	}

	s.handleResponse(resp)
}

func (s *Server) handleResponse(resp *gameon.MessageCollection) {
	switch len(resp.Messages) {
	case 0:
		logrus.Debugf("Response contains no messages")
	case 1:
		logrus.Debugf("Dispatching 1 response message")
	default:
		logrus.Debugf("Dispatching %d response message", len(resp.Messages))
	}

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
	logrus.WithFields(messageToFields(msg)).Debugf("Sending message")

	bytes, err := formatMessage(msg)
	if err != nil {
		logrus.WithError(err).Errorf("Error formatting message")
		return
	}

	err = session.Conn.WriteMessage(websocket.TextMessage, bytes)
	if err != nil {
		logrus.WithError(err).Errorf("Error sending message")
		session.Close()
	}
}

func (s *Server) broadcastMessage(msg *gameon.Message, sessions []*Session) {
	logrus.WithFields(messageToFields(msg)).Debugf("Broadcasting message")

	bytes, err := formatMessage(msg)
	if err != nil {
		logrus.WithError(err).Errorf("Error formatting message")
		return
	}

	for _, session := range sessions {
		err := session.Conn.WriteMessage(websocket.TextMessage, bytes)
		if err != nil {
			logrus.WithError(err).Errorf("Error broadcasting message")
			session.Close()
		}
	}

}
