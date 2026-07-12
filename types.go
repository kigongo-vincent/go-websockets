package websockets

import "time"

type MessageType string

const (
	MessageTypeChat      MessageType = "chat"
	MessageTypeTyping    MessageType = "typing"
	MessageTypeDelivered MessageType = "delivered"
	MessageTypeRead      MessageType = "read"
	MessageTypePresence  MessageType = "presence"
	MessageTypeError     MessageType = "error"
)

type IncomingMessage[T any] struct {
	Type   MessageType `json:"type"`
	To     uint        `json:"to"`
	RoomID uint        `json:"room_id"`
	Data   T           `json:"data"`
}

type OutgoingMessage[T any] struct {
	ID        string      `json:"id"`
	Type      MessageType `json:"type"`
	From      uint        `json:"from"`
	To        uint        `json:"to"`
	RoomID    uint        `json:"room_id"`
	Data      T           `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
}

type Message[T any] struct {
	ID        string
	Type      MessageType
	From      uint
	To        uint
	RoomID    uint
	Data      T
	Timestamp time.Time
}

type PresenceStatus string

const (
	PresenceOnline  PresenceStatus = "online"
	PresenceOffline PresenceStatus = "offline"
)

type PresenceEvent struct {
	UserID    uint
	Status    PresenceStatus
	Timestamp time.Time
}

type ErrorMessage struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
