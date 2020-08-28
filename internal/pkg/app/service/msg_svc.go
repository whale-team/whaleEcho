package service

import (
	"context"

	"github.com/whale-team/whaleEcho/internal/pkg/app/entity"
	"github.com/whale-team/whaleEcho/pkg/subjects"
)

func (s service) DispatchMessage(ctx context.Context, msg *entity.Message) error {
	s.dispatcher.DispatchMessage(ctx, msg)
	return nil
}

func (s service) PublishRoomMessage(ctx context.Context, data []byte) error {
	if err := s.broker.Publish(ctx, subjects.RoomMsgSubject, data); err != nil {
		return err
	}
	return nil
}
