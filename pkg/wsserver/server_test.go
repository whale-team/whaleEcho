package wsserver

import (
	"context"
	"testing"
	"time"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
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

	ctx := context.Background()
	conn, _, _, err := ws.DefaultDialer.Dial(ctx, "ws://127.0.0.1:1333")
	assert.Nil(t, err)

	msg := "hello!"
	wsutil.WriteMessage(conn, ws.StateClientSide, ws.OpText, []byte(msg))

	data, op, err := wsutil.ReadServerData(conn)
	assert.Nil(t, err)
	assert.Equal(t, op, ws.OpText)
	assert.Equal(t, data, []byte(msg))

	go func() {
		time.Sleep(1000 * time.Millisecond)
		err = conn.Close()
		assert.Nil(t, err)
	}()

	conn2, _, _, _ := ws.DefaultDialer.Dial(ctx, "ws://127.0.0.1:1333")
	conn2.Close()

	err = server.Shutdown()
	assert.Nil(t, err)
}

func TestShutdown(t *testing.T) {

}
