package service

import (
	"context"
	"errors"

	"github.com/whale-team/whaleEcho/internal/pkg/app/entity"
	"github.com/whale-team/whaleEcho/internal/pkg/repository/db"
)

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
		err = s.repo.DeleteRoom(ctx, room.UID)
		return err
	}

	return nil
}

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

	if err = s.checkRunningRoom(ctx, room); err != nil {
		return err
	}

	if memberCount > room.MembersLimit {
		_, err = s.repo.DecrMember(ctx, roomUID)
		return err
	}

	if err := s.dispatcher.JoinUserToRoom(ctx, roomUID, user); err != nil {
		_, err = s.repo.DecrMember(ctx, roomUID)
		return err
	}

	return nil
}

func (s service) checkRunningRoom(ctx context.Context, room *entity.Room) error {
	var (
		runingRoom = s.dispatcher.GetRoom(ctx, room.UID)
		err        error
	)
	if runingRoom == nil {
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
