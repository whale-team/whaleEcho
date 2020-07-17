package entity_test

import (
	"time"
	"github.com/whale-team/whaleEcho/internal/pkg/app/entity"
	"github.com/whale-team/whaleEcho/internal/pkg/app/entity/roomcenter"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/mock"

	. "github.com/onsi/ginkgo"
	// . "github.com/onsi/gomega"
)

type MockReceiver struct {
	id string
	mock.Mock
}

func (r *MockReceiver) Receive(msg *entity.Message) error {
	args := r.Called(msg)
	return args.Error(0)
}

func (r *MockReceiver) GetID() string {
	return r.id
}

type MockSub struct {
	mock.Mock
}

func (m *MockSub) Unsubscribe() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockSub) Drain() error {
	return nil
}

func (m *MockSub) IsValid() error {
	return nil
}

var _ = Describe("Room Entity", func() {

	closedMsg := &entity.Message{Msg: &nats.Msg{Data: []byte("closed")}}

	Describe("#Run", func() {
		room := roomcenter.NewRoom()
		sub := &MockSub{}
		room.Subscribe = sub
		room.SetClosedMsg(closedMsg)

		Context("HappyCase", func() {
			ch := make(chan *nats.Msg, 1)
			room.SetMsgChannel(ch)
			room.Run()
			mock := &MockReceiver{id: "testing"}
			room.Join(mock)

			testMsg := &nats.Msg{Data: []byte("testing")}

			It("should Publish message to receiver", func() {
				mock.On("Receive", &entity.Message{Msg: testMsg}).Return(nil)
				n := 10000
				for i := 0; i < n; i++ {
					go func() { ch <- testMsg }()
				}

				time.Sleep(200 * time.Millisecond)
				mock.AssertExpectations(GinkgoT())
				mock.AssertNumberOfCalls(GinkgoT(), "Receive", n)
			})

			It("should close after notifing closed message", func() {
				mock.On("Receive", closedMsg).Return(nil)
				sub.On("Unsubscribe").Return(nil)
				room.Close()
				mock.AssertExpectations(GinkgoT())
				sub.AssertExpectations(GinkgoT())
			})
		})
	})

	Describe("#Leave", func() {
		room := roomcenter.NewRoom()
		sub := &MockSub{}
		room.Subscribe = sub
		testMsg := &nats.Msg{Data: []byte("testing")}
		room.SetClosedMsg(closedMsg)

		Context("HappyCase", func() {
			ch := make(chan *nats.Msg, 1)
			room.SetMsgChannel(ch)
			room.Run()
			mock := &MockReceiver{id: "testing"}
			room.Join(mock)
			It("should not receive message", func() {
				room.Leave(mock)
				sub.On("Unsubscribe").Return(nil)

				ch <- testMsg
				room.Close()
				mock.AssertNumberOfCalls(GinkgoT(), "Receive", 0)
				sub.AssertExpectations(GinkgoT())
			})
		})
	})

})