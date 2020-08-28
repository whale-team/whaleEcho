package app

import (
	"context"

	"github.com/whale-team/whaleEcho/internal/pkg/app/entity"
)

// Broker broadcast message on specific subject
type Broker interface {
	Publish(ctx context.Context, subject string, data []byte) error
}

// Dispatcher is responsbile for managing client connections
type Dispatcher interface {
	RegisterRoom(ctx context.Context, room *entity.Room) error
	CloseRoom(ctx context.Context, roomUID string) error
	DispatchMessage(ctx context.Context, msg *entity.Message) error
	JoinUserToRoom(ctx context.Context, roomUID string, user *entity.User) error
	GetRoom(ctx context.Context, roomUID string) *entity.Room
}

// Repositorier store room state
type Repositorier interface {
	CreateRoom(ctx context.Context, room *entity.Room) error
	IncrMember(ctx context.Context, roomUID string) (int64, error)
	DecrMember(ctx context.Context, roomUID string) (int64, error)
	GetRoom(ctx context.Context, roomUID string, room *entity.Room) error
	DeleteRoom(ctx context.Context, roomUID string) error
}

// Servicer domain logic
type Servicer interface {
	CreateRoom(ctx context.Context, room *entity.Room) error
	JoinRoom(ctx context.Context, roomUID string, user *entity.User) error
	CloseRoom(ctx context.Context, roomUID string) error
	PublishRoomMessage(ctx context.Context, data []byte) error
	LeaveRoom(ctx context.Context, roomUID string, userID string) error
	DispatchMessage(ctx context.Context, msg *entity.Message) error
}
