package db

import (
	"strconv"

	"github.com/whale-team/whaleEcho/internal/pkg/app/entity"
)

func RoomToMap(room *entity.Room) map[string]interface{} {
	return map[string]interface{}{
		"uid":          room.UID,
		"membersCount": room.MembersCount,
		"membersLimit": room.MembersLimit,
		"name":         room.Name,
	}
}

func MapToRoom(source map[string]string, room *entity.Room) error {
	room.UID = source["uid"]
	room.Name = source["name"]
	memberCount, err := strconv.ParseInt(source["membersCount"], 10, 64)
	if err != nil {
		return err
	}
	room.MembersCount = int64(memberCount)
	limit, err := strconv.ParseInt(source["membersLimit"], 10, 64)
	if err != nil {
		return err
	}
	room.MembersLimit = int64(limit)
	return nil
}
