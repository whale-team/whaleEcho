package converter

import (
	"github.com/whale-team/whaleEcho/internal/pkg/app/entity"
	"github.com/whale-team/whaleEcho/pkg/echoproto"
	"github.com/whale-team/whaleEcho/pkg/wserrors"
	"google.golang.org/protobuf/proto"
)

// UnmarshalMessage convert proto message to entity message
func UnmarshalMessage(data []byte, msg *entity.Message) error {
	var (
		msgProto = echoproto.Message{}
		err      error
	)
	if err = proto.Unmarshal(data, &msgProto); err != nil {
		return wserrors.Wrapf(wserrors.ErrInputInvalid, "delivery: umarshal message occurs error, err:%+v ", err)
	}

	msg.Data = data
	msg.RoomUID = msgProto.RoomUid
	msg.SenderName = msgProto.SenderName

	return nil
}

// UnmarshalUser convert proto user to entity user
func UnmarshalUser(data []byte, user *entity.User) error {
	var (
		userProto = echoproto.User{}
		err       error
	)
	if err = proto.Unmarshal(data, &userProto); err != nil {
		return wserrors.Wrapf(wserrors.ErrInputInvalid, "delivery: umarshal user occurs error, err:%+v ", err)
	}
	user.UID = userProto.Uid
	user.Name = userProto.Name
	user.RoomUID = userProto.RoomUid
	return nil
}

// UnmarshalRoom convert proto room to entity room
func UnmarshalRoom(data []byte, room *entity.Room) error {
	var (
		roomProto = echoproto.Room{}
		err       error
	)
	if err = proto.Unmarshal(data, &roomProto); err != nil {
		return wserrors.Wrapf(wserrors.ErrInputInvalid, "delivery: umarshal room occurs error, err:%+v ", err)
	}
	room.UID = roomProto.Uid
	room.Name = roomProto.Name
	room.MembersLimit = roomProto.MembersLimit
	return nil
}
