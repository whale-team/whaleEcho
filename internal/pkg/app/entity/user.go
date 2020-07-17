package entity

import (
	"net"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

// User ...
type User struct {
	conn net.Conn
	Name string
	ID   string
}

func (u User) Receive(msg Message) error {
	return wsutil.WriteMessage(u.conn, ws.StateServerSide, ws.OpText, msg.Data())
}

func (u User) GetID() string {
	return u.ID
}
