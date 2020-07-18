package roomscenter

import (
	"context"

	"github.com/nats-io/nats.go"
	"github.com/whale-team/whaleEcho/internal/pkg/app/entity"
	"github.com/whale-team/whaleEcho/internal/pkg/app/mapper"
	"github.com/whale-team/whaleEcho/internal/pkg/app/msgbroker"
)

type AsyncHandler interface {
	ErrHandle(error, *nats.Msg)
	OpenRoom(context.Context, *nats.Msg, *entity.Room) error
	CloseRoom(context.Context, *nats.Msg) (string, error)
}

func NewDefaultHandler(broker msgbroker.MsgBroker) AsyncHandler {
	return &DefaultHandler{
		broker: broker,
	}
}

type DefaultHandler struct {
	broker msgbroker.MsgBroker
}

func (h DefaultHandler) OpenRoom(ctx context.Context, msg *nats.Msg, room *entity.Room) error {
	err := mapper.UnmarshalRoom(msg.Data, room)
	if err != nil {
		return err
	}

	msgCh := make(chan *nats.Msg, 1)
	room.SetMsgChannel(msgCh)
	sub, err := h.broker.BindChannelWithSubject(ctx, room.Subject(), msgCh)
	if err != nil {
		return err
	}
	room.Subscribe = sub
	return nil
}

func (h DefaultHandler) CloseRoom(ctx context.Context, msg *nats.Msg) (string, error) {
	room := &entity.Room{}
	err := mapper.UnmarshalRoom(msg.Data, room)
	if err != nil {
		return "", nil
	}
	return room.UID, nil
}

func (h DefaultHandler) ErrHandle(err error, msg *nats.Msg) {

}
