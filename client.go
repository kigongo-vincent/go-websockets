package websockets

import (
	"time"

	"github.com/gofiber/contrib/websocket"
)

type Client struct {
	ID          string
	UserID      uint
	Conn        *websocket.Conn
	Send        chan any
	ConnectedAt time.Time
}

func NewClient(id string, userID uint, conn *websocket.Conn,
) *Client {
	return &Client{
		ID:          id,
		UserID:      userID,
		Conn:        conn,
		Send:        make(chan any, 64),
		ConnectedAt: time.Now(),
	}
}

func (c *Client) WritePump() {
	defer c.Conn.Close()
	for msg := range c.Send {
		if err := c.Conn.WriteJSON(msg); err != nil {
			return
		}
	}
}
