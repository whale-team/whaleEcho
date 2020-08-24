package converter

import (
	"github.com/whale-team/whaleEcho/internal/pkg/app/entity"
	"github.com/whale-team/whaleEcho/pkg/echoproto"
	"google.golang.org/protobuf/proto"
)

func UnmarshalMessage(data []byte, msg *entity.Message) error {
	var (
		msgProto = echoproto.Message{}
		err      error
	)
	if err = proto.Unmarshal(data, &msgProto); err != nil {
		return err
	}

	msg.Data = data
	msg.RoomUID = msgProto.RoomUid
	msg.SenderName = msgProto.SenderName

	return nil
}

func UnmarshalUser(data []byte, user *entity.User) error {
	var (
		userProto = echoproto.User{}
		err       error
	)
	if err = proto.Unmarshal(data, &userProto); err != nil {
		return err
	}
	user.UID = userProto.Uid
	user.Name = userProto.Name
	user.RoomUID = userProto.RoomUid
	return nil
}

func UnmarshalRoom(data []byte, room *entity.Room) error {
	var (
		roomProto = echoproto.Room{}
		err       error
	)
	if err = proto.Unmarshal(data, &roomProto); err != nil {
		return err
	}
	room.UID = roomProto.Uid
	room.Name = roomProto.Name
	room.MembersLimit = roomProto.MembersLimit
	return nil
}
