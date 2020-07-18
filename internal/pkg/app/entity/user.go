package entity

import (
	"net"

	"github.com/gobwas/ws/wsutil"
)

func NewUser(conn net.Conn) *User {
	return &User{
		conn: conn,
	}
}

// User ...
type User struct {
	conn net.Conn
	Name string
	ID   int64
}

func (u *User) BindConn(conn net.Conn) {
	u.conn = conn
}

func (u User) Receive(msg *Message) error {
	return wsutil.WriteServerBinary(u.conn, msg.GetData())
}

func (u User) GetID() int64 {
	return u.ID
}
