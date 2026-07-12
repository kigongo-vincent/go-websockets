package websockets

import (
	"github.com/gofiber/contrib/websocket"
)

type Handler[T any] struct {
	Service Service[T]
}

func NewHandler[T any](
	service Service[T],
) *Handler[T] {
	return &Handler[T]{
		Service: service,
	}
}

func (h *Handler[T]) Connect(
	conn *websocket.Conn,
) {
	h.Service.Connect(conn)
}
