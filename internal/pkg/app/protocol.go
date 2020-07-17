package app

import (
	"context"
)

// Servicer service facade interface
type Servicer interface {
	roomServicer
	msgServicer
}

type msgServicer interface {
	SendMsgToRoom(ctx context.Context, msg entity.Message, room entity.Room) error
}

type roomServicer interface {
	CreateRoom(ctx context.Context, room *entity.Room) error
	JoinRoom(ctx context.Context, roomID string, user entity.User) error
	LeaveRoom(ctx context.Context, roomID string, user entity.User) error
}

// MessageBroker message pub/sub interface
type MessageBroker interface {
	BindChWithSubject(ctx context.Context, subject string, ch chan<- entity.Message) (entity.Subscriber, error)
	PublishMessage(ctx context.Context, subject string, msg entity.Message) error
}
