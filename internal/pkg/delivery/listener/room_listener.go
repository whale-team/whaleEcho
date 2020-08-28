package listener

import (
	"context"

	"github.com/whale-team/whaleEcho/internal/pkg/app/entity"
	"github.com/whale-team/whaleEcho/internal/pkg/delivery/converter"
)

func (ln Listener) CreateRoom(ctx context.Context, data []byte) error {
	room := &entity.Room{}

	if err := converter.UnmarshalRoom(data, room); err != nil {
		return err
	}

	return ln.svc.CreateRoom(ctx, room)
}

func (ln Listener) DispatchMessage(ctx context.Context, data []byte) error {
	msg := &entity.Message{}

	if err := converter.UnmarshalMessage(data, msg); err != nil {
		return err
	}

	return ln.svc.DispatchMessage(ctx, msg)
}
