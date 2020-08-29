package service

import (
	"context"
	"errors"

	"github.com/rs/zerolog/log"
	"github.com/whale-team/whaleEcho/internal/pkg/app/entity"
	"github.com/whale-team/whaleEcho/internal/pkg/repository/db"
	"github.com/whale-team/whaleEcho/pkg/wserrors"
)

// CreateRoom create room is not exist,
//  register the room on dispatcher so that room can be assoicated with client connection
func (s service) CreateRoom(ctx context.Context, room *entity.Room) error {
	err := s.repo.GetRoom(ctx, room.UID, room)
	if err != nil && !errors.Is(err, db.ErrNotFound) {
		return wserrors.Wrapf(wserrors.ErrRoomNotFound, "svc: CreateRoom room not found, roomUID:%s", room.UID)
	}

	if errors.Is(err, db.ErrNotFound) {
		if err := s.repo.CreateRoom(ctx, room); err != nil {
			return wserrors.Wrapf(wserrors.ErrInternal, "svc: CreateRoom repo create room occur err, err:%+v roomUID:%s", err, room.UID)
		}
	}

	if err := s.dispatcher.RegisterRoom(ctx, room); err != nil {
		return wserrors.Wrapf(wserrors.ErrInternal, "svc: CreateRoom dispatcher register room occur err, err:%+v roomUID:%s", err, room.UID)
	}

	return nil
}

// JoinRoom join a user to the room, increase room's membersCount
//  ensure room has been registered in dispatcher
func (s service) JoinRoom(ctx context.Context, roomUID string, user *entity.User) error {
	var (
		err         error
		memberCount int64
		room        = &entity.Room{}
	)

	if err = s.repo.GetRoom(ctx, roomUID, room); err != nil {
		return wserrors.Wrapf(wserrors.ErrRoomNotFound, "svc: JoinRoom get room not found, roomUID:%s", roomUID)
	}

	if memberCount, err = s.repo.IncrMember(ctx, roomUID); err != nil {
		return wserrors.Wrapf(wserrors.ErrInternal, "svc: JoinRoom increase member failed, roomUID:%s", roomUID)
	}

	if memberCount > room.MembersLimit {
		if _, err = s.repo.DecrMember(ctx, roomUID); err != nil {
			log.Error().Stack().Err(err).Msgf("svc: JoinRoom decrease member failed, roomUID:%s", roomUID)
		}
		return wserrors.ErrRoomOutOfLimit
	}

	if err = s.JoinAndCheckRoom(ctx, room, user); err != nil {
		log.Error().Stack().Err(err).Msgf("svc: JoinRoom join user to room through dispatcher failed, roomUID:%s", roomUID)
		_, err = s.repo.DecrMember(ctx, roomUID)
	}
	if err != nil {
		log.Error().Stack().Err(err).Msgf("svc: JoinRoom decrease member failed, roomUID:%s", roomUID)
		return err
	}

	return nil
}

// JoinAndCheckRoom join user to room after checking if room is registered in dispatcher
func (s service) JoinAndCheckRoom(ctx context.Context, room *entity.Room, user *entity.User) error {
	var err error

	if err = s.checkRoomInDispatcher(ctx, room); err != nil {
		return err
	}

	if err := s.dispatcher.JoinUserToRoom(ctx, room.UID, user); err != nil {
		return err
	}

	return nil
}

func (s service) checkRoomInDispatcher(ctx context.Context, room *entity.Room) error {
	var (
		registeredRoom = s.dispatcher.GetRoom(ctx, room.UID)
		err            error
	)

	if registeredRoom == nil {
		err = s.dispatcher.RegisterRoom(ctx, room)
	}
	if err != nil {
		return err
	}
	return nil
}

func (s service) LeaveRoom(ctx context.Context, roomUID string, userID string) error {
	var err error

	if err = s.dispatcher.RemoveUserFromRoom(ctx, roomUID, userID); err != nil {
		return wserrors.Wrapf(wserrors.ErrInternal, "svc: LeaveRoom dispatcher remove user from room occur error, err:%+v", err)
	}
	_, err = s.repo.DecrMember(ctx, roomUID)
	if err != nil {
		return wserrors.Wrapf(wserrors.ErrInternal, "svc: LeaveRoom repo decrease room's members occur error, err:%+v", err)
	}

	return nil
}

// CloseRoom close room delete room from redis and dispatcher
func (s service) CloseRoom(ctx context.Context, roomUID string) error {
	var (
		err  error
		room = entity.Room{}
	)

	err = s.repo.GetRoom(ctx, roomUID, &room)
	if err != nil && !wserrors.Is(err, db.ErrNotFound) {
		return wserrors.Wrapf(wserrors.ErrInternal, "svc: CloseRoom repo get room occur error, err:%+v, roomUID:%s", err, roomUID)
	}
	if room.UID != "" {
		err = s.repo.DeleteRoom(ctx, roomUID)
	}
	if err != nil {
		return wserrors.Wrapf(wserrors.ErrInternal, "svc: CloseRoom repo delete room occur error, err:%+v, roomUID:%s", err, roomUID)
	}

	if err = s.dispatcher.CloseRoom(ctx, roomUID); err != nil {
		return wserrors.Wrapf(wserrors.ErrInternal, "svc: CloseRoom dispatcher close room occur error, err:%+v, roomUID:%s", err, roomUID)
	}

	return nil
}
