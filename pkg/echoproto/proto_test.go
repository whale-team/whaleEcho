package echoproto

import (
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
)

func TestCommandProto(t *testing.T) {
	room := &Room{
		Name: "test",
	}
	data, err := proto.Marshal(room)
	assert.Nil(t, err)

	command := &Command{
		Type:    CommandType_CreateRoom,
		Payload: data,
	}

	data, err = proto.Marshal(command)
	assert.Nil(t, err)

	command2 := &Command{}
	err = proto.Unmarshal(data, command2)
	assert.Nil(t, err)
	assert.Equal(t, command2.Type, command.Type)

	room2 := &Room{}
	err = proto.Unmarshal(command2.Payload, room2)
	assert.Nil(t, err)
	assert.Equal(t, room2.Name, room.Name)
}
