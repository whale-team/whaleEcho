package service

import (
	"context"

	"github.com/whale-team/whaleEcho/internal/pkg/app/entity"
	"github.com/whale-team/whaleEcho/internal/pkg/app/msgbroker"
	"github.com/whale-team/whaleEcho/internal/pkg/app/roomscenter"
)

// New construct a message servicer
func New(broker msgbroker.MsgBroker, rooms *roomscenter.Center) Servicer {
	return &service{
		msgBroker: broker,
		rooms:     rooms,
	}
}

type service struct {
	msgBroker msgbroker.MsgBroker
	rooms     *roomscenter.Center
}

// Servicer provide message servicer interface
type Servicer interface {
	JoinRoom(ctx context.Context, roomUID string, user *entity.User) error
	PublishText(ctx context.Context, msg *entity.Message) error
	LeaveRoom(ctx context.Context, roomUID string, user *entity.User) error
	LeaveAllRooms(ctx context.Context, user *entity.User) error
}
