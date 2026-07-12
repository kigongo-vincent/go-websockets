package websockets

import (
	"context"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
)

type service[T any] struct {
	hub Hub[T]

	resolver UserResolver

	store MessageStore[T]

	presence PresenceStore

	hooks Hooks[T]
}

func NewService[T any](
	deps Dependencies[T],
) Service[T] {

	return &service[T]{
		hub: deps.Hub,

		resolver: deps.GetUserID,

		store: deps.Store,

		presence: deps.Presence,

		hooks: deps.Hooks,
	}
}

func (s *service[T]) Connect(
	conn *websocket.Conn,
) {

	userID, err := s.userID(conn)

	if err != nil {
		conn.Close()
		return
	}

	client := NewClient(
		uuid.NewString(),
		userID,
		conn,
	)

	s.register(client)

	go client.WritePump()

	s.read(client)
}

func (s *service[T]) userID(
	conn *websocket.Conn,
) (uint, error) {

	return s.resolver(conn)
}

func (s *service[T]) register(
	client *Client,
) {

	s.hub.Register(client)

	if s.hooks.OnConnect != nil {
		s.hooks.OnConnect(client)
	}
}

func (s *service[T]) read(
	client *Client,
) {
	defer s.disconnect(client)

	for {
		if err := s.receive(client); err != nil {
			return
		}
	}
}

func (s *service[T]) receive(
	client *Client,
) error {

	var input IncomingMessage[T]

	if err := client.Conn.ReadJSON(&input); err != nil {
		return err
	}

	msg := s.message(
		client.UserID,
		input,
	)

	s.send(msg)

	return nil
}

func (s *service[T]) message(
	userID uint,
	input IncomingMessage[T],
) Message[T] {

	return Message[T]{
		ID: uuid.NewString(),

		Type: input.Type,

		From: userID,

		To: input.To,

		RoomID: input.RoomID,

		Data: input.Data,

		Timestamp: time.Now(),
	}
}

func (s *service[T]) send(
	msg Message[T],
) {

	if s.before(msg) != nil {
		return
	}

	s.persist(msg)

	s.hub.Send(msg)

	s.after(msg)
}

func (s *service[T]) before(
	msg Message[T],
) error {

	if s.hooks.BeforeSend == nil {
		return nil
	}

	return s.hooks.BeforeSend(&msg)
}

func (s *service[T]) after(
	msg Message[T],
) {

	if s.hooks.AfterSend != nil {
		s.hooks.AfterSend(&msg)
	}
}

func (s *service[T]) persist(
	msg Message[T],
) {

	if s.store == nil {
		return
	}

	s.store.Create(
		context.Background(),
		msg,
	)
}

func (s *service[T]) disconnect(
	client *Client,
) {

	s.hub.Unregister(client)

	if s.hooks.OnDisconnect != nil {
		s.hooks.OnDisconnect(client)
	}
}

func (s *service[T]) Send(
	ctx context.Context,
	msg Message[T],
) error {

	return s.hub.Send(msg)
}

func (s *service[T]) IsOnline(
	userID uint,
) bool {

	return s.hub.Connected(userID)
}

func (s *service[T]) Clients(
	userID uint,
) []*Client {

	return s.hub.UserClients(userID)
}
