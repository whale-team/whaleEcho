package handler_test

import (
	"time"

	"github.com/gorilla/websocket"
	. "github.com/onsi/ginkgo"
	"github.com/stretchr/testify/assert"
	"github.com/whale-team/whaleEcho/internal/pkg/app/roomscenter"
	"github.com/whale-team/whaleEcho/pkg/echoproto"
	"google.golang.org/protobuf/proto"
)

func AssertSendMessage(msg *echoproto.Message, room *echoproto.Room, command *echoproto.Command, conn *websocket.Conn) []byte {
	msg.Room = room
	msg.Sender = suite.users[0]
	msg.Type = echoproto.MessageType_Text
	msg.SentAt = time.Now().Unix()
	msgData, err := proto.Marshal(msg)
	assert.Nil(suite.T, err)
	command.Type = echoproto.CommandType_SendMessage
	command.Payload = msgData
	err = suite.SendCommand(conn, command)
	assert.Nil(suite.T, err)
	AssertRespStatus(conn, echoproto.Status_OK)
	return command.Payload
}

func AssertRespStatus(conn *websocket.Conn, status echoproto.Status) {
	resp, err := suite.ReadResp(conn)
	assert.Nil(suite.T, err)
	assert.Equal(suite.T, echoproto.Status_OK, resp.Status)
}

func AssertJoinRoom(conn *websocket.Conn, command *echoproto.Command, room *echoproto.Room) {
	roomData, err := proto.Marshal(room)
	assert.Nil(suite.T, err)
	command.Type = echoproto.CommandType_JoinRoom
	command.Payload = roomData
	err = suite.SendCommand(conn, command)
	assert.Nil(suite.T, err)
	AssertRespStatus(conn, echoproto.Status_OK)
}

var _ = Describe("Room Handler", func() {
	mockHandler := &MockHandler{handler: roomscenter.NewDefaultHandler(suite.broker)}
	suite.center.SetAsyncHandler(mockHandler)
	command := echoproto.Command{}

	AssertOpenRoom(suite.rooms[0], 0, mockHandler)
	AssertOpenRoom(suite.rooms[1], 0, mockHandler)
	AssertOpenRoom(suite.rooms[2], 0, mockHandler)

	Describe("#JoinRoom", func() {
		Context("when thers is a room number zero", func() {
			conn, err := suite.Dial()
			assert.Nil(suite.T, err)
			It("should response ok", func() {
				room := suite.rooms[1]
				room.Participant = suite.users[0]
				AssertJoinRoom(conn, &command, room)
				conn.Close()
			})
		})
	})

	Describe("#SendMessage", func() {
		Context("when there are two participants, using room number zero", func() {
			room := suite.rooms[0]
			conn, err := suite.Dial()
			assert.Nil(suite.T, err)
			conn2, err := suite.Dial()
			assert.Nil(suite.T, err)
			BeforeEach(func() {
				room.Participant = suite.users[1]
				AssertJoinRoom(conn, &command, room)
				room.Participant = suite.users[2]
				AssertJoinRoom(conn2, &command, room)
			})

			It("shoul receive message", func() {
				sentData := AssertSendMessage(suite.messages[1], room, &command, conn)

				_, recData, err := conn2.ReadMessage()
				assert.Nil(suite.T, err)
				assert.Equal(suite.T, sentData, recData)

				conn.Close()
				conn2.Close()
			})
		})
	})

	Describe("#LeaveRoom", func() {
		Context("when theare are two people, use room 2", func() {
			conn, err := suite.Dial()
			assert.Nil(suite.T, err)
			conn2, err := suite.Dial()
			assert.Nil(suite.T, err)
			room := suite.rooms[2]

			var sentData []byte
			BeforeEach(func() {
				room.Participant = suite.users[0]
				AssertJoinRoom(conn, &command, room)
				room.Participant = suite.users[1]
				AssertJoinRoom(conn2, &command, room)
			})

			It("should close the connection", func() {
				By("Sent message")

				command.Type = echoproto.CommandType_LeaveRoom
				room.Participant = suite.users[0]
				roomData, err := proto.Marshal(room)
				assert.Nil(suite.T, err)
				command.Payload = roomData
				err = suite.SendCommand(conn, &command)
				assert.Nil(suite.T, err)
				AssertRespStatus(conn, echoproto.Status_OK)
				sentData = AssertSendMessage(suite.messages[0], room, &command, conn2)

				By("Rec message")

				// conn2.SetReadDeadline(time.Now().Add(1 * time.Second))
				_, recData, err := conn2.ReadMessage()
				assert.Nil(suite.T, err)
				assert.Equal(suite.T, sentData, recData)

				By("Rec timeout")

				conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
				_, _, err = conn.ReadMessage()
				assert.NotNil(suite.T, err)
				conn.Close()
				conn2.Close()
			})

		})

	})
})
