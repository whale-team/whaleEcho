package service

import (
	"context"

	"github.com/whale-team/whaleEcho/internal/pkg/app/entity"
	"github.com/whale-team/whaleEcho/internal/pkg/app/msgbroker"
	"github.com/whale-team/whaleEcho/internal/pkg/app/roomscenter"
)

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

// Servicer service facade interface
type Servicer interface {
	JoinRoom(ctx context.Context, roomUID string, user *entity.User) error
	PublishText(ctx context.Context, msg *entity.Message) error
	LeaveRoom(ctx context.Context, roomUID string, user *entity.User) error
}
