package delivery_test

import (
	"fmt"
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
	time.Sleep(5 * time.Millisecond)
	assert.Contains(suite.T, suite.buf.String(), "receive message on subject(rooms.open")
}

func joinRoom(suite *testSuite, userData *echoproto.User, port string) *websocket.Conn {
	conn, err := dial(addr, port)
	assert.NoError(suite.T, err)
	command, err := newCommand(userData, echoproto.CommandType_JoinRoom)
	assert.NoError(suite.T, err)
	err = suite.SendCommand(conn, &command)
	assert.NoError(suite.T, err)
	time.Sleep(5 * time.Millisecond)
	msg, err := suite.ReadResp(conn)
	assert.NoError(suite.T, err)
	assert.Equal(suite.T, msg.Status, echoproto.Status_OK)
	assert.Contains(suite.T, suite.buf.String(), "join_room")
	return conn
}

var _ = Describe("Room Delivery Test", func() {
	port := "13333"
	suite, err := setupSuite(addr, port)
	if err != nil {
		fmt.Printf("setup suite failed, err:%+v", err)
		return
	}

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
				err = suite.repo.GetRoom(suite.ctx, suite.rooms[0].Uid, room)
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
				joinRoom(suite, userData, port)
				room := &entity.Room{}
				err = suite.repo.GetRoom(suite.ctx, userData.RoomUid, room)
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
				joinRoom(suite, userData, port)
				room := &entity.Room{}
				err = suite.repo.GetRoom(suite.ctx, userData.RoomUid, room)
				assert.NoError(suite.T, err)
				assert.GreaterOrEqual(suite.T, room.MembersCount, int64(1))
				room = suite.rms.GetRoom(userData.RoomUid)
				assert.NotNil(suite.T, room)
				assert.GreaterOrEqual(suite.T, room.CurrentMembersCount(), 1)
			})
		})
	})
})
