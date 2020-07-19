package echoproto

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

// func TestCommandProto(t *testing.T) {
// 	room := &Room{
// 		Name: "test",
// 	}
// 	data, err := proto.Marshal(room)
// 	assert.Nil(t, err)

// 	command := &Command{
// 		Type:    CommandType_CreateRoom,
// 		Payload: data,
// 	}

// 	data, err = proto.Marshal(command)
// 	assert.Nil(t, err)

// 	command2 := &Command{}
// 	err = proto.Unmarshal(data, command2)
// 	assert.Nil(t, err)
// 	assert.Equal(t, command2.Type, command.Type)

// 	room2 := &Room{}
// 	err = proto.Unmarshal(command2.Payload, room2)
// 	assert.Nil(t, err)
// 	assert.Equal(t, room2.Name, room.Name)
// }

func TestChiness(t *testing.T) {
	msg := &Message{
		Uid:  "12345678901234567890",
		Text: "你好嗎，我很好，你好嗎，我很好，你好嗎，我很好，你好嗎，我很好，你好嗎，我很好，你好嗎，我很好，你好嗎，我很好，你好嗎，我很好，",
		Type: MessageType_Text,
		Sender: &User{
			Id:   1,
			Name: "隔壁的小王",
		},
	}
	data, err := proto.Marshal(msg)
	assert.Nil(t, err)
	msg2 := &Message{}

	err = proto.Unmarshal(data, msg2)
	assert.Nil(t, err)
	assert.Equal(t, msg.Text, msg2.Text)
}

func BenchmarkProto(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		msg := &Message{
			Uid:  "12345678901234567890",
			Text: "你好嗎，我很好，你好嗎，我很好，你好嗎，我很好，你好嗎，我很好，你好嗎，我很好，你好嗎，我很好，你好嗎，我很好，你好嗎，我很好，",
			Type: MessageType_Text,
			Sender: &User{
				Id:   1,
				Name: "隔壁的小王",
			},
			Room: &Room{
				Uid: "14293403902afkasdlgk23423fsdf",
			},
		}
		data, _ := proto.Marshal(msg)
		msg2 := &Message{}
		proto.Unmarshal(data, msg2)
	}
}

func BenchmarkJson(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		msg := &Message{
			Uid:  "12345678901234567890",
			Text: "你好嗎，我很好，你好嗎，我很好，你好嗎，我很好，你好嗎，我很好，你好嗎，我很好，你好嗎，我很好，你好嗎，我很好，你好嗎，我很好，",
			Type: MessageType_Text,
			Sender: &User{
				Id:   1,
				Name: "隔壁的小王",
			},
			Room: &Room{
				Uid: "14293403902afkasdlgk23423fsdf",
			},
		}
		data, _ := json.Marshal(msg)
		msg2 := &Message{}
		json.Unmarshal(data, msg2)
	}
}
