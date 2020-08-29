package listener

import (
	"context"

	"github.com/whale-team/whaleEcho/internal/pkg/app/entity"
	"github.com/whale-team/whaleEcho/internal/pkg/delivery/converter"
)

func (ln Listener) CreateRoom(ctx context.Context, data []byte) error {
	var (
		room = &entity.Room{}
		err  error
	)

	if err = converter.UnmarshalRoom(data, room); err != nil {
		return err
	}

	return ln.svc.CreateRoom(ctx, room)
}

func (ln Listener) DispatchMessage(ctx context.Context, data []byte) error {
	var (
		msg = &entity.Message{}
		err error
	)

	if err = converter.UnmarshalMessage(data, msg); err != nil {
		return err
	}

	return ln.svc.DispatchMessage(ctx, msg)
}

func (ln Listener) CloseRoom(ctx context.Context, data []byte) error {
	var (
		room = &entity.Room{}
		err  error
	)

	if err = converter.UnmarshalRoom(data, room); err != nil {
		return err
	}

	return ln.svc.CloseRoom(ctx, room.UID)
}
