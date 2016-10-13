//
package main

import (
	"encoding/json"
)

type UserInfo struct {
	Username string `json:"username,omitempty"`
	UserID   string `json:"userId,omitempty"`
}

type Hello struct {
	UserInfo
	Version  int  `json:"version,omitempty"`
	Recovery bool `json:"recovery,omitempty"`
}

type Goodbye struct {
	UserInfo
}

type PlayerMessage struct {
	UserInfo
	Content string `json:"content,omitempty`
}

type Response struct {
	Destination string           `json:"destination,omitempty"`
	Recipient   string           `json:"recipient,omitempty"`
	Payload     *json.RawMessage `json:"payload,omitempty"`
}

type Location struct {
	Type        string `json:"type,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	// FullName string
	// exits is array of "NSEW": "description"
	// commands is array of "/command": "description"
	// inventory is array of string
}

type Event struct {
	Type     string            `json:"type,omitempty"`
	Bookmark string            `json:"bookmark,omitempty"`
	Content  map[string]string `json:"content,omitempty"`
}
