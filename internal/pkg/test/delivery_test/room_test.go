package delivery_test

import (
	"time"

	"github.com/gorilla/websocket"
	. "github.com/onsi/ginkgo"
	"github.com/stretchr/testify/assert"
	"github.com/whale-team/whaleEcho/internal/pkg/app/entity"
	"github.com/whale-team/whaleEcho/pkg/echoproto"
	"github.com/whale-team/whaleEcho/pkg/subjects"
	"google.golang.org/protobuf/proto"
)

func createRoom(suite *testSuite, roomData *echoproto.Room) {
	defer GinkgoRecover()

	data, err := proto.Marshal(roomData)
	assert.NoError(suite.T, err)
	err = suite.stan.Publish(suite.ctx, subjects.OpenRoomSubject, data)
	assert.NoError(suite.T, err)
	time.Sleep(10 * time.Millisecond)
	assert.Contains(suite.T, suite.buf.String(), "receive message on subject(rooms.open")
}

func joinRoom(suite *testSuite, userData *echoproto.User, status echoproto.Status) *websocket.Conn {
	conn, err := dial(suite.addr, suite.port)
	assert.NoError(suite.T, err)
	command, err := newCommand(userData, echoproto.CommandType_JoinRoom)
	assert.NoError(suite.T, err)
	err = suite.SendCommand(conn, &command)
	assert.NoError(suite.T, err)
	time.Sleep(10 * time.Millisecond)
	msg, err := suite.ReadResp(conn)
	assert.NoError(suite.T, err)
	assert.Equal(suite.T, msg.Status, status)
	assert.Contains(suite.T, suite.buf.String(), "join_room")
	return conn
}

func leaveRoom(suite *testSuite, conn *websocket.Conn, userData *echoproto.User, status echoproto.Status) {
	var err error
	if conn == nil {
		conn, err = dial(suite.addr, suite.port)
		assert.NoError(suite.T, err)
	}
	command, err := newCommand(userData, echoproto.CommandType_LeaveRoom)
	assert.NoError(suite.T, err)
	err = suite.SendCommand(conn, &command)
	assert.NoError(suite.T, err)
	time.Sleep(10 * time.Millisecond)
	msg, err := suite.ReadResp(conn)
	assert.NoError(suite.T, err)
	assert.Equal(suite.T, msg.Status, status)
	assert.Contains(suite.T, suite.buf.String(), "leave_room")
}

func closeRoom(suite *testSuite, room *echoproto.Room) {
	data, err := proto.Marshal(room)
	assert.NoError(suite.T, err)
	suite.stan.Publish(suite.ctx, subjects.CloseRoomSubject, data)
}

