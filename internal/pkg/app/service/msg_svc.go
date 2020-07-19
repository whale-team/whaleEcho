package service

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/whale-team/whaleEcho/internal/pkg/app/entity"
	"github.com/whale-team/whaleEcho/pkg/natspool"
	"github.com/whale-team/whaleEcho/pkg/wserrors"
)

// PublishText method used to publish text message
func (svc service) PublishText(ctx context.Context, msg *entity.Message) error {
	if err := svc.msgBroker.PublishMessage(ctx, msg.Subject(), msg.GetRawData()); err != nil {
		if wserrors.Is(err, natspool.ErrGetConnTimeout) {
			return wserrors.Wrapf(wserrors.ErrSysBusy, "svc: PublishMessage msgbroker get connection timeout")
		}

		return wserrors.Wrapf(wserrors.ErrInternal, "svc: PublishMessage failed, err:%+v, subject:%s msg:%+v", err, msg.Subject(), msg)
	}
	return nil
}

// JoinRoom method used to join user to specific room
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
		return wserrors.Wrapf(wserrors.ErrInternal, "svc: JoinRoom failed, err:%+v, roomUID:%s user:%+v", err, roomUID, user)
	}
	return nil
}

// LeaveRoom method used to remove user from a specific room
func (svc service) LeaveRoom(ctx context.Context, roomUID string, user *entity.User) error {
	svc.rooms.LeaveRoom(roomUID, user)
	return nil
}
