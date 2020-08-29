package service

import (
	"context"

	"github.com/whale-team/whaleEcho/internal/pkg/app/entity"
	"github.com/whale-team/whaleEcho/pkg/subjects"
)

// DispatchMessage dispatch message to all users
func (s service) DispatchMessage(ctx context.Context, msg *entity.Message) error {
	s.dispatcher.DispatchMessage(ctx, msg)
	return nil
}

// PublishRoomMessage publish message to room
func (s service) PublishRoomMessage(ctx context.Context, data []byte) error {
	if err := s.broker.Publish(ctx, subjects.RoomMsgSubject, data); err != nil {
		return err
	}
	return nil
}