var _ = Describe("Room Delivery Test", func() {
	AfterEach(func() {
		suite.clear()
	})

	Describe("Listener#CreateRoom", func() {
		Context("when create a new room", func() {
			roomData := suite.rooms[0]
			It("should create a room in runtime store and redis db", func() {
				createRoom(suite, roomData)
				assert.GreaterOrEqual(suite.T, suite.rms.Len(), 1)
				room := &entity.Room{}
				err := suite.repo.GetRoom(suite.ctx, suite.rooms[0].Uid, room)
				assert.NoError(suite.T, err)
				assert.Equal(suite.T, room.UID, roomData.Uid)
				assert.Equal(suite.T, room.MembersLimit, roomData.MembersLimit)
			})
		})
	})

	Describe("WSHandler#JoinRoom", func() {
		Context("when join a created room", func() {
			roomData := suite.rooms[0]
			userData := suite.users[0]
			It("shuld increase room members count", func() {
				createRoom(suite, roomData)
				joinRoom(suite, userData, echoproto.Status_OK)
				room := &entity.Room{}
				err := suite.repo.GetRoom(suite.ctx, userData.RoomUid, room)
				assert.NoError(suite.T, err)
				assert.GreaterOrEqual(suite.T, room.MembersCount, int64(1))
				room = suite.rms.GetRoom(userData.RoomUid)
				assert.NotNil(suite.T, room)
				assert.GreaterOrEqual(suite.T, room.CurrentMembersCount(), 1)
			})
		})

		Context("when join a not running but created room", func() {
			userData := suite.users[1]
			room := &entity.Room{
				UID:          userData.RoomUid,
				MembersLimit: 10,
			}

			BeforeEach(func() {
				err := suite.repo.CreateRoom(suite.ctx, room)
				assert.NoError(suite.T, err)
			})

			It("should increase room members count", func() {
				joinRoom(suite, userData, echoproto.Status_OK)
				room := &entity.Room{}
				err := suite.repo.GetRoom(suite.ctx, userData.RoomUid, room)
				assert.NoError(suite.T, err)
				assert.GreaterOrEqual(suite.T, room.MembersCount, int64(1))
				room = suite.rms.GetRoom(userData.RoomUid)
				assert.NotNil(suite.T, room)
				assert.GreaterOrEqual(suite.T, room.CurrentMembersCount(), 1)
			})
		})

		Context("when join a room which is out of limit", func() {
			userData := suite.users[1]
			room := &entity.Room{
				UID:          userData.RoomUid,
				MembersLimit: 0,
			}

			BeforeEach(func() {
				suite.rms.CreateRoom(room)
				err := suite.repo.CreateRoom(suite.ctx, room)
				assert.NoError(suite.T, err)
			})

			It("should return not allow message", func() {
				joinRoom(suite, userData, echoproto.Status_NotAllow)
				room := &entity.Room{}
				err := suite.repo.GetRoom(suite.ctx, userData.RoomUid, room)
				assert.NoError(suite.T, err)
				assert.Equal(suite.T, room.MembersCount, int64(0))
				room = suite.rms.GetRoom(userData.RoomUid)
				assert.NotNil(suite.T, room)
				assert.Equal(suite.T, room.CurrentMembersCount(), 0)
			})
		})
	})

	Describe("WSHandler#LeaveRoom", func() {
		Context("when user levea room", func() {
			var (
				room = suite.rooms[0]
				user = suite.users[0]
				conn = &websocket.Conn{}
			)

			JustBeforeEach(func() {
				createRoom(suite, room)
				user.RoomUid = room.Uid
				conn = joinRoom(suite, user, echoproto.Status_OK)
			})

			It("should close connection", func() {
				r := suite.rms.GetRoom(room.Uid)
				assert.NotNil(suite.T, r)
				assert.Equal(suite.T, r.CurrentMembersCount(), 1)
				err := suite.repo.GetRoom(suite.ctx, room.Uid, r)
				assert.NoError(suite.T, err)
				assert.Equal(suite.T, r.MembersCount, int64(1))
				leaveRoom(suite, conn, user, echoproto.Status_OK)
				r = suite.rms.GetRoom(room.Uid)
				assert.NotNil(suite.T, r)
				assert.Equal(suite.T, r.CurrentMembersCount(), 0)
				err = suite.repo.GetRoom(suite.ctx, room.Uid, r)
				assert.NoError(suite.T, err)
				assert.Equal(suite.T, r.MembersCount, int64(0))
			})
		})
	})

	Describe("Listener#CloseRoom", func() {
		Context("when receive close room subject", func() {
			var (
				room = suite.rooms[0]
				user = suite.users[0]
				conn = &websocket.Conn{}
			)

			JustBeforeEach(func() {
				createRoom(suite, room)
				user.RoomUid = room.Uid
				conn = joinRoom(suite, user, echoproto.Status_OK)
			})

			It("should close room and publish close room message", func() {
				closeRoom(suite, room)
				msg, err := suite.ReadResp(conn)
				assert.NoError(suite.T, err)
				assert.NotNil(suite.T, msg)
				assert.Equal(suite.T, msg.RoomUid, room.Uid)
				assert.Contains(suite.T, msg.Text, "closed")
				assert.Contains(suite.T, msg.SenderName, "whale")
			})
		})
	})
})
