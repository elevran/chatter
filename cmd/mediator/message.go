package main

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/elevran/chatter/pkg/gameon"
)

func formatMessage(msg *gameon.Message) ([]byte, error) {
	var buf bytes.Buffer

	buf.WriteString(msg.Direction)
	buf.WriteRune(',')

	if msg.Recipient != "" {
		buf.WriteString(msg.Recipient)
		buf.WriteRune(',')
	}

	buf.Write(msg.Payload)
	return buf.Bytes(), nil
}

func parseMessage(data []byte) (*gameon.Message, error) {
	parts := strings.SplitN(string(data), ",", 3)

	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid websocket message format: %s", string(data))
	}

	msg := new(gameon.Message)
	msg.Direction = parts[0]

	if strings.HasPrefix(parts[1], "{") {
		// case 1: <direction>,{...}
		msg.Payload = data[len(parts[0])+1:]
	} else {
		// case 2: <direction>,<recipient>,{...}
		msg.Recipient = parts[1]
		msg.Payload = data[len(parts[0])+len(parts[1])+2:]
	}

	return msg, nil
}
