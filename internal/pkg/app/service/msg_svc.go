package service

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/whale-team/whaleEcho/internal/pkg/app/entity"
)

func (svc service) PublishText(ctx context.Context, msg *entity.Message) error {
	if err := svc.msgBroker.PublishMessage(ctx, msg.Subject(), msg.ToMsgData()); err != nil {
		return err
	}
	return nil
}

func (svc service) JoinRoom(ctx context.Context, roomUID string, user *entity.User) error {
	if !svc.rooms.HasRoom(roomUID) {
		log.Warn().Msg("join not existed room")
		room := entity.NewRoom()
		room.UID = roomUID
		if err := svc.rooms.AddRoom(room); err != nil {
			return err
		}
		log.Warn().Msgf("add room(%s", roomUID)
	}

	if err := svc.rooms.JoinRoom(roomUID, user); err != nil {
		return err
	}
	return nil
}

func (svc service) LeaveRoom(ctx context.Context, roomUID string, user *entity.User) error {
	svc.rooms.LeaveRoom(roomUID, user)
	return nil
}
