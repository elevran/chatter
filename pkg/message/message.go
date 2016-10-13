package message

import "encoding/json"

// Message is a generic GameOn! message, holding a direction (player, room, ...),
// a recipient (playerID, roomID, *, ...), and a message-type specific payload.
type Message struct {
	Direction string          `json:"direction,omitempty"`
	Recipient string          `json:"recipient,omitempty"`
	Payload   json.RawMessage `json:"payload,omitempty"`
}

// Command is the payload provided for a [client --> mediator --> room] chat/command message.
type Command struct {
	Username string `json:"username,omitempty"`
	UserID   string `json:"userId,omitempty"`
	Content  string `json:"content,omitempty"`
}

// Hello is the payload provided for a [mediator --> room] hello message.
type Hello struct {
	Username string `json:"username,omitempty"`
	UserID   string `json:"userId,omitempty"`
	Version  int    `json:"version,omitempty"`
	Recovery bool   `json:"recovery,omitempty"`
}

// Goodbye is the payload provided for a [mediator --> room] goodbye message.
type Goodbye struct {
	Username string `json:"username,omitempty"`
	UserID   string `json:"userId,omitempty"`
}

// Ack is the payload provided for a [room --> mediator] ack message.
type Ack struct {
	Version []int `json:"version,omitempty"`
}

// Location is the payload provided for a [room --> mediator --> client] location message.
type Location struct {
	Type        string            `json:"type,omitempty"`
	Name        string            `json:"name,omitempty"`
	FullName    string            `json:"fullName,omitempty"`
	Description string            `json:"description,omitempty"`
	Exits       map[string]string `json:"exits,omitempty"`
	Commands    map[string]string `json:"commands,omitempty"`
	Inventory   []string          `json:"roomInventory,omitempty"`
}

// Chat is the payload provided for a [room --> mediator --> client] chat message.
type Chat struct {
	Content  Content
	Bookmark string `json:"bookmark,omitempty"`
}

// Event is the payload provided for a [room --> mediator --> client] event message.
type Event struct {
	Content  Content
	Bookmark string `json:"bookmark,omitempty"`
}

// Content part of a chat / event message
type Content struct {
	Broadcast string
	Private   map[string]string
}

type typedContent struct {
	Type     string  `json:"type,omitempty"`
	Content  Content `json:"content,omitempty"`
	Bookmark string  `json:"bookmark,omitempty"`
}

func (c *Chat) MarshalJSON() ([]byte, error) {
	t := &typedContent{
		Type:     "chat",
		Content:  c.Content,
		Bookmark: c.Bookmark,
	}

	return json.Marshal(t)
}

func (e *Event) MarshalJSON() ([]byte, error) {
	t := &typedContent{
		Type:     "event",
		Content:  e.Content,
		Bookmark: e.Bookmark,
	}
	return json.Marshal(t)
}

func (c *Content) MarshalJSON() ([]byte, error) {
	m := make(map[string]string, len(c.Private)+1)
	for user, msg := range c.Private {
		m[user] = msg
	}

	if c.Broadcast != "" {
		m["*"] = c.Broadcast
	}

	return json.Marshal(m)
}

func (c *Content) UnmarshalJSON(data []byte) error {
	m := make(map[string]string)
	err := json.Unmarshal(data, &m)
	if err != nil {
		return err
	}

	c.Private = make(map[string]string, len(m))
	for user, msg := range m {
		c.Private[user] = msg
	}

	delete(c.Private, "*")
	c.Broadcast = m["*"]

	return nil
}
