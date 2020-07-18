package mapper

import (
	"github.com/whale-team/whaleEcho/internal/pkg/app/entity"
	"github.com/whale-team/whaleEcho/internal/pkg/app/entity/value"
	"github.com/whale-team/whaleEcho/pkg/echoproto"
	"google.golang.org/protobuf/proto"
)

func UnmarshalRoom(data []byte, room *entity.Room) error {
	protoRoom := echoproto.Room{}
	if err := proto.Unmarshal(data, &protoRoom); err != nil {
		return err
	}
	room.ID = protoRoom.Id
	room.UID = protoRoom.Uid
	room.Limit = protoRoom.ParticipantsLimit
	return nil
}

func UnmarshalRoomAndUser(data []byte, room *entity.Room, user *entity.User) error {
	protoRoom := echoproto.Room{}
	if err := proto.Unmarshal(data, &protoRoom); err != nil {
		return err
	}
	room.ID = protoRoom.Id
	room.UID = protoRoom.Uid
	room.Limit = protoRoom.ParticipantsLimit

	user.ID = protoRoom.Participant.Id
	user.Name = protoRoom.Participant.Name
	return nil
}

func UnmarshalMessage(data []byte, msg *entity.Message) error {
	protoMsg := echoproto.Message{}
	if err := proto.Unmarshal(data, &protoMsg); err != nil {
		return err
	}
	msg.SetRawData(data)
	msg.Room.ID = protoMsg.Room.Id
	msg.Room.UID = protoMsg.Room.Uid
	msg.Sender.Name = protoMsg.Sender.Name
	msg.Sender.ID = protoMsg.Sender.Id
	msg.Type = value.MessageType(protoMsg.Type)
	msg.Text = protoMsg.Text
	msg.File = protoMsg.File
	msg.FileType = protoMsg.FileType
	msg.SentAt = protoMsg.SentAt
	return nil
}
