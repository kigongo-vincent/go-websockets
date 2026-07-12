package websockets

import (
	"github.com/gofiber/contrib/websocket"
)

type Config[T any] struct {
	GetUserID func(
		conn *websocket.Conn,
	) (uint, error)

	Hooks Hooks[T]
}
