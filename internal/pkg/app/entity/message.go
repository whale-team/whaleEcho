package entity

import (
	"time"

	"github.com/nats-io/nats.go"
	"github.com/whale-team/whaleEcho/internal/pkg/app/entity/value"
	"github.com/whale-team/whaleEcho/pkg/echoproto"
	"google.golang.org/protobuf/proto"
)

var sysMessageSet = map[value.SysMsgType]*SysMessage{
	value.CloseRoom: RoomCloseMessage,
}

// NewSysMessage construct a sys message, return pointer
func NewSysMessage(text string) *SysMessage {
	return &SysMessage{
		Message: Message{
			Text: text,
		},
	}
}

// SysMessage represent message sent from system
type SysMessage struct {
	Message
}

// GetData marshal message to bytes data
func (m *SysMessage) GetData() []byte {
	msgProto := &echoproto.Message{}
	msgProto.Text = m.Text
	msgProto.Sender = &echoproto.User{
		Name: "system",
	}
	msgProto.SentAt = time.Now().Unix()

	data, _ := proto.Marshal(msgProto)
	return data
}

// RoomCloseMessage 房間關閉的系統訊息
var RoomCloseMessage = NewSysMessage("room is closed.")

// Message represent websocket message received from client
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

// SetRawData method set the raw data of message
func (m *Message) SetRawData(data []byte) {
	m.rawData = data
}

// GetData method get data from nats message
func (m *Message) GetData() []byte {
	return m.Msg.Data
}

// Subject method to get message's subject
func (m *Message) Subject() string {
	return m.Room.Subject()
}

// SentAtTime method parse unix to time
func (m *Message) SentAtTime() time.Time {
	return time.Unix(m.SentAt, 0)
}

// GetRawData method return raw data of message
func (m *Message) GetRawData() []byte {
	return m.rawData
}
