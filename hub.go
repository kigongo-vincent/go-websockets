package websockets

import (
	"errors"
	"sync"
)

type hub[T any] struct {
	mu sync.RWMutex

	clients map[uint]map[string]*Client
}

func NewHub[T any]() Hub[T] {
	return &hub[T]{
		clients: make(map[uint]map[string]*Client),
	}
}

func (h *hub[T]) Register(
	client *Client,
) error {

	h.mu.Lock()
	defer h.mu.Unlock()

	h.initUser(client.UserID)

	h.clients[client.UserID][client.ID] = client

	return nil
}

func (h *hub[T]) initUser(
	userID uint,
) {

	if h.clients[userID] == nil {
		h.clients[userID] = make(map[string]*Client)
	}
}

func (h *hub[T]) Unregister(
	client *Client,
) error {

	h.mu.Lock()
	defer h.mu.Unlock()

	users := h.clients[client.UserID]

	if users == nil {
		return nil
	}

	delete(
		users,
		client.ID,
	)

	return nil
}

func (h *hub[T]) DisconnectUser(
	userID uint,
) error {

	h.mu.Lock()
	defer h.mu.Unlock()

	for _, client := range h.clients[userID] {
		client.Conn.Close()
	}

	delete(
		h.clients,
		userID,
	)

	return nil
}

func (h *hub[T]) Send(
	msg Message[T],
) error {

	h.mu.RLock()
	defer h.mu.RUnlock()

	users := h.clients[msg.To]

	if len(users) == 0 {
		return errors.New("user offline")
	}

	payload := h.output(msg)

	for _, client := range users {
		client.Send <- payload
	}

	return nil
}

func (h *hub[T]) output(
	msg Message[T],
) OutgoingMessage[T] {

	return OutgoingMessage[T]{
		ID: msg.ID,

		Type: msg.Type,

		From: msg.From,

		To: msg.To,

		RoomID: msg.RoomID,

		Data: msg.Data,

		Timestamp: msg.Timestamp,
	}
}

func (h *hub[T]) SendUser(
	userID uint,
	data T,
) error {

	h.mu.RLock()
	defer h.mu.RUnlock()

	users := h.clients[userID]

	if len(users) == 0 {
		return ErrOffline
	}

	for _, client := range users {
		client.Send <- data
	}

	return nil
}

func (h *hub[T]) UserClients(
	userID uint,
) []*Client {

	h.mu.RLock()
	defer h.mu.RUnlock()

	result := make([]*Client, 0)

	for _, client := range h.clients[userID] {
		result = append(result, client)
	}

	return result
}

func (h *hub[T]) Connected(
	userID uint,
) bool {

	return len(
		h.UserClients(userID),
	) > 0
}

func (h *hub[T]) Count() int {

	h.mu.RLock()
	defer h.mu.RUnlock()

	count := 0

	for _, users := range h.clients {
		count += len(users)
	}

	return count
}
