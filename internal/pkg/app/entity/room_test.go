package entity

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
)

type MockBinaryWriterCloser struct {
	mock.Mock
}

func (mock *MockBinaryWriterCloser) WriteBinary(data []byte) error {
	time.Sleep(1 * time.Millisecond)
	mock.Called(data)
	return nil
}

func (mock *MockBinaryWriterCloser) Close() error {
	mock.Called()
	return nil
}

func TestPublishMessage(t *testing.T) {
	writer := new(MockBinaryWriterCloser)
	room := &Room{
		UID: "abcs12332",
	}

	users := make([]*User, 10)
	for i := range users {
		users[i] = &User{
			Name:    fmt.Sprintf("user_%d", i),
			UID:     fmt.Sprintf("user_uid_%d", i),
			RoomUID: room.UID,
			conn:    writer,
		}
		room.JoinUser(users[i])
	}

	msg := []byte("testing")
	writer.On("WriteBinary", msg)
	room.publishMessage(msg)
	writer.AssertNumberOfCalls(t, "WriteBinary", 10)
}

func TestPublishMessageFaster(t *testing.T) {
	writer := new(MockBinaryWriterCloser)
	room := &Room{
		UID: "abcs12332",
	}

	users := make([]*User, 10)
	for i := range users {
		users[i] = &User{
			Name:    fmt.Sprintf("user_%d", i),
			UID:     fmt.Sprintf("user_uid_%d", i),
			RoomUID: room.UID,
			conn:    writer,
		}
		room.JoinUser(users[i])
	}

	msg := []byte("testing")
	writer.On("WriteBinary", msg)
	room.publishMessageFaster(msg)
	writer.AssertNumberOfCalls(t, "WriteBinary", 10)
}
func BenchmarkPublishMessage(b *testing.B) {
	writer := new(MockBinaryWriterCloser)
	room := &Room{
		UID: "abcs12332",
	}

	usersCount := 4
	users := make([]*User, usersCount)
	for i := range users {
		users[i] = &User{
			Name:    fmt.Sprintf("user_%d", i),
			UID:     fmt.Sprintf("user_uid_%d", i),
			RoomUID: room.UID,
			conn:    writer,
		}
		room.JoinUser(users[i])
	}

	msg := []byte("testing")
	writer.On("WriteBinary", msg)

	b.Run("normal", func(b *testing.B) {
		room.publishMessage(msg)
	})
	b.Run("faster", func(b *testing.B) {
		room.publishMessageFaster(msg)
	})
}

// func TestDivideSlice(t *testing.T) {

// 	test := []string{
// 		"1", "2", "3", "4", "5", "6", "7",
// 	}

// 	res := DivideSlice(test, 4)
// 	t.Log(res)

// }
