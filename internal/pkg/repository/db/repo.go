package db

import (
	"context"
	"errors"

	"github.com/go-redis/redis/v8"
	"github.com/whale-team/whaleEcho/internal/pkg/app/entity"
)

var (
	ErrNotFound = errors.New("record not found")
)

const roomKey = "room"

type Repo struct {
	redisDB *redis.Client
}

func (repo *Repo) CreateRoom(ctx context.Context, room *entity.Room) error {
	roomMap := RoomToMap(room)

	if _, err := repo.redisDB.HMSet(ctx, roomKey+"."+room.UID, roomMap).Result(); err != nil {
		return err
	}

	return nil
}

func (repo *Repo) IncrMember(ctx context.Context, roomUID string) (int64, error) {
	return repo.redisDB.HIncrBy(ctx, roomKey+"."+roomUID, "membersCount", 1).Result()
}

func (repo *Repo) DecrMember(ctx context.Context, roomUID string) (int64, error) {

	return repo.redisDB.HIncrBy(ctx, roomKey+"."+roomUID, "membersCount", -1).Result()
}

func (repo *Repo) GetRoom(ctx context.Context, roomUID string, room *entity.Room) error {
	mapData, err := repo.redisDB.HGetAll(ctx, roomKey+"."+roomUID).Result()
	if err != nil {
		return err
	}
	if len(mapData) == 0 {
		return ErrNotFound
	}

	err = MapToRoom(mapData, room)
	if err != nil {
		return err
	}

	return nil
}

func (repo *Repo) DeleteRoom(ctx context.Context, roomUID string) error {
	_, err := repo.redisDB.Del(ctx, roomKey+"."+roomUID).Result()
	return err
}
