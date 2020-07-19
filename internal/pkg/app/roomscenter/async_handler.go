package roomscenter

import (
	"context"

	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"
	"github.com/whale-team/whaleEcho/internal/pkg/app/entity"
	"github.com/whale-team/whaleEcho/internal/pkg/app/mapper"
	"github.com/whale-team/whaleEcho/internal/pkg/app/msgbroker"
	"github.com/whale-team/whaleEcho/pkg/natspool"
	"github.com/whale-team/whaleEcho/pkg/wserrors"
)

// AsyncHandler provide async handler function that used to handle nats event
type AsyncHandler interface {
	ErrHandle(error, *nats.Msg, *entity.Room)
	OpenRoom(context.Context, *nats.Msg, *entity.Room) error
	CloseRoom(context.Context, *nats.Msg, *entity.Room) error
}

// NewDefaultHandler construct default async handler
func NewDefaultHandler(broker msgbroker.MsgBroker) AsyncHandler {
	return &DefaultHandler{
		broker: broker,
	}
}

// DefaultHandler represent default async handler
type DefaultHandler struct {
	broker msgbroker.MsgBroker
}

// OpenRoom method used to create a new room
func (h DefaultHandler) OpenRoom(ctx context.Context, msg *nats.Msg, room *entity.Room) error {
	err := mapper.UnmarshalRoom(msg.Data, room)
	if err != nil {
		return wserrors.Wrapf(ErrOpenRoom, "err:%+v", err)
	}

	msgCh := make(chan *nats.Msg, 1)
	room.SetMsgChannel(msgCh)

	for i := 0; i < 3; i++ {
		sub, err := h.broker.BindChannelWithSubject(ctx, room.Subject(), msgCh)
		if err != nil {
			if wserrors.Is(err, natspool.ErrGetConnTimeout) {
				continue
			}
			return wserrors.Wrapf(ErrOpenRoom, "err:%+v", err)
		}
		room.Subscribe = sub
		break
	}
	return nil
}

// CloseRoom method used to close a room
func (h DefaultHandler) CloseRoom(ctx context.Context, msg *nats.Msg, room *entity.Room) error {
	err := mapper.UnmarshalRoom(msg.Data, room)
	if err != nil {
		return wserrors.Wrapf(ErrCloseRoom, "err:%+v", err)
	}
	return nil
}

// ErrHandle handle async method error
func (h DefaultHandler) ErrHandle(err error, msg *nats.Msg, room *entity.Room) {
	if wserrors.Is(wserrors.Cause(err), ErrOpenRoom) {
		log.Error().Err(err).Msgf("roomsCenter: open room failed, room:%+v", room)
	} else if wserrors.Is(wserrors.Cause(err), ErrCloseRoom) {
		log.Error().Err(err).Msgf("roomsCenter: close room failed, room:%+v", room)
	}
}