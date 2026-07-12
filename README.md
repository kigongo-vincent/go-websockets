# Go Fiber WebSocket Manager

A generic, production-focused WebSocket management library for **Go Fiber** applications.

`go-websockets` provides a reusable WebSocket infrastructure that handles:

- Client connection management
- User connection tracking
- Multiple connections per user
- Real-time message delivery
- Connection lifecycle events
- User identification
- Hook-based customization

The library is designed to be embedded into any Go Fiber application that requires real-time communication.

It does **not** enforce authentication, persistence, or business logic. Those remain the responsibility of the consuming application.

---

# Table of Contents

- [Installation](#installation)
- [Architecture](#architecture)
- [Core Concepts](#core-concepts)
- [Quick Start](#quick-start)
- [Configuration](#configuration)
- [Manager API](#manager-api)
- [Message System](#message-system)
- [Client System](#client-system)
- [Hooks](#hooks)
- [Hub API](#hub-api)
- [Connection Lifecycle](#connection-lifecycle)
- [Authentication Integration](#authentication-integration)
- [Usage Patterns](#usage-patterns)
- [Rules and Design Principles](#rules-and-design-principles)
- [Scaling](#scaling)
- [Complete Example](#complete-example)

---

# Installation

Install using Go modules:

```bash
go get github.com/kigongo-vincent/go-websockets
```

Import:

```go
import ws "github.com/kigongo-vincent/go-websockets"
```

---

# Requirements

| Dependency | Requirement     |
| ---------- | --------------- |
| Go         | 1.24+           |
| Fiber      | v2              |
| WebSocket  | Fiber WebSocket |

---

# Architecture

The library follows a layered architecture.

```
                  Your Application

                         |
                         |
                         v

                    Manager[T]

                         |
        +----------------+----------------+

        |                                 |

        v                                 v

    Service[T]                         Hub[T]

        |                                 |

        |                                 |

        v                                 v

 WebSocket Lifecycle              Connected Clients


                         |

                         v

                      Client
```

---

# Core Concepts

## Manager

The Manager is the main entry point.

It creates and exposes:

- WebSocket handlers
- Message sending capabilities
- Connection management

Applications normally interact only with the Manager.

---

## Service

The service controls the WebSocket lifecycle.

Responsibilities:

| Responsibility     | Description                            |
| ------------------ | -------------------------------------- |
| Accept connections | Handles incoming WebSocket connections |
| Identify users     | Uses `GetUserID`                       |
| Register clients   | Adds users to the hub                  |
| Execute hooks      | Calls lifecycle callbacks              |
| Handle disconnects | Cleans connections                     |

---

## Hub

The Hub is the connection registry.

Responsibilities:

| Function       | Purpose                      |
| -------------- | ---------------------------- |
| Store clients  | Maintains active connections |
| Group users    | Maps users to connections    |
| Send messages  | Delivers messages            |
| Track presence | Knows online users           |

---

# Quick Start

## Define Payload

The library uses generics.

Any Go type can be used.

Example:

```go
type ChatPayload struct {
	Text string `json:"text"`
}
```

---

## Create Manager

```go
socket := ws.NewManager[ChatPayload](
	ws.Config[ChatPayload]{

		GetUserID: func(
			conn *websocket.Conn,
		)(uint,error){

			return 1,nil
		},

	},
)
```

---

## Register Route

```go
app.Get(
	"/ws",
	socket.Handler(),
)
```

---

# Configuration

Configuration controls how the WebSocket manager behaves.

```go
type Config[T any] struct {

	GetUserID func(
		conn *websocket.Conn,
	)(uint,error)

	Hooks Hooks[T]
}
```

---

# Config Reference

| Field     | Required | Description                      |
| --------- | -------- | -------------------------------- |
| GetUserID | Yes      | Resolves the connected user's ID |
| Hooks     | No       | Lifecycle callbacks              |

---

# GetUserID

## Purpose

Determines which user owns the WebSocket connection.

Signature:

```go
GetUserID func(
	conn *websocket.Conn,
)(uint,error)
```

---

## Rules

| Rule                       | Description                                 |
| -------------------------- | ------------------------------------------- |
| Must return user ID        | Every connection belongs to a user          |
| Authentication is external | Library does not validate tokens            |
| Errors reject connection   | Failed identification prevents registration |

---

## JWT Example

```go
GetUserID: func(
	conn *websocket.Conn,
)(uint,error){

	token :=
		conn.Headers(
			"Authorization",
		)


	return ValidateToken(token)
}
```

---

# Manager API

## NewManager

Creates a WebSocket manager.

```go
func NewManager[T any](
	config Config[T],
) *Manager[T]
```

---

## Parameters

| Parameter | Description           |
| --------- | --------------------- |
| T         | Payload type          |
| Config    | Manager configuration |

---

## Example

```go
socket :=
ws.NewManager[Notification](
	config,
)
```

---

# Handler()

## Signature

```go
func (m *Manager[T]) Handler() fiber.Handler
```

---

## Purpose

Returns the Fiber WebSocket handler.

---

## Usage

```go
app.Get(
	"/ws",
	socket.Handler(),
)
```

---

## Rules

| Rule                      | Explanation                          |
| ------------------------- | ------------------------------------ |
| Register as GET route     | WebSocket handshake uses GET         |
| Use Fiber v2              | Handler type depends on Fiber v2     |
| Do not wrap unnecessarily | Handler manages connection lifecycle |

---

# SendUser()

## Signature

```go
func (m *Manager[T]) SendUser(
	userID uint,
	data T,
) error
```

---

## Purpose

Send data to a specific user.

---

## Behavior

Example:

```
User 10

Browser
   |
WebSocket

Mobile
   |
WebSocket
```

Calling:

```go
socket.SendUser(
	10,
	message,
)
```

sends to:

```
Browser
Mobile
```

---

## Rules

| Rule                      | Description             |
| ------------------------- | ----------------------- |
| User must be online       | Otherwise ErrOffline    |
| All sessions receive data | Not only one connection |
| Payload must match T      | Generic type safety     |

---

# Message System

The library supports structured messages.

---

# Message[T]

Internal message structure.

Used mainly inside hooks.

```go
type Message[T any] struct {

	ID string

	Type string

	From uint

	To uint

	RoomID uint

	Data T

	Timestamp time.Time
}
```

---

# Message Fields

| Field     | Description        |
| --------- | ------------------ |
| ID        | Message identifier |
| Type      | Message category   |
| From      | Sender             |
| To        | Receiver           |
| RoomID    | Optional room      |
| Data      | Actual payload     |
| Timestamp | Creation time      |

---

# Client

Represents a connected user session.

```go
type Client struct {

	ID string

	UserID uint

	Conn *websocket.Conn

	Send chan any
}
```

---

# Client Fields

| Field  | Description              |
| ------ | ------------------------ |
| ID     | Unique connection ID     |
| UserID | Owner of connection      |
| Conn   | WebSocket connection     |
| Send   | Outgoing message channel |

---

# Hooks

Hooks allow custom behavior without modifying the library.

```go
type Hooks[T any] struct {

	OnConnect func(
		client *Client,
	) error


	BeforeSend func(
		message *Message[T],
	) error


	OnDisconnect func(
		client *Client,
	)

}
```

---

# OnConnect

Runs after successful connection.

Common uses:

- Presence updates
- Logging
- Analytics

Example:

```go
OnConnect: func(
	client *ws.Client,
) error {

	log.Println(
		"online",
		client.UserID,
	)

	return nil
}
```

---

# BeforeSend

Runs before message delivery.

Uses:

- Logging
- Auditing
- Message transformation

Example:

```go
BeforeSend: func(
	message *ws.Message[ChatPayload],
) error {

	log.Println(
		message,
	)

	return nil
}
```

---

# OnDisconnect

Runs when a connection closes.

Uses:

- Presence updates
- Cleanup

Example:

```go
OnDisconnect: func(
	client *ws.Client,
){

	log.Println(
		"offline",
		client.UserID,
	)

}
```

---

# Hub API

The Hub manages active connections.

---

| Method         | Purpose                      |
| -------------- | ---------------------------- |
| Register       | Add client                   |
| Unregister     | Remove client                |
| DisconnectUser | Disconnect all user sessions |
| Send           | Send structured message      |
| SendUser       | Send payload                 |
| UserClients    | Get user sessions            |
| Connected      | Check online state           |
| Count          | Total connections            |

---

# Register

```go
Register(
	client *Client,
) error
```

Adds a client connection.

Normally called internally.

---

# Unregister

```go
Unregister(
	client *Client,
) error
```

Removes disconnected clients.

---

# DisconnectUser

```go
DisconnectUser(
	userID uint,
) error
```

Disconnects all sessions.

Useful for:

- Logout everywhere
- Account suspension
- Security events

---

# UserClients

```go
UserClients(
	userID uint,
) []*Client
```

Returns all active connections.

---

# Connected

```go
Connected(
	userID uint,
) bool
```

Checks if user is online.

Example:

```go
if socket.Connected(10){

	fmt.Println(
		"online",
	)

}
```

---

# Count

```go
Count() int
```

Returns total active connections.

---

# Connection Lifecycle

```
Client Connects

       |

       v

Fiber WebSocket Route

       |

       v

GetUserID()

       |

       v

Client Created

       |

       v

Hub.Register()

       |

       v

OnConnect()

       |

       v

Connection Active


       |

       v

Disconnect


       |

       v

OnDisconnect()


       |

       v

Hub.Unregister()
```

---

# Authentication Integration

The library intentionally does not know:

- Users
- Passwords
- JWT
- Sessions
- Permissions

You provide authentication.

Supported approaches:

| Method      | Supported |
| ----------- | --------- |
| JWT         | Yes       |
| Cookies     | Yes       |
| Sessions    | Yes       |
| API tokens  | Yes       |
| Custom auth | Yes       |

---

# Usage Patterns

## Notification Service

```go
func Notify(
	userID uint,
	message string,
){

	socket.SendUser(
		userID,
		Notification{
			Message:message,
		},
	)

}
```

---

## Chat

```go
type Chat struct {

	Text string `json:"text"`

}
```

---

## Live Dashboard

```go
type Event struct {

	Type string `json:"type"`

	Data any `json:"data"`

}
```

---

# Rules and Design Principles

## Rule 1

The library manages transport only.

Do:

- Send messages
- Manage connections
- Track clients

Do not:

- Store business data
- Manage users
- Handle permissions

---

## Rule 2

Authentication belongs to the application.

The library only asks:

"Who is this connection?"

---

## Rule 3

One Manager per WebSocket domain.

Example:

```
Chat Manager

Notification Manager

Game Manager
```

---

## Rule 4

Keep business logic outside hooks.

Hooks should:

- Notify
- Log
- Trigger events

Avoid:

- Heavy database operations
- Long processing

---

# Scaling

Single instance:

```
Application

    |

   Hub

    |

 Users
```

Multiple instances:

```
             Load Balancer


              /       \


          Server A   Server B


             |          |

           Hub        Hub
```

For distributed systems add:

- Redis Pub/Sub
- Message broker
- Shared event bus

---

# Recommended Structure

```
project

├── main.go

├── websocket

│     setup.go


├── services

│     notifications.go


├── handlers

├── database

└── models
```

---

# License

MIT

---

# Author

Kigongo Vincent

GitHub:

https://github.com/kigongo-vincent
