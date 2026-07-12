package websockets

import (
	"context"

	"github.com/gofiber/contrib/websocket"
)

type UserResolver func(
	*websocket.Conn,
) (uint, error)

type Hub[T any] interface {
	Register(
		*Client,
	) error

	Unregister(
		*Client,
	) error

	DisconnectUser(
		uint,
	) error

	Send(
		Message[T],
	) error

	SendUser(
		uint,
		T,
	) error

	UserClients(
		uint,
	) []*Client

	Connected(
		uint,
	) bool

	Count() int
}

type Hooks[T any] struct {
	OnConnect func(
		*Client,
	) error

	OnDisconnect func(
		*Client,
	)

	BeforeSend func(
		*Message[T],
	) error

	AfterSend func(
		*Message[T],
	)

	OnError func(
		error,
	)
}

type Dependencies[T any] struct {
	Hub Hub[T]

	GetUserID UserResolver

	Store MessageStore[T]

	Presence PresenceStore

	Hooks Hooks[T]
}

type Service[T any] interface {
	Connect(
		*websocket.Conn,
	)

	Send(
		context.Context,
		Message[T],
	) error

	IsOnline(
		uint,
	) bool

	Clients(
		uint,
	) []*Client
}

type MessageStore[T any] interface {
	Create(
		context.Context,
		Message[T],
	) error

	ListRoom(
		context.Context,
		uint,
		int,
		int,
	) ([]Message[T], error)
}

type PresenceStore interface {
	SetOnline(
		context.Context,
		uint,
	) error

	SetOffline(
		context.Context,
		uint,
	) error

	IsOnline(
		context.Context,
		uint,
	) (bool, error)
}
