package websockets

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

type Manager[T any] struct {
	hub Hub[T]

	service Service[T]
}

func NewManager[T any](
	config Config[T],
) *Manager[T] {

	hub := NewHub[T]()

	service := NewService(
		Dependencies[T]{

			Hub: hub,

			GetUserID: config.GetUserID,

			Hooks: config.Hooks,
		},
	)

	return &Manager[T]{

		hub: hub,

		service: service,
	}
}

func (m *Manager[T]) Handler() fiber.Handler {

	return websocket.New(
		m.service.Connect,
	)
}

func (m *Manager[T]) SendUser(
	userID uint,
	data T,
) error {

	return m.hub.SendUser(
		userID,
		data,
	)
}
