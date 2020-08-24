package service

import (
	"context"

	"github.com/whale-team/whaleEcho/internal/pkg/app/entity"
	"github.com/whale-team/whaleEcho/pkg/subjects"
)

func (s service) DispatchMessage(ctx context.Context, roomUID string, msg *entity.Message) error {
	return nil
}

func (s service) PublishMessage(ctx context.Context, msg *entity.Message) error {
	if err := s.broker.Publish(ctx, subjects.RoomMsgSubject, msg.Data); err != nil {
		return err
	}
	return nil
}
