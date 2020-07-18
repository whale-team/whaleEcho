package wsserver

import (
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestEchoHandle(t *testing.T) {
	server := New()
	server.Handler = EchoHandle
	server.ErrHandler = ErrHandle
	server.ConnCloseHandler = ConnCloseHandle
	server.ConnBuildHandleFunc = ConnBuildHandle
	go server.ListenAndServe("127.0.0.1", "1333")

	time.Sleep(50 * time.Millisecond)

	conn, _, err := websocket.DefaultDialer.Dial("ws://127.0.0.1:1333", nil)
	assert.Nil(t, err)

	msg := "hello!"
	err = conn.WriteMessage(websocket.TextMessage, []byte(msg))

	assert.Nil(t, err)

	_, data, err := conn.ReadMessage()
	assert.Nil(t, err)

	assert.Equal(t, []byte(msg), data)

	go func() {
		time.Sleep(1000 * time.Millisecond)
		err = conn.Close()
		assert.Nil(t, err)
		t.Log("connection closed")
	}()

	conn2, _, _ := websocket.DefaultDialer.Dial("ws://127.0.0.1:1333", nil)
	conn2.Close()

	err = server.Shutdown()
	assert.Nil(t, err)
}

func TestShutdown(t *testing.T) {

}
