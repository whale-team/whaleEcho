package entity

import (
	"time"

	"github.com/nats-io/nats.go"
	"github.com/whale-team/whaleEcho/internal/pkg/app/entity/value"
	"github.com/whale-team/whaleEcho/pkg/echoproto"
	"google.golang.org/protobuf/proto"
)

type SysMessage Message

// RoomCloseMessage 房間關閉的系統訊息
var RoomCloseMessage = &Message{
	Text: "room is closed",
}

// Message represent websocket message
type Message struct {
	*nats.Msg
	UID      string
	Text     string
	File     []byte
	FileType string
	Type     value.MessageType
	Room     Room
	Sender   User
	SentAt   int64

	rawData []byte
}

func (m *Message) SetRawData(data []byte) {
	m.rawData = data
}

func (m *Message) GetData() []byte {
	if m.Msg == nil {
		return m.toData()
	}

	return m.Msg.Data
}

func (m Message) toData() []byte {
	msgProto := &echoproto.Message{}
	msgProto.Text = m.Text

	data, _ := proto.Marshal(msgProto)
	return data
}

func (m Message) Subject() string {
	return m.Room.Subject()
}

func (m Message) SentAtTime() time.Time {
	return time.Unix(m.SentAt, 0)
}

func (m *Message) ToMsgData() []byte {
	return m.rawData
}
