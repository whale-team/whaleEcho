package wsserver

import (
	"context"
	"encoding/json"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

type Context struct {
	ID        string
	opcode    ws.OpCode
	payload   []byte
	Conn      *onceCloseConn
	Context   context.Context
	ctxCancel context.CancelFunc
}

func (c *Context) read() error {
	data, opcode, err := wsutil.ReadClientData(c.Conn)
	if err != nil {
		return err
	}
	c.opcode = opcode
	c.payload = data
	return nil
}

func (c Context) BindJSON(dest interface{}) error {
	return json.Unmarshal(c.payload, dest)
}

func (c Context) GetPayload() []byte {
	return c.payload
}

func (c Context) Close() error {
	c.payload = nil
	c.opcode = ws.OpClose
	c.ctxCancel()

	return c.Conn.OnceClose()
}

func (c Context) WriteBinary(data []byte) error {
	return wsutil.WriteServerBinary(c.Conn, data)
}

func (c Context) WriteText(data []byte) error {
	return wsutil.WriteServerText(c.Conn, data)
}
