package entity

import (
	"net"

	"github.com/gobwas/ws/wsutil"
)

// NewUser construct a user struct
func NewUser(conn net.Conn) *User {
	return &User{
		conn: conn,
	}
}

// User represent mapping between client connection and its identitiy
type User struct {
	conn net.Conn
	Name string
	ID   int64
}

// BindConn binding websocket connection
func (u *User) BindConn(conn net.Conn) {
	u.conn = conn
}

// Receive message
func (u User) Receive(msg MsgData) error {
	return wsutil.WriteServerBinary(u.conn, msg.GetData())
}

// GetID return user id
func (u User) GetID() int64 {
	return u.ID
}

// MsgData represent GetData interface
type MsgData interface {
	GetData() []byte
}
