package entity

import (
	"github.com/vicxu416/wsserver"
)

// User represent mapping between client connection and its identitiy
type User struct {
	connCtx *wsserver.Context
	Name    string
	UID     string
	RoomUID string
}

// BindConn binding websocket connection
func (u *User) BindConn(connCtx *wsserver.Context) {
	u.connCtx = connCtx
}

func (u *User) CloseConn() error {
	return u.connCtx.Close()
}

// IsValid return false if connCtx is nil
func (u *User) IsValid() bool {
	return u.connCtx != nil
}

// Receive message
func (u User) Receive(msg []byte) error {
	return u.connCtx.WriteBinary(msg)
}
