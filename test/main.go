package main

import (
	websocket "github.com/kigongo-vincent/go-websockets"
)

func main() {
	socket := websocket.NewManager[ChatPayload](
		websocket.Config[ChatPayload]{
			GetUserID: func(
				conn *websocket.Conn,
			) (uint, error) {
				return 1, nil
			},
		},
	)

	app.Get(
		"/ws",
		socket.Handler(),
	)
}
