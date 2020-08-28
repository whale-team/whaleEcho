package delivery_test

import (
	"context"
	"time"

	"github.com/gorilla/websocket"
	. "github.com/onsi/ginkgo"
	"github.com/stretchr/testify/assert"
	"github.com/whale-team/whaleEcho/pkg/echoproto"
	"github.com/whale-team/whaleEcho/pkg/subjects"
	"google.golang.org/protobuf/proto"
)

func publishMessage(suite *testSuite, conn *websocket.Conn, message *echoproto.Message, status echoproto.Status) *websocket.Conn {
	var err error
	if conn == nil {
		conn, err = dial(suite.addr, suite.port)
	}
	assert.NoError(suite.T, err)
	command, err := newCommand(message, echoproto.CommandType_SendMessage)
	assert.NoError(suite.T, err)
	err = suite.SendCommand(conn, &command)
	assert.NoError(suite.T, err)
	time.Sleep(5 * time.Millisecond)
	msg, err := suite.ReadResp(conn)
	assert.NoError(suite.T, err)
	assert.Equal(suite.T, msg.Status, status)
	assert.Contains(suite.T, suite.buf.String(), "send_message")
	return conn
}

var _ = Describe("Message Delivery Test", func() {
	AfterEach(func() {
		suite.clear()
	})

	Describe("Handler#PublishMessage", func() {
		Context("when client send message", func() {
			var msg = &echoproto.Message{}
			var i int
			suite.stan.Subscribe(subjects.RoomMsgSubject, func(ctx context.Context, data []byte) error {
				proto.Unmarshal(data, msg)
				i++
				return nil
			})
			It("should receive message", func() {
				msgData := suite.messages[0]
				publishMessage(suite, nil, msgData, echoproto.Status_OK)
				assert.Equal(suite.T, msgData.Uid, msg.Uid)
				assert.Equal(suite.T, i, 1)
			})
		})
	})

	FDescribe("Listener#DispatcherMessage", func() {
		Context("when client send message", func() {
			var (
				room  = suite.rooms[0]
				user1 = suite.users[0]
				user2 = suite.users[1]
				conn1 = &websocket.Conn{}
				conn2 = &websocket.Conn{}
			)

			BeforeEach(func() {
				createRoom(suite, room)
				user1.RoomUid = room.Uid
				user2.RoomUid = room.Uid
				conn1 = joinRoom(suite, user1, echoproto.Status_OK)
				conn2 = joinRoom(suite, user2, echoproto.Status_OK)
			})

			It("should dispatch message to user2", func() {
				msgData := suite.messages[0]
				msgData.RoomUid = room.Uid
				msgData.SenderName = user1.Name
				publishMessage(suite, conn1, msgData, echoproto.Status_OK)

				msg, err := suite.ReadResp(conn2)
				assert.NoError(suite.T, err)
				assert.Equal(suite.T, msg.Text, msgData.Text)
				msg, err = suite.ReadResp(conn1)
				assert.NoError(suite.T, err)
				assert.Equal(suite.T, msg.Text, msgData.Text)
			})
		})
	})
})
