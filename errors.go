package main

import "errors"

var (
	ErrUnauthorized = errors.New(
		"unauthorized",
	)

	ErrOffline = errors.New(
		"user offline",
	)

	ErrInvalidMessage = errors.New(
		"invalid message",
	)
)
