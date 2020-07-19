package handler_test

import (
	"context"
	"time"

	"github.com/nats-io/nats.go"
	. "github.com/onsi/ginkgo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/whale-team/whaleEcho/internal/pkg/app/entity"
	"github.com/whale-team/whaleEcho/internal/pkg/app/roomscenter"
	"github.com/whale-team/whaleEcho/pkg/echoproto"
	"github.com/whale-team/whaleEcho/pkg/subjects"
	"google.golang.org/protobuf/proto"
)

type MockHandler struct {
	handler roomscenter.AsyncHandler
	mock.Mock
}

func (h *MockHandler) OpenRoom(ctx context.Context, msg *nats.Msg, room *entity.Room) error {
	err := h.handler.OpenRoom(ctx, msg, room)
	h.Called(msg.Data)
	return err
}

func (h *MockHandler) CloseRoom(ctx context.Context, msg *nats.Msg, room *entity.Room) error {
	err := h.handler.CloseRoom(ctx, msg, room)
	h.Called(msg.Data)
	return err
}

func (h *MockHandler) ErrHandle(err error, msg *nats.Msg, room *entity.Room) {
	h.handler.ErrHandle(err, msg, room)
	h.Called()
}

var AssertOpenRoom = func(room *echoproto.Room, size int64, mockHandler *MockHandler) {
	data, err := proto.Marshal(room)
	assert.Nil(suite.T, err)

	if mockHandler != nil {
		mockHandler.On("OpenRoom", data)
	}
	err = suite.broker.PublishMessage(suite.Ctx, subjects.OpenRoomSubject, data)
	assert.Nil(suite.T, err)
	time.Sleep(10 * time.Millisecond)

	if mockHandler != nil {
		mockHandler.AssertCalled(suite.T, "OpenRoom", data)
		mockHandler.AssertNotCalled(suite.T, "ErrHandle")
	}

	if size != 0 {
		assert.GreaterOrEqual(suite.T, suite.center.Size(), size)
	}
}

var _ = Describe("Rooms Center", func() {
	mockHandler := &MockHandler{handler: roomscenter.NewDefaultHandler(suite.broker)}
	suite.center.SetAsyncHandler(mockHandler)

	Describe("Publish/Subscribe Spec", func() {
		Context("When Publish Create Room Message", func() {

			It("should create a room", func() {
				AssertOpenRoom(suite.rooms[0], 1, mockHandler)
			})
		})

		Context("When Publish Close Room Message", func() {

			BeforeEach(func() {
				AssertOpenRoom(suite.rooms[1], 0, mockHandler)
			})

			It("should close a room", func() {
				preSize := suite.center.Size()
				data, err := proto.Marshal(suite.rooms[1])
				assert.Nil(suite.T, err)
				mockHandler.On("CloseRoom", data)
				err = suite.broker.PublishMessage(suite.Ctx, subjects.CloseRoomSubject, data)
				assert.Nil(suite.T, err)
				time.Sleep(10 * time.Millisecond)
				mockHandler.AssertNumberOfCalls(suite.T, "CloseRoom", 1)
				mockHandler.AssertNotCalled(suite.T, "ErrHandle")
				assert.Equal(suite.T, preSize-1, suite.center.Size())
			})
		})
	})
})
