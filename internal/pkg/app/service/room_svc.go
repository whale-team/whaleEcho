package service

import (
	"context"
	"errors"

	"github.com/whale-team/whaleEcho/internal/pkg/app/entity"
	"github.com/whale-team/whaleEcho/internal/pkg/repository/db"
)

var (
	ErrRoomOutOfLimit = errors.New("cannot join room, number of mebmers out of limit")
)

// CreateRoom create room is not exist,
//  register the room on dispatcher so that room can be assoicated with client connection
func (s service) CreateRoom(ctx context.Context, room *entity.Room) error {
	err := s.repo.GetRoom(ctx, room.UID, room)
	if err != nil && !errors.Is(err, db.ErrNotFound) {
		return err
	}

	if errors.Is(err, db.ErrNotFound) {
		if err := s.repo.CreateRoom(ctx, room); err != nil {
			return err
		}
	}

	if err := s.dispatcher.RegisterRoom(ctx, room); err != nil {
		return err
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
		return err
	}

	if memberCount, err = s.repo.IncrMember(ctx, roomUID); err != nil {
		return err
	}

	if memberCount > room.MembersLimit {
		if _, err = s.repo.DecrMember(ctx, roomUID); err != nil {

		}
		return ErrRoomOutOfLimit
	}

	if err = s.JoinAndCheckRoom(ctx, room, user); err != nil {
		// log message join room failed
		_, err = s.repo.DecrMember(ctx, roomUID)
	}
	if err != nil {
		// log decrease mebmer failed
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
	return nil
}

func (s service) CloseRoom(ctx context.Context, roomUID string) error {
	return nil
}
